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

func WaitDeployedTransaction(ctx context.Context, tx *types.Transaction, shouldWaitContractReturn bool) (err error) {

	if tx == nil {
		return fmt.Errorf("transaction is nil")
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
		return fmt.Errorf("transaction %s failed with status %d", txHex, receipt.Status)
	}

	logger.LogInfo("✅ Transaction %s mined successfully in block %d", txHex, receipt.BlockNumber.Uint64())

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
			return nil
		} else {
			logger.LogWarn("Attemp %d/%d: Contract code at address %s is empty. Retrying in %s...", attempt, maxRetries, contractAddress.Hex(), retryDelay)
		}

		if attempt < maxRetries {
			select {
			case <-time.After(retryDelay):
				// continue to next attempt
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return fmt.Errorf("contract at address %s not active after %d attempts", contractAddress.Hex(), maxRetries)
}

func WaitTREXSuiteDeployment(ctx context.Context, trexFactoryInstance *contracts.TREXFactory, salt string) (*contracts.TREXFactoryTREXSuiteDeployed, error) {

	logger.LogInfo("⏳ Waiting for TREX Suite deployment to be completed...")

	eventChan := make(chan *contracts.TREXFactoryTREXSuiteDeployed)

	sub, err := trexFactoryInstance.WatchTREXSuiteDeployed(&bind.WatchOpts{Context: ctx}, eventChan, nil, []string{salt})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to TREX Suite deployment events: %v", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case event := <-eventChan:
			return event, nil
		case err := <-sub.Err():
			return nil, fmt.Errorf("error while waiting for TREX Suite deployment event: %v", err)
		case <-time.After(2 * time.Minute):
			return nil, fmt.Errorf("timeout while waiting for TREX Suite deployment event")
		}
	}
}
