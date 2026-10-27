// Package testenv gives integration tests the real infrastructure they run on:
// a private PostgreSQL database built from the committed schema, and a client
// on the local EVM node.
//
// It is imported by tests only. The server binary must never link it, which
// `go list -deps . | grep internal/testenv` checks: it must print nothing.
//
// Two environment variables drive it:
//
//	TEST_DATABASE_URL  a PostgreSQL URL whose user may create databases
//	ETH_TEST_RPC       the node to dial, ws://127.0.0.1:8545 when unset
//
// Setting either of them means an integration run was asked for: from then on
// an unreachable database or node fails the test instead of skipping it.
package testenv

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
)

// HardhatAccount0Key is the private key of the first account of a Hardhat
// node. It is published by Hardhat itself and funds nothing outside a local
// node; production reads its signer from PRIVATE_KEY.
const HardhatAccount0Key = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

// DefaultNodeURL is dialled when ETH_TEST_RPC is unset. It is a WebSocket URL
// because production dials WS_RPC_URL: eth_subscribe does not exist over HTTP,
// so a subscription bug would look different over the other transport.
const DefaultNodeURL = "ws://127.0.0.1:8545"

// A database name is spliced into DROP and CREATE DATABASE, which take no
// parameters: only plain lower-case identifiers are accepted. The prefix makes
// sure a typo can never drop a real database, such as the one the URL names.
var databaseName = regexp.MustCompile(`^tooken_it_[a-z0-9_]{1,53}$`)

// ScratchDatabase builds the database name from tools/db/baseline.sql and
// returns a connection to it. The database is dropped when the test ends.
//
// Each package passes its own name. Packages may run one after the other in a
// single `go test` invocation, but a name shared by two of them would let one
// drop the other's rows; and a name is rebuilt on every run, so a suite killed
// halfway never poisons the next one.
//
// The test is skipped when TEST_DATABASE_URL is unset, even when ETH_TEST_RPC
// is set, so that a run against the node alone stays possible.
func ScratchDatabase(t testing.TB, name string) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	if !databaseName.MatchString(name) {
		t.Fatalf("scratch database name %q must be tooken_it_ followed by lower-case letters, digits or underscores", name)
	}

	adminDSN, scratchDSN, err := scratchDSNs(dsn, name)
	if err != nil {
		Unavailable(t, "TEST_DATABASE_URL is not a URL a database can be derived from: %v", err)
	}

	admin, err := sql.Open("postgres", adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	if err := admin.Ping(); err != nil {
		Unavailable(t, "database server not reachable: %v", err)
	}
	if _, err := admin.Exec(`DROP DATABASE IF EXISTS ` + name); err != nil {
		Unavailable(t, "cannot drop %s (the test user needs CREATEDB): %v", name, err)
	}
	if _, err := admin.Exec(`CREATE DATABASE ` + name); err != nil {
		Unavailable(t, "cannot create %s (the test user needs CREATEDB): %v", name, err)
	}
	// Registered first, so it runs last: after the connection below is closed,
	// and even when a later step fails the test.
	t.Cleanup(func() { dropDatabase(t, adminDSN, name) })

	baselinePath, err := BaselinePath()
	if err != nil {
		t.Fatal(err)
	}
	baseline, err := os.ReadFile(baselinePath)
	if err != nil {
		t.Fatalf("read %s: %v", baselinePath, err)
	}

	db, err := sql.Open("postgres", scratchDSN)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	if _, err := db.Exec(string(baseline)); err != nil {
		t.Fatalf("apply %s: %v", baselinePath, err)
	}

	return db
}

// BaselinePath is the absolute path of tools/db/baseline.sql, the schema every
// scratch database is built from.
//
// `go test` runs each test binary in the directory of its package, so the
// module root is the nearest directory above it holding go.mod. The location
// of this source file cannot be used instead: -trimpath rewrites it to an
// import path, which names no directory on disk.
func BaselinePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("working directory: %w", err)
	}
	root, err := moduleRoot(dir)
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "tools", "db", "baseline.sql"), nil
}

// moduleRoot walks up from dir to the first directory holding go.mod.
func moduleRoot(dir string) (string, error) {
	for start := dir; ; {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no go.mod in %s or above it: run the tests from inside the module", start)
		}
		dir = parent
	}
}

// dropDatabase only reports a failure: the next run drops the database
// before creating it again, so a leftover never poisons it.
func dropDatabase(t testing.TB, adminDSN, name string) {
	admin, err := sql.Open("postgres", adminDSN)
	if err != nil {
		t.Logf("scratch database %s left behind: %v", name, err)
		return
	}
	defer admin.Close()
	if _, err := admin.Exec(`DROP DATABASE IF EXISTS ` + name); err != nil {
		t.Logf("scratch database %s left behind: %v", name, err)
	}
}

// scratchDSNs turns the configured URL into one on the maintenance database and
// one on the scratch database, keeping credentials, host and sslmode.
func scratchDSNs(dsn, name string) (adminDSN, scratchDSN string, err error) {
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
	scratch.Path = "/" + name

	return admin.String(), scratch.String(), nil
}

// DialNode connects to the node named by ETH_TEST_RPC, or DefaultNodeURL. The
// client is closed when the test ends, after the cleanups registered later,
// so a test that wires it into a global can restore that global first.
func DialNode(ctx context.Context, t testing.TB) *ethclient.Client {
	t.Helper()

	rpc := os.Getenv("ETH_TEST_RPC")
	if rpc == "" {
		rpc = DefaultNodeURL
	}
	client, err := ethclient.DialContext(ctx, rpc)
	if err != nil {
		Unavailable(t, "no local EVM node at %s: %v", rpc, err)
	}
	t.Cleanup(client.Close)
	if _, err := client.ChainID(ctx); err != nil {
		Unavailable(t, "local EVM node at %s not reachable: %v", rpc, err)
	}
	return client
}

// Unavailable skips the test when no integration run was asked for, and fails
// it when one was: with TEST_DATABASE_URL or ETH_TEST_RPC set, a missing
// database or node is a broken run, and a skip would let it pass unnoticed.
func Unavailable(t testing.TB, format string, args ...any) {
	t.Helper()
	if os.Getenv("TEST_DATABASE_URL") != "" || os.Getenv("ETH_TEST_RPC") != "" {
		t.Fatalf(format, args...)
	}
	t.Skipf(format, args...)
}
