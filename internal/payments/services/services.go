package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	orderDB "github.com/TookenOrg/tooken-services/internal/orders/database"
	"github.com/TookenOrg/tooken-services/internal/payments/database"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

var (
	ErrOrderNotAwaitingPayment = database.ErrOrderNotAwaitingPayment // 409
	ErrOrderReservationExpired = database.ErrOrderReservationExpired // 409
)

func (s *Service) ConfirmIssuanceOrderPayment(ctx context.Context, orderReference string, userID int) (order server.IssuanceOrder, err error) {

	// 1 - Get issuance order
	order, err = orderDB.GetIssuanceOrderByRef(ctx, orderReference, true, 0)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = orderDB.ErrOrderNotFound
		}
		return
	}

	// 2 - Confirm the payment in the database
	err = database.ConfirmIssuanceOrderPayment(ctx, order.RealEstateId, orderReference, userID)
	if err != nil {
		return
	}

	// 3 - Get updated issuance order
	order, err = orderDB.GetIssuanceOrderByRef(ctx, orderReference, true, 0)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = orderDB.ErrOrderNotFound
		}
		return
	}

	return
}
