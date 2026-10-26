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
//   - a mint towards a wallet that is not verified fails on chain, so the
//     delivery checks the KYC before it signs anything;
//   - a burn does NOT require the holder to be verified, so that check applies
//     to the mint only.
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

	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/ethereum/go-ethereum"
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

	t.Run("an unknown token is refused before anything is signed", func(t *testing.T) {
		unknown := common.HexToAddress("0x00000000000000000000000000000000000B0B0B")

		txHash, sent, err := svc.SendMintShares(ctx, unknown.Hex(), holder.Hex(), 1, func(string) error { return nil })
		if err == nil {
			t.Fatal("an unknown token must be refused")
		}
		// The whole safety of the delivery rests on this: no hash means nothing
		// was ever broadcast.
		if txHash != "" || sent {
			t.Fatalf("nothing must have been signed or sent, got hash %q sent=%v", txHash, sent)
		}
	})

	t.Run("a mint of zero shares is refused", func(t *testing.T) {
		// Asserted so the test above cannot pass through a guard wired to refuse
		// everything for the wrong reason.
		if _, _, err := svc.SendMintShares(ctx, tokenAddr, holder.Hex(), 0, func(string) error { return nil }); err == nil {
			t.Fatal("a mint of zero shares must be refused")
		}
	})

	t.Run("a failing record stops the mint before it is sent", func(t *testing.T) {
		recordErr := errors.New("database is down")

		txHash, sent, err := svc.SendMintShares(ctx, tokenAddr, holder.Hex(), 1, func(string) error { return recordErr })
		if !errors.Is(err, recordErr) {
			t.Fatalf("the record failure must surface, got: %v", err)
		}
		if sent {
			t.Fatal("nothing must be sent when the hash could not be recorded")
		}
		// The hash still comes back: the attempt is real, and the caller may
		// put it on the timeline.
		if txHash == "" {
			t.Fatal("the signed transaction hash must be returned even when the record failed")
		}
		// And the node must never have seen it: a sent=false that had in fact
		// broadcast would be the double mint this rule exists to prevent.
		if _, _, err := fixture.Client.TransactionByHash(ctx, common.HexToHash(txHash)); !errors.Is(err, ethereum.NotFound) {
			t.Fatalf("the node must not know a mint whose hash was not recorded, got: %v", err)
		}
	})

	// ---- mint ------------------------------------------------------------------

	t.Run("Mint credits a verified investor and records the transaction", func(t *testing.T) {
		before, err := suite.Token.BalanceOf(fixture.Call, holder)
		failOnErr(t, err, "BalanceOf before")

		txHash, err := mintSharesForTest(ctx, svc, tokenAddr, holder.Hex(), 1000)
		if err != nil {
			t.Fatalf("Mint: %v", err)
		}
		if txHash == "" {
			t.Fatal("Mint must report the transaction it sent")
		}

		after, err := suite.Token.BalanceOf(fixture.Call, holder)
		failOnErr(t, err, "BalanceOf after")

		// 1000 tokens at 18 decimals: the conversion is the part no unit test covers
		// end to end, and getting it wrong by a factor of 10^18 is invisible on-chain.
		expected := new(big.Int).Mul(big.NewInt(1000), oneToken)
		if gained := new(big.Int).Sub(after, before); gained.Cmp(expected) != 0 {
			t.Fatalf("minted amount: expected %s, got %s", expected, gained)
		}

		// blk.eth_transaction is no longer written for a mint: the hash now
		// lives on the order itself, in iss.issuance_orders.delivery_tx_hash
		// and on its timeline.
	})

	t.Run("Mint towards a wallet that is not verified fails", func(t *testing.T) {
		_, err := mintSharesForTest(ctx, svc, tokenAddr, stranger.Hex(), 1)
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

		// The point of this sub-test is *where* the failure happens. The chain
		// refuses at gas estimation, which is before the transaction is signed,
		// so no hash ever exists for this attempt. That is what lets the
		// delivery read "no hash" as "nothing was broadcast".
		txHash, sent, err := svc.SendMintShares(ctx, tokenAddr, stranger.Hex(), 1, func(string) error { return nil })
		if err == nil {
			t.Fatal("a non-verified recipient must be refused")
		}
		if txHash != "" || sent {
			t.Fatalf("the refusal must happen before signing, got hash %q sent=%v", txHash, sent)
		}
		// Naming the reason is what makes the failure actionable for the operator.
		if !strings.Contains(err.Error(), "Identity is not verified") {
			t.Fatalf("the refusal must name the reason, got: %v", err)
		}
	})

	// ---- what MintReceipt makes of a transaction ---------------------------------
	//
	// found=false is what keeps an order from being failed, so it must only ever
	// mean "not mined yet"; succeeded must require the exact Transfer.

	t.Run("MintReceipt does not find a hash the chain never saw", func(t *testing.T) {
		unknown := "0x" + strings.Repeat("ab", 32)
		found, succeeded, err := svc.MintReceipt(ctx, unknown, holder.Hex(), 1)
		if err != nil || found || succeeded {
			t.Fatalf("an unknown hash is undecided: want (false, false, nil), got (%v, %v, %v)", found, succeeded, err)
		}
	})

	t.Run("MintReceipt requires the Transfer to this wallet, for this amount", func(t *testing.T) {
		txHash, err := mintSharesForTest(ctx, svc, tokenAddr, holder.Hex(), 1)
		if err != nil {
			t.Fatalf("Mint: %v", err)
		}

		for _, tc := range []struct {
			name   string
			wallet string
			shares int
		}{
			{"another wallet", stranger.Hex(), 1},
			{"another amount", holder.Hex(), 2},
		} {
			found, succeeded, err := svc.MintReceipt(ctx, txHash, tc.wallet, tc.shares)
			if err != nil || !found || succeeded {
				t.Fatalf("%s: a mined mint that does not match is a failure: want (true, false, nil), got (%v, %v, %v)",
					tc.name, found, succeeded, err)
			}
		}
	})

	t.Run("MintReceipt reports a mined mint that reverted", func(t *testing.T) {
		// A fixed gas limit skips the estimation that would refuse the
		// non-verified wallet, so the mint is mined and reverts on chain.
		opts, err := utils.GenerateTransactOpts(ctx)
		failOnErr(t, err, "GenerateTransactOpts")
		opts.NoSend = true
		opts.GasLimit = 500_000
		tx, err := suite.Token.Mint(opts, stranger, oneToken)
		failOnErr(t, err, "sign reverting mint")
		// Hardhat reports the revert as a send error, but mines the transaction.
		_ = fixture.Client.SendTransaction(ctx, tx)

		found, succeeded, err := svc.MintReceipt(ctx, tx.Hash().Hex(), stranger.Hex(), 1)
		if err != nil || !found || succeeded {
			t.Fatalf("a reverted mint is a failure: want (true, false, nil), got (%v, %v, %v)", found, succeeded, err)
		}
	})

	t.Run("the mint guard does not leak onto the burn guard", func(t *testing.T) {
		// Same wallet, same moment: refused for a mint, accepted for a burn. Taking
		// shares back from someone who has just been struck off is the whole point of
		// a permissioned security, so the identity check must stay on the mint side.
		if _, _, err := svc.SendMintShares(ctx, tokenAddr, stranger.Hex(), 1, func(string) error { return nil }); err == nil {
			t.Fatal("the mint must refuse a non-verified recipient")
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
		// pre-flight isVerified check must apply to the mint only.
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
		if _, err := mintSharesForTest(ctx, svc, tokenAddr, holder.Hex(), 1); err == nil {
			t.Fatal("minting to a struck-off holder must fail")
		}
	})
}
