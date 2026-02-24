package handlers

import (
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetActiveRealEstates(gCtx *gin.Context) {
	logger.LogInfo("🚀 Starting get real estates")

	realEstates, err := h.realEstateSvc.GetActiveRealEstates(gCtx.Request.Context())

	if err != nil {
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{
			Message: err.Error(),
		})
		return
	}

	gCtx.JSON(http.StatusOK, realEstates)
}

func (h *Handler) GetRealEstateById(gCtx *gin.Context, id int) {
	logger.LogInfo("🚀 Starting get one real estate id %d", id)

	realEstate, err := h.realEstateSvc.GetRealEstateById(gCtx.Request.Context(), id)

	if err != nil {
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{
			Message: err.Error(),
		})
		return
	}

	gCtx.JSON(http.StatusOK, realEstate)
}
