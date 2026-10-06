package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	assetManagementDb "github.com/TookenOrg/tooken-services/internal/assets_managements/database"
	"github.com/TookenOrg/tooken-services/internal/orders/database"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/avast/retry-go/v4"
	"github.com/shopspring/decimal"
)

var (
	ErrQuantityInvalid = errors.New("quantity can't be 0 or negative")
	// ErrRealEstateNotFound is returned when the requested real estate does not exist.
	ErrRealEstateNotFound = errors.New("real estate not found")

	ErrRealEstateNotAvailable = errors.New("real estate not available for issuance")

	ErrNotEnoughShares = errors.New("not enough shares available for issuance")

	ErrStakeLimit = errors.New("stake limit exceeded")

	ErrInvestorNotEligible = errors.New("investor not eligible")

	ErrIdempotencyKeyReused = errors.New("idempotency key reused")
)

const (
	StatusRealEstateFundraising = 4
)

// pricedOrder is what the investor owes, frozen at order time.
// Every amount is already rounded to 2 decimals, exactly as the database CHECKs expect.
type pricedOrder struct {
	UnitPrice   decimal.Decimal // copied as is from shares_config (8 decimals)
	GrossAmount decimal.Decimal // round(quantity × unitPrice, 2)
	FeeAmount   decimal.Decimal // sum of Fees[i].Amount; zero when there is no fee
	AmountDue   decimal.Decimal // GrossAmount + FeeAmount
	Fees        []pricedFee     // empty when there is no fee (U4), never nil
}

// pricedFee is one line of iss.issuance_order_fees.
type pricedFee struct {
	Code       string          // "ENTRY"
	Rate       decimal.Decimal // percentage, e.g. 2.0000 for 2 %
	BaseAmount decimal.Decimal // what the rate applies to: GrossAmount for ENTRY
	Amount     decimal.Decimal // round(BaseAmount × Rate / 100, 2)
}

// Issuance order is for order on primary market
func (s *Service) CreateIssuanceOrder(ctx context.Context, req server.CreateIssuanceOrderRequest, userId int, idempotencyKey server.IdempotencyKey) (order server.IssuanceOrder, created bool, err error) {

	// 0 - Check idempotency key

	// 1 - check request
	if req.Quantity <= 0 {
		err = ErrQuantityInvalid
		return
	}

	// false: an order can only be placed on a publicly available asset. A draft
	// or a deleted asset must be as unreachable here as it is on the listing.
	realEstate, err := assetManagementDb.GetRealEstateById(ctx, req.RealEstateId, false)
	if err != nil {
		if err == sql.ErrNoRows {
			return server.IssuanceOrder{}, false, ErrRealEstateNotFound
		}
		return server.IssuanceOrder{}, false, err
	}

	if realEstate.StatusId != StatusRealEstateFundraising {
		return server.IssuanceOrder{}, false, ErrRealEstateNotAvailable
	}

	pricedOrder := priceOrder(int64(req.Quantity), realEstate.SharesConfig.PricePerShare.Decimal, nil)

	// 2 - Insert with a unique reference, retrying on the rare reference collision
	var (
		orderReference string
		createdAt      time.Time
	)
	err = retry.Do(
		func() error {
			orderReference, err = generateIssuanceOrderReference(time.Now().UTC())
			if err != nil {
				return err
			}
			_, createdAt, err = database.InsertIssuranceOrder(ctx, req.RealEstateId, req.Quantity, userId, orderReference)
			return err
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

	order.CreatedAt = createdAt
	order.CreatedBy = userId
	order.OrderRef = orderReference
	order.RealEstateId = req.RealEstateId
	order.TokenQuantity = req.Quantity

	// TODO, structToString()
	logger.LogDebug("%v", realEstate)

	return

}

func (s *Service) FetchIssuanceOrder(ctx context.Context, orderRef string) (order server.IssuanceOrder, err error) {
	return database.GetIssuanceOrderByRef(ctx, orderRef)
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

// priceOrder freezes what the investor owes. Rounding must match the
// database CHECKs (round half away from zero, 2 decimals).
func priceOrder(quantity int64, unitPrice decimal.Decimal, entryFeeRate *decimal.Decimal) pricedOrder {

	pricedOrder := pricedOrder{
		UnitPrice: unitPrice,
		Fees:      []pricedFee{},
	}

	total := unitPrice.Mul(decimal.NewFromInt(quantity)).Round(2)
	pricedOrder.GrossAmount = total
	amountDue := total
	if entryFeeRate != nil && entryFeeRate.GreaterThan(decimal.Zero) {
		entryFee := total.Mul(*entryFeeRate).Div(decimal.NewFromInt(100)).Round(2)
		if entryFee.GreaterThan(decimal.Zero) {
			amountDue = total.Add(entryFee)
			pricedOrder.FeeAmount = entryFee
			pricedOrder.Fees = append(pricedOrder.Fees, pricedFee{
				Code:       "ENTRY",
				Rate:       *entryFeeRate,
				BaseAmount: total,
				Amount:     entryFee,
			})
		}
	}

	pricedOrder.AmountDue = amountDue

	return pricedOrder
}
