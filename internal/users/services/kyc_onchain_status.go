package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	blkDatabase "github.com/TookenOrg/tooken-services/internal/blockchain/database"
	blkGlobals "github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
		fetchIsVerified:   fetchSharedRegistryVerification,
	}
	return reader.read(ctx, userID)
}

func fetchSharedRegistryVerification(ctx context.Context, wallet string) (bool, error) {
	if blkGlobals.EthClient == nil {
		return false, ErrKYCOnChainUnavailable
	}
	storage, err := blkDatabase.GetContractRoleByName(ctx, blkGlobals.SharedIdentityRegistryStorageName)
	if err != nil {
		return false, err
	}
	if !common.IsHexAddress(storage.Address) || common.HexToAddress(storage.Address) == (common.Address{}) {
		return false, errors.New("shared identity storage is not configured")
	}
	irs, err := contracts.NewIdentityRegistryStorage(common.HexToAddress(storage.Address), blkGlobals.EthClient)
	if err != nil {
		return false, err
	}
	opts := &bind.CallOpts{Context: ctx}
	registries, err := irs.LinkedIdentityRegistries(opts)
	if err != nil {
		return false, err
	}
	if len(registries) == 0 {
		return false, errors.New("no registry is bound to the shared identity storage")
	}
	// Reads do not require agent permissions or a signing key. The write
	// resolver checks both, so use a registry linked by the shared storage.
	registry, err := contracts.NewIdentityRegistry(registries[0], blkGlobals.EthClient)
	if err != nil {
		return false, err
	}
	return registry.IsVerified(opts, common.HexToAddress(wallet))
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
