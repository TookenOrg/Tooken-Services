package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// SendMintShares mints shares whole units of the token at tokenAddress to
// wallet. The transaction is signed first and handed to record, which must
// store its hash; only then is it sent. sent says whether the mint may have
// reached the network: when it is false, nothing was broadcast.
func (s *Service) SendMintShares(ctx context.Context, tokenAddress, wallet string, shares int, record func(txHash string) error) (string, bool, error) {

	if shares <= 0 {
		return "", false, fmt.Errorf("shares must be greater than zero, got %d", shares)
	}
	if !common.IsHexAddress(wallet) {
		return "", false, fmt.Errorf("not a valid wallet address: %s", wallet)
	}
	if record == nil {
		return "", false, errors.New("a record function is required: the hash must be stored before sending")
	}

	// The token says how many decimals one share is worth on chain.
	token, err := database.GetTokenByAddress(ctx, tokenAddress)
	if err != nil {
		return "", false, err
	}

	amount, err := sharesToBaseUnits(strconv.Itoa(shares), int(token.NbDecimal))
	if err != nil {
		return "", false, err
	}

	tokenInstance, err := contracts.NewToken(common.HexToAddress(tokenAddress), globals.EthClient)
	if err != nil {
		return "", false, fmt.Errorf("bind token %s: %w", tokenAddress, err)
	}

	logger.LogInfo("💌 Minting [%d] shares ([%s] base units) of token [%s] for wallet [%s]...",
		shares, amount.String(), tokenAddress, wallet)

	var (
		txHash string
		sent   bool
	)

	_, err = utils.Transact(ctx, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		opts.NoSend = true

		tx, err := tokenInstance.Mint(opts, common.HexToAddress(wallet), amount)
		if err != nil {
			// Estimation or signing refused: no transaction exists, so no hash does.
			return nil, fmt.Errorf("sign mint of %d shares to %s: %w", shares, wallet, err)
		}

		txHash = tx.Hash().Hex()

		if err := record(txHash); err != nil {
			return tx, fmt.Errorf("record mint hash %s before sending: %w", txHash, err)
		}

		if err := globals.EthClient.SendTransaction(opts.Context, tx); err != nil {
			// opts.Context may be the one that just expired: ask with a fresh one.
			lookupCtx, cancel := context.WithTimeout(ctx, mintLookupTimeout)
			defer cancel()

			_, _, byHashErr := globals.EthClient.TransactionByHash(lookupCtx, tx.Hash())
			if !mintMayHaveBeenSent(byHashErr) {
				return tx, fmt.Errorf("send mint of %d shares to %s: %w", shares, wallet, err)
			}

			logger.LogWarn("⚠️ Sending mint [%s] returned an error (%s) and the node cannot rule it out: treating it as sent",
				txHash, err.Error())
		}

		sent = true
		return tx, nil
	})
	if err != nil {
		return txHash, sent, err
	}

	logger.LogInfo("⏳ Mint [%s] sent, waiting for it to be mined", txHash)

	return txHash, true, nil
}

// mintLookupTimeout bounds the question asked after a failed send.
const mintLookupTimeout = 10 * time.Second

// mintMayHaveBeenSent reads the lookup made after a failed send. Only a node
// answering "not found" proves the mint never left; any other error leaves the
// question open, and an open question must count as sent so the order keeps
// its hash instead of being failed.
func mintMayHaveBeenSent(lookupErr error) bool {
	return !errors.Is(lookupErr, ethereum.NotFound)
}

// MintReceipt waits for a mint to be mined and says what the chain made of it.
//
// found is false when the transaction is still not mined once the wait is over.
// That is not a failure: a transaction abandoned here may well be mined a
// minute later, so the caller must leave the order alone rather than mark it
// failed, which would open the door to a second mint.
//
// succeeded is true only when the receipt carries status 1 *and* a
// Transfer(0x0 -> wallet, shares) among its logs. A receipt alone says the
// transaction did not revert; the log says the tokens were actually created
// for this wallet, in this transaction.
func (s *Service) MintReceipt(ctx context.Context, txHash, wallet string, shares int) (bool, bool, error) {

	if !common.IsHexAddress(wallet) {
		return false, false, fmt.Errorf("not a valid wallet address: %s", wallet)
	}
	if shares <= 0 {
		return false, false, fmt.Errorf("shares must be greater than zero, got %d", shares)
	}

	hash := common.HexToHash(txHash)

	waitCtx, cancel := context.WithTimeout(ctx, utils.TxWaitTimeout())
	defer cancel()

	tx, _, err := globals.EthClient.TransactionByHash(waitCtx, hash)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) || waitCtx.Err() != nil {
			// The node never saw it, or we gave up asking: undecided, not failed.
			return false, false, nil
		}
		return false, false, fmt.Errorf("look up mint %s: %w", txHash, err)
	}

	logger.LogInfo("⏳ Waiting up to %s for mint %s to be mined...", utils.TxWaitTimeout(), txHash)

	receipt, err := bind.WaitMined(waitCtx, globals.EthClient, tx)
	if err != nil {
		if waitCtx.Err() != nil {
			logger.LogWarn("⚠️ Mint %s is still not mined: leaving the order as it is", txHash)
			return false, false, nil
		}
		return false, false, fmt.Errorf("wait for mint %s: %w", txHash, err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		logger.LogWarn("⚠️ Mint %s reverted in block %d", txHash, receipt.BlockNumber.Uint64())
		return true, false, nil
	}

	minted, err := mintTransferFound(ctx, receipt, common.HexToAddress(wallet), shares)
	if err != nil {
		return true, false, err
	}
	if !minted {
		logger.LogWarn("⚠️ Mint %s succeeded but carries no Transfer of %d shares to %s", txHash, shares, wallet)
		return true, false, nil
	}

	logger.LogInfo("✅ Mint %s mined in block %d: %d shares are in %s", txHash, receipt.BlockNumber.Uint64(), shares, wallet)

	return true, true, nil
}

// mintTransferFound looks for the proof that this very transaction created the
// shares: a Transfer from the zero address to wallet, for the exact amount.
// The emitting contract names the token, which is where the decimals come from.
func mintTransferFound(ctx context.Context, receipt *types.Receipt, wallet common.Address, shares int) (bool, error) {
	var zero common.Address

	for _, entry := range receipt.Logs {
		if entry == nil {
			continue
		}

		filterer, err := contracts.NewTokenFilterer(entry.Address, globals.EthClient)
		if err != nil {
			continue
		}

		transfer, err := filterer.ParseTransfer(*entry)
		if err != nil {
			continue // not a Transfer event
		}

		if transfer.From != zero || transfer.To != wallet {
			continue
		}

		token, err := database.GetTokenByAddress(ctx, entry.Address.Hex())
		if err != nil {
			continue // a Transfer from a contract that is not one of our tokens
		}

		amount, err := sharesToBaseUnits(strconv.Itoa(shares), int(token.NbDecimal))
		if err != nil {
			return false, err
		}

		if transfer.Value.Cmp(amount) == 0 {
			return true, nil
		}
	}

	return false, nil
}
