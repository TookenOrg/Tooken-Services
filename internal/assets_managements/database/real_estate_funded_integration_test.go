//go:build integration

// Marking an asset funded is a guarded status move: only a live asset in
// fundraising goes to funded, and the move happens at most once.
//
// Run it on a THROWAWAY database built from tools/db/baseline.sql.
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/globals"
)

func seedAssetWithStatus(t *testing.T, db *sql.DB, status int) int {
	t.Helper()

	run := time.Now().UnixNano()
	var issuer, token, asset int
	if err := db.QueryRow(`INSERT INTO ass.issuer (name, legal_form) VALUES ($1, 'SA') RETURNING id`,
		fmt.Sprintf("funded-%d", run)).Scan(&issuer); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`
	    INSERT INTO blk.token (symbol, token_name, address, nb_decimal)
	    VALUES ('FND', $1::text, '0x' || substr(md5($1::text), 1, 32) || '00000000', 0)
	    RETURNING id`, fmt.Sprintf("funded-%d", run)).Scan(&token); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`
	    INSERT INTO ass.real_estate (title, estate_type, issuer_id, status_id, token_id)
	    SELECT $1, min(id), $2, $3, $4 FROM ass.real_estate_type
	    RETURNING id`, fmt.Sprintf("Villa Belair %d", run), issuer, status, token).Scan(&asset); err != nil {
		t.Fatal(err)
	}
	return asset
}

func assetStatus(t *testing.T, db *sql.DB, id int) int {
	t.Helper()

	var status int
	if err := db.QueryRow(`SELECT status_id FROM ass.real_estate WHERE id = $1`, id).Scan(&status); err != nil {
		t.Fatal(err)
	}
	return status
}

func TestMarkRealEstateFundedMovesOnlyAFundraisingAsset(t *testing.T) {
	db := openSchemaProbeDB(t)
	globals.DB = db
	ctx := context.Background()

	for _, tc := range []struct {
		name   string
		status int
		want   error
		after  int
	}{
		{"fundraising becomes funded", StatusFundraising, nil, StatusFunded},
		{"already funded", StatusFunded, sql.ErrNoRows, StatusFunded},
		{"not yet open", 3, sql.ErrNoRows, 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			id := seedAssetWithStatus(t, db, tc.status)
			if err := MarkRealEstateFunded(ctx, id); !errors.Is(err, tc.want) {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
			if got := assetStatus(t, db, id); got != tc.after {
				t.Fatalf("status = %d, want %d", got, tc.after)
			}
		})
	}

	t.Run("two calls at once", func(t *testing.T) {
		id := seedAssetWithStatus(t, db, StatusFundraising)
		errs := make([]error, 2)
		var wg sync.WaitGroup
		for i := range errs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				errs[i] = MarkRealEstateFunded(ctx, id)
			}()
		}
		wg.Wait()

		won := 0
		for _, err := range errs {
			if err == nil {
				won++
			} else if !errors.Is(err, sql.ErrNoRows) {
				t.Fatalf("unexpected error: %v", err)
			}
		}
		if won != 1 {
			t.Fatalf("%d calls moved the asset, want exactly 1", won)
		}
	})
}
