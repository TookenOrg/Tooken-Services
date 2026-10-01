//go:build integration

package utils

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Default Hardhat / anvil account #0 (publicly known dev key, no 0x prefix).
const hardhatAccount0Key = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

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
	logger.Init(true)

	rpc := os.Getenv("ETH_TEST_RPC")
	if rpc == "" {
		rpc = "http://127.0.0.1:8545"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpc)
	if err != nil {
		t.Skipf("no local EVM node at %s: %v", rpc, err)
	}
	if _, err := client.ChainID(ctx); err != nil {
		t.Skipf("local EVM node at %s not reachable: %v", rpc, err)
	}

	// Wire the package globals + signer exactly like main.go does at startup.
	globals.EthClient = client
	if os.Getenv("PRIVATE_KEY") == "" {
		t.Setenv("PRIVATE_KEY", hardhatAccount0Key)
	}
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
