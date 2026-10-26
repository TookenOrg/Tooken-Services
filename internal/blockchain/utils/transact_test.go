package utils

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// fakeNode answers the three JSON-RPC calls GenerateTransactOpts makes. While
// frozen is set, it never answers the nonce request, like a node that hangs.
type fakeNode struct {
	frozen atomic.Bool
}

func (n *fakeNode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     json.RawMessage `json:"id"`
		Method string          `json:"method"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result string
	switch req.Method {
	case "eth_getTransactionCount":
		if n.frozen.Load() {
			<-r.Context().Done()
			return
		}
		result = "0x5"
	case "net_version":
		result = "31337"
	case "eth_maxPriorityFeePerGas":
		result = "0x1"
	default:
		http.Error(w, "unexpected method "+req.Method, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": result})
}

// useFakeNode points the package at a fake node and a throwaway key, and
// shortens sendTimeout to timeout, all restored when the test ends.
func useFakeNode(t *testing.T, timeout time.Duration) *fakeNode {
	t.Helper()

	node := &fakeNode{}
	srv := httptest.NewServer(node)
	t.Cleanup(srv.Close)

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(client.Close)

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("PRIVATE_KEY", hex.EncodeToString(crypto.FromECDSA(key)))

	previousClient, previousTimeout := globals.EthClient, sendTimeout
	globals.EthClient, sendTimeout = client, timeout
	t.Cleanup(func() { globals.EthClient, sendTimeout = previousClient, previousTimeout })

	return node
}

func TestTransactSendsOneTransactionAtATime(t *testing.T) {
	useFakeNode(t, 5*time.Second)

	var inside, most atomic.Int32
	var wg sync.WaitGroup

	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := Transact(context.Background(), func(opts *bind.TransactOpts) (*types.Transaction, error) {
				now := inside.Add(1)
				defer inside.Add(-1)
				for {
					seen := most.Load()
					if now <= seen || most.CompareAndSwap(seen, now) {
						break
					}
				}
				// Long enough for the others to pile up behind the lock.
				time.Sleep(20 * time.Millisecond)
				return nil, nil
			})
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()

	if got := most.Load(); got != 1 {
		t.Fatalf("%d transactions were built at the same time, want 1", got)
	}
}

func TestTransactBoundsWhatItHandsOver(t *testing.T) {
	useFakeNode(t, 5*time.Second)

	_, err := Transact(context.Background(), func(opts *bind.TransactOpts) (*types.Transaction, error) {
		if opts.Nonce == nil || opts.Nonce.Uint64() != 5 {
			t.Errorf("nonce = %v, want the node's 5", opts.Nonce)
		}
		deadline, ok := opts.Context.Deadline()
		if !ok {
			t.Fatal("opts.Context has no deadline: bind would wait on a frozen node forever")
		}
		if left := time.Until(deadline); left > 5*time.Second {
			t.Errorf("opts.Context has %s left, want at most the 5s send timeout", left)
		}
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTransactReleasesTheLockWhenTheNodeFreezes(t *testing.T) {
	const timeout = 200 * time.Millisecond

	t.Run("frozen while reading the nonce", func(t *testing.T) {
		node := useFakeNode(t, timeout)
		node.frozen.Store(true)

		start := time.Now()
		_, err := Transact(context.Background(), func(*bind.TransactOpts) (*types.Transaction, error) {
			t.Error("send must not run without options")
			return nil, nil
		})
		if err == nil {
			t.Fatal("want an error from a node that never answers")
		}
		if waited := time.Since(start); waited > 10*timeout {
			t.Fatalf("gave up after %s, want about %s", waited, timeout)
		}

		node.frozen.Store(false)
		expectLockFree(t)
	})

	t.Run("frozen while sending", func(t *testing.T) {
		useFakeNode(t, timeout)

		start := time.Now()
		_, err := Transact(context.Background(), func(opts *bind.TransactOpts) (*types.Transaction, error) {
			// What bind does when SendTransaction never comes back.
			<-opts.Context.Done()
			return nil, opts.Context.Err()
		})
		if err == nil {
			t.Fatal("want the send timeout as an error")
		}
		if waited := time.Since(start); waited > 10*timeout {
			t.Fatalf("gave up after %s, want about %s", waited, timeout)
		}

		expectLockFree(t)
	})
}

// expectLockFree fails when the next transaction cannot get through.
func expectLockFree(t *testing.T) {
	t.Helper()

	done := make(chan error, 1)
	go func() {
		_, err := Transact(context.Background(), func(*bind.TransactOpts) (*types.Transaction, error) {
			return nil, nil
		})
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("next transaction: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the lock is still held: every following transaction is stuck")
	}
}
