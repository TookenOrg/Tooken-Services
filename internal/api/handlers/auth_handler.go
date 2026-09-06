package handlers

import (
	"errors"
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	authServices "github.com/TookenOrg/tooken-services/internal/auth/services"
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

	errCode, jwtToken, err := h.authSvc.SignUp(gCtx.Request.Context(), string(req.Email), req.Password, req.FullName)
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

	gCtx.JSON(http.StatusCreated, server.APIResponse{
		Message: "User created successfully.",
		Data:    jwtToken,
	})
}

func (h *Handler) PostAuthSignIn(gCtx *gin.Context) {

	logger.LogInfo("🚀 Starting sign-in")

	var req server.SignInRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, jwtToken, refreshToken, tokenExpireAt, err := h.authSvc.SignIn(gCtx.Request.Context(), string(req.Email), req.Password)
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
		Data: &server.JwtToken{
			Token:        jwtToken,
			User:         user,
			RefreshToken: refreshToken,
			ExpiresAt:    tokenExpireAt,
		},
	})
}

func (h *Handler) GetUserById(gCtx *gin.Context, userId int64) {
	gCtx.JSON(http.StatusBadRequest, server.APIResponse{
		Message: "Not implemented.",
	})
}
