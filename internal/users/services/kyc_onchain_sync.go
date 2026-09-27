package services

import (
	"context"
	"errors"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	blkDatabase "github.com/TookenOrg/tooken-services/internal/blockchain/database"
	blkServices "github.com/TookenOrg/tooken-services/internal/blockchain/services"
	"github.com/TookenOrg/tooken-services/internal/users/database"
	"github.com/TookenOrg/tooken-services/internal/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
)

const kycClaimTopic = 7

type kycOnChainSyncService struct {
	blockchainSvc *blkServices.Service
}

func newKYCOnChainSyncService(blockchainSvc *blkServices.Service) *kycOnChainSyncService {
	return &kycOnChainSyncService{
		blockchainSvc: blockchainSvc,
	}
}

func (s *kycOnChainSyncService) SyncApprovedKYCOnChain(ctx context.Context, kycVerification server.KycVerification) (isVerifiedOnChain bool, err error) {
	if kycVerification.KycStatus != nil && *kycVerification.KycStatus == database.StatusVerified {
		done, err := s.returnIfAlreadyVerifiedOnChain(ctx, kycVerification)
		if err != nil {
			return false, err
		}
		if done {
			return true, nil
		}
	}

	numericCountryCode, err := utils.CountryAlpha2ToNumeric(kycVerification.DeclaredCountryCode)
	if err != nil {
		return false, ErrMessageInvalidCountryCode
	}

	walletAddress, identityAddress, err := blkServices.EnsureWalletAndIdentity(ctx, kycVerification.UserId)
	if err != nil {
		if errors.Is(err, blkServices.ErrIdentityWithoutWallet) {
			return false, ErrMessageIncoherentUserState
		}
		return false, err
	}
	logger.LogInfo(
		"KYC verification ID=%d has wallet=%s identity=%s country_numeric=%d",
		kycVerification.Id,
		walletAddress.Hex(),
		identityAddress.Hex(),
		numericCountryCode,
	)

	err = blkServices.RegisterIdentityInSharedRegistry(ctx, walletAddress, identityAddress, numericCountryCode)
	if err != nil {
		return false, err
	}

	_, err = s.blockchainSvc.AddClaimToIdentity(ctx, server.AddClaimRequest{
		UserId:     kycVerification.UserId,
		ClaimTopic: kycClaimTopic,
	})
	if err != nil {
		return false, err
	}

	isVerifiedOnChain, err = blkServices.FetchIsVerifiedOnSharedIdentityRegistry(ctx, walletAddress.Hex())
	if err != nil {
		return false, err
	}

	if isVerifiedOnChain {
		logger.LogInfo("KYC verification ID=%d has been successfully verified on-chain", kycVerification.Id)
	} else {
		logger.LogWarn("KYC verification ID=%d failed to be verified on-chain", kycVerification.Id)
	}

	return isVerifiedOnChain, nil
}

func (s *kycOnChainSyncService) returnIfAlreadyVerifiedOnChain(ctx context.Context, kycVerification server.KycVerification) (bool, error) {
	logger.LogInfo("KYC verification ID=%d is already verified in DB for user ID=%d, checking on-chain status", kycVerification.Id, kycVerification.UserId)

	userWalletAddrPtr, err := blkDatabase.GetWalletByUserId(ctx, kycVerification.UserId)
	if err != nil {
		return false, err
	}
	if userWalletAddrPtr == nil {
		return false, ErrMessageInvalidWalletAddress
	}

	isVerifiedOnChain, err := blkServices.FetchIsVerifiedOnSharedIdentityRegistry(ctx, userWalletAddrPtr.Hex())
	if err != nil {
		return false, err
	}
	if isVerifiedOnChain {
		logger.LogInfo("KYC verification ID=%d is already verified on-chain", kycVerification.Id)
		return true, nil
	}

	logger.LogInfo("KYC verification ID=%d is not verified on-chain", kycVerification.Id)
	return false, nil
}
