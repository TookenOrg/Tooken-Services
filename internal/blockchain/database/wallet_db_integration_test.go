//go:build integration

// Coverage of blk.user_wallet against a real PostgreSQL.
//
// This file exists because of two defects that shipped through a compiler, a
// vet run and the whole integration suite without being noticed, and were only
// caught by hand:
//
//   - sql.ErrNoRows was wrapped as a failure, so "this user has no wallet yet" —
//     the nominal case when an identity is created — was reported as a database
//     error, and CreateIdentity refused every new investor;
//   - the address was scanned straight into a common.Address. That type does
//     implement sql.Scanner, which makes the code look right, but its Scan wants
//     20 raw bytes while the column is a varchar holding "0x…". Every read
//     failed with "can't scan string into Address".
//
// Neither is visible without a database. Hence these tests.
//
// Run them with the same disposable database as the handler suite:
//
//	TEST_DATABASE_URL="postgres://…" go test -tags integration ./internal/blockchain/database/...
package database

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/lib/pq"
)

// Probe rows live on user ids no fixture uses, and carry a label that makes
// them recognisable and easy to remove. blk.user_wallet has no foreign key on
// user_id, so these ids need not exist in usr.users.
const (
	probeUser      = 900001
	probeOtherUser = 900002
	probeLabel     = "wallet-db-probe"
)

func openProbeDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	globals.DB = db

	// Clean before as well as after: a suite interrupted mid-run must not make
	// the next one fail. The repository expects two consecutive runs on a
	// database that is never emptied to give the same result.
	cleanProbeRows(t, db)
	t.Cleanup(func() {
		cleanProbeRows(t, db)
		db.Close()
	})

	return db
}

func cleanProbeRows(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(`DELETE FROM blk.user_wallet WHERE label = $1`, probeLabel); err != nil {
		t.Fatal(err)
	}
}

