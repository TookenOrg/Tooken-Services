package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/gin-gonic/gin"
)

// useRoutes stands in for InitAuth for the duration of the test.
func useRoutes(t *testing.T, routes map[string]map[string]AuthMode) {
	t.Helper()
	previous := protectedRoutes.Load()
	if routes == nil {
		protectedRoutes.Store(nil)
	} else {
		protectedRoutes.Store(&routes)
	}
	t.Cleanup(func() { protectedRoutes.Store(previous) })
}

// guardedEngine serves GET /orders/:ref behind the middleware.
func guardedEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AutoAuthMiddleware())
	r.GET("/orders/:ref", func(c *gin.Context) { c.String(http.StatusOK, "reached") })
	return r
}

func TestToGinPath(t *testing.T) {
	tests := map[string]string{
		"/payments/add":                   "/payments/add",
		"/assets/real-estates/{id}":       "/assets/real-estates/:id",
		"/orders/{orderRef}":              "/orders/:orderRef",
		"/a/{first}/b/{second}":           "/a/:first/b/:second",
		"/assets/real-estates/{id}/media": "/assets/real-estates/:id/media",
	}

	for in, want := range tests {
		if got := toGinPath(in); got != want {
			t.Errorf("toGinPath(%q) = %q, want %q", in, got, want)
		}
	}
}

// A protected route carrying a path parameter used to be reachable without any
// token: the map was keyed on the OpenAPI template while the middleware
// compared it to the concrete request path, so the two never matched.
func TestProtectedRouteWithPathParameterRequiresAToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	useRoutes(t, map[string]map[string]AuthMode{
		"/orders/:orderRef": {"GET": AuthRequired},
		"/open/:id":         {"POST": AuthRequired}, // protected on POST only
	})

	r := gin.New()
	r.Use(AutoAuthMiddleware())
	reached := func(c *gin.Context) { c.String(http.StatusOK, "reached") }
	r.GET("/orders/:orderRef", reached)
	r.GET("/open/:id", reached)

	tests := []struct {
		name   string
		method string
		url    string
		want   int
	}{
		{"protected route with a parameter", "GET", "/orders/ORD-42", http.StatusUnauthorized},
		{"unprotected method on the same shape", "GET", "/open/7", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(tt.method, tt.url, nil))
			if w.Code != tt.want {
				t.Errorf("want %d, got %d (%s)", tt.want, w.Code, w.Body.String())
			}
		})
	}
}

// An optional-auth route serves anonymous callers, but must still reject a
// token it cannot trust rather than quietly answering with the public view.
func TestOptionalAuthRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	useRoutes(t, map[string]map[string]AuthMode{
		"/assets/:id": {"GET": AuthOptional},
	})

	r := gin.New()
	r.Use(AutoAuthMiddleware())
	r.GET("/assets/:id", func(c *gin.Context) {
		if _, ok := GetUserClaims(c); ok {
			c.String(http.StatusOK, "identified")
			return
		}
		c.String(http.StatusOK, "anonymous")
	})

	tests := []struct {
		name   string
		header string
		want   int
		body   string
	}{
		{"no token at all", "", http.StatusOK, "anonymous"},
		{"garbage token", "Bearer not-a-jwt", http.StatusUnauthorized, ""},
		{"header without the scheme", "some-token", http.StatusUnauthorized, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/assets/7", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.want {
				t.Fatalf("want %d, got %d (%s)", tt.want, w.Code, w.Body.String())
			}
			if tt.body != "" && w.Body.String() != tt.body {
				t.Errorf("want body %q, got %q", tt.body, w.Body.String())
			}
		})
	}
}

// A valid token on an optional-auth route must reach the handler as claims,
// carrying the role decided in migration 000009.
func TestOptionalAuthRouteParsesAValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	token, err := utils.GenerateJWT(7, "manager@tooken.lu", utils.RoleManager)
	if err != nil {
		t.Fatalf("GenerateJWT: %v", err)
	}

	useRoutes(t, map[string]map[string]AuthMode{
		"/assets/:id": {"GET": AuthOptional},
	})

	r := gin.New()
	r.Use(AutoAuthMiddleware())
	r.GET("/assets/:id", func(c *gin.Context) {
		claims, ok := GetUserClaims(c)
		if !ok {
			c.String(http.StatusOK, "anonymous")
			return
		}
		c.String(http.StatusOK, claims.Role)
	})

	req := httptest.NewRequest("GET", "/assets/7", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK || w.Body.String() != utils.RoleManager {
		t.Errorf("want 200 %q, got %d %q", utils.RoleManager, w.Code, w.Body.String())
	}
}

