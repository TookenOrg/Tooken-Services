package handlers

import (
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetRealEstates(gCtx *gin.Context) {
	logger.LogInfo("🚀 Starting get real estates")

	realEstates, err := h.realEstateSvc.GetRealEstates(gCtx.Request.Context())

	if err != nil {
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{
			Message: err.Error(),
		})
		return
	}

	gCtx.JSON(http.StatusOK, realEstates)
}
