// internal/blockchain/services/service.go
package services

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	blkServices "github.com/TookenOrg/tooken-services/internal/blockchain/services"
)

type Service struct {
	kycOnChainSync kycOnChainSynchronizer
}

type kycOnChainSynchronizer interface {
	SyncApprovedKYCOnChain(ctx context.Context, kycVerification server.KycVerification) (server.KycOnChainSyncResponse, bool, error)
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
