//go:build integration

// Direct coverage of the shared locks against a real PostgreSQL. Every case
// runs inside a transaction rolled back at the end: nothing is left behind.
// The waits and the lock ordering are covered by the callers' tests
// (internal/orders/database, internal/payments/database).
//
//	TEST_DATABASE_URL="postgres://...?sslmode=disable" \
//	    go test -tags integration -race ./internal/locks/
package locks

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

var (
	locksRun = time.Now().UnixNano() % 1_000_000_000
	locksSeq atomic.Int64
)

// locksUnique stays short: order_reference is a varchar(32).
func locksUnique(t *testing.T, prefix string) string {
	t.Helper()
	return fmt.Sprintf("%s-%d-%d", prefix, locksRun, locksSeq.Add(1))
}

// withTx runs fn inside a transaction that is always rolled back.
func withTx(t *testing.T, fn func(tx *sql.Tx)) {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	fn(tx)
}

func mustScan(t *testing.T, tx *sql.Tx, dest any, query string, args ...any) {
	t.Helper()
	if err := tx.QueryRow(query, args...).Scan(dest); err != nil {
		t.Fatalf("%v\n%s", err, query)
	}
}

// newLocksAsset writes a fundraising asset; withConfig adds its shares config
// (1 500 shares at 200 EUR, 2 % of fee, 50 per investor).
func newLocksAsset(t *testing.T, tx *sql.Tx, withConfig bool) int {
	t.Helper()
	var issuer, token, id int
	mustScan(t, tx, &issuer, `INSERT INTO ass.issuer (name, legal_form) VALUES ($1, 'SA') RETURNING id`,
		locksUnique(t, "issuer"))
	mustScan(t, tx, &token, `
	    INSERT INTO blk.token (symbol, token_name, address, nb_decimal)
	    VALUES ('LCK', $1::text, '0x' || substr(md5($1::text), 1, 32) || '00000000', 0)
	    RETURNING id`, locksUnique(t, "token"))
	mustScan(t, tx, &id, `
	    INSERT INTO ass.real_estate (title, estate_type, issuer_id, status_id, token_id)
	    SELECT $1, min(id), $2, 4, $3 FROM ass.real_estate_type
	    RETURNING id`, locksUnique(t, "Villa Belair"), issuer, token)
	if withConfig {
		if _, err := tx.Exec(`
		    INSERT INTO ass.real_estate_shares_config
		        (real_estate_id, total_shares, price_per_share, entry_fee_rate, max_shares_per_investor, currency_code)
		    VALUES ($1, 1500, 200, 2, 50, 'EUR')`, id); err != nil {
			t.Fatal(err)
		}
	}
	return id
}

// newLocksOrder writes an order of 10 shares (amount_due 2 040.00 EUR).
// expiresIn is added to now(); "" stores NULL (allowed outside status 2).
func newLocksOrder(t *testing.T, tx *sql.Tx, asset, status int, expiresIn string) (int64, string, int) {
	t.Helper()
	var user int
	mustScan(t, tx, &user, `
	    INSERT INTO usr.users (full_name, email, password) VALUES ('Lock probe', $1, 'x') RETURNING id`,
		locksUnique(t, "user")+"@probe.local")
	ref := locksUnique(t, "ISS")
	var expires any
	if expiresIn != "" {
		expires = expiresIn
	}
	var id int64
	mustScan(t, tx, &id, `
	    INSERT INTO iss.issuance_orders (
	        user_id, asset_id, quantity, status_id, order_reference,
	        unit_price, currency_code, gross_amount, fee_amount, amount_due,
	        created_at, reservation_expires_at)
	    VALUES ($1, $2, 10, $3, $4, 200, 'EUR', 2000, 40, 2040,
	            now() - interval '1 hour', now() + $5::interval)
	    RETURNING id`, user, asset, status, ref, expires)
	return id, ref, user
}

func TestLockIssuanceOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("existing order: locked", func(t *testing.T) {
		withTx(t, func(tx *sql.Tx) {
			_, ref, _ := newLocksOrder(t, tx, newLocksAsset(t, tx, true), 2, "15 minutes")
			if err := LockIssuanceOrder(ctx, tx, ref); err != nil {
				t.Fatalf("error = %v", err)
			}
		})
	})

	t.Run("unknown order: ErrOrderNotFound", func(t *testing.T) {
		withTx(t, func(tx *sql.Tx) {
			if err := LockIssuanceOrder(ctx, tx, locksUnique(t, "ISS-UNKNOWN")); !errors.Is(err, ErrOrderNotFound) {
				t.Fatalf("error = %v, want %v", err, ErrOrderNotFound)
			}
		})
	})

	t.Run("technical error: wrapped, not a 404", func(t *testing.T) {
		withTx(t, func(tx *sql.Tx) {
			cancelled, cancel := context.WithCancel(ctx)
			cancel()
			err := LockIssuanceOrder(cancelled, tx, "ISS-ANY")
			if err == nil || errors.Is(err, ErrOrderNotFound) {
				t.Fatalf("error = %v, want a technical error", err)
			}
		})
	})
}

func TestReadLockedIssuanceOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("every column is read", func(t *testing.T) {
		withTx(t, func(tx *sql.Tx) {
			asset := newLocksAsset(t, tx, true)
			id, ref, user := newLocksOrder(t, tx, asset, 2, "15 minutes")

			o, err := ReadLockedIssuanceOrder(ctx, tx, ref)
			if err != nil {
				t.Fatal(err)
			}
			if o.Id != id || o.UserId != user || o.AssetId != asset || o.Quantity != 10 ||
				o.StatusId != 2 || o.OrderReference != ref || o.CurrencyCode != "EUR" {
				t.Fatalf("order = %+v", o)
			}
			if o.UnitPrice.String() != "200" || o.GrossAmount.String() != "2000" ||
				o.FeeAmount.String() != "40" || o.AmountDue.String() != "2040" {
				t.Fatalf("amounts = %s x, %s + %s = %s", o.UnitPrice, o.GrossAmount, o.FeeAmount, o.AmountDue)
			}
			if !o.ReservationExpiresAt.Valid || o.ReservationExpired {
				t.Fatalf("reservation = %v, expired %v, want a running reservation", o.ReservationExpiresAt, o.ReservationExpired)
			}
		})
	})

	for _, c := range []struct {
		name      string
		status    int
		expiresIn string
		want      bool
	}{
		{"reservation running: not expired", 2, "15 minutes", false},
		{"reservation lapsed: expired", 2, "-1 second", true},
		{"no reservation date: not expired (COALESCE)", 6, "", false},
	} {
		t.Run(c.name, func(t *testing.T) {
			withTx(t, func(tx *sql.Tx) {
				_, ref, _ := newLocksOrder(t, tx, newLocksAsset(t, tx, true), c.status, c.expiresIn)
				o, err := ReadLockedIssuanceOrder(ctx, tx, ref)
				if err != nil {
					t.Fatal(err)
				}
				if o.ReservationExpired != c.want {
					t.Fatalf("ReservationExpired = %v, want %v", o.ReservationExpired, c.want)
				}
				if o.ReservationExpiresAt.Valid != (c.expiresIn != "") {
					t.Fatalf("ReservationExpiresAt.Valid = %v", o.ReservationExpiresAt.Valid)
				}
			})
		})
	}

	// clock_timestamp(), not now(): inside a transaction opened before the
	// reservation lapsed, the reading must still see the real time.
	t.Run("expiry judged on the real time, not the start of the transaction", func(t *testing.T) {
		withTx(t, func(tx *sql.Tx) {
			_, ref, _ := newLocksOrder(t, tx, newLocksAsset(t, tx, true), 2, "1 second")
			time.Sleep(1500 * time.Millisecond)
			o, err := ReadLockedIssuanceOrder(ctx, tx, ref)
			if err != nil {
				t.Fatal(err)
			}
			if !o.ReservationExpired {
				t.Fatal("ReservationExpired = false: the expiry was judged on the transaction start")
			}
		})
	})

	t.Run("unknown order: sql.ErrNoRows, wrapped", func(t *testing.T) {
		withTx(t, func(tx *sql.Tx) {
			_, err := ReadLockedIssuanceOrder(ctx, tx, locksUnique(t, "ISS-UNKNOWN"))
			if !errors.Is(err, sql.ErrNoRows) || !strings.Contains(err.Error(), "read locked issuance order") {
				t.Fatalf("error = %v, want a wrapped sql.ErrNoRows", err)
			}
		})
	})
}

