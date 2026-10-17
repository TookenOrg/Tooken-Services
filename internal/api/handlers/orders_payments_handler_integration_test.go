//go:build integration

// End-to-end coverage of POST /assets/real-estate/issuance/orders/{orderRef}/confirm-payment
// (M3-5) through the real router, middleware, handler, service and database.
// What the transaction proves on its own (locks, clock, rollbacks) lives in
// internal/payments/database; this file proves the contract a client sees.
//
// Writes append-only rows: run on a THROWAWAY database (see
// orders_handler_integration_test.go for the setup).
//
//	go test -tags integration -race -run 'TestConfirmPaymentEndpoint' ./internal/api/handlers/
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/internal/globals"
)

func confirmPath(ref string) string {
	return ordersPath + "/" + ref + "/confirm-payment"
}

func (h *orderHarness) confirm(token, ref string) (int, string) {
	req := httptest.NewRequest(http.MethodPost, globals.BaseURL+confirmPath(ref), nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	h.r.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

// rawOrderRef writes an order directly in SQL, like rawOrder, and returns its
// reference. expiresIn is added to now() for reservation_expires_at.
func (h *orderHarness) rawOrderRef(t *testing.T, user, asset int, status int, expiresIn string) string {
	t.Helper()
	ref := h.unique("ISS")
	if _, err := h.db.Exec(`
	    INSERT INTO iss.issuance_orders (
	        user_id, asset_id, quantity, status_id, order_reference,
	        unit_price, currency_code, gross_amount, fee_amount, amount_due,
	        created_at, reservation_expires_at, delivery_tx_hash, refund_reason)
	    VALUES (
	        $1, $2, 10, $3::int, $4::text,
	        200, 'EUR', 2000, 40, 2040,
	        now() - interval '1 hour', now() + $5::interval,
	        CASE WHEN $3::int IN (5, 10, 11, 12) THEN '0x' || md5($4::text) || md5($4::text || '#') END,
	        CASE WHEN $3::int IN (13, 14) THEN 'PLATFORM' END)`,
		user, asset, status, ref, expiresIn); err != nil {
		t.Fatalf("creating raw order (status %d): %v", status, err)
	}
	return ref
}

func (h *orderHarness) admin(t *testing.T) orderUser {
	return h.newUser(t, authUtils.RoleAdmin, "none", "NULL")
}

func (h *orderHarness) statusOf(t *testing.T, ref string) int {
	t.Helper()
	return h.count(t, `SELECT status_id FROM iss.issuance_orders WHERE order_reference = $1`, ref)
}

func (h *orderHarness) paymentsOf(t *testing.T, ref string) int {
	t.Helper()
	return h.count(t, `
	    SELECT count(*) FROM iss.issuance_order_payments p
	    JOIN iss.issuance_orders o ON o.id = p.order_id
	    WHERE o.order_reference = $1`, ref)
}

func (h *orderHarness) historyOf(t *testing.T, ref string) int {
	t.Helper()
	return h.count(t, `
	    SELECT count(*) FROM iss.issuance_order_status_history s
	    JOIN iss.issuance_orders o ON o.id = s.order_id
	    WHERE o.order_reference = $1 AND s.from_status_id = 2 AND s.to_status_id = 4`, ref)
}

// expectUntouched checks that a refused call left the order as it was.
func (h *orderHarness) expectUntouched(t *testing.T, ref string, status int) {
	t.Helper()
	if got := h.statusOf(t, ref); got != status {
		t.Fatalf("status = %d, want %d (unchanged)", got, status)
	}
	if n := h.paymentsOf(t, ref); n != 0 {
		t.Fatalf("want no payment, found %d", n)
	}
	if n := h.historyOf(t, ref); n != 0 {
		t.Fatalf("want no history 2 -> 4, found %d", n)
	}
}

func (h *orderHarness) expectConfirmedOnce(t *testing.T, ref string) {
	t.Helper()
	if got := h.statusOf(t, ref); got != 4 {
		t.Fatalf("status = %d, want 4", got)
	}
	if n := h.paymentsOf(t, ref); n != 1 {
		t.Fatalf("want exactly 1 payment, found %d", n)
	}
	if n := h.historyOf(t, ref); n != 1 {
		t.Fatalf("want exactly 1 history 2 -> 4, found %d", n)
	}
}

// decodeConfirmed decodes the 200 of the confirmation: the {message, data}
// envelope (P7), with the order in PAYMENT_CONFIRMED.
func decodeConfirmed(t *testing.T, body, ref string) server.IssuanceOrder {
	t.Helper()
	var resp server.CreateIssuanceOrderResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("decoding %s: %v", body, err)
	}
	if resp.Message == "" {
		t.Fatalf("no message in %s", body)
	}
	if resp.Data == nil {
		t.Fatalf("no order in %s", body)
	}
	o := *resp.Data
	if o.OrderRef != ref || o.StatusCode != string(server.PAYMENTCONFIRMED) {
		t.Fatalf("order = %s %s, want %s PAYMENT_CONFIRMED", o.OrderRef, o.StatusCode, ref)
	}
	return o
}

func expectMessage(t *testing.T, body string) {
	t.Helper()
	var resp server.APIResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil || resp.Message == "" {
		t.Fatalf("want an APIResponse with a message, got %s", body)
	}
}

func TestConfirmPaymentEndpoint(t *testing.T) {
	h := newOrderHarness(t)
	admin := h.admin(t)

	t.Run("who may confirm", func(t *testing.T) {
		asset := h.newAsset(t, villaBelair())
		investor := h.investor(t)
		ref := h.rawOrderRef(t, investor.id, asset, 2, "15 minutes")

		t.Run("01 anonymous is a 401", func(t *testing.T) {
			code, body := h.confirm("", ref)
			expectCode(t, code, body, http.StatusUnauthorized)
			h.expectUntouched(t, ref, 2)
		})

		t.Run("02 an invalid token is a 401", func(t *testing.T) {
			code, body := h.confirm("not-a-jwt", ref)
			expectCode(t, code, body, http.StatusUnauthorized)
			h.expectUntouched(t, ref, 2)
		})

		t.Run("03 a USER cannot confirm his own order: 403", func(t *testing.T) {
			code, body := h.confirm(investor.token, ref)
			expectCode(t, code, body, http.StatusForbidden)
			expectMessage(t, body)
			h.expectUntouched(t, ref, 2)
		})

		t.Run("04 a MANAGER is a 403 (P3: ADMIN only)", func(t *testing.T) {
			manager := h.newUser(t, authUtils.RoleManager, "none", "NULL")
			code, body := h.confirm(manager.token, ref)
			expectCode(t, code, body, http.StatusForbidden)
			h.expectUntouched(t, ref, 2)
		})

		t.Run("05 the role is checked before the order: a USER on an unknown order is a 403", func(t *testing.T) {
			code, body := h.confirm(investor.token, h.unique("ISS-UNKNOWN"))
			expectCode(t, code, body, http.StatusForbidden)
		})
	})

	t.Run("06 nominal: 200, enveloped order, payment recorded by the ADMIN", func(t *testing.T) {
		asset := h.newAsset(t, villaBelair())
		investor := h.investor(t)
		ref := h.rawOrderRef(t, investor.id, asset, 2, "15 minutes")
		reservedBefore := h.reservedOn(t, asset)

		code, body := h.confirm(admin.token, ref)
		expectCode(t, code, body, http.StatusOK)
		o := decodeConfirmed(t, body, ref)
		if o.UserId != investor.id || o.AmountDue != "2040.00" || o.StatusId == nil || *o.StatusId != 4 {
			t.Fatalf("order of user %d, %s due, status %v; want user %d, 2040.00, 4", o.UserId, o.AmountDue, o.StatusId, investor.id)
		}
		h.expectConfirmedOnce(t, ref)

		if n := h.count(t, `
		    SELECT count(*) FROM iss.issuance_order_payments p
		    JOIN iss.issuance_orders o ON o.id = p.order_id
		    WHERE o.order_reference = $1 AND p.recorded_by_user_id = $2 AND p.amount = o.amount_due`,
			ref, admin.id); n != 1 {
			t.Fatal("the payment is not recorded by the ADMIN for the amount due")
		}
		if n := h.count(t, `
		    SELECT count(*) FROM iss.issuance_order_status_history s
		    JOIN iss.issuance_orders o ON o.id = s.order_id
		    WHERE o.order_reference = $1 AND s.actor_user_id = $2`, ref, admin.id); n != 1 {
			t.Fatal("the history line does not name the ADMIN")
		}
		if after := h.reservedOn(t, asset); after != reservedBefore {
			t.Fatalf("reserved = %d after, want %d: the shares stay reserved", after, reservedBefore)
		}
	})

	t.Run("07 double click: 200 again, nothing more written", func(t *testing.T) {
		ref := h.rawOrderRef(t, h.investor(t).id, h.newAsset(t, villaBelair()), 2, "15 minutes")
		for range 2 {
			code, body := h.confirm(admin.token, ref)
			expectCode(t, code, body, http.StatusOK)
			decodeConfirmed(t, body, ref)
		}
		h.expectConfirmedOnce(t, ref)
	})

	t.Run("08 already paid, reservation since lapsed: 200, not 409", func(t *testing.T) {
		ref := h.rawOrderRef(t, h.investor(t).id, h.newAsset(t, villaBelair()), 4, "-1 minute")
		code, body := h.confirm(admin.token, ref)
		expectCode(t, code, body, http.StatusOK)
		decodeConfirmed(t, body, ref)
		h.expectUntouched(t, ref, 4)
	})

	t.Run("09 order placed through the API, then confirmed", func(t *testing.T) {
		asset := h.newAsset(t, villaBelair())
		investor := h.investor(t)
		code, body := h.post(orderCall{token: investor.token, key: newKey(), body: orderBody(asset, 10)})
		expectCode(t, code, body, http.StatusCreated)
		placed := decodeOrder(t, body)

		code, body = h.confirm(admin.token, placed.OrderRef)
		expectCode(t, code, body, http.StatusOK)
		o := decodeConfirmed(t, body, placed.OrderRef)
		if o.AmountDue != placed.AmountDue || o.TokenQuantity != placed.TokenQuantity {
			t.Fatalf("confirmed %s for %d shares, placed %s for %d", o.AmountDue, o.TokenQuantity, placed.AmountDue, placed.TokenQuantity)
		}
		h.expectConfirmedOnce(t, placed.OrderRef)

		// The investor sees the new status on the detail route.
		code, body = h.get(investor.token, ordersPath+"/"+placed.OrderRef)
		expectCode(t, code, body, http.StatusOK)
		if d := decodeOrderDetail(t, body); d.StatusCode != string(server.PAYMENTCONFIRMED) {
			t.Fatalf("detail status = %s, want PAYMENT_CONFIRMED", d.StatusCode)
		}
	})

	t.Run("refusals", func(t *testing.T) {
		asset := h.newAsset(t, villaBelair())
		investor := h.investor(t)

		t.Run("10 unknown reference is a 404", func(t *testing.T) {
			code, body := h.confirm(admin.token, h.unique("ISS-UNKNOWN"))
			expectCode(t, code, body, http.StatusNotFound)
			expectMessage(t, body)
		})

		t.Run("11 lapsed reservation is a 409, nothing written", func(t *testing.T) {
			ref := h.rawOrderRef(t, investor.id, asset, 2, "-1 minute")
			code, body := h.confirm(admin.token, ref)
			expectCode(t, code, body, http.StatusConflict)
			expectMessage(t, body)
			h.expectUntouched(t, ref, 2)
		})

		// A still-valid token of a deleted ADMIN: the payment insert breaks its
		// foreign key. A technical failure is a 500 with the generic message,
		// and the rollback leaves nothing behind.
		t.Run("11b a technical failure is a 500, nothing written", func(t *testing.T) {
			ghost := h.admin(t)
			if _, err := h.db.Exec(`DELETE FROM usr.users WHERE id = $1`, ghost.id); err != nil {
				t.Fatalf("deleting the ADMIN: %v", err)
			}
			ref := h.rawOrderRef(t, investor.id, asset, 2, "15 minutes")
			code, body := h.confirm(ghost.token, ref)
			expectCode(t, code, body, http.StatusInternalServerError)
			expectMessage(t, body)
			h.expectUntouched(t, ref, 2)
		})

		for _, c := range []struct {
			name   string
			status int
		}{
			{"12 a CANCELLED order is a 409", 6},
			{"13 an EXPIRED order is a 409", 7},
			{"14 a CLOSED order is a 409", 5},
			{"15 a REFUNDED order is a 409", 14},
		} {
			t.Run(c.name, func(t *testing.T) {
				ref := h.rawOrderRef(t, investor.id, asset, c.status, "15 minutes")
				code, body := h.confirm(admin.token, ref)
				expectCode(t, code, body, http.StatusConflict)
				expectMessage(t, body)
				h.expectUntouched(t, ref, c.status)
			})
		}
	})
}

func TestConfirmPaymentConcurrency(t *testing.T) {
	h := newOrderHarness(t)
	admin := h.admin(t)

	t.Run("16 ten simultaneous confirmations: ten 200, one payment", func(t *testing.T) {
		ref := h.rawOrderRef(t, h.investor(t).id, h.newAsset(t, villaBelair()), 2, "15 minutes")
		var (
			wg    sync.WaitGroup
			start = make(chan struct{})
			codes = make([]int, 10)
		)
		for i := range codes {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				codes[i], _ = h.confirm(admin.token, ref)
			}()
		}
		close(start)
		wg.Wait()
		for i, c := range codes {
			if c != http.StatusOK {
				t.Fatalf("call %d: code %d, want 200 (codes = %v)", i, c, codes)
			}
		}
		h.expectConfirmedOnce(t, ref)
	})
}
