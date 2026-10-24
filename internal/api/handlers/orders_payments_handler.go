package handlers

import (
	"context"
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/internal/middleware"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ConfirmIssuanceOrderPayment(gCtx *gin.Context, orderRef string) {

	if !middleware.RequireRole(gCtx, authUtils.RoleAdmin) {
		return
	}

	claims, exists := middleware.GetUserClaims(gCtx)
	if !exists {
		gCtx.JSON(http.StatusUnauthorized, server.APIResponse{Message: "JWT not valid"})
		return
	}
	userID := claims.UserID

	logger.LogInfo("🚀 Starting the confirmation of issuance order payment for orderRef: %s", orderRef)

	order, err := h.paymentsSvc.ConfirmIssuanceOrderPayment(gCtx.Request.Context(), orderRef, userID)
	if err != nil {
		respondIssuanceOrderError(gCtx, err, "Failed to confirm issuance order payment")
		return
	}
	logger.LogInfo("✅ Successfully confirmed issuance order payment for orderRef: %s", orderRef)

	logger.LogInfo("Starting async delivery for orderRef %s ....", orderRef)
	go h.deliverySvc.Deliver(context.WithoutCancel(gCtx.Request.Context()), orderRef)

	gCtx.JSON(http.StatusOK, server.CreateIssuanceOrderResponse{
		Data:    &order,
		Message: "Successfully confirmed issuance order payment for orderRef: " + orderRef,
	})
}
