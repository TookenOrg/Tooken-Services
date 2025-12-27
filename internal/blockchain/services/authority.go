package services

import (
	"context"
	"fmt"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Service) ConfigureAuthority(ctx context.Context) (server.TxHashName, error) {

	trexContracts, err := configureITREXAuthority(ctx)
	if err != nil {
		return server.TxHashName{}, fmt.Errorf("failed to configure ITREX authority: %w", err)
	}

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return server.TxHashName{}, fmt.Errorf("failed to get transact opts: %w", err)
	}

	// Add and use TREX version
	trexAuthorityInstance, err := createTrexAuthorityInstance(ctx)
	if err != nil {
		return server.TxHashName{}, fmt.Errorf("failed to create authority instance: %w", err)
	}

	version := defineAuthorityVersion()

	logger.LogInfo("💌 Configuring Authority with TREX version %d.%d.%d...", version.Major, version.Minor, version.Patch)
	tx, err := trexAuthorityInstance.AddAndUseTREXVersion(
		auth,
		version,
		trexContracts,
	)
	if err != nil {
		return server.TxHashName{}, fmt.Errorf("failed to add and use TREX version: %w", err)
	}

	if err = utils.WaitDeployedTransaction(ctx, tx, false); err != nil {
		return server.TxHashName{}, fmt.Errorf("failed to wait for transaction confirmation: %w", err)
	}
	logger.LogInfo("📬 Authority configured with TREX version %d.%d.%d successfully", version.Major, version.Minor, version.Patch)

	// DB
	// Insert in DB

	return server.TxHashName{
		TransactionHash: tx.Hash().Hex(),
		OperationName:   "Configure Authority with TREX version",
	}, nil

}

func createTrexAuthorityInstance(ctx context.Context) (*contracts.TREXImplementationAuthorityTransactor, error) {
	authorityImplAddress, err := retreiveTrexAuthorityAddress(ctx)
	if err != nil {
		return nil, fmt.Errorf("authority implementation address not found")
	}

	// Create instance of Authority
	return contracts.NewTREXImplementationAuthorityTransactor(authorityImplAddress, globals.EthClient)
}

func retreiveTrexAuthorityAddress(ctx context.Context) (trexAuthorityAddr common.Address, err error) {
	contractDetails, err := database.GetImplementationContractByName(ctx, globals.ImplTrexAuthorityName)
	return common.HexToAddress(contractDetails.Address), err
}

func defineAuthorityVersion() contracts.ITREXImplementationAuthorityVersion {
	return contracts.ITREXImplementationAuthorityVersion{
		Major: 1,
		Minor: 0,
		Patch: 0,
	}
}

func configureITREXAuthority(ctx context.Context) (contracts.ITREXImplementationAuthorityTREXContracts, error) {

	allContractsImplementations, err := database.GetAllContractImplementations(ctx)
	if err != nil {
		return contracts.ITREXImplementationAuthorityTREXContracts{}, fmt.Errorf("failed to get all contract implementations: %w", err)
	}

	tokenImplAddrPtr := utils.FindContractByName(allContractsImplementations, globals.ImplTokenName)
	if tokenImplAddrPtr == nil {
		return contracts.ITREXImplementationAuthorityTREXContracts{}, fmt.Errorf("token implementation address not found")
	}
	tokenImplAddr := *tokenImplAddrPtr

	claimTopicsRegistryAdrrPtr := utils.FindContractByName(allContractsImplementations, globals.ImplClaimTopicRegistryName)
	if claimTopicsRegistryAdrrPtr == nil {
		return contracts.ITREXImplementationAuthorityTREXContracts{}, fmt.Errorf("claim topics registry implementation address not found")
	}
	claimTopicsRegistryAddr := *claimTopicsRegistryAdrrPtr

	identityRegistryPtr := utils.FindContractByName(allContractsImplementations, globals.ImplIdentityRegistryName)
	if identityRegistryPtr == nil {
		return contracts.ITREXImplementationAuthorityTREXContracts{}, fmt.Errorf("identity registry implementation address not found")
	}
	identityRegistryAddr := *identityRegistryPtr

	identityRegistryStoragePtr := utils.FindContractByName(allContractsImplementations, globals.ImplIdentityRegistryStorage)
	if identityRegistryStoragePtr == nil {
		return contracts.ITREXImplementationAuthorityTREXContracts{}, fmt.Errorf("identity registry storage implementation address not found")
	}
	identityRegistryStorageAddr := *identityRegistryStoragePtr

	trustedIssuersRegistryPtr := utils.FindContractByName(allContractsImplementations, globals.ImplTrustedIssuerRegistryName)
	if trustedIssuersRegistryPtr == nil {
		return contracts.ITREXImplementationAuthorityTREXContracts{}, fmt.Errorf("trusted issuers registry implementation address not found")
	}
	trustedIssuersRegistryAddr := *trustedIssuersRegistryPtr

	modularCompliancePtr := utils.FindContractByName(allContractsImplementations, globals.ImplModularComplianceName)
	if modularCompliancePtr == nil {
		return contracts.ITREXImplementationAuthorityTREXContracts{}, fmt.Errorf("modular compliance implementation address not found")
	}
	modularComplianceAddr := *modularCompliancePtr

	return contracts.ITREXImplementationAuthorityTREXContracts{
		TokenImplementation: tokenImplAddr,
		CtrImplementation:   claimTopicsRegistryAddr,
		IrImplementation:    identityRegistryAddr,
		IrsImplementation:   identityRegistryStorageAddr,
		TirImplementation:   trustedIssuersRegistryAddr,
		McImplementation:    modularComplianceAddr,
	}, nil
}
