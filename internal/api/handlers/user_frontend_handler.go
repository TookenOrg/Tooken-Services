package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/internal/middleware"
	"github.com/TookenOrg/tooken-services/internal/users/services"
	"github.com/TookenOrg/tooken-services/internal/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetFavoritesRealEstatesForUser(gCtx *gin.Context, userId int) {
	if !requireOwner(gCtx, userId) {
		return
	}
	favorites, err := h.userSvc.GetFavoritesRealEstatesForUser(gCtx.Request.Context(), userId, utils.IsStaff(gCtx))
	if err != nil {
		respondFavoriteError(gCtx, err, "Unable to retrieve favorites.")
		return
	}
	gCtx.JSON(http.StatusOK, favorites)
}

func (h *Handler) GetCurrentUserKycOnChain(gCtx *gin.Context) {
	claims, ok := middleware.GetUserClaims(gCtx)
	if !ok {
		gCtx.JSON(http.StatusUnauthorized, server.APIResponse{Message: "Authentication required."})
		return
	}
	status, err := h.userSvc.GetCurrentUserKycOnChain(gCtx.Request.Context(), claims.UserID)
	switch {
	case errors.Is(err, services.ErrUserNotFound):
		gCtx.JSON(http.StatusNotFound, server.APIResponse{Message: "User not found."})
	case errors.Is(err, services.ErrKYCOnChainUnavailable):
		logger.LogError("Unable to retrieve on-chain KYC status: %v", err)
		gCtx.JSON(http.StatusServiceUnavailable, server.APIResponse{Message: "On-chain KYC status is temporarily unavailable."})
	case err != nil:
		logger.LogError("Unable to retrieve on-chain KYC status: %v", err)
		gCtx.JSON(http.StatusInternalServerError, server.APIResponse{Message: "Unable to retrieve on-chain KYC status."})
	default:
		gCtx.JSON(http.StatusOK, status)
	}
}

func decodeUserRoleRequest(gCtx *gin.Context) (server.UpdateUserRoleRequest, error) {
	var body server.UpdateUserRoleRequest
	decoder := json.NewDecoder(http.MaxBytesReader(gCtx.Writer, gCtx.Request.Body, 1024*1024))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		return body, err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return body, errors.New("expected one JSON object")
	}
	return body, nil
}

func (h *Handler) UpdateUserRole(gCtx *gin.Context, userId int) {
	if !middleware.RequireRole(gCtx, authUtils.RoleAdmin) {
		return
	}
	body, err := decodeUserRoleRequest(gCtx)
	if err != nil {
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{Message: "A JSON object with role USER, MANAGER or ADMIN is required."})
		return
	}
	claims, _ := middleware.GetUserClaims(gCtx)
	user, err := h.userSvc.UpdateUserRole(gCtx.Request.Context(), claims.UserID, userId, server.UserRole(body.Role))
	switch {
	case errors.Is(err, services.ErrInvalidRole):
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{Message: err.Error()})
	case errors.Is(err, services.ErrRoleForbidden):
		gCtx.JSON(http.StatusForbidden, server.APIResponse{Message: "Your account must currently be an administrator."})
	case errors.Is(err, services.ErrRoleConflict):
		gCtx.JSON(http.StatusConflict, server.APIResponse{Message: "You cannot demote yourself or the last administrator."})
	case errors.Is(err, services.ErrUserNotFound):
		gCtx.JSON(http.StatusNotFound, server.APIResponse{Message: "User not found."})
	case err != nil:
		logger.LogError("Unable to update user role: %v", err)
		gCtx.JSON(http.StatusInternalServerError, server.APIResponse{Message: "Unable to update user role."})
	default:
		gCtx.JSON(http.StatusOK, user)
	}
}
