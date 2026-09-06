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

	payload := `{
      "title": "Villa Belair",
      "description": "Nice",
      "estate_type_id": 1,
      "address": {"street":"12 rue de la Gare","postal_code":"L-1611","city":"Luxembourg","country_code":"lu","latitude":"49.611622","longitude":"6.131935"},
      "specification": {"surface_area":"128.50","bedroom_number":4,"energy_class":"c"},
      "configuration": {"total_shares":"1500","price_per_share":"199.99","currency_code":"eur","entry_fee_rate":"2.5"},
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
			"fractional shares":   `{"title":"x","address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"},"configuration":{"total_shares":"10.5","price_per_share":"1"}}`,
			"zero price":          `{"title":"x","address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"},"configuration":{"total_shares":"10","price_per_share":"0"}}`,
			"bad country":         `{"title":"x","address":{"street":"a","postal_code":"b","city":"c","country_code":"LUX"}}`,
			"half coordinates":    `{"title":"x","address":{"street":"a","postal_code":"b","city":"c","country_code":"LU","latitude":"49.1"}}`,
			"two covers":          `{"title":"x","address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"},"media":[{"url":"a","is_cover":true},{"url":"b","is_cover":true}]}`,
			"unknown estate type": `{"title":"x","estate_type_id":999,"address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"}}`,
			"empty title":         `{"title":"   ","address":{"street":"a","postal_code":"b","city":"c","country_code":"LU"}}`,
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
		code, body := do("PATCH", "/assets/real-estates/"+itoa(id), mgr, `{"description":""}`)
		if code != 200 {
			t.Fatalf("want 200 got %d: %s", code, body)
		}

		var re map[string]any
		json.Unmarshal([]byte(body), &re)
		if re["description"] != nil {
			t.Errorf("description not cleared: %v", re["description"])
		}
		if re["title"] != "Villa Belair v2" {
			t.Errorf("title lost: %v", re["title"])
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

	t.Run("guards on a tokenized asset", func(t *testing.T) {
		var tokenID, assetID int
		db.QueryRow(`INSERT INTO blk.token (address, token_name, symbol, nb_decimal) VALUES ('0xguard','G','G',0) RETURNING id`).Scan(&tokenID)
		db.QueryRow(`INSERT INTO ass.real_estate (title, status_id, token_id) VALUES ('Tokenized guard', 3, $1) RETURNING id`, tokenID).Scan(&assetID)
		db.Exec(`INSERT INTO ass.real_estate_shares_config (real_estate_id, total_shares, price_per_share) VALUES ($1, 1000, 10)`, assetID)

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
		var assetID int
		db.QueryRow(`INSERT INTO ass.real_estate (title, status_id) VALUES ('Reserved guard', 3) RETURNING id`).Scan(&assetID)
		db.Exec(`INSERT INTO ass.real_estate_shares_config (real_estate_id, total_shares, price_per_share) VALUES ($1, 1000, 10)`, assetID)
		// status 2 = RESERVED, which carries counts_as_reserved.
		db.Exec(`INSERT INTO iss.issuance_orders (asset_id, quantity, status_id, order_ref) VALUES ($1, 400, 2, 'ORD-RESERVED-1')`, assetID)

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

	t.Run("draft is invisible to the public", func(t *testing.T) {
		var draftID int
		db.QueryRow(`INSERT INTO ass.real_estate (title, status_id) VALUES ('Secret draft', 1) RETURNING id`).Scan(&draftID)

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
		if bytes.Contains([]byte(listMgr), []byte("Villa Belair")) {
			t.Errorf("deleted asset still listed")
		}
	})
}

func itoa(i int) string {
	b, _ := json.Marshal(i)
	return string(b)
}
