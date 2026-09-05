package services

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/assets_managements/database"
)

func (s *Service) GetActiveRealEstates(ctx context.Context) (realEstates []server.RealEstate, err error) {
	// 1 - Call DB to fetch real_estate
	return database.GetActiveRealEstates(ctx)
}

func (s *Service) GetActiveRealEstateById(ctx context.Context, id int) (realEstate server.RealEstate, err error) {
	return database.GetActiveRealEstateById(ctx, id)
}
