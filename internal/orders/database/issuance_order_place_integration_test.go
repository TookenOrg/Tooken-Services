//go:build integration

// Coverage of PlaceIssuanceOrder against a real
// PostgreSQL.
//
// Unlike order_model_integration_test.go, these cases cannot run inside a
// transaction rolled back at the end: PlaceIssuanceOrder opens and commits its
// own. What it writes can never be deleted either (the history is append-only
// and referenced with RESTRICT). Run it on a THROWAWAY database only:
//
//	docker run -d --rm --name tooken-place -e POSTGRES_PASSWORD=p -p 55432:5432 postgres:16
//	docker cp tools/db/baseline.sql tooken-place:/b.sql
//	docker exec tooken-place psql -U postgres -q -v ON_ERROR_STOP=1 -f /b.sql
//	TEST_DATABASE_URL="postgres://postgres:p@localhost:55432/postgres?sslmode=disable" \
//	    go test -tags integration -run TestPlaceIssuanceOrder ./internal/orders/database/
//
// Every case creates its own investors and assets with unique names, so two
// runs on the same database give the same result.
//
// The technical failures (a statement that errors halfway) are produced with
// lock_timeout: a second transaction holds a lock on the table a given step
// reads, and the call under test gives up after a few hundred milliseconds.
// That reaches every step deterministically, and proves the rollback: the
// order written at step 7 must vanish when step 8 or 9 fails.
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
	"github.com/lib/pq"
)

const placeTTL = 15 * time.Minute

// placeRun makes every name of this run unique, across runs on the same database.
var (
	placeRun = time.Now().UnixNano() % 1_000_000_000
	placeSeq atomic.Int64
)

func placeUnique(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, placeRun, placeSeq.Add(1))
}

