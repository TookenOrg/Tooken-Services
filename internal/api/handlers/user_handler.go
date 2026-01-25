package handlers

import (
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
)

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
