package handlers

import (
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/api/server"
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
	if err != nil {
		if errCode != 0 {
			gCtx.JSON(errCode, err.Error())
		} else {

			gCtx.JSON(http.StatusBadRequest, err.Error())
		}
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
		gCtx.JSON(http.StatusBadRequest, server.APIResponse{
			Message: "Invalid email or password.",
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
