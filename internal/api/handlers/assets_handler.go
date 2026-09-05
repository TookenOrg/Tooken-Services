package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetActiveRealEstates(gCtx *gin.Context) {
	logger.LogDebug("🚀 Starting get real estates")

	realEstates, err := h.realEstateSvc.GetActiveRealEstates(gCtx.Request.Context())

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

	realEstate, err := h.realEstateSvc.GetActiveRealEstateById(gCtx.Request.Context(), id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			gCtx.JSON(http.StatusNotFound, server.APIResponse{Message: "Real estate not found."})
		} else {
			logger.LogError("Failed to retrieve real estate %d: %v", id, err)
			gCtx.JSON(http.StatusInternalServerError, server.APIResponse{
				Message: "Unable to retrieve real estate.",
			})
		}
		return
	}

	gCtx.JSON(http.StatusOK, realEstate)
}
