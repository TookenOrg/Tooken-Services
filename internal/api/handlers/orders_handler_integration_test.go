//go:build integration

// End-to-end coverage of POST /assets/real-estate/issuance/orders through the
// real router, middleware, handler, service and
// database. What the database layer proves on its own (every refusal, every
// rollback) lives in internal/orders/database; this file proves the contract a
// client sees: status codes, headers, the response body.
//
// iss.issuance_order_status_history is append-only: the orders written here can
// never be deleted. Run it on a THROWAWAY database built from the baseline:
//
//	docker run -d --rm --name zz-orders -e POSTGRES_PASSWORD=pw -e POSTGRES_DB=tk \
//	    -p 55432:5432 postgres:16
//	docker exec -i zz-orders psql -U postgres -d tk < tools/db/baseline.sql
//	TEST_DATABASE_URL="postgres://postgres:pw@localhost:55432/tk?sslmode=disable" \
//	    go test -tags integration -run 'TestIssuanceOrder' ./internal/api/handlers/
//
// The concurrent cases (17, 20, 21) live in their own function so that they can
// be repeated alone with -count=20:
//
//	go test -tags integration -race -count=20 -run 'TestIssuanceOrderConcurrency' ./internal/api/handlers/
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/middleware"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const ordersPath = "/assets/real-estate/issuance/orders"

// orderHarness is the router plus the database the fixtures write to.
type orderHarness struct {
	t   *testing.T
	db  *sql.DB
	r   *gin.Engine
	run int64
	seq atomic.Int64
}

func newOrderHarness(t *testing.T) *orderHarness {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	logger.Init(false)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	globals.DB = db

	if err := middleware.InitAuth("../../../api/openapi.yaml"); err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	g := r.Group(globals.BaseURL)
	g.Use(middleware.AutoAuthMiddleware())
	server.RegisterHandlers(g, NewHandler())

	return &orderHarness{t: t, db: db, r: r, run: time.Now().UnixNano() % 1_000_000_000}
}

func (h *orderHarness) unique(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, h.run, h.seq.Add(1))
}

// ---------------------------------------------------------------- fixtures

type orderUser struct {
	id    int
	token string
}

// newUser writes a user with the given role and KYC. kycExpires is an SQL
// expression for kyc_expires_at ("NULL" = never expires).
func (h *orderHarness) newUser(t *testing.T, role, kycStatus, kycExpires string) orderUser {
	t.Helper()
	verifiedAt := "NULL"
	if kycStatus == "verified" {
		verifiedAt = "now() - interval '1 day'"
	}
	email := h.unique("order") + "@tooken.test"
	var id int
	if err := h.db.QueryRow(fmt.Sprintf(`
	    INSERT INTO usr.users (full_name, email, password, role, kyc_status, kyc_verified_at, kyc_expires_at)
	    VALUES ('Order test', $1, 'x', $2, $3, %s, %s)
	    RETURNING id`, verifiedAt, kycExpires), email, role, kycStatus).Scan(&id); err != nil {
		t.Fatalf("creating user: %v", err)
	}
	token, err := authUtils.GenerateJWT(id, email, role)
	if err != nil {
		t.Fatal(err)
	}
	return orderUser{id: id, token: token}
}

// investor is the reference investor: a USER whose KYC is verified for a year.
func (h *orderHarness) investor(t *testing.T) orderUser {
	return h.newUser(t, authUtils.RoleUser, "verified", "now() + interval '1 year'")
}

type orderAsset struct {
	status int     // ass.real_estate_status id: 1 draft, 3 published, 4 fundraising
	total  int64   // total_shares
	price  string  // price_per_share
	rate   *string // entry_fee_rate, nil = NULL
	cap    *int64  // max_shares_per_investor, nil = NULL
}

// villaBelair is the reference asset of these tests: 1 500 shares at 200.00 EUR,
// 2 % of entry fee, 50 shares per investor at most.
func villaBelair() orderAsset {
	rate, limit := "2", int64(50)
	return orderAsset{status: 4, total: 1500, price: "200", rate: &rate, cap: &limit}
}

