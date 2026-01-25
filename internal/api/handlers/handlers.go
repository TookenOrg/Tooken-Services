package handlers

import (
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	assetsService "github.com/TookenOrg/tooken-services/internal/assets_managements/services"
	authService "github.com/TookenOrg/tooken-services/internal/auth/services"
	blockchainService "github.com/TookenOrg/tooken-services/internal/blockchain/services"
	paymentsService "github.com/TookenOrg/tooken-services/internal/payments/services"
	usersService "github.com/TookenOrg/tooken-services/internal/users/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	blockchainSvc *blockchainService.Service
	paymentsSvc   *paymentsService.Service
	authSvc       *authService.Service
	userSvc       *usersService.Service
	realEstateSvc *assetsService.Service
}

func NewHandler() server.ServerInterface {
	return &Handler{
		blockchainSvc: blockchainService.NewService(),
		paymentsSvc:   paymentsService.NewService(),
		authSvc:       authService.NewService(),
		realEstateSvc: assetsService.NewService(),
		userSvc:       usersService.NewService(),
	}
}

func (h *Handler) GetHealth(gCtx *gin.Context) {
	gCtx.JSON(200, gin.H{"status": "ok", "version": "v1.0.0"})
}

func (h *Handler) PaymentAdd(gCtx *gin.Context) {
	gCtx.JSON(http.StatusAccepted, server.APIResponse{
		Message: "The deployment of all implementations has been started asynchronously.",
	})
}
