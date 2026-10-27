//go:build integration

// Shared chain + database environment for the integration tests of this package.
//
// The point of these tests is to travel the production path, and the production path
// reads the database: which token to reach the IdentityRegistry through, which address
// the shared whitelist is anchored at, which ONCHAINID authority to deploy proxies
// from. Rebuilding that context from local variables is exactly what let B4 live.
//
// Each run builds its own database from tools/db/baseline.sql and drops it at the end
// (see testenv.ScratchDatabase).
// Sharing one database across packages does not work here: the handler suite leaves
// rows in blk.token, and "no token deployed yet" is a case this package has to be able
// to reproduce. Isolation also means these tests can neither be broken by, nor break,
// another package's fixtures.
//
// Requires a node and a PostgreSQL user allowed to create databases:
//
//	cd tools/hardhat && npx hardhat node
//	TEST_DATABASE_URL="postgres://…" go test -tags integration ./internal/blockchain/services/...
package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	chainglobals "github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	appglobals "github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/testenv"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/lib/pq"
)

const (
	// scratchDatabaseName is rebuilt on every run, so a suite interrupted halfway
	// never poisons the next one.
	scratchDatabaseName = "tooken_it_blockchain_services"

	// blk tables CHECK that a transaction hash is 0x followed by 64 hex characters,
	// so a marker has to be made of hex digits too.
	probeTxHashPrefix = "0xfeedface"
)

// chainDBEnv is a live node plus a private database, both wired into the package
// globals the production code reads.
type chainDBEnv struct {
	DB      *sql.DB
	Fixture *trexFixture
}

// newBareChainDBEnv gives a private database built from the baseline and a client on
// the local node — and nothing else.
//
// No contract is deployed, no row is seeded: the point is to let the *production*
// initialisation routes do that work, which is the only way to prove they can. Use
// newChainDBEnv instead when a test needs a T-REX infrastructure it did not build.
func newBareChainDBEnv(ctx context.Context, t *testing.T) *chainDBEnv {
	t.Helper()

	db := testenv.ScratchDatabase(t, scratchDatabaseName)
	logger.Init(true)
	client := testenv.DialNode(ctx, t)

	previousDB := appglobals.DB
	appglobals.DB = db
	previousClient := chainglobals.EthClient
	chainglobals.EthClient = client

	// GenerateTransactOpts signs with this key, exactly like production does.
	t.Setenv("PRIVATE_KEY", testenv.HardhatAccount0Key)

	t.Cleanup(func() {
		chainglobals.EthClient = previousClient
		appglobals.DB = previousDB
	})

	return &chainDBEnv{DB: db}
}

// newChainDBEnv adds a deployed T-REX infrastructure to the bare environment, and
// announces in the database the two contracts production looks up by name.
func newChainDBEnv(ctx context.Context, t *testing.T) *chainDBEnv {
	t.Helper()

	env := newBareChainDBEnv(ctx, t)
	env.Fixture = newTREXFixture(ctx, t)

	previousClient := chainglobals.EthClient
	chainglobals.EthClient = env.Fixture.Client
	t.Cleanup(func() { chainglobals.EthClient = previousClient })

	env.SetImplementationRow(t, chainglobals.ImplIdentityAuthorityName, env.Fixture.IdentityAuthorityRef)
	env.SetContractRole(t, chainglobals.ClaimIssuerName, env.Fixture.ClaimIssuer)

	return env
}

// NewUser inserts an investor and returns its id.
//
// blk.user_wallet and blk.identity carry a foreign key to usr.users since 000026:
// an ONCHAINID can no longer be deployed for somebody who does not exist. Tests
// therefore have to create their investor rather than invent an id — which is the
// point, since production has to do the same.
func (e *chainDBEnv) NewUser(t *testing.T) int {
	t.Helper()

	var id int
	if err := e.DB.QueryRow(`
		INSERT INTO usr.users (full_name, email, password)
		VALUES ($1, $2, 'not-a-real-hash')
		RETURNING id`,
		"Test Investor",
		fmt.Sprintf("investor-%d@test.local", time.Now().UnixNano()),
	).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}

// SetContractRole makes a contract role point at an address, or removes it when given
// the zero address. blk.contract_role is UNIQUE on contract_name, so the row is
// replaced rather than duplicated.
func (e *chainDBEnv) SetContractRole(t *testing.T, name string, address common.Address) {
	t.Helper()

	if _, err := e.DB.Exec(`DELETE FROM blk.contract_role WHERE contract_name = $1`, name); err != nil {
		t.Fatal(err)
	}
	if address == (common.Address{}) {
		return
	}
	if _, err := e.DB.Exec(
		`INSERT INTO blk.contract_role (tx_hash, address, contract_name) VALUES ($1, $2, $3)`,
		probeTxHash(name), address.Hex(), name); err != nil {
		t.Fatal(err)
	}
}

// SetSharedIRSRow points the shared investor whitelist anchor at an address, or
// removes it when given the zero address.
func (e *chainDBEnv) SetSharedIRSRow(t *testing.T, address common.Address) {
	t.Helper()
	e.SetContractRole(t, chainglobals.SharedIdentityRegistryStorageName, address)
}

// SetImplementationRow announces a deployed implementation under the name production
// looks it up by.
func (e *chainDBEnv) SetImplementationRow(t *testing.T, name string, address common.Address) {
	t.Helper()

	if _, err := e.DB.Exec(`DELETE FROM blk.contract_implementation WHERE contract_name = $1`, name); err != nil {
		t.Fatal(err)
	}
	if _, err := e.DB.Exec(
		`INSERT INTO blk.contract_implementation (tx_hash, address, contract_name) VALUES ($1, $2, $3)`,
		probeTxHash(name), address.Hex(), name); err != nil {
		t.Fatal(err)
	}
}

// InsertTokenRow announces a deployed token in blk.token, which is how Mint and Burn
// find its decimals before touching the chain.
func (e *chainDBEnv) InsertTokenRow(t *testing.T, address common.Address, name, symbol string, decimals int) {
	t.Helper()

	if _, err := e.DB.Exec(`
		INSERT INTO blk.token (symbol, token_name, salt, address, nb_decimal)
		VALUES ($1, $2, $3, $4, $5)`,
		symbol, name, "probe-"+symbol, address.Hex(), decimals); err != nil {
		t.Fatal(err)
	}
}

// probeTxHash builds a recognisable and schema-valid transaction hash: the columns
// CHECK for 0x followed by 64 hex characters, and several tables are UNIQUE on it, so
// it has to look real and stay distinct per seed.
func probeTxHash(seed string) string {
	digest := fmt.Sprintf("%x", []byte(seed))
	hash := probeTxHashPrefix + digest
	if len(hash) > 66 {
		return hash[:66]
	}
	return hash + strings.Repeat("0", 66-len(hash))
}
