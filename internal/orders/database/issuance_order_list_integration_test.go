//go:build integration

// Coverage of ListIssuanceOrders and GetIssuanceOrderByRef against a real
// PostgreSQL: who sees which orders, in which order, and the EXPIRED status
// shown at read time (TICKET-M3-4 §4.1).
//
// Same constraints as issuance_order_place_integration_test.go: what is written
// here can never be deleted, so run it on a THROWAWAY database only:
//
//	TEST_DATABASE_URL="postgres://postgres:p@localhost:55432/postgres?sslmode=disable" \
//	    go test -tags integration -run TestListIssuanceOrders ./internal/orders/database/
//
// The database is shared with the other cases, so "every order" is checked as
// "contains these orders", never as an exact count.
package database

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/shopspring/decimal"
)

// datedOrder writes an order with a chosen created_at, to control the sort.
// createdAt and expiresAt are SQL timestamptz expressions.
func datedOrder(t *testing.T, user, asset, status int, createdAt, expiresAt string) string {
	t.Helper()
	ref := placeUnique("DATED")
	if _, err := globals.DB.Exec(`
	    INSERT INTO iss.issuance_orders (
	        user_id, asset_id, quantity, status_id, order_reference,
	        unit_price, currency_code, gross_amount, fee_amount, amount_due,
	        created_at, reservation_expires_at)
	    VALUES ($1, $2, 1, $3, $4, 1, 'EUR', 1, 0, 1, `+createdAt+`, `+expiresAt+`)`,
		user, asset, status, ref); err != nil {
		t.Fatalf("creating dated order: %v", err)
	}
	return ref
}

// storedStatus reads status_id as written in the table, not as shown.
func storedStatus(t *testing.T, ref string) int {
	t.Helper()
	var status int
	if err := globals.DB.QueryRow(
		`SELECT status_id FROM iss.issuance_orders WHERE order_reference = $1`, ref).Scan(&status); err != nil {
		t.Fatalf("reading stored status of %s: %v", ref, err)
	}
	return status
}

// refOfUser returns the reference of the only order of user, written by rawOrder.
func refOfUser(t *testing.T, user int) string {
	t.Helper()
	var ref string
	if err := globals.DB.QueryRow(
		`SELECT order_reference FROM iss.issuance_orders WHERE user_id = $1`, user).Scan(&ref); err != nil {
		t.Fatalf("reading order of user %d: %v", user, err)
	}
	return ref
}

func refsOf(orders []server.IssuanceOrder) map[string]server.IssuanceOrder {
	byRef := make(map[string]server.IssuanceOrder, len(orders))
	for _, o := range orders {
		byRef[o.OrderRef] = o
	}
	return byRef
}

func listOf(t *testing.T, viewAll bool, user int) []server.IssuanceOrder {
	t.Helper()
	orders, err := ListIssuanceOrders(context.Background(), viewAll, user)
	if err != nil {
		t.Fatalf("listing orders (viewAll=%v, user=%d): %v", viewAll, user, err)
	}
	return orders
}

func expectStatus(t *testing.T, o server.IssuanceOrder, wantID int, wantCode string) {
	t.Helper()
	if o.StatusId == nil || *o.StatusId != wantID || o.StatusCode != wantCode {
		t.Fatalf("order %s: want status %d %s, got %v %s", o.OrderRef, wantID, wantCode, o.StatusId, o.StatusCode)
	}
	if o.StatusLabel == nil || *o.StatusLabel == "" {
		t.Fatalf("order %s: status label is empty", o.OrderRef)
	}
}