// openPlaceDB opens the test database and makes it the one PlaceIssuanceOrder
// uses. extra is appended to the DSN as run-time parameters (lock_timeout).
func openPlaceDB(t *testing.T, extra string) *sql.DB {
	t.Helper()

	// openOrderModelDB skips when TEST_DATABASE_URL is unset and checks that
	// the schema carries the order model.
	openOrderModelDB(t)
	dsn := os.Getenv("TEST_DATABASE_URL")
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

// useDB points globals.DB at db for the duration of the case.
func useDB(t *testing.T, db *sql.DB) {
	t.Helper()
	previous := globals.DB
	globals.DB = db
	t.Cleanup(func() { globals.DB = previous })
}

// ---------------------------------------------------------------- fixtures

type kycFixture struct {
	status  string
	expires string // SQL expression for kyc_expires_at
}

var (
	kycVerified        = kycFixture{"verified", "now() + interval '1 year'"}
	kycVerifiedNoLimit = kycFixture{"verified", "NULL"}
)

func newInvestor(t *testing.T, db *sql.DB, k kycFixture) int {
	t.Helper()
	verifiedAt := "NULL"
	if k.status == "verified" {
		verifiedAt = "now() - interval '1 day'"
	}
	var id int
	query := fmt.Sprintf(`
	    INSERT INTO usr.users (full_name, email, password, kyc_status, kyc_verified_at, kyc_expires_at)
	    VALUES ('Place probe', $1, 'not-a-real-hash', $2, %s, %s)
	    RETURNING id`, verifiedAt, k.expires)
	if err := db.QueryRow(query, placeUnique("place")+"@probe.local", k.status).Scan(&id); err != nil {
		t.Fatalf("creating investor: %v", err)
	}
	return id
}

type assetFixture struct {
	status   int     // ass.real_estate_status id
	deleted  bool    // soft-deleted (status 7, no token)
	noConfig bool    // no row in real_estate_shares_config
	total    int64   // total_shares
	price    string  // price_per_share
	rate     *string // entry_fee_rate, nil = NULL
	cap      *int64  // max_shares_per_investor, nil = NULL
	currency string
}

func fundraising() assetFixture {
	return assetFixture{status: 4, total: 1500, price: "199.995", rate: strPtr("2.5"), currency: "EUR"}
}

func strPtr(s string) *string { return &s }
func i64Ptr(i int64) *int64   { return &i }

func newAsset(t *testing.T, db *sql.DB, a assetFixture) int {
	t.Helper()

	var issuer int
	if err := db.QueryRow(`
	    INSERT INTO ass.issuer (name, legal_form) VALUES ($1, 'SA') RETURNING id`,
		placeUnique("issuer")).Scan(&issuer); err != nil {
		t.Fatalf("creating issuer: %v", err)
	}

	// A public status needs a token (real_estate_active_requires_token_ck); a
	// deleted asset must not have one (real_estate_deleted_not_tokenized_ck).
	var token sql.NullInt64
	status := a.status
	if a.deleted {
		status = 7
	} else if status >= 3 && status <= 5 {
		var id int64
		if err := db.QueryRow(`
		    INSERT INTO blk.token (symbol, token_name, address, nb_decimal)
		    VALUES ('PRB', $1::text, '0x' || substr(md5($1::text), 1, 32) || '00000000', 0)
		    RETURNING id`, placeUnique("token")).Scan(&id); err != nil {
			t.Fatalf("creating token: %v", err)
		}
		token = sql.NullInt64{Int64: id, Valid: true}
	}

	var id int
	if err := db.QueryRow(`
	    INSERT INTO ass.real_estate (title, estate_type, issuer_id, status_id, token_id, deleted_at)
	    SELECT $1, min(id), $2, $3, $4, CASE WHEN $5 THEN now() END FROM ass.real_estate_type
	    RETURNING id`, placeUnique("Villa Belair"), issuer, status, token, a.deleted).Scan(&id); err != nil {
		t.Fatalf("creating asset: %v", err)
	}

	if !a.noConfig {
		currency := a.currency
		if currency == "" {
			currency = "EUR"
		}
		if _, err := db.Exec(`
		    INSERT INTO ass.real_estate_shares_config
		        (real_estate_id, total_shares, price_per_share, entry_fee_rate, max_shares_per_investor, currency_code)
		    VALUES ($1, $2, $3, $4, $5, $6)`,
			id, a.total, a.price, a.rate, a.cap, currency); err != nil {
			t.Fatalf("creating shares config: %v", err)
		}
	}
	return id
}

// rawOrder writes an order directly in SQL, in any status, to set up the stock.
// expiresIn is a PostgreSQL interval added to now() for reservation_expires_at.
func rawOrder(t *testing.T, db execer, user, asset int, quantity int64, status int, expiresIn string) {
	t.Helper()
	if _, err := db.Exec(rawOrderSQL, user, asset, quantity, status, placeUnique("RAW"), expiresIn, nil); err != nil {
		t.Fatalf("creating raw order (status %d): %v", status, err)
	}
}

type execer interface {
	Exec(query string, args ...any) (sql.Result, error)
}

const rawOrderSQL = `
    INSERT INTO iss.issuance_orders (
        user_id, asset_id, quantity, status_id, order_reference,
        unit_price, currency_code, gross_amount, fee_amount, amount_due,
        created_at, reservation_expires_at, delivery_tx_hash, refund_reason, idempotency_key)
    VALUES (
        $1, $2, $3::bigint, $4::int, $5::text,
        1, 'EUR', $3::numeric, 0, $3::numeric,
        now() - interval '1 hour', now() + $6::interval,
        CASE WHEN $4::int IN (5, 10, 11, 12) THEN '0x' || md5($5::text) || md5($5::text || '#') END,
        CASE WHEN $4::int IN (13, 14) THEN 'PLATFORM' END,
        $7)`

func placeInput(user, asset int, quantity int) PlaceOrderDTO {
	return PlaceOrderDTO{
		UserID:         user,
		RealEstateID:   asset,
		Quantity:       quantity,
		IdempotencyKey: placeUnique("key"),
		OrderReference: placeUnique("ISS"),
		ReservationTTL: placeTTL,
	}
}

// ---------------------------------------------------------------- assertions

func countOf(t *testing.T, db *sql.DB, query string, args ...any) int {
	t.Helper()
	var n int
	if err := db.QueryRow(query, args...).Scan(&n); err != nil {
		t.Fatalf("%v\n%s", err, query)
	}
	return n
}

// expectNothingWritten checks that a refused or failed call left no trace: no
// order under its key or its reference, hence no fee and no history either.
func expectNothingWritten(t *testing.T, db *sql.DB, in PlaceOrderDTO) {
	t.Helper()
	n := countOf(t, db, `
	    SELECT count(*) FROM iss.issuance_orders
	    WHERE (user_id = $1 AND idempotency_key = $2) OR order_reference = $3`,
		in.UserID, in.IdempotencyKey, in.OrderReference)
	if n != 0 {
		t.Fatalf("want nothing written, found %d order(s)", n)
	}
}

// expectReferenceUnused checks that no order carries the reference: used when
// the key itself legitimately exists on an earlier order.
func expectReferenceUnused(t *testing.T, db *sql.DB, reference string) {
	t.Helper()
	if n := countOf(t, db, `SELECT count(*) FROM iss.issuance_orders WHERE order_reference = $1`, reference); n != 0 {
		t.Fatalf("want nothing written, found %d order(s) with reference %s", n, reference)
	}
}

func reservedOn(t *testing.T, db *sql.DB, asset int) int64 {
	t.Helper()
	onAsset, _, err := readReservedSharesNoTx(db, asset)
	if err != nil {
		t.Fatal(err)
	}
	return onAsset
}

func readReservedSharesNoTx(db *sql.DB, asset int) (int64, int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()
	return readReservedShares(context.Background(), tx, asset, 0)
}

func expectErrorIs(t *testing.T, err, want error) {
	t.Helper()
	if !errors.Is(err, want) {
		t.Fatalf("error = %v, want %v", err, want)
	}
}

// expectTechnicalError checks a failure that is none of the business refusals:
// the service must answer 500, not a 4xx that would blame the investor.
func expectTechnicalError(t *testing.T, err error, wantInMessage string) {
	t.Helper()
	if err == nil {
		t.Fatal("want an error, got none")
	}
	for _, business := range []error{
		ErrOrderAssetNotFound, ErrOrderAssetNotOpen, ErrOrderIdempotencyKeyReuse,
		ErrOrderInvestorNotEligible, ErrOrderNotEnoughShares, ErrOrderStakeLimit,
		ErrDuplicateOrderReference,
	} {
		if errors.Is(err, business) {
			t.Fatalf("want a technical error, got business refusal %v", err)
		}
	}
	if !strings.Contains(err.Error(), wantInMessage) {
		t.Fatalf("error %q does not mention %q", err, wantInMessage)
	}
}

// holdLock runs stmt in a transaction left open until the case ends (or until
// release is called), so that the call under test meets the lock it takes.
func holdLock(t *testing.T, db *sql.DB, stmt string, args ...any) (release func()) {
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
	return release
}

// waitForLockWaiters blocks until n sessions wait on a lock, so a concurrent
// case acts once the call under test is provably queued — not after a sleep.
func waitForLockWaiters(t *testing.T, db *sql.DB, n int) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if countOf(t, db, `
		    SELECT count(*) FROM pg_stat_activity
		    WHERE datname = current_database() AND wait_event_type = 'Lock'`) >= n {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %d lock waiter(s)", n)
}

type placeResult struct {
	placed PlacedOrderDTO
	err    error
}

func placeAsync(in PlaceOrderDTO) <-chan placeResult {
	out := make(chan placeResult, 1)
	go func() {
		p, err := PlaceIssuanceOrder(context.Background(), in)
		out <- placeResult{p, err}
	}()
	return out
}

// ---------------------------------------------------------------- cases

func TestPlaceIssuanceOrder(t *testing.T) {
	ctx := context.Background()
	db := openPlaceDB(t, "")
	useDB(t, db)

	t.Run("nominal order: amounts, fee line, reservation and history", func(t *testing.T) {
		a := fundraising()
		a.cap = i64Ptr(50)
		asset := newAsset(t, db, a)
		alice := newInvestor(t, db, kycVerified)
		in := placeInput(alice, asset, 7)

		placed, err := PlaceIssuanceOrder(ctx, in)
		if err != nil {
			t.Fatal(err)
		}
		if placed.Replayed || placed.OrderID == 0 || placed.OrderReference != in.OrderReference {
			t.Fatalf("placed = %+v", placed)
		}

		var row string
		if err := db.QueryRow(`
		    SELECT concat_ws(' | ', user_id, asset_id, quantity, status_id, order_reference, idempotency_key,
		                     unit_price, currency_code, gross_amount, fee_amount, amount_due,
		                     reservation_expires_at - created_at)
		    FROM iss.issuance_orders WHERE id = $1`, placed.OrderID).Scan(&row); err != nil {
			t.Fatal(err)
		}
		want := fmt.Sprintf("%d | %d | 7 | 2 | %s | %s | 199.99500000 | EUR | 1399.97 | 35.00 | 1434.97 | 00:15:00",
			alice, asset, in.OrderReference, in.IdempotencyKey)
		if row != want {
			t.Fatalf("order row\n got %s\nwant %s", row, want)
		}

		var fee string
		if err := db.QueryRow(`
		    SELECT string_agg(concat_ws(' | ', fee_code, rate, base_amount, amount), ' ; ')
		    FROM iss.issuance_order_fees WHERE order_id = $1`, placed.OrderID).Scan(&fee); err != nil {
			t.Fatal(err)
		}
		if fee != "ENTRY | 2.5000 | 1399.97 | 35.00" {
			t.Fatalf("fee lines = %q", fee)
		}

		var history string
		var instants int
		if err := db.QueryRow(`
		    SELECT string_agg(concat(coalesce(from_status_id::text, 'NULL'), '->', to_status_id, ' by ', actor_user_id),
		                      ', ' ORDER BY created_at, id),
		           count(DISTINCT created_at)
		    FROM iss.issuance_order_status_history WHERE order_id = $1`, placed.OrderID).Scan(&history, &instants); err != nil {
			t.Fatal(err)
		}
		if want := fmt.Sprintf("NULL->1 by %d, 1->2 by %d", alice, alice); history != want {
			t.Fatalf("history = %q, want %q", history, want)
		}
		if instants != 1 {
			t.Fatalf("history rows carry %d distinct instants, want 1 (same transaction)", instants)
		}
	})

	t.Run("fees", func(t *testing.T) {
		cases := []struct {
			name      string
			price     string
			rate      *string
			quantity  int
			wantFee   string
			wantLines int
		}{
			{"no entry fee (NULL rate)", "100", nil, 3, "0.00", 0},
			{"zero entry fee writes no line (U4)", "100", strPtr("0"), 3, "0.00", 0},
			{"fee rounding to zero writes no line", "0.01", strPtr("2"), 1, "0.00", 0},
			{"fee on the rounded gross (0.195 -> 0.20 -> 0.01)", "0.195", strPtr("2.5"), 1, "0.01", 1},
			{"100 % fee", "10", strPtr("100"), 1, "10.00", 1},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				a := fundraising()
				a.price, a.rate = c.price, c.rate
				asset := newAsset(t, db, a)
				placed, err := PlaceIssuanceOrder(ctx, placeInput(newInvestor(t, db, kycVerified), asset, c.quantity))
				if err != nil {
					t.Fatal(err)
				}
				var fee string
				if err := db.QueryRow(`SELECT fee_amount::text FROM iss.issuance_orders WHERE id = $1`, placed.OrderID).Scan(&fee); err != nil {
					t.Fatal(err)
				}
				if fee != c.wantFee {
					t.Fatalf("fee_amount = %s, want %s", fee, c.wantFee)
				}
				if n := countOf(t, db, `SELECT count(*) FROM iss.issuance_order_fees WHERE order_id = $1`, placed.OrderID); n != c.wantLines {
					t.Fatalf("%d fee line(s), want %d", n, c.wantLines)
				}
			})
		}
	})

	t.Run("currency and reservation length come from the asset and the input", func(t *testing.T) {
		a := fundraising()
		a.currency = "CHF"
		asset := newAsset(t, db, a)
		in := placeInput(newInvestor(t, db, kycVerified), asset, 1)
		in.ReservationTTL = 90 * time.Second

		placed, err := PlaceIssuanceOrder(ctx, in)
		if err != nil {
			t.Fatal(err)
		}
		var currency, length string
		if err := db.QueryRow(`
		    SELECT currency_code, (reservation_expires_at - created_at)::text
		    FROM iss.issuance_orders WHERE id = $1`, placed.OrderID).Scan(&currency, &length); err != nil {
			t.Fatal(err)
		}
		if currency != "CHF" || length != "00:01:30" {
			t.Fatalf("currency %s, reservation %s", currency, length)
		}
	})

	t.Run("idempotency", func(t *testing.T) {
		t.Run("replay returns the same order and writes nothing", func(t *testing.T) {
			asset := newAsset(t, db, fundraising())
			in := placeInput(newInvestor(t, db, kycVerified), asset, 3)
			first, err := PlaceIssuanceOrder(ctx, in)
			if err != nil {
				t.Fatal(err)
			}

			// The service generates a fresh reference on every call: the replay
			// must still answer with the stored one.
			again := in
			again.OrderReference = placeUnique("ISS")
			second, err := PlaceIssuanceOrder(ctx, again)
			if err != nil {
				t.Fatal(err)
			}
			if !second.Replayed || second.OrderID != first.OrderID || second.OrderReference != in.OrderReference {
				t.Fatalf("replay = %+v, first = %+v", second, first)
			}
			if n := countOf(t, db, `SELECT count(*) FROM iss.issuance_orders WHERE asset_id = $1`, asset); n != 1 {
				t.Fatalf("%d orders, want 1", n)
			}
			if n := countOf(t, db, `SELECT count(*) FROM iss.issuance_order_status_history WHERE order_id = $1`, first.OrderID); n != 2 {
				t.Fatalf("%d history rows, want 2", n)
			}
		})

		// A replay answers with the order as it is, without re-running any
		// business check.
		for _, c := range []struct {
			name   string
			change string
		}{
			{"replay after the fundraising closed", `UPDATE ass.real_estate SET status_id = 3 WHERE id = $2`},
			{"replay after the reservation expired", `UPDATE iss.issuance_orders SET created_at = now() - interval '1 hour', reservation_expires_at = now() - interval '1 minute' WHERE id = $3`},
			{"replay after the KYC was revoked", `UPDATE usr.users SET kyc_status = 'revoked', kyc_verified_at = NULL WHERE id = $1`},
			{"replay after the asset sold out", `UPDATE ass.real_estate_shares_config SET total_shares = 3 WHERE real_estate_id = $2`},
		} {
			t.Run(c.name, func(t *testing.T) {
				asset := newAsset(t, db, fundraising())
				alice := newInvestor(t, db, kycVerified)
				in := placeInput(alice, asset, 3)
				first, err := PlaceIssuanceOrder(ctx, in)
				if err != nil {
					t.Fatal(err)
				}
				args := []any{alice, asset, first.OrderID}
				if _, err := db.Exec(onlyUsedArgs(c.change, args)); err != nil {
					t.Fatalf("changing the world: %v", err)
				}
				second, err := PlaceIssuanceOrder(ctx, in)
				if err != nil {
					t.Fatal(err)
				}
				if !second.Replayed || second.OrderID != first.OrderID {
					t.Fatalf("replay = %+v", second)
				}
			})
		}

		t.Run("same key, other quantity is refused (U5, 422)", func(t *testing.T) {
			asset := newAsset(t, db, fundraising())
			in := placeInput(newInvestor(t, db, kycVerified), asset, 3)
			if _, err := PlaceIssuanceOrder(ctx, in); err != nil {
				t.Fatal(err)
			}
			other := in
			other.Quantity = 4
			other.OrderReference = placeUnique("ISS")
			_, err := PlaceIssuanceOrder(ctx, other)
			expectErrorIs(t, err, ErrOrderIdempotencyKeyReuse)
			if n := countOf(t, db, `SELECT count(*) FROM iss.issuance_orders WHERE asset_id = $1`, asset); n != 1 {
				t.Fatalf("%d orders, want 1", n)
			}
		})

		t.Run("same key, other asset is refused", func(t *testing.T) {
			first := newAsset(t, db, fundraising())
			second := newAsset(t, db, fundraising())
			in := placeInput(newInvestor(t, db, kycVerified), first, 3)
			if _, err := PlaceIssuanceOrder(ctx, in); err != nil {
				t.Fatal(err)
			}
			other := in
			other.RealEstateID = second
			other.OrderReference = placeUnique("ISS")
			_, err := PlaceIssuanceOrder(ctx, other)
			expectErrorIs(t, err, ErrOrderIdempotencyKeyReuse)
			expectReferenceUnused(t, db, other.OrderReference)
		})

		t.Run("the key is scoped to its investor", func(t *testing.T) {
			asset := newAsset(t, db, fundraising())
			alice := placeInput(newInvestor(t, db, kycVerified), asset, 3)
			bob := placeInput(newInvestor(t, db, kycVerified), asset, 3)
			bob.IdempotencyKey = alice.IdempotencyKey

			a, err := PlaceIssuanceOrder(ctx, alice)
			if err != nil {
				t.Fatal(err)
			}
			b, err := PlaceIssuanceOrder(ctx, bob)
			if err != nil {
				t.Fatal(err)
			}
			if b.Replayed || a.OrderID == b.OrderID {
				t.Fatalf("Bob got Alice's order: %+v / %+v", a, b)
			}
		})

		t.Run("concurrent requests with the same key create one order", func(t *testing.T) {
			asset := newAsset(t, db, fundraising())
			in := placeInput(newInvestor(t, db, kycVerified), asset, 2)

			const n = 10
			results := make([]<-chan placeResult, n)
			for i := range results {
				call := in
				call.OrderReference = placeUnique("ISS")
				results[i] = placeAsync(call)
			}
			created, replayed := 0, 0
			var id int64
			for _, ch := range results {
				r := <-ch
				if r.err != nil {
					t.Fatal(r.err)
				}
				if id != 0 && r.placed.OrderID != id {
					t.Fatalf("two orders for one key: %d and %d", id, r.placed.OrderID)
				}
				id = r.placed.OrderID
				if r.placed.Replayed {
					replayed++
				} else {
					created++
				}
			}
			if created != 1 || replayed != n-1 {
				t.Fatalf("created %d, replayed %d", created, replayed)
			}
		})

		// Two assets are two locks: nothing serialises these two calls, so the
		// second can pass step 2 before the first commits. The unique index
		// is then the last line, and must still read as a reused key.
		t.Run("same key on another asset, racing: the unique index answers 422", func(t *testing.T) {
			first := newAsset(t, db, fundraising())
			second := newAsset(t, db, fundraising())
			alice := newInvestor(t, db, kycVerified)
			key := placeUnique("key")

			blocker, err := db.Begin()
			if err != nil {
				t.Fatal(err)
			}
			defer blocker.Rollback()
			if _, err := blocker.Exec(rawOrderSQL, alice, first, 1, statusAwaitingPayment, placeUnique("RAW"), "10 minutes", key); err != nil {
				t.Fatal(err)
			}

			in := placeInput(alice, second, 1)
			in.IdempotencyKey = key
			result := placeAsync(in)
			waitForLockWaiters(t, db, 1) // queued on the unique index
			if err := blocker.Commit(); err != nil {
				t.Fatal(err)
			}

			r := <-result
			expectErrorIs(t, r.err, ErrOrderIdempotencyKeyReuse)
			expectReferenceUnused(t, db, in.OrderReference)
		})
	})

	t.Run("asset visibility and status (U7)", func(t *testing.T) {
		alice := newInvestor(t, db, kycVerified)
		cases := []struct {
			name  string
			asset func() int
			want  error
		}{
			{"unknown asset", func() int { return 2_000_000_000 }, ErrOrderAssetNotFound},
			{"deleted asset", func() int { a := fundraising(); a.deleted = true; return newAsset(t, db, a) }, ErrOrderAssetNotFound},
			{"draft asset", func() int { a := fundraising(); a.status = 1; return newAsset(t, db, a) }, ErrOrderAssetNotFound},
			{"asset in review", func() int { a := fundraising(); a.status = 2; return newAsset(t, db, a) }, ErrOrderAssetNotFound},
			{"closed asset", func() int { a := fundraising(); a.status = 6; return newAsset(t, db, a) }, ErrOrderAssetNotFound},
			{"published, fundraising not open yet", func() int { a := fundraising(); a.status = 3; return newAsset(t, db, a) }, ErrOrderAssetNotOpen},
			{"funded, fundraising over", func() int { a := fundraising(); a.status = 5; return newAsset(t, db, a) }, ErrOrderAssetNotOpen},
			{"fundraising without share configuration", func() int { a := fundraising(); a.noConfig = true; return newAsset(t, db, a) }, ErrOrderAssetNotOpen},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				in := placeInput(alice, c.asset(), 1)
				_, err := PlaceIssuanceOrder(ctx, in)
				expectErrorIs(t, err, c.want)
				expectNothingWritten(t, db, in)
			})
		}

		t.Run("the asset is checked before the investor", func(t *testing.T) {
			a := fundraising()
			a.status = 3
			in := placeInput(newInvestor(t, db, kycFixture{"none", "NULL"}), newAsset(t, db, a), 1)
			_, err := PlaceIssuanceOrder(ctx, in)
			expectErrorIs(t, err, ErrOrderAssetNotOpen)
		})
	})

	t.Run("KYC (U7, 403)", func(t *testing.T) {
		asset := newAsset(t, db, fundraising())
		cases := []struct {
			name string
			kyc  kycFixture
			want error
		}{
			{"no KYC", kycFixture{"none", "NULL"}, ErrOrderInvestorNotEligible},
			{"KYC pending", kycFixture{"pending", "NULL"}, ErrOrderInvestorNotEligible},
			{"KYC approved but not yet verified", kycFixture{"approved", "NULL"}, ErrOrderInvestorNotEligible},
			{"KYC rejected", kycFixture{"rejected", "NULL"}, ErrOrderInvestorNotEligible},
			{"KYC expired status", kycFixture{"expired", "NULL"}, ErrOrderInvestorNotEligible},
			{"KYC revoked", kycFixture{"revoked", "NULL"}, ErrOrderInvestorNotEligible},
			{"verified, expiry date passed", kycFixture{"verified", "now() - interval '1 hour'"}, ErrOrderInvestorNotEligible},
			{"verified, no expiry date", kycVerifiedNoLimit, nil},
			{"verified, expiry in the future", kycVerified, nil},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				in := placeInput(newInvestor(t, db, c.kyc), asset, 1)
				_, err := PlaceIssuanceOrder(ctx, in)
				if c.want == nil {
					if err != nil {
						t.Fatal(err)
					}
					return
				}
				expectErrorIs(t, err, c.want)
				expectNothingWritten(t, db, in)
			})
		}

		t.Run("unknown investor", func(t *testing.T) {
			in := placeInput(2_000_000_000, asset, 1)
			_, err := PlaceIssuanceOrder(ctx, in)
			expectErrorIs(t, err, ErrOrderInvestorNotEligible)
		})
	})

	t.Run("stock", func(t *testing.T) {
		// Each case: an asset of 1500 shares, a set of existing orders by
		// another investor, then Alice orders.
		type existing struct {
			quantity int64
			status   int
			expires  string
		}
		cases := []struct {
			name     string
			orders   []existing
			quantity int
			want     error
			message  string
		}{
			{"the exact last shares", []existing{{1495, 2, "10 minutes"}}, 5, nil, ""},
			{"one share too many", []existing{{1495, 2, "10 minutes"}}, 6, ErrOrderNotEnoughShares, "only 5 shares left, 6 requested"},
			{"sold out", []existing{{1500, 4, "0"}}, 1, ErrOrderNotEnoughShares, "only 0 shares left, 1 requested"},
			{"expired reservations free their shares", []existing{{1495, 2, "-1 minute"}}, 10, nil, ""},
			{"cancelled, expired and refunded orders hold nothing", []existing{{500, 6, "0"}, {500, 7, "0"}, {499, 14, "0"}}, 1500, nil, ""},
			{"paid, closed, delivering and reversing orders hold shares", []existing{{300, 4, "0"}, {300, 5, "0"}, {300, 8, "0"}, {300, 11, "0"}, {295, 12, "0"}}, 6, ErrOrderNotEnoughShares, "only 5 shares left, 6 requested"},
			{"a whole asset in one order", nil, 1500, nil, ""},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				asset := newAsset(t, db, fundraising())
				bob := newInvestor(t, db, kycVerified)
				for _, o := range c.orders {
					rawOrder(t, db, bob, asset, o.quantity, o.status, o.expires)
				}
				in := placeInput(newInvestor(t, db, kycVerified), asset, c.quantity)
				_, err := PlaceIssuanceOrder(ctx, in)
				if c.want == nil {
					if err != nil {
						t.Fatal(err)
					}
					return
				}
				expectErrorIs(t, err, c.want)
				if !strings.Contains(err.Error(), c.message) {
					t.Fatalf("message %q does not say %q", err, c.message)
				}
				expectNothingWritten(t, db, in)
			})
		}

		t.Run("orders on another asset do not count", func(t *testing.T) {
			busy := newAsset(t, db, fundraising())
			rawOrder(t, db, newInvestor(t, db, kycVerified), busy, 1500, 4, "0")
			free := newAsset(t, db, fundraising())
			if _, err := PlaceIssuanceOrder(ctx, placeInput(newInvestor(t, db, kycVerified), free, 1500)); err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("per-investor cap (U8)", func(t *testing.T) {
		type existing struct {
			mine     bool
			quantity int64
			status   int
			expires  string
		}
		cases := []struct {
			name     string
			cap      *int64
			orders   []existing
			quantity int
			want     error
			message  string
		}{
			{"exactly at the cap", i64Ptr(50), []existing{{true, 45, 2, "10 minutes"}}, 5, nil, ""},
			{"above the cap", i64Ptr(50), []existing{{true, 45, 2, "10 minutes"}}, 10, ErrOrderStakeLimit, "you already hold 45 of 50 shares allowed on this asset, 10 requested"},
			{"a single order above the cap", i64Ptr(50), nil, 51, ErrOrderStakeLimit, "you already hold 0 of 50"},
			{"a closed order counts, a cancelled one does not (case 19)", i64Ptr(50), []existing{{true, 45, 5, "0"}, {true, 30, 6, "0"}}, 10, ErrOrderStakeLimit, "you already hold 45 of 50"},
			{"an expired reservation does not count", i64Ptr(50), []existing{{true, 45, 2, "-1 minute"}}, 50, nil, ""},
			{"other investors' orders do not count", i64Ptr(50), []existing{{false, 900, 4, "0"}}, 50, nil, ""},
			{"no cap: a large order passes (case 16)", nil, []existing{{true, 400, 4, "0"}}, 1000, nil, ""},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				a := fundraising()
				a.cap = c.cap
				asset := newAsset(t, db, a)
				alice := newInvestor(t, db, kycVerified)
				bob := newInvestor(t, db, kycVerified)
				for _, o := range c.orders {
					owner := bob
					if o.mine {
						owner = alice
					}
					rawOrder(t, db, owner, asset, o.quantity, o.status, o.expires)
				}
				in := placeInput(alice, asset, c.quantity)
				_, err := PlaceIssuanceOrder(ctx, in)
				if c.want == nil {
					if err != nil {
						t.Fatal(err)
					}
					return
				}
				expectErrorIs(t, err, c.want)
				if !strings.Contains(err.Error(), c.message) {
					t.Fatalf("message %q does not say %q", err, c.message)
				}
				expectNothingWritten(t, db, in)
			})
		}
	})

	t.Run("concurrency", func(t *testing.T) {
		// Case 17 (P7), repeated: one test that passes once proves nothing.
		t.Run("two investors race for the last 5 shares", func(t *testing.T) {
			for i := 0; i < 20; i++ {
				asset := newAsset(t, db, fundraising())
				rawOrder(t, db, newInvestor(t, db, kycVerified), asset, 1495, 4, "0")
				alice := placeAsync(placeInput(newInvestor(t, db, kycVerified), asset, 5))
				bob := placeAsync(placeInput(newInvestor(t, db, kycVerified), asset, 5))
				ra, rb := <-alice, <-bob

				ok, refused := 0, 0
				for _, r := range []placeResult{ra, rb} {
					switch {
					case r.err == nil:
						ok++
					case errors.Is(r.err, ErrOrderNotEnoughShares):
						refused++
					default:
						t.Fatalf("round %d: unexpected error %v", i, r.err)
					}
				}
				if ok != 1 || refused != 1 {
					t.Fatalf("round %d: %d accepted, %d refused", i, ok, refused)
				}
				if got := reservedOn(t, db, asset); got != 1500 {
					t.Fatalf("round %d: %d shares reserved, want 1500", i, got)
				}
			}
		})

		// The price changes while the order waits for the lock.
		t.Run("a price changed while waiting is the price charged", func(t *testing.T) {
			a := fundraising()
			a.price, a.rate = "200", nil
			asset := newAsset(t, db, a)

			manager, err := db.Begin()
			if err != nil {
				t.Fatal(err)
			}
			defer manager.Rollback()
			if _, err := manager.Exec(`UPDATE ass.real_estate SET updated_at = now() WHERE id = $1`, asset); err != nil {
				t.Fatal(err)
			}
			if _, err := manager.Exec(`UPDATE ass.real_estate_shares_config SET price_per_share = 250 WHERE real_estate_id = $1`, asset); err != nil {
				t.Fatal(err)
			}

			result := placeAsync(placeInput(newInvestor(t, db, kycVerified), asset, 10))
			waitForLockWaiters(t, db, 1)
			if err := manager.Commit(); err != nil {
				t.Fatal(err)
			}

			r := <-result
			if r.err != nil {
				t.Fatal(r.err)
			}
			var unit, gross string
			if err := db.QueryRow(`SELECT unit_price::text, gross_amount::text FROM iss.issuance_orders WHERE id = $1`, r.placed.OrderID).Scan(&unit, &gross); err != nil {
				t.Fatal(err)
			}
			if unit != "250.00000000" || gross != "2500.00" {
				t.Fatalf("charged %s (gross %s), want the new price 250", unit, gross)
			}
		})

		// Case 21: the fundraising closes while the order waits.
		t.Run("a fundraising closed while waiting refuses the order", func(t *testing.T) {
			asset := newAsset(t, db, fundraising())
			manager, err := db.Begin()
			if err != nil {
				t.Fatal(err)
			}
			defer manager.Rollback()
			if _, err := manager.Exec(`UPDATE ass.real_estate SET status_id = 3 WHERE id = $1`, asset); err != nil {
				t.Fatal(err)
			}

			in := placeInput(newInvestor(t, db, kycVerified), asset, 1)
			result := placeAsync(in)
			waitForLockWaiters(t, db, 1)
			if err := manager.Commit(); err != nil {
				t.Fatal(err)
			}

			r := <-result
			expectErrorIs(t, r.err, ErrOrderAssetNotOpen)
			expectNothingWritten(t, db, in)
		})

		t.Run("orders on two assets do not wait for each other", func(t *testing.T) {
			busy := newAsset(t, db, fundraising())
			holdLock(t, db, `SELECT 1 FROM ass.real_estate WHERE id = $1 FOR UPDATE`, busy)

			free := newAsset(t, db, fundraising())
			in := placeInput(newInvestor(t, db, kycVerified), free, 1)
			done := make(chan error, 1)
			go func() {
				_, err := PlaceIssuanceOrder(ctx, in)
				done <- err
			}()
			select {
			case err := <-done:
				if err != nil {
					t.Fatal(err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("an order on another asset waited for the lock")
			}
		})
	})

	t.Run("failures roll everything back", func(t *testing.T) {
		t.Run("reference collision is reported for the retry", func(t *testing.T) {
			asset := newAsset(t, db, fundraising())
			first := placeInput(newInvestor(t, db, kycVerified), asset, 1)
			if _, err := PlaceIssuanceOrder(ctx, first); err != nil {
				t.Fatal(err)
			}
			second := placeInput(newInvestor(t, db, kycVerified), asset, 1)
			second.OrderReference = first.OrderReference
			_, err := PlaceIssuanceOrder(ctx, second)
			expectErrorIs(t, err, ErrDuplicateOrderReference)
			if n := countOf(t, db, `SELECT count(*) FROM iss.issuance_orders WHERE user_id = $1`, second.UserID); n != 0 {
				t.Fatalf("%d orders left behind", n)
			}
		})

		t.Run("an amount too large for the column fails without a trace", func(t *testing.T) {
			// 999 999 999 999 x 10 000 000 = 1e19 > numeric(20,2).
			a := fundraising()
			a.total, a.price, a.rate = 10_000_000, "999999999999", nil
			in := placeInput(newInvestor(t, db, kycVerified), newAsset(t, db, a), 10_000_000)
			_, err := PlaceIssuanceOrder(ctx, in)
			expectTechnicalError(t, err, "insert order")
			expectNothingWritten(t, db, in)
		})

		t.Run("a cancelled context opens nothing", func(t *testing.T) {
			cancelled, cancel := context.WithCancel(ctx)
			cancel()
			in := placeInput(newInvestor(t, db, kycVerified), newAsset(t, db, fundraising()), 1)
			if _, err := PlaceIssuanceOrder(cancelled, in); !errors.Is(err, context.Canceled) {
				t.Fatalf("error = %v, want context.Canceled", err)
			}
			expectNothingWritten(t, db, in)
		})

		// One case per step: a concurrent transaction holds a lock on what the
		// step reads or writes, and the call gives up after lock_timeout.
		impatient := openPlaceDB(t, "lock_timeout=300")
		steps := []struct {
			name    string
			lock    string // takes the lock; $1 is the asset when used
			usesID  bool
			message string
		}{
			{"step 1, the asset lock", `SELECT 1 FROM ass.real_estate WHERE id = $1 FOR UPDATE`, true, "lock real estate"},
			{"step 2, the idempotency lookup", `LOCK TABLE iss.issuance_orders IN ACCESS EXCLUSIVE MODE`, false, "find order by idempotency key"},
			{"step 3, the asset read", `LOCK TABLE ass.real_estate_shares_config IN ACCESS EXCLUSIVE MODE`, false, "read locked real estate"},
			{"step 4, the KYC read", `LOCK TABLE usr.users IN ACCESS EXCLUSIVE MODE`, false, "check KYC"},
			{"step 5, the reserved shares", `LOCK TABLE iss.issuance_order_statuses IN ACCESS EXCLUSIVE MODE`, false, "read reserved shares"},
			{"step 8, the fee line", `LOCK TABLE iss.issuance_order_fees IN ACCESS EXCLUSIVE MODE`, false, "insert fee"},
			{"step 9, the history", `LOCK TABLE iss.issuance_order_status_history IN ACCESS EXCLUSIVE MODE`, false, "insert creation history"},
		}
		for _, s := range steps {
			t.Run(s.name, func(t *testing.T) {
				asset := newAsset(t, db, fundraising())
				in := placeInput(newInvestor(t, db, kycVerified), asset, 1)

				var args []any
				if s.usesID {
					args = append(args, asset)
				}
				release := holdLock(t, db, s.lock, args...)

				globals.DB = impatient
				_, err := PlaceIssuanceOrder(ctx, in)
				globals.DB = db
				release()

				expectTechnicalError(t, err, s.message)
				var pqErr *pq.Error
				if !errors.As(err, &pqErr) || pqErr.Code != "55P03" {
					t.Fatalf("want lock_not_available (55P03), got %v", err)
				}
				// Steps 8 and 9 fail after the order was inserted at step 7:
				// the rollback must take it away too.
				expectNothingWritten(t, db, in)
			})
		}
	})
}

// onlyUsedArgs inlines the $1..$3 the statement uses, so one table of changes
// can mix statements that need the investor, the asset or the order.
func onlyUsedArgs(stmt string, args []any) string {
	for i := len(args); i >= 1; i-- {
		stmt = strings.ReplaceAll(stmt, fmt.Sprintf("$%d", i), fmt.Sprint(args[i-1]))
	}
	return stmt
}
