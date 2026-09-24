//go:build integration

// Coverage of the invariant the whole product rests on: no asset can be on sale
// without a token behind it.
//
// The rule was stated in the design, implemented in Go and covered by the
// handler suite — and absent from the database until migration 000021 added
// real_estate_active_requires_token_ck. Until then a manual UPDATE, a future
// endpoint or a service bug could put an asset back on sale with nothing behind
// it, silently.
//
// 000021 added the constraint NOT VALID, so the assets inherited from before
// tokenising publication would not block it, and 000022 validated it once they
// were cleared. Both halves matter and neither shows up in Go code, which is
// why they are checked here rather than trusted:
//
//   - the constraint rejects new writes (what NOT VALID already did);
//   - it is marked valid, meaning the rows already stored were read and comply.
//
// Run with the disposable database used by the rest of the integration suite:
//
//	TEST_DATABASE_URL="postgres://…" go test -tags integration ./internal/assets_managements/database/...
package database

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

const activeRequiresTokenConstraint = "real_estate_active_requires_token_ck"

func openSchemaProbeDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	return db
}

// A constraint left NOT VALID is often misread as "disabled". It is not: it is
// enforced on every INSERT and UPDATE from the moment it exists, and only skips
// re-reading the rows already stored. What 000022 adds is that second half.
//
// The distinction is invisible in any query result — it lives in
// pg_constraint.convalidated — so it can be lost by a rebuild without anything
// failing. That is precisely how it was nearly lost: the VALIDATE had been run
// by hand and existed in no file, so every database rebuilt from the repository
// came back unvalidated.
func TestActiveRequiresTokenIsValidated(t *testing.T) {
	db := openSchemaProbeDB(t)

	var convalidated bool
	err := db.QueryRow(`
		SELECT convalidated FROM pg_constraint
		WHERE conname = $1`, activeRequiresTokenConstraint).Scan(&convalidated)
	if err == sql.ErrNoRows {
		t.Fatalf("%s does not exist: the invariant is not in the schema at all",
			activeRequiresTokenConstraint)
	}
	if err != nil {
		t.Fatal(err)
	}

	if !convalidated {
		t.Fatalf("%s exists but is NOT VALID: new writes are checked, "+
			"the rows already stored are not", activeRequiresTokenConstraint)
	}
}

// The constraint has to refuse, not merely exist. This writes an asset the way
// a stray UPDATE would — active, with no token — and expects the database to
// say no on its own, without any Go code being involved.
//
// The transaction is always rolled back: the point is the refusal, not the row.
func TestActiveAssetWithoutTokenIsRejected(t *testing.T) {
	db := openSchemaProbeDB(t)

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	estateType, issuer := ensurePrerequisites(t, tx)

	// status_id 3 is 'published', which the trigger projects to active = true.
	// Going through the status rather than writing active directly is what a
	// real publication does — and writing active has no effect anyway, the
	// trigger recomputes it.
	var id int64
	err = tx.QueryRow(`
		INSERT INTO ass.real_estate (title, estate_type, issuer_id, status_id, token_id)
		VALUES ('constraint probe', $1, $2, 3, NULL)
		RETURNING id`, estateType, issuer).Scan(&id)

	if err == nil {
		t.Fatal("an asset was published with no token behind it: " +
			"the invariant is not enforced by the database")
	}
	if !strings.Contains(err.Error(), activeRequiresTokenConstraint) {
		t.Fatalf("the insert failed, but not on the invariant: %v", err)
	}
}

// The mirror assertion, and the one that keeps the test above honest. A
// constraint that refused everything would make the first test pass just as
// well; this proves the refusal is about the missing token and nothing else.
//
// A draft is not on sale, so it is allowed to have no token — that is the whole
// point of being able to prepare an asset before tokenising it.
func TestDraftWithoutTokenIsAccepted(t *testing.T) {
	db := openSchemaProbeDB(t)

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	estateType, issuer := ensurePrerequisites(t, tx)

	// status_id 1 is 'draft', projected to active = false.
	var id int64
	if err := tx.QueryRow(`
		INSERT INTO ass.real_estate (title, estate_type, issuer_id, status_id, token_id)
		VALUES ('constraint probe draft', $1, $2, 1, NULL)
		RETURNING id`, estateType, issuer).Scan(&id); err != nil {
		t.Fatalf("a draft without a token must be accepted: %v", err)
	}

	var active bool
	if err := tx.QueryRow(
		`SELECT active FROM ass.real_estate WHERE id = $1`, id).Scan(&active); err != nil {
		t.Fatal(err)
	}
	if active {
		t.Fatal("a draft must not be active: the projection trigger is not doing its job")
	}
}

// An asset needs a type and an issuer to exist at all, and neither is something
// these tests are about. They are created inside the caller's transaction when
// the database does not already carry them, so the tests run on a database
// rebuilt from the baseline as well as on one the handler fixtures have seeded —
// and nothing survives the rollback either way.
//
// This is not over-engineering: tools/db/baseline.sql seeds three reference
// tables and ass.real_estate_type is not among them, so a freshly rebuilt
// database has none.
func ensurePrerequisites(t *testing.T, tx *sql.Tx) (estateType, issuer int64) {
	t.Helper()

	if err := tx.QueryRow(`SELECT id FROM ass.real_estate_type LIMIT 1`).Scan(&estateType); err != nil {
		if err != sql.ErrNoRows {
			t.Fatal(err)
		}
		if err := tx.QueryRow(`
			INSERT INTO ass.real_estate_type (name) VALUES ('constraint probe type')
			RETURNING id`).Scan(&estateType); err != nil {
			t.Fatalf("could not provide an estate type: %v", err)
		}
	}

	if err := tx.QueryRow(`SELECT id FROM ass.issuer LIMIT 1`).Scan(&issuer); err != nil {
		if err != sql.ErrNoRows {
			t.Fatal(err)
		}
		var status int64
		if err := tx.QueryRow(`SELECT id FROM ass.issuer_status LIMIT 1`).Scan(&status); err != nil {
			t.Fatalf("no issuer status in the reference table: %v", err)
		}
		if err := tx.QueryRow(`
			INSERT INTO ass.issuer (name, legal_form, country_code, status_id)
			VALUES ('Constraint Probe', 'SA', 'LU', $1)
			RETURNING id`, status).Scan(&issuer); err != nil {
			t.Fatalf("could not provide an issuer: %v", err)
		}
	}

	return estateType, issuer
}
