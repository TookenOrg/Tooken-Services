//go:build integration

// Schema-level coverage of the MANUAL and ONCHAIN payment modes on
// iss.issuance_order_payments against a real PostgreSQL built from baseline.sql.
package database

import (
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/testenv"
	"github.com/lib/pq"
)

const (
	paymentCheckViolation  = "23514"
	paymentUniqueViolation = "23505"
)

type paymentModeFixture struct {
	db                     *sql.DB
	admin, investor, asset int64
}

func newPaymentModeFixture(t *testing.T) paymentModeFixture {
	t.Helper()

	db := testenv.ScratchDatabase(t, "tooken_it_payment_mode")
	f := paymentModeFixture{db: db}

	if err := db.QueryRow(`
	    INSERT INTO usr.users (full_name, email, password, role)
	    VALUES ('Payment mode admin', 'admin@payment-mode.local', 'not-a-real-hash', 'ADMIN')
	    RETURNING id`).Scan(&f.admin); err != nil {
		t.Fatalf("creating admin user: %v", err)
	}

	if err := db.QueryRow(`
	    INSERT INTO usr.users (full_name, email, password, role)
	    VALUES ('Payment mode investor', 'investor@payment-mode.local', 'not-a-real-hash', 'USER')
	    RETURNING id`).Scan(&f.investor); err != nil {
		t.Fatalf("creating investor user: %v", err)
	}

	var issuer int64
	if err := db.QueryRow(`
	    INSERT INTO ass.issuer (name, legal_form)
	    VALUES ('Villa Belair SAS', 'SAS')
	    RETURNING id`).Scan(&issuer); err != nil {
		t.Fatalf("creating issuer: %v", err)
	}

	if err := db.QueryRow(`
	    INSERT INTO ass.real_estate (title, estate_type, issuer_id)
	    SELECT 'Tooken Villa Belair #12', min(id), $1 FROM ass.real_estate_type
	    RETURNING id`, issuer).Scan(&f.asset); err != nil {
		t.Fatalf("creating asset: %v", err)
	}

	if _, err := db.Exec(`
	    INSERT INTO ass.real_estate_shares_config (real_estate_id, total_shares, price_per_share)
	    VALUES ($1, 1000, 200)`, f.asset); err != nil {
		t.Fatalf("creating shares config: %v", err)
	}

	return f
}

func (f paymentModeFixture) newOrder(t *testing.T, reference string) int64 {
	t.Helper()

	var id int64
	if err := f.db.QueryRow(`
	    INSERT INTO iss.issuance_orders
	        (user_id, asset_id, quantity, status_id, order_reference,
	         unit_price, currency_code, gross_amount, fee_amount, amount_due, reservation_expires_at)
	    VALUES ($1, $2, 10, 2, $3, 200.00000000, 'EUR', 2000.00, 40.00, 2040.00, now() + interval '15 minutes')
	    RETURNING id`, f.investor, f.asset, reference).Scan(&id); err != nil {
		t.Fatalf("creating order %s: %v", reference, err)
	}
	return id
}

func expectPaymentRefused(t *testing.T, db *sql.DB, wantCode, wantConstraint, query string, args ...any) {
	t.Helper()

	_, err := db.Exec(query, args...)
	if err == nil {
		t.Fatalf("want refused (%s %s), got accepted\n%s", wantCode, wantConstraint, query)
	}
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		t.Fatalf("want a PostgreSQL error, got %v", err)
	}
	if string(pqErr.Code) != wantCode {
		t.Fatalf("want code %s, got %s (%s)", wantCode, pqErr.Code, pqErr.Message)
	}
	if wantConstraint != "" && pqErr.Constraint != wantConstraint {
		t.Fatalf("want constraint %s, got %q (%s)", wantConstraint, pqErr.Constraint, pqErr.Message)
	}
}

func mustExecPayment(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()
	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("want accepted, got %v\n%s", err, query)
	}
}

