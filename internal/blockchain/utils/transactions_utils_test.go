package utils

import (
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TestMinGasTipCap covers the unit conversion and the fallbacks.
//
// The unit is the whole point: ETH_MIN_TIP_GWEI is read in gwei but the protocol
// works in wei, so a missing 10^9 makes the floor a billion times too small and
// the node's suggestion always wins — the function looks fine and does nothing.
//
// The fallbacks matter just as much: the most common case is the variable being
// unset, and a floor of 0 there would silently reinstate the behaviour this code
// exists to prevent.
func TestMinGasTipCap(t *testing.T) {
	const gwei = 1_000_000_000

	tests := []struct {
		name  string
		unset bool
		env   string
		want  int64
	}{
		{name: "unset falls back to the default", unset: true, want: gwei},
		{name: "empty falls back to the default", env: "", want: gwei},
		{name: "whole gwei", env: "2", want: 2 * gwei},

		// Decimal.BigInt truncates towards zero, so scaling has to happen first.
		{name: "half a gwei is not truncated to zero", env: "0.5", want: gwei / 2},
		{name: "the smallest unit the chain knows", env: "0.000000001", want: 1},

		{name: "surrounding spaces are tolerated", env: " 2 ", want: 2 * gwei},

		// Anything unusable falls back rather than failing: a mistyped variable
		// must not stop the server from sending transactions.
		{name: "not a number falls back", env: "abc", want: gwei},
		{name: "negative falls back", env: "-1", want: gwei},
		{name: "zero falls back", env: "0", want: gwei},

		// "1e9" parses cleanly and would mean 1 ETH per unit of gas — one
		// transaction would drain the wallet. The tip is paid in full, unlike the
		// fee cap, so an absurd value is treated as the typo it almost certainly is.
		{name: "absurdly large falls back", env: "1e9", want: gwei},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Setenv registers the restore; Unsetenv then produces the "absent"
			// state without leaking it into the following tests.
			t.Setenv("ETH_MIN_TIP_GWEI", tt.env)
			if tt.unset {
				os.Unsetenv("ETH_MIN_TIP_GWEI")
			}

			got := minGasTipCap()

			// *big.Int must be compared with Cmp: == compares the pointers.
			if want := big.NewInt(tt.want); got.Cmp(want) != 0 {
				t.Fatalf("minGasTipCap() = %s wei, want %s wei", got, want)
			}
		})
	}
}

// TestTxWaitTimeout covers the fallbacks around ETH_TX_WAIT.
//
// The default matters more than it looks: the previous hard-coded two minutes is
// ten blocks on a 12-second chain, which is why a Sepolia deployment was declared
// failed while its transaction was still perfectly alive in the mempool.
func TestTxWaitTimeout(t *testing.T) {
	tests := []struct {
		name  string
		unset bool
		env   string
		want  time.Duration
	}{
		{name: "unset falls back to the default", unset: true, want: 5 * time.Minute},
		{name: "empty falls back to the default", env: "", want: 5 * time.Minute},
		{name: "minutes", env: "10m", want: 10 * time.Minute},
		{name: "seconds", env: "90s", want: 90 * time.Second},
		{name: "composite duration", env: "1m30s", want: 90 * time.Second},
		{name: "surrounding spaces are tolerated", env: " 90s ", want: 90 * time.Second},

		// A bare number is not a Go duration: "300" is rejected, not read as
		// seconds. Falling back beats waiting 300 nanoseconds.
		{name: "a unitless number falls back", env: "300", want: 5 * time.Minute},
		{name: "plain words fall back", env: "cinq minutes", want: 5 * time.Minute},
		{name: "zero falls back", env: "0s", want: 5 * time.Minute},
		{name: "negative falls back", env: "-1m", want: 5 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ETH_TX_WAIT", tt.env)
			if tt.unset {
				os.Unsetenv("ETH_TX_WAIT")
			}

			if got := txWaitTimeout(); got != tt.want {
				t.Fatalf("txWaitTimeout() = %s, want %s", got, tt.want)
			}
		})
	}
}

// TestTxDetailsToAddressHex covers the distinction between a call and a deployment.
//
// A contract creation carries no recipient: the protocol encodes that as a nil `to`,
// and go-ethereum returns a nil *common.Address for it. Reading .Hex() on that nil is
// a panic, and it sat on every deploy path of this package — POST /contract/identity
// crashed there, after the ONCHAINID had already been deployed, which is the worst
// possible moment: the chain had moved, the database had not.
//
// The deployment case is therefore the one that matters; the call case is asserted
// too, because a fix that returned the created address for everything would lose the
// recipient of every ordinary transaction while looking just as green.
func TestTxDetailsToAddressHex(t *testing.T) {
	created := common.HexToAddress("0xAAAAaaaAAaAAaAaaAAAAaAAAaAaAaaAaaAaAAAAa")
	recipient := common.HexToAddress("0xBBbBBbbBbbBbBBbBbBBbbbBbBBbbBbbBBBbBBbbb")

	t.Run("a deployment answers with the address it created", func(t *testing.T) {
		details := txDetails{
			// types.NewContractCreation builds a transaction with no recipient,
			// exactly like the bindings' Deploy* helpers do.
			Tx:              types.NewContractCreation(0, big.NewInt(0), 21000, big.NewInt(1), nil),
			ContractAddress: created,
		}

		if got := details.ToAddressHex(); got != created.Hex() {
			t.Fatalf("a deployment must report the created contract, got %s want %s", got, created.Hex())
		}
	})

	t.Run("a call answers with its recipient", func(t *testing.T) {
		details := txDetails{
			Tx: types.NewTransaction(0, recipient, big.NewInt(0), 21000, big.NewInt(1), nil),
			// A receipt for a call carries the zero address here; it must be ignored.
			ContractAddress: common.Address{},
		}

		if got := details.ToAddressHex(); got != recipient.Hex() {
			t.Fatalf("a call must report its recipient, got %s want %s", got, recipient.Hex())
		}
	})
}
