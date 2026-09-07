package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/auth/database"
	"github.com/TookenOrg/tooken-services/internal/auth/utils"
)

// ErrInvalidCredentials: the caller gave an unknown email or a wrong password.
//
// Every other failure must NOT reach the client as this error. Reporting a
// database outage as "invalid email or password" sends the user hunting for a
// password problem that does not exist, and hides the incident from whoever
// reads the logs.
var ErrInvalidCredentials = errors.New("invalid email or password")

func (s *Service) SignIn(ctx context.Context, email, password string) (user server.User, jwtToken string, refreshToken string, expiresAt time.Time, err error) {

	// 1 - Get user
	userDto, err := database.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrInvalidCredentials
		} else {
			// Typically a missing column or an unreachable database. Wrapped so
			// the cause survives all the way to the log.
			err = fmt.Errorf("cannot read the account: %w", err)
		}
		return
	}

	// 2 - Compare password
	if !utils.CheckPasswordHash(password, userDto.HashedPassword) {
		err = ErrInvalidCredentials
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
