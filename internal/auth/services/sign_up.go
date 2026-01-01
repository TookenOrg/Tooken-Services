package services

import (
	"context"
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/auth/database"
	"github.com/TookenOrg/tooken-services/internal/auth/utils"
)

func (s *Service) SignUp(ctx context.Context, email, password, fullName string) (errCode int, jwtToken string, err error) {
	// 1 - idempotency
	isIdempotent, errCode, err := checkIdempotencySignUp(ctx, email)
	if err != nil {
		return
	}
	if !isIdempotent {
		return
	}

	// 2 - hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return
	}

	// 3 - Insert new user
	userId, err := database.InsertUser(ctx, email, hashedPassword, fullName)
	if err != nil {
		return
	}

	// 4 - generate jwt for session
	jwtToken, err = utils.GenerateJWT(userId, email)
	if err != nil {
		return
	}

	// 5 - return
	return
}

func checkIdempotencySignUp(ctx context.Context, email string) (isIdempotent bool, errCode int, err error) {
	exists, err := database.CheckEmailExists(ctx, email)
	if err != nil {
		errCode = http.StatusBadRequest
		return
	}
	if exists {
		errCode = http.StatusConflict
		return
	}

	isIdempotent = true
	return
}
