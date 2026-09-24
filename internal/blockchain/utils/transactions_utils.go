package utils

import (
	"context"
	"fmt"
	"math/big"
	"time"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateTransactOpts(ctx context.Context) (opts *bind.TransactOpts, err error) {
	privateKey, publicKey := getECDSAKeys()
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	nonce, err := globals.EthClient.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return
	}

	gasPrice, err := globals.EthClient.SuggestGasPrice(ctx)
	if err != nil {
		return
	}

	chainID, err := globals.EthClient.NetworkID(ctx)
	if err != nil {
		return
	}

	opts, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return
	}

	opts.From = fromAddress
	opts.Nonce = big.NewInt(int64(nonce))
	opts.Value = big.NewInt(0)         // in wei
	opts.GasLimit = uint64(15_000_000) // in units
	opts.GasPrice = gasPrice

	return
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

	logger.LogInfo("⏳ Waiting for transaction %s to be mined...", txHex)

	waitCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
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
