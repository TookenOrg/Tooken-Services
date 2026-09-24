package handlers

import (
	"errors"
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/internal/middleware"
	"github.com/TookenOrg/tooken-services/internal/users/services"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func (h *Handler) GetCurrentUser(gCtx *gin.Context) {

	claims, exists := middleware.GetUserClaims(gCtx)
	if !exists {
		gCtx.JSON(http.StatusUnauthorized, logger.LogError("JWT not valid"))
		return
	}

	logger.LogInfo("🚀 Starting getting the current user")

	user, err := h.userSvc.GetUserByID(gCtx.Request.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			gCtx.JSON(http.StatusNotFound, err.Error())
			return
		}
		gCtx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	gCtx.JSON(http.StatusOK, user)
}

func (h *Handler) AddFavoritesRealEstateForUser(gCtx *gin.Context, userId, realEstateId int) {

	logger.LogInfo("🚀 Starting adding favorite real-estate for one user")

	err := h.userSvc.AddFavoritesRealEstateForUser(gCtx.Request.Context(), userId, realEstateId)
	if err != nil {
		gCtx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	gCtx.JSON(http.StatusOK, server.APIResponse{
		Message: "Favorite added successfully.",
	})
}

func (h *Handler) RemoveFavoritesRealEstateForUser(gCtx *gin.Context, userId, realEstateId int) {

	logger.LogInfo("🚀 Starting removing favorite real-estate for one user")

	err := h.userSvc.RemoveFavoritesRealEstateForUser(gCtx.Request.Context(), userId, realEstateId)
	if err != nil {
		gCtx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	gCtx.JSON(http.StatusOK, server.APIResponse{
		Message: "Favorite removed successfully.",
	})
}

func (h *Handler) PostKycVerifications(gCtx *gin.Context) {

	claims, exists := middleware.GetUserClaims(gCtx)
	if !exists {
		gCtx.JSON(http.StatusUnauthorized, logger.LogError("JWT not valid"))
		return
	}

	userId := claims.UserID
	logger.LogInfo("🚀 Starting creating a new KYC verification for user ID:%d", userId)

	var req server.KycVerificationRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := h.userSvc.PostKycVerifications(gCtx.Request.Context(), userId, &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrMessageKycAlreadyExists):
			gCtx.JSON(http.StatusConflict, err.Error())
			return
		case errors.Is(err, services.ErrMessageInvalidCountryCode), errors.Is(err, services.ErrMessageInvalidFullName):
			gCtx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		gCtx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	gCtx.JSON(http.StatusCreated, server.APIResponse{
		Message: "KYC verification submitted successfully.",
	})
}

func (h *Handler) GetListKycVerifications(gCtx *gin.Context, params server.GetListKycVerificationsParams) {

	claims, exists := middleware.GetUserClaims(gCtx)
	if !exists {
		gCtx.JSON(http.StatusUnauthorized, logger.LogError("JWT not valid"))
		return
	}

	if claims.Role != authUtils.RoleAdmin {
		gCtx.JSON(http.StatusForbidden, logger.LogError("Forbidden"))
		return
	}

	logger.LogInfo("🚀 Starting getting the list of KYC verifications")

	var statusStr *string
	if params.Status != nil {
		statusStr = lo.ToPtr(string(*params.Status))
	}

	kycVerification, err := h.userSvc.GetListKycVerifications(gCtx.Request.Context(), params.UserId, statusStr)
	if err != nil {
		gCtx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	gCtx.JSON(http.StatusOK, kycVerification)
}
