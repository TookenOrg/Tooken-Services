// internal/blockchain/services/service.go
package services

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/users/database"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) AddFavoritesRealEstateForUser(ctx context.Context, userId, realEstateId int) (err error) {
	// 1 - Insert
	return database.AddFavoritesRealEstateByUserId(ctx, userId, realEstateId)
}

func (s *Service) RemoveFavoritesRealEstateForUser(ctx context.Context, userId, realEstateId int) (err error) {
	// 1 - Remove
	return database.RemoveFavoritesRealEstateByUserId(ctx, userId, realEstateId)
}
