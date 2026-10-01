package utils

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
)

func GenerateTransactOpts(ctx context.Context) (opts *bind.TransactOpts, err error) {
	privateKey, publicKey := getECDSAKeys()
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	nonce, err := globals.EthClient.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return
	}

	chainID, err := globals.EthClient.NetworkID(ctx)
	if err != nil {
		return
	}

	gasTipCap, err := globals.EthClient.SuggestGasTipCap(ctx)
	if err != nil {
		return
	}

	opts, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return
	}

	// The node's suggestion tracks recent blocks, but a chain whose blocks are
	// half empty suggests a tip near zero — 0.001 gwei on Sepolia. That buys no
	// priority the moment blocks fill up, so keep whichever value is higher.
	if floor := minGasTipCap(); gasTipCap.Cmp(floor) < 0 {
		gasTipCap = floor
	}

	opts.From = fromAddress
	opts.Nonce = big.NewInt(int64(nonce))
	opts.Value = big.NewInt(0) // in wei
	opts.GasTipCap = gasTipCap

	// GasPrice, GasFeeCap and GasLimit are deliberately left unset.
	//
	// Setting GasPrice makes bind build a legacy transaction whose price is frozen
	// at signing time. The base fee moves by up to 12.5% per block, so such a
	// transaction drops below the network floor within a few blocks and can no
	// longer enter any block — it is not slow, it is forbidden. That is what
	// stalled the Sepolia install sequence on 2026-09-28.
	//
	// Left nil, bind builds a dynamic-fee transaction with
	// GasFeeCap = GasTipCap + 2*baseFee, which tolerates a doubling of the base
	// fee, and estimates GasLimit instead of reserving a fixed 15M. The cap is a
	// ceiling, not a price: only baseFee + tip is ever actually paid.

	return
}

const (
	// defaultMinTipGwei is the tip floor, in gwei, applied when ETH_MIN_TIP_GWEI
	// is unset or unusable.
	defaultMinTipGwei = 1

	// maxMinTipGwei rejects values that can only be a typo. 1000 gwei is already
	// an order of magnitude above mainnet peaks, and unlike the fee cap — which is
	// merely a ceiling — the tip is paid in full on every single transaction.
	maxMinTipGwei = 1000
)

// weiPerGwei scales the unit operators think in (gwei) to the unit the protocol
// works in (wei).
var weiPerGwei = decimal.New(1, 9)

// minGasTipCap returns the floor applied to the node's suggested tip, in wei.
//
// It reads ETH_MIN_TIP_GWEI, a value expressed in gwei that may be fractional
// ("1", "0.5", "2.5"). Any unusable value falls back to the default instead of
// failing: a mistyped environment variable must not stop the server from sending
// transactions. The fallback is always logged, so it never passes unnoticed.
func minGasTipCap() *big.Int {
	fallback := decimal.NewFromInt(defaultMinTipGwei).Mul(weiPerGwei).BigInt()

	raw := strings.TrimSpace(os.Getenv("ETH_MIN_TIP_GWEI"))
	if raw == "" {
		return fallback
	}

	gwei, err := decimal.NewFromString(raw)
	if err != nil {
		logger.LogWarn("⚠️ ETH_MIN_TIP_GWEI=%q is not a number, using %d gwei", raw, defaultMinTipGwei)
		return fallback
	}

	if !gwei.IsPositive() {
		logger.LogWarn("⚠️ ETH_MIN_TIP_GWEI=%q must be strictly positive, using %d gwei", raw, defaultMinTipGwei)
		return fallback
	}

	if gwei.GreaterThan(decimal.NewFromInt(maxMinTipGwei)) {
		logger.LogWarn("⚠️ ETH_MIN_TIP_GWEI=%q exceeds the %d gwei ceiling, using %d gwei", raw, maxMinTipGwei, defaultMinTipGwei)
		return fallback
	}

	// Scale before converting: Decimal.BigInt truncates towards zero, so
	// converting "0.5" first would yield 0 instead of 500000000 wei.
	return gwei.Mul(weiPerGwei).BigInt()
}

// defaultTxWait is how long to wait for a transaction to be mined when
// ETH_TX_WAIT is unset or unusable.
//
// The previous value was two minutes, which is ten blocks on a 12-second chain —
// enough on a local node that mines instantly, too little on a public network
// where a transaction may sit through a burst of congestion before landing.
const defaultTxWait = 5 * time.Minute

// txWaitTimeout returns how long to wait for a transaction to be mined.
//
// It reads ETH_TX_WAIT as a Go duration ("5m", "90s"). Like minGasTipCap, any
// unusable value falls back to the default and says so, rather than failing.
func txWaitTimeout() time.Duration {
	raw := strings.TrimSpace(os.Getenv("ETH_TX_WAIT"))
	if raw == "" {
		return defaultTxWait
	}

	wait, err := time.ParseDuration(raw)
	if err != nil {
		logger.LogWarn("⚠️ ETH_TX_WAIT=%q is not a duration, using %s", raw, defaultTxWait)
		return defaultTxWait
	}

	if wait <= 0 {
		logger.LogWarn("⚠️ ETH_TX_WAIT=%q must be strictly positive, using %s", raw, defaultTxWait)
		return defaultTxWait
	}

	return wait
}

