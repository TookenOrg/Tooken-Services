package database

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
)

func GetUserByID(ctx context.Context, userID int) (user *server.User, err error) {
	query := `
        SELECT id, full_name, email, created_at, role, kyc_status, kyc_expires_at, kyc_verified_at, country_code, updated_at, last_connexion
        FROM usr.users
        WHERE id = $1
    `
	user = &server.User{}
	err = globals.DB.QueryRow(query, userID).Scan(
		&user.Id,
		&user.FullName,
		&user.Email,
		&user.CreatedAt,
		&user.Role,
		&user.KycStatus,
		&user.KycExpiresAt,
		&user.KycVerifiedAt,
		&user.CountryCode,
		&user.UpdatedAt,
		&user.LastConnexion,
	)
	if err != nil {
		return
	}

	logger.LogDebug("User retrieved with ID [%d]", userID)
	return
}
