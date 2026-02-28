package handlers

import (
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/middleware"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateIssuanceOrder(gCtx *gin.Context) {

	claims, exists := middleware.GetUserClaims(gCtx)
	if !exists {
		gCtx.JSON(http.StatusBadRequest, logger.LogError("JWT not valid"))
		return
	}

	var req server.CreateIssuanceOrderRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userId := claims.UserID

	logger.LogInfo("🚀 Starting creating one issuance order for userId %d and real estate id %d, quantity: %d", userId, req.RealEstateId, req.Quantity)

	orderCreated, err := h.orderSvc.CreateIssuanceOrder(gCtx.Request.Context(), req, userId)
	if err != nil {
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{
			Message: err.Error(),
		})
		return
	}

	resp := server.CreateIssuanceOrderResponse{
		Data:    &orderCreated,
		Message: "Order created",
	}

	gCtx.JSON(http.StatusCreated, resp)
}

func (h *Handler) FetchIssuanceOrder(gCtx *gin.Context, orderRef string) {
	logger.LogInfo("🚀 Starting fetch one issuance order for orderRef = [%s]", orderRef)

	order, err := h.orderSvc.FetchIssuanceOrder(gCtx.Request.Context(), orderRef)
	if err != nil {
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{
			Message: err.Error(),
		})
		return
	}

	resp := server.CreateIssuanceOrderResponse{
		Data:    &order,
		Message: "Order fetched",
	}

	gCtx.JSON(http.StatusOK, resp)
}
