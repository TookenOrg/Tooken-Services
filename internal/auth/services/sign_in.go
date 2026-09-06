package services

import (
	"context"
	"errors"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/auth/database"
	"github.com/TookenOrg/tooken-services/internal/auth/utils"
)

func (s *Service) SignIn(ctx context.Context, email, password string) (user server.User, jwtToken string, refreshToken string, expiresAt time.Time, err error) {

	// 1 - Get user
	userDto, err := database.GetUserByEmail(ctx, email)
	if err != nil {
		return
	}

	// 2 - Compare password
	if !utils.CheckPasswordHash(password, userDto.HashedPassword) {
		err = errors.New("Invalid email or password")
		return
	}

	utils.MapOneToOne(userDto, &user)

	// 3 - TODO: Update last connexion timestamp

	// 4 - generate jwt for session
	jwtToken, err = utils.GenerateJWT(userDto.Id, email, userDto.Role)
	if err != nil {
		return
	}

	// 5 - generate refresh token (8h)
	refreshToken, err = utils.GenerateRefreshToken()
	if err != nil {
		return
	}

	// 6 - ExpiresAt pour le frontend
	claims, _ := utils.ParseJWT(jwtToken) // parse pour extraire ExpiresAt
	expiresAt = claims.ExpiresAt.Time

	// 7 - return
	return
}
