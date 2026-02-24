package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	assetManagementDb "github.com/TookenOrg/tooken-services/internal/assets_managements/database"
	"github.com/TookenOrg/tooken-services/internal/orders/database"
	"github.com/TookenOrg/tooken-services/pkg/logger"
)

// Issuance order is for order on primary market
func (s *Service) CreateIssuanceOrder(ctx context.Context, req server.CreateIssuanceOrderRequest, userId int) (order server.IssuanceOrder, err error) {

	// 1 - check request
	if req.Quantity <= 0 {
		err = errors.New("Quantity can't be 0 or negative")
		return
	}

	realEstate, err := assetManagementDb.GetRealEstateById(ctx, req.RealEstateId)
	if err != nil {
		if err == sql.ErrNoRows {
			return server.IssuanceOrder{}, logger.LogError("Real Estate not found with id %d", req.RealEstateId)
		}
		return server.IssuanceOrder{}, err
	}

	// 2 - Insert in DB without ref
	idGenerated, createdAt, err := database.InsertIssuranceOrder(ctx, req.RealEstateId, req.Quantity, userId)

	// 3 - Generate a unique order reference
	orderReference := generateIssuanceOrderReference(int64(idGenerated), createdAt)

	// 4 - Update the record with reference
	err = database.UpdateReferenceIssuanceOrder(ctx, idGenerated, orderReference)

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

func generateIssuanceOrderReference(id int64, createdAt time.Time) string {
	return fmt.Sprintf(
		"ISS-%s-%06d",
		createdAt.Format("20060102"),
		id,
	)
}
