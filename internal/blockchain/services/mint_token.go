package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"

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

	amtWei, err := controlInputMintBurn(tokenAddr, to, humanAmount, tokenInfos.NbDecimal)
	if err != nil {
		err = fmt.Errorf("input data for mint are incorrect: %w", err)
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

// ErrInvalidMintBurnInput marks a refusal caused by the request itself, as
// opposed to a server-side failure. It is what lets the handler answer 400
// rather than 500 without inspecting error strings.
var ErrInvalidMintBurnInput = errors.New("invalid mint or burn input")

// ValidateMintBurnInput checks a mint or burn request without touching the
// chain, so that the handler can answer 400 before accepting the job.
//
// Mint and burn answer 202 and finish in a goroutine: without this, a refused
// amount produced a "started successfully" response and a line in the log that
// nobody reads. The caller would believe the operation was under way.
func (s *Service) ValidateMintBurnInput(ctx context.Context, tokenAddr, to string, humanAmount float64) error {
	tokenInfos, err := database.GetTokenByAddress(ctx, tokenAddr)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: unknown token contract address %q", ErrInvalidMintBurnInput, tokenAddr)
	}

	// Any other database error is a server failure, not a bad request: it must
	// not be reported to the caller as something they can fix.
	if err != nil {
		return err
	}

	if _, err := controlInputMintBurn(tokenAddr, to, humanAmount, tokenInfos.NbDecimal); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidMintBurnInput, err)
	}

	return nil
}

// controlInputMintBurn validates the inputs and returns the converted amount,
// so that the conversion happens exactly once. It returns the reason of a
// refusal rather than a bare boolean: mint and burn answer 202 before running,
// so the log is the only place where that reason can still be read.
func controlInputMintBurn(tokenAddr, to string, humanAmount float64, nbDecimal int64) (amtWei *big.Int, err error) {
	if !common.IsHexAddress(tokenAddr) {
		return nil, fmt.Errorf("invalid token contract address %q", tokenAddr)
	}

	if !common.IsHexAddress(to) {
		return nil, fmt.Errorf("invalid wallet address %q", to)
	}

	return utils.ConvertFloatToWei(humanAmount, nbDecimal)
}
