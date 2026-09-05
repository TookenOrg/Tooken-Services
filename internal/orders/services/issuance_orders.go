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
)

// Issuance order is for order on primary market
func (s *Service) CreateIssuanceOrder(ctx context.Context, req server.CreateIssuanceOrderRequest, userId int) (order server.IssuanceOrder, err error) {

	// 1 - check request
	if req.Quantity <= 0 {
		err = errors.New("Quantity can't be 0 or negative")
		return
	}

	realEstate, err := assetManagementDb.GetActiveRealEstateById(ctx, req.RealEstateId)
	if err != nil {
		if err == sql.ErrNoRows {
			return server.IssuanceOrder{}, logger.LogError("Real Estate not found with id %d", req.RealEstateId)
		}
		return
	}

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
