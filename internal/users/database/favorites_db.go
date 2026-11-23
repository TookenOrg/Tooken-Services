package database

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
)

func GetFavoriteRealEstateIDs(ctx context.Context, userID int) ([]int, error) {
	rows, err := globals.DB.QueryContext(ctx, `
		SELECT real_estate_id
		FROM rel.user_asset_favorite
		WHERE user_id = $1
		ORDER BY real_estate_id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// AddFavoritesRealEstateByUserId inserts a favorite only when the target asset
// is visible to the caller, so probing an unpublished or soft-deleted asset ID
// returns sql.ErrNoRows just like an identifier that never existed.
func AddFavoritesRealEstateByUserId(ctx context.Context, userId, realEstateId int, includeUnpublished bool) (err error) {
	query := `
        WITH visible_asset AS (
            SELECT id
            FROM ass.real_estate
            WHERE id = $2
              AND deleted_at IS NULL
              AND (active = true OR $3)
        )
        INSERT INTO rel.user_asset_favorite
            (user_id, real_estate_id)
        SELECT $1, id FROM visible_asset
        RETURNING id
    `
	var id int
	err = globals.DB.QueryRowContext(ctx, query, userId, realEstateId, includeUnpublished).Scan(&id)
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
