//go:build integration

// Coverage of the order model laid down by migrations 000027, 000028 and 000029,
// against a real PostgreSQL.
//
// Everything here is written in plain SQL, on purpose: what is under test is the
// schema, not the Go that will one day write to it. A constraint that only holds
// because today's service happens to respect it is not a constraint, and the
// services of M3-3 onwards will be written against what these tests pin down.
//
// Every case runs inside its own transaction, rolled back at the end. That is not
// a convenience, it is the only clean way out: the order history is append-only
// (I1) and references its order with RESTRICT, so once a case has written one
// history row, neither the row nor its order can ever be deleted again. A
// rollback leaves nothing behind, needs no cleanup order, and makes two
// consecutive runs on the same database give the same result.
//
// Each case also creates its own investors, issuer and asset (PROGRESS.md §45.3:
// never rely on ids that happen to exist).
//
// Run it against a database that carries the migrations:
//
//	TEST_DATABASE_URL="postgres://…" go test -tags integration ./internal/orders/database/...
package database

import (
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/lib/pq"
)

// Two transaction hashes that pass I2. Their exact value does not matter, only
// that they are well formed and distinct.
var (
	hashA = "0x" + strings.Repeat("a", 64)
	hashB = "0x" + strings.Repeat("b", 64)
)

const (
	codeCheckViolation  = "23514"
	codeUniqueViolation = "23505"
	codeRaiseException  = "P0001"
)

// Status ids of the target referential (TICKET-M3-1 §2.1). The constraints are
// CHECKs on ids, so the tests name ids too.
const (
	statusCreated         = 1
	statusAwaitingPayment = 2
	statusClosed          = 5
	statusRefundPending   = 13
)

func openOrderModelDB(t *testing.T) *sql.DB {
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

	// A database below 000029 would fail every case below on a missing column,
	// naming the column rather than the cause. Say it once, here.
	var hasCap bool
	if err := db.QueryRow(`
	    SELECT EXISTS (
	        SELECT 1 FROM information_schema.columns
	        WHERE table_schema = 'ass'
	          AND table_name   = 'real_estate_shares_config'
	          AND column_name  = 'max_shares_per_investor')`).Scan(&hasCap); err != nil {
		t.Fatal(err)
	}
	if !hasCap {
		t.Fatal("the test database does not carry migration 000029: apply tools/db/migrations first")
	}

	return db
}

// orderFixture is what every case needs: two investors and one priced asset of
// 1000 shares, created inside the case's own transaction.
type orderFixture struct {
	tx    *sql.Tx
	alice int64
	bob   int64
	asset int64
}

// withFixture runs fn inside a transaction that is always rolled back.
func withFixture(t *testing.T, db *sql.DB, fn func(t *testing.T, f orderFixture)) {
	t.Helper()

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	f := orderFixture{tx: tx}
	for _, u := range []struct {
		dst   *int64
		email string
	}{
		{&f.alice, "alice@order-model-probe.local"},
		{&f.bob, "bob@order-model-probe.local"},
	} {
		if err := tx.QueryRow(`
		    INSERT INTO usr.users (full_name, email, password)
		    VALUES ('Order model probe', $1, 'not-a-real-hash')
		    RETURNING id`, u.email).Scan(u.dst); err != nil {
			t.Fatal(err)
		}
	}

	var issuer int64
	if err := tx.QueryRow(`
	    INSERT INTO ass.issuer (name, legal_form) VALUES ('Order model probe', 'SA')
	    RETURNING id`).Scan(&issuer); err != nil {
		t.Fatal(err)
	}

	// estate_type is a foreign key to a referential the baseline seeds. Taking an
	// existing row rather than inventing one keeps the fixture honest.
	if err := tx.QueryRow(`
	    INSERT INTO ass.real_estate (title, estate_type, issuer_id)
	    SELECT 'Order model probe', min(id), $1 FROM ass.real_estate_type
	    RETURNING id`, issuer).Scan(&f.asset); err != nil {
		t.Fatalf("creating the asset (is ass.real_estate_type seeded?): %v", err)
	}
	mustExecTx(t, tx, `
	    INSERT INTO ass.real_estate_shares_config (real_estate_id, total_shares, price_per_share)
	    VALUES ($1, 1000, 200)`, f.asset)

	fn(t, f)
}

// mustExecTx runs a statement that must be accepted. It goes through a
// savepoint like expectRefused, so a failure is reported without poisoning the
// rest of the case.
func mustExecTx(t *testing.T, tx *sql.Tx, query string, args ...any) {
	t.Helper()
	if err := execInSavepoint(tx, query, args...); err != nil {
		t.Fatalf("want accepted, got %v\n%s", err, query)
	}
}