func (h *orderHarness) newAsset(t *testing.T, a orderAsset) int {
	t.Helper()
	var issuer int
	if err := h.db.QueryRow(`
	    INSERT INTO ass.issuer (name, legal_form) VALUES ($1, 'SA') RETURNING id`,
		h.unique("issuer")).Scan(&issuer); err != nil {
		t.Fatalf("creating issuer: %v", err)
	}

	// A public status needs a token (real_estate_active_requires_token_ck).
	var token sql.NullInt64
	if a.status >= 3 && a.status <= 5 {
		var id int64
		if err := h.db.QueryRow(`
		    INSERT INTO blk.token (symbol, token_name, address, nb_decimal)
		    VALUES ('ORD', $1::text, '0x' || substr(md5($1::text), 1, 32) || '00000000', 0)
		    RETURNING id`, h.unique("token")).Scan(&id); err != nil {
			t.Fatalf("creating token: %v", err)
		}
		token = sql.NullInt64{Int64: id, Valid: true}
	}

	var id int
	if err := h.db.QueryRow(`
	    INSERT INTO ass.real_estate (title, estate_type, issuer_id, status_id, token_id)
	    SELECT $1, min(id), $2, $3, $4 FROM ass.real_estate_type
	    RETURNING id`, h.unique("Villa Belair"), issuer, a.status, token).Scan(&id); err != nil {
		t.Fatalf("creating asset: %v", err)
	}
	if _, err := h.db.Exec(`
	    INSERT INTO ass.real_estate_shares_config
	        (real_estate_id, total_shares, price_per_share, entry_fee_rate, max_shares_per_investor)
	    VALUES ($1, $2, $3, $4, $5)`, id, a.total, a.price, a.rate, a.cap); err != nil {
		t.Fatalf("creating shares config: %v", err)
	}
	return id
}

// rawOrder writes an order directly in SQL, in any status, to set up the
// stock. expiresIn is a PostgreSQL interval added to now() for
// reservation_expires_at ("-1 minute" for a lapsed reservation).
func (h *orderHarness) rawOrder(t *testing.T, user, asset int, quantity int64, status int, expiresIn string) {
	t.Helper()
	if _, err := h.db.Exec(`
	    INSERT INTO iss.issuance_orders (
	        user_id, asset_id, quantity, status_id, order_reference,
	        unit_price, currency_code, gross_amount, fee_amount, amount_due,
	        created_at, reservation_expires_at, delivery_tx_hash, refund_reason)
	    VALUES (
	        $1, $2, $3::bigint, $4::int, $5::text,
	        1, 'EUR', $3::numeric, 0, $3::numeric,
	        now() - interval '1 hour', now() + $6::interval,
	        CASE WHEN $4::int IN (5, 10, 11, 12) THEN '0x' || md5($5::text) || md5($5::text || '#') END,
	        CASE WHEN $4::int IN (13, 14) THEN 'PLATFORM' END)`,
		user, asset, quantity, status, h.unique("RAW"), expiresIn); err != nil {
		t.Fatalf("creating raw order (status %d): %v", status, err)
	}
}

// ---------------------------------------------------------------- requests

type orderCall struct {
	token string
	key   string // "" = no Idempotency-Key header
	body  string
}

func (h *orderHarness) post(c orderCall) (int, string) {
	req := httptest.NewRequest(http.MethodPost, globals.BaseURL+ordersPath, strings.NewReader(c.body))
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if c.key != "" {
		req.Header.Set("Idempotency-Key", c.key)
	}
	w := httptest.NewRecorder()
	h.r.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

func orderBody(asset, quantity int) string {
	return fmt.Sprintf(`{"realEstateId": %d, "quantity": %d}`, asset, quantity)
}

func newKey() string { return uuid.NewString() }

func expectCode(t *testing.T, code int, body string, want int) {
	t.Helper()
	if code != want {
		t.Fatalf("want %d, got %d: %s", want, code, body)
	}
}

func decodeOrder(t *testing.T, body string) server.IssuanceOrder {
	t.Helper()
	var resp server.CreateIssuanceOrderResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("decoding %s: %v", body, err)
	}
	if resp.Data == nil {
		t.Fatalf("no order in %s", body)
	}
	return *resp.Data
}

func (h *orderHarness) count(t *testing.T, query string, args ...any) int {
	t.Helper()
	var n int
	if err := h.db.QueryRow(query, args...).Scan(&n); err != nil {
		t.Fatalf("%v\n%s", err, query)
	}
	return n
}

func (h *orderHarness) ordersOf(t *testing.T, user int) int {
	t.Helper()
	return h.count(t, `SELECT count(*) FROM iss.issuance_orders WHERE user_id = $1`, user)
}

