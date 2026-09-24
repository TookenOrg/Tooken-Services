//go:build integration

// Coverage of the initialisation sequence, on a database that starts empty.
//
// This is the test that did not exist, and whose absence let three blocking defects
// live: deployed implementations that were never recorded, an identity factory whose
// address was returned and thrown away, and a suite deployment that waited for an
// event it had already let pass. None of them could be caught by a unit test, and the
// existing integration tests all started from a T-REX infrastructure built in Go
// rather than through the production routes.
//
// Everything here goes through the real services, in the order an operator calls
// them, and no row is ever seeded by hand. If a step needs a manual INSERT to pass,
// the environment is not installable — which was exactly the situation on 2026-09-18.
//
// 🔴 It also runs on an auto-mining node, on purpose: that is the configuration that
// made CreateToken time out, and the one the default Hardhat setup gives.
//
// Requires both a node and a database:
//
//	cd tools/hardhat && npx hardhat node
//	TEST_DATABASE_URL="postgres://…" go test -tags integration ./internal/blockchain/services/...
package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	chainglobals "github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/lib/pq"
)

// The nine implementations DeployAllImplementations is expected to record, under the
// names the rest of the code looks them up by. Taken from globals so a rename cannot
// silently desynchronise the test from production.
func expectedImplementationNames() []string {
	return []string{
		chainglobals.ImplClaimsTopicRegistryName,
		chainglobals.ImplTrustedIssuerRegistryName,
		chainglobals.ImplIdentityRegistryStorageName,
		chainglobals.ImplIdentityRegistryName,
		chainglobals.ImplModularComplianceName,
		chainglobals.ImplTokenName,
		chainglobals.ImplIdentityName,
		chainglobals.ImplIdentityAuthorityName,
		chainglobals.ImplTrexAuthorityName,
	}
}

func TestInitSequenceIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	env := newBareChainDBEnv(ctx, t)
	svc := NewService()

	// The sub-tests below are a sequence, not a set: each one is the precondition of
	// the next, exactly like the six calls an operator makes. Running them in order on
	// a database that started empty *is* the assertion.

	// ---- step 1: the implementations ------------------------------------------

	t.Run("DeployAllImplementations records the nine implementations", func(t *testing.T) {
		deployed, err := svc.DeployAllImplementations(ctx)
		if err != nil {
			t.Fatalf("DeployAllImplementations: %v", err)
		}
		if len(deployed) != len(expectedImplementationNames()) {
			t.Fatalf("expected %d contracts in the response, got %d", len(expectedImplementationNames()), len(deployed))
		}

		var rows int
		if err := env.DB.QueryRow(`SELECT count(*) FROM blk.contract_implementation`).Scan(&rows); err != nil {
			t.Fatal(err)
		}
		if rows != len(expectedImplementationNames()) {
			t.Fatalf("blk.contract_implementation must hold %d rows, found %d",
				len(expectedImplementationNames()), rows)
		}
	})

	t.Run("each implementation is readable under the name production uses", func(t *testing.T) {
		// Counting rows is not enough: a name typed by hand instead of taken from
		// globals fills the table just as well, and only shows up two calls later.
		for _, name := range expectedImplementationNames() {
			contract, err := database.GetImplementationContractByName(ctx, name)
			if err != nil {
				t.Fatalf("%s must be readable by name: %v", name, err)
			}
			if !common.IsHexAddress(contract.Address) || contract.Address == (common.Address{}).Hex() {
				t.Fatalf("%s recorded an unusable address: %q", name, contract.Address)
			}
		}
	})

	// ---- step 2: the identity factory -----------------------------------------

	var identityFactoryAddress string

	t.Run("DeployIdentityFactory records the role", func(t *testing.T) {
		details, err := svc.DeployIdentityFactory(ctx)
		if err != nil {
			t.Fatalf("DeployIdentityFactory: %v", err)
		}
		identityFactoryAddress = details.Address

		role, err := database.GetContractRoleByName(ctx, chainglobals.IdentityFactoryName)
		if err != nil {
			t.Fatalf("the identity factory must be recorded as a contract role: %v", err)
		}
		if !strings.EqualFold(role.Address, details.Address) {
			t.Fatalf("the recorded address %q must be the one returned %q", role.Address, details.Address)
		}
	})

	t.Run("DeployIdentityFactory replayed answers the same address without redeploying", func(t *testing.T) {
		// blk.contract_role is UNIQUE on contract_name: a second deployment would
		// either violate the constraint or silently orphan the first factory — which
		// is the contract every investor ONCHAINID is created by.
		details, err := svc.DeployIdentityFactory(ctx)
		if err != nil {
			t.Fatalf("a replayed init route must not fail: %v", err)
		}
		if !strings.EqualFold(details.Address, identityFactoryAddress) {
			t.Fatalf("a replay must answer the existing factory %q, got %q", identityFactoryAddress, details.Address)
		}

		var rows int
		if err := env.DB.QueryRow(
			`SELECT count(*) FROM blk.contract_role WHERE contract_name = $1`,
			chainglobals.IdentityFactoryName).Scan(&rows); err != nil {
			t.Fatal(err)
		}
		if rows != 1 {
			t.Fatalf("exactly one identity factory must be recorded, found %d", rows)
		}
	})

	// ---- step 3: the authority ------------------------------------------------

	t.Run("ConfigureAuthority succeeds with no manual seeding", func(t *testing.T) {
		// 🔴 This is the step that stopped the installation on 2026-09-18: it reads
		// every implementation by name, and the table was empty.
		result, err := svc.ConfigureAuthority(ctx)
		if err != nil {
			t.Fatalf("ConfigureAuthority: %v", err)
		}
		if result.TransactionHash == "" {
			t.Fatal("ConfigureAuthority must report the transaction it sent")
		}
	})

	// ---- step 4: the TREX factory ---------------------------------------------

	t.Run("DeployAndInitTrexFactory succeeds and records the role", func(t *testing.T) {
		// 🔴 The second stop: it looks the identity factory up in blk.contract_role.
		if err := svc.DeployAndInitTrexFactory(ctx); err != nil {
			t.Fatalf("DeployAndInitTrexFactory: %v", err)
		}

		role, err := database.GetContractRoleByName(ctx, chainglobals.TrexFactoryName)
		if err != nil {
			t.Fatalf("the TREX factory must be recorded as a contract role: %v", err)
		}
		if !common.IsHexAddress(role.Address) {
			t.Fatalf("the TREX factory recorded an unusable address: %q", role.Address)
		}
	})

	// ---- step 5: the shared singletons ----------------------------------------

	t.Run("DeployTrexSuite records the claim issuer and the compliance module", func(t *testing.T) {
		if _, err := svc.DeployTrexSuite(ctx); err != nil {
			t.Fatalf("DeployTrexSuite: %v", err)
		}
		for _, name := range []string{chainglobals.ClaimIssuerName, chainglobals.TransferRestrictionModuleName} {
			if _, err := database.GetContractRoleByName(ctx, name); err != nil {
				t.Fatalf("%s must be recorded as a contract role: %v", name, err)
			}
		}
	})

	// ---- step 6: a token, on an auto-mining node -------------------------------

	t.Run("CreateToken succeeds on an auto-mining node", func(t *testing.T) {
		// 🔴 The third stop, and the one that would have survived into production as an
		// intermittent failure. The previous implementation subscribed to a *future*
		// event after sending the transaction; on a node that mines immediately the
		// event was already gone, and the call waited two minutes for nothing.
		//
		// Nothing here slows the node down on purpose: the default Hardhat automine is
		// the configuration this test exists to cover.
		token, err := svc.CreateToken(ctx, server.CreateTokenRequest{
			TokenName: "Tooken Init Sequence",
			Symbol:    "TKIS",
			NbDecimal: 18,
		})
		if err != nil {
			t.Fatalf("CreateToken: %v", err)
		}
		if !common.IsHexAddress(token.Address) {
			t.Fatalf("CreateToken must answer a token address, got %q", token.Address)
		}

		// And the whole point of the sequence: the platform can now onboard investors.
		if _, err := database.GetContractRoleByName(ctx, chainglobals.SharedIdentityRegistryStorageName); err != nil {
			t.Fatalf("the shared investor whitelist must be anchored: %v", err)
		}
	})

	t.Run("CreateToken refuses a salt already spent", func(t *testing.T) {
		// The guard written for TICKET-13. It fired for real during the manual demo,
		// after the timeout above had left a suite on-chain with no row in blk.token.
		_, err := svc.CreateToken(ctx, server.CreateTokenRequest{
			TokenName: "Tooken Init Sequence",
			Symbol:    "TKIS",
			NbDecimal: 18,
		})
		if err == nil {
			t.Fatal("deploying twice under the same salt must be refused")
		}
	})
}
