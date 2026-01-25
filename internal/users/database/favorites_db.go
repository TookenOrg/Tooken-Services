package database

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
)

func AddFavoritesRealEstateByUserId(ctx context.Context, userId, realEstateId int) (err error) {
	query := `
        INSERT INTO rel.user_asset_favorite
            (user_id, real_estate_id)
        VALUES ($1, $2)
        RETURNING id
    `
	var id int
	err = globals.DB.QueryRow(query, userId, realEstateId).Scan(&id)
	if err != nil {
		return
	}

	logger.LogDebug("Favorite added on %d for real-estate %d. Id [%d]", userId, realEstateId, id)
	return
}

func RemoveFavoritesRealEstateByUserId(ctx context.Context, userId, realEstateId int) (err error) {
	query := `
        DELETE FROM rel.user_asset_favorite
        WHERE user_id = $1
		AND real_estate_id = $2
    RETURNING id
	`

	var id int64
	err = globals.DB.QueryRow(query, userId, realEstateId).Scan(&id)
	if err != nil {
		return
	}

	logger.LogDebug(
		"Favorite removed (id=%d): user %d -> real-estate %d",
		id,
		userId,
		realEstateId,
	)

	return
}
