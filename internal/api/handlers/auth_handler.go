package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	authServices "github.com/TookenOrg/tooken-services/internal/auth/services"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *Handler) PostAuthSignUp(gCtx *gin.Context) {

	logger.LogInfo("🚀 Starting sign-up")

	var req server.SignUpRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Binding does not enforce the required fields: a missing one arrives
	// empty. An empty email and password would make an account anyone can
	// sign in to.
	if req.Email == "" || req.Password == "" || strings.TrimSpace(req.FullName) == "" {
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{
			Message: "email, password and fullName are required.",
		})
		return
	}
	if len(req.Password) > authUtils.MaxPasswordBytes {
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{
			Message: fmt.Sprintf("The password must not exceed %d bytes.", authUtils.MaxPasswordBytes),
		})
		return
	}

	errCode, session, err := h.authSvc.SignUp(gCtx.Request.Context(), string(req.Email), req.Password, req.FullName)
	if err != nil || errCode >= http.StatusBadRequest {
		if errCode == 0 {
			errCode = http.StatusBadRequest
		}

		message := "Unable to create the account."
		switch {
		case errors.Is(err, authServices.ErrEmailAlreadyUsed):
			message = "An account already exists for this email."
		case errCode == http.StatusInternalServerError:
			// The database error is logged, never echoed: it would disclose the
			// schema to an anonymous caller.
			logger.LogError("Sign-up failed for %s: %v", req.Email, err)
		}

		gCtx.JSON(errCode, server.APIResponse{Message: message})
		return
	}

	gCtx.JSON(http.StatusCreated, server.SignUpResponse{
		Message: "User created successfully.",
		Data:    &session,
	})
}

func (h *Handler) PostAuthSignIn(gCtx *gin.Context) {

	logger.LogInfo("🚀 Starting sign-in")

	var req server.SignInRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	session, err := h.authSvc.SignIn(gCtx.Request.Context(), string(req.Email), req.Password)
	if err != nil {
		if errors.Is(err, authServices.ErrInvalidCredentials) {
			gCtx.JSON(http.StatusUnauthorized, server.APIResponse{
				Message: "Invalid email or password.",
			})
			return
		}

		logger.LogError("Sign-in failed for %s: %v", req.Email, err)
		gCtx.JSON(http.StatusInternalServerError, server.APIResponse{
			Message: "Unable to sign you in right now.",
		})
		return
	}

	gCtx.JSON(http.StatusOK, server.SignInResponse{
		Message: "User connected successfully.",
		Data:    &session,
	})
}

func (h *Handler) PostAuthRefresh(gCtx *gin.Context) {

	var req server.RefreshTokenRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{
			Message: "The body must be a JSON object holding the refreshToken.",
		})
		return
	}

	session, err := h.authSvc.Refresh(gCtx.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, authServices.ErrInvalidRefreshToken) {
			logger.LogInfo("Refresh refused: %v", err)
			gCtx.JSON(http.StatusUnauthorized, server.APIResponse{
				Message: "Session expired, please sign in again.",
			})
			return
		}

		logger.LogError("Refresh failed: %v", err)
		gCtx.JSON(http.StatusInternalServerError, server.APIResponse{
			Message: "Unable to refresh the session right now.",
		})
		return
	}

	gCtx.JSON(http.StatusOK, server.RefreshTokenResponse{
		Message: "Session refreshed.",
		Data:    &session,
	})
}