type txDetails struct {
	Tx              *types.Transaction
	BlockNumber     big.Int
	ReceiptStatus   uint64
	ContractAddress common.Address // set by the receipt when the transaction created a contract
}

// ToAddressHex is the address the transaction acted on: its recipient for a call, and
// the address it created for a deployment.
//
// Transaction.To() is nil for a contract creation — that is how the protocol encodes
// "no recipient" — so calling .Hex() on it panics on every deploy path. Persisting the
// created address is also the more useful answer: an eth_transaction row pointing at
// nothing says nothing.
func (d txDetails) ToAddressHex() string {
	if to := d.Tx.To(); to != nil {
		return to.Hex()
	}
	return d.ContractAddress.Hex()
}

func WaitDeployedTransaction(ctx context.Context, tx *types.Transaction, shouldWaitContractReturn bool) (txDetails txDetails, err error) {

	if tx == nil {
		return txDetails, fmt.Errorf("transaction is nil")
	}

	txHex := tx.Hash().Hex()

	wait := txWaitTimeout()

	logger.LogInfo("⏳ Waiting up to %s for transaction %s to be mined...", wait, txHex)

	waitCtx, cancel := context.WithTimeout(ctx, wait)
	defer cancel()

	receipt, err := bind.WaitMined(waitCtx, globals.EthClient, tx)
	if err != nil {
		return
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		err = fmt.Errorf("transaction %s failed with status %d", txHex, receipt.Status)
		return
	}

	logger.LogInfo("✅ Transaction %s mined successfully in block %d", txHex, receipt.BlockNumber.Uint64())

	txDetails.Tx = tx
	txDetails.BlockNumber = *receipt.BlockNumber
	txDetails.ReceiptStatus = receipt.Status
	txDetails.ContractAddress = receipt.ContractAddress

	if !shouldWaitContractReturn {
		return
	}

	logger.LogInfo("⏳ Waiting for contract execution result of transaction %s...", txHex)

	contractAddress := receipt.ContractAddress
	maxRetries := 10
	retryDelay := 1 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		code, err := globals.EthClient.CodeAt(ctx, contractAddress, nil)
		if err != nil {
			logger.LogWarn("Attemp %d/%d: Failed to get contract code at address %s: %v", attempt, maxRetries, contractAddress.Hex(), err)
		} else if len(code) > 0 {
			logger.LogInfo("✅ Contract at address %s is successfully deployed and active.", contractAddress.Hex())
			return txDetails, nil
		} else {
			logger.LogWarn("Attemp %d/%d: Contract code at address %s is empty. Retrying in %s...", attempt, maxRetries, contractAddress.Hex(), retryDelay)
		}

		if attempt < maxRetries {
			select {
			case <-time.After(retryDelay):
				// continue to next attempt
			case <-ctx.Done():
				return txDetails, ctx.Err()
			}
		}
	}
	return txDetails, fmt.Errorf("contract at address %s not active after %d attempts", contractAddress.Hex(), maxRetries)
}

// WaitTREXSuiteDeployment returns the TREXSuiteDeployed event emitted by a suite
// deployment transaction.
//
// It waits for the receipt, then reads the logs of the block that transaction landed
// in. The earlier version subscribed with WatchTREXSuiteDeployed instead, which only
// ever reports *future* logs: between sending the transaction and subscribing, the
// node may already have mined it, and the event was then missed for good — two
// minutes of waiting followed by a timeout, for a deployment that had succeeded.
//
// The receipt is what makes this reliable: it proves the transaction is mined and
// says in which block, so the log cannot be missed whatever the pace of the chain.
func WaitTREXSuiteDeployment(ctx context.Context, trexFactoryInstance *contracts.TREXFactory, tx *types.Transaction, salt string) (*contracts.TREXFactoryTREXSuiteDeployed, error) {

	logger.LogInfo("⏳ Waiting for TREX Suite deployment to be completed...")

	txDetails, err := WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return nil, err
	}

	// A single block: the one the receipt names. Leaving End at nil would scan up to
	// the head of the chain for a log we already know the location of.
	blockNumber := txDetails.BlockNumber.Uint64()

	iterator, err := trexFactoryInstance.FilterTREXSuiteDeployed(
		&bind.FilterOpts{Start: blockNumber, End: &blockNumber, Context: ctx},
		nil,
		[]string{salt},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read the TREX Suite deployment logs of block %d: %w", blockNumber, err)
	}
	defer iterator.Close()

	if !iterator.Next() {
		if err := iterator.Error(); err != nil {
			return nil, fmt.Errorf("failed to read the TREX Suite deployment logs of block %d: %w", blockNumber, err)
		}
		// The transaction succeeded but emitted nothing for this salt. Saying so beats
		// answering a zero-valued event that the caller would dereference.
		return nil, fmt.Errorf("no TREX Suite deployment event for salt %q in block %d", salt, blockNumber)
	}

	logger.LogInfo("📬 TREX Suite deployed in block %d", blockNumber)

	return iterator.Event, nil
}
