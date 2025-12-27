package services

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Service) DeployIdentityFactory(ctx context.Context) (details server.ContractDetails, err error) {

	// Get Identity Authority Implementation
	implementationIdentityAuthorityDetails, err := database.GetImplementationContractByName(ctx, globals.ImplIdentityAuthorityName)
	if err != nil {
		return
	}

	implIdentityAuthorityAddr := common.HexToAddress(implementationIdentityAuthorityDetails.Address)

	identityFactory, err := deployIdentityFactory(ctx, implIdentityAuthorityAddr)
	if err != nil {
		return server.ContractDetails{}, err
	}

	details = server.ContractDetails{
		Address: identityFactory.Hex(),
		Name:    globals.IdentityFactoryName,
	}

	// TODO: Store deployed contract details in the database

	return
}

func deployIdentityFactory(ctx context.Context, implementationAuthorityAddrCommon common.Address) (idAddr common.Address, err error) {

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Deploying Identity Factory...")
	idAddr, tx, _, err := contracts.DeployIdFactory(auth, globals.EthClient, implementationAuthorityAddrCommon)
	if err != nil {
		return
	}
	if err = utils.WaitDeployedTransaction(ctx, tx, true); err != nil {
		return
	}
	logger.LogInfo("📬 Identity Factory deployed at address: %s", idAddr.Hex())
	return
}
