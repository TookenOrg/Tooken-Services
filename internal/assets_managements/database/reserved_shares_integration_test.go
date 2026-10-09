//go:build integration

// The three reads of reserved shares must agree (TICKET-M3-3 §2.5): the stock
// an order is checked against, tokens_sold shown on the asset, and the guard
// state a manager's write is checked against. Two of them once counted lapsed
// reservations and the third did not: the order passed while the asset showed
// more shares sold than it holds.
//
// The orders written here can never be deleted (the status history is
// append-only): run it on a THROWAWAY database built from tools/db/baseline.sql.
package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/globals"
)

func TestReservedSharesHaveOneDefinition(t *testing.T) {
	db := openSchemaProbeDB(t)
	globals.DB = db
	ctx := context.Background()
	run := time.Now().UnixNano()
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}

	var issuer, token, asset, user int
	must(db.QueryRow(`INSERT INTO ass.issuer (name, legal_form) VALUES ($1, 'SA') RETURNING id`,
		fmt.Sprintf("reserved-%d", run)).Scan(&issuer))
	must(db.QueryRow(`
	    INSERT INTO blk.token (symbol, token_name, address, nb_decimal)
	    VALUES ('RSV', $1::text, '0x' || substr(md5($1::text), 1, 32) || '00000000', 0)
	    RETURNING id`, fmt.Sprintf("reserved-%d", run)).Scan(&token))
	must(db.QueryRow(`
	    INSERT INTO ass.real_estate (title, estate_type, issuer_id, status_id, token_id)
	    SELECT $1, min(id), $2, 4, $3 FROM ass.real_estate_type
	    RETURNING id`, fmt.Sprintf("Villa Belair %d", run), issuer, token).Scan(&asset))
	_, err := db.Exec(`
	    INSERT INTO ass.real_estate_shares_config (real_estate_id, total_shares, price_per_share)
	    VALUES ($1, 1500, 200)`, asset)
	must(err)
	must(db.QueryRow(`
	    INSERT INTO usr.users (full_name, email, password) VALUES ('Reserved', $1, 'x')
	    RETURNING id`, fmt.Sprintf("reserved-%d@tooken.test", run)).Scan(&user))

	// One order per kind of status. Only the marked ones hold shares.
	orders := []struct {
		quantity  int
		status    int
		expiresIn string
		holds     bool
	}{
		{1, 2, "10 minutes", true},    // AWAITING_PAYMENT, reservation running
		{1000, 2, "-1 second", false}, // AWAITING_PAYMENT, reservation lapsed
		{20, 4, "-1 hour", true},      // PAYMENT_CONFIRMED: paid, the date no longer matters
		{300, 5, "-1 hour", true},     // CLOSED (U8)
		{40, 6, "-1 hour", false},     // CANCELLED
		{50, 7, "-1 hour", false},     // EXPIRED
	}
	want := 0
	for i, o := range orders {
		ref := fmt.Sprintf("RSV-%d-%d", run, i)
		_, err := db.Exec(`
		    INSERT INTO iss.issuance_orders (
		        user_id, asset_id, quantity, status_id, order_reference,
		        unit_price, currency_code, gross_amount, fee_amount, amount_due,
		        created_at, reservation_expires_at, delivery_tx_hash)
		    VALUES (
		        $1, $2, $3::bigint, $4::int, $5::text,
		        1, 'EUR', $3::numeric, 0, $3::numeric,
		        now() - interval '2 hours', now() + $6::interval,
		        CASE WHEN $4::int IN (5, 10, 11, 12) THEN '0x' || md5($5::text) || md5($5::text || '#') END)`,
			user, asset, o.quantity, o.status, ref, o.expiresIn)
		if err != nil {
			t.Fatalf("order %d (status %d): %v", i, o.status, err)
		}
		if o.holds {
			want += o.quantity
		}
	}

	t.Run("the shared definition", func(t *testing.T) {
		var got int
		must(db.QueryRow(`SELECT COALESCE(SUM(o.quantity), 0)`+ReservedOrders+` AND o.asset_id = $1`, asset).Scan(&got))
		if got != want {
			t.Fatalf("ReservedOrders sums %d, want %d", got, want)
		}
	})

	t.Run("tokens_sold on the asset", func(t *testing.T) {
		re, err := GetRealEstateById(ctx, asset, true)
		must(err)
		if re.Progression.TokensSold != int64(want) {
			t.Fatalf("tokens_sold = %d, want %d", re.Progression.TokensSold, want)
		}
	})

	t.Run("the guard state of a manager's write", func(t *testing.T) {
		s, err := GetRealEstateGuardState(ctx, asset)
		must(err)
		if s.ReservedShares != int64(want) {
			t.Fatalf("reserved = %d, want %d", s.ReservedShares, want)
		}
	})
}
