//go:build integration

// Coverage of the payment service against a real PostgreSQL: the translation
// of a missing order into a 404, the propagation of the refusals, and the
// order read back after the confirmation. What the transaction does is covered
// by internal/payments/database.
//
// Writes committed, append-only rows: run on a THROWAWAY database.
//
//	TEST_DATABASE_URL="postgres://...?sslmode=disable" \
//	    go test -tags integration -race ./internal/payments/services/
package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/globals"
	orderDB "github.com/TookenOrg/tooken-services/internal/orders/database"
	"github.com/TookenOrg/tooken-services/internal/payments/database"
	_ "github.com/lib/pq"
)

var (
	svcRun = time.Now().UnixNano() % 1_000_000_000
	svcSeq atomic.Int64
)

func svcUnique(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, svcRun, svcSeq.Add(1))
}

func openServiceDB(t *testing.T) *sql.DB {
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
	previous := globals.DB
	globals.DB = db
	t.Cleanup(func() { globals.DB = previous })
	return db
}

type svcFixture struct {
	admin    int
	investor int
	asset    int
}

func newSvcFixture(t *testing.T, db *sql.DB) svcFixture {
	t.Helper()
	var f svcFixture
	for _, u := range []struct {
		role string
		id   *int
	}{{"ADMIN", &f.admin}, {"USER", &f.investor}} {
		if err := db.QueryRow(`
		    INSERT INTO usr.users (full_name, email, password, role)
		    VALUES ('Service probe', $1, 'x', $2) RETURNING id`,
			svcUnique("svc")+"@probe.local", u.role).Scan(u.id); err != nil {
			t.Fatalf("creating user: %v", err)
		}
	}

	var issuer, token int
	if err := db.QueryRow(`INSERT INTO ass.issuer (name, legal_form) VALUES ($1, 'SA') RETURNING id`,
		svcUnique("issuer")).Scan(&issuer); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`
	    INSERT INTO blk.token (symbol, token_name, address, nb_decimal)
	    VALUES ('SVC', $1::text, '0x' || substr(md5($1::text), 1, 32) || '00000000', 0)
	    RETURNING id`, svcUnique("token")).Scan(&token); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`
	    INSERT INTO ass.real_estate (title, estate_type, issuer_id, status_id, token_id)
	    SELECT $1, min(id), $2, 4, $3 FROM ass.real_estate_type
	    RETURNING id`, svcUnique("Villa Belair"), issuer, token).Scan(&f.asset); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`
	    INSERT INTO ass.real_estate_shares_config
	        (real_estate_id, total_shares, price_per_share, entry_fee_rate, max_shares_per_investor, currency_code)
	    VALUES ($1, 1500, 200, 2, 50, 'EUR')`, f.asset); err != nil {
		t.Fatal(err)
	}
	return f
}

// order writes an order of 10 shares (amount_due 2 040.00 EUR) and returns
// its reference.
func (f svcFixture) order(t *testing.T, db *sql.DB, status int, expiresIn string) string {
	t.Helper()
	ref := svcUnique("ISS")
	if _, err := db.Exec(`
	    INSERT INTO iss.issuance_orders (
	        user_id, asset_id, quantity, status_id, order_reference,
	        unit_price, currency_code, gross_amount, fee_amount, amount_due,
	        created_at, reservation_expires_at)
	    VALUES ($1, $2, 10, $3, $4, 200, 'EUR', 2000, 40, 2040,
	            now() - interval '1 hour', now() + $5::interval)`,
		f.investor, f.asset, status, ref, expiresIn); err != nil {
		t.Fatalf("creating order: %v", err)
	}
	return ref
}

func TestConfirmIssuanceOrderPaymentService(t *testing.T) {
	ctx := context.Background()
	db := openServiceDB(t)
	svc := NewService()

	t.Run("confirmed: the order returned is the one after the payment", func(t *testing.T) {
		f := newSvcFixture(t, db)
		ref := f.order(t, db, 2, "15 minutes")

		order, err := svc.ConfirmIssuanceOrderPayment(ctx, ref, f.admin)
		if err != nil {
			t.Fatalf("error = %v", err)
		}
		if order.OrderRef != ref || order.StatusCode != "PAYMENT_CONFIRMED" ||
			order.StatusId == nil || *order.StatusId != 4 {
			t.Fatalf("order = %s %s (%v), want %s PAYMENT_CONFIRMED (4)", order.OrderRef, order.StatusCode, order.StatusId, ref)
		}
		if order.UserId != f.investor || order.AmountDue != "2040.00" {
			t.Fatalf("order of user %d for %s, want user %d for 2040.00", order.UserId, order.AmountDue, f.investor)
		}
	})

	t.Run("already paid: success, order returned as is", func(t *testing.T) {
		f := newSvcFixture(t, db)
		ref := f.order(t, db, 4, "15 minutes")

		order, err := svc.ConfirmIssuanceOrderPayment(ctx, ref, f.admin)
		if err != nil {
			t.Fatalf("error = %v", err)
		}
		if order.StatusCode != "PAYMENT_CONFIRMED" {
			t.Fatalf("status = %s, want PAYMENT_CONFIRMED", order.StatusCode)
		}
	})

	t.Run("unknown reference: ErrOrderNotFound", func(t *testing.T) {
		newSvcFixture(t, db)
		_, err := svc.ConfirmIssuanceOrderPayment(ctx, svcUnique("ISS-UNKNOWN"), 1)
		if !errors.Is(err, orderDB.ErrOrderNotFound) || !errors.Is(err, database.ErrOrderNotFound) {
			t.Fatalf("error = %v, want ErrOrderNotFound", err)
		}
	})

	for _, c := range []struct {
		name      string
		status    int
		expiresIn string
		want      error
	}{
		{"lapsed reservation: ErrOrderReservationExpired passes through", 2, "-1 minute", database.ErrOrderReservationExpired},
		{"cancelled order: ErrOrderNotAwaitingPayment passes through", 6, "15 minutes", database.ErrOrderNotAwaitingPayment},
	} {
		t.Run(c.name, func(t *testing.T) {
			f := newSvcFixture(t, db)
			ref := f.order(t, db, c.status, c.expiresIn)
			_, err := svc.ConfirmIssuanceOrderPayment(ctx, ref, f.admin)
			if !errors.Is(err, c.want) {
				t.Fatalf("error = %v, want %v", err, c.want)
			}
		})
	}

	// Only sql.ErrNoRows is a 404: a broken database must stay a 500.
	t.Run("technical error on the first read: not masked as a 404", func(t *testing.T) {
		closed, err := sql.Open("postgres", os.Getenv("TEST_DATABASE_URL"))
		if err != nil {
			t.Fatal(err)
		}
		closed.Close()
		previous := globals.DB
		globals.DB = closed
		defer func() { globals.DB = previous }()

		_, err = svc.ConfirmIssuanceOrderPayment(ctx, "ISS-ANY", 1)
		if err == nil || errors.Is(err, orderDB.ErrOrderNotFound) {
			t.Fatalf("error = %v, want a technical error", err)
		}
	})
}
