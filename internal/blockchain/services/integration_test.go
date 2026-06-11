//go:build integration

package services

import (
	"context"
	"os"
	"testing"
	"time"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Default Hardhat / anvil account #0 (publicly known dev key, no 0x prefix).
const hardhatAccount0Key = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

// TestDeployTransferRestrictModuleIntegration exercises the real on-chain path
// against a local EVM node (Hardhat): build the transact options, deploy a real
// contract through the generated binding, and wait for it via WaitDeployedTransaction.
//
// It is gated behind the `integration` build tag so the default CI `go test ./...`
// never requires a node. To run it:
//
//	cd tools/hardhat && npm install && npx hardhat node    # in one terminal
//	go test -tags integration ./internal/blockchain/...    # in another
//
// Override the endpoint with ETH_TEST_RPC and the signer with PRIVATE_KEY.
func TestDeployTransferRestrictModuleIntegration(t *testing.T) {
	logger.Init(true) // production wires this in main.go; helpers below call logger.LogInfo

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

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		t.Fatalf("GenerateTransactOpts: %v", err)
	}

	addr, tx, _, err := contracts.DeployTransferRestrictModule(auth, globals.EthClient)
	if err != nil {
		t.Fatalf("DeployTransferRestrictModule: %v", err)
	}

	if _, err := utils.WaitDeployedTransaction(ctx, tx, true); err != nil {
		t.Fatalf("WaitDeployedTransaction: %v", err)
	}

	code, err := client.CodeAt(ctx, addr, nil)
	if err != nil {
		t.Fatalf("CodeAt: %v", err)
	}
	if len(code) == 0 {
		t.Fatalf("expected deployed bytecode at %s, got none", addr.Hex())
	}
	t.Logf("TransferRestrictModule deployed at %s (%d bytes of code)", addr.Hex(), len(code))
}
