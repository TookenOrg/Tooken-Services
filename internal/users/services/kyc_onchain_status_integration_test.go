//go:build integration

package services

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	blkGlobals "github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/testenv"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestKYCOnChainStatusReadsSharedRegistryWithoutSigner(t *testing.T) {
	db := testenv.ScratchDatabase(t, "tooken_it_user_kyc_status")
	oldDB, oldClient := globals.DB, blkGlobals.EthClient
	globals.DB = db
	t.Cleanup(func() {
		globals.DB = oldDB
		blkGlobals.EthClient = oldClient
	})
	logger.Init(false)
	t.Setenv("PRIVATE_KEY", "")
	storage := common.HexToAddress("0x0000000000000000000000000000000000000001")
	registry := common.HexToAddress("0x0000000000000000000000000000000000000002")
	wallet := common.HexToAddress("0x0000000000000000000000000000000000000003")
	var userID int
	if err := db.QueryRow(`INSERT INTO usr.users (full_name, email, password)
		VALUES ('Read-only KYC', 'readonly@example.test', 'x') RETURNING id`).Scan(&userID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO blk.user_wallet (user_id, wallet_address, label)
		VALUES ($1, $2, 'Main')`, userID, wallet.Hex()); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO blk.contract_role (tx_hash, address, contract_name)
		VALUES ($1, $2, $3)`, "0x"+strings.Repeat("0", 64), storage.Hex(), blkGlobals.SharedIdentityRegistryStorageName); err != nil {
		t.Fatal(err)
	}
	storageABI, err := contracts.IdentityRegistryStorageMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}
	registryABI, err := contracts.IdentityRegistryMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}
	var verified, rpcFailure, noRegistries atomic.Bool
	rpcServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			ID     json.RawMessage   `json:"id"`
			Method string            `json:"method"`
			Params []json.RawMessage `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if request.Method != "eth_call" || len(request.Params) == 0 {
			t.Errorf("unexpected RPC method %q", request.Method)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var call struct {
			To    common.Address `json:"to"`
			Input string         `json:"input"`
			Data  string         `json:"data"`
		}
		if err := json.Unmarshal(request.Params[0], &call); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if rpcFailure.Load() {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"jsonrpc": "2.0", "id": request.ID,
				"error": map[string]any{"code": -32000, "message": "RPC unavailable"},
			})
			return
		}
		var output []byte
		var encodeErr error
		switch call.To {
		case storage:
			linked := []common.Address{registry}
			if noRegistries.Load() {
				linked = nil
			}
			output, encodeErr = storageABI.Methods["linkedIdentityRegistries"].Outputs.Pack(linked)
		case registry:
			input := call.Input
			if input == "" {
				input = call.Data
			}
			expected, packErr := registryABI.Pack("isVerified", wallet)
			if packErr != nil || input != "0x"+hex.EncodeToString(expected) {
				t.Errorf("unexpected registry call input %s", input)
			}
			output, encodeErr = registryABI.Methods["isVerified"].Outputs.Pack(verified.Load())
		default:
			t.Errorf("unexpected contract %s", call.To)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if encodeErr != nil {
			t.Error(encodeErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"jsonrpc": "2.0", "id": request.ID, "result": "0x" + hex.EncodeToString(output),
		})
	}))
	defer rpcServer.Close()
	client, err := ethclient.DialContext(t.Context(), rpcServer.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	blkGlobals.EthClient = client
	service := &Service{}
	for _, value := range []bool{false, true} {
		verified.Store(value)
		status, err := service.GetCurrentUserKycOnChain(t.Context(), userID)
		if err != nil {
			t.Fatal(err)
		}
		if status.IsVerified != value || status.WalletAddress == nil || *status.WalletAddress != wallet.Hex() {
			t.Fatalf("status = %+v, want verified=%v", status, value)
		}
	}
	noRegistries.Store(true)
	if _, err := service.GetCurrentUserKycOnChain(t.Context(), userID); !errors.Is(err, ErrKYCOnChainUnavailable) {
		t.Fatalf("missing registry error = %v", err)
	}
	noRegistries.Store(false)
	rpcFailure.Store(true)
	if _, err := service.GetCurrentUserKycOnChain(t.Context(), userID); !errors.Is(err, ErrKYCOnChainUnavailable) {
		t.Fatalf("RPC error = %v", err)
	}
	if _, err := db.Exec(`DELETE FROM blk.contract_role WHERE contract_name = $1`,
		blkGlobals.SharedIdentityRegistryStorageName); err != nil {
		t.Fatal(err)
	}
	if _, err := service.GetCurrentUserKycOnChain(t.Context(), userID); !errors.Is(err, ErrKYCOnChainUnavailable) {
		t.Fatalf("missing shared storage error = %v", err)
	}
}
