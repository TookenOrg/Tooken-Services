package services

import (
	"context"
	"math/big"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/models"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func deployTokenProxy(ctx context.Context, trexAuthorityImplementationAddr common.Address, complianceSuite models.ComplianceSuite, tokenName string, symbol string, nbDecimal int) (tokenAddr common.Address, tx *types.Transaction, tokenProxyInstance *contracts.TokenProxy, err error) {

	identityRegistryDetails, err := database.GetContractInstanceByName(ctx, globals.IdentityRegistry)
	if err != nil {
		return
	}
	irAddress := common.HexToAddress(identityRegistryDetails.Address)

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Deploying new token proxy for token [%s]", tokenName)
	tokenAddr, tx, tokenProxyInstance, err = contracts.DeployTokenProxy(
		auth,
		globals.EthClient,
		trexAuthorityImplementationAddr,
		irAddress,
		complianceSuite.Address,
		tokenName,
		symbol,
		uint8(nbDecimal),
		common.Address{},
	)

	deployedTxDetails, err := utils.WaitDeployedTransaction(ctx, tx, true)
	if err != nil {
		return
	}
	logger.LogInfo("📬 Token Proxy deployed at address [%s]", tx.Hash().Hex())
	logger.LogInfo("Token Proxy generate new token at address %s", tokenAddr)

	err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "TOKEN_PROXY", deployedTxDetails.Tx.To().Hex(), deployedTxDetails.BlockNumber.Int64(), big.Int{})
	if err != nil {
		return
	}

	return
}