func TestLockRealEstate(t *testing.T) {
	ctx := context.Background()

	t.Run("existing asset: locked", func(t *testing.T) {
		withTx(t, func(tx *sql.Tx) {
			if err := LockRealEstate(ctx, tx, newLocksAsset(t, tx, true)); err != nil {
				t.Fatalf("error = %v", err)
			}
		})
	})

	t.Run("unknown asset: ErrOrderAssetNotFound", func(t *testing.T) {
		withTx(t, func(tx *sql.Tx) {
			if err := LockRealEstate(ctx, tx, -1); !errors.Is(err, ErrOrderAssetNotFound) {
				t.Fatalf("error = %v, want %v", err, ErrOrderAssetNotFound)
			}
		})
	})

	t.Run("technical error: wrapped, not a 404", func(t *testing.T) {
		withTx(t, func(tx *sql.Tx) {
			cancelled, cancel := context.WithCancel(ctx)
			cancel()
			err := LockRealEstate(cancelled, tx, 1)
			if err == nil || errors.Is(err, ErrOrderAssetNotFound) {
				t.Fatalf("error = %v, want a technical error", err)
			}
		})
	})
}

func TestReadLockedRealEstate(t *testing.T) {
	ctx := context.Background()

	t.Run("asset with its shares config", func(t *testing.T) {
		withTx(t, func(tx *sql.Tx) {
			a, err := ReadLockedRealEstate(ctx, tx, newLocksAsset(t, tx, true))
			if err != nil {
				t.Fatal(err)
			}
			if a.StatusID != 4 || !a.Active || a.Deleted || !a.Tokenized || !a.HasSharesConfig ||
				a.TotalShares != 1500 || a.PricePerShare.String() != "200" || a.CurrencyCode != "EUR" ||
				!a.EntryFeeRate.Valid || a.EntryFeeRate.Decimal.String() != "2" ||
				!a.MaxSharesPerInvestor.Valid || a.MaxSharesPerInvestor.Int64 != 50 {
				t.Fatalf("asset = %+v", a)
			}
		})
	})

	t.Run("asset without shares config: zero values", func(t *testing.T) {
		withTx(t, func(tx *sql.Tx) {
			a, err := ReadLockedRealEstate(ctx, tx, newLocksAsset(t, tx, false))
			if err != nil {
				t.Fatal(err)
			}
			if a.HasSharesConfig || a.TotalShares != 0 || !a.PricePerShare.IsZero() || a.CurrencyCode != "" ||
				a.EntryFeeRate.Valid || a.MaxSharesPerInvestor.Valid {
				t.Fatalf("asset = %+v, want no config", a)
			}
		})
	})

	t.Run("unknown asset: sql.ErrNoRows, wrapped", func(t *testing.T) {
		withTx(t, func(tx *sql.Tx) {
			_, err := ReadLockedRealEstate(ctx, tx, -1)
			if !errors.Is(err, sql.ErrNoRows) || !strings.Contains(err.Error(), "read locked real estate") {
				t.Fatalf("error = %v, want a wrapped sql.ErrNoRows", err)
			}
		})
	})
}
