package handlers

import (
	"errors"
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/middleware"
	"github.com/TookenOrg/tooken-services/internal/users/services"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
)

// accountGoneMessage answers a valid token whose account was deleted since it
// was issued. It is a 401, not a 404: the caller did nothing wrong with the
// resource, its session points at nobody, and signing in again is the way out.
const accountGoneMessage = "Your account no longer exists. Sign in again."

func (h *Handler) GetCurrentUser(gCtx *gin.Context) {

	claims, exists := middleware.GetUserClaims(gCtx)
	if !exists {
		gCtx.JSON(http.StatusUnauthorized, server.APIResponse{Message: "Authentication required."})
		return
	}

	logger.LogInfo("🚀 Starting getting the current user")

	user, err := h.userSvc.GetUserByID(gCtx.Request.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			gCtx.JSON(http.StatusNotFound, server.APIResponse{Message: "User not found."})
			return
		}
		logger.LogError("Failed to retrieve the current user %d: %v", claims.UserID, err)
		gCtx.JSON(http.StatusInternalServerError, server.APIResponse{Message: "Unable to retrieve the current user."})
		return
	}
	gCtx.JSON(http.StatusOK, user)
}

// requireOwner lets the request through only when the token was issued to
// userID, and answers 401 or 403 otherwise.
//
// Favorites are personal: no role, administrator included, changes the
// favorites of somebody else. The route must declare bearerAuth in
// openapi.yaml, or there are no claims to compare the path to.
func requireOwner(gCtx *gin.Context, userID int) bool {
	claims, ok := middleware.GetUserClaims(gCtx)
	if !ok {
		gCtx.AbortWithStatusJSON(http.StatusUnauthorized, server.APIResponse{Message: "Authentication required."})
		return false
	}
	if claims.UserID != userID {
		gCtx.AbortWithStatusJSON(http.StatusForbidden, server.APIResponse{Message: "You can only change your own favorites."})
		return false
	}
	return true
}

// respondFavoriteError maps the favorite sentinels onto status codes. Anything
// unrecognised is a 500 whose detail is logged, never sent.
func respondFavoriteError(gCtx *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, services.ErrUserNotFound):
		gCtx.JSON(http.StatusUnauthorized, server.APIResponse{Message: accountGoneMessage})
	case errors.Is(err, services.ErrUnknownRealEstate):
		gCtx.JSON(http.StatusNotFound, server.APIResponse{Message: "Real estate not found."})
	case errors.Is(err, services.ErrFavoriteNotFound):
		gCtx.JSON(http.StatusNotFound, server.APIResponse{Message: "This real estate is not in your favorites."})
	case errors.Is(err, services.ErrFavoriteExists):
		gCtx.JSON(http.StatusConflict, server.APIResponse{Message: "This real estate already is in your favorites."})
	default:
		logger.LogError("%s: %v", fallback, err)
		gCtx.JSON(http.StatusInternalServerError, server.APIResponse{Message: fallback})
	}
}

func (h *Handler) AddFavoritesRealEstateForUser(gCtx *gin.Context, userId, realEstateId int) {

	if !requireOwner(gCtx, userId) {
		return
	}

	logger.LogInfo("🚀 Starting adding favorite real-estate for one user")

	err := h.userSvc.AddFavoritesRealEstateForUser(gCtx.Request.Context(), userId, realEstateId)
	if err != nil {
		respondFavoriteError(gCtx, err, "Unable to add the favorite.")
		return
	}

	gCtx.JSON(http.StatusOK, server.APIResponse{
		Message: "Favorite added successfully.",
	})
}

func (h *Handler) RemoveFavoritesRealEstateForUser(gCtx *gin.Context, userId, realEstateId int) {

	if !requireOwner(gCtx, userId) {
		return
	}

	logger.LogInfo("🚀 Starting removing favorite real-estate for one user")

	err := h.userSvc.RemoveFavoritesRealEstateForUser(gCtx.Request.Context(), userId, realEstateId)
	if err != nil {
		respondFavoriteError(gCtx, err, "Unable to remove the favorite.")
		return
	}

	gCtx.JSON(http.StatusOK, server.APIResponse{
		Message: "Favorite removed successfully.",
	})
}
