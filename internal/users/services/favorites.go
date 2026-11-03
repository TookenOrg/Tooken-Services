package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	assetServices "github.com/TookenOrg/tooken-services/internal/assets_managements/services"
	"github.com/TookenOrg/tooken-services/internal/users/database"
	"github.com/lib/pq"
)

func (s *Service) GetFavoritesRealEstatesForUser(ctx context.Context, userID int, isStaff bool) ([]server.RealEstateSummary, error) {
	if _, err := s.GetUserByID(ctx, userID); err != nil {
		return nil, err
	}
	ids, err := database.GetFavoriteRealEstateIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []server.RealEstateSummary{}, nil
	}
	// Reuse the catalogue's visibility and safe summary mapping in one query,
	// rather than loading every favorite's sensitive detail separately.
	estates, err := assetServices.NewService(nil).GetRealEstates(ctx, isStaff)
	if err != nil {
		return nil, err
	}
	favorites := make([]server.RealEstateSummary, 0, len(ids))
	for _, estate := range estates {
		if _, ok := ids[estate.Id]; ok {
			favorites = append(favorites, estate)
		}
	}
	return favorites, nil
}

var (
	// ErrFavoriteExists is returned when the real estate already is a
	// favorite of the user.
	ErrFavoriteExists = errors.New("the real estate already is a favorite of the user")
	// ErrFavoriteNotFound is returned when removing a real estate that is not
	// a favorite of the user.
	ErrFavoriteNotFound = errors.New("the real estate is not a favorite of the user")
	// ErrUnknownRealEstate is returned when the real estate does not exist.
	ErrUnknownRealEstate = errors.New("the real estate does not exist")
)

func (s *Service) AddFavoritesRealEstateForUser(ctx context.Context, userId, realEstateId int) error {
	return translateFavoriteWriteError(database.AddFavoritesRealEstateByUserId(ctx, userId, realEstateId))
}

func (s *Service) RemoveFavoritesRealEstateForUser(ctx context.Context, userId, realEstateId int) error {
	err := database.RemoveFavoritesRealEstateByUserId(ctx, userId, realEstateId)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrFavoriteNotFound
	}
	return err
}

// translateFavoriteWriteError turns the constraints of rel.user_asset_favorite
// into sentinels, so that a refusal the caller can act on is not answered as
// a server failure.
//
// The insert is left to fail rather than checked beforehand: a lookup then an
// insert would let two identical requests both pass the lookup, while the
// constraints decide atomically.
func translateFavoriteWriteError(err error) error {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return err
	}

	switch {
	case pqErr.Code.Name() == "unique_violation" && pqErr.Constraint == "uq_user_asset_favorite":
		return ErrFavoriteExists
	case pqErr.Code.Name() == "foreign_key_violation" && pqErr.Constraint == "fk_fav_asset":
		return ErrUnknownRealEstate
	case pqErr.Code.Name() == "foreign_key_violation" && pqErr.Constraint == "fk_fav_user":
		// The token is valid but its account was deleted since it was issued.
		return ErrUserNotFound
	}
	return err
}
