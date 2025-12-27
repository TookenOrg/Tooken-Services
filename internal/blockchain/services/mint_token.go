package services

import (
	"context"
	"errors"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Service) Mint(ctx context.Context, tokenAddr, to string, humanAmount float64) (txHashName server.TxHashName, err error) {

	ok := controlInputMint(tokenAddr, to, humanAmount)
	if !ok {
		err = errors.New("Input data for mint are incorrect")
		return
	}

	amtWei, err := utils.ConvertFloatToWei(humanAmount, 0)
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

	logger.LogInfo("💌 Minting [%f] tokens (converted to [%s] wei) in Token [%s]...", humanAmount, amtWei.String(), tokenAddr)
	tx, err := tokenInstance.Mint(auth, common.HexToAddress(to), amtWei)
	if err != nil {
		return
	}

	err = utils.WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return
	}
	logger.LogInfo("📬 [%s] tokens minted on Token contract [%s]", amtWei.String(), tokenAddr)

	txHashName.OperationName = "Mint"
	txHashName.TransactionHash = tx.Hash().Hex()

	return
}

func controlInputMint(tokenAddr, to string, humanAmount float64) (ok bool) {
	_, err := utils.ConvertFloatToWei(humanAmount, 0)
	return common.IsHexAddress(tokenAddr) && common.IsHexAddress(to) && err != nil
}
