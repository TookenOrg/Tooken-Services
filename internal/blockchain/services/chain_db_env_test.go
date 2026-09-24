//go:build integration

// Shared chain + database environment for the integration tests of this package.
//
// The point of these tests is to travel the production path, and the production path
// reads the database: which token to reach the IdentityRegistry through, which address
// the shared whitelist is anchored at, which ONCHAINID authority to deploy proxies
// from. Rebuilding that context from local variables is exactly what let B4 live.
//
// Each run builds its own database from tools/db/baseline.sql and drops it at the end.
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
	"net/url"
	"os"
	"strings"
	"testing"

	chainglobals "github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	appglobals "github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
)

const (
	// scratchDatabaseName is rebuilt on every run, so a suite interrupted halfway
	// never poisons the next one.
	scratchDatabaseName = "tooken_it_blockchain_services"
	baselineSQLPath     = "../../../tools/db/baseline.sql"

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

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	logger.Init(true)

	client := dialTestNode(ctx, t)
	db := createScratchDatabase(t, dsn)

	previousDB := appglobals.DB
	appglobals.DB = db
	previousClient := chainglobals.EthClient
	chainglobals.EthClient = client

	// GenerateTransactOpts signs with this key, exactly like production does.
	t.Setenv("PRIVATE_KEY", hardhatAccount0Key)

	t.Cleanup(func() {
		chainglobals.EthClient = previousClient
		appglobals.DB = previousDB
		db.Close()
		dropScratchDatabase(t, dsn)
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

// dialTestNode connects to the local EVM node, skipping the test — rather than
// failing it — when none answers.
//
// It dials WebSocket by default because production does: SetupEthClient reads
// WS_RPC_URL. The transport is not a detail — eth_subscribe does not exist over
// HTTP, so a subscription bug looks like "notifications not supported" there and
// like a silent timeout in production. Testing over the transport production uses is
// what makes the difference visible.
func dialTestNode(ctx context.Context, t *testing.T) *ethclient.Client {
	t.Helper()

	rpc := os.Getenv("ETH_TEST_RPC")
	if rpc == "" {
		rpc = "ws://127.0.0.1:8545"
	}
	client, err := ethclient.DialContext(ctx, rpc)
	if err != nil {
		t.Skipf("no local EVM node at %s: %v", rpc, err)
	}
	if _, err := client.ChainID(ctx); err != nil {
		t.Skipf("local EVM node at %s not reachable: %v", rpc, err)
	}
	return client
}

// createScratchDatabase rebuilds this package's own database from the committed
// baseline, so the tests read the real schema — constraints, triggers and all —
// rather than a hand-made approximation of it.
func createScratchDatabase(t *testing.T, dsn string) *sql.DB {
	t.Helper()

	adminDSN, scratchDSN, err := scratchDSNs(dsn)
	if err != nil {
		t.Skipf("TEST_DATABASE_URL is not a URL a database can be derived from: %v", err)
	}

	admin, err := sql.Open("postgres", adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	if err := admin.Ping(); err != nil {
		t.Skipf("database server not reachable: %v", err)
	}

	if _, err := admin.Exec(`DROP DATABASE IF EXISTS ` + scratchDatabaseName); err != nil {
		t.Skipf("cannot drop %s (the test user needs CREATEDB): %v", scratchDatabaseName, err)
	}
	if _, err := admin.Exec(`CREATE DATABASE ` + scratchDatabaseName); err != nil {
		t.Skipf("cannot create %s (the test user needs CREATEDB): %v", scratchDatabaseName, err)
	}

	baseline, err := os.ReadFile(baselineSQLPath)
	if err != nil {
		t.Fatalf("read %s: %v", baselineSQLPath, err)
	}

	db, err := sql.Open("postgres", scratchDSN)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(string(baseline)); err != nil {
		db.Close()
		t.Fatalf("apply %s: %v", baselineSQLPath, err)
	}

	return db
}

func dropScratchDatabase(t *testing.T, dsn string) {
	t.Helper()

	adminDSN, _, err := scratchDSNs(dsn)
	if err != nil {
		return
	}
	admin, err := sql.Open("postgres", adminDSN)
	if err != nil {
		return
	}
	defer admin.Close()
	_, _ = admin.Exec(`DROP DATABASE IF EXISTS ` + scratchDatabaseName)
}

// scratchDSNs turns the configured DSN into one pointing at the maintenance database
// and one pointing at this package's scratch database, keeping every other parameter
// — credentials, host, sslmode — as configured.
func scratchDSNs(dsn string) (adminDSN, scratchDSN string, err error) {
	parsed, err := url.Parse(dsn)
	if err != nil {
		return "", "", err
	}
	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		return "", "", fmt.Errorf("unsupported scheme %q", parsed.Scheme)
	}

	admin := *parsed
	admin.Path = "/postgres"

	scratch := *parsed
	scratch.Path = "/" + scratchDatabaseName

	return admin.String(), scratch.String(), nil
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
