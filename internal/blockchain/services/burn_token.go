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

func (s *Service) Burn(ctx context.Context, tokenAddr, to string, humanAmount float64) (txHashName server.TxHashName, err error) {

	tokenInfos, err := database.GetTokenByAddress(ctx, tokenAddr)
	if err != nil {
		return
	}

	amtWei, err := utils.ConvertFloatToWei(humanAmount, tokenInfos.NbDecimal)
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
	err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "BURN_TOKEN", deployedTxDetails.ToAddressHex(), deployedTxDetails.BlockNumber.Int64(), *amtWei)
	if err != nil {
		return
	}
	logger.LogInfo("📬 [%s] tokens burned on Token contract [%s]", amtWei.String(), tokenAddr)

	txHashName.OperationName = "Burn"
	txHashName.TransactionHash = tx.Hash().Hex()

	return
}

// ErrInvalidBurnInput marks a refusal caused by the request itself, as
// opposed to a server-side failure. It is what lets the handler answer 400
// rather than 500 without inspecting error strings.
var ErrInvalidBurnInput = errors.New("invalid burn input")

// ValidateBurnInput refuses a burn request before the handler accepts the job.
//
// BurnTokenAsync answers 202 and finishes in a goroutine, so this is the last moment
// where the caller can still be told the request is wrong.
//
// Unlike ValidateMintInput it does not touch the chain, and that is deliberate: a
// burn does not require the holder to be verified. Recovering shares from someone
// struck off the registry — a court order, a lost key, a fraud found after the fact —
// is exactly what an agent's burn power is for. Adding an identity check here would
// break recovery while trying to secure issuance.
func (s *Service) ValidateBurnInput(ctx context.Context, tokenAddr, to string, humanAmount float64) error {

	logger.LogInfo("🔍 Validating burn input for token [%s] and wallet [%s] with amount [%f]...", tokenAddr, to, humanAmount)

	tokenInfos, err := database.GetTokenByAddress(ctx, tokenAddr)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: unknown token contract address %q", ErrInvalidBurnInput, tokenAddr)
	}

	// Any other database error is a server failure, not a bad request: it must
	// not be reported to the caller as something they can fix.
	if err != nil {
		return err
	}

	if _, err := controlInputBurn(tokenAddr, to, humanAmount, tokenInfos.NbDecimal); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidBurnInput, err)
	}

	return nil
}

// controlInputBurn validates the inputs and returns the converted amount,
// so that the conversion happens exactly once. It returns the reason of a
// refusal rather than a bare boolean: mint and burn answer 202 before running,
// so the log is the only place where that reason can still be read.
func controlInputBurn(tokenAddr, to string, humanAmount float64, nbDecimal int64) (amtWei *big.Int, err error) {
	if !common.IsHexAddress(tokenAddr) {
		return nil, fmt.Errorf("invalid token contract address %q", tokenAddr)
	}

	if !common.IsHexAddress(to) {
		return nil, fmt.Errorf("invalid wallet address %q", to)
	}

	return utils.ConvertFloatToWei(humanAmount, nbDecimal)
}
