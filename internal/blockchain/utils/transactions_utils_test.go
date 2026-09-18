package utils

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

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