// reservedOn is what the order placement counts against the stock.
func (h *orderHarness) reservedOn(t *testing.T, asset int) int {
	t.Helper()
	return h.count(t, `
	    SELECT COALESCE(SUM(o.quantity), 0)
	    FROM iss.issuance_orders o
	    JOIN iss.issuance_order_statuses s ON s.id = o.status_id
	    WHERE o.asset_id = $1 AND s.counts_as_reserved
	      AND (o.status_id <> 2 OR o.reservation_expires_at > now())`, asset)
}

// ---------------------------------------------------------------- cases

func TestIssuanceOrderEndpoint(t *testing.T) {
	h := newOrderHarness(t)

	// A reservation of 20 minutes rather than the default 15 proves the
	// variable reaches the order, not just the default.
	t.Setenv("ORDER_RESERVATION_TTL", "20m")

	t.Run("who may order", func(t *testing.T) {
		asset := h.newAsset(t, villaBelair())

		t.Run("01 anonymous is a 401", func(t *testing.T) {
			code, body := h.post(orderCall{key: newKey(), body: orderBody(asset, 1)})
			expectCode(t, code, body, http.StatusUnauthorized)
		})

		t.Run("01b a forged token is a 401", func(t *testing.T) {
			code, body := h.post(orderCall{token: "not.a.jwt", key: newKey(), body: orderBody(asset, 1)})
			expectCode(t, code, body, http.StatusUnauthorized)
		})

		// U6: MANAGER and ADMIN do not order, even with a verified KYC.
		for _, role := range []string{authUtils.RoleManager, authUtils.RoleAdmin} {
			t.Run("02 a "+role+" is a 403", func(t *testing.T) {
				staff := h.newUser(t, role, "verified", "NULL")
				code, body := h.post(orderCall{token: staff.token, key: newKey(), body: orderBody(asset, 1)})
				expectCode(t, code, body, http.StatusForbidden)
				if n := h.ordersOf(t, staff.id); n != 0 {
					t.Fatalf("want no order, found %d", n)
				}
			})
		}

		for _, kyc := range []struct{ name, status, expires string }{
			{"03 a USER without KYC", "none", "NULL"},
			{"03b a USER with a pending KYC", "pending", "NULL"},
			{"03c a USER with an approved, not yet verified, KYC", "approved", "NULL"},
			{"03d a USER with a rejected KYC", "rejected", "NULL"},
			{"03e a USER with a revoked KYC", "revoked", "NULL"},
			{"04 a USER verified but past kyc_expires_at", "verified", "now() - interval '1 second'"},
		} {
			t.Run(kyc.name+" is a 403, nothing written", func(t *testing.T) {
				u := h.newUser(t, authUtils.RoleUser, kyc.status, kyc.expires)
				code, body := h.post(orderCall{token: u.token, key: newKey(), body: orderBody(asset, 1)})
				expectCode(t, code, body, http.StatusForbidden)
				if n := h.ordersOf(t, u.id); n != 0 {
					t.Fatalf("want no order, found %d", n)
				}
			})
		}

		t.Run("04b a verified KYC without expiry may order", func(t *testing.T) {
			u := h.newUser(t, authUtils.RoleUser, "verified", "NULL")
			code, body := h.post(orderCall{token: u.token, key: newKey(), body: orderBody(asset, 1)})
			expectCode(t, code, body, http.StatusCreated)
		})
	})

	t.Run("malformed requests are a 400, nothing written", func(t *testing.T) {
		asset := h.newAsset(t, villaBelair())
		alice := h.investor(t)

		for _, c := range []struct {
			name string
			call orderCall
		}{
			{"05 without Idempotency-Key", orderCall{body: orderBody(asset, 1)}},
			{"05b an Idempotency-Key that is not a UUID", orderCall{key: "my-key-1", body: orderBody(asset, 1)}},
			{"06 quantity 0", orderCall{key: newKey(), body: orderBody(asset, 0)}},
			{"06b a negative quantity", orderCall{key: newKey(), body: orderBody(asset, -3)}},
			{"06c a decimal quantity", orderCall{key: newKey(), body: fmt.Sprintf(`{"realEstateId": %d, "quantity": 1.5}`, asset)}},
			{"06d a body that is not JSON", orderCall{key: newKey(), body: `realEstateId=1`}},
		} {
			t.Run(c.name, func(t *testing.T) {
				c.call.token = alice.token
				code, body := h.post(c.call)
				expectCode(t, code, body, http.StatusBadRequest)
			})
		}
		if n := h.ordersOf(t, alice.id); n != 0 {
			t.Fatalf("want no order, found %d", n)
		}
	})

	t.Run("the asset must be open", func(t *testing.T) {
		alice := h.investor(t)

		t.Run("07 an unknown asset is a 404", func(t *testing.T) {
			code, body := h.post(orderCall{token: alice.token, key: newKey(), body: orderBody(999_999_999, 1)})
			expectCode(t, code, body, http.StatusNotFound)
		})

		// U7: what the public cannot see does not exist for it.
		t.Run("08 a draft is a 404", func(t *testing.T) {
			a := villaBelair()
			a.status = 1
			code, body := h.post(orderCall{token: alice.token, key: newKey(), body: orderBody(h.newAsset(t, a), 1)})
			expectCode(t, code, body, http.StatusNotFound)
		})

		t.Run("09 a published asset is a 409 that says the fundraising is not open", func(t *testing.T) {
			a := villaBelair()
			a.status = 3
			code, body := h.post(orderCall{token: alice.token, key: newKey(), body: orderBody(h.newAsset(t, a), 1)})
			expectCode(t, code, body, http.StatusConflict)
			if !strings.Contains(body, "fundraising is not open") {
				t.Fatalf("message does not explain the refusal: %s", body)
			}
		})

		t.Run("09b a funded asset is a 409", func(t *testing.T) {
			a := villaBelair()
			a.status = 5
			code, body := h.post(orderCall{token: alice.token, key: newKey(), body: orderBody(h.newAsset(t, a), 1)})
			expectCode(t, code, body, http.StatusConflict)
		})

		if n := h.ordersOf(t, alice.id); n != 0 {
			t.Fatalf("want no order, found %d", n)
		}
	})

	t.Run("the nominal order", func(t *testing.T) {
		asset := h.newAsset(t, villaBelair())
		alice := h.investor(t)
		key := newKey()
		call := orderCall{token: alice.token, key: key, body: orderBody(asset, 10)}

		var first server.IssuanceOrder
		t.Run("10 is created, priced and reserved", func(t *testing.T) {
			code, body := h.post(call)
			expectCode(t, code, body, http.StatusCreated)
			first = decodeOrder(t, body)

			for _, f := range []struct{ name, got, want string }{
				{"statusCode", first.StatusCode, "AWAITING_PAYMENT"},
				{"unitPrice", first.UnitPrice, "200.00000000"},
				{"currencyCode", first.CurrencyCode, "EUR"},
				{"grossAmount", first.GrossAmount, "2000.00"},
				{"feeAmount", first.FeeAmount, "40.00"},
				{"amountDue", first.AmountDue, "2040.00"},
			} {
				if f.got != f.want {
					t.Errorf("%s = %q, want %q", f.name, f.got, f.want)
				}
			}
			if first.RealEstateId != asset || first.TokenQuantity != 10 || first.UserId != alice.id {
				t.Errorf("order = asset %d, quantity %d, user %d", first.RealEstateId, first.TokenQuantity, first.UserId)
			}
			if !strings.HasPrefix(first.OrderRef, "ISS-") {
				t.Errorf("orderRef = %q", first.OrderRef)
			}
			want := server.IssuanceOrderFee{Code: server.ENTRY, Rate: "2.0000", BaseAmount: "2000.00", Amount: "40.00"}
			if len(first.Fees) != 1 || first.Fees[0] != want {
				t.Errorf("fees = %+v, want [%+v]", first.Fees, want)
			}
			// Both dates come from the database clock, in the same transaction.
			if ttl := first.ReservationExpiresAt.Sub(first.CreatedAt); ttl != 20*time.Minute {
				t.Errorf("reservation lasts %s, want the 20m of ORDER_RESERVATION_TTL", ttl)
			}
		})

		t.Run("10b one order, one fee line and two history rows in the database", func(t *testing.T) {
			if n := h.ordersOf(t, alice.id); n != 1 {
				t.Fatalf("want 1 order, found %d", n)
			}
			if n := h.count(t, `
			    SELECT count(*) FROM iss.issuance_order_fees f
			    JOIN iss.issuance_orders o ON o.id = f.order_id
			    WHERE o.order_reference = $1 AND f.fee_code = 'ENTRY'
			      AND f.rate = 2 AND f.base_amount = 2000 AND f.amount = 40`, first.OrderRef); n != 1 {
				t.Fatalf("want 1 ENTRY line of 40.00, found %d", n)
			}
			var history string
			if err := h.db.QueryRow(`
			    SELECT string_agg(COALESCE(h.from_status_id::text, 'NULL') || '>' || h.to_status_id
			                      || ' by ' || (h.actor_user_id = o.user_id), ', ' ORDER BY h.id)
			    FROM iss.issuance_order_status_history h
			    JOIN iss.issuance_orders o ON o.id = h.order_id
			    WHERE o.order_reference = $1`, first.OrderRef).Scan(&history); err != nil {
				t.Fatal(err)
			}
			if history != "NULL>1 by true, 1>2 by true" {
				t.Fatalf("history = %q, want CREATED then AWAITING_PAYMENT, by the investor", history)
			}
		})

		t.Run("10c the order is readable at the reference returned", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, globals.BaseURL+ordersPath+"/"+first.OrderRef, nil)
			req.Header.Set("Authorization", "Bearer "+alice.token)
			w := httptest.NewRecorder()
			h.r.ServeHTTP(w, req)
			expectCode(t, w.Code, w.Body.String(), http.StatusOK)
			if got := decodeOrder(t, w.Body.String()); got.AmountDue != "2040.00" || len(got.Fees) != 1 {
				t.Fatalf("read back = %+v", got)
			}
		})

		t.Run("11 a replay is a 200 with the same order, nothing created", func(t *testing.T) {
			code, body := h.post(call)
			expectCode(t, code, body, http.StatusOK)
			replay := decodeOrder(t, body)
			if replay.OrderRef != first.OrderRef || !replay.CreatedAt.Equal(first.CreatedAt) || replay.AmountDue != first.AmountDue {
				t.Fatalf("replay = %s created %s, want %s created %s",
					replay.OrderRef, replay.CreatedAt, first.OrderRef, first.CreatedAt)
			}
			if n := h.ordersOf(t, alice.id); n != 1 {
				t.Fatalf("want still 1 order, found %d", n)
			}
		})

		t.Run("12 the same key with another quantity is a 422, nothing created", func(t *testing.T) {
			code, body := h.post(orderCall{token: alice.token, key: key, body: orderBody(asset, 11)})
			expectCode(t, code, body, http.StatusUnprocessableEntity)
			if n := h.ordersOf(t, alice.id); n != 1 {
				t.Fatalf("want still 1 order, found %d", n)
			}
		})

		t.Run("12b the same key on another asset is a 422", func(t *testing.T) {
			other := h.newAsset(t, villaBelair())
			code, body := h.post(orderCall{token: alice.token, key: key, body: orderBody(other, 10)})
			expectCode(t, code, body, http.StatusUnprocessableEntity)
		})

		t.Run("12c another investor may use the same key", func(t *testing.T) {
			bob := h.investor(t)
			code, body := h.post(orderCall{token: bob.token, key: key, body: orderBody(asset, 10)})
			expectCode(t, code, body, http.StatusCreated)
		})
	})

	t.Run("stock", func(t *testing.T) {
		t.Run("13 5 shares left, 10 ordered: 409 that cites the 5", func(t *testing.T) {
			asset := h.newAsset(t, villaBelair())
			h.rawOrder(t, h.investor(t).id, asset, 1495, 4, "1 hour")
			alice := h.investor(t)
			code, body := h.post(orderCall{token: alice.token, key: newKey(), body: orderBody(asset, 10)})
			expectCode(t, code, body, http.StatusConflict)
			if !strings.Contains(body, "only 5 shares left") {
				t.Fatalf("message does not cite what is left: %s", body)
			}
			if n := h.ordersOf(t, alice.id); n != 0 {
				t.Fatalf("want no order, found %d", n)
			}
		})

		t.Run("13b the exact last 5 shares are a 201", func(t *testing.T) {
			asset := h.newAsset(t, villaBelair())
			h.rawOrder(t, h.investor(t).id, asset, 1495, 4, "1 hour")
			code, body := h.post(orderCall{token: h.investor(t).token, key: newKey(), body: orderBody(asset, 5)})
			expectCode(t, code, body, http.StatusCreated)
		})

		t.Run("14 lapsed reservations free their shares", func(t *testing.T) {
			asset := h.newAsset(t, villaBelair())
			h.rawOrder(t, h.investor(t).id, asset, 1495, 2, "-1 minute")

			code, body := h.post(orderCall{token: h.investor(t).token, key: newKey(), body: orderBody(asset, 10)})
			expectCode(t, code, body, http.StatusCreated)

			// The detail shown on the asset must agree with the stock the order
			// was checked against: 10 sold, not 1 505.
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/assets/real-estates/%d", globals.BaseURL, asset), nil)
			w := httptest.NewRecorder()
			h.r.ServeHTTP(w, req)
			expectCode(t, w.Code, w.Body.String(), http.StatusOK)
			var re server.RealEstate
			if err := json.Unmarshal(w.Body.Bytes(), &re); err != nil {
				t.Fatal(err)
			}
			if re.Progression == nil || re.Progression.TokensSold == nil {
				t.Fatalf("no tokens_sold in %s", w.Body.String())
			}
			if sold := *re.Progression.TokensSold; sold != 10 {
				t.Fatalf("tokens_sold = %d, want 10: the lapsed 1 495 are still counted", sold)
			}
		})
	})

	t.Run("per-investor cap", func(t *testing.T) {
		t.Run("15 45 held, 10 ordered, cap 50: 409", func(t *testing.T) {
			asset := h.newAsset(t, villaBelair())
			alice := h.investor(t)
			h.rawOrder(t, alice.id, asset, 45, 2, "10 minutes")
			code, body := h.post(orderCall{token: alice.token, key: newKey(), body: orderBody(asset, 10)})
			expectCode(t, code, body, http.StatusConflict)
			if !strings.Contains(body, "45 of 50") {
				t.Fatalf("message does not cite the cap: %s", body)
			}
		})

		t.Run("15b 45 held, 5 ordered, cap 50: 201", func(t *testing.T) {
			asset := h.newAsset(t, villaBelair())
			alice := h.investor(t)
			h.rawOrder(t, alice.id, asset, 45, 2, "10 minutes")
			code, body := h.post(orderCall{token: alice.token, key: newKey(), body: orderBody(asset, 5)})
			expectCode(t, code, body, http.StatusCreated)
		})

		t.Run("16 no cap: an order of 1 000 is a 201", func(t *testing.T) {
			a := villaBelair()
			a.cap = nil
			code, body := h.post(orderCall{token: h.investor(t).token, key: newKey(), body: orderBody(h.newAsset(t, a), 1000)})
			expectCode(t, code, body, http.StatusCreated)
		})

		// U8: a closed order still holds its shares, a cancelled one does not.
		t.Run("19 a closed 45 counts, a cancelled 30 does not: 409", func(t *testing.T) {
			asset := h.newAsset(t, villaBelair())
			alice := h.investor(t)
			h.rawOrder(t, alice.id, asset, 45, 5, "0")
			h.rawOrder(t, alice.id, asset, 30, 6, "0")
			code, body := h.post(orderCall{token: alice.token, key: newKey(), body: orderBody(asset, 10)})
			expectCode(t, code, body, http.StatusConflict)
			if !strings.Contains(body, "45 of 50") {
				t.Fatalf("want the closed 45 counted, the cancelled 30 not: %s", body)
			}
		})
	})

	t.Run("fees", func(t *testing.T) {
		zero := "0"
		for _, c := range []struct {
			name string
			rate *string
		}{
			{"18 no fee rate", nil},
			{"18b a fee rate of 0 (U4)", &zero},
		} {
			t.Run(c.name+": feeAmount 0.00, fees empty, no line in the database", func(t *testing.T) {
				a := villaBelair()
				a.rate = c.rate
				code, body := h.post(orderCall{token: h.investor(t).token, key: newKey(), body: orderBody(h.newAsset(t, a), 10)})
				expectCode(t, code, body, http.StatusCreated)
				o := decodeOrder(t, body)
				if o.FeeAmount != "0.00" || o.AmountDue != "2000.00" {
					t.Errorf("feeAmount %q, amountDue %q", o.FeeAmount, o.AmountDue)
				}
				// [] and not null: the contract says "empty when there is no fee".
				if !strings.Contains(body, `"fees":[]`) {
					t.Errorf("fees is not an empty array: %s", body)
				}
				if n := h.count(t, `
				    SELECT count(*) FROM iss.issuance_order_fees f
				    JOIN iss.issuance_orders o ON o.id = f.order_id
				    WHERE o.order_reference = $1`, o.OrderRef); n != 0 {
					t.Errorf("want no fee line, found %d", n)
				}
			})
		}
	})
}