func TestPaymentMode(t *testing.T) {
	f := newPaymentModeFixture(t)

	const insertPaymentRow = `
	    INSERT INTO iss.issuance_order_payments
	        (order_id, kind, provider, amount, currency_code, recorded_by_user_id, tx_hash)
	    VALUES ($1, $2, $3, 2040.00, 'EUR', $4, $5)`

	var (
		hashA     = "0x" + strings.Repeat("a", 64)
		hashB     = "0x" + strings.Repeat("b", 64)
		hashC     = "0x" + strings.Repeat("c", 64)
		hashD     = "0x" + strings.Repeat("d", 64)
		hashE     = "0x" + strings.Repeat("e", 64)
		hashUpper = "0x" + strings.Repeat("A", 64)
	)

	t.Run("manual payment with admin and no tx_hash is accepted", func(t *testing.T) {
		order := f.newOrder(t, "ORD-MANUAL-OK")
		mustExecPayment(t, f.db, insertPaymentRow, order, "PAYMENT", "MANUAL", f.admin, nil)
	})

	t.Run("onchain payment with tx_hash and no admin is accepted", func(t *testing.T) {
		order := f.newOrder(t, "ORD-ONCHAIN-OK")
		mustExecPayment(t, f.db, insertPaymentRow, order, "PAYMENT", "ONCHAIN", nil, hashA)
	})

	t.Run("onchain refund on same order with distinct tx_hash is accepted", func(t *testing.T) {
		order := f.newOrder(t, "ORD-ONCHAIN-REFUND-OK")
		mustExecPayment(t, f.db, insertPaymentRow, order, "PAYMENT", "ONCHAIN", nil, hashC)
		mustExecPayment(t, f.db, insertPaymentRow, order, "REFUND", "ONCHAIN", nil, hashB)
	})

	t.Run("onchain payment with uppercase hex tx_hash is accepted", func(t *testing.T) {
		order := f.newOrder(t, "ORD-ONCHAIN-UPPER-OK")
		mustExecPayment(t, f.db, insertPaymentRow, order, "PAYMENT", "ONCHAIN", nil, hashUpper)
	})

	t.Run("manual payment carrying a tx_hash is refused", func(t *testing.T) {
		order := f.newOrder(t, "ORD-MANUAL-WITH-HASH")
		expectPaymentRefused(t, f.db, paymentCheckViolation, "issuance_order_payments_mode_ck",
			insertPaymentRow, order, "PAYMENT", "MANUAL", f.admin, hashA)
	})

	t.Run("manual payment without admin is refused", func(t *testing.T) {
		order := f.newOrder(t, "ORD-MANUAL-NO-ADMIN")
		expectPaymentRefused(t, f.db, paymentCheckViolation, "issuance_order_payments_mode_ck",
			insertPaymentRow, order, "PAYMENT", "MANUAL", nil, nil)
	})

	t.Run("onchain payment without tx_hash is refused", func(t *testing.T) {
		order := f.newOrder(t, "ORD-ONCHAIN-NO-HASH")
		expectPaymentRefused(t, f.db, paymentCheckViolation, "issuance_order_payments_mode_ck",
			insertPaymentRow, order, "PAYMENT", "ONCHAIN", nil, nil)
	})

	t.Run("onchain payment carrying an admin is refused", func(t *testing.T) {
		order := f.newOrder(t, "ORD-ONCHAIN-WITH-ADMIN")
		expectPaymentRefused(t, f.db, paymentCheckViolation, "issuance_order_payments_mode_ck",
			insertPaymentRow, order, "PAYMENT", "ONCHAIN", f.admin, hashA)
	})

	t.Run("onchain payment with truncated tx_hash is refused", func(t *testing.T) {
		order := f.newOrder(t, "ORD-ONCHAIN-SHORT-HASH")
		expectPaymentRefused(t, f.db, paymentCheckViolation, "issuance_order_payments_tx_hash_ck",
			insertPaymentRow, order, "PAYMENT", "ONCHAIN", nil, "0x123")
	})

	t.Run("onchain payment with tx_hash missing 0x prefix is refused", func(t *testing.T) {
		order := f.newOrder(t, "ORD-ONCHAIN-NO-PREFIX")
		expectPaymentRefused(t, f.db, paymentCheckViolation, "issuance_order_payments_tx_hash_ck",
			insertPaymentRow, order, "PAYMENT", "ONCHAIN", nil, strings.Repeat("a", 64))
	})

	t.Run("duplicate tx_hash on another order payment is refused", func(t *testing.T) {
		first := f.newOrder(t, "ORD-DUP-HASH-PAY-1")
		second := f.newOrder(t, "ORD-DUP-HASH-PAY-2")
		mustExecPayment(t, f.db, insertPaymentRow, first, "PAYMENT", "ONCHAIN", nil, hashD)
		expectPaymentRefused(t, f.db, paymentUniqueViolation, "issuance_order_payments_tx_hash_uidx",
			insertPaymentRow, second, "PAYMENT", "ONCHAIN", nil, hashD)
	})

	t.Run("duplicate tx_hash reused as a refund on another order is refused", func(t *testing.T) {
		first := f.newOrder(t, "ORD-DUP-HASH-REF-1")
		second := f.newOrder(t, "ORD-DUP-HASH-REF-2")
		mustExecPayment(t, f.db, insertPaymentRow, first, "PAYMENT", "ONCHAIN", nil, hashE)
		expectPaymentRefused(t, f.db, paymentUniqueViolation, "issuance_order_payments_tx_hash_uidx",
			insertPaymentRow, second, "REFUND", "ONCHAIN", nil, hashE)
	})

	t.Run("unknown provider falls through mode_ck to provider_ck", func(t *testing.T) {
		order := f.newOrder(t, "ORD-UNKNOWN-PROVIDER")
		expectPaymentRefused(t, f.db, paymentCheckViolation, "issuance_order_payments_provider_ck",
			insertPaymentRow, order, "PAYMENT", "STRIPE", f.admin, nil)
	})
}
