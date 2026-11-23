package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/TookenOrg/tooken-services/internal/api/server"
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
	return s.getRealEstateReader().GetRealEstatesByIDs(ctx, ids, isStaff)
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

func (s *Service) AddFavoritesRealEstateForUser(ctx context.Context, userId, realEstateId int, isStaff bool) error {
	return translateFavoriteWriteError(database.AddFavoritesRealEstateByUserId(ctx, userId, realEstateId, isStaff))
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
	if errors.Is(err, sql.ErrNoRows) {
		return ErrUnknownRealEstate
	}

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
