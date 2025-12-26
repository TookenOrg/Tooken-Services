package services

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Service) DeployAllImplementations(ctx context.Context) (deployedContracts []server.ContractDetails, err error) {

	// 1. Deploy ClaimTopicsRegistry Implementation
	ctrAddr, err := deployClaimTopicsRegistryImplementation(ctx)
	if err != nil {
		return
	}
	deployedContracts = append(deployedContracts, server.ContractDetails{
		Address: ctrAddr.Hex(),
		Name:    globals.ImplClaimTopicRegistryName,
	})

	// 2. Deploy TrustedIssuerRegistry Implementation
	tirAddr, err := deployTrustedIssuersRegistryImplementation(ctx)
	if err != nil {
		return
	}
	deployedContracts = append(deployedContracts, server.ContractDetails{
		Address: tirAddr.Hex(),
		Name:    globals.ImplTrustedIssuersRegistryName,
	})

	// 3. Deploy IdentityRegistryStorage Implementation
	irsAddr, err := deployIdentityRegistryStorageImplementation(ctx)
	if err != nil {
		return
	}
	deployedContracts = append(deployedContracts, server.ContractDetails{
		Address: irsAddr.Hex(),
		Name:    globals.ImplIdentityRegistryStorage,
	})

	// 4. Deploy IdentityRegistry Implementation
	irAddr, err := deployIdentityRegistryImplementation(ctx)
	if err != nil {
		return
	}
	deployedContracts = append(deployedContracts, server.ContractDetails{
		Address: irAddr.Hex(),
		Name:    globals.ImplIdentityRegistryName,
	})

	// 5. Deploy ModularCompliance Implementation
	mcAddr, err := deployModularComplianceImplementation(ctx)
	if err != nil {
		return
	}
	deployedContracts = append(deployedContracts, server.ContractDetails{
		Address: mcAddr.Hex(),
		Name:    globals.ImplModularComplianceName,
	})

	// 6. Deploy Token Implementation
	tokenAddr, err := deployTokenImplementation(ctx)
	if err != nil {
		return
	}
	deployedContracts = append(deployedContracts, server.ContractDetails{
		Address: tokenAddr.Hex(),
		Name:    globals.ImplTokenName,
	})

	// 7. Deploy Identity Implementation
	identityAddr, err := deployIdentityImplementation(ctx)
	if err != nil {
		return
	}
	deployedContracts = append(deployedContracts, server.ContractDetails{
		Address: identityAddr.Hex(),
		Name:    globals.ImplIdentityName,
	})

	// 8. Deploy IdentityAuthority Implementation
	identityAuthorityAddr, err := deployIdentityAuthorityImplementation(ctx, identityAddr)
	if err != nil {
		return
	}
	deployedContracts = append(deployedContracts, server.ContractDetails{
		Address: identityAuthorityAddr.Hex(),
		Name:    globals.ImplIdentityAuthorityName,
	})

	// 9. Deploy TREXImplementationAuthority
	authorityAddr, err := deployTREXImplementationAuthority(ctx)
	if err != nil {
		return
	}
	deployedContracts = append(deployedContracts, server.ContractDetails{
		Address: authorityAddr.Hex(),
		Name:    globals.ImplTrexAuthorityName,
	})

	logger.LogInfo("✅ All TREX implementation contracts deployed")

	// DB
	// err = database.InsertContractsImplementation(ctx, deployedContracts)
	// if err != nil {
	// 	return
	// }

	return
}

func deployClaimTopicsRegistryImplementation(ctx context.Context) (ctrAddr common.Address, err error) {
	logger.LogInfo("💌 Deploying ClaimTopicsRegistry Implementation...")
	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}
	ctrAddr, tx, _, err := contracts.DeployClaimTopicsRegistry(auth, globals.EthClient)
	if err != nil {
		return
	}

	if err = utils.WaitDeployedTransaction(ctx, tx, true); err != nil {
		return
	}
	logger.LogInfo("📬 ClaimTopicsRegistry Implementation deployed at address: %s", ctrAddr.Hex())
	return
}

func deployTrustedIssuersRegistryImplementation(ctx context.Context) (tirAddr common.Address, err error) {
	logger.LogInfo("💌 Deploying TrustedIssuerRegistry Implementation...")
	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}
	tirAddr, tx, _, err := contracts.DeployTrustedIssuersRegistry(auth, globals.EthClient)
	if err != nil {
		return
	}
	if err = utils.WaitDeployedTransaction(ctx, tx, true); err != nil {
		return
	}
	logger.LogInfo("📬 TrustedIssuerRegistry Implementation deployed at address: %s", tirAddr.Hex())
	return
}

