//go:build integration

package services

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/globals"
	_ "github.com/lib/pq"
)

// The race itself cannot be timed from a test: a second manager would have to
// act between the guard read and the UPDATE. What can be checked is the answer
// given once it has happened, which is the only part that has a contract.
func TestAfterLostFundraisingRace(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	globals.DB = db
	ctx := context.Background()

	mustExec := func(t *testing.T, query string, args ...any) {
		t.Helper()
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("fixture failed: %v", err)
		}
	}
	mustScanID := func(t *testing.T, query string, args ...any) int {
		t.Helper()
		var id int
		if err := db.QueryRow(query, args...).Scan(&id); err != nil {
			t.Fatalf("fixture failed: %v", err)
		}
		return id
	}

	mustExec(t, `INSERT INTO ass.real_estate_type (id, name) VALUES (1, 'Apartment') ON CONFLICT (id) DO NOTHING`)
	issuerID := mustScanID(t, `
	    INSERT INTO ass.issuer (name, legal_form, status_id)
	    VALUES ('Fundraising race probe', 'SA', 2) RETURNING id`)

	// Assets are deleted before their tokens (real_estate_token_id_fkey is
	// RESTRICT) and both before the issuer, in one cleanup so the order holds.
	var assets, tokens []int
	trackAsset := func(id int) int {
		assets = append(assets, id)
		return id
	}
	t.Cleanup(func() {
		for _, id := range assets {
			_, _ = db.Exec(`DELETE FROM ass.real_estate WHERE id = $1`, id)
		}
		for _, id := range tokens {
			_, _ = db.Exec(`DELETE FROM blk.token WHERE id = $1`, id)
		}
		_, _ = db.Exec(`DELETE FROM ass.issuer WHERE id = $1`, issuerID)
	})

	assetIn := func(t *testing.T, status int, deleted bool) int {
		t.Helper()
		// A deleted asset carries status 7 and no token
		// (real_estate_deleted_status_ck, real_estate_deleted_not_tokenized_ck).
		if deleted {
			return trackAsset(mustScanID(t, `
			    INSERT INTO ass.real_estate (title, estate_type, issuer_id, status_id, deleted_at)
			    VALUES ('Race probe', 1, $1, 7, now()) RETURNING id`, issuerID))
		}
		tokenID := mustScanID(t, `
		    INSERT INTO blk.token (address, token_name, symbol, nb_decimal)
		    VALUES ('0x' || lpad(md5(random()::text), 40, '0'), 'race-' || md5(random()::text), 'RACE', 0)
		    RETURNING id`)
		tokens = append(tokens, tokenID)
		return trackAsset(mustScanID(t, `
		    INSERT INTO ass.real_estate (title, estate_type, issuer_id, status_id, token_id)
		    VALUES ('Race probe', 1, $1, $2, $3) RETURNING id`, issuerID, status, tokenID))
	}

	s := &Service{}

	t.Run("another manager opened it first: 200 with the asset", func(t *testing.T) {
		id := assetIn(t, 4, false)
		re, err := s.afterLostFundraisingRace(ctx, id)
		if err != nil {
			t.Fatalf("want no error, got %v", err)
		}
		if re.Status == nil || *re.Status != "fundraising" {
			t.Errorf("want fundraising, got %v", re.Status)
		}
	})

	t.Run("it was funded meanwhile: 409", func(t *testing.T) {
		id := assetIn(t, 5, false)
		if _, err := s.afterLostFundraisingRace(ctx, id); !errors.Is(err, ErrRealEstateConflict) {
			t.Fatalf("want a conflict, got %v", err)
		}
	})

	t.Run("it was deleted meanwhile: 404", func(t *testing.T) {
		id := assetIn(t, 0, true)
		if _, err := s.afterLostFundraisingRace(ctx, id); !errors.Is(err, ErrRealEstateNotFound) {
			t.Fatalf("want not found, got %v", err)
		}
	})
}
