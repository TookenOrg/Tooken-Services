package services

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/assets_managements/database"
)

func (s *Service) GetActiveRealEstates(ctx context.Context) (realEstates []server.RealEstateSummary, err error) {
	// 1 - Call DB to fetch real_estate
	dtos, err := database.GetActiveRealEstates(ctx)
	if err != nil {
		return nil, err
	}

	// 2 - Translate the schema into the API contract
	return toServerRealEstates(dtos), nil
}

func (s *Service) GetActiveRealEstateById(ctx context.Context, id int) (realEstate server.RealEstate, err error) {
	dto, err := database.GetActiveRealEstateById(ctx, id)
	if err != nil {
		return server.RealEstate{}, err
	}

	return toServerRealEstate(dto), nil
}