func TestListIssuanceOrders(t *testing.T) {
	db := openPlaceDB(t, "")
	useDB(t, db)
	ctx := context.Background()

	place := func(t *testing.T, user, asset int) string {
		t.Helper()
		in := placeInput(user, asset, 2)
		if _, err := PlaceIssuanceOrder(ctx, in); err != nil {
			t.Fatalf("placing order: %v", err)
		}
		return in.OrderReference
	}

	t.Run("who sees what", func(t *testing.T) {
		asset := newAsset(t, db, fundraising())
		paul := newInvestor(t, db, kycVerified)
		marie := newInvestor(t, db, kycVerified)
		paulRefs := []string{place(t, paul, asset), place(t, paul, asset)}
		marieRef := place(t, marie, asset)

		t.Run("01 a user sees only his orders", func(t *testing.T) {
			got := listOf(t, false, paul)
			if len(got) != 2 {
				t.Fatalf("want 2 orders, got %d", len(got))
			}
			for _, o := range got {
				if o.UserId != paul {
					t.Fatalf("order %s belongs to user %d, not %d", o.OrderRef, o.UserId, paul)
				}
			}
			byRef := refsOf(got)
			for _, ref := range paulRefs {
				if _, ok := byRef[ref]; !ok {
					t.Fatalf("order %s missing from Paul's list", ref)
				}
			}
		})

		t.Run("02 every order, for staff", func(t *testing.T) {
			byRef := refsOf(listOf(t, true, 0))
			for _, ref := range append(paulRefs, marieRef) {
				if _, ok := byRef[ref]; !ok {
					t.Fatalf("order %s missing from the full list", ref)
				}
			}
		})

		t.Run("02b staff list ignores the user id", func(t *testing.T) {
			byRef := refsOf(listOf(t, true, paul))
			if _, ok := byRef[marieRef]; !ok {
				t.Fatalf("Marie's order %s missing when viewAll is true", marieRef)
			}
		})
	})

	t.Run("03 newest first", func(t *testing.T) {
		asset := newAsset(t, db, fundraising())
		user := newInvestor(t, db, kycVerified)
		at10 := datedOrder(t, user, asset, 6, "now() - interval '2 hours'", "now() - interval '1 hour'")
		at11 := datedOrder(t, user, asset, 6, "now() - interval '1 hour'", "now()")

		got := listOf(t, false, user)
		if len(got) != 2 || got[0].OrderRef != at11 || got[1].OrderRef != at10 {
			t.Fatalf("want [%s %s], got %v", at11, at10, refsInOrder(got))
		}
	})

	t.Run("03b same created_at, highest id first", func(t *testing.T) {
		asset := newAsset(t, db, fundraising())
		user := newInvestor(t, db, kycVerified)
		first := datedOrder(t, user, asset, 6, "'2026-01-01 10:00:00+00'", "'2026-01-01 10:15:00+00'")
		second := datedOrder(t, user, asset, 6, "'2026-01-01 10:00:00+00'", "'2026-01-01 10:15:00+00'")

		got := listOf(t, false, user)
		if len(got) != 2 || got[0].OrderRef != second || got[1].OrderRef != first {
			t.Fatalf("want [%s %s], got %v", second, first, refsInOrder(got))
		}
	})

	t.Run("04 no order is an empty slice, not nil", func(t *testing.T) {
		nobody := newInvestor(t, db, kycVerified)
		got := listOf(t, false, nobody)
		if got == nil || len(got) != 0 {
			t.Fatalf("want an empty non-nil slice, got %#v", got)
		}
	})

	t.Run("expiry shown at read time", func(t *testing.T) {
		asset := newAsset(t, db, fundraising())

		t.Run("05 awaiting payment past its reservation is EXPIRED", func(t *testing.T) {
			user := newInvestor(t, db, kycVerified)
			rawOrder(t, db, user, asset, 1, statusAwaitingPayment, "-1 minute")
			ref := refOfUser(t, user)

			got := listOf(t, false, user)
			if len(got) != 1 {
				t.Fatalf("want 1 order, got %d", len(got))
			}
			expectStatus(t, got[0], 7, "EXPIRED")

			detail, err := GetIssuanceOrderByRef(ctx, ref, false, user)
			if err != nil {
				t.Fatal(err)
			}
			expectStatus(t, detail, 7, "EXPIRED")

			if stored := storedStatus(t, ref); stored != statusAwaitingPayment {
				t.Fatalf("a read must not write: stored status is %d, want %d", stored, statusAwaitingPayment)
			}
		})

		t.Run("06 awaiting payment within its reservation stays AWAITING_PAYMENT", func(t *testing.T) {
			user := newInvestor(t, db, kycVerified)
			rawOrder(t, db, user, asset, 1, statusAwaitingPayment, "10 minutes")

			got := listOf(t, false, user)
			if len(got) != 1 {
				t.Fatalf("want 1 order, got %d", len(got))
			}
			expectStatus(t, got[0], statusAwaitingPayment, "AWAITING_PAYMENT")
		})

		t.Run("07 paid in time stays PAYMENT_CONFIRMED after the reservation", func(t *testing.T) {
			user := newInvestor(t, db, kycVerified)
			rawOrder(t, db, user, asset, 1, 4, "-1 minute")

			got := listOf(t, false, user)
			if len(got) != 1 {
				t.Fatalf("want 1 order, got %d", len(got))
			}
			expectStatus(t, got[0], 4, "PAYMENT_CONFIRMED")
		})

		t.Run("07b a cancelled order past its reservation stays CANCELLED", func(t *testing.T) {
			user := newInvestor(t, db, kycVerified)
			rawOrder(t, db, user, asset, 1, 6, "-1 minute")

			got := listOf(t, false, user)
			if len(got) != 1 {
				t.Fatalf("want 1 order, got %d", len(got))
			}
			expectStatus(t, got[0], 6, "CANCELLED")
		})
	})

	t.Run("08 fees are listed and add up to feeAmount", func(t *testing.T) {
		asset := newAsset(t, db, fundraising())
		user := newInvestor(t, db, kycVerified)
		place(t, user, asset)

		got := listOf(t, false, user)
		if len(got) != 1 {
			t.Fatalf("want 1 order, got %d", len(got))
		}
		o := got[0]
		if len(o.Fees) != 1 || o.Fees[0].Code != "ENTRY" {
			t.Fatalf("want one ENTRY fee, got %+v", o.Fees)
		}
		sum := decimal.Zero
		for _, f := range o.Fees {
			sum = sum.Add(decimal.RequireFromString(f.Amount))
		}
		if !sum.Equal(decimal.RequireFromString(o.FeeAmount)) || sum.IsZero() {
			t.Fatalf("fees add up to %s, feeAmount is %s", sum, o.FeeAmount)
		}
	})

	t.Run("08b an order without fee lines has fees: [] not null", func(t *testing.T) {
		asset := newAsset(t, db, fundraising())
		user := newInvestor(t, db, kycVerified)
		rawOrder(t, db, user, asset, 1, statusAwaitingPayment, "10 minutes")

		got := listOf(t, false, user)
		if len(got) != 1 || got[0].Fees == nil || len(got[0].Fees) != 0 {
			t.Fatalf("want fees to be an empty non-nil slice, got %#v", got)
		}
	})

	t.Run("09 the title of the asset is shown", func(t *testing.T) {
		asset := newAsset(t, db, fundraising())
		user := newInvestor(t, db, kycVerified)
		place(t, user, asset)

		var want string
		if err := db.QueryRow(`SELECT title FROM ass.real_estate WHERE id = $1`, asset).Scan(&want); err != nil {
			t.Fatal(err)
		}
		got := listOf(t, false, user)
		if len(got) != 1 || got[0].RealEstateTitle == nil || *got[0].RealEstateTitle != want {
			t.Fatalf("want title %q, got %+v", want, got)
		}
		if !strings.HasPrefix(want, "Villa Belair") {
			t.Fatalf("unexpected fixture title %q", want)
		}
	})

	// Zero rows is a read result, not a business error: the database reports
	// sql.ErrNoRows and the service decides it is ErrOrderNotFound (L8).
	t.Run("10 an unknown reference is sql.ErrNoRows", func(t *testing.T) {
		for _, viewAll := range []bool{false, true} {
			_, err := GetIssuanceOrderByRef(ctx, placeUnique("UNKNOWN"), viewAll, 0)
			if !errors.Is(err, sql.ErrNoRows) {
				t.Fatalf("viewAll=%v: want sql.ErrNoRows, got %v", viewAll, err)
			}
		}
	})

	t.Run("detail by reference", func(t *testing.T) {
		asset := newAsset(t, db, fundraising())
		marie := newInvestor(t, db, kycVerified)
		paul := newInvestor(t, db, kycVerified)
		rawOrder(t, db, marie, asset, 1, statusAwaitingPayment, "15 minutes")
		marieRef := refOfUser(t, marie)
		// Paul's order is the most recent one: a query ignoring the reference
		// would return it instead of Marie's.
		rawOrder(t, db, paul, asset, 2, statusAwaitingPayment, "15 minutes")

		t.Run("10b the owner gets their order", func(t *testing.T) {
			got, err := GetIssuanceOrderByRef(ctx, marieRef, false, marie)
			if err != nil {
				t.Fatal(err)
			}
			if got.OrderRef != marieRef {
				t.Fatalf("want %s, got %s", marieRef, got.OrderRef)
			}
		})

		t.Run("10c someone else's order is sql.ErrNoRows", func(t *testing.T) {
			_, err := GetIssuanceOrderByRef(ctx, marieRef, false, paul)
			if !errors.Is(err, sql.ErrNoRows) {
				t.Fatalf("want sql.ErrNoRows, got %v", err)
			}
		})

		t.Run("10d staff gets the requested order, whoever owns it", func(t *testing.T) {
			got, err := GetIssuanceOrderByRef(ctx, marieRef, true, paul)
			if err != nil {
				t.Fatal(err)
			}
			if got.OrderRef != marieRef {
				t.Fatalf("want %s, got %s", marieRef, got.OrderRef)
			}
		})
	})
}

func refsInOrder(orders []server.IssuanceOrder) []string {
	refs := make([]string, len(orders))
	for i, o := range orders {
		refs[i] = o.OrderRef
	}
	return refs
}
