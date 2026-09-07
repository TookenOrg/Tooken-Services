package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/gin-gonic/gin"
)

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

	protectedRoutes = map[string]map[string]AuthMode{
		"/orders/:orderRef": {"GET": AuthRequired},
		"/open/:id":         {"POST": AuthRequired}, // protected on POST only
	}

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

	protectedRoutes = map[string]map[string]AuthMode{
		"/assets/:id": {"GET": AuthOptional},
	}

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

	protectedRoutes = map[string]map[string]AuthMode{
		"/assets/:id": {"GET": AuthOptional},
	}

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