func mustScanIDTx(t *testing.T, tx *sql.Tx, query string, args ...any) int64 {
	t.Helper()
	var id int64
	if err := tx.QueryRow(query, args...).Scan(&id); err != nil {
		t.Fatalf("%v\n%s", err, query)
	}
	return id
}

// expectRefused runs a statement that the schema must reject, and checks that it
// is rejected for the right reason: the SQLSTATE, and when given, the name of the
// constraint. A refusal for an unrelated reason — a typo, a missing NOT NULL —
// would otherwise pass for the invariant under test.
//
// The statement runs under a savepoint: a failed statement aborts the whole
// transaction in PostgreSQL, and the case usually goes on afterwards.
func expectRefused(t *testing.T, tx *sql.Tx, wantCode, wantConstraint, query string, args ...any) {
	t.Helper()

	err := execInSavepoint(tx, query, args...)
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

func execInSavepoint(tx *sql.Tx, query string, args ...any) error {
	if _, err := tx.Exec(`SAVEPOINT probe`); err != nil {
		return err
	}
	if _, err := tx.Exec(query, args...); err != nil {
		if _, rbErr := tx.Exec(`ROLLBACK TO SAVEPOINT probe`); rbErr != nil {
			return errors.Join(err, rbErr)
		}
		return err
	}
	_, err := tx.Exec(`RELEASE SAVEPOINT probe`)
	return err
}

// The insert most cases start from. Columns not listed stay NULL, which is a
// draft order as InsertIssuranceOrder writes it today.
const insertBareOrder = `
    INSERT INTO iss.issuance_orders (user_id, asset_id, quantity, status_id, order_reference)
    VALUES ($1, $2, $3, $4, $5)`

// A fully priced order: 10 shares of Villa Belair #12 at 200 EUR, 2 % of fees.
const insertPricedOrder = `
    INSERT INTO iss.issuance_orders
        (user_id, asset_id, quantity, status_id, order_reference,
         unit_price, currency_code, gross_amount, fee_amount, amount_due, reservation_expires_at)
    VALUES ($1, $2, 10, $3, $4, 200.00000000, 'EUR', 2000.00, 40.00, 2040.00, now() + interval '15 minutes')
    RETURNING id`

func TestOrderModel(t *testing.T) {
	db := openOrderModelDB(t)

	// 1 — the referential is the one §2.1 describes.
	t.Run("01 the referential has 13 statuses and no PAYMENT_PENDING", func(t *testing.T) {
		want := map[int]struct {
			code     string
			final    bool
			reserved bool
		}{
			1:  {"CREATED", false, false},
			2:  {"AWAITING_PAYMENT", false, true},
			4:  {"PAYMENT_CONFIRMED", false, true},
			5:  {"CLOSED", true, true},
			6:  {"CANCELLED", true, false},
			7:  {"EXPIRED", true, false},
			8:  {"DELIVERY_IN_PROGRESS", false, true},
			9:  {"DELIVERY_FAILED", false, true},
			10: {"DELIVERED", false, true},
			11: {"REVERSAL_IN_PROGRESS", false, true},
			12: {"REVERSAL_FAILED", false, true},
			13: {"REFUND_PENDING", false, false},
			14: {"REFUNDED", true, false},
		}

		rows, err := db.Query(`SELECT id, code, is_final, counts_as_reserved FROM iss.issuance_order_statuses`)
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()

		seen := 0
		for rows.Next() {
			var (
				id              int
				code            string
				final, reserved bool
			)
			if err := rows.Scan(&id, &code, &final, &reserved); err != nil {
				t.Fatal(err)
			}
			seen++
			w, ok := want[id]
			if !ok {
				t.Errorf("unexpected status %d %s", id, code)
				continue
			}
			if code != w.code || final != w.final || reserved != w.reserved {
				t.Errorf("status %d: got (%s final=%t reserved=%t), want (%s final=%t reserved=%t)",
					id, code, final, reserved, w.code, w.final, w.reserved)
			}
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		if seen != len(want) {
			t.Errorf("got %d statuses, want %d", seen, len(want))
		}

		// The sequence must clear the ids written by hand, or the next status
		// created without an explicit id collides with an existing one.
		var next int64
		if err := db.QueryRow(`
		    SELECT CASE WHEN is_called THEN last_value + 1 ELSE last_value END
		    FROM iss.issuance_order_statuses_id_seq`).Scan(&next); err != nil {
			t.Fatal(err)
		}
		if next != 15 {
			t.Errorf("next status id = %d, want 15", next)
		}
	})

	// 2 — what the production code writes today must keep working.
	t.Run("02 a CREATED order with no amount is accepted", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			mustExecTx(t, f.tx, insertBareOrder, f.alice, f.asset, 5, statusCreated, "OM-02")
		})
	})

	// 3 — the nominal reservation.
	t.Run("03 a complete AWAITING_PAYMENT order is accepted", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			mustScanIDTx(t, f.tx, insertPricedOrder, f.alice, f.asset, statusAwaitingPayment, "OM-03")
		})
	})

	// 4 and 5 — I4: a reservation has an end, and the end is after the start.
	t.Run("04 AWAITING_PAYMENT without reservation_expires_at is refused", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			expectRefused(t, f.tx, codeCheckViolation, "issuance_orders_awaiting_payment_ck", `
			    INSERT INTO iss.issuance_orders
			        (user_id, asset_id, quantity, status_id, order_reference,
			         unit_price, currency_code, gross_amount, fee_amount, amount_due)
			    VALUES ($1, $2, 10, 2, 'OM-04', 200, 'EUR', 2000.00, 40.00, 2040.00)`,
				f.alice, f.asset)
		})
	})

	t.Run("05 AWAITING_PAYMENT expiring at its creation is refused", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			// now() is frozen for the whole transaction, so created_at and
			// reservation_expires_at are exactly equal here — the boundary itself.
			expectRefused(t, f.tx, codeCheckViolation, "issuance_orders_awaiting_payment_ck", `
			    INSERT INTO iss.issuance_orders
			        (user_id, asset_id, quantity, status_id, order_reference,
			         unit_price, currency_code, gross_amount, fee_amount, amount_due, created_at, reservation_expires_at)
			    VALUES ($1, $2, 10, 2, 'OM-05', 200, 'EUR', 2000.00, 40.00, 2040.00, now(), now())`,
				f.alice, f.asset)
		})
	})

	// 6 to 8 — I3: all five amounts or none, and they add up.
	t.Run("06 a half-priced order is refused", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			expectRefused(t, f.tx, codeCheckViolation, "issuance_orders_amounts_ck", `
			    INSERT INTO iss.issuance_orders
			        (user_id, asset_id, quantity, status_id, order_reference, unit_price, currency_code)
			    VALUES ($1, $2, 10, 1, 'OM-06', 200, 'EUR')`,
				f.alice, f.asset)
		})
	})

	t.Run("07 the gross amount is rounded half away from zero", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			// 7 × 199.995 = 1399.965. Banker's rounding would give 1399.96;
			// PostgreSQL's round() on numeric gives 1399.97, and the schema
			// arbitrates so that every writer agrees with it.
			const q = `
			    INSERT INTO iss.issuance_orders
			        (user_id, asset_id, quantity, status_id, order_reference,
			         unit_price, currency_code, gross_amount, fee_amount, amount_due)
			    VALUES ($1, $2, 7, 1, $3, 199.99500000, 'EUR', $4, 0.00, $4)`
			expectRefused(t, f.tx, codeCheckViolation, "issuance_orders_amounts_ck", q, f.alice, f.asset, "OM-07a", "1399.96")
			mustExecTx(t, f.tx, q, f.alice, f.asset, "OM-07b", "1399.97")
		})
	})

	t.Run("08 amount_due that is not gross plus fee is refused", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			expectRefused(t, f.tx, codeCheckViolation, "issuance_orders_amounts_ck", `
			    INSERT INTO iss.issuance_orders
			        (user_id, asset_id, quantity, status_id, order_reference,
			         unit_price, currency_code, gross_amount, fee_amount, amount_due)
			    VALUES ($1, $2, 10, 1, 'OM-08', 200, 'EUR', 2000.00, 40.00, 2039.00)`,
				f.alice, f.asset)
		})
	})

	// 9 and 10 — I8 and the one-entry-fee index.
	t.Run("09 a fee line is rounded by the same rule", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			order := mustScanIDTx(t, f.tx, insertPricedOrder, f.alice, f.asset, statusAwaitingPayment, "OM-09")
			// 2.5 % of 1399.97 = 34.99925, which rounds to 35.00.
			const q = `
			    INSERT INTO iss.issuance_order_fees (order_id, fee_code, rate, base_amount, amount)
			    VALUES ($1, 'ENTRY', 2.5000, 1399.97, $2)`
			expectRefused(t, f.tx, codeCheckViolation, "issuance_order_fees_amount_ck", q, order, "34.99")
			mustExecTx(t, f.tx, q, order, "35.00")
		})
	})

	t.Run("10 an order pays its entry fee once", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			order := mustScanIDTx(t, f.tx, insertPricedOrder, f.alice, f.asset, statusAwaitingPayment, "OM-10")
			const q = `
			    INSERT INTO iss.issuance_order_fees (order_id, fee_code, rate, base_amount, amount)
			    VALUES ($1, 'ENTRY', 2.0000, 2000.00, 40.00)`
			mustExecTx(t, f.tx, q, order)
			expectRefused(t, f.tx, codeUniqueViolation, "issuance_order_fees_order_code_uidx", q, order)
		})
	})

	// 11 and 12 — payments.
	t.Run("11 one payment and one refund per order", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			order := mustScanIDTx(t, f.tx, insertPricedOrder, f.alice, f.asset, statusAwaitingPayment, "OM-11")
			const q = `
			    INSERT INTO iss.issuance_order_payments
			        (order_id, kind, provider, amount, currency_code, recorded_by_user_id)
			    VALUES ($1, $2, 'MANUAL', 2040.00, 'EUR', $3)`
			mustExecTx(t, f.tx, q, order, "PAYMENT", f.bob)
			expectRefused(t, f.tx, codeUniqueViolation, "issuance_order_payments_one_payment_uidx", q, order, "PAYMENT", f.bob)
			// The partial indexes are per kind: a refund is a different row.
			mustExecTx(t, f.tx, q, order, "REFUND", f.bob)
		})
	})

	t.Run("12 malformed payments are refused", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			order := mustScanIDTx(t, f.tx, insertPricedOrder, f.alice, f.asset, statusAwaitingPayment, "OM-12")
			const q = `
			    INSERT INTO iss.issuance_order_payments
			        (order_id, kind, provider, amount, currency_code, recorded_by_user_id)
			    VALUES ($1, $2, $3, $4, 'EUR', $5)`
			for _, c := range []struct {
				name, kind, provider, amount, constraint string
			}{
				{"zero amount", "PAYMENT", "MANUAL", "0", "issuance_order_payments_amount_ck"},
				{"unknown provider", "PAYMENT", "STRIPE", "10", "issuance_order_payments_provider_ck"},
				{"unknown kind", "GIFT", "MANUAL", "10", "issuance_order_payments_kind_ck"},
			} {
				t.Run(c.name, func(t *testing.T) {
					expectRefused(t, f.tx, codeCheckViolation, c.constraint, q, order, c.kind, c.provider, c.amount, f.bob)
				})
			}
		})
	})

	// 13 to 15 — I5, the uniqueness of the delivery hash, and I2.
	t.Run("13 a closed order carries its delivery hash", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			expectRefused(t, f.tx, codeCheckViolation, "issuance_orders_delivered_hash_ck",
				insertBareOrder, f.alice, f.asset, 1, statusClosed, "OM-13a")
			mustExecTx(t, f.tx, `
			    INSERT INTO iss.issuance_orders
			        (user_id, asset_id, quantity, status_id, order_reference, delivery_tx_hash)
			    VALUES ($1, $2, 1, 5, 'OM-13b', $3)`, f.alice, f.asset, hashA)
		})
	})

	t.Run("14 one mint does not deliver two orders", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			const q = `
			    INSERT INTO iss.issuance_orders
			        (user_id, asset_id, quantity, status_id, order_reference, delivery_tx_hash)
			    VALUES ($1, $2, 1, 5, $3, $4)`
			mustExecTx(t, f.tx, q, f.alice, f.asset, "OM-14a", hashA)
			expectRefused(t, f.tx, codeUniqueViolation, "issuance_orders_delivery_tx_hash_uidx", q, f.bob, f.asset, "OM-14b", hashA)
		})
	})

	t.Run("15 a malformed hash is refused", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			const q = `
			    INSERT INTO iss.issuance_orders
			        (user_id, asset_id, quantity, status_id, order_reference, delivery_tx_hash)
			    VALUES ($1, $2, 1, 1, $3, $4)`
			expectRefused(t, f.tx, codeCheckViolation, "issuance_orders_delivery_tx_hash_ck", q, f.alice, f.asset, "OM-15a", "0x123")
			expectRefused(t, f.tx, codeCheckViolation, "issuance_orders_delivery_tx_hash_ck", q, f.alice, f.asset, "OM-15b", strings.Repeat("a", 64))
		})
	})

	// 16 — I6 and I7: a refund says why, with a known reason.
	t.Run("16 a refund carries a known reason", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			const q = `
			    INSERT INTO iss.issuance_orders
			        (user_id, asset_id, quantity, status_id, order_reference, refund_reason)
			    VALUES ($1, $2, 1, $3, $4, $5)`
			expectRefused(t, f.tx, codeCheckViolation, "issuance_orders_refund_has_reason_ck", q, f.alice, f.asset, statusRefundPending, "OM-16a", nil)
			expectRefused(t, f.tx, codeCheckViolation, "issuance_orders_refund_reason_ck", q, f.alice, f.asset, statusRefundPending, "OM-16b", "OOPS")
			mustExecTx(t, f.tx, q, f.alice, f.asset, statusRefundPending, "OM-16c", "PLATFORM")
		})
	})

	// 17 — the idempotency key is unique per author, not globally.
	t.Run("17 an idempotency key is unique per investor", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			const q = `
			    INSERT INTO iss.issuance_orders
			        (user_id, asset_id, quantity, status_id, order_reference, idempotency_key)
			    VALUES ($1, $2, 1, 1, $3, 'retry-1')`
			mustExecTx(t, f.tx, q, f.alice, f.asset, "OM-17a")
			expectRefused(t, f.tx, codeUniqueViolation, "issuance_orders_idempotency_uidx", q, f.alice, f.asset, "OM-17b")
			mustExecTx(t, f.tx, q, f.bob, f.asset, "OM-17c")
		})
	})

	// 18 and 19 — the history.
	t.Run("18 the history is append-only", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			order := mustScanIDTx(t, f.tx, insertPricedOrder, f.alice, f.asset, statusAwaitingPayment, "OM-18")
			row := mustScanIDTx(t, f.tx, `
			    INSERT INTO iss.issuance_order_status_history (order_id, from_status_id, to_status_id, reason)
			    VALUES ($1, 1, 2, 'probe') RETURNING id`, order)

			// A RAISE EXCEPTION carries no constraint name: the SQLSTATE and the
			// message are what identify the trigger.
			for _, q := range []string{
				`UPDATE iss.issuance_order_status_history SET reason = 'rewritten' WHERE id = $1`,
				`DELETE FROM iss.issuance_order_status_history WHERE id = $1`,
			} {
				err := execInSavepoint(f.tx, q, row)
				var pqErr *pq.Error
				if !errors.As(err, &pqErr) || string(pqErr.Code) != codeRaiseException ||
					!strings.Contains(pqErr.Message, "append-only") {
					t.Errorf("want the append-only refusal, got %v\n%s", err, q)
				}
			}
		})
	})

	t.Run("19 a transition goes somewhere", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			order := mustScanIDTx(t, f.tx, insertPricedOrder, f.alice, f.asset, statusAwaitingPayment, "OM-19")
			expectRefused(t, f.tx, codeCheckViolation, "issuance_order_status_history_transition_ck", `
			    INSERT INTO iss.issuance_order_status_history (order_id, from_status_id, to_status_id)
			    VALUES ($1, 2, 2)`, order)
		})
	})

	// 20 — I10 at its four boundaries, on an asset of 1000 shares.
	t.Run("20 the per-investor cap stays within the offer", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			const q = `UPDATE ass.real_estate_shares_config SET max_shares_per_investor = $1 WHERE real_estate_id = $2`
			expectRefused(t, f.tx, codeCheckViolation, "real_estate_shares_config_max_shares_ck", q, 0, f.asset)
			expectRefused(t, f.tx, codeCheckViolation, "real_estate_shares_config_max_shares_ck", q, 1001, f.asset)
			mustExecTx(t, f.tx, q, 1000, f.asset)
			mustExecTx(t, f.tx, q, nil, f.asset)
		})
	})

	// 21 — T4: CREATED no longer holds shares out of the stock.
	t.Run("21 only a reservation counts as sold", func(t *testing.T) {
		withFixture(t, db, func(t *testing.T, f orderFixture) {
			mustExecTx(t, f.tx, insertBareOrder, f.alice, f.asset, 5, statusCreated, "OM-21a")
			mustScanIDTx(t, f.tx, insertPricedOrder, f.bob, f.asset, statusAwaitingPayment, "OM-21b")

			// The same predicate as the tokens_sold aggregate of
			// assets_managements/database/real_estate_columns.go: what reserves
			// a share is carried by the referential, not by the query.
			var sold int64
			if err := f.tx.QueryRow(`
			    SELECT COALESCE(SUM(o.quantity), 0)
			    FROM iss.issuance_orders o
			    JOIN iss.issuance_order_statuses s ON s.id = o.status_id
			    WHERE s.counts_as_reserved AND o.asset_id = $1`, f.asset).Scan(&sold); err != nil {
				t.Fatal(err)
			}
			if sold != 10 {
				t.Errorf("tokens sold = %d, want 10: the CREATED order must not reserve", sold)
			}
		})
	})
}
