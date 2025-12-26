package services

import (
	"context"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/models"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func deployTokenProxy(ctx context.Context, trexAuthorityImplementationAddr common.Address, complianceSuite models.ComplianceSuite, tokenName string, symbol string, nbDecimal int) (tokenAddr common.Address, tx *types.Transaction, tokenProxyInstance *contracts.TokenProxy, err error) {

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Deploying new token proxy for token [%s]", tokenName)
	tokenAddr, tx, tokenProxyInstance, err = contracts.DeployTokenProxy(
		auth,
		globals.EthClient,
		trexAuthorityImplementationAddr,
		globals.IdentityRegistryAddress,
		complianceSuite.Address,
		tokenName,
		symbol,
		uint8(nbDecimal),
		common.HexToAddress(auth.From.Hex()),
	)

	if err = utils.WaitDeployedTransaction(ctx, tx, true); err != nil {
		return
	}
	logger.LogInfo("📬 Token Proxy deployed at address [%s]", tx.Hash().Hex())
	logger.LogInfo("Token Proxy generate new token at address %s", tokenAddr)

	return
}