func insertProbeWallet(t *testing.T, db *sql.DB, userID int, address string, active bool) {
	t.Helper()
	_, err := db.Exec(`
		INSERT INTO blk.user_wallet (user_id, wallet_address, label, is_active)
		VALUES ($1, $2, $3, $4)`,
		userID, address, probeLabel, active)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetWalletByUserId(t *testing.T) {
	db := openProbeDB(t)
	ctx := context.Background()

	// The case that lets through, and the reason this test matters most.
	//
	// An investor who has no wallet is not an anomaly: it is what every identity
	// creation starts from. The function must say "nothing here" without saying
	// "something went wrong", because isIdempotentIdentity treats any error as a
	// refusal. Written on its own, a guard that always refuses would look just
	// as green as a correct one.
	t.Run("a user without a wallet is not an error", func(t *testing.T) {
		got, err := GetWalletByUserId(ctx, probeUser)
		if err != nil {
			t.Fatalf("no wallet must not be an error, got %v", err)
		}
		if got != nil {
			t.Fatalf("expected no address, got %s", got.Hex())
		}
	})

	// The case that catches the scan defect. The assertion is on the value, not
	// on the absence of error: a read that returns the wrong address is worse
	// than one that fails.
	t.Run("an active wallet is returned, converted from its text column", func(t *testing.T) {
		const address = "0x5272DeDe69BD79fDCb8d66ad173f15dF984D4ED9"
		insertProbeWallet(t, db, probeUser, address, true)

		got, err := GetWalletByUserId(ctx, probeUser)
		if err != nil {
			t.Fatal(err)
		}
		if got == nil {
			t.Fatal("expected an address, got none")
		}
		if *got != common.HexToAddress(address) {
			t.Fatalf("got %s, want %s", got.Hex(), address)
		}

		cleanProbeRows(t, db)
	})

	// The column is plain text, so nothing forces a casing. An address written
	// in lower case is the same address; if the read ever stopped going through
	// HexToAddress, this is where it would show.
	t.Run("casing in the column does not change the address", func(t *testing.T) {
		const canonical = "0x5272DeDe69BD79fDCb8d66ad173f15dF984D4ED9"
		insertProbeWallet(t, db, probeUser, strings.ToLower(canonical), true)

		got, err := GetWalletByUserId(ctx, probeUser)
		if err != nil {
			t.Fatal(err)
		}
		if got == nil || *got != common.HexToAddress(canonical) {
			t.Fatalf("got %v, want %s", got, canonical)
		}

		cleanProbeRows(t, db)
	})

	// D7 allows only one active wallet per investor, so a deactivated one must
	// not be mistaken for a live one. Without this filter a retired wallet would
	// keep an investor from ever getting a new identity.
	t.Run("a deactivated wallet is invisible", func(t *testing.T) {
		insertProbeWallet(t, db, probeUser, "0x1111111111111111111111111111111111111111", false)

		got, err := GetWalletByUserId(ctx, probeUser)
		if err != nil {
			t.Fatal(err)
		}
		if got != nil {
			t.Fatalf("a deactivated wallet must not be returned, got %s", got.Hex())
		}

		cleanProbeRows(t, db)
	})

	// A missing WHERE on user_id would still pass every test above, since they
	// only ever hold one row. This one fails if the query stops discriminating.
	t.Run("another investor's wallet is never returned", func(t *testing.T) {
		insertProbeWallet(t, db, probeOtherUser, "0x2222222222222222222222222222222222222222", true)

		got, err := GetWalletByUserId(ctx, probeUser)
		if err != nil {
			t.Fatal(err)
		}
		if got != nil {
			t.Fatalf("expected nothing for user %d, got %s", probeUser, got.Hex())
		}

		cleanProbeRows(t, db)
	})
}

func TestInsertWallet(t *testing.T) {
	db := openProbeDB(t)
	ctx := context.Background()

	// Reading back through GetWalletByUserId rather than through a raw SELECT is
	// deliberate: it proves the two functions agree on the format written into
	// the column. A writer storing raw bytes and a reader expecting hex would
	// both look correct in isolation.
	t.Run("what is written is what is read", func(t *testing.T) {
		const address = "0x3333333333333333333333333333333333333333"

		id, err := InsertWallet(probeUser, address, probeLabel)
		if err != nil {
			t.Fatal(err)
		}
		if id == 0 {
			t.Fatal("expected the generated id to be returned")
		}

		got, err := GetWalletByUserId(ctx, probeUser)
		if err != nil {
			t.Fatal(err)
		}
		if got == nil || *got != common.HexToAddress(address) {
			t.Fatalf("got %v, want %s", got, address)
		}

		cleanProbeRows(t, db)
	})

	// A wallet is born usable. If the default ever changed, GetWalletByUserId
	// would stop seeing freshly created wallets and every second identity
	// creation would succeed instead of being refused.
	t.Run("a new wallet is active", func(t *testing.T) {
		const address = "0x4444444444444444444444444444444444444444"

		if _, err := InsertWallet(probeUser, address, probeLabel); err != nil {
			t.Fatal(err)
		}

		var active bool
		if err := db.QueryRow(
			`SELECT is_active FROM blk.user_wallet WHERE wallet_address = $1`,
			address).Scan(&active); err != nil {
			t.Fatal(err)
		}
		if !active {
			t.Fatal("a newly inserted wallet must be active")
		}

		cleanProbeRows(t, db)
	})
}

// Migration 000023 removed every key column from blk.user_wallet, because the
// platform never needs an investor's private key: mint, burn and forcedTransfer
// are signed by the token agent, and a wallet whose key is gone is recovered
// through the investor's ONCHAINID.
//
// The decision only holds as long as there is nowhere to write a key. A column
// that exists for a key eventually receives one, so the absence is the thing
// worth testing — not the code that happens not to write there today.
func TestUserWalletHoldsNoPrivateKey(t *testing.T) {
	db := openProbeDB(t)

	rows, err := db.Query(`
		SELECT column_name
		FROM information_schema.columns
		WHERE table_schema = 'blk' AND table_name = 'user_wallet'
		  AND column_name ILIKE '%private%key%'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var found []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		found = append(found, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}

	if len(found) > 0 {
		t.Fatalf("blk.user_wallet must hold no private key, found: %s",
			strings.Join(found, ", "))
	}
}
