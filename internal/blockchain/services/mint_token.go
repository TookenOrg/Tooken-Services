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
	amtWei, err := utils.ConvertFloatToWei(humanAmount, tokenInfos.NbDecimal)
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

	err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "MINT_TOKEN", deployedTxDetails.ToAddressHex(), deployedTxDetails.BlockNumber.Int64(), *amtWei)
	if err != nil {
		return
	}

	txHashName.OperationName = "Mint"
	txHashName.TransactionHash = tx.Hash().Hex()

	return
}

// ErrInvalidMintInput marks a refusal caused by the request itself, as
// opposed to a server-side failure. It is what lets the handler answer 400
// rather than 500 without inspecting error strings.
var ErrInvalidMintInput = errors.New("invalid mint input")

// ValidateMintInput refuses a mint request before the handler accepts the job.
//
// MintTokenAsync answers 202 and finishes in a goroutine, so this is the last moment
// where the caller can still be told the request is wrong: past it, a refusal is only
// a line in a log nobody reads, behind a "started successfully" response.
//
// It costs one database read (the decimals) and two eth_call (the token's registry,
// then isVerified). None of it is a transaction, but it is not free either: it is a
// per-request check, not something to run in a loop.
//
// Two kinds of failure come out of here, and the difference is what the handler
// branches on:
//
//   - a refusal the caller can act on — unknown token, bad amount, recipient not
//     verified — is wrapped in ErrInvalidMintInput and becomes a 400;
//   - a failure on our side — database down, node unreachable — is returned as is,
//     with its cause, and becomes a 500. Reporting it as a bad request would send the
//     caller fixing a payload that was never the problem.
//
// The identity check lives here and deliberately NOT in ValidateBurnInput: taking
// shares back from a holder who has just been struck off the registry is legitimate,
// and is what a permissioned security exists for. Verified by integration test —
// same wallet, same moment, mint refused and burn accepted.
func (s *Service) ValidateMintInput(ctx context.Context, tokenAddr, to string, humanAmount float64) error {

	logger.LogInfo("🔍 Validating mint input for token [%s] and wallet [%s] with amount [%f]...", tokenAddr, to, humanAmount)

	tokenInfos, err := database.GetTokenByAddress(ctx, tokenAddr)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: unknown token contract address %q", ErrInvalidMintInput, tokenAddr)
	}

	// Any other database error is a server failure, not a bad request: it must
	// not be reported to the caller as something they can fix.
	if err != nil {
		return err
	}

	if _, err := controlInputMint(tokenAddr, to, humanAmount, tokenInfos.NbDecimal); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidMintInput, err)
	}

	// Asked before the 202 rather than left to the chain: minting to an unverified
	// wallet reverts, and that revert would happen after the caller has been told the
	// operation started.
	isVerified, err := fetchIsVerifiedByTokenAddress(ctx, tokenAddr, to)
	if err != nil {
		return fmt.Errorf("could not check whether %s is verified for token %s: %w", to, tokenAddr, err)
	}
	if !isVerified {
		return fmt.Errorf("%w: wallet address %q is not verified", ErrInvalidMintInput, to)
	}

	return nil
}

// controlInputMint validates the inputs and returns the converted amount,
// so that the conversion happens exactly once. It returns the reason of a
// refusal rather than a bare boolean: mint and burn answer 202 before running,
// so the log is the only place where that reason can still be read.
func controlInputMint(tokenAddr, to string, humanAmount float64, nbDecimal int64) (amtWei *big.Int, err error) {
	if !common.IsHexAddress(tokenAddr) {
		return nil, fmt.Errorf("invalid token contract address %q", tokenAddr)
	}

	if !common.IsHexAddress(to) {
		return nil, fmt.Errorf("invalid wallet address %q", to)
	}

	return utils.ConvertFloatToWei(humanAmount, nbDecimal)
}
