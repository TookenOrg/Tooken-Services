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

// refuseKyc answers a refusal the caller can act on. The KYC sentinels carry
// no internal detail, so their text is safe to show.
func refuseKyc(gCtx *gin.Context, code int, err error) {
	gCtx.JSON(code, server.APIResponse{Message: err.Error()})
}

// failKyc answers an error nobody expected. The detail goes to the log only:
// a database error names tables and connection state a client has no use for.
func failKyc(gCtx *gin.Context, action string, err error) {
	logger.LogError("Unable to %s: %v", action, err)
	gCtx.JSON(http.StatusInternalServerError, server.APIResponse{Message: "Unable to " + action + "."})
}

func invalidKycPayload(gCtx *gin.Context, err error) {
	gCtx.JSON(http.StatusBadRequest, server.APIResponse{Message: "Invalid payload: " + err.Error()})
}

func (h *Handler) PostKycVerifications(gCtx *gin.Context) {

	claims, exists := middleware.GetUserClaims(gCtx)
	if !exists {
		gCtx.JSON(http.StatusUnauthorized, server.APIResponse{Message: "Authentication required."})
		return
	}

	userId := claims.UserID
	logger.LogInfo("🚀 Starting creating a new KYC verification for user ID:%d", userId)

	var req server.KycVerificationRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		invalidKycPayload(gCtx, err)
		return
	}

	err := h.userSvc.PostKycVerifications(gCtx.Request.Context(), userId, &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			gCtx.JSON(http.StatusUnauthorized, server.APIResponse{Message: accountGoneMessage})
		case errors.Is(err, services.ErrMessageKycAlreadyExists), errors.Is(err, services.ErrMessageKycStillValid):
			refuseKyc(gCtx, http.StatusConflict, err)
		case errors.Is(err, services.ErrMessageInvalidCountryCode), errors.Is(err, services.ErrMessageInvalidFullName):
			refuseKyc(gCtx, http.StatusBadRequest, err)
		default:
			failKyc(gCtx, "submit the KYC verification", err)
		}
		return
	}

	gCtx.JSON(http.StatusCreated, server.APIResponse{
		Message: "KYC verification submitted successfully.",
	})
}

func (h *Handler) GetListKycVerifications(gCtx *gin.Context, params server.GetListKycVerificationsParams) {

	if !middleware.RequireRole(gCtx, authUtils.RoleAdmin) {
		return
	}

	logger.LogInfo("🚀 Starting getting the list of KYC verifications")

	var statusStr *string
	if params.Status != nil {
		// The generated wrapper does not enforce the enum: without this check an
		// unknown status matches nothing and reads as "nobody is in that state".
		switch *params.Status {
		case server.GetListKycVerificationsParamsStatusSubmitted, server.GetListKycVerificationsParamsStatusApproved:
		default:
			gCtx.JSON(http.StatusBadRequest, server.APIResponse{Message: "Invalid status: use submitted or approved."})
			return
		}
		statusStr = lo.ToPtr(string(*params.Status))
	}

	kycVerification, err := h.userSvc.GetListKycVerifications(gCtx.Request.Context(), params.UserId, statusStr)
	if err != nil {
		failKyc(gCtx, "list the KYC verifications", err)
		return
	}

	gCtx.JSON(http.StatusOK, kycVerification)
}

func (h *Handler) GetKycVerificationById(gCtx *gin.Context, verificationId int) {

	if !middleware.RequireRole(gCtx, authUtils.RoleAdmin) {
		return
	}

	logger.LogInfo("🚀 Starting getting KYC verification by ID: %d", verificationId)

	kycVerification, err := h.userSvc.GetKycVerificationById(gCtx.Request.Context(), verificationId)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrKycVerificationNotFound):
			refuseKyc(gCtx, http.StatusNotFound, err)
		default:
			failKyc(gCtx, "retrieve the KYC verification", err)
		}
		return
	}

	gCtx.JSON(http.StatusOK, kycVerification)
}

