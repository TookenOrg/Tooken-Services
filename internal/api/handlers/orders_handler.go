package handlers

import (
	"errors"
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/internal/middleware"
	"github.com/TookenOrg/tooken-services/internal/orders/services"
	"github.com/TookenOrg/tooken-services/internal/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateIssuanceOrder(gCtx *gin.Context, params server.CreateIssuanceOrderParams) {

	if !middleware.RequireRole(gCtx, authUtils.RoleUser) {
		return
	}

	claims, exists := middleware.GetUserClaims(gCtx)
	if !exists {
		gCtx.JSON(http.StatusUnauthorized, server.APIResponse{Message: "JWT not valid"})
		return
	}
	userId := claims.UserID

	var req server.CreateIssuanceOrderRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{Message: err.Error()})
		return
	}

	logger.LogInfo("🚀 Starting creating one issuance order for userId %d and real estate id %d, quantity: %d. Idempotency key: %s", userId, req.RealEstateId, req.Quantity, params.IdempotencyKey)

	order, created, err := h.orderSvc.CreateIssuanceOrder(gCtx.Request.Context(), req, userId, params.IdempotencyKey)
	if err != nil {
		respondIssuanceOrderError(gCtx, err, "Unable to create issuance order.")
		return
	}

	if created {
		logger.LogInfo("✅ Successfully created issuance order for userId %d and real estate id %d, quantity: %d. Idempotency key: %s", userId, req.RealEstateId, req.Quantity, params.IdempotencyKey)
		resp := server.CreateIssuanceOrderResponse{
			Data:    &order,
			Message: "Order created",
		}
		gCtx.JSON(http.StatusCreated, resp)
	} else {
		logger.LogInfo("ℹ️ Issuance order already exists for userId %d and real estate id %d, quantity: %d. Idempotency key: %s", userId, req.RealEstateId, req.Quantity, params.IdempotencyKey)
		resp := server.CreateIssuanceOrderResponse{
			Data:    &order,
			Message: "Order already exists",
		}
		gCtx.JSON(http.StatusOK, resp)
	}
}

func (h *Handler) ListIssuanceOrders(gCtx *gin.Context) {
	logger.LogInfo("🚀 Starting list issuance orders for the authenticated user.")

	isStaff := utils.IsStaff(gCtx)
	logger.LogInfo("Fetching issuance orders for isStaff = %v", isStaff)

	claims, exists := middleware.GetUserClaims(gCtx)
	var userID int
	if exists {
		userID = claims.UserID
	}

	orders, err := h.orderSvc.ListIssuanceOrders(gCtx.Request.Context(), isStaff, userID)
	if err != nil {
		respondIssuanceOrderError(gCtx, err, "list issuance orders failed")
		return
	}

	resp := server.IssuanceOrderListResponse{
		Data:    orders,
		Message: "Orders fetched",
	}

	gCtx.JSON(http.StatusOK, resp)
}

func (h *Handler) FetchIssuanceOrder(gCtx *gin.Context, orderRef string) {
	logger.LogInfo("🚀 Starting fetch one issuance order for orderRef = [%s]", orderRef)

	isStaff := utils.IsStaff(gCtx)
	logger.LogInfo("Fetching issuance orders for isStaff = %v", isStaff)

	claims, exists := middleware.GetUserClaims(gCtx)
	var userID int
	if exists {
		userID = claims.UserID
	}

	order, err := h.orderSvc.FetchIssuanceOrder(gCtx.Request.Context(), orderRef, isStaff, userID)
	if err != nil {
		respondIssuanceOrderError(gCtx, err, "fetch issuance order failed")
		return
	}

	resp := server.CreateIssuanceOrderResponse{
		Data:    &order,
		Message: "Order fetched",
	}

	gCtx.JSON(http.StatusOK, resp)
}

// respondIssuanceOrderError maps the service sentinels onto status codes.
func respondIssuanceOrderError(gCtx *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, services.ErrRealEstateNotFound),
		errors.Is(err, services.ErrOrderNotFound):
		gCtx.JSON(http.StatusNotFound, server.APIResponse{Message: err.Error()})
	case errors.Is(err, services.ErrQuantityInvalid):
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{Message: err.Error()})
	case errors.Is(err, services.ErrNotEnoughShares),
		errors.Is(err, services.ErrRealEstateNotAvailable),
		errors.Is(err, services.ErrStakeLimit):
		gCtx.JSON(http.StatusConflict, server.APIResponse{Message: err.Error()})
	case errors.Is(err, services.ErrIdempotencyKeyReused):
		gCtx.JSON(http.StatusUnprocessableEntity, server.APIResponse{Message: err.Error()})
	case errors.Is(err, services.ErrInvestorNotEligible):
		gCtx.JSON(http.StatusForbidden, server.APIResponse{Message: err.Error()})
	default:
		logger.LogError("%s: %v", fallback, err)
		gCtx.JSON(http.StatusInternalServerError, server.APIResponse{Message: fallback})
	}
}
