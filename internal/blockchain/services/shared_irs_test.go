//go:build integration

package services

import (
	"context"
	"math/big"
	"testing"
	"time"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/ethereum/go-ethereum/common"
)

// TestSharedIRSReuseIntegration proves the core promise of the CreateToken refactor:
// a SECOND token deployed through the factory while REUSING the shared
// IdentityRegistryStorage (SHARED_IRS) shares one investor whitelist with the first
// token. An investor KYC'd once (via the IR the factory bound to the shared IRS) is
// verified and can receive shares on BOTH tokens.
//
// It exercises the real production helpers buildTokenDetails + defineClaimSuiteDetails
// (the same construction CreateToken uses, through the shared fixture) and confirms
// that the factory retains ownership of the shared IRS and binds EVERY token's IR to
// it: both IRs become agents of the shared storage, so any token can register
// investors into the one shared whitelist and all tokens verify against it.
//
// Gated by the `integration` build tag. Requires a local node:
//
//	cd tools/hardhat && npm install && npx hardhat node
//	go test -tags integration -run TestSharedIRSReuseIntegration -v ./internal/blockchain/...
func TestSharedIRSReuseIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	t.Log("deploying T-REX infrastructure")
	fixture := newTREXFixture(ctx, t)

	// ---- Token A: fresh IRS (factory deploys + binds A's IR to it) ----------
	t.Log("deploying token A with a fresh IRS")
	tokenA := fixture.DeploySuite(t, "SharedTokenA", "Tooken Property A", "TKPA", common.Address{})
	if tokenA.IRSAddr == (common.Address{}) {
		t.Fatal("token A must have a non-zero IdentityRegistryStorage")
	}
	sharedIRS := tokenA.IRSAddr

	// ---- Token B: REUSE the shared IRS (this is the CreateToken reuse path) --
	t.Log("deploying token B reusing token A's IRS (shared whitelist)")
	tokenB := fixture.DeploySuite(t, "SharedTokenB", "Tooken Property B", "TKPB", sharedIRS)
	if tokenB.IRSAddr != sharedIRS {
		t.Fatalf("token B must reuse the shared IRS: got %s, want %s", tokenB.IRSAddr.Hex(), sharedIRS.Hex())
	}
	t.Logf("  shared IRS = %s (token A IR=%s, token B IR=%s)", sharedIRS.Hex(), tokenA.IRAddr.Hex(), tokenB.IRAddr.Hex())

	// The factory retains ownership of the shared IRS and binds every token's IR to it,
	// so BOTH token A's and token B's IRs must be agents of the shared storage. This is
	// what makes the shared whitelist writable from any token and readable by all.
	linked, err := tokenA.IRS.LinkedIdentityRegistries(fixture.Call)
	failOnErr(t, err, "LinkedIdentityRegistries")
	contains := func(list []common.Address, a common.Address) bool {
		for _, x := range list {
			if x == a {
				return true
			}
		}
		return false
	}
	if !contains(linked, tokenA.IRAddr) {
		t.Fatalf("token A IR %s must be bound to the shared IRS, linked=%v", tokenA.IRAddr.Hex(), linked)
	}
	if !contains(linked, tokenB.IRAddr) {
		t.Fatalf("token B IR %s must also be bound to the shared IRS on reuse, linked=%v", tokenB.IRAddr.Hex(), linked)
	}
	t.Logf("  both token IRs are bound to the shared IRS: %v", linked)

	// ---- KYC an investor ONCE, through the bound IR (token A) ---------------
	t.Log("registering + KYC'ing an investor once, via the bound IR")
	investor, investorIdAddr, investorID := fixture.NewInvestorIdentity(t)

	tx, err := tokenA.IR.RegisterIdentity(fixture.Auth, investor, investorIdAddr, countryCode)
	failOnErr(t, err, "RegisterIdentity via token A IR")
	fixture.Mine(tx, "RegisterIdentity")

	fixture.AddKYCClaim(t, investorIdAddr, investorID)

	// ---- The shared whitelist must make the investor verified on BOTH tokens -
	vA, err := tokenA.IR.IsVerified(fixture.Call, investor)
	failOnErr(t, err, "IsVerified token A")
	if !vA {
		t.Fatal("investor must be verified on token A")
	}
	vB, err := tokenB.IR.IsVerified(fixture.Call, investor)
	failOnErr(t, err, "IsVerified token B")
	if !vB {
		t.Fatal("investor must be verified on token B via the SHARED IRS — reuse is broken")
	}
	t.Log("  investor verified on BOTH token A and token B ✓")

	// ---- Mint on BOTH tokens must succeed for the single KYC'd investor -----
	wei18 := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	amount := new(big.Int).Mul(big.NewInt(500), wei18)

	for _, tc := range []struct {
		name  string
		token *contracts.Token
	}{{"token A", tokenA.Token}, {"token B", tokenB.Token}} {
		tx, err = tc.token.Mint(fixture.Auth, investor, amount)
		failOnErr(t, err, "Mint "+tc.name)
		fixture.Mine(tx, "Mint "+tc.name)
		bal, err := tc.token.BalanceOf(fixture.Call, investor)
		failOnErr(t, err, "BalanceOf "+tc.name)
		if bal.Cmp(amount) != 0 {
			t.Fatalf("%s: expected balance %s, got %s", tc.name, amount, bal)
		}
		t.Logf("  ✅ %s minted %s to the shared-whitelist investor", tc.name, amount)
	}
}
