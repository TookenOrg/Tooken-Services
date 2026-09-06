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

func (s *Service) Burn(ctx context.Context, tokenAddr, to string, humanAmount float64) (txHashName server.TxHashName, err error) {

	tokenInfos, err := database.GetTokenByAddress(ctx, tokenAddr)
	if err != nil {
		return
	}

	amtWei, err := controlInputMintBurn(tokenAddr, to, humanAmount, tokenInfos.NbDecimal)
	if err != nil {
		err = fmt.Errorf("input data for burn are incorrect: %w", err)
		return
	}

	tokenInstance, err := contracts.NewToken(common.HexToAddress(tokenAddr), globals.EthClient)
	if err != nil {
		return
	}

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Burning [%f] tokens (converted to [%s] wei) in Token [%s] for wallet [%s]...", humanAmount, amtWei.String(), tokenAddr, to)
	tx, err := tokenInstance.Burn(auth, common.HexToAddress(to), amtWei)
	if err != nil {
		return
	}

	deployedTxDetails, err := utils.WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return
	}
	err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "BURN_TOKEN", deployedTxDetails.Tx.To().Hex(), deployedTxDetails.BlockNumber.Int64(), *amtWei)
	if err != nil {
		return
	}
	logger.LogInfo("📬 [%s] tokens burned on Token contract [%s]", amtWei.String(), tokenAddr)

	txHashName.OperationName = "Burn"
	txHashName.TransactionHash = tx.Hash().Hex()

	return
}
