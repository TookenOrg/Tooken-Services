//go:build integration

// Coverage of the investor registration path, through the production code.
//
// This is the test that did not exist, and whose absence let B4 live: the suite was
// green while POST /contract/identity crashed on a nil pointer, because every test
// rebuilt its own registry from local instances instead of calling the service. A
// test that reconstructs its context validates the blockchain, not the application.
//
// Everything here goes through registerIdentity — the real one, reading the real
// blk.contract_role row, against a real chain.
//
// Requires both a node and a database:
//
//	cd tools/hardhat && npx hardhat node
//	TEST_DATABASE_URL="postgres://…" go test -tags integration ./internal/blockchain/services/...
package services

import (
	"context"
	"database/sql"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	chainglobals "github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	appglobals "github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/lib/pq"
)

const registerProbeTxName = "REGISTER_IDENTITY"

func TestRegisterIdentityIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	logger.Init(true)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	previousDB := appglobals.DB
	appglobals.DB = db

	// blk.contract_role is UNIQUE on contract_name, so the SHARED_IRS row has to be
	// moved aside and put back rather than duplicated.
	var savedTxHash, savedAddress sql.NullString
	if err := db.QueryRow(
		`SELECT tx_hash, address FROM blk.contract_role WHERE contract_name = $1`,
		chainglobals.SharedIdentityRegistryStorageName,
	).Scan(&savedTxHash, &savedAddress); err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}

	writtenTxHashes := []string{}
	t.Cleanup(func() {
		_, _ = db.Exec(`DELETE FROM blk.contract_role WHERE contract_name = $1`,
			chainglobals.SharedIdentityRegistryStorageName)
		if savedAddress.Valid {
			_, _ = db.Exec(
				`INSERT INTO blk.contract_role (tx_hash, address, contract_name) VALUES ($1, $2, $3)`,
				savedTxHash.String, savedAddress.String, chainglobals.SharedIdentityRegistryStorageName)
		}
		for _, hash := range writtenTxHashes {
			_, _ = db.Exec(`DELETE FROM blk.eth_transaction WHERE tx_hash = $1`, hash)
		}
		appglobals.DB = previousDB
		db.Close()
	})

	setSharedIRSRow := func(t *testing.T, address common.Address) {
		t.Helper()
		if _, err := db.Exec(`DELETE FROM blk.contract_role WHERE contract_name = $1`,
			chainglobals.SharedIdentityRegistryStorageName); err != nil {
			t.Fatal(err)
		}
		if address == (common.Address{}) {
			return
		}
		if _, err := db.Exec(
			`INSERT INTO blk.contract_role (tx_hash, address, contract_name) VALUES ($1, $2, $3)`,
			"0x"+strings.Repeat("11", 32), address.Hex(),
			chainglobals.SharedIdentityRegistryStorageName); err != nil {
			t.Fatal(err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	fixture := newTREXFixture(ctx, t)
	previousClient := chainglobals.EthClient
	chainglobals.EthClient = fixture.Client
	t.Cleanup(func() { chainglobals.EthClient = previousClient })

	// The suite the platform will be allowed to write into...
	shared := fixture.DeploySuite(t, "RegisterIdentityShared", "TKRS", common.Address{})
	// ...and one the platform is deliberately NOT an agent of, to prove the guard bites.
	foreign := fixture.DeploySuite(t, "RegisterIdentityForeign", "TKRF", common.Address{})

	fixture.MakePlatformIRSAgent(t, shared)

	investor, investorIdentityAddr, investorIdentity := fixture.NewInvestorIdentity(t)

	// ---- the two ways the resolution must fail loudly ----------------------

	t.Run("no shared IRS recorded is an explicit error", func(t *testing.T) {
		setSharedIRSRow(t, common.Address{})

		_, err := registerIdentity(ctx, investor, investorIdentityAddr, 250)
		if err == nil {
			t.Fatal("registering with no SHARED_IRS row must fail")
		}
		// resolveSharedIRS answers the zero address here; binding it would have failed
		// far away from the cause, or silently.
		if !strings.Contains(err.Error(), "no shared IdentityRegistryStorage") {
			t.Fatalf("the error must name the missing IRS, got: %v", err)
		}
	})

	t.Run("an IRS the platform cannot write into is an explicit error", func(t *testing.T) {
		setSharedIRSRow(t, foreign.IRSAddr)

		_, err := registerIdentity(ctx, investor, investorIdentityAddr, 250)
		if err == nil {
			t.Fatal("registering into an IRS the platform is not an agent of must fail")
		}
		if !strings.Contains(err.Error(), "not an agent") {
			t.Fatalf("the error must name the missing agent right, got: %v", err)
		}

		// And it must have failed *before* sending anything.
		stored, err := foreign.IRS.StoredIdentity(fixture.Call, investor)
		failOnErr(t, err, "foreign IRS.storedIdentity")
		if stored != (common.Address{}) {
			t.Fatalf("nothing must have been written to the foreign IRS, got %s", stored.Hex())
		}
	})

	// ---- the happy path, through the production function -------------------

	setSharedIRSRow(t, shared.IRSAddr)

	t.Run("registers the investor in the shared whitelist", func(t *testing.T) {
		tx, err := registerIdentity(ctx, investor, investorIdentityAddr, 250)
		if err != nil {
			t.Fatalf("registerIdentity: %v", err)
		}
		writtenTxHashes = append(writtenTxHashes, tx.Hash().Hex())

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
		if err := db.QueryRow(
			`SELECT count(*) FROM blk.eth_transaction WHERE tx_hash = $1 AND tx_name = $2`,
			tx.Hash().Hex(), registerProbeTxName).Scan(&count); err != nil {
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

	// ---- the promise of writing into the storage: the investor precedes the token

	t.Run("a token created afterwards verifies the investor without re-registering", func(t *testing.T) {
		later := fixture.DeploySuite(t, "RegisterIdentityLater", "TKRL", shared.IRSAddr)

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

	// ---- the trap of the one-off setup -------------------------------------

	t.Run("the factory can still deploy suites after the ownership round-trip", func(t *testing.T) {
		owner, err := shared.IRS.Owner(fixture.Call)
		failOnErr(t, err, "IRS.owner")
		if owner != fixture.FactoryAddr {
			t.Fatalf("IRS ownership must be back with the factory: owner=%s factory=%s",
				owner.Hex(), fixture.FactoryAddr.Hex())
		}
		// The assertion above states the rule; this one proves it bites. Without the
		// hand-back, bindIdentityRegistry is onlyOwner and this call reverts.
		fixture.DeploySuite(t, "RegisterIdentityAfterSetup", "TKRA", shared.IRSAddr)
	})
}
