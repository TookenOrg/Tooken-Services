//go:build integration

// End-to-end coverage of the real estate write endpoints, against a real
// PostgreSQL migrated to the latest version. Nothing here is mocked: the
// transactions, the triggers of migrations 000006 and 000008 and the CHECK
// constraints of 000011 are part of what is being tested, and none of them
// exists in a fake.
//
// Run it with a disposable database:
//
//	docker run -d --name pgtest -e POSTGRES_PASSWORD=pw -e POSTGRES_DB=tooken \
//	    -p 55444:5432 postgres:16
//	# apply tools/db/migrations in order, then:
//	TEST_DATABASE_URL="postgres://postgres:pw@localhost:55444/tooken?sslmode=disable" \
//	    go test -tags integration ./internal/api/handlers/
package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/middleware"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func TestRealEstateWriteEndpoints(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	logger.Init(true)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	globals.DB = db
	if err := middleware.InitAuth("../../../api/openapi.yaml"); err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	g := r.Group(globals.BaseURL)
	g.Use(middleware.AutoAuthMiddleware())
	server.RegisterHandlers(g, NewHandler())

	// The referential is empty in a freshly migrated database, and every write
	// below points at type 1. Seeding it here keeps the suite runnable from the
	// migrations alone, without a manual step to forget.
	if _, err := db.Exec(`
	    INSERT INTO ass.real_estate_type (id, name) VALUES (1, 'Apartment')
	    ON CONFLICT (id) DO NOTHING`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`
	    INSERT INTO ass.payment_frequency_type (id, name) VALUES (1, 'Year')
	    ON CONFLICT (id) DO NOTHING`); err != nil {
		t.Fatal(err)
	}

	mgr, _ := authUtils.GenerateJWT(1, "m@t.lu", authUtils.RoleManager)
	usr, _ := authUtils.GenerateJWT(2, "u@t.lu", authUtils.RoleUser)

	do := func(method, url, token, body string) (int, string) {
		var b *bytes.Buffer = bytes.NewBufferString(body)
		req := httptest.NewRequest(method, globals.BaseURL+url, b)
		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code, w.Body.String()
	}

	// An issuer is part of what a published offer must name, so the referential
	// needs one before an asset can be born visible.
	if _, err := db.Exec(`
	    INSERT INTO ass.issuer (id, name, legal_form) VALUES (1, 'Tooken RE I', 'SA')
	    ON CONFLICT (id) DO NOTHING`); err != nil {
		t.Fatal(err)
	}

	payload := `{
      "title": "Villa Belair",
      "description": "Nice",
      "estate_type_id": 1,
      "issuer_id": 1,
      "address": {"street":"12 rue de la Gare","postal_code":"L-1611","city":"Luxembourg","country_code":"lu","latitude":"49.611622","longitude":"6.131935"},
      "specification": {"surface_area":"128.50","bedroom_number":4,"energy_class":"c"},
      "configuration": {"total_shares":"1500","price_per_share":"199.99","currency_code":"eur","yield":"4.25","payment_frequency":1,"payment_frequency_type_id":1,"entry_fee_rate":"2.5"},
      "media": [{"url":"https://x/1.jpg","is_cover":true},{"url":"https://x/2.jpg"}]
    }`

	t.Run("anonymous cannot create", func(t *testing.T) {
		code, body := do("POST", "/assets/real-estates", "", payload)
		t.Logf("%d %s", code, body)
		if code != http.StatusUnauthorized {
			t.Errorf("want 401 got %d", code)
		}
	})
	t.Run("plain user cannot create", func(t *testing.T) {
		code, body := do("POST", "/assets/real-estates", usr, payload)
		t.Logf("%d %s", code, body)
		if code != http.StatusForbidden {
			t.Errorf("want 403 got %d", code)
		}
	})

	var id int
	t.Run("manager creates", func(t *testing.T) {
		code, body := do("POST", "/assets/real-estates", mgr, payload)
		t.Logf("%d %s", code, body)
		if code != http.StatusCreated {
			t.Fatalf("want 201 got %d", code)
		}
		var re map[string]any
		json.Unmarshal([]byte(body), &re)
		id = int(re["id"].(float64))
		if re["active"] != true || re["status"] != "published" {
			t.Errorf("asset not born visible: %v %v", re["active"], re["status"])
		}
		cfg := re["configuration"].(map[string]any)
		if cfg["total_valuation"] != "299985" || cfg["price_per_share"] != "199.99" || cfg["currency_code"] != "EUR" {
			t.Errorf("config: %v", cfg)
		}
		if re["address"] == nil {
			t.Errorf("manager should see the address")
		}
	})

	t.Run("public detail hides the address", func(t *testing.T) {
		code, body := do("GET", "/assets/real-estates/"+itoa(id), "", "")
		t.Logf("%d %s", code, body)
		if code != 200 {
			t.Fatalf("want 200 got %d", code)
		}
		if bytes.Contains([]byte(body), []byte("rue de la Gare")) {
			t.Errorf("address leaked to anonymous")
		}
		if !bytes.Contains([]byte(body), []byte(`"city":"Luxembourg"`)) {
			t.Errorf("coarse location missing")
		}
	})

	t.Run("invalid payloads", func(t *testing.T) {
		cases := map[string]string{
			"fractional shares": `{"title":"x","address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"},"configuration":{"total_shares":"10.5","price_per_share":"1"}}`,
			"zero price":        `{"title":"x","address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"},"configuration":{"total_shares":"10","price_per_share":"0"}}`,
			// A price used to be accepted alone and stored as EUR, a currency
			// the manager never chose and had no reason to double-check.
			"price without currency": `{"title":"x","address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"},"configuration":{"total_shares":"10","price_per_share":"1"}}`,
			// A yield with no rhythm cannot be compared to another offer.
			"yield without a rhythm": `{"title":"x","address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"},"configuration":{"total_shares":"10","price_per_share":"1","currency_code":"EUR","yield":"4"}}`,
			"bad country":            `{"title":"x","address":{"street":"a","postal_code":"b","city":"c","country_code":"LUX"}}`,
			"half coordinates":       `{"title":"x","address":{"street":"a","postal_code":"b","city":"c","country_code":"LU","latitude":"49.1"}}`,
			"two covers":             `{"title":"x","address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"},"media":[{"url":"a","is_cover":true},{"url":"b","is_cover":true}]}`,
			"unknown estate type":    `{"title":"x","estate_type_id":999,"address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"}}`,
			"empty title":            `{"title":"   ","address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"}}`,
			"no estate type":         `{"title":"x","address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"}}`,
			"title over 255":         `{"title":"` + strings.Repeat("A", 256) + `","estate_type_id":1,"address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"}}`,
		}
		for name, p := range cases {
			code, body := do("POST", "/assets/real-estates", mgr, p)
			t.Logf("%-18s -> %d %s", name, code, body)
			if code != http.StatusBadRequest {
				t.Errorf("%s: want 400 got %d", name, code)
			}
		}
	})

	// The whole point of a PATCH: an editor that shows one section must not
	// erase the sections it never displayed.
	t.Run("patch keeps what it does not mention", func(t *testing.T) {
		p := `{"title":"Villa Belair v2","configuration":{"price_per_share":"250"}}`
		code, body := do("PATCH", "/assets/real-estates/"+itoa(id), mgr, p)
		t.Logf("%d %s", code, body)
		if code != 200 {
			t.Fatalf("want 200 got %d", code)
		}

		var re map[string]any
		json.Unmarshal([]byte(body), &re)

		if re["title"] != "Villa Belair v2" {
			t.Errorf("title not applied: %v", re["title"])
		}

		cfg := re["configuration"].(map[string]any)
		if cfg["price_per_share"] != "250" {
			t.Errorf("price not applied: %v", cfg["price_per_share"])
		}
		// Untouched, and the derived valuation follows the new price.
		if cfg["total_shares"] != "1500" {
			t.Errorf("total_shares lost by the patch: %v", cfg["total_shares"])
		}
		if cfg["total_valuation"] != "375000" {
			t.Errorf("valuation not recomputed: %v", cfg["total_valuation"])
		}
		if cfg["currency_code"] != "EUR" {
			t.Errorf("currency lost by the patch: %v", cfg["currency_code"])
		}

		if re["specification"] == nil {
			t.Errorf("specification erased by a patch that never mentioned it")
		}
		if re["media"] == nil || len(re["media"].([]any)) != 2 {
			t.Errorf("gallery erased by a patch that never mentioned it: %v", re["media"])
		}
		if re["address"].(map[string]any)["street"] != "12 rue de la Gare" {
			t.Errorf("address erased by a patch that never mentioned it: %v", re["address"])
		}
		if re["description"] != "Nice" {
			t.Errorf("description erased by a patch that never mentioned it: %v", re["description"])
		}
	})

	// A patch touching one address field must not blank the others.
	t.Run("patch merges inside a section", func(t *testing.T) {
		code, body := do("PATCH", "/assets/real-estates/"+itoa(id), mgr, `{"address":{"city":"Esch"}}`)
		if code != 200 {
			t.Fatalf("want 200 got %d: %s", code, body)
		}

		var re map[string]any
		json.Unmarshal([]byte(body), &re)
		addr := re["address"].(map[string]any)

		if addr["city"] != "Esch" {
			t.Errorf("city not applied: %v", addr["city"])
		}
		if addr["street"] != "12 rue de la Gare" || addr["postal_code"] != "L-1611" {
			t.Errorf("sibling fields lost: %v", addr)
		}
		if addr["latitude"] != "49.611622" {
			t.Errorf("coordinates lost: %v", addr["latitude"])
		}
	})

	// Clearing has to stay possible, otherwise a value entered by mistake could
	// never be removed. An empty string is the way to say it.
	t.Run("empty string clears an optional value", func(t *testing.T) {
		code, body := do("PATCH", "/assets/real-estates/"+itoa(id), mgr,
			`{"specification":{"energy_class":""}}`)
		if code != 200 {
			t.Fatalf("want 200 got %d: %s", code, body)
		}

		var re map[string]any
		json.Unmarshal([]byte(body), &re)
		spec, _ := re["specification"].(map[string]any)
		if spec != nil && spec["energy_class"] != nil {
			t.Errorf("energy class not cleared: %v", spec["energy_class"])
		}
		if re["title"] != "Villa Belair v2" {
			t.Errorf("title lost: %v", re["title"])
		}
	})

	// The counterpart of that freedom: what a published asset promises to an
	// investor cannot be taken away by a patch. Antony's arbitration was to
	// refuse rather than silently unpublish — an asset can be mid-fundraising,
	// and pulling it off the site over a mistyped field is worse than asking
	// for the field back.
	t.Run("a published asset cannot lose what investors rely on", func(t *testing.T) {
		for _, patch := range []string{`{"description":""}`, `{"configuration":{"yield":""}}`} {
			code, body := do("PATCH", "/assets/real-estates/"+itoa(id), mgr, patch)
			if code != http.StatusConflict {
				t.Errorf("%s: want 409 got %d: %s", patch, code, body)
			}
		}
	})

	// The merged asset is validated as a whole: a patch cannot slip the asset
	// into a state a create would have refused.
	t.Run("patch validates the merged result", func(t *testing.T) {
		code, body := do("PATCH", "/assets/real-estates/"+itoa(id), mgr, `{"configuration":{"total_shares":"10.5"}}`)
		t.Logf("%d %s", code, body)
		if code != http.StatusBadRequest {
			t.Errorf("want 400 got %d", code)
		}

		code, body = do("PATCH", "/assets/real-estates/"+itoa(id), mgr, `{"address":{"country_code":"LUX"}}`)
		t.Logf("%d %s", code, body)
		if code != http.StatusBadRequest {
			t.Errorf("want 400 got %d", code)
		}
	})

	t.Run("patch replaces the gallery as a whole", func(t *testing.T) {
		code, body := do("PATCH", "/assets/real-estates/"+itoa(id), mgr, `{"media":[{"url":"https://x/only.jpg"}]}`)
		if code != 200 {
			t.Fatalf("want 200 got %d: %s", code, body)
		}

		var re map[string]any
		json.Unmarshal([]byte(body), &re)
		media := re["media"].([]any)
		if len(media) != 1 || media[0].(map[string]any)["url"] != "https://x/only.jpg" {
			t.Errorf("gallery not replaced: %v", media)
		}
	})

	t.Run("delete then 404", func(t *testing.T) {
		code, body := do("DELETE", "/assets/real-estates/"+itoa(id), mgr, "")
		t.Logf("delete %d %s", code, body)
		if code != http.StatusNoContent {
			t.Fatalf("want 204 got %d", code)
		}
		code, _ = do("GET", "/assets/real-estates/"+itoa(id), mgr, "")
		if code != http.StatusNotFound {
			t.Errorf("deleted asset still readable by a manager: %d", code)
		}
		code, _ = do("DELETE", "/assets/real-estates/"+itoa(id), mgr, "")
		if code != http.StatusNotFound {
			t.Errorf("second delete: want 404 got %d", code)
		}
		code, _ = do("PATCH", "/assets/real-estates/"+itoa(id), mgr, payload)
		if code != http.StatusNotFound {
			t.Errorf("patch of a deleted asset: want 404 got %d", code)
		}
	})

	// Fixtures must never fail silently: a rejected INSERT would leave the guard
	// with nothing to guard, and the assertion below would pass for the wrong
	// reason.
	mustExec := func(t *testing.T, query string, args ...any) {
		t.Helper()
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatalf("fixture failed: %v", err)
		}
	}
	mustScanID := func(t *testing.T, query string, args ...any) int {
		t.Helper()
		var id int
		if err := db.QueryRow(query, args...).Scan(&id); err != nil {
			t.Fatalf("fixture failed: %v", err)
		}
		return id
	}

	t.Run("guards on a tokenized asset", func(t *testing.T) {
		tokenID := mustScanID(t, `INSERT INTO blk.token (address, token_name, symbol, nb_decimal)
		    VALUES ('0x' || lpad(md5(random()::text), 40, '0'),'G','G',0) RETURNING id`)
		assetID := mustScanID(t, `INSERT INTO ass.real_estate (title, estate_type, status_id, token_id)
		    VALUES ('Tokenized guard', 1, 3, $1) RETURNING id`, tokenID)
		mustExec(t, `INSERT INTO ass.real_estate_shares_config (real_estate_id, total_shares, price_per_share) VALUES ($1, 1000, 10)`, assetID)

		p := `{"configuration":{"total_shares":"500"}}`
		code, body := do("PATCH", "/assets/real-estates/"+itoa(assetID), mgr, p)
		t.Logf("shrink shares on tokenized -> %d %s", code, body)
		if code != http.StatusConflict {
			t.Errorf("want 409 got %d", code)
		}

		code, body = do("DELETE", "/assets/real-estates/"+itoa(assetID), mgr, "")
		t.Logf("delete tokenized -> %d %s", code, body)
		if code != http.StatusConflict {
			t.Errorf("want 409 got %d", code)
		}
	})

	// A share already reserved by an order is a commitment: the asset can grow,
	// but it cannot shrink under what investors were promised.
	t.Run("total shares cannot fall under the reserved ones", func(t *testing.T) {
		assetID := mustScanID(t, `INSERT INTO ass.real_estate (title, description, imageurl, estate_type, issuer_id, status_id)
		    VALUES ('Reserved guard', 'A guard', 'https://x/g.jpg', 1, 1, 3) RETURNING id`)
		mustExec(t, `INSERT INTO ass.real_estate_shares_config (real_estate_id, total_shares, price_per_share, yield) VALUES ($1, 1000, 10, 4)`, assetID)
		// status 2 = RESERVED, which carries counts_as_reserved. The reference is
		// derived from the asset so the suite can run twice on the same database:
		// order_ref is unique.
		mustExec(t, `INSERT INTO iss.issuance_orders (asset_id, quantity, status_id, order_ref)
		    VALUES ($1, 400, 2, $2)`, assetID, "ORD-RESERVED-"+itoa(assetID))

		base := `{"configuration":{"total_shares":"%s"}}`

		code, body := do("PATCH", "/assets/real-estates/"+itoa(assetID), mgr, fmt.Sprintf(base, "300"))
		t.Logf("shrink under reserved -> %d %s", code, body)
		if code != http.StatusConflict {
			t.Errorf("want 409 got %d", code)
		}

		code, body = do("PATCH", "/assets/real-estates/"+itoa(assetID), mgr, fmt.Sprintf(base, "2000"))
		t.Logf("grow above reserved -> %d %s", code, body)
		if code != http.StatusOK {
			t.Errorf("growing the offer must stay allowed, got %d", code)
		}
	})

	// Everything below reproduces what the production table used to reject with
	// a 500, before migration 000012 aligned it with the contract.
	t.Run("the shape the contract promises is the shape the table has", func(t *testing.T) {
		t.Run("an asset without description or image is created", func(t *testing.T) {
			p := `{"title":"Bare asset","estate_type_id":1,"address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"}}`
			code, body := do("POST", "/assets/real-estates", mgr, p)
			if code != http.StatusCreated {
				t.Errorf("want 201 got %d: %s", code, body)
			}
		})

		t.Run("a title of 200 characters is accepted", func(t *testing.T) {
			p := `{"title":"` + strings.Repeat("A", 200) + `","estate_type_id":1,"address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"}}`
			code, body := do("POST", "/assets/real-estates", mgr, p)
			if code != http.StatusCreated {
				t.Errorf("want 201 got %d: %s", code, body)
			}
		})

		// Without the foreign key this was stored as-is, which is worse than a
		// 500: the registry ends up pointing at a type that does not exist.
		t.Run("an unknown estate type is refused by the database too", func(t *testing.T) {
			var stored int
			err := db.QueryRow(`INSERT INTO ass.real_estate (title, estate_type, status_id)
			    VALUES ('Ghost type', 999999, 3) RETURNING id`).Scan(&stored)
			if err == nil {
				t.Errorf("row %d was stored with a type that does not exist", stored)
			}
		})

		// A registry compares instants. Two rows written from servers on
		// different offsets are not comparable if the zone is dropped.
		t.Run("every date of the table carries its zone", func(t *testing.T) {
			rows, err := db.Query(`
			    SELECT column_name, data_type
			    FROM information_schema.columns
			    WHERE table_schema = 'ass' AND table_name = 'real_estate'
			      AND data_type LIKE 'timestamp%'`)
			if err != nil {
				t.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				var name, kind string
				if err := rows.Scan(&name, &kind); err != nil {
					t.Fatal(err)
				}
				if kind != "timestamp with time zone" {
					t.Errorf("%s is %q, want timestamp with time zone", name, kind)
				}
			}
		})
	})

	// An explicit id does not advance the identity sequence, so seeding rows or
	// restoring a dump leaves the next generated id pointing at a row that
	// already exists. This is what migration 000013 repairs.
	t.Run("a sequence left behind by explicit ids", func(t *testing.T) {
		// Exactly the state a seeded table is in: rows carry ids 1..N while the
		// sequence never moved, so the next generated id is one that exists.
		mustExec(t, `INSERT INTO ass.real_estate (id, title, estate_type, status_id)
		    VALUES (1, 'Seeded with its id', 1, 3)
		    ON CONFLICT (id) DO NOTHING`)
		mustExec(t, `SELECT setval(pg_get_serial_sequence('ass.real_estate', 'id'), 1, false)`)

		payload := `{"title":"After the gap","estate_type_id":1,"address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"}}`

		// The caller sent no id, so this must not come back as a 409 telling him
		// a value of his is already taken.
		code, body := do("POST", "/assets/real-estates", mgr, payload)
		if code != http.StatusInternalServerError {
			t.Errorf("want 500 got %d: %s", code, body)
		}
		if strings.Contains(body, "already used") {
			t.Errorf("a sequence gap must not be blamed on the payload: %s", body)
		}

		migration, err := os.ReadFile("../../../tools/db/migrations/000013_resync_identity_sequences.up.sql")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(string(migration)); err != nil {
			t.Fatalf("migration 000013 failed: %v", err)
		}

		if code, body := do("POST", "/assets/real-estates", mgr, payload); code != http.StatusCreated {
			t.Errorf("want 201 after the repair, got %d: %s", code, body)
		}
	})

	t.Run("draft is invisible to the public", func(t *testing.T) {
		draftID := mustScanID(t, `INSERT INTO ass.real_estate (title, estate_type, status_id) VALUES ('Secret draft', 1, 1) RETURNING id`)

		code, _ := do("GET", "/assets/real-estates/"+itoa(draftID), "", "")
		if code != http.StatusNotFound {
			t.Errorf("anonymous sees a draft: %d", code)
		}
		code, _ = do("GET", "/assets/real-estates/"+itoa(draftID), usr, "")
		if code != http.StatusNotFound {
			t.Errorf("plain user sees a draft: %d", code)
		}
		code, body := do("GET", "/assets/real-estates/"+itoa(draftID), mgr, "")
		if code != http.StatusOK {
			t.Errorf("manager cannot see a draft: %d %s", code, body)
		}

		_, listAnon := do("GET", "/assets/real-estates/active", "", "")
		_, listMgr := do("GET", "/assets/real-estates/active", mgr, "")
		if bytes.Contains([]byte(listAnon), []byte("Secret draft")) {
			t.Errorf("draft leaked in the public listing")
		}
		if !bytes.Contains([]byte(listMgr), []byte("Secret draft")) {
			t.Errorf("draft missing from the manager listing")
		}
		// Matching on the id, not on the title: the database is not wiped
		// between runs, so an identical title from a previous run would make
		// this assertion fail on a perfectly correct listing.
		var listed []map[string]any
		json.Unmarshal([]byte(listMgr), &listed)
		for _, e := range listed {
			if n, ok := e["id"].(float64); ok && int(n) == id {
				t.Errorf("deleted asset still listed: %v", e)
			}
		}
	})

	// A patch merges onto the stored configuration, so an asset that already
	// has a currency keeps it. One that has no configuration at all is really
	// creating it, and must name the currency like a create would.
	t.Run("a patch creating a configuration must name the currency", func(t *testing.T) {
		bareID := mustScanID(t, `
		    INSERT INTO ass.real_estate (title, estate_type, status_id)
		    VALUES ('No configuration yet', 1, 1) RETURNING id`)

		code, body := do("PATCH", "/assets/real-estates/"+itoa(bareID), mgr,
			`{"configuration":{"total_shares":"100","price_per_share":"5","payment_frequency":1,"payment_frequency_type_id":1}}`)
		if code != http.StatusBadRequest {
			t.Fatalf("want 400 got %d: %s", code, body)
		}
		if !strings.Contains(body, "currency_code") {
			t.Errorf("the missing field is not named: %s", body)
		}

		code, body = do("PATCH", "/assets/real-estates/"+itoa(bareID), mgr,
			`{"configuration":{"total_shares":"100","price_per_share":"5","currency_code":"chf",
			  "payment_frequency":1,"payment_frequency_type_id":1}}`)
		if code != http.StatusOK {
			t.Fatalf("want 200 got %d: %s", code, body)
		}
		if !strings.Contains(body, `"currency_code":"CHF"`) {
			t.Errorf("the chosen currency was not kept: %s", body)
		}
	})

	// The lifecycle Antony arbitrated: a stub is born as a draft, it becomes
	// visible only once it carries what an investor needs to decide, and the
	// button that makes it visible says exactly what is missing when it cannot.
	t.Run("publication", func(t *testing.T) {
		stub := `{"title":"Half filled","estate_type_id":1,
		  "address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"}}`

		code, body := do("POST", "/assets/real-estates", mgr, stub)
		if code != http.StatusCreated {
			t.Fatalf("want 201 got %d: %s", code, body)
		}

		var draft map[string]any
		json.Unmarshal([]byte(body), &draft)
		draftID := int(draft["id"].(float64))

		t.Run("an incomplete asset is born a draft", func(t *testing.T) {
			if draft["status"] != "draft" {
				t.Errorf("want draft got %v", draft["status"])
			}
			if draft["active"] != false {
				t.Errorf("an incomplete asset is active: %v", draft["active"])
			}
			if draft["published_at"] != nil {
				t.Errorf("a draft carries a publication date: %v", draft["published_at"])
			}
			if code, _ := do("GET", "/assets/real-estates/"+itoa(draftID), "", ""); code != http.StatusNotFound {
				t.Errorf("a draft is visible to the public: %d", code)
			}
		})

		t.Run("publishing names everything that is missing", func(t *testing.T) {
			code, body := do("POST", "/assets/real-estates/"+itoa(draftID)+"/publish", mgr, "")
			if code != http.StatusConflict {
				t.Fatalf("want 409 got %d: %s", code, body)
			}
			for _, field := range []string{
				"issuer_id", "description", "media",
				"configuration.total_shares", "configuration.price_per_share",
				"configuration.currency_code", "configuration.yield",
				"configuration.payment_frequency", "configuration.payment_frequency_type_id",
			} {
				if !strings.Contains(body, field) {
					t.Errorf("%s is required but not named: %s", field, body)
				}
			}
		})

		t.Run("only a manager can publish", func(t *testing.T) {
			if code, _ := do("POST", "/assets/real-estates/"+itoa(draftID)+"/publish", "", ""); code != http.StatusUnauthorized {
				t.Errorf("anonymous: want 401 got %d", code)
			}
			if code, _ := do("POST", "/assets/real-estates/"+itoa(draftID)+"/publish", usr, ""); code != http.StatusForbidden {
				t.Errorf("plain user: want 403 got %d", code)
			}
		})

		t.Run("completed then published", func(t *testing.T) {
			complete := `{"description":"Now browsable","issuer_id":1,
			  "media":[{"url":"https://x/h.jpg","is_cover":true}],
			  "configuration":{"total_shares":"800","price_per_share":"125","currency_code":"EUR","yield":"3.9",
			    "payment_frequency":4,"payment_frequency_type_id":1}}`

			code, body := do("PATCH", "/assets/real-estates/"+itoa(draftID), mgr, complete)
			if code != http.StatusOK {
				t.Fatalf("patch: want 200 got %d: %s", code, body)
			}

			// Completing does not publish on its own: making an asset visible
			// stays an explicit decision.
			var patched map[string]any
			json.Unmarshal([]byte(body), &patched)
			if patched["status"] != "draft" {
				t.Errorf("a patch published the asset by itself: %v", patched["status"])
			}

			code, body = do("POST", "/assets/real-estates/"+itoa(draftID)+"/publish", mgr, "")
			if code != http.StatusOK {
				t.Fatalf("publish: want 200 got %d: %s", code, body)
			}

			var published map[string]any
			json.Unmarshal([]byte(body), &published)
			if published["status"] != "published" || published["active"] != true {
				t.Errorf("not published: %v", body)
			}
			if published["published_at"] == nil {
				t.Errorf("no publication date: %s", body)
			}

			if code, _ := do("GET", "/assets/real-estates/"+itoa(draftID), "", ""); code != http.StatusOK {
				t.Errorf("a published asset stays hidden from the public: %d", code)
			}

			// Two managers clicking the same button must not produce a failure,
			// and the date of the first publication is the one that counts.
			code, body = do("POST", "/assets/real-estates/"+itoa(draftID)+"/publish", mgr, "")
			if code != http.StatusOK {
				t.Fatalf("republish: want 200 got %d: %s", code, body)
			}
			var again map[string]any
			json.Unmarshal([]byte(body), &again)
			if again["published_at"] != published["published_at"] {
				t.Errorf("publication date moved: %v then %v", published["published_at"], again["published_at"])
			}
		})

		// The correction that matters most: an asset published before these
		// requirements existed must stay editable, otherwise every patch would
		// be refused, including the one completing it.
		t.Run("an asset published incomplete can still be patched", func(t *testing.T) {
			legacyID := mustScanID(t, `
			    INSERT INTO ass.real_estate (title, estate_type, status_id)
			    VALUES ('Legacy listing', 1, 3) RETURNING id`)

			code, body := do("PATCH", "/assets/real-estates/"+itoa(legacyID), mgr,
				`{"title":"Legacy listing renamed"}`)
			if code != http.StatusOK {
				t.Fatalf("want 200 got %d: %s", code, body)
			}
		})

		t.Run("a deleted asset cannot be published", func(t *testing.T) {
			goneID := mustScanID(t, `
			    INSERT INTO ass.real_estate (title, estate_type, status_id, deleted_at)
			    VALUES ('Gone', 1, 7, now()) RETURNING id`)

			code, _ := do("POST", "/assets/real-estates/"+itoa(goneID)+"/publish", mgr, "")
			if code != http.StatusNotFound {
				t.Errorf("want 404 got %d", code)
			}
			code, _ = do("POST", "/assets/real-estates/999999/publish", mgr, "")
			if code != http.StatusNotFound {
				t.Errorf("unknown id: want 404 got %d", code)
			}
		})
	})
}

func itoa(i int) string {
	b, _ := json.Marshal(i)
	return string(b)
}
