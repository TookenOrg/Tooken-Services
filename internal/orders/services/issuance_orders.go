package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/orders/database"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/avast/retry-go/v4"
)

var (
	ErrQuantityInvalid = errors.New("quantity can't be 0 or negative")

	// The refusals of the order placement are the database's own sentinels:
	// they are decided under the asset lock, in database.PlaceIssuanceOrder,
	// which cannot import this package. Aliasing keeps their detailed message
	// ("only 5 shares left, 10 requested") for the 409.
	ErrRealEstateNotFound     = database.ErrOrderAssetNotFound       // 404
	ErrRealEstateNotAvailable = database.ErrOrderAssetNotOpen        // 409
	ErrNotEnoughShares        = database.ErrOrderNotEnoughShares     // 409
	ErrStakeLimit             = database.ErrOrderStakeLimit          // 409
	ErrInvestorNotEligible    = database.ErrOrderInvestorNotEligible // 403
	ErrIdempotencyKeyReused   = database.ErrOrderIdempotencyKeyReuse // 422
	ErrOrderNotFound          = errors.New("order not found")
)

// defaultReservationTTL is how long an unpaid order holds its shares when
// ORDER_RESERVATION_TTL is unset or unusable (U1).
const defaultReservationTTL = 15 * time.Minute

// reservationTTL reads ORDER_RESERVATION_TTL as a Go duration ("15m", "1h").
// Like ETH_TX_WAIT, an unusable value falls back to the default and says so.
func reservationTTL() time.Duration {
	raw := strings.TrimSpace(os.Getenv("ORDER_RESERVATION_TTL"))
	if raw == "" {
		return defaultReservationTTL
	}

	ttl, err := time.ParseDuration(raw)
	if err != nil {
		logger.LogWarn("⚠️ ORDER_RESERVATION_TTL=%q is not a duration, using %s", raw, defaultReservationTTL)
		return defaultReservationTTL
	}

	// Zero or negative would break issuance_orders_awaiting_payment_ck
	// (reservation_expires_at > created_at): every order would fail with a 500.
	if ttl <= 0 {
		logger.LogWarn("⚠️ ORDER_RESERVATION_TTL=%q must be strictly positive, using %s", raw, defaultReservationTTL)
		return defaultReservationTTL
	}

	return ttl
}

// Issuance order is for order on primary market
func (s *Service) CreateIssuanceOrder(ctx context.Context, req server.CreateIssuanceOrderRequest, userId int, idempotencyKey server.IdempotencyKey) (order server.IssuanceOrder, created bool, err error) {

	// 1 - check request
	if req.Quantity <= 0 {
		err = ErrQuantityInvalid
		return
	}

	// 2 - Place the order, retrying the whole transaction on the rare
	// reference collision.
	var (
		orderReference string
		placed         database.PlacedOrderDTO
	)
	ttl := reservationTTL()
	err = retry.Do(
		func() error {
			var genErr error
			orderReference, genErr = generateIssuanceOrderReference(time.Now().UTC())
			if genErr != nil {
				return genErr
			}
			var placeErr error
			placed, placeErr = database.PlaceIssuanceOrder(ctx, database.PlaceOrderDTO{
				UserID:         userId,
				RealEstateID:   req.RealEstateId,
				Quantity:       req.Quantity,
				IdempotencyKey: idempotencyKey.String(),
				OrderReference: orderReference,
				ReservationTTL: ttl,
			})
			return placeErr
		},
		retry.Attempts(3),
		retry.Context(ctx),
		retry.RetryIf(func(err error) bool {
			return errors.Is(err, database.ErrDuplicateOrderReference)
		}),
		retry.OnRetry(func(n uint, err error) {
			logger.LogWarn("Order reference collision on %s, retrying (attempt %d)", orderReference, n+1)
		}),
		retry.LastErrorOnly(true),
	)
	if err != nil {
		return
	}
	created = !placed.Replayed

	// 3 - Answer with the order as the database holds it: amounts, fees,
	// dates from its clock. On a replay it is the existing order, as it is
	// now, not the reference generated for this call.
	order, err = database.GetIssuanceOrderByRef(ctx, placed.OrderReference, false, userId)
	if err != nil {
		// The order is committed: a retry with the same key replays it.
		err = fmt.Errorf("read placed order %s: %w", placed.OrderReference, err)
		return
	}
	order.CreatedBy = userId

	return
}

func (s *Service) ListIssuanceOrders(ctx context.Context, isStaff bool, userID int) (orders []server.IssuanceOrder, err error) {
	return database.ListIssuanceOrders(ctx, isStaff, userID)
}

func (s *Service) FetchIssuanceOrder(ctx context.Context, orderRef string, isStaff bool, userID int) (order server.IssuanceOrder, err error) {
	order, err = database.GetIssuanceOrderByRef(ctx, orderRef, isStaff, userID)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrOrderNotFound
		return
	}
	return
}

// generateIssuanceOrderReference builds "ISS-YYYYMMDD-XXXXXX": chronologically
// sortable and unpredictable (random suffix leaks no order volume).
func generateIssuanceOrderReference(t time.Time) (string, error) {
	// Crockford base32 alphabet without ambiguous chars (0/O, 1/I, U).
	const alphabet = "23456789ABCDEFGHJKLMNPQRSTVWXYZ"

	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return "", logger.LogError("failed to generate order reference: %v", err)
	}

	for i, b := range buf {
		buf[i] = alphabet[int(b)%len(alphabet)]
	}

	return fmt.Sprintf("ISS-%s-%s", t.Format("20060102"), string(buf)), nil
}
