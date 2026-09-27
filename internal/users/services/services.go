// internal/blockchain/services/service.go
package services

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	blkServices "github.com/TookenOrg/tooken-services/internal/blockchain/services"
	"github.com/TookenOrg/tooken-services/internal/users/database"
)

type Service struct {
	kycOnChainSync kycOnChainSynchronizer
}

type kycOnChainSynchronizer interface {
	SyncApprovedKYCOnChain(ctx context.Context, kycVerification server.KycVerification) error
}

func NewService() *Service {
	return &Service{
		kycOnChainSync: newKYCOnChainSyncService(blkServices.NewService()),
	}
}

func (s *Service) getKYCOnChainSync() kycOnChainSynchronizer {
	if s.kycOnChainSync == nil {
		s.kycOnChainSync = newKYCOnChainSyncService(blkServices.NewService())
	}
	return s.kycOnChainSync
}

func (s *Service) AddFavoritesRealEstateForUser(ctx context.Context, userId, realEstateId int) (err error) {
	// 1 - Insert
	return database.AddFavoritesRealEstateByUserId(ctx, userId, realEstateId)
}

func (s *Service) RemoveFavoritesRealEstateForUser(ctx context.Context, userId, realEstateId int) (err error) {
	// 1 - Remove
	return database.RemoveFavoritesRealEstateByUserId(ctx, userId, realEstateId)
}
