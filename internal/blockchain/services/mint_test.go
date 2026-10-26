package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum"
)

// The guards of SendMintShares run before anything is looked up or signed, and
// what they must guarantee is not just a refusal: they must leave no hash and
// nothing sent. The whole delivery reads "no hash" as "nothing was broadcast",
// so a guard that refused while reporting a hash would be a silent trap.
func TestSendMintSharesRefusesBadInputBeforeSigningAnything(t *testing.T) {
	svc := NewService()
	wallet := "0x1111111111111111111111111111111111111111"
	token := "0x2222222222222222222222222222222222222222"
	record := func(string) error { return nil }

	cases := []struct {
		name   string
		token  string
		wallet string
		shares int
		record func(string) error
	}{
		{"no shares at all", token, wallet, 0, record},
		{"a negative number of shares", token, wallet, -5, record},
		{"a wallet that is not an address", token, "not-an-address", 1, record},
		{"an empty wallet", token, "", 1, record},
		{"no way to record the hash", token, wallet, 1, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txHash, sent, err := svc.SendMintShares(context.Background(), tc.token, tc.wallet, tc.shares, tc.record)

			if err == nil {
				t.Fatal("the mint must be refused")
			}
			if txHash != "" {
				t.Fatalf("a refusal must leave no hash, got %q", txHash)
			}
			if sent {
				t.Fatal("a refusal must send nothing")
			}
		})
	}
}

// MintReceipt is asked whether a mint worked. Answering that on an input it
// cannot check would be worse than refusing: found=true, succeeded=false reads
// as "the chain rejected it", which marks the order failed.
func TestMintReceiptRefusesBadInputWithoutAnsweringAboutTheChain(t *testing.T) {
	svc := NewService()
	hash := "0xab12cd34ef56ab12cd34ef56ab12cd34ef56ab12cd34ef56ab12cd34ef56ab12"

	cases := []struct {
		name   string
		wallet string
		shares int
	}{
		{"a wallet that is not an address", "not-an-address", 1},
		{"an empty wallet", "", 1},
		{"no shares at all", "0x1111111111111111111111111111111111111111", 0},
		{"a negative number of shares", "0x1111111111111111111111111111111111111111", -5},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			found, succeeded, err := svc.MintReceipt(context.Background(), hash, tc.wallet, tc.shares)

			if err == nil {
				t.Fatal("the input must be refused")
			}
			if found || succeeded {
				t.Fatalf("a refused input must not report on the chain, got found=%v succeeded=%v", found, succeeded)
			}
		})
	}
}

// After a failed send, only a node that answers "not found" proves the mint
// never left. Reading any other lookup error as "not sent" would fail the order
// and clear its hash while the mint may still be mined.
func TestMintMayHaveBeenSentOnlyRulesOutWhatTheNodeDoesNotKnow(t *testing.T) {
	cases := []struct {
		name      string
		lookupErr error
		want      bool
	}{
		{"the node knows the transaction", nil, true},
		{"the node does not know it", ethereum.NotFound, false},
		{"a wrapped not found", fmt.Errorf("lookup: %w", ethereum.NotFound), false},
		{"the lookup timed out", context.DeadlineExceeded, true},
		{"the node could not be reached", errors.New("connection refused"), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := mintMayHaveBeenSent(tc.lookupErr); got != tc.want {
				t.Fatalf("mintMayHaveBeenSent(%v) = %v, want %v", tc.lookupErr, got, tc.want)
			}
		})
	}
}
