//go:build integration

// Coverage of the sign-in endpoint against a real PostgreSQL. Two things are
// checked that no unit test can prove: that the role stored in the database
// reaches both the JWT and the response body, and that a database failure is
// reported as a failure instead of being disguised as a wrong password.
//
// Run it with the same disposable database as the real estate suite, see
// real_estate_write_integration_test.go.
package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
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

func TestSignInCarriesTheRole(t *testing.T) {
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

	signIn := func(email, password string) (int, string) {
		body := `{"email":"` + email + `","password":"` + password + `"}`
		req := httptest.NewRequest(http.MethodPost, globals.BaseURL+"/auth/sign-in",
			strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code, w.Body.String()
	}

	// A manager account created through the real sign-up endpoint, then promoted
	// in the database the way an operator would.
	const email = "role-probe@tooken.test"
	if _, err := db.Exec(`DELETE FROM usr.users WHERE email = $1`, email); err != nil {
		t.Fatal(err)
	}
	signUpBody := `{"email":"` + email + `","password":"Sup3r-secret!","full_name":"Role Probe"}`
	req := httptest.NewRequest(http.MethodPost, globals.BaseURL+"/auth/sign-up",
		strings.NewReader(signUpBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("sign-up: got %d %s", w.Code, w.Body.String())
	}
	if _, err := db.Exec(`UPDATE usr.users SET role = 'MANAGER' WHERE email = $1`, email); err != nil {
		t.Fatal(err)
	}

	// Before the fix the service reported the conflict through its status code
	// alone, leaving err nil, so the handler fell through to the success branch
	// and answered 201 with an empty token.
	t.Run("signing up twice is a 409, not a 201 with no token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, globals.BaseURL+"/auth/sign-up",
			strings.NewReader(signUpBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusConflict {
			t.Fatalf("got %d %s, want 409", w.Code, w.Body.String())
		}
	})

	t.Run("the role reaches the token and the body", func(t *testing.T) {
		code, body := signIn(email, "Sup3r-secret!")
		if code != http.StatusOK {
			t.Fatalf("got %d %s", code, body)
		}
		var resp struct {
			Data struct {
				Token string `json:"token"`
				User  struct {
					Role string `json:"role"`
				} `json:"user"`
			} `json:"data"`
		}
		if err := json.Unmarshal([]byte(body), &resp); err != nil {
			t.Fatal(err)
		}
		if resp.Data.User.Role != authUtils.RoleManager {
			t.Errorf("response body role: got %q, want MANAGER", resp.Data.User.Role)
		}
		raw, err := base64.RawURLEncoding.DecodeString(strings.Split(resp.Data.Token, ".")[1])
		if err != nil {
			t.Fatal(err)
		}
		var claims struct {
			Role string `json:"role"`
		}
		if err := json.Unmarshal(raw, &claims); err != nil {
			t.Fatal(err)
		}
		if claims.Role != authUtils.RoleManager {
			t.Errorf("jwt claim role: got %q, want MANAGER, payload=%s", claims.Role, raw)
		}
	})

	t.Run("a wrong password is a 401", func(t *testing.T) {
		if code, body := signIn(email, "not-the-password"); code != http.StatusUnauthorized {
			t.Errorf("got %d %s, want 401", code, body)
		}
	})

	t.Run("an unknown email is a 401", func(t *testing.T) {
		if code, body := signIn("nobody@tooken.test", "whatever"); code != http.StatusUnauthorized {
			t.Errorf("got %d %s, want 401", code, body)
		}
	})

	// The regression that started all this: when the database cannot answer, the
	// endpoint used to reply "Invalid email or password", which points at the
	// user instead of at the incident.
	t.Run("a database outage is a 500, not a 401", func(t *testing.T) {
		broken, err := sql.Open("postgres", "postgres://nobody:nobody@127.0.0.1:1/none?sslmode=disable")
		if err != nil {
			t.Fatal(err)
		}
		globals.DB = broken
		defer func() {
			globals.DB = db
			broken.Close()
		}()

		if code, body := signIn(email, "Sup3r-secret!"); code != http.StatusInternalServerError {
			t.Errorf("got %d %s, want 500", code, body)
		}
	})
}