// waitForLockWaiter blocks until a session waits on a lock, so a concurrent
// case acts once the request under test is provably queued, not after a sleep.
func (h *orderHarness) waitForLockWaiter(t *testing.T) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if h.count(t, `
		    SELECT count(*) FROM pg_stat_activity
		    WHERE datname = current_database() AND wait_event_type = 'Lock'`) >= 1 {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("timed out waiting for the request to queue on the lock")
}

type orderResponse struct {
	code int
	body string
}

func (h *orderHarness) postAsync(c orderCall) <-chan orderResponse {
	out := make(chan orderResponse, 1)
	go func() {
		code, body := h.post(c)
		out <- orderResponse{code, body}
	}()
	return out
}

// The cases that must be run with -count=20.
func TestIssuanceOrderConcurrency(t *testing.T) {
	h := newOrderHarness(t)

	t.Run("17 two investors race for the last 5 shares: one 201, one 409", func(t *testing.T) {
		asset := h.newAsset(t, villaBelair())
		h.rawOrder(t, h.investor(t).id, asset, 1495, 4, "1 hour")
		alice, bob := h.investor(t), h.investor(t)

		var (
			wg    sync.WaitGroup
			start = make(chan struct{})
			codes = make([]int, 2)
		)
		for i, u := range []orderUser{alice, bob} {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				codes[i], _ = h.post(orderCall{token: u.token, key: newKey(), body: orderBody(asset, 5)})
			}()
		}
		close(start)
		wg.Wait()

		created, refused := 0, 0
		for _, c := range codes {
			switch c {
			case http.StatusCreated:
				created++
			case http.StatusConflict:
				refused++
			}
		}
		if created != 1 || refused != 1 {
			t.Fatalf("codes = %v, want one 201 and one 409", codes)
		}
		if got := h.reservedOn(t, asset); got != 1500 {
			t.Fatalf("reserved = %d, want exactly the 1 500 shares", got)
		}
	})

	// What the order waited for is what it is priced with.
	t.Run("20 a price changed while the order waits is the price charged", func(t *testing.T) {
		asset := h.newAsset(t, villaBelair())
		manager, err := h.db.Begin()
		if err != nil {
			t.Fatal(err)
		}
		defer manager.Rollback()
		if _, err := manager.Exec(`SELECT 1 FROM ass.real_estate WHERE id = $1 FOR UPDATE`, asset); err != nil {
			t.Fatal(err)
		}
		if _, err := manager.Exec(`UPDATE ass.real_estate_shares_config SET price_per_share = 250 WHERE real_estate_id = $1`, asset); err != nil {
			t.Fatal(err)
		}

		result := h.postAsync(orderCall{token: h.investor(t).token, key: newKey(), body: orderBody(asset, 10)})
		h.waitForLockWaiter(t)
		if err := manager.Commit(); err != nil {
			t.Fatal(err)
		}

		r := <-result
		expectCode(t, r.code, r.body, http.StatusCreated)
		o := decodeOrder(t, r.body)
		if o.UnitPrice != "250.00000000" || o.GrossAmount != "2500.00" {
			t.Fatalf("priced at %s (gross %s), want 250.00000000 (2500.00)", o.UnitPrice, o.GrossAmount)
		}
	})

	t.Run("21 a fundraising closed while the order waits is a 409, nothing written", func(t *testing.T) {
		asset := h.newAsset(t, villaBelair())
		manager, err := h.db.Begin()
		if err != nil {
			t.Fatal(err)
		}
		defer manager.Rollback()
		if _, err := manager.Exec(`SELECT 1 FROM ass.real_estate WHERE id = $1 FOR UPDATE`, asset); err != nil {
			t.Fatal(err)
		}
		if _, err := manager.Exec(`UPDATE ass.real_estate SET status_id = 3 WHERE id = $1`, asset); err != nil {
			t.Fatal(err)
		}

		alice := h.investor(t)
		result := h.postAsync(orderCall{token: alice.token, key: newKey(), body: orderBody(asset, 10)})
		h.waitForLockWaiter(t)
		if err := manager.Commit(); err != nil {
			t.Fatal(err)
		}

		r := <-result
		expectCode(t, r.code, r.body, http.StatusConflict)
		if n := h.ordersOf(t, alice.id); n != 0 {
			t.Fatalf("want no order, found %d", n)
		}
	})
}