func (h *Handler) ApproveKycVerification(gCtx *gin.Context, verificationId int) {

	if !middleware.RequireRole(gCtx, authUtils.RoleAdmin) {
		return
	}
	claims, _ := middleware.GetUserClaims(gCtx)

	logger.LogInfo("🚀 Starting approving KYC verification by ID: %d", verificationId)

	var req server.KycVerificationApproveRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		invalidKycPayload(gCtx, err)
		return
	}
	expiresAt := req.ExpiresAt

	kycVerification, err := h.userSvc.ApproveKycVerification(gCtx.Request.Context(), verificationId, expiresAt, claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrKycVerificationNotFound):
			refuseKyc(gCtx, http.StatusNotFound, err)
		case errors.Is(err, services.ErrMessageInvalidExpiresAt):
			refuseKyc(gCtx, http.StatusBadRequest, err)
		case errors.Is(err, services.ErrMessageInvalidStatus):
			refuseKyc(gCtx, http.StatusConflict, err)
		default:
			failKyc(gCtx, "approve the KYC verification", err)
		}
		return
	}

	gCtx.JSON(http.StatusOK, kycVerification)
}

func (h *Handler) RejectKycVerification(gCtx *gin.Context, verificationId int) {

	if !middleware.RequireRole(gCtx, authUtils.RoleAdmin) {
		return
	}
	claims, _ := middleware.GetUserClaims(gCtx)

	logger.LogInfo("🚀 Starting rejecting KYC verification by ID: %d", verificationId)

	var req server.KycVerificationRejectRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		invalidKycPayload(gCtx, err)
		return
	}
	reason := req.Reason

	kycVerification, err := h.userSvc.RejectKycVerification(gCtx.Request.Context(), verificationId, reason, claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrKycVerificationNotFound):
			refuseKyc(gCtx, http.StatusNotFound, err)
		case errors.Is(err, services.ErrMessageInvalidReason):
			refuseKyc(gCtx, http.StatusBadRequest, err)
		case errors.Is(err, services.ErrMessageInvalidStatus):
			refuseKyc(gCtx, http.StatusConflict, err)
		default:
			failKyc(gCtx, "reject the KYC verification", err)
		}
		return
	}

	gCtx.JSON(http.StatusOK, kycVerification)
}

func (h *Handler) RevokeKycVerification(gCtx *gin.Context, verificationId int) {

	if !middleware.RequireRole(gCtx, authUtils.RoleAdmin) {
		return
	}
	claims, _ := middleware.GetUserClaims(gCtx)

	logger.LogInfo("🚀 Starting revoking KYC verification by ID: %d", verificationId)

	var req server.KycVerificationRevokeRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		invalidKycPayload(gCtx, err)
		return
	}
	reason := req.Reason

	kycVerification, err := h.userSvc.RevokeKycVerification(gCtx.Request.Context(), verificationId, reason, claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrKycVerificationNotFound):
			refuseKyc(gCtx, http.StatusNotFound, err)
		case errors.Is(err, services.ErrMessageInvalidReason):
			refuseKyc(gCtx, http.StatusBadRequest, err)
		case errors.Is(err, services.ErrMessageInvalidStatus):
			refuseKyc(gCtx, http.StatusConflict, err)
		default:
			failKyc(gCtx, "revoke the KYC verification", err)
		}
		return
	}

	gCtx.JSON(http.StatusOK, kycVerification)
}

func (h *Handler) SyncApprovedKycVerificationOnChain(gCtx *gin.Context, verificationId int) {

	if !middleware.RequireRole(gCtx, authUtils.RoleAdmin) {
		return
	}

	logger.LogInfo("🚀 Starting syncing approved KYC verification on-chain by ID: %d", verificationId)

	kycVerification, err := h.userSvc.SyncApprovedKycVerificationOnChain(gCtx.Request.Context(), verificationId)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrKycVerificationNotFound):
			refuseKyc(gCtx, http.StatusNotFound, err)
		// The filing was valid when approved; what blocks it now is the state
		// it is in, hence a conflict rather than a bad request.
		case errors.Is(err, services.ErrMessageInvalidStatus),
			errors.Is(err, services.ErrMessageInvalidExpiresAt),
			errors.Is(err, services.ErrMessageInvalidCountryCode),
			errors.Is(err, services.ErrMessageIncoherentUserState),
			errors.Is(err, services.ErrMessageInvalidWalletAddress),
			errors.Is(err, services.ErrMessageFailedToVerifyOnChain):
			refuseKyc(gCtx, http.StatusConflict, err)
		default:
			failKyc(gCtx, "sync the KYC verification on-chain", err)
		}
		return
	}

	gCtx.JSON(http.StatusOK, kycVerification)
}