// RequireRole must tell "no token" from "wrong role": only the first is worth
// logging in again for.
func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name   string
		claims *utils.CustomClaims
		want   int
	}{
		{"anonymous", nil, http.StatusUnauthorized},
		{"plain user", &utils.CustomClaims{Role: utils.RoleUser}, http.StatusForbidden},
		{"manager", &utils.CustomClaims{Role: utils.RoleManager}, http.StatusOK},
		{"admin", &utils.CustomClaims{Role: utils.RoleAdmin}, http.StatusOK},
		{"empty role", &utils.CustomClaims{}, http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/assets", nil)

			if tt.claims != nil {
				c.Set("user_claims", tt.claims)
			}

			if RequireRole(c, utils.RoleManager, utils.RoleAdmin) {
				c.Status(http.StatusOK)
			}

			if w.Code != tt.want {
				t.Errorf("want %d, got %d (%s)", tt.want, w.Code, w.Body.String())
			}
		})
	}
}

// A failed reload must not leave the routes open. The list used to be
// assigned whatever the load returned, nil on failure, and a protected route
// then answered 200 to an anonymous caller.
func TestFailedInitAuthKeepsTheRoutesProtected(t *testing.T) {
	useRoutes(t, map[string]map[string]AuthMode{"/orders/:ref": {"GET": AuthRequired}})
	r := guardedEngine()

	if err := InitAuth("./missing-openapi.yaml"); err == nil {
		t.Fatal("InitAuth must report a spec it cannot read")
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/orders/1", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("after a failed InitAuth: want 401, got %d (%s)", w.Code, w.Body.String())
	}
}

// Without InitAuth nothing tells a protected route from a public one: the
// middleware refuses everything rather than serving everything.
func TestMiddlewareWithoutInitAuthRefusesEveryRequest(t *testing.T) {
	useRoutes(t, nil)
	r := guardedEngine()

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/orders/1", nil))
	if w.Code != http.StatusInternalServerError || w.Body.String() == "reached" {
		t.Fatalf("want 500 without reaching the handler, got %d (%s)", w.Code, w.Body.String())
	}
}

// Every refusal of the middleware uses the {"message": ...} envelope of the
// rest of the API, so the client reads all its errors the same way.
func TestRefusalsUseTheMessageEnvelope(t *testing.T) {
	useRoutes(t, map[string]map[string]AuthMode{"/orders/:ref": {"GET": AuthRequired}})
	r := guardedEngine()

	for name, header := range map[string]string{
		"no header":        "",
		"no Bearer scheme": "some-token",
		"invalid token":    "Bearer not-a-jwt",
	} {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/orders/1", nil)
			if header != "" {
				req.Header.Set("Authorization", header)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var body map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatalf("body is not a JSON object: %s", w.Body.String())
			}
			if w.Code != http.StatusUnauthorized || body["message"] == "" || len(body) != 1 {
				t.Fatalf(`want 401 {"message": ...}, got %d %s`, w.Code, w.Body.String())
			}
		})
	}
}

// A header without the scheme is answered with the expected format, so the
// caller can fix the request without reading the code.
func TestAMalformedHeaderIsToldTheExpectedFormat(t *testing.T) {
	useRoutes(t, map[string]map[string]AuthMode{"/orders/:ref": {"GET": AuthRequired}})
	req := httptest.NewRequest("GET", "/orders/1", nil)
	req.Header.Set("Authorization", "some-token")
	w := httptest.NewRecorder()
	guardedEngine().ServeHTTP(w, req)

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if w.Code != http.StatusUnauthorized || !strings.Contains(body["message"], "Bearer <token>") {
		t.Fatalf("want 401 naming the Bearer format, got %d %s", w.Code, w.Body.String())
	}
}
