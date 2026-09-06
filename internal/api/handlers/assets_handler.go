package handlers

import (
	"errors"
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/assets_managements/services"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/internal/middleware"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
)

// isStaff reports whether the caller manages the catalogue.
//
// The two real estate read endpoints are declared with an optional security
// block, so claims are present only when a token was supplied: no token means
// an anonymous visitor, and the public view.
func isStaff(gCtx *gin.Context) bool {
	claims, ok := middleware.GetUserClaims(gCtx)
	if !ok {
		return false
	}

	return claims.HasRole(authUtils.RoleManager, authUtils.RoleAdmin)
}

// respondRealEstateError maps the service sentinels onto status codes.
//
// Anything unrecognised is a 500: an unexpected failure, a dropped database
// connection for instance, must never be reported as the caller's mistake.
func respondRealEstateError(gCtx *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, services.ErrRealEstateNotFound):
		gCtx.JSON(http.StatusNotFound, server.APIResponse{Message: "Real estate not found."})
	case errors.Is(err, services.ErrInvalidRealEstate):
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{Message: err.Error()})
	case errors.Is(err, services.ErrRealEstateConflict):
		gCtx.JSON(http.StatusConflict, server.APIResponse{Message: err.Error()})
	default:
		logger.LogError("%s: %v", fallback, err)
		gCtx.JSON(http.StatusInternalServerError, server.APIResponse{Message: fallback})
	}
}

func (h *Handler) GetActiveRealEstates(gCtx *gin.Context) {
	logger.LogDebug("🚀 Starting get real estates")

	realEstates, err := h.realEstateSvc.GetRealEstates(gCtx.Request.Context(), isStaff(gCtx))

	if err != nil {
		logger.LogError("Failed to retrieve active real estates: %v", err)
		gCtx.JSON(http.StatusInternalServerError, server.APIResponse{
			Message: "Unable to retrieve real estates.",
		})
		return
	}

	gCtx.JSON(http.StatusOK, realEstates)
}

func (h *Handler) GetRealEstateById(gCtx *gin.Context, id int) {
	logger.LogDebug("🚀 Starting get one real estate id %d", id)

	realEstate, err := h.realEstateSvc.GetRealEstateById(gCtx.Request.Context(), id, isStaff(gCtx))

	if err != nil {
		respondRealEstateError(gCtx, err, "Unable to retrieve real estate.")
		return
	}

	gCtx.JSON(http.StatusOK, realEstate)
}

func (h *Handler) CreateRealEstate(gCtx *gin.Context) {
	logger.LogDebug("🚀 Starting create real estate")

	if !middleware.RequireRole(gCtx, authUtils.RoleManager, authUtils.RoleAdmin) {
		return
	}

	var req server.RealEstateWriteRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{
			Message: "Invalid payload: " + err.Error(),
		})
		return
	}

	realEstate, err := h.realEstateSvc.CreateRealEstate(gCtx.Request.Context(), req)
	if err != nil {
		respondRealEstateError(gCtx, err, "Unable to create real estate.")
		return
	}

	logger.LogInfo("✅ Real estate %d created", realEstate.Id)
	gCtx.JSON(http.StatusCreated, realEstate)
}

func (h *Handler) PatchRealEstate(gCtx *gin.Context, id int) {
	logger.LogDebug("🚀 Starting patch real estate id %d", id)

	if !middleware.RequireRole(gCtx, authUtils.RoleManager, authUtils.RoleAdmin) {
		return
	}

	var req server.RealEstatePatchRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{
			Message: "Invalid payload: " + err.Error(),
		})
		return
	}

	realEstate, err := h.realEstateSvc.PatchRealEstate(gCtx.Request.Context(), id, req)
	if err != nil {
		respondRealEstateError(gCtx, err, "Unable to update real estate.")
		return
	}

	logger.LogInfo("✅ Real estate %d updated", id)
	gCtx.JSON(http.StatusOK, realEstate)
}

func (h *Handler) DeleteRealEstate(gCtx *gin.Context, id int) {
	logger.LogDebug("🚀 Starting delete real estate id %d", id)

	if !middleware.RequireRole(gCtx, authUtils.RoleManager, authUtils.RoleAdmin) {
		return
	}

	if err := h.realEstateSvc.DeleteRealEstate(gCtx.Request.Context(), id); err != nil {
		respondRealEstateError(gCtx, err, "Unable to delete real estate.")
		return
	}

	logger.LogInfo("✅ Real estate %d deleted", id)
	gCtx.Status(http.StatusNoContent)
}
