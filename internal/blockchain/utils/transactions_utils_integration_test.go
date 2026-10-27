//go:build integration

package utils

import (
	"context"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/testenv"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

// TestGenerateTransactOptsUsesDynamicFees is a regression guard, not a feature test.
//
// GenerateTransactOpts used to set GasPrice from SuggestGasPrice and GasLimit to a
// flat 15,000,000. Both look harmless and both are traps:
//
//   - a GasPrice makes bind build a *legacy* transaction whose price is frozen at
//     signing time. The base fee moves by up to 12.5% per block, so the transaction
//     drops below the network floor within a few blocks and becomes impossible to
//     include — on 2026-09-28 one sat in the Sepolia mempool for 25 minutes and had
//     to be evicted by hand, freezing the nine-contract install sequence behind it;
//   - a fixed GasLimit skips estimation and reserves a quarter of a Sepolia block
//     for a deployment that consumes 0.9% of it.
//
// Leaving all three fields unset is what makes bind produce a dynamic-fee
// transaction with GasFeeCap = GasTipCap + 2*baseFee and an estimated limit. This
// test exists so that anyone re-adding those lines "to be explicit" sees it fail.
//
// Gated behind the `integration` tag: it needs a node for the nonce, chain ID and
// tip suggestion. To run it:
//
//	cd tools/hardhat && npm install && npx hardhat node    # in one terminal
//	go test -tags integration ./internal/blockchain/...    # in another
func TestGenerateTransactOptsUsesDynamicFees(t *testing.T) {
	ctx := useTestNode(t)
	t.Setenv("ETH_MIN_TIP_GWEI", "2")

	opts, err := GenerateTransactOpts(ctx)
	if err != nil {
		t.Fatalf("GenerateTransactOpts: %v", err)
	}

	if opts.GasPrice != nil {
		t.Errorf("GasPrice must stay nil, got %s — bind would build a legacy transaction", opts.GasPrice)
	}

	if opts.GasFeeCap != nil {
		t.Errorf("GasFeeCap must stay nil so bind derives it from the current base fee, got %s", opts.GasFeeCap)
	}

	if opts.GasLimit != 0 {
		t.Errorf("GasLimit must stay 0 so bind estimates it, got %d", opts.GasLimit)
	}

	if opts.GasTipCap == nil {
		t.Fatal("GasTipCap must be set, otherwise the floor never applies")
	}

	// The floor is 2 gwei here; a local node typically suggests less, so the floor
	// is what must win. A node suggesting more is also correct — never less.
	floor := big.NewInt(2_000_000_000)
	if opts.GasTipCap.Cmp(floor) < 0 {
		t.Errorf("GasTipCap = %s wei, must be at least the %s wei floor", opts.GasTipCap, floor)
	}
}

// useTestNode wires the package globals and signer to the local node, the way
// main.go does at startup, and returns a context bounded for the test.
func useTestNode(t *testing.T) context.Context {
	t.Helper()
	logger.Init(true)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)

	client := testenv.DialNode(ctx, t)

	previous := globals.EthClient
	globals.EthClient = client
	t.Cleanup(func() { globals.EthClient = previous })

	if os.Getenv("PRIVATE_KEY") == "" {
		t.Setenv("PRIVATE_KEY", testenv.HardhatAccount0Key)
	}
	return ctx
}

// TestTransactKeepsConcurrentNoncesApart sends real transactions from many
// goroutines at once, as a mint and a KYC registration would. Each must get a
// nonce of its own: without the shared lock, the node answers "nonce too low".
func TestTransactKeepsConcurrentNoncesApart(t *testing.T) {
	ctx := useTestNode(t)

	chainID, err := globals.EthClient.ChainID(ctx)
	if err != nil {
		t.Fatal(err)
	}

	const senders = 10
	txs := make([]*types.Transaction, senders)
	errs := make([]error, senders)
	var wg sync.WaitGroup

	for i := range senders {
		wg.Add(1)
		go func() {
			defer wg.Done()
			txs[i], errs[i] = Transact(ctx, func(opts *bind.TransactOpts) (*types.Transaction, error) {
				return sendZeroValueToSelf(opts, chainID)
			})
		}()
	}
	wg.Wait()

	nonces := map[uint64]bool{}
	for i := range senders {
		if errs[i] != nil {
			t.Fatalf("sender %d: %v", i, errs[i])
		}
		if nonces[txs[i].Nonce()] {
			t.Fatalf("nonce %d was handed out twice", txs[i].Nonce())
		}
		nonces[txs[i].Nonce()] = true

		receipt, err := bind.WaitMined(ctx, globals.EthClient, txs[i])
		if err != nil {
			t.Fatalf("sender %d: %v", i, err)
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			t.Fatalf("sender %d: status %d", i, receipt.Status)
		}
	}
}

// sendZeroValueToSelf signs and broadcasts the cheapest transaction there is
// with the nonce opts carries.
func sendZeroValueToSelf(opts *bind.TransactOpts, chainID *big.Int) (*types.Transaction, error) {
	head, err := globals.EthClient.HeaderByNumber(opts.Context, nil)
	if err != nil {
		return nil, err
	}
	feeCap := new(big.Int).Add(opts.GasTipCap, new(big.Int).Mul(head.BaseFee, big.NewInt(2)))

	tx, err := opts.Signer(opts.From, types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     opts.Nonce.Uint64(),
		GasTipCap: opts.GasTipCap,
		GasFeeCap: feeCap,
		Gas:       21_000,
		To:        &opts.From,
		Value:     big.NewInt(0),
	}))
	if err != nil {
		return nil, err
	}
	return tx, globals.EthClient.SendTransaction(opts.Context, tx)
}
