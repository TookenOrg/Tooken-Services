//go:build integration

// Coverage of Mint and Burn through the production services, against a real chain
// and a real database.
//
// Neither had any integration coverage: only controlInputMintBurn was unit-tested,
// so everything past the input check — the database lookup for the decimals, the
// conversion to wei, the transaction, the eth_transaction trace — ran for the first
// time in production.
//
// Two things are settled here that were open questions:
//
//   - a mint towards a wallet that is not verified fails, and it fails *inside* the
//     goroutine: the caller of POST /contract/token/mint has already been told 202.
//     This is the debt the pre-flight check of TICKET-M2-1 morceau 4 will close;
//   - a burn does NOT require the holder to be verified. That was asserted nowhere,
//     and it matters: the pre-flight check must therefore apply to the mint only.
//
// Requires both a node and a database:
//
//	cd tools/hardhat && npx hardhat node
//	TEST_DATABASE_URL="postgres://…" go test -tags integration ./internal/blockchain/services/...
package services

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	_ "github.com/lib/pq"
)

func TestMintAndBurnIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	env := newChainDBEnv(ctx, t)
	fixture := env.Fixture
	svc := NewService()

	suite := fixture.DeploySuite(t, "MintBurn", "Tooken MintBurn", "TKMB", common.Address{})
	env.SetSharedIRSRow(t, suite.IRSAddr)
	env.InsertTokenRow(t, suite.TokenAddr, "Tooken MintBurn", "TKMB", 18)

	tokenAddr := suite.TokenAddr.Hex()
	oneToken := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

	// A holder who went through the whole journey, and a wallet that did not.
	holder, holderIdentity, holderIdentityInstance := fixture.NewInvestorIdentity(t)
	if _, err := registerIdentity(ctx, holder, holderIdentity, 250); err != nil {
		t.Fatalf("registerIdentity: %v", err)
	}
	fixture.AddKYCClaim(t, holderIdentity, holderIdentityInstance)

	stranger := common.HexToAddress("0x000000000000000000000000000000000000dEaD")

	// ---- the input guard, before anything reaches the chain --------------------

	t.Run("an unknown token is refused as a bad request", func(t *testing.T) {
		unknown := common.HexToAddress("0x00000000000000000000000000000000000B0B0B")

		err := svc.ValidateMintInput(ctx, unknown.Hex(), holder.Hex(), 1)
		if err == nil {
			t.Fatal("an unknown token must be refused")
		}
		// The handler branches on this sentinel to answer 400 rather than 500.
		if !errors.Is(err, ErrInvalidMintInput) {
			t.Fatalf("the refusal must be marked as a bad request, got: %v", err)
		}
	})

	t.Run("a known token with sane inputs passes the guard", func(t *testing.T) {
		// Asserted so the test above cannot pass through a guard wired to refuse
		// everything.
		if err := svc.ValidateMintInput(ctx, tokenAddr, holder.Hex(), 1); err != nil {
			t.Fatalf("a valid mint request must pass the guard, got: %v", err)
		}
	})

	// ---- mint ------------------------------------------------------------------

	t.Run("Mint credits a verified investor and records the transaction", func(t *testing.T) {
		before, err := suite.Token.BalanceOf(fixture.Call, holder)
		failOnErr(t, err, "BalanceOf before")

		result, err := svc.Mint(ctx, tokenAddr, holder.Hex(), 1000)
		if err != nil {
			t.Fatalf("Mint: %v", err)
		}
		if result.OperationName != "Mint" || result.TransactionHash == "" {
			t.Fatalf("Mint must report what it did, got %+v", result)
		}

		after, err := suite.Token.BalanceOf(fixture.Call, holder)
		failOnErr(t, err, "BalanceOf after")

		// 1000 tokens at 18 decimals: the conversion is the part no unit test covers
		// end to end, and getting it wrong by a factor of 10^18 is invisible on-chain.
		expected := new(big.Int).Mul(big.NewInt(1000), oneToken)
		if gained := new(big.Int).Sub(after, before); gained.Cmp(expected) != 0 {
			t.Fatalf("minted amount: expected %s, got %s", expected, gained)
		}

		var count int
		if err := env.DB.QueryRow(
			`SELECT count(*) FROM blk.eth_transaction WHERE tx_hash = $1 AND tx_name = $2`,
			result.TransactionHash, "MINT_TOKEN").Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("the mint must be traced in blk.eth_transaction, found %d rows", count)
		}
	})

	t.Run("Mint towards a wallet that is not verified fails", func(t *testing.T) {
		_, err := svc.Mint(ctx, tokenAddr, stranger.Hex(), 1)
		if err == nil {
			t.Fatal("minting to a non-verified wallet must fail")
		}
		if !strings.Contains(err.Error(), "Identity is not verified") {
			t.Fatalf("the failure must name the reason, got: %v", err)
		}

		balance, err := suite.Token.BalanceOf(fixture.Call, stranger)
		failOnErr(t, err, "BalanceOf stranger")
		if balance.Sign() != 0 {
			t.Fatalf("the non-verified wallet must hold nothing, got %s", balance)
		}

		// The point of this sub-test is *where* the failure happens. Mint runs in a
		// goroutine after the handler has answered 202, so this error never reaches
		// the caller: the pre-flight guard has to catch it first.
		preflight := svc.ValidateMintInput(ctx, tokenAddr, stranger.Hex(), 1)
		if preflight == nil {
			t.Fatal("the pre-flight guard must refuse a non-verified recipient, before the 202")
		}
		// The handler branches on this sentinel to answer 400 rather than 500: an
		// unverified recipient is a bad request, not a server failure.
		if !errors.Is(preflight, ErrInvalidMintInput) {
			t.Fatalf("the refusal must be marked as a bad request, got: %v", preflight)
		}
		// Naming the wallet is what makes the answer actionable for the operator.
		if !strings.Contains(preflight.Error(), stranger.Hex()) {
			t.Fatalf("the refusal must name the wallet, got: %v", preflight)
		}
	})

	t.Run("the mint guard does not leak onto the burn guard", func(t *testing.T) {
		// Same wallet, same moment: refused for a mint, accepted for a burn. Taking
		// shares back from someone who has just been struck off is the whole point of
		// a permissioned security, so the identity check must stay on the mint side.
		if err := svc.ValidateMintInput(ctx, tokenAddr, stranger.Hex(), 1); err == nil {
			t.Fatal("the mint guard must refuse a non-verified recipient")
		}
		if err := svc.ValidateBurnInput(ctx, tokenAddr, stranger.Hex(), 1); err != nil {
			t.Fatalf("the burn guard must not check the identity, got: %v", err)
		}
	})

	// ---- burn ------------------------------------------------------------------

	t.Run("Burn takes shares back and records the transaction", func(t *testing.T) {
		before, err := suite.Token.BalanceOf(fixture.Call, holder)
		failOnErr(t, err, "BalanceOf before burn")

		result, err := svc.Burn(ctx, tokenAddr, holder.Hex(), 400)
		if err != nil {
			t.Fatalf("Burn: %v", err)
		}
		if result.OperationName != "Burn" || result.TransactionHash == "" {
			t.Fatalf("Burn must report what it did, got %+v", result)
		}

		after, err := suite.Token.BalanceOf(fixture.Call, holder)
		failOnErr(t, err, "BalanceOf after burn")

		expected := new(big.Int).Mul(big.NewInt(400), oneToken)
		if lost := new(big.Int).Sub(before, after); lost.Cmp(expected) != 0 {
			t.Fatalf("burned amount: expected %s, got %s", expected, lost)
		}

		var count int
		if err := env.DB.QueryRow(
			`SELECT count(*) FROM blk.eth_transaction WHERE tx_hash = $1 AND tx_name = $2`,
			result.TransactionHash, "BURN_TOKEN").Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("the burn must be traced in blk.eth_transaction, found %d rows", count)
		}
	})

	t.Run("Burn refuses to take more than the holder owns", func(t *testing.T) {
		balance, err := suite.Token.BalanceOf(fixture.Call, holder)
		failOnErr(t, err, "BalanceOf")
		tooMuch := new(big.Int).Div(balance, oneToken).Int64() + 1

		if _, err := svc.Burn(ctx, tokenAddr, holder.Hex(), float64(tooMuch)); err == nil {
			t.Fatal("burning more than the balance must fail")
		}
	})

	// ---- 🔴 the asymmetry that decides where the pre-flight guard applies -------

	t.Run("Burn does not require the holder to be verified", func(t *testing.T) {
		// A permissioned security has to stay recoverable: a court order, a lost key
		// or a fraud discovered after the fact must not be blocked by the fact that
		// the holder has just been struck off the registry.
		//
		// Measured rather than assumed, because it decides a design point: the
		// pre-flight isVerified check of morceau 4 must apply to the mint only.
		tx, err := suite.IR.DeleteIdentity(fixture.Auth, holder)
		failOnErr(t, err, "IR.deleteIdentity")
		fixture.Mine(tx, "IR.deleteIdentity")

		verified, err := suite.IR.IsVerified(fixture.Call, holder)
		failOnErr(t, err, "IsVerified after removal")
		if verified {
			t.Fatal("the holder must no longer be verified, otherwise this proves nothing")
		}

		before, err := suite.Token.BalanceOf(fixture.Call, holder)
		failOnErr(t, err, "BalanceOf before recovery burn")
		if before.Sign() == 0 {
			t.Fatal("the holder must still own shares for this to mean anything")
		}

		if _, err := svc.Burn(ctx, tokenAddr, holder.Hex(), 100); err != nil {
			t.Fatalf("burning from a struck-off holder must still work, got: %v", err)
		}

		after, err := suite.Token.BalanceOf(fixture.Call, holder)
		failOnErr(t, err, "BalanceOf after recovery burn")
		expected := new(big.Int).Mul(big.NewInt(100), oneToken)
		if lost := new(big.Int).Sub(before, after); lost.Cmp(expected) != 0 {
			t.Fatalf("recovered amount: expected %s, got %s", expected, lost)
		}

		// And the counterpart, on the same wallet in the same state: minting to it is
		// refused. Same wallet, same moment, opposite answers — that is the asymmetry.
		if _, err := svc.Mint(ctx, tokenAddr, holder.Hex(), 1); err == nil {
			t.Fatal("minting to a struck-off holder must fail")
		}
	})
}
