package services

import (
	"context"
	"errors"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Service) Mint(ctx context.Context, tokenAddr, to string, humanAmount float64) (txHashName server.TxHashName, err error) {

	tokenInfos, err := database.GetTokenByAddress(ctx, tokenAddr)
	if err != nil {
		return
	}

	ok := controlInputMintBurn(tokenAddr, to, humanAmount, tokenInfos.NbDecimal)
	if !ok {
		err = errors.New("Input data for mint are incorrect")
		return
	}

	amtWei, err := utils.ConvertFloatToWei(humanAmount, tokenInfos.NbDecimal)
	if err != nil {
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

	logger.LogInfo("💌 Minting [%f] tokens (converted to [%s] wei) in Token [%s] for wallet [%s]...", humanAmount, amtWei.String(), tokenAddr, to)
	tx, err := tokenInstance.Mint(auth, common.HexToAddress(to), amtWei)
	if err != nil {
		return
	}

	deployedTxDetails, err := utils.WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return
	}
	logger.LogInfo("📬 [%s] tokens minted on Token contract [%s]", amtWei.String(), tokenAddr)

	err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "MINT_TOKEN", deployedTxDetails.Tx.To().Hex(), deployedTxDetails.BlockNumber.Int64(), *amtWei)
	if err != nil {
		return
	}

	txHashName.OperationName = "Mint"
	txHashName.TransactionHash = tx.Hash().Hex()

	return
}

func controlInputMintBurn(tokenAddr, to string, humanAmount float64, nbDecimal int64) (ok bool) {
	_, err := utils.ConvertFloatToWei(humanAmount, nbDecimal)
	ok = err == nil && common.IsHexAddress(tokenAddr) && common.IsHexAddress(to)
	return
}