func deployIdentityRegistryStorageImplementation(ctx context.Context) (irsAddr common.Address, err error) {
	logger.LogInfo("💌 Deploying IdentityRegistryStorage Implementation...")
	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}
	irsAddr, tx, _, err := contracts.DeployIdentityRegistryStorage(auth, globals.EthClient)
	if err != nil {
		return
	}
	if err = utils.WaitDeployedTransaction(ctx, tx, true); err != nil {
		return
	}
	logger.LogInfo("📬 IdentityRegistryStorage Implementation deployed at address: %s", irsAddr.Hex())
	return
}

func deployIdentityRegistryImplementation(ctx context.Context) (irAddr common.Address, err error) {
	logger.LogInfo("💌 Deploying IdentityRegistry Implementation...")
	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}
	irAddr, tx, _, err := contracts.DeployIdentityRegistry(auth, globals.EthClient)
	if err != nil {
		return
	}
	if err = utils.WaitDeployedTransaction(ctx, tx, true); err != nil {
		return
	}
	logger.LogInfo("📬 IdentityRegistry Implementation deployed at address: %s", irAddr.Hex())
	return
}

func deployModularComplianceImplementation(ctx context.Context) (mcAddr common.Address, err error) {
	logger.LogInfo("💌 Deploying ModularCompliance Implementation...")
	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}
	mcAddr, tx, _, err := contracts.DeployModularCompliance(auth, globals.EthClient)
	if err != nil {
		return
	}
	if err = utils.WaitDeployedTransaction(ctx, tx, true); err != nil {
		return
	}
	logger.LogInfo("📬 ModularCompliance Implementation deployed at address: %s", mcAddr.Hex())
	return
}

func deployTokenImplementation(ctx context.Context) (tokenAddr common.Address, err error) {
	logger.LogInfo("💌 Deploying Token Implementation...")
	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}
	tokenAddr, tx, _, err := contracts.DeployToken(auth, globals.EthClient)
	if err != nil {
		return
	}
	if err = utils.WaitDeployedTransaction(ctx, tx, true); err != nil {
		return
	}
	logger.LogInfo("📬 Token Implementation deployed at address: %s", tokenAddr.Hex())
	return
}

func deployIdentityImplementation(ctx context.Context) (identityAddr common.Address, err error) {
	logger.LogInfo("💌 Deploying Identity Implementation...")
	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}
	identityAddr, tx, _, err := contracts.DeployIdentity(auth, globals.EthClient, auth.From, false)
	if err != nil {
		return
	}
	if err = utils.WaitDeployedTransaction(ctx, tx, true); err != nil {
		return
	}
	logger.LogInfo("📬 Identity Implementation deployed at address: %s", identityAddr.Hex())
	return
}

func deployIdentityAuthorityImplementation(ctx context.Context, identity common.Address) (identityAuthAddr common.Address, err error) {
	logger.LogInfo("💌 Deploying IdentityAuthority Implementation...")
	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}
	identityAuthAddr, tx, _, err := contracts.DeployImplementationAuthority(auth, globals.EthClient, identity)
	if err != nil {
		return
	}
	if err = utils.WaitDeployedTransaction(ctx, tx, true); err != nil {
		return
	}
	logger.LogInfo("📬 IdentityAuthority Implementation deployed at address: %s", identityAuthAddr.Hex())
	return
}

func deployTREXImplementationAuthority(ctx context.Context) (authorityAddr common.Address, err error) {
	logger.LogInfo("💌 Deploying TREXImplementationAuthority Implementation...")
	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}
	authorityAddr, tx, _, err := contracts.DeployTREXImplementationAuthority(
		auth,
		globals.EthClient,
		true,             // isReference - To use AddAndUseTREXVersion
		common.Address{}, // trexFactory address
		common.Address{}, // identityFactory address
	)
	if err != nil {
		return
	}
	if err = utils.WaitDeployedTransaction(ctx, tx, true); err != nil {
		return
	}
	logger.LogInfo("📬 TREXImplementationAuthority Implementation deployed at address: %s", authorityAddr.Hex())
	return
}
