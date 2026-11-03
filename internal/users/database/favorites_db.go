package database

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
)

func GetFavoriteRealEstateIDs(ctx context.Context, userID int) (map[int]struct{}, error) {
	rows, err := globals.DB.QueryContext(ctx, `
		SELECT real_estate_id FROM rel.user_asset_favorite WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := make(map[int]struct{})
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids[id] = struct{}{}
	}
	return ids, rows.Err()
}

func AddFavoritesRealEstateByUserId(ctx context.Context, userId, realEstateId int) (err error) {
	query := `
        INSERT INTO rel.user_asset_favorite
            (user_id, real_estate_id)
        VALUES ($1, $2)
        RETURNING id
    `
	var id int
	err = globals.DB.QueryRowContext(ctx, query, userId, realEstateId).Scan(&id)
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
	err = globals.DB.QueryRowContext(ctx, query, userId, realEstateId).Scan(&id)
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
