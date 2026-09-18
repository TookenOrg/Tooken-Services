//go:build integration

// Coverage of the investor registration path, through the production code.
//
// This is the test that did not exist, and whose absence let B4 live: the suite was
// green while POST /contract/identity crashed on a nil pointer, because every test
// rebuilt its own registry from local instances instead of calling the service. A
// test that reconstructs its context validates the blockchain, not the application.
//
// Everything here goes through registerIdentity — the real one, reading the real
// blk.token and blk.contract_role rows, against a real chain.
//
// Requires both a node and a database:
//
//	cd tools/hardhat && npx hardhat node
//	TEST_DATABASE_URL="postgres://…" go test -tags integration ./internal/blockchain/services/...
package services

import (
	"context"
	"math/big"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	_ "github.com/lib/pq"
)

func TestRegisterIdentityIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	env := newChainDBEnv(ctx, t)
	fixture := env.Fixture

	// The suite the platform created itself: its registry is bound to the storage,
	// and the platform is one of its agents.
	shared := fixture.DeploySuite(t, "RegisterIdentityShared", "Tooken Shared", "TKRS", common.Address{})
	env.SetSharedIRSRow(t, shared.IRSAddr)

	investor, investorIdentityAddr, investorIdentity := fixture.NewInvestorIdentity(t)

	// ---- 🔑 the claim of this design: T-REX already granted every right --------

	t.Run("the authorisation chain exists with no setup whatsoever", func(t *testing.T) {
		// bindIdentityRegistry, called by the factory when the suite was deployed,
		// made the IdentityRegistry an agent of the storage — and listed it.
		linked, err := shared.IRS.LinkedIdentityRegistries(fixture.Call)
		failOnErr(t, err, "IRS.linkedIdentityRegistries")
		if len(linked) == 0 || linked[0] != shared.IRAddr {
			t.Fatalf("the storage must list the registry bound to it, got %v", linked)
		}

		irIsStorageAgent, err := shared.IRS.IsAgent(fixture.Call, shared.IRAddr)
		failOnErr(t, err, "IRS.isAgent(IR)")
		if !irIsStorageAgent {
			t.Fatal("the IdentityRegistry must be an agent of its storage — T-REX binds it at deployment")
		}

		// buildTokenDetails passed the platform in IrAgents.
		platformIsIRAgent, err := shared.IR.IsAgent(fixture.Call, fixture.Deployer)
		failOnErr(t, err, "IR.isAgent(platform)")
		if !platformIsIRAgent {
			t.Fatal("the platform must be an agent of the IdentityRegistry — buildTokenDetails puts it in IrAgents")
		}

		// And the platform is deliberately NOT an agent of the storage. This is the
		// assertion that keeps the implementation on the standard path: the day
		// someone makes the platform write into the storage directly, this fails.
		platformIsStorageAgent, err := shared.IRS.IsAgent(fixture.Call, fixture.Deployer)
		failOnErr(t, err, "IRS.isAgent(platform)")
		if platformIsStorageAgent {
			t.Fatal("the platform must not be an agent of the storage: writing goes through the IdentityRegistry, as T-REX intends")
		}

		owner, err := shared.IRS.Owner(fixture.Call)
		failOnErr(t, err, "IRS.owner")
		if owner != fixture.FactoryAddr {
			t.Fatalf("the factory must keep the storage ownership, got %s", owner.Hex())
		}
	})

	// ---- the three ways the resolution must fail loudly ------------------------

	t.Run("no shared whitelist anchor is an explicit error", func(t *testing.T) {
		env.SetSharedIRSRow(t, common.Address{}) // deletes the row

		_, err := registerIdentity(ctx, investor, investorIdentityAddr, 250)
		if err == nil {
			t.Fatal("registering with no SHARED_IRS anchor must fail")
		}
		if !strings.Contains(err.Error(), "no shared IdentityRegistryStorage recorded") {
			t.Fatalf("the error must name the missing anchor, got: %v", err)
		}
	})

	t.Run("a whitelist with no registry bound to it is an explicit error", func(t *testing.T) {
		// A storage deployed on its own, outside any suite: nothing has been bound to
		// it, so there is no door at all. Without the guard the loop would simply find
		// nothing and the failure would surface much later.
		orphanIRS := fixture.NewStandaloneIRS(t)
		env.SetSharedIRSRow(t, orphanIRS)

		_, err := registerIdentity(ctx, investor, investorIdentityAddr, 250)
		if err == nil {
			t.Fatal("registering through a storage with no bound registry must fail")
		}
		if !strings.Contains(err.Error(), "no IdentityRegistry is bound") {
			t.Fatalf("the error must say no registry is bound, got: %v", err)
		}
	})

	t.Run("a whitelist whose registries refuse the platform is an explicit error", func(t *testing.T) {
		// The platform is an agent of every registry it creates, so this cannot happen
		// by accident — which is exactly why the guard would otherwise never be
		// exercised, and would rot. Here the right is revoked on purpose, on a suite
		// deployed with its own storage so it is the *only* door of that whitelist.
		revoked := fixture.DeploySuite(t, "RegisterIdentityRevoked", "Tooken Revoked", "TKRV", common.Address{})
		tx, err := revoked.IR.RemoveAgent(fixture.Auth, fixture.Deployer)
		failOnErr(t, err, "IR.removeAgent(platform)")
		fixture.Mine(tx, "IR.removeAgent")

		env.SetSharedIRSRow(t, revoked.IRSAddr)

		_, err = registerIdentity(ctx, investor, investorIdentityAddr, 250)
		if err == nil {
			t.Fatal("registering with no usable registry must fail")
		}
		// Without the guard the chain answers "AgentRole: caller does not have the
		// Agent role", which names neither the registry nor the platform.
		if !strings.Contains(err.Error(), "is an agent of none of the") {
			t.Fatalf("the error must say no door accepts the platform, got: %v", err)
		}

		// And it must have failed before sending anything.
		stored, err := revoked.IRS.StoredIdentity(fixture.Call, investor)
		failOnErr(t, err, "IRS.storedIdentity after the refusal")
		if stored != (common.Address{}) {
			t.Fatalf("nothing must have been written, got %s", stored.Hex())
		}
	})

	// ---- the happy path, through the production function -----------------------

	env.SetSharedIRSRow(t, shared.IRSAddr)

	t.Run("registers the investor in the shared whitelist", func(t *testing.T) {
		tx, err := registerIdentity(ctx, investor, investorIdentityAddr, 250)
		if err != nil {
			t.Fatalf("registerIdentity: %v", err)
		}

		stored, err := shared.IRS.StoredIdentity(fixture.Call, investor)
		failOnErr(t, err, "IRS.storedIdentity")
		if stored != investorIdentityAddr {
			t.Fatalf("stored identity: expected %s, got %s", investorIdentityAddr.Hex(), stored.Hex())
		}

		country, err := shared.IRS.StoredInvestorCountry(fixture.Call, investor)
		failOnErr(t, err, "IRS.storedInvestorCountry")
		if country != 250 {
			t.Fatalf("stored country: expected 250, got %d", country)
		}

		// The transaction must be traceable, like every other on-chain write.
		var count int
		if err := env.DB.QueryRow(
			`SELECT count(*) FROM blk.eth_transaction WHERE tx_hash = $1 AND tx_name = $2`,
			tx.Hash().Hex(), "REGISTER_IDENTITY").Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("the registration must be recorded in blk.eth_transaction, found %d rows", count)
		}
	})

	t.Run("a registered investor without a claim is still not verified", func(t *testing.T) {
		// Registration alone is not KYC: the claim is the other half, and the token
		// checks both. Asserting this before adding the claim is what proves the next
		// sub-test is measuring the claim and not the registration.
		verified, err := shared.IR.IsVerified(fixture.Call, investor)
		failOnErr(t, err, "IsVerified before the claim")
		if verified {
			t.Fatal("an investor with no KYC claim must not be verified")
		}
	})

	t.Run("with the KYC claim the investor becomes verified and can be minted to", func(t *testing.T) {
		fixture.AddKYCClaim(t, investorIdentityAddr, investorIdentity)

		verified, err := shared.IR.IsVerified(fixture.Call, investor)
		failOnErr(t, err, "IsVerified after the claim")
		if !verified {
			t.Fatal("a registered investor holding a KYC claim must be verified")
		}

		amount := new(big.Int).Mul(big.NewInt(1000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
		before, err := shared.Token.BalanceOf(fixture.Call, investor)
		failOnErr(t, err, "BalanceOf before")
		tx, err := shared.Token.Mint(fixture.Auth, investor, amount)
		failOnErr(t, err, "Mint")
		fixture.Mine(tx, "Mint")
		after, err := shared.Token.BalanceOf(fixture.Call, investor)
		failOnErr(t, err, "BalanceOf after")

		if gained := new(big.Int).Sub(after, before); gained.Cmp(amount) != 0 {
			t.Fatalf("investor balance: expected +%s, got +%s", amount, gained)
		}
	})

	t.Run("a wallet that was never registered is refused by the chain", func(t *testing.T) {
		stranger := common.HexToAddress("0x000000000000000000000000000000000000dEaD")

		verified, err := shared.IR.IsVerified(fixture.Call, stranger)
		failOnErr(t, err, "IsVerified stranger")
		if verified {
			t.Fatal("an unknown wallet must not be verified")
		}

		_, err = shared.Token.Mint(fixture.Auth, stranger, big.NewInt(1))
		if err == nil {
			t.Fatal("minting to a non-KYC wallet must revert")
		}
		if !strings.Contains(err.Error(), "Identity is not verified") {
			t.Fatalf("the revert must name the reason, got: %v", err)
		}

		balance, err := shared.Token.BalanceOf(fixture.Call, stranger)
		failOnErr(t, err, "BalanceOf stranger")
		if balance.Sign() != 0 {
			t.Fatalf("the non-KYC wallet must hold nothing, got %s", balance)
		}
	})

	// ---- one whitelist, every token --------------------------------------------

	t.Run("a token created afterwards verifies the investor without re-registering", func(t *testing.T) {
		later := fixture.DeploySuite(t, "RegisterIdentityLater", "Tooken Later", "TKRL", shared.IRSAddr)

		verified, err := later.IR.IsVerified(fixture.Call, investor)
		failOnErr(t, err, "IsVerified on the later token")
		if !verified {
			t.Fatal("an investor registered before a token must be verified on it")
		}

		amount := new(big.Int).Mul(big.NewInt(500), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
		tx, err := later.Token.Mint(fixture.Auth, investor, amount)
		failOnErr(t, err, "Mint on the later token")
		fixture.Mine(tx, "Mint on the later token")

		balance, err := later.Token.BalanceOf(fixture.Call, investor)
		failOnErr(t, err, "BalanceOf on the later token")
		if balance.Cmp(amount) != 0 {
			t.Fatalf("later token balance: expected %s, got %s", amount, balance)
		}
	})

	t.Run("the door used is one the whitelist itself listed", func(t *testing.T) {
		// Which registry is used must not matter — they all write into the same
		// storage — but it must be one the storage knows about. A registry picked
		// anywhere else could be bound to another ledger, and produce a ghost KYC.
		later := fixture.DeploySuite(t, "RegisterIdentitySecondDoor", "Tooken Door", "TKRD", shared.IRSAddr)

		second, secondIdentityAddr, secondIdentity := fixture.NewInvestorIdentity(t)

		tx, err := registerIdentity(ctx, second, secondIdentityAddr, 442)
		if err != nil {
			t.Fatalf("registerIdentity: %v", err)
		}

		linked, err := shared.IRS.LinkedIdentityRegistries(fixture.Call)
		failOnErr(t, err, "IRS.linkedIdentityRegistries")
		to := tx.To()
		if to == nil {
			t.Fatal("a registration is a call, it must have a recipient")
		}
		if !slices.Contains(linked, *to) {
			t.Fatalf("the registration went to %s, which the storage does not list: %v", to.Hex(), linked)
		}
		if later.IRAddr == shared.IRAddr {
			t.Fatal("the second suite must have its own registry, otherwise this proves nothing")
		}

		// Written once, readable from the first token's registry too.
		stored, err := shared.IRS.StoredIdentity(fixture.Call, second)
		failOnErr(t, err, "IRS.storedIdentity for the second investor")
		if stored != secondIdentityAddr {
			t.Fatalf("stored identity: expected %s, got %s", secondIdentityAddr.Hex(), stored.Hex())
		}

		fixture.AddKYCClaim(t, secondIdentityAddr, secondIdentity)
		verified, err := shared.IR.IsVerified(fixture.Call, second)
		failOnErr(t, err, "IsVerified on the first token")
		if !verified {
			t.Fatal("an investor registered through one token must be verified on the others")
		}
	})

	t.Run("the factory can still deploy suites", func(t *testing.T) {
		// Registration touches no ownership, so this can never break — which is
		// precisely the property the standard path buys. The assertion stays because
		// an implementation that starts moving ownership around would fail it.
		owner, err := shared.IRS.Owner(fixture.Call)
		failOnErr(t, err, "IRS.owner")
		if owner != fixture.FactoryAddr {
			t.Fatalf("storage ownership must never leave the factory: owner=%s factory=%s",
				owner.Hex(), fixture.FactoryAddr.Hex())
		}
		fixture.DeploySuite(t, "RegisterIdentityAfterAll", "Tooken After", "TKRA", shared.IRSAddr)
	})
}
