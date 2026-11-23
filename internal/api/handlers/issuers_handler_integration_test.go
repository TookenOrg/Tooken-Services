//go:build integration

// HTTP coverage of the issuer back office: /assets/issuers and
// /assets/issuers/{id}, through the real router, the spec-driven auth
// middleware, the handler, the service and PostgreSQL.
//
// The database is shared with other suites and is never wiped: every issuer
// written here carries a per-run name prefix, every lookup filters on it, and
// the rows (with the assets and tokens hung on them) are removed when the test
// ends.
package handlers

import (
	"crypto/md5"
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
	_ "github.com/lib/pq"
)

// ---------------------------------------------------------------- shared helpers

// callAPI sends one request through the router of the harness. A body, when
// given, is sent as JSON.
func callAPI(h *orderHarness, method, path, token, body string) (int, string) {
	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(method, globals.BaseURL+path, nil)
	} else {
		req = httptest.NewRequest(method, globals.BaseURL+path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	h.r.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

// expectAPIMessage checks that a refusal comes in the APIResponse envelope and
// that its message names the reason, so a client can tell two refusals sharing
// a status code apart.
func expectAPIMessage(t *testing.T, body string, fragments ...string) {
	t.Helper()
	expectMessage(t, body)
	var resp server.APIResponse
	_ = json.Unmarshal([]byte(body), &resp)
	for _, fragment := range fragments {
		if !strings.Contains(resp.Message, fragment) {
			t.Fatalf("message %q does not mention %q", resp.Message, fragment)
		}
	}
}

// expectGenericFailure checks that a 500 says the operation failed and
// nothing more. Database and driver errors read "sql: ..." or "pq: ...": a
// message holding a colon is one of them leaking to the client.
func expectGenericFailure(t *testing.T, body string) {
	t.Helper()
	expectAPIMessage(t, body, "Unable to")
	var resp server.APIResponse
	_ = json.Unmarshal([]byte(body), &resp)
	if strings.Contains(resp.Message, ":") {
		t.Fatalf("the 500 leaks an internal error: %q", resp.Message)
	}
}

// withClosedDatabase runs fn while globals.DB is a pool that has already been
// closed, so the first query a handler reaches fails as it would on a lost
// database. Tests of this package never run in parallel, which is what makes
// swapping the global safe; it is put back before returning.
func withClosedDatabase(t *testing.T, fn func()) {
	t.Helper()
	broken, err := sql.Open("postgres", os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	broken.Close()

	live := globals.DB
	globals.DB = broken
	defer func() { globals.DB = live }()
	fn()
}

// newTrackedUser writes a user with the order fixtures, then removes him when
// the test ends, together with every KYC filing he holds or decided. His
// favorites go with him through their ON DELETE CASCADE.
func newTrackedUser(t *testing.T, h *orderHarness, role string) orderUser {
	t.Helper()
	u := h.newUser(t, role, "none", "NULL")
	t.Cleanup(func() {
		if _, err := h.db.Exec(`
		    DELETE FROM usr.kyc_verification
		    WHERE user_id = $1 OR decided_by = $1 OR revoked_by = $1`, u.id); err != nil {
			t.Logf("cleanup of the KYC filings of user %d: %v", u.id, err)
		}
		if _, err := h.db.Exec(`DELETE FROM usr.users WHERE id = $1`, u.id); err != nil {
			t.Logf("cleanup of user %d: %v", u.id, err)
		}
	})
	return u
}

// assetFixtures writes issuers and assets straight in SQL, to set up the state
// a request is checked against without going through the rules under test.
//
// Every issuer it writes — and every issuer a test creates through the API —
// is named after prefix, which is what the cleanup keys on: the assets under
// those issuers go first (the foreign key is ON DELETE RESTRICT), then the
// tokens of the published ones, then the issuers.
type assetFixtures struct {
	h      *orderHarness
	prefix string
}

func newAssetFixtures(t *testing.T, h *orderHarness) *assetFixtures {
	t.Helper()
	f := &assetFixtures{h: h, prefix: h.unique("Issuer IT")}
	t.Cleanup(func() {
		pattern := f.prefix + "%"
		for _, query := range []string{
			`DELETE FROM ass.real_estate WHERE issuer_id IN (SELECT id FROM ass.issuer WHERE name LIKE $1)`,
			`DELETE FROM blk.token WHERE token_name LIKE $1`,
			`DELETE FROM ass.issuer WHERE name LIKE $1`,
		} {
			if _, err := h.db.Exec(query, pattern); err != nil {
				t.Logf("cleanup %q: %v", query, err)
			}
		}
	})
	return f
}

// name returns an issuer name the cleanup will find.
func (f *assetFixtures) name(label string) string {
	return f.prefix + " " + label
}

// registration returns a registration number no other issuer holds.
func (f *assetFixtures) registration() string {
	return f.h.unique("B")
}

// lei returns a LEI no other issuer holds: twenty upper-case hexadecimal
// characters, which the LEI format accepts.
func (f *assetFixtures) lei() string {
	sum := md5.Sum([]byte(f.h.unique("lei")))
	return strings.ToUpper(fmt.Sprintf("%x", sum))[:20]
}

// issuer writes a Luxembourg issuer in the given status. An empty registration
// is stored as NULL.
func (f *assetFixtures) issuer(t *testing.T, label string, status int, registration string) int {
	t.Helper()
	var reg sql.NullString
	if registration != "" {
		reg = sql.NullString{String: registration, Valid: true}
	}
	var id int
	if err := f.h.db.QueryRow(`
	    INSERT INTO ass.issuer (name, legal_form, registration_number, country_code, status_id)
	    VALUES ($1, 'SA', $2, 'LU', $3)
	    RETURNING id`, f.name(label), reg, status).Scan(&id); err != nil {
		t.Fatalf("seeding issuer %q: %v", label, err)
	}
	return id
}

// activeIssuer is the usual starting point: an issuer assets can be published
// under, identifiable by its registration number.
func (f *assetFixtures) activeIssuer(t *testing.T, label string) int {
	t.Helper()
	return f.issuer(t, label, 2, f.registration())
}

// asset writes a real estate under the issuer, published or draft. A published
// asset needs a token: real_estate_active_requires_token_ck refuses a visible
// row without one.
func (f *assetFixtures) asset(t *testing.T, issuer int, published bool) int {
	t.Helper()
	status := 1
	var token sql.NullInt64
	if published {
		status = 3
		var id int64
		if err := f.h.db.QueryRow(`
		    INSERT INTO blk.token (symbol, token_name, address, nb_decimal)
		    VALUES ('ISS', $1::text, '0x' || md5($1::text) || '00000000', 0)
		    RETURNING id`, f.h.unique(f.prefix+" token")).Scan(&id); err != nil {
			t.Fatalf("seeding token: %v", err)
		}
		token = sql.NullInt64{Int64: id, Valid: true}
	}

	var id int
	if err := f.h.db.QueryRow(`
	    INSERT INTO ass.real_estate (title, estate_type, issuer_id, status_id, token_id)
	    SELECT $1, min(id), $2, $3, $4 FROM ass.real_estate_type
	    RETURNING id`, f.h.unique("Villa Belair"), issuer, status, token).Scan(&id); err != nil {
		t.Fatalf("seeding asset: %v", err)
	}
	return id
}

// deletedAsset writes a soft-deleted real estate under the issuer: it still
// references the issuer, but nothing reads it any more.
func (f *assetFixtures) deletedAsset(t *testing.T, issuer int) int {
	t.Helper()
	var id int
	if err := f.h.db.QueryRow(`
	    INSERT INTO ass.real_estate (title, estate_type, issuer_id, status_id, deleted_at)
	    SELECT $1, min(id), $2, 7, now() FROM ass.real_estate_type
	    RETURNING id`, f.h.unique("Villa Belair (deleted)"), issuer).Scan(&id); err != nil {
		t.Fatalf("seeding deleted asset: %v", err)
	}
	return id
}

// storedIssuer is an issuer row as the database holds it, read without the API.
type storedIssuer struct {
	name         string
	legalForm    string
	registration sql.NullString
	country      string
	lei          sql.NullString
	status       int
}

func (f *assetFixtures) stored(t *testing.T, id int) storedIssuer {
	t.Helper()
	var s storedIssuer
	if err := f.h.db.QueryRow(`
	    SELECT name, legal_form, registration_number, country_code, lei_code, status_id
	    FROM ass.issuer WHERE id = $1`, id).Scan(
		&s.name, &s.legalForm, &s.registration, &s.country, &s.lei, &s.status); err != nil {
		t.Fatalf("reading issuer %d: %v", id, err)
	}
	return s
}

// countNamed counts the issuers stored under an exact name: it is how a test
// proves that a refused creation wrote nothing.
func (f *assetFixtures) countNamed(t *testing.T, name string) int {
	t.Helper()
	return f.h.count(t, `SELECT count(*) FROM ass.issuer WHERE name = $1`, name)
}

func decodeIssuer(t *testing.T, body string) server.Issuer {
	t.Helper()
	var issuer server.Issuer
	if err := json.Unmarshal([]byte(body), &issuer); err != nil {
		t.Fatalf("decoding %s: %v", body, err)
	}
	if issuer.Id == 0 {
		t.Fatalf("no issuer in %s", body)
	}
	return issuer
}

func issuerPath(id int) string {
	return fmt.Sprintf("/assets/issuers/%d", id)
}

// unknownID is an identifier no issuer or real estate can hold: the columns
// are SERIAL, and the largest INTEGER is far beyond anything a test database
// reaches.
const unknownID = 2147483647

// ---------------------------------------------------------------- cases

func TestIssuerRoutesAreReservedToManagersAndAdmins(t *testing.T) {
	h := newOrderHarness(t)
	f := newAssetFixtures(t, h)
	user := newTrackedUser(t, h, authUtils.RoleUser)
	manager := newTrackedUser(t, h, authUtils.RoleManager)
	admin := newTrackedUser(t, h, authUtils.RoleAdmin)

	target := f.activeIssuer(t, "Tooken Villa Belair SA")
	before := f.stored(t, target)
	createdName := f.name("never created")
	createBody := string(mustJSON(t, map[string]any{
		"name": createdName, "legal_form": "SA", "country_code": "LU",
		"registration_number": f.registration(),
	}))
	patchBody := string(mustJSON(t, map[string]any{"name": f.name("never renamed")}))

	routes := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{"list", http.MethodGet, "/assets/issuers", ""},
		{"detail", http.MethodGet, issuerPath(target), ""},
		{"create", http.MethodPost, "/assets/issuers", createBody},
		{"patch", http.MethodPatch, issuerPath(target), patchBody},
		{"dissolve", http.MethodDelete, issuerPath(target), ""},
	}

	for _, route := range routes {
		t.Run("an anonymous caller is refused with 401 on "+route.name, func(t *testing.T) {
			code, body := callAPI(h, route.method, route.path, "", route.body)
			expectCode(t, code, body, http.StatusUnauthorized)
			expectAPIMessage(t, body, "Authorization header required")
		})
		t.Run("a USER is refused with 403 on "+route.name, func(t *testing.T) {
			code, body := callAPI(h, route.method, route.path, user.token, route.body)
			expectCode(t, code, body, http.StatusForbidden)
			expectAPIMessage(t, body, "not allowed")
		})
	}

	t.Run("no refused write reached the database", func(t *testing.T) {
		if n := f.countNamed(t, createdName); n != 0 {
			t.Fatalf("%d issuer(s) created by a refused request", n)
		}
		if after := f.stored(t, target); after != before {
			t.Fatalf("issuer changed by refused requests: before %+v, after %+v", before, after)
		}
	})

	t.Run("a MANAGER and an ADMIN are both allowed", func(t *testing.T) {
		for role, token := range map[string]string{"MANAGER": manager.token, "ADMIN": admin.token} {
			code, body := callAPI(h, http.MethodGet, issuerPath(target), token, "")
			if code != http.StatusOK {
				t.Fatalf("%s: want 200, got %d: %s", role, code, body)
			}
			code, body = callAPI(h, http.MethodGet, "/assets/issuers", token, "")
			if code != http.StatusOK {
				t.Fatalf("%s list: want 200, got %d: %s", role, code, body)
			}
		}
	})
}

func TestReadIssuers(t *testing.T) {
	h := newOrderHarness(t)
	f := newAssetFixtures(t, h)
	manager := newTrackedUser(t, h, authUtils.RoleManager)

	t.Run("the list includes dissolved issuers and counts only live assets", func(t *testing.T) {
		carrying := f.activeIssuer(t, "Tooken RE I SA")
		f.asset(t, carrying, false)
		f.asset(t, carrying, true)
		f.deletedAsset(t, carrying)
		dissolved := f.issuer(t, "Tooken RE Liquidated SA", 4, f.registration())

		code, body := callAPI(h, http.MethodGet, "/assets/issuers", manager.token, "")
		expectCode(t, code, body, http.StatusOK)
		var listed []server.Issuer
		if err := json.Unmarshal([]byte(body), &listed); err != nil {
			t.Fatalf("decoding the list: %v", err)
		}

		byID := map[int]server.Issuer{}
		for _, issuer := range listed {
			byID[issuer.Id] = issuer
		}
		got, ok := byID[carrying]
		if !ok {
			t.Fatalf("issuer %d missing from the list", carrying)
		}
		if got.RealEstateCount == nil || *got.RealEstateCount != 2 {
			t.Fatalf("real_estate_count = %v, want 2 (the soft-deleted asset does not count)", got.RealEstateCount)
		}
		gone, ok := byID[dissolved]
		if !ok {
			t.Fatalf("dissolved issuer %d missing from the list", dissolved)
		}
		if gone.Status != "dissolved" || gone.StatusId != 4 {
			t.Fatalf("dissolved issuer listed as %q (%d)", gone.Status, gone.StatusId)
		}
	})

	t.Run("the detail returns the issuer as stored", func(t *testing.T) {
		id := f.activeIssuer(t, "Tooken Villa Belair SA")
		f.asset(t, id, false)
		stored := f.stored(t, id)

		code, body := callAPI(h, http.MethodGet, issuerPath(id), manager.token, "")
		expectCode(t, code, body, http.StatusOK)
		got := decodeIssuer(t, body)
		if got.Id != id || got.Name != stored.name || got.LegalForm != "SA" || got.CountryCode != "LU" {
			t.Fatalf("detail = %+v, want the stored %+v", got, stored)
		}
		if got.RegistrationNumber == nil || *got.RegistrationNumber != stored.registration.String {
			t.Fatalf("registration_number = %v, want %q", got.RegistrationNumber, stored.registration.String)
		}
		if got.StatusId != 2 || got.Status != "active" {
			t.Fatalf("status = %q (%d), want active (2)", got.Status, got.StatusId)
		}
		if got.RealEstateCount == nil || *got.RealEstateCount != 1 {
			t.Fatalf("real_estate_count = %v, want 1", got.RealEstateCount)
		}
		if got.CreatedAt == nil || got.UpdatedAt == nil {
			t.Fatalf("timestamps missing: %s", body)
		}
	})

	t.Run("an unknown issuer answers 404", func(t *testing.T) {
		code, body := callAPI(h, http.MethodGet, issuerPath(unknownID), manager.token, "")
		expectCode(t, code, body, http.StatusNotFound)
		expectAPIMessage(t, body, "Issuer not found")
	})
}

func TestCreateIssuer(t *testing.T) {
	h := newOrderHarness(t)
	f := newAssetFixtures(t, h)
	manager := newTrackedUser(t, h, authUtils.RoleManager)

	create := func(t *testing.T, payload map[string]any) (int, string) {
		t.Helper()
		return callAPI(h, http.MethodPost, "/assets/issuers", manager.token, string(mustJSON(t, payload)))
	}

	t.Run("an issuer is created active by default and stored normalised", func(t *testing.T) {
		name := f.name("Tooken Villa Belair SA")
		registration := f.registration()
		lei := f.lei()

		code, body := create(t, map[string]any{
			"name":                "  " + name + "  ",
			"legal_form":          " SA ",
			"country_code":        " lu ",
			"registration_number": "  " + registration + "  ",
			"lei_code":            strings.ToLower(lei),
		})
		expectCode(t, code, body, http.StatusCreated)
		got := decodeIssuer(t, body)
		if got.StatusId != 2 || got.Status != "active" {
			t.Fatalf("status = %q (%d), want active (2)", got.Status, got.StatusId)
		}
		if got.RealEstateCount == nil || *got.RealEstateCount != 0 {
			t.Fatalf("real_estate_count = %v, want 0", got.RealEstateCount)
		}

		stored := f.stored(t, got.Id)
		want := storedIssuer{
			name: name, legalForm: "SA", country: "LU", status: 2,
			registration: sql.NullString{String: registration, Valid: true},
			lei:          sql.NullString{String: lei, Valid: true},
		}
		if stored != want {
			t.Fatalf("stored %+v, want %+v", stored, want)
		}
		if got.Name != stored.name || got.CountryCode != stored.country || got.LeiCode == nil || *got.LeiCode != stored.lei.String {
			t.Fatalf("answer %+v differs from the stored row %+v", got, stored)
		}
	})

	t.Run("a draft can be created without a registration number", func(t *testing.T) {
		name := f.name("Tooken RE II SCSp (incorporating)")
		code, body := create(t, map[string]any{
			"name": name, "legal_form": "SCSp", "country_code": "LU", "status_id": 1,
		})
		expectCode(t, code, body, http.StatusCreated)
		got := decodeIssuer(t, body)
		if got.Status != "draft" || got.RegistrationNumber != nil {
			t.Fatalf("draft created as %q with registration %v", got.Status, got.RegistrationNumber)
		}
		if stored := f.stored(t, got.Id); stored.status != 1 || stored.registration.Valid {
			t.Fatalf("stored %+v, want a draft without registration number", stored)
		}
	})

	// Lengths are counted in characters, not bytes: a name written with accents
	// is not longer than the same name without.
	t.Run("the name limit counts characters, not bytes", func(t *testing.T) {
		pad := func(total int) string {
			base := f.name("")
			return base + strings.Repeat("é", total-len([]rune(base)))
		}

		atLimit := pad(255)
		code, body := create(t, map[string]any{
			"name": atLimit, "legal_form": "SA", "country_code": "LU", "status_id": 1,
		})
		expectCode(t, code, body, http.StatusCreated)

		overLimit := pad(256)
		code, body = create(t, map[string]any{
			"name": overLimit, "legal_form": "SA", "country_code": "LU", "status_id": 1,
		})
		expectCode(t, code, body, http.StatusBadRequest)
		expectAPIMessage(t, body, "name must be at most 255 characters, got 256")
		if n := f.countNamed(t, overLimit); n != 0 {
			t.Fatalf("%d issuer(s) written with a refused name", n)
		}
	})

	t.Run("an invalid payload is refused with 400 and writes nothing", func(t *testing.T) {
		cases := []struct {
			name    string
			payload func(name string) map[string]any
			reason  []string
		}{
			{
				"an active issuer without registration number",
				func(n string) map[string]any {
					return map[string]any{"name": n, "legal_form": "SA", "country_code": "LU"}
				},
				[]string{"an active issuer requires registration_number", "draft (status_id 1)"},
			},
			{
				"a blank registration number counts as missing",
				func(n string) map[string]any {
					return map[string]any{"name": n, "legal_form": "SA", "country_code": "LU", "registration_number": "   "}
				},
				[]string{"an active issuer requires registration_number"},
			},
			{
				"an issuer created dissolved",
				func(n string) map[string]any {
					return map[string]any{"name": n, "legal_form": "SA", "country_code": "LU", "registration_number": f.registration(), "status_id": 4}
				},
				[]string{"cannot be created dissolved"},
			},
			{
				"an unknown status",
				func(n string) map[string]any {
					return map[string]any{"name": n, "legal_form": "SA", "country_code": "LU", "status_id": 9}
				},
				[]string{"status_id 9 does not refer to an existing status"},
			},
			{
				"a legal form left blank",
				func(n string) map[string]any {
					return map[string]any{"name": n, "legal_form": "  ", "country_code": "LU", "status_id": 1}
				},
				[]string{"legal_form is required"},
			},
			{
				"a registration number longer than 64 characters",
				func(n string) map[string]any {
					return map[string]any{"name": n, "legal_form": "SA", "country_code": "LU", "registration_number": strings.Repeat("B", 65)}
				},
				[]string{"registration_number must be at most 64 characters, got 65"},
			},
			{
				"a three-letter country code",
				func(n string) map[string]any {
					return map[string]any{"name": n, "legal_form": "SA", "country_code": "LUX", "status_id": 1}
				},
				[]string{"country_code must be an ISO 3166-1 alpha-2 code", `"LUX"`},
			},
			{
				"a LEI shorter than 20 characters",
				func(n string) map[string]any {
					return map[string]any{"name": n, "legal_form": "SA", "country_code": "LU", "status_id": 1, "lei_code": "529900T8BM49AURSDO5"}
				},
				[]string{"lei_code must be 20 upper-case alphanumerics"},
			},
			{
				"a LEI with punctuation",
				func(n string) map[string]any {
					return map[string]any{"name": n, "legal_form": "SA", "country_code": "LU", "status_id": 1, "lei_code": "529900T8BM49-URSDO55"}
				},
				[]string{"lei_code must be 20 upper-case alphanumerics"},
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				name := f.name(c.name)
				code, body := create(t, c.payload(name))
				expectCode(t, code, body, http.StatusBadRequest)
				expectAPIMessage(t, body, c.reason...)
				if n := f.countNamed(t, name); n != 0 {
					t.Fatalf("%d issuer(s) written by a refused creation", n)
				}
			})
		}

		t.Run("a blank name", func(t *testing.T) {
			code, body := create(t, map[string]any{"name": "   ", "legal_form": "SA", "country_code": "LU", "status_id": 1})
			expectCode(t, code, body, http.StatusBadRequest)
			expectAPIMessage(t, body, "name is required")
		})

		t.Run("a body that is not JSON", func(t *testing.T) {
			code, body := callAPI(h, http.MethodPost, "/assets/issuers", manager.token, `{"name": "`+f.name("truncated"))
			expectCode(t, code, body, http.StatusBadRequest)
			expectAPIMessage(t, body, "Invalid payload")
			if n := f.countNamed(t, f.name("truncated")); n != 0 {
				t.Fatalf("%d issuer(s) written from a truncated body", n)
			}
		})
	})

	t.Run("a registration number already used in the same country is refused with 409", func(t *testing.T) {
		registration := f.registration()
		f.issuer(t, "Tooken RE I SA", 2, registration)

		name := f.name("Tooken RE I bis SA")
		code, body := create(t, map[string]any{
			"name": name, "legal_form": "SA", "country_code": "LU", "registration_number": registration,
		})
		expectCode(t, code, body, http.StatusConflict)
		expectAPIMessage(t, body, "registration number is already used by another issuer in this country")
		if n := f.countNamed(t, name); n != 0 {
			t.Fatalf("%d duplicate issuer(s) written", n)
		}
	})

	t.Run("the same registration number is accepted in another country", func(t *testing.T) {
		registration := f.registration()
		f.issuer(t, "Tooken RE I SA (LU)", 2, registration)

		code, body := create(t, map[string]any{
			"name": f.name("Tooken RE I SAS (FR)"), "legal_form": "SAS", "country_code": "FR", "registration_number": registration,
		})
		expectCode(t, code, body, http.StatusCreated)
		if got := decodeIssuer(t, body); got.CountryCode != "FR" {
			t.Fatalf("country_code = %q, want FR", got.CountryCode)
		}
	})

	t.Run("a LEI already held by another issuer is refused with 409, whatever its case", func(t *testing.T) {
		lei := f.lei()
		code, body := create(t, map[string]any{
			"name": f.name("Tooken RE III SA"), "legal_form": "SA", "country_code": "LU",
			"registration_number": f.registration(), "lei_code": lei,
		})
		expectCode(t, code, body, http.StatusCreated)

		name := f.name("Tooken RE III bis SA")
		code, body = create(t, map[string]any{
			"name": name, "legal_form": "SA", "country_code": "LU",
			"registration_number": f.registration(), "lei_code": strings.ToLower(lei),
		})
		expectCode(t, code, body, http.StatusConflict)
		expectAPIMessage(t, body, "LEI code already identifies another issuer")
		if n := f.countNamed(t, name); n != 0 {
			t.Fatalf("%d issuer(s) written with a duplicate LEI", n)
		}
	})
}

func TestPatchIssuer(t *testing.T) {
	h := newOrderHarness(t)
	f := newAssetFixtures(t, h)
	admin := newTrackedUser(t, h, authUtils.RoleAdmin)

	patch := func(t *testing.T, id int, payload map[string]any) (int, string) {
		t.Helper()
		return callAPI(h, http.MethodPatch, issuerPath(id), admin.token, string(mustJSON(t, payload)))
	}

	// expectRefused checks the code and the reason, then that the row is
	// exactly as it was before the request.
	expectRefused := func(t *testing.T, id int, before storedIssuer, code int, body string, want int, reason ...string) {
		t.Helper()
		expectCode(t, code, body, want)
		expectAPIMessage(t, body, reason...)
		if after := f.stored(t, id); after != before {
			t.Fatalf("refused patch changed the issuer: before %+v, after %+v", before, after)
		}
	}

	t.Run("an omitted field keeps its value and an empty optional field is cleared", func(t *testing.T) {
		id := f.activeIssuer(t, "Tooken Villa Belair SA")
		lei := f.lei()
		if _, err := h.db.Exec(`UPDATE ass.issuer SET lei_code = $2 WHERE id = $1`, id, lei); err != nil {
			t.Fatal(err)
		}
		before := f.stored(t, id)

		renamed := f.name("Tooken Villa Belair SARL")
		code, body := patch(t, id, map[string]any{"name": renamed, "legal_form": "SARL", "lei_code": ""})
		expectCode(t, code, body, http.StatusOK)
		got := decodeIssuer(t, body)
		if got.Name != renamed || got.LegalForm != "SARL" || got.LeiCode != nil {
			t.Fatalf("answer %+v, want the new name and legal form and no LEI", got)
		}

		after := f.stored(t, id)
		want := before
		want.name, want.legalForm, want.lei = renamed, "SARL", sql.NullString{}
		if after != want {
			t.Fatalf("stored %+v, want %+v", after, want)
		}
	})

	t.Run("a draft is activated when its registration number comes in the same request", func(t *testing.T) {
		id := f.issuer(t, "Tooken RE II SCSp", 1, "")
		registration := f.registration()

		code, body := patch(t, id, map[string]any{"status_id": 2, "registration_number": registration})
		expectCode(t, code, body, http.StatusOK)
		if got := decodeIssuer(t, body); got.Status != "active" {
			t.Fatalf("status = %q, want active", got.Status)
		}
		if stored := f.stored(t, id); stored.status != 2 || stored.registration.String != registration {
			t.Fatalf("stored %+v, want active with %q", stored, registration)
		}
	})

	t.Run("a draft without registration number cannot be activated", func(t *testing.T) {
		id := f.issuer(t, "Tooken RE II SCSp", 1, "")
		before := f.stored(t, id)

		code, body := patch(t, id, map[string]any{"status_id": 2})
		expectRefused(t, id, before, code, body, http.StatusConflict,
			"cannot be activated while it lacks registration_number")
	})

	t.Run("an active issuer cannot lose its registration number", func(t *testing.T) {
		id := f.activeIssuer(t, "Tooken Villa Belair SA")
		before := f.stored(t, id)

		code, body := patch(t, id, map[string]any{"registration_number": ""})
		expectRefused(t, id, before, code, body, http.StatusConflict,
			"cannot lose what identifies it", "registration_number")
	})

	// Draft counts as a withdrawal like the two other statuses: an asset can
	// only stay published under an ACTIVE issuer.
	t.Run("an issuer carrying a published asset cannot be withdrawn", func(t *testing.T) {
		for _, target := range []struct {
			status int
			label  string
		}{{1, "draft"}, {3, "suspended"}, {4, "dissolved"}} {
			t.Run("to "+target.label, func(t *testing.T) {
				id := f.activeIssuer(t, "Tooken Villa Belair SA")
				f.asset(t, id, true)
				before := f.stored(t, id)

				code, body := patch(t, id, map[string]any{"status_id": target.status})
				expectRefused(t, id, before, code, body, http.StatusConflict,
					"carries 1 published asset(s)", "before setting it to "+target.label)
			})
		}
	})

	t.Run("a patch that keeps the issuer active is not blocked by its published assets", func(t *testing.T) {
		id := f.activeIssuer(t, "Tooken Villa Belair SA")
		f.asset(t, id, true)

		renamed := f.name("Tooken Villa Belair SA (renamed)")
		code, body := patch(t, id, map[string]any{"name": renamed})
		expectCode(t, code, body, http.StatusOK)
		if stored := f.stored(t, id); stored.name != renamed || stored.status != 2 {
			t.Fatalf("stored %+v, want renamed and still active", stored)
		}
	})

	t.Run("an issuer whose assets are all drafts can be suspended", func(t *testing.T) {
		id := f.activeIssuer(t, "Tooken Villa Belair SA")
		f.asset(t, id, false)

		code, body := patch(t, id, map[string]any{"status_id": 3})
		expectCode(t, code, body, http.StatusOK)
		if got := decodeIssuer(t, body); got.Status != "suspended" {
			t.Fatalf("status = %q, want suspended", got.Status)
		}
	})

	t.Run("an issuer still carrying a draft asset cannot be dissolved by a patch", func(t *testing.T) {
		id := f.activeIssuer(t, "Tooken Villa Belair SA")
		f.asset(t, id, false)
		before := f.stored(t, id)

		code, body := patch(t, id, map[string]any{"status_id": 4})
		expectRefused(t, id, before, code, body, http.StatusConflict, "still carries 1 asset(s)")
	})

	t.Run("a dissolved issuer is final", func(t *testing.T) {
		id := f.issuer(t, "Tooken RE Liquidated SA", 4, f.registration())
		before := f.stored(t, id)

		code, body := patch(t, id, map[string]any{"name": f.name("Tooken RE Revived SA")})
		expectRefused(t, id, before, code, body, http.StatusConflict, "dissolved issuer is final")

		code, body = patch(t, id, map[string]any{"status_id": 2})
		expectRefused(t, id, before, code, body, http.StatusConflict, "dissolved issuer is final")
	})

	t.Run("a registration number taken by another issuer of the same country is refused", func(t *testing.T) {
		taken := f.registration()
		f.issuer(t, "Tooken RE I SA", 2, taken)
		id := f.activeIssuer(t, "Tooken RE IV SA")
		before := f.stored(t, id)

		code, body := patch(t, id, map[string]any{"registration_number": taken})
		expectRefused(t, id, before, code, body, http.StatusConflict,
			"registration number is already used by another issuer in this country")
	})

	t.Run("an invalid merged issuer is refused with 400", func(t *testing.T) {
		id := f.activeIssuer(t, "Tooken Villa Belair SA")
		before := f.stored(t, id)

		for _, c := range []struct {
			name    string
			payload map[string]any
			reason  string
		}{
			{"an empty name", map[string]any{"name": ""}, "name is required"},
			{"a one-letter country code", map[string]any{"country_code": "L"}, "country_code must be an ISO 3166-1 alpha-2 code"},
			{"an unknown status", map[string]any{"status_id": 7}, "status_id 7 does not refer to an existing status"},
			{"a malformed LEI", map[string]any{"lei_code": "NOT-A-LEI"}, "lei_code must be 20 upper-case alphanumerics"},
		} {
			t.Run(c.name, func(t *testing.T) {
				code, body := patch(t, id, c.payload)
				expectRefused(t, id, before, code, body, http.StatusBadRequest, c.reason)
			})
		}

		t.Run("a body that does not match the schema", func(t *testing.T) {
			code, body := callAPI(h, http.MethodPatch, issuerPath(id), admin.token, `{"name": 42}`)
			expectRefused(t, id, before, code, body, http.StatusBadRequest, "Invalid payload")
		})
	})

	t.Run("an unknown issuer answers 404", func(t *testing.T) {
		code, body := patch(t, unknownID, map[string]any{"name": f.name("ghost")})
		expectCode(t, code, body, http.StatusNotFound)
		expectAPIMessage(t, body, "Issuer not found")
	})
}

func TestDissolveIssuer(t *testing.T) {
	h := newOrderHarness(t)
	f := newAssetFixtures(t, h)
	manager := newTrackedUser(t, h, authUtils.RoleManager)

	dissolve := func(id int) (int, string) {
		return callAPI(h, http.MethodDelete, issuerPath(id), manager.token, "")
	}

	t.Run("dissolving keeps the row, and doing it twice still answers 204", func(t *testing.T) {
		id := f.activeIssuer(t, "Tooken RE Wound Up SA")

		code, body := dissolve(id)
		expectCode(t, code, body, http.StatusNoContent)
		if body != "" {
			t.Fatalf("204 with a body: %s", body)
		}
		if stored := f.stored(t, id); stored.status != 4 {
			t.Fatalf("status_id = %d, want 4 (dissolved)", stored.status)
		}

		code, body = dissolve(id)
		expectCode(t, code, body, http.StatusNoContent)

		code, body = callAPI(h, http.MethodGet, issuerPath(id), manager.token, "")
		expectCode(t, code, body, http.StatusOK)
		if got := decodeIssuer(t, body); got.Status != "dissolved" {
			t.Fatalf("status = %q, want dissolved", got.Status)
		}
	})

	t.Run("an issuer still carrying an asset cannot be dissolved", func(t *testing.T) {
		id := f.activeIssuer(t, "Tooken Villa Belair SA")
		f.asset(t, id, false)
		before := f.stored(t, id)

		code, body := dissolve(id)
		expectCode(t, code, body, http.StatusConflict)
		expectAPIMessage(t, body, "still carries 1 asset(s)", "reassign them")
		if after := f.stored(t, id); after != before {
			t.Fatalf("refused dissolution changed the issuer: before %+v, after %+v", before, after)
		}
	})

	t.Run("an unknown issuer answers 404", func(t *testing.T) {
		code, body := dissolve(unknownID)
		expectCode(t, code, body, http.StatusNotFound)
		expectAPIMessage(t, body, "Issuer not found")
	})
}

// An unexpected failure must be a 500 with a generic message: the caller did
// nothing wrong, and the driver error describes the infrastructure.
func TestIssuerRoutesAnswer500WhenTheDatabaseFails(t *testing.T) {
	h := newOrderHarness(t)
	f := newAssetFixtures(t, h)
	manager := newTrackedUser(t, h, authUtils.RoleManager)
	id := f.activeIssuer(t, "Tooken Villa Belair SA")

	for _, c := range []struct {
		name     string
		path     string
		fallback string
	}{
		{"list", "/assets/issuers", "Unable to retrieve issuers."},
		{"detail", issuerPath(id), "Unable to retrieve issuer."},
	} {
		t.Run(c.name, func(t *testing.T) {
			withClosedDatabase(t, func() {
				code, body := callAPI(h, http.MethodGet, c.path, manager.token, "")
				expectCode(t, code, body, http.StatusInternalServerError)
				expectAPIMessage(t, body, c.fallback)
				if strings.Contains(body, "database is closed") {
					t.Fatalf("the driver error leaked to the client: %s", body)
				}
			})
		})
	}
}
