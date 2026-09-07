package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/auth/database"
	"github.com/TookenOrg/tooken-services/internal/auth/utils"
)

// ErrEmailAlreadyUsed: an account already exists for that email.
var ErrEmailAlreadyUsed = errors.New("email already used")

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
		errCode = http.StatusInternalServerError
		return
	}

	// 3 - Insert new user
	userId, role, err := database.InsertUser(ctx, email, hashedPassword, fullName)
	if err != nil {
		errCode = http.StatusInternalServerError
		return
	}

	// 4 - generate jwt for session
	jwtToken, err = utils.GenerateJWT(userId, email, role)
	if err != nil {
		return
	}

	// 5 - return
	return
}

func checkIdempotencySignUp(ctx context.Context, email string) (isIdempotent bool, errCode int, err error) {
	exists, err := database.CheckEmailExists(ctx, email)
	if err != nil {
		// Nothing the caller can fix by editing the request.
		errCode = http.StatusInternalServerError
		return
	}
	if exists {
		errCode = http.StatusConflict
		err = ErrEmailAlreadyUsed
		return
	}

	isIdempotent = true
	return
}
