//go:build integration

// Coverage of ConfirmIssuanceOrderPayment against a real PostgreSQL.
//
// The function opens and commits its own transaction, and the status history
// is append-only: what these cases write can never be deleted. Run them on a
// THROWAWAY database built from the baseline:
//
//	docker exec -i <container> psql -U postgres -d <db> -v ON_ERROR_STOP=1 < tools/db/baseline.sql
//	TEST_DATABASE_URL="postgres://...?sslmode=disable" \
//	    go test -tags integration -race -run 'TestConfirmIssuanceOrderPayment' ./internal/payments/database/
//
// Every case creates its own users, assets and orders with unique names, so
// two runs on the same database give the same result.
//
// The technical failures are produced with lock_timeout: a second transaction
// holds a SHARE lock on the table a given step writes to, and the call under
// test gives up after a few hundred milliseconds. That reaches each write
// deterministically and proves the rollback of the writes made before it.
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/locks"
	_ "github.com/lib/pq"
)

var (
	confirmRun = time.Now().UnixNano() % 1_000_000_000
	confirmSeq atomic.Int64
)

func confirmUnique(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, confirmRun, confirmSeq.Add(1))
}

// openConfirmDB opens the test database. extra is appended to the DSN as
// run-time parameters (lock_timeout).
func openConfirmDB(t *testing.T, extra string) *sql.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	if extra != "" {
		sep := "?"
		if strings.Contains(dsn, "?") {
			sep = "&"
		}
		dsn += sep + extra
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// useConfirmDB points globals.DB at db for the duration of the case.
func useConfirmDB(t *testing.T, db *sql.DB) {
	t.Helper()
	previous := globals.DB
	globals.DB = db
	t.Cleanup(func() { globals.DB = previous })
}

// ---------------------------------------------------------------- fixtures

func newConfirmUser(t *testing.T, db *sql.DB, role string) int {
	t.Helper()
	var id int
	if err := db.QueryRow(`
	    INSERT INTO usr.users (full_name, email, password, role)
	    VALUES ('Payment probe', $1, 'not-a-real-hash', $2)
	    RETURNING id`, confirmUnique("pay")+"@probe.local", role).Scan(&id); err != nil {
		t.Fatalf("creating user: %v", err)
	}
	return id
}

// newConfirmAsset writes a fundraising asset (status 4) of 1 500 shares.
func newConfirmAsset(t *testing.T, db *sql.DB) int {
	t.Helper()
	var issuer int
	if err := db.QueryRow(`
	    INSERT INTO ass.issuer (name, legal_form) VALUES ($1, 'SA') RETURNING id`,
		confirmUnique("issuer")).Scan(&issuer); err != nil {
		t.Fatalf("creating issuer: %v", err)
	}
	var token int64
	if err := db.QueryRow(`
	    INSERT INTO blk.token (symbol, token_name, address, nb_decimal)
	    VALUES ('PAY', $1::text, '0x' || substr(md5($1::text), 1, 32) || '00000000', 0)
	    RETURNING id`, confirmUnique("token")).Scan(&token); err != nil {
		t.Fatalf("creating token: %v", err)
	}
	var id int
	if err := db.QueryRow(`
	    INSERT INTO ass.real_estate (title, estate_type, issuer_id, status_id, token_id)
	    SELECT $1, min(id), $2, 4, $3 FROM ass.real_estate_type
	    RETURNING id`, confirmUnique("Villa Belair"), issuer, token).Scan(&id); err != nil {
		t.Fatalf("creating asset: %v", err)
	}
	if _, err := db.Exec(`
	    INSERT INTO ass.real_estate_shares_config
	        (real_estate_id, total_shares, price_per_share, entry_fee_rate, max_shares_per_investor, currency_code)
	    VALUES ($1, 1500, 200, 2, 50, 'EUR')`, id); err != nil {
		t.Fatalf("creating shares config: %v", err)
	}
	return id
}

type confirmOrder struct {
	id  int64
	ref string
}

// newConfirmOrder writes an order of 10 shares at 200.00 EUR plus 40.00 of
// fees (amount_due 2 040.00), in any status. expiresIn is a PostgreSQL
// interval added to now() for reservation_expires_at; "" stores NULL (allowed
// outside status 2).
func newConfirmOrder(t *testing.T, db *sql.DB, user, asset, status int, expiresIn string) confirmOrder {
	t.Helper()
	ref := confirmUnique("ISS")
	var expires any
	if expiresIn != "" {
		expires = expiresIn
	}
	var id int64
	if err := db.QueryRow(`
	    INSERT INTO iss.issuance_orders (
	        user_id, asset_id, quantity, status_id, order_reference,
	        unit_price, currency_code, gross_amount, fee_amount, amount_due,
	        created_at, updated_at, reservation_expires_at, delivery_tx_hash, refund_reason)
	    VALUES (
	        $1, $2, 10, $3::int, $4::text,
	        200, 'EUR', 2000, 40, 2040,
	        now() - interval '1 hour', now() - interval '1 hour',
	        now() + $5::interval,
	        CASE WHEN $3::int IN (5, 10, 11, 12) THEN '0x' || md5($4::text) || md5($4::text || '#') END,
	        CASE WHEN $3::int IN (13, 14) THEN 'PLATFORM' END)
	    RETURNING id`, user, asset, status, ref, expires).Scan(&id); err != nil {
		t.Fatalf("creating order (status %d): %v", status, err)
	}
	return confirmOrder{id: id, ref: ref}
}

// confirmFixture is one ADMIN, one investor, one asset and one order awaiting
// payment for 15 more minutes.
type confirmFixture struct {
	admin    int
	investor int
	asset    int
	order    confirmOrder
}

func newConfirmFixture(t *testing.T, db *sql.DB) confirmFixture {
	t.Helper()
	f := confirmFixture{
		admin:    newConfirmUser(t, db, "ADMIN"),
		investor: newConfirmUser(t, db, "USER"),
		asset:    newConfirmAsset(t, db),
	}
	f.order = newConfirmOrder(t, db, f.investor, f.asset, 2, "15 minutes")
	return f
}

// ---------------------------------------------------------------- assertions

func confirmCount(t *testing.T, db *sql.DB, query string, args ...any) int {
	t.Helper()
	var n int
	if err := db.QueryRow(query, args...).Scan(&n); err != nil {
		t.Fatalf("%v\n%s", err, query)
	}
	return n
}

func statusOf(t *testing.T, db *sql.DB, order confirmOrder) int {
	t.Helper()
	return confirmCount(t, db, `SELECT status_id FROM iss.issuance_orders WHERE id = $1`, order.id)
}

func paymentsOf(t *testing.T, db *sql.DB, order confirmOrder) int {
	t.Helper()
	return confirmCount(t, db, `SELECT count(*) FROM iss.issuance_order_payments WHERE order_id = $1`, order.id)
}

func historyOf(t *testing.T, db *sql.DB, order confirmOrder) int {
	t.Helper()
	return confirmCount(t, db, `SELECT count(*) FROM iss.issuance_order_status_history WHERE order_id = $1`, order.id)
}

// expectUntouched checks that a refused or failed call left the order as it
// was: same status, no payment, no history.
func expectUntouched(t *testing.T, db *sql.DB, order confirmOrder, status int) {
	t.Helper()
	if got := statusOf(t, db, order); got != status {
		t.Fatalf("status = %d, want %d (unchanged)", got, status)
	}
	if n := paymentsOf(t, db, order); n != 0 {
		t.Fatalf("want no payment, found %d", n)
	}
	if n := historyOf(t, db, order); n != 0 {
		t.Fatalf("want no history, found %d", n)
	}
}

// expectConfirmedOnce checks the result of a successful confirmation, however
// many times it was asked: status 4, one payment, one history line.
func expectConfirmedOnce(t *testing.T, db *sql.DB, order confirmOrder) {
	t.Helper()
	if got := statusOf(t, db, order); got != 4 {
		t.Fatalf("status = %d, want 4", got)
	}
	if n := paymentsOf(t, db, order); n != 1 {
		t.Fatalf("want exactly 1 payment, found %d", n)
	}
	if n := historyOf(t, db, order); n != 1 {
		t.Fatalf("want exactly 1 history line, found %d", n)
	}
}

func expectIs(t *testing.T, err, want error) {
	t.Helper()
	if !errors.Is(err, want) {
		t.Fatalf("error = %v, want %v", err, want)
	}
}

// expectTechnical checks a failure that is none of the business refusals: the
// handler must answer 500, not a 4xx.
func expectTechnical(t *testing.T, err error, wantInMessage string) {
	t.Helper()
	if err == nil {
		t.Fatal("want an error, got none")
	}
	for _, business := range []error{
		ErrOrderNotFound, ErrOrderNotAwaitingPayment, ErrOrderReservationExpired, locks.ErrOrderAssetNotFound,
	} {
		if errors.Is(err, business) {
			t.Fatalf("want a technical error, got business refusal %v", err)
		}
	}
	if !strings.Contains(err.Error(), wantInMessage) {
		t.Fatalf("error %q does not mention %q", err, wantInMessage)
	}
}

// holdLock runs stmt in a transaction left open until release (or the end of
// the case), so that the call under test meets the lock it takes.
func holdLock(t *testing.T, db *sql.DB, stmt string, args ...any) (tx *sql.Tx, release func()) {
	t.Helper()
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(stmt, args...); err != nil {
		tx.Rollback()
		t.Fatalf("taking lock: %v\n%s", err, stmt)
	}
	var once sync.Once
	release = func() { once.Do(func() { tx.Rollback() }) }
	t.Cleanup(release)
	return tx, release
}

// waitForLockWaiters blocks until n sessions wait on a lock, so a concurrent
// case acts once the call under test is provably queued, not after a sleep.
func waitForLockWaiters(t *testing.T, db *sql.DB, n int) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if confirmCount(t, db, `
		    SELECT count(*) FROM pg_stat_activity
		    WHERE datname = current_database() AND wait_event_type = 'Lock'`) >= n {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %d lock waiter(s)", n)
}

func confirmAsync(asset int, ref string, admin int) <-chan error {
	out := make(chan error, 1)
	go func() {
		out <- ConfirmIssuanceOrderPayment(context.Background(), asset, ref, admin)
	}()
	return out
}

// ---------------------------------------------------------------- cases

func TestConfirmIssuanceOrderPayment(t *testing.T) {
	ctx := context.Background()
	db := openConfirmDB(t, "")
	useConfirmDB(t, db)

	t.Run("01 nominal: status, payment and history written together", func(t *testing.T) {
		f := newConfirmFixture(t, db)

		if err := ConfirmIssuanceOrderPayment(ctx, f.asset, f.order.ref, f.admin); err != nil {
			t.Fatalf("error = %v", err)
		}
		expectConfirmedOnce(t, db, f.order)

		var refreshed bool
		if err := db.QueryRow(`
		    SELECT updated_at > now() - interval '1 minute'
		    FROM iss.issuance_orders WHERE id = $1`, f.order.id).Scan(&refreshed); err != nil {
			t.Fatal(err)
		}
		if !refreshed {
			t.Fatal("updated_at was not refreshed")
		}

		var (
			kind, provider, currency, amount string
			recordedBy                       int
			externalRef, idemKey             sql.NullString
		)
		if err := db.QueryRow(`
		    SELECT kind, provider, amount::text, currency_code, recorded_by_user_id,
		           external_reference, idempotency_key
		    FROM iss.issuance_order_payments WHERE order_id = $1`, f.order.id).Scan(
			&kind, &provider, &amount, &currency, &recordedBy, &externalRef, &idemKey); err != nil {
			t.Fatal(err)
		}
		if kind != "PAYMENT" || provider != "MANUAL" || amount != "2040.00" || currency != "EUR" || recordedBy != f.admin {
			t.Fatalf("payment = (%s, %s, %s %s, by %d), want (PAYMENT, MANUAL, 2040.00 EUR, by %d)",
				kind, provider, amount, currency, recordedBy, f.admin)
		}
		if externalRef.Valid || idemKey.Valid {
			t.Fatalf("external_reference = %v, idempotency_key = %v, want both NULL", externalRef, idemKey)
		}

		var (
			from, to, actor int
			reason          sql.NullString
			txHash          sql.NullString
		)
		if err := db.QueryRow(`
		    SELECT from_status_id, to_status_id, actor_user_id, reason, tx_hash
		    FROM iss.issuance_order_status_history WHERE order_id = $1`, f.order.id).Scan(
			&from, &to, &actor, &reason, &txHash); err != nil {
			t.Fatal(err)
		}
		if from != 2 || to != 4 || actor != f.admin {
			t.Fatalf("history = %d -> %d by %d, want 2 -> 4 by %d", from, to, actor, f.admin)
		}
		if !reason.Valid || reason.String == "" {
			t.Fatal("history reason is empty")
		}
		if txHash.Valid {
			t.Fatalf("tx_hash = %q, want NULL: no transaction on a payment", txHash.String)
		}
	})

	t.Run("02 double click: second call succeeds and writes nothing", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		for i := range 2 {
			if err := ConfirmIssuanceOrderPayment(ctx, f.asset, f.order.ref, f.admin); err != nil {
				t.Fatalf("call %d: error = %v", i+1, err)
			}
		}
		expectConfirmedOnce(t, db, f.order)
	})

	t.Run("03 an order already paid by someone else: success, nothing written", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		paid := newConfirmOrder(t, db, f.investor, f.asset, 4, "15 minutes")
		if err := ConfirmIssuanceOrderPayment(ctx, f.asset, paid.ref, f.admin); err != nil {
			t.Fatalf("error = %v", err)
		}
		expectUntouched(t, db, paid, 4)
	})

	// The 13 h 51 bug: reservation_expires_at is never cleared, so a double
	// click after the end of the reservation must still be a success.
	t.Run("04 already paid, reservation since lapsed: success, not expired", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		paid := newConfirmOrder(t, db, f.investor, f.asset, 4, "-1 minute")
		if err := ConfirmIssuanceOrderPayment(ctx, f.asset, paid.ref, f.admin); err != nil {
			t.Fatalf("error = %v, want nil", err)
		}
		expectUntouched(t, db, paid, 4)
	})

	t.Run("05 reservation lapsed: expired, nothing written", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		lapsed := newConfirmOrder(t, db, f.investor, f.asset, 2, "-1 minute")
		err := ConfirmIssuanceOrderPayment(ctx, f.asset, lapsed.ref, f.admin)
		expectIs(t, err, ErrOrderReservationExpired)
		expectUntouched(t, db, lapsed, 2)
	})

	t.Run("06 every other status: not awaiting payment, nothing written", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		// 8 and 9 are absent: a confirmed payment starts the delivery on its
		// own, so a second confirm-payment lands on an order that already
		// moved to DELIVERY_IN_PROGRESS or DELIVERY_FAILED. Those two are
		// covered by 08 below, where they answer like 4 does.
		for _, status := range []int{1, 5, 6, 7, 10, 11, 12, 13, 14} {
			t.Run(fmt.Sprintf("status %d", status), func(t *testing.T) {
				o := newConfirmOrder(t, db, f.investor, f.asset, status, "15 minutes")
				err := ConfirmIssuanceOrderPayment(ctx, f.asset, o.ref, f.admin)
				expectIs(t, err, ErrOrderNotAwaitingPayment)
				expectUntouched(t, db, o, status)
			})
		}
	})

	// The delivery moves the order out of PAYMENT_CONFIRMED by itself, within
	// milliseconds of the answer. An ADMIN who clicks twice must not be told
	// their payment failed because the delivery was quick, so these two
	// statuses answer exactly like 4: nothing written, no error.
	t.Run("08 the delivery already started or failed: idempotent, nothing written", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		for _, status := range []int{8, 9} {
			t.Run(fmt.Sprintf("status %d", status), func(t *testing.T) {
				o := newConfirmOrder(t, db, f.investor, f.asset, status, "15 minutes")
				if err := ConfirmIssuanceOrderPayment(ctx, f.asset, o.ref, f.admin); err != nil {
					t.Fatalf("error = %v, want nil", err)
				}
				expectUntouched(t, db, o, status)
			})
		}
	})

	t.Run("07 cancelled with a lapsed reservation: not awaiting payment, not expired", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		o := newConfirmOrder(t, db, f.investor, f.asset, 6, "-1 minute")
		err := ConfirmIssuanceOrderPayment(ctx, f.asset, o.ref, f.admin)
		expectIs(t, err, ErrOrderNotAwaitingPayment)
		expectUntouched(t, db, o, 6)
	})

	t.Run("08 no reservation date outside status 2: not awaiting payment", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		o := newConfirmOrder(t, db, f.investor, f.asset, 6, "")
		err := ConfirmIssuanceOrderPayment(ctx, f.asset, o.ref, f.admin)
		expectIs(t, err, ErrOrderNotAwaitingPayment)
		expectUntouched(t, db, o, 6)
	})

	t.Run("09 unknown reference: order not found", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		err := ConfirmIssuanceOrderPayment(ctx, f.asset, confirmUnique("ISS-UNKNOWN"), f.admin)
		expectIs(t, err, ErrOrderNotFound)
		expectIs(t, err, locks.ErrOrderNotFound)
	})

	t.Run("10 unknown asset: asset not found, the order is untouched", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		err := ConfirmIssuanceOrderPayment(ctx, -1, f.order.ref, f.admin)
		expectIs(t, err, locks.ErrOrderAssetNotFound)
		expectUntouched(t, db, f.order, 2)
	})

	t.Run("11 the shares stay reserved after the payment", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		reserved := func() int {
			return confirmCount(t, db, `
			    SELECT COALESCE(SUM(o.quantity), 0)
			    FROM iss.issuance_orders o
			    JOIN iss.issuance_order_statuses s ON s.id = o.status_id
			    WHERE o.asset_id = $1 AND s.counts_as_reserved
			      AND (o.status_id <> 2 OR o.reservation_expires_at > now())`, f.asset)
		}
		before := reserved()
		if err := ConfirmIssuanceOrderPayment(ctx, f.asset, f.order.ref, f.admin); err != nil {
			t.Fatal(err)
		}
		if after := reserved(); after != before || after != 10 {
			t.Fatalf("reserved = %d before, %d after, want 10 both times", before, after)
		}
	})

	t.Run("12 unknown admin: technical error, everything rolled back", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		err := ConfirmIssuanceOrderPayment(ctx, f.asset, f.order.ref, -1)
		expectTechnical(t, err, "insert payment")
		expectUntouched(t, db, f.order, 2)
	})

	t.Run("13 ten concurrent confirmations: ten successes, one payment", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		var (
			wg    sync.WaitGroup
			start = make(chan struct{})
			errs  = make([]error, 10)
		)
		for i := range errs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				errs[i] = ConfirmIssuanceOrderPayment(ctx, f.asset, f.order.ref, f.admin)
			}()
		}
		close(start)
		wg.Wait()
		for i, err := range errs {
			if err != nil {
				t.Fatalf("call %d: error = %v", i, err)
			}
		}
		expectConfirmedOnce(t, db, f.order)
	})

	// The asset first, the order second, like PlaceIssuanceOrder. While the
	// call waits on the asset, the order must still be free.
	t.Run("14 the asset is locked before the order", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		_, release := holdLock(t, db, `SELECT 1 FROM ass.real_estate WHERE id = $1 FOR UPDATE`, f.asset)

		result := confirmAsync(f.asset, f.order.ref, f.admin)
		waitForLockWaiters(t, db, 1)

		probe, err := db.Begin()
		if err != nil {
			t.Fatal(err)
		}
		_, err = probe.Exec(`SELECT 1 FROM iss.issuance_orders WHERE id = $1 FOR UPDATE NOWAIT`, f.order.id)
		probe.Rollback()
		if err != nil {
			t.Fatalf("the order is locked while the asset is still awaited: %v", err)
		}

		release()
		if err := <-result; err != nil {
			t.Fatalf("error = %v", err)
		}
		expectConfirmedOnce(t, db, f.order)
	})

	// What changed while the call waited is what it judges: the order is read
	// after the lock, not before.
	t.Run("15 an order cancelled while the call waits: not awaiting payment", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		canceller, _ := holdLock(t, db, `SELECT 1 FROM ass.real_estate WHERE id = $1 FOR UPDATE`, f.asset)
		if _, err := canceller.Exec(`UPDATE iss.issuance_orders SET status_id = 6 WHERE id = $1`, f.order.id); err != nil {
			t.Fatal(err)
		}

		result := confirmAsync(f.asset, f.order.ref, f.admin)
		waitForLockWaiters(t, db, 1)
		if err := canceller.Commit(); err != nil {
			t.Fatal(err)
		}

		expectIs(t, <-result, ErrOrderNotAwaitingPayment)
		expectUntouched(t, db, f.order, 6)
	})

	// now() is frozen at the start of the transaction, before the lock wait.
	// The reservation must be judged on the real time.
	t.Run("16 a reservation that lapses while the call waits: expired", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		o := newConfirmOrder(t, db, f.investor, f.asset, 2, "2 seconds")
		_, release := holdLock(t, db, `SELECT 1 FROM ass.real_estate WHERE id = $1 FOR UPDATE`, f.asset)

		result := confirmAsync(f.asset, o.ref, f.admin)
		waitForLockWaiters(t, db, 1)

		deadline := time.Now().Add(10 * time.Second)
		for confirmCount(t, db, `
		    SELECT count(*) FROM iss.issuance_orders
		    WHERE id = $1 AND reservation_expires_at < clock_timestamp()`, o.id) == 0 {
			if time.Now().After(deadline) {
				t.Fatal("the reservation never lapsed")
			}
			time.Sleep(50 * time.Millisecond)
		}
		release()

		expectIs(t, <-result, ErrOrderReservationExpired)
		expectUntouched(t, db, o, 2)
	})

	t.Run("technical failures roll everything back", func(t *testing.T) {
		impatient := openConfirmDB(t, "lock_timeout=300")

		for _, c := range []struct {
			name, table, message string
		}{
			{"17 the payment cannot be written", "iss.issuance_order_payments", "insert payment"},
			{"18 the history cannot be written, the payment is undone", "iss.issuance_order_status_history", "insert order status history"},
			{"19 the status cannot be written, payment and history are undone", "iss.issuance_orders", "update issuance order status"},
		} {
			t.Run(c.name, func(t *testing.T) {
				f := newConfirmFixture(t, db)
				useConfirmDB(t, impatient)
				_, release := holdLock(t, db, "LOCK TABLE "+c.table+" IN SHARE MODE")

				err := ConfirmIssuanceOrderPayment(ctx, f.asset, f.order.ref, f.admin)
				release()
				expectTechnical(t, err, c.message)
				expectUntouched(t, db, f.order, 2)
			})
		}

		t.Run("20 the order lock times out", func(t *testing.T) {
			f := newConfirmFixture(t, db)
			useConfirmDB(t, impatient)
			_, release := holdLock(t, db, `SELECT 1 FROM iss.issuance_orders WHERE id = $1 FOR UPDATE`, f.order.id)

			err := ConfirmIssuanceOrderPayment(ctx, f.asset, f.order.ref, f.admin)
			release()
			expectTechnical(t, err, "lock issuance order")
			expectUntouched(t, db, f.order, 2)
		})

		t.Run("21 the asset lock times out", func(t *testing.T) {
			f := newConfirmFixture(t, db)
			useConfirmDB(t, impatient)
			_, release := holdLock(t, db, `SELECT 1 FROM ass.real_estate WHERE id = $1 FOR UPDATE`, f.asset)

			err := ConfirmIssuanceOrderPayment(ctx, f.asset, f.order.ref, f.admin)
			release()
			expectTechnical(t, err, "lock real estate")
			expectUntouched(t, db, f.order, 2)
		})
	})

	t.Run("22 a cancelled context: nothing written", func(t *testing.T) {
		f := newConfirmFixture(t, db)
		cancelled, cancel := context.WithCancel(ctx)
		cancel()
		if err := ConfirmIssuanceOrderPayment(cancelled, f.asset, f.order.ref, f.admin); err == nil {
			t.Fatal("want an error on a cancelled context")
		}
		expectUntouched(t, db, f.order, 2)
	})
}

// TestUpdateIssuanceOrderStatusToPaymentConfirmed covers the guard the
// confirmation never reaches (the order is locked first): an UPDATE that
// matches no row.
func TestUpdateIssuanceOrderStatusToPaymentConfirmed(t *testing.T) {
	db := openConfirmDB(t, "")
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	err = UpdateIssuanceOrderStatusToPaymentConfirmed(context.Background(), tx, confirmUnique("ISS-UNKNOWN"))
	expectIs(t, err, ErrOrderNotFound)
}
