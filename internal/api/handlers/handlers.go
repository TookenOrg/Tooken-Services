package handlers

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	assetsService "github.com/TookenOrg/tooken-services/internal/assets_managements/services"
	authService "github.com/TookenOrg/tooken-services/internal/auth/services"
	blockchainService "github.com/TookenOrg/tooken-services/internal/blockchain/services"
	deliveryService "github.com/TookenOrg/tooken-services/internal/delivery/services"
	ordersService "github.com/TookenOrg/tooken-services/internal/orders/services"
	paymentsService "github.com/TookenOrg/tooken-services/internal/payments/services"
	usersService "github.com/TookenOrg/tooken-services/internal/users/services"
	"github.com/gin-gonic/gin"
)

// orderDeliverer is the delivery seen from the handler: one call, no answer.
// Declaring it here rather than taking the concrete service lets a test put a
// no-op in its place, so the handler suite can check the HTTP contract without
// a chain running behind it.
type orderDeliverer interface {
	Deliver(ctx context.Context, orderRef string)
}

type Handler struct {
	blockchainSvc *blockchainService.Service
	paymentsSvc   *paymentsService.Service
	authSvc       *authService.Service
	userSvc       *usersService.Service
	realEstateSvc *assetsService.Service
	orderSvc      *ordersService.Service
	deliverySvc   orderDeliverer
}

func NewHandler() server.ServerInterface {

	// Share the same blockchain service instance across all dependent services.
	blockchainSvc := blockchainService.NewService()

	return &Handler{
		blockchainSvc: blockchainSvc,
		paymentsSvc:   paymentsService.NewService(),
		authSvc:       authService.NewService(),
		realEstateSvc: assetsService.NewService(blockchainSvc),
		orderSvc:      ordersService.NewService(),
		userSvc:       usersService.NewService(),
		deliverySvc:   deliveryService.NewService(blockchainSvc),
	}
}

func (h *Handler) GetHealth(gCtx *gin.Context) {
	gCtx.JSON(200, gin.H{"status": "ok", "version": "v1.0.0"})
}
