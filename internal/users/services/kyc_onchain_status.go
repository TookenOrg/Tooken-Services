package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	blkDatabase "github.com/TookenOrg/tooken-services/internal/blockchain/database"
	blkServices "github.com/TookenOrg/tooken-services/internal/blockchain/services"
	"github.com/ethereum/go-ethereum/common"
)

var ErrKYCOnChainUnavailable = errors.New("on-chain KYC status unavailable")

type kycOnChainStatusReader struct {
	getWalletByUserID func(context.Context, int) (*common.Address, error)
	fetchIsVerified   func(context.Context, string) (bool, error)
}

func (s *Service) GetCurrentUserKycOnChain(ctx context.Context, userID int) (server.KycOnChainStatus, error) {
	if _, err := s.GetUserByID(ctx, userID); err != nil {
		return server.KycOnChainStatus{}, err
	}
	reader := kycOnChainStatusReader{
		getWalletByUserID: blkDatabase.GetWalletByUserId,
		fetchIsVerified:   blkServices.FetchIsVerifiedOnSharedIdentityRegistry,
	}
	return reader.read(ctx, userID)
}

func (s kycOnChainStatusReader) read(ctx context.Context, userID int) (server.KycOnChainStatus, error) {
	status := server.KycOnChainStatus{UserId: userID}
	wallet, err := s.getWalletByUserID(ctx, userID)
	if err != nil {
		return status, err
	}
	if wallet == nil {
		return status, nil
	}
	address := wallet.Hex()
	status.WalletAddress = &address
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	verified, err := s.fetchIsVerified(ctx, address)
	if err != nil {
		return status, fmt.Errorf("%w: %v", ErrKYCOnChainUnavailable, err)
	}
	status.IsVerified = verified
	return status, nil
}
