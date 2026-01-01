package services

import (
	"context"
	"errors"

	"github.com/TookenOrg/tooken-services/internal/auth/database"
	"github.com/TookenOrg/tooken-services/internal/auth/utils"
)

func (s *Service) SignIn(ctx context.Context, email, password string) (jwtToken string, err error) {

	// 3 - Get user
	user, err := database.GetUserByEmail(ctx, email)
	if err != nil {
		return
	}

	// 4 - Compare password
	if !utils.CheckPasswordHash(password, *user.HashedPassword) {
		err = errors.New("Invalid email or password")
		return
	}

	// 5 - TODO: Update last connexion timestamp

	// 4 - generate jwt for session
	jwtToken, err = utils.GenerateJWT(user.Id, email)
	if err != nil {
		return
	}

	// 5 - return
	return
}
