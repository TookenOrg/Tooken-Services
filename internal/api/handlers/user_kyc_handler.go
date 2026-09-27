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

func (h *Handler) GetKycVerificationById(gCtx *gin.Context, verificationId int) {

	claims, exists := middleware.GetUserClaims(gCtx)
	if !exists {
		gCtx.JSON(http.StatusUnauthorized, logger.LogError("JWT not valid"))
		return
	}

	if claims.Role != authUtils.RoleAdmin {
		gCtx.JSON(http.StatusForbidden, logger.LogError("Forbidden"))
		return
	}

	logger.LogInfo("🚀 Starting getting KYC verification by ID: %d", verificationId)

	kycVerification, err := h.userSvc.GetKycVerificationById(gCtx.Request.Context(), verificationId)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrKycVerificationNotFound):
			gCtx.JSON(http.StatusNotFound, err.Error())
			return
		}
		gCtx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	gCtx.JSON(http.StatusOK, kycVerification)
}

func (h *Handler) ApproveKycVerification(gCtx *gin.Context, verificationId int) {

	claims, exists := middleware.GetUserClaims(gCtx)
	if !exists {
		gCtx.JSON(http.StatusUnauthorized, logger.LogError("JWT not valid"))
		return
	}

	if claims.Role != authUtils.RoleAdmin {
		gCtx.JSON(http.StatusForbidden, logger.LogError("Forbidden"))
		return
	}

	logger.LogInfo("🚀 Starting approving KYC verification by ID: %d", verificationId)

	var req server.KycVerificationApproveRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	expiresAt := req.ExpiresAt

	kycVerification, err := h.userSvc.ApproveKycVerification(gCtx.Request.Context(), verificationId, expiresAt, claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrKycVerificationNotFound):
			gCtx.JSON(http.StatusNotFound, err.Error())
			return
		case errors.Is(err, services.ErrMessageInvalidExpiresAt):
			gCtx.JSON(http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, services.ErrMessageInvalidStatus):
			gCtx.JSON(http.StatusConflict, err.Error())
			return
		}
		gCtx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	gCtx.JSON(http.StatusOK, kycVerification)
}

func (h *Handler) RejectKycVerification(gCtx *gin.Context, verificationId int) {

	claims, exists := middleware.GetUserClaims(gCtx)
	if !exists {
		gCtx.JSON(http.StatusUnauthorized, logger.LogError("JWT not valid"))
		return
	}

	if claims.Role != authUtils.RoleAdmin {
		gCtx.JSON(http.StatusForbidden, logger.LogError("Forbidden"))
		return
	}

	logger.LogInfo("🚀 Starting rejecting KYC verification by ID: %d", verificationId)

	var req server.KycVerificationRejectRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	reason := req.Reason

	kycVerification, err := h.userSvc.RejectKycVerification(gCtx.Request.Context(), verificationId, reason, claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrKycVerificationNotFound):
			gCtx.JSON(http.StatusNotFound, err.Error())
			return
		case errors.Is(err, services.ErrMessageInvalidReason):
			gCtx.JSON(http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, services.ErrMessageInvalidStatus):
			gCtx.JSON(http.StatusConflict, err.Error())
			return
		}
		gCtx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	gCtx.JSON(http.StatusOK, kycVerification)
}

func (h *Handler) RevokeKycVerification(gCtx *gin.Context, verificationId int) {

	claims, exists := middleware.GetUserClaims(gCtx)
	if !exists {
		gCtx.JSON(http.StatusUnauthorized, logger.LogError("JWT not valid"))
		return
	}

	if claims.Role != authUtils.RoleAdmin {
		gCtx.JSON(http.StatusForbidden, logger.LogError("Forbidden"))
		return
	}

	logger.LogInfo("🚀 Starting revoking KYC verification by ID: %d", verificationId)

	var req server.KycVerificationRevokeRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	reason := req.Reason

	kycVerification, err := h.userSvc.RevokeKycVerification(gCtx.Request.Context(), verificationId, reason, claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrKycVerificationNotFound):
			gCtx.JSON(http.StatusNotFound, err.Error())
			return
		case errors.Is(err, services.ErrMessageInvalidReason):
			gCtx.JSON(http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, services.ErrMessageInvalidStatus):
			gCtx.JSON(http.StatusConflict, err.Error())
			return
		}
		gCtx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	gCtx.JSON(http.StatusOK, kycVerification)
}

func (h *Handler) SyncApprovedKycVerificationOnChain(gCtx *gin.Context, verificationId int) {
	claims, exists := middleware.GetUserClaims(gCtx)
	if !exists {
		gCtx.JSON(http.StatusUnauthorized, logger.LogError("JWT not valid"))
		return
	}

	if claims.Role != authUtils.RoleAdmin {
		gCtx.JSON(http.StatusForbidden, logger.LogError("Forbidden"))
		return
	}

	logger.LogInfo("🚀 Starting syncing approved KYC verification on-chain by ID: %d", verificationId)

	kycVerification, err := h.userSvc.SyncApprovedKycVerificationOnChain(gCtx.Request.Context(), verificationId)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrKycVerificationNotFound):
			gCtx.JSON(http.StatusNotFound, err.Error())
			return
		case errors.Is(err, services.ErrMessageInvalidStatus):
			gCtx.JSON(http.StatusConflict, err.Error())
			return
		case errors.Is(err, services.ErrMessageInvalidExpiresAt):
			gCtx.JSON(http.StatusConflict, err.Error())
			return
		case errors.Is(err, services.ErrMessageInvalidCountryCode):
			gCtx.JSON(http.StatusConflict, err.Error())
			return
		case errors.Is(err, services.ErrMessageIncoherentUserState):
			gCtx.JSON(http.StatusConflict, err.Error())
			return
		case errors.Is(err, services.ErrMessageInvalidWalletAddress):
			gCtx.JSON(http.StatusConflict, err.Error())
			return
		}
		gCtx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	gCtx.JSON(http.StatusOK, kycVerification)
}
