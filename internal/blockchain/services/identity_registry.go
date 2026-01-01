package services

import (
	"context"
	"math/big"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func registerIdentity(ctx context.Context, userWallet, identityAddress common.Address, countryCode int) (tx *types.Transaction, err error) {

	irInstance, err := getIdentityRegistryInstance(ctx)

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Registration of new identity in identity registry...")
	tx, err = irInstance.RegisterIdentity(auth, userWallet, identityAddress, uint16(countryCode))
	if err != nil {
		return
	}
	deployedTxDetails, err := utils.WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return
	}

	logger.LogInfo("📬 Identity registred on transaction: %s", tx.Hash().Hex())

	err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "REGISTER_IDENTITY", deployedTxDetails.Tx.To().Hex(), deployedTxDetails.BlockNumber.Int64(), big.Int{})
	if err != nil {
		return
	}

	return
}

func addAgentOnIdentityRegistry(ctx context.Context, tokenAddr common.Address) (tx *types.Transaction, err error) {

	irInstance, err := getIdentityRegistryInstance(ctx)
	if err != nil {
		return
	}

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Add agent on identity registry...")
	tx, err = irInstance.AddAgent(auth, tokenAddr)
	if err != nil {
		return
	}
	deployedTxDetails, err := utils.WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return
	}

	logger.LogInfo("📬 Agent added on identity registry: %s", tx.Hash().Hex())

	err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "ADD_AGENT_IDENTITY_REGISTRY", deployedTxDetails.Tx.To().Hex(), deployedTxDetails.BlockNumber.Int64(), big.Int{})
	if err != nil {
		return
	}

	return
}

func getIdentityRegistryInstance(ctx context.Context) (irInstance *contracts.IdentityRegistry, err error) {

	identityRegistryDetails, err := database.GetContractInstanceByName(ctx, globals.IdentityRegistry)
	if err != nil {
		return
	}
	irAddress := common.HexToAddress(identityRegistryDetails.Address)

	irInstance, err = contracts.NewIdentityRegistry(irAddress, globals.EthClient)
	if err != nil {
		return
	}

	return
}
