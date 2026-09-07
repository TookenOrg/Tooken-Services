package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/assets_managements/database"
)

// GetRealEstates serves the listing.
//
// isStaff comes from the caller's role, not from the request: a MANAGER or an
// ADMIN also sees the assets that are not published yet, everyone else sees the
// public catalogue.
func (s *Service) GetRealEstates(ctx context.Context, isStaff bool) (realEstates []server.RealEstateSummary, err error) {
	// 1 - Call DB to fetch real_estate
	dtos, err := database.GetRealEstates(ctx, isStaff)
	if err != nil {
		return nil, err
	}

	// 2 - Translate the schema into the API contract
	return toServerRealEstates(dtos), nil
}

// GetRealEstateById serves the detail.
//
// isStaff decides two things at once: whether an unpublished asset is readable
// at all, and whether the street address is part of the answer. An asset the
// caller may not read reports ErrRealEstateNotFound, so a draft is
// indistinguishable from an identifier that never existed.
func (s *Service) GetRealEstateById(ctx context.Context, id int, isStaff bool) (realEstate server.RealEstate, err error) {
	dto, err := database.GetRealEstateById(ctx, id, isStaff)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.RealEstate{}, ErrRealEstateNotFound
		}
		return server.RealEstate{}, err
	}

	return toServerRealEstate(dto, isStaff), nil
}
