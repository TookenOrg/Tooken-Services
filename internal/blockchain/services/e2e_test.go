//go:build integration

package services

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestTokenLifecycleIntegration reproduces the full ERC-3643 / T-REX flow on a local
// EVM (Hardhat), directly through the contract bindings (no DB), to prove that:
//
//  1. the T-REX infrastructure can be deployed and configured,
//  2. a security token can be issued by the shared factory,
//  3. a KYC'd investor identity can be registered, and
//  4. that investor wallet can then RECEIVE shares of the token (mint succeeds and
//     the balance reflects it) — i.e. the permissioning chain works end to end.
//
// It registers the investor through IdentityRegistry.registerIdentity, the path the
// T-REX documentation describes. The path production takes — writing into the shared
// storage — is covered by TestRegisterIdentityIntegration, which also goes through
// the service and the database. Keeping both means a divergence between the protocol
// and our use of it cannot pass unnoticed.
//
// Gated by the `integration` build tag. Requires a local node:
//
//	cd tools/hardhat && npm install && npx hardhat node
//	go test -tags integration -run TestTokenLifecycleIntegration -v ./internal/blockchain/...
func TestTokenLifecycleIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// ---- Steps 1 to 5: the T-REX infrastructure ----------------------------
	t.Log("STEP 1-5 — deploying the T-REX infrastructure (implementations, ONCHAINID, authority, factory, claim issuer)")
	fixture := newTREXFixture(ctx, t)

	// ---- Step 6: deploy the token suite via the factory ---------------------
	t.Log("STEP 6 — deploying the token suite via deployTREXSuite (ONCHAINID=0, KYC topic 7)")
	suite := fixture.DeploySuite(t, "E2ETokenLifecycle", "Tooken Property #1", "TKP1", common.Address{})
	t.Logf("  token=%s  identityRegistry=%s", suite.TokenAddr.Hex(), suite.IRAddr.Hex())

	// ---- Step 7: a KYC'd investor identity, registered in the IR ------------
	t.Log("STEP 7 — creating an investor identity, registering it and adding a KYC claim")
	investor, investorIdAddr, investorID := fixture.NewInvestorIdentity(t)

	// register the wallet -> identity link (deployer is an IR agent).
	tx, err := suite.IR.RegisterIdentity(fixture.Auth, investor, investorIdAddr, countryCode)
	failOnErr(t, err, "RegisterIdentity")
	fixture.Mine(tx, "RegisterIdentity")

	fixture.AddKYCClaim(t, investorIdAddr, investorID)

	// ---- Step 8: the investor must be eligible ------------------------------
	t.Log("STEP 8 — checking the investor is verified")
	verified, err := suite.IR.IsVerified(fixture.Call, investor)
	failOnErr(t, err, "IsVerified")
	if !verified {
		t.Fatal("investor is NOT verified — KYC claim / registration failed")
	}
	t.Logf("  investor %s is verified ✓", investor.Hex())

	// ---- Step 9: mint shares to the investor and assert the balance ---------
	t.Log("STEP 9 — minting shares to the investor wallet and asserting the balance")
	wei18 := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	amount := new(big.Int).Mul(big.NewInt(1000), wei18) // 1000 tokens

	before, err := suite.Token.BalanceOf(fixture.Call, investor)
	failOnErr(t, err, "BalanceOf before")
	tx, err = suite.Token.Mint(fixture.Auth, investor, amount)
	failOnErr(t, err, "Mint")
	fixture.Mine(tx, "Mint")
	after, err := suite.Token.BalanceOf(fixture.Call, investor)
	failOnErr(t, err, "BalanceOf after")

	gained := new(big.Int).Sub(after, before)
	if gained.Cmp(amount) != 0 {
		t.Fatalf("investor balance: expected +%s, got +%s", amount, gained)
	}
	t.Logf("  ✅ investor received %s base units (balance %s)", amount, after)

	// ---- Step 10: a NON-KYC investor must be rejected ----------------------
	t.Log("STEP 10 — a non-KYC investor must NOT be able to receive shares")
	investor2Key, _ := crypto.GenerateKey()
	investor2 := crypto.PubkeyToAddress(investor2Key.PublicKey)

	// Sanity: this wallet was never registered nor KYC'd.
	v2, err := suite.IR.IsVerified(fixture.Call, investor2)
	failOnErr(t, err, "IsVerified investor2")
	if v2 {
		t.Fatal("investor2 should NOT be verified")
	}

	// Attempt to mint to it — must fail, either at gas estimation (the binding
	// returns an error) or, if sent, on-chain (the transaction reverts).
	tx2, mintErr := suite.Token.Mint(fixture.Auth, investor2, amount)
	if mintErr == nil {
		rcpt, werr := bind.WaitMined(ctx, fixture.Client, tx2)
		switch {
		case werr != nil:
			mintErr = werr
		case rcpt.Status != types.ReceiptStatusSuccessful:
			mintErr = fmt.Errorf("transaction reverted (status 0)")
		}
	}
	if mintErr == nil {
		t.Fatal("expected mint to a non-KYC investor to FAIL, but it succeeded")
	}
	t.Logf("  mint to non-KYC investor correctly rejected: %v", mintErr)

	// The non-KYC wallet must still hold nothing.
	bal2, err := suite.Token.BalanceOf(fixture.Call, investor2)
	failOnErr(t, err, "BalanceOf investor2")
	if bal2.Sign() != 0 {
		t.Fatalf("non-KYC investor balance should be 0, got %s", bal2)
	}
	t.Logf("  ✅ non-KYC investor balance is 0 as expected")
}
