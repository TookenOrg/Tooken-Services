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

// Issuer endpoints. Every one of them is back office: an issuer is the legal
// vehicle behind an asset, and the public catalogue only ever shows the trimmed
// view already embedded in a real estate.
//
// The role is checked here and the route protection in openapi.yaml: the spec
// decides that a token is required, the handler decides which roles are enough.

// respondIssuerError maps the service sentinels onto status codes. Anything
// unrecognised stays a 500: an unexpected failure must never be reported as the
// caller's mistake.
func respondIssuerError(gCtx *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, services.ErrIssuerNotFound):
		gCtx.JSON(http.StatusNotFound, server.APIResponse{Message: "Issuer not found."})
	case errors.Is(err, services.ErrInvalidIssuer):
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{Message: err.Error()})
	case errors.Is(err, services.ErrIssuerConflict):
		gCtx.JSON(http.StatusConflict, server.APIResponse{Message: err.Error()})
	default:
		logger.LogError("%s: %v", fallback, err)
		gCtx.JSON(http.StatusInternalServerError, server.APIResponse{Message: fallback})
	}
}

func (h *Handler) GetIssuers(gCtx *gin.Context) {
	logger.LogDebug("🚀 Starting get issuers")

	if !middleware.RequireRole(gCtx, authUtils.RoleManager, authUtils.RoleAdmin) {
		return
	}

	issuers, err := h.realEstateSvc.GetIssuers(gCtx.Request.Context())
	if err != nil {
		respondIssuerError(gCtx, err, "Unable to retrieve issuers.")
		return
	}

	gCtx.JSON(http.StatusOK, issuers)
}

func (h *Handler) GetIssuerById(gCtx *gin.Context, id int) {
	logger.LogDebug("🚀 Starting get issuer id %d", id)

	if !middleware.RequireRole(gCtx, authUtils.RoleManager, authUtils.RoleAdmin) {
		return
	}

	issuer, err := h.realEstateSvc.GetIssuerById(gCtx.Request.Context(), id)
	if err != nil {
		respondIssuerError(gCtx, err, "Unable to retrieve issuer.")
		return
	}

	gCtx.JSON(http.StatusOK, issuer)
}

func (h *Handler) CreateIssuer(gCtx *gin.Context) {
	logger.LogDebug("🚀 Starting create issuer")

	if !middleware.RequireRole(gCtx, authUtils.RoleManager, authUtils.RoleAdmin) {
		return
	}

	var req server.IssuerWriteRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{
			Message: "Invalid payload: " + err.Error(),
		})
		return
	}

	issuer, err := h.realEstateSvc.CreateIssuer(gCtx.Request.Context(), req)
	if err != nil {
		respondIssuerError(gCtx, err, "Unable to create issuer.")
		return
	}

	logger.LogInfo("✅ Issuer %d created", issuer.Id)
	gCtx.JSON(http.StatusCreated, issuer)
}

func (h *Handler) PatchIssuer(gCtx *gin.Context, id int) {
	logger.LogDebug("🚀 Starting patch issuer id %d", id)

	if !middleware.RequireRole(gCtx, authUtils.RoleManager, authUtils.RoleAdmin) {
		return
	}

	var req server.IssuerPatchRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{
			Message: "Invalid payload: " + err.Error(),
		})
		return
	}

	issuer, err := h.realEstateSvc.PatchIssuer(gCtx.Request.Context(), id, req)
	if err != nil {
		respondIssuerError(gCtx, err, "Unable to update issuer.")
		return
	}

	logger.LogInfo("✅ Issuer %d updated", id)
	gCtx.JSON(http.StatusOK, issuer)
}

func (h *Handler) DeleteIssuer(gCtx *gin.Context, id int) {
	logger.LogDebug("🚀 Starting dissolve issuer id %d", id)

	if !middleware.RequireRole(gCtx, authUtils.RoleManager, authUtils.RoleAdmin) {
		return
	}

	if err := h.realEstateSvc.DeleteIssuer(gCtx.Request.Context(), id); err != nil {
		respondIssuerError(gCtx, err, "Unable to dissolve issuer.")
		return
	}

	logger.LogInfo("✅ Issuer %d dissolved", id)
	gCtx.Status(http.StatusNoContent)
}
