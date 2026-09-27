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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const kycClaimTopic = 7

type kycOnChainSyncService struct {
	addClaimToIdentity               func(context.Context, server.AddClaimRequest) (*types.Transaction, error)
	ensureWalletAndIdentity          func(context.Context, int) (common.Address, common.Address, *string, error)
	fetchIsVerified                  func(context.Context, string) (bool, error)
	getWalletByUserID                func(context.Context, int) (*common.Address, error)
	registerIdentityInSharedRegistry func(context.Context, common.Address, common.Address, int) (*string, error)
}

func newKYCOnChainSyncService(blockchainSvc *blkServices.Service) *kycOnChainSyncService {
	return &kycOnChainSyncService{
		addClaimToIdentity:               blockchainSvc.AddClaimToIdentity,
		ensureWalletAndIdentity:          blkServices.EnsureWalletAndIdentity,
		fetchIsVerified:                  blkServices.FetchIsVerifiedOnSharedIdentityRegistry,
		getWalletByUserID:                blkDatabase.GetWalletByUserId,
		registerIdentityInSharedRegistry: blkServices.RegisterIdentityInSharedRegistry,
	}
}

func (s *kycOnChainSyncService) SyncApprovedKYCOnChain(ctx context.Context, kycVerification server.KycVerification) (response server.KycOnChainSyncResponse, isVerifiedOnChain bool, err error) {
	numericCountryCode, err := utils.CountryAlpha2ToNumeric(kycVerification.DeclaredCountryCode)
	if err != nil {
		return response, false, ErrMessageInvalidCountryCode
	}

	walletAddress, identityAddress, identityTxHash, err := s.ensureWalletAndIdentity(ctx, kycVerification.UserId)
	if err != nil {
		if errors.Is(err, blkServices.ErrIdentityWithoutWallet) {
			return response, false, ErrMessageIncoherentUserState
		}
		return response, false, err
	}

	transactions := server.KycOnChainSyncTransactions{
		Identity: identityTxHash,
	}
	response = server.KycOnChainSyncResponse{
		UserId:          kycVerification.UserId,
		WalletAddress:   walletAddress.Hex(),
		IdentityAddress: identityAddress.Hex(),
		CountryCode:     kycVerification.DeclaredCountryCode,
		CountryNumeric:  numericCountryCode,
		Transactions:    transactions,
	}

	logger.LogInfo(
		"KYC verification ID=%d has wallet=%s identity=%s country_numeric=%d",
		kycVerification.Id,
		walletAddress.Hex(),
		identityAddress.Hex(),
		numericCountryCode,
	)

	if kycVerification.KycStatus != nil && *kycVerification.KycStatus == database.StatusVerified {
		done, err := s.returnIfAlreadyVerifiedOnChain(ctx, kycVerification)
		if err != nil {
			return response, false, err
		}
		if done {
			return response, true, nil
		}
	}

	claimTx, err := s.addClaimToIdentity(ctx, server.AddClaimRequest{
		UserId:     kycVerification.UserId,
		ClaimTopic: kycClaimTopic,
	})
	if err != nil {
		return response, false, err
	}
	if claimTx != nil {
		hash := claimTx.Hash().Hex()
		response.Transactions.Claim = &hash
	}

	registerIdentityTxHash, err := s.registerIdentityInSharedRegistry(ctx, walletAddress, identityAddress, numericCountryCode)
	if err != nil {
		return response, false, err
	}
	response.Transactions.RegisterIdentity = registerIdentityTxHash

	isVerifiedOnChain, err = s.fetchIsVerified(ctx, walletAddress.Hex())
	if err != nil {
		return response, false, err
	}

	if isVerifiedOnChain {
		logger.LogInfo("KYC verification ID=%d has been successfully verified on-chain", kycVerification.Id)
	} else {
		logger.LogWarn("KYC verification ID=%d failed to be verified on-chain", kycVerification.Id)
	}

	return response, isVerifiedOnChain, nil
}

func (s *kycOnChainSyncService) returnIfAlreadyVerifiedOnChain(ctx context.Context, kycVerification server.KycVerification) (bool, error) {
	logger.LogInfo("KYC verification ID=%d is already verified in DB for user ID=%d, checking on-chain status", kycVerification.Id, kycVerification.UserId)

	userWalletAddrPtr, err := s.getWalletByUserID(ctx, kycVerification.UserId)
	if err != nil {
		return false, err
	}
	if userWalletAddrPtr == nil {
		return false, ErrMessageInvalidWalletAddress
	}

	isVerifiedOnChain, err := s.fetchIsVerified(ctx, userWalletAddrPtr.Hex())
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
