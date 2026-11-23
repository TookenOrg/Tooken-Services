package services

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	assetServices "github.com/TookenOrg/tooken-services/internal/assets_managements/services"
	blkServices "github.com/TookenOrg/tooken-services/internal/blockchain/services"
)

type Service struct {
	kycOnChainSync   kycOnChainSynchronizer
	realEstateReader realEstateSummaryReader
}

type kycOnChainSynchronizer interface {
	SyncApprovedKYCOnChain(ctx context.Context, kycVerification server.KycVerification) (server.KycOnChainSyncResponse, bool, error)
}

type realEstateSummaryReader interface {
	GetRealEstatesByIDs(ctx context.Context, ids []int, isStaff bool) ([]server.RealEstateSummary, error)
}

func NewService() *Service {
	return NewServiceWithDependencies(blkServices.NewService(), assetServices.NewService(nil))
}

func NewServiceWithDependencies(blkSvc *blkServices.Service, realEstateReader realEstateSummaryReader) *Service {
	if blkSvc == nil {
		blkSvc = blkServices.NewService()
	}
	if realEstateReader == nil {
		realEstateReader = assetServices.NewService(nil)
	}
	return &Service{
		kycOnChainSync:   newKYCOnChainSyncService(blkSvc),
		realEstateReader: realEstateReader,
	}
}

func (s *Service) getKYCOnChainSync() kycOnChainSynchronizer {
	if s.kycOnChainSync == nil {
		s.kycOnChainSync = newKYCOnChainSyncService(blkServices.NewService())
	}
	return s.kycOnChainSync
}

func (s *Service) getRealEstateReader() realEstateSummaryReader {
	if s.realEstateReader == nil {
		s.realEstateReader = assetServices.NewService(nil)
	}
	return s.realEstateReader
}
