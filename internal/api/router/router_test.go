package router

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/handlers"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/gin-gonic/gin"
)

const (
	specPath      = "../../../api/openapi.yaml"
	allowedOrigin = "https://app.tooken.test"
)

func newTestEngine(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = io.Discard

	engine, err := New(handlers.NewHandler(), specPath, []string{allowedOrigin})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return engine
}

func serve(engine *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func TestNewServesTheAPIUnderItsBaseURL(t *testing.T) {
	engine := newTestEngine(t)

	w := serve(engine, httptest.NewRequest(http.MethodGet, globals.BaseURL+"/health", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"status":"ok"`) {
		t.Fatalf("GET health = %d %s, want 200 with status ok", w.Code, w.Body.String())
	}

	w = serve(engine, httptest.NewRequest(http.MethodGet, "/health", nil))
	if w.Code != http.StatusNotFound {
		t.Fatalf("GET /health outside the base URL = %d, want 404", w.Code)
	}
}

// The middleware reads which operations are protected from the spec given to
// New: a route declared with bearerAuth must refuse a call without a token
// before any handler — and so any database — is reached.
func TestNewProtectsTheRoutesTheSpecDeclaresProtected(t *testing.T) {
	engine := newTestEngine(t)

	for _, route := range []struct{ method, path string }{
		{http.MethodGet, "/users/me"},
		{http.MethodPost, "/trex/init/implementations"},
		{http.MethodPost, "/assets/real-estate/issuance/orders"},
		{http.MethodPost, "/assets/real-estate/issuance/orders/ISS-ANY/confirm-payment"},
	} {
		w := serve(engine, httptest.NewRequest(route.method, globals.BaseURL+route.path, nil))
		if w.Code != http.StatusUnauthorized {
			t.Errorf("%s %s without a token = %d, want 401", route.method, route.path, w.Code)
		}
	}
}

// Without a token, the handlers refuse on their own, so the test above would
// pass with no middleware at all. The middleware is the only place a token
// becomes claims: a valid USER session on a staff route is forbidden, not
// anonymous, and a forged token is named as such.
func TestNewTurnsTheBearerTokenIntoClaims(t *testing.T) {
	engine := newTestEngine(t)
	route := globals.BaseURL + "/trex/init/implementations"

	token, err := authUtils.GenerateJWT(42, "investor@tooken.test", authUtils.RoleUser)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, route, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	if w := serve(engine, req); w.Code != http.StatusForbidden {
		t.Fatalf("a USER token on a staff route = %d %s, want 403", w.Code, w.Body.String())
	}

	req = httptest.NewRequest(http.MethodPost, route, nil)
	req.Header.Set("Authorization", "Bearer not-a-jwt")
	w := serve(engine, req)
	if w.Code != http.StatusUnauthorized || !strings.Contains(w.Body.String(), "Invalid or expired token") {
		t.Fatalf("a forged token = %d %s, want 401 naming the token", w.Code, w.Body.String())
	}
}

func TestNewAppliesTheCORSPolicy(t *testing.T) {
	engine := newTestEngine(t)

	preflight := func(origin string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodOptions, globals.BaseURL+"/assets/real-estates/active", nil)
		req.Header.Set("Origin", origin)
		req.Header.Set("Access-Control-Request-Method", http.MethodGet)
		req.Header.Set("Access-Control-Request-Headers", "Authorization")
		return serve(engine, req)
	}

	w := preflight(allowedOrigin)
	if w.Code != http.StatusNoContent {
		t.Fatalf("preflight from an allowed origin = %d, want 204", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != allowedOrigin {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, allowedOrigin)
	}
	if got := w.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("Access-Control-Allow-Credentials = %q, want true", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Headers"); !strings.Contains(got, "Authorization") {
		t.Fatalf("Access-Control-Allow-Headers = %q, want Authorization allowed", got)
	}

	w = preflight("https://evil.example")
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("an unknown origin was allowed: Access-Control-Allow-Origin = %q", got)
	}
	if w.Code != http.StatusForbidden {
		t.Fatalf("preflight from an unknown origin = %d, want 403", w.Code)
	}
}

func TestNewServesTheSpecAndItsDocumentation(t *testing.T) {
	engine := newTestEngine(t)

	w := serve(engine, httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil))
	if w.Code != http.StatusOK || !strings.HasPrefix(w.Body.String(), "openapi:") {
		t.Fatalf("GET /openapi.yaml = %d, want 200 with the spec", w.Code)
	}

	w = serve(engine, httptest.NewRequest(http.MethodGet, "/docs/", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "Tooken API Docs") {
		t.Fatalf("GET /docs/ = %d, want 200 with the documentation UI", w.Code)
	}
}

// Each of these used to stop the server with a panic or a log.Fatal from
// inside the construction; New reports them instead, so main decides.
func TestNewRefusesAConfigurationItCannotServe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	for _, c := range []struct {
		name    string
		spec    string
		origins []string
		want    string
	}{
		{"no allowed origin", specPath, nil, "CORS"},
		{"an origin without its scheme", specPath, []string{"app.tooken.test"}, "CORS"},
		{"a spec that does not exist", "./missing-openapi.yaml", []string{allowedOrigin}, "authentication"},
	} {
		t.Run(c.name, func(t *testing.T) {
			engine, err := New(handlers.NewHandler(), c.spec, c.origins)
			if err == nil {
				t.Fatalf("New accepted the configuration and returned %v", engine)
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Fatalf("error %q does not name the %s setting", err, c.want)
			}
		})
	}
}
