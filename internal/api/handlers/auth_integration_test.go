//go:build integration

// Coverage of the auth endpoints against a real PostgreSQL. What no unit test
// can prove is checked here: that the role stored in the database reaches both
// the JWT and the response body, that a database failure is reported as a
// failure instead of being disguised as a wrong password, and that the
// sessions handed out work, and only where they should, through the real
// middleware.
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
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/middleware"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// authRouter serves the whole API behind the real middleware, on the database
// named by TEST_DATABASE_URL.
func authRouter(t *testing.T) (*gin.Engine, *sql.DB) {
	t.Helper()
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
	return r, db
}

// send serves one request; bearer is left out when empty.
func send(r *gin.Engine, method, path, bearer, body string) (int, string) {
	req := httptest.NewRequest(method, globals.BaseURL+path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

// decodeSession reads the session of a sign-up, sign-in or refresh answer.
// encoding/json matches keys regardless of case, so the key itself is checked
// on the raw body: the generated client of the frontend will not be as lenient.
func decodeSession(t *testing.T, body string) server.JwtToken {
	t.Helper()
	var keys map[string]json.RawMessage
	if err := json.Unmarshal([]byte(body), &keys); err != nil {
		t.Fatalf("body is not a JSON object: %v: %s", err, body)
	}
	if _, ok := keys["data"]; !ok {
		t.Fatalf(`no "data" key in %s`, body)
	}
	var session server.JwtToken
	if err := json.Unmarshal(keys["data"], &session); err != nil {
		t.Fatalf("data is not a session: %v: %s", err, body)
	}
	if session.Token == "" || session.RefreshToken == "" || session.User.Id == 0 {
		t.Fatalf("incomplete session: %s", body)
	}
	return session
}

func signUpFresh(t *testing.T, r *gin.Engine, db *sql.DB, email string) server.JwtToken {
	t.Helper()
	if _, err := db.Exec(`DELETE FROM usr.users WHERE email = $1`, email); err != nil {
		t.Fatal(err)
	}
	code, body := send(r, http.MethodPost, "/auth/sign-up", "",
		`{"email":"`+email+`","password":"Sup3r-secret!","fullName":"Session Probe"}`)
	if code != http.StatusCreated {
		t.Fatalf("sign-up %s: got %d %s", email, code, body)
	}
	return decodeSession(t, body)
}

func TestSignUpOpensASession(t *testing.T) {
	r, db := authRouter(t)

	const email = "sign-up-session@tooken.test"
	before := time.Now()
	session := signUpFresh(t, r, db, email)

	if string(session.User.Email) != email || session.User.FullName != "Session Probe" {
		t.Errorf("user: got %s %q, want %s %q", session.User.Email, session.User.FullName, email, "Session Probe")
	}
	if session.User.Role == nil || string(*session.User.Role) != authUtils.RoleUser {
		t.Errorf("user role: got %v, want USER", session.User.Role)
	}
	// Read the way GET /users/me reads it, so the KYC state is there too.
	if session.User.KycStatus == nil || *session.User.KycStatus != "none" {
		t.Errorf("user kyc_status: got %v, want none", session.User.KycStatus)
	}
	if session.User.UpdatedAt.IsZero() {
		t.Error("user updated_at is zero")
	}
	if d := session.ExpiresAt.Sub(before); d < 59*time.Minute || d > 61*time.Minute {
		t.Errorf("expiresAt is %s away, want about an hour", d)
	}

	if code, body := send(r, http.MethodGet, "/users/me", session.Token, ""); code != http.StatusOK {
		t.Errorf("GET /users/me with the access token: got %d %s, want 200", code, body)
	}
	if code, body := send(r, http.MethodGet, "/users/me", session.RefreshToken, ""); code != http.StatusUnauthorized {
		t.Errorf("GET /users/me with the refresh token: got %d %s, want 401", code, body)
	}
}

// The body is bound without validation, so a missing field arrives empty. An
// empty email and password used to make an account anyone could sign in to,
// and the answer, which renders the email, came back as a 201 with no body.
// A password bcrypt cannot hash is the client's mistake too, not a 500.
func TestSignUpValidatesTheBody(t *testing.T) {
	r, db := authRouter(t)

	const incomplete = `email IN ('', 'no-password@tooken.test', 'no-name@tooken.test', 'long-password@tooken.test')`
	if _, err := db.Exec(`DELETE FROM usr.users WHERE ` + incomplete); err != nil {
		t.Fatal(err)
	}

	tests := map[string]string{
		"empty object":   `{}`,
		"no password":    `{"email":"no-password@tooken.test","fullName":"No Password"}`,
		"no email":       `{"password":"Sup3r-secret!","fullName":"No Email"}`,
		"blank fullName": `{"email":"no-name@tooken.test","password":"Sup3r-secret!","fullName":"  "}`,
		// 25 runes, but 75 bytes for bcrypt: "€" is three bytes in UTF-8.
		"password over 72 bytes": `{"email":"long-password@tooken.test","password":"` +
			strings.Repeat("€", 25) + `","fullName":"Long Password"}`,
	}

	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			code, resp := send(r, http.MethodPost, "/auth/sign-up", "", body)
			if code != http.StatusBadRequest {
				t.Fatalf("got %d %q, want 400", code, resp)
			}
			expectMessage(t, resp)
		})
	}

	var n int
	if err := db.QueryRow(`SELECT count(*) FROM usr.users WHERE ` + incomplete).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("%d account(s) created from an invalid body", n)
	}
}

func TestRefresh(t *testing.T) {
	r, db := authRouter(t)

	const email = "refresh-session@tooken.test"
	session := signUpFresh(t, r, db, email)

	refresh := func(token string) (int, string) {
		return send(r, http.MethodPost, "/auth/refresh", "", `{"refreshToken":"`+token+`"}`)
	}

	t.Run("the new access token carries the current role", func(t *testing.T) {
		if _, err := db.Exec(`UPDATE usr.users SET role = 'MANAGER' WHERE email = $1`, email); err != nil {
			t.Fatal(err)
		}
		// Token times are whole seconds: a refresh token issued within the
		// same second as the original would be byte for byte identical, and
		// an extended session would go unnoticed.
		time.Sleep(1100 * time.Millisecond)
		code, body := refresh(session.RefreshToken)
		if code != http.StatusOK {
			t.Fatalf("got %d %s, want 200", code, body)
		}
		renewed := decodeSession(t, body)

		claims, err := authUtils.ParseJWT(renewed.Token)
		if err != nil {
			t.Fatalf("renewed access token refused: %v", err)
		}
		if claims.Role != authUtils.RoleManager || renewed.User.Role == nil || string(*renewed.User.Role) != authUtils.RoleManager {
			t.Errorf("role: claim %q, body %v, want MANAGER", claims.Role, renewed.User.Role)
		}
		// The end of the session is fixed at sign-in: refreshing must not
		// push it back.
		original, err := authUtils.ParseRefreshToken(session.RefreshToken)
		if err != nil {
			t.Fatal(err)
		}
		returned, err := authUtils.ParseRefreshToken(renewed.RefreshToken)
		if err != nil {
			t.Fatalf("returned refresh token refused: %v", err)
		}
		if !returned.ExpiresAt.Equal(original.ExpiresAt.Time) {
			t.Errorf("session end moved from %s to %s", original.ExpiresAt, returned.ExpiresAt)
		}
		if code, body := send(r, http.MethodGet, "/users/me", renewed.Token, ""); code != http.StatusOK {
			t.Errorf("GET /users/me with the renewed token: got %d %s, want 200", code, body)
		}
	})

	t.Run("an access token is not a refresh token", func(t *testing.T) {
		code, body := refresh(session.Token)
		if code != http.StatusUnauthorized {
			t.Fatalf("got %d %s, want 401", code, body)
		}
		expectMessage(t, body)
	})

	t.Run("a forged token is a 401", func(t *testing.T) {
		code, body := refresh("not.a.token")
		if code != http.StatusUnauthorized {
			t.Fatalf("got %d %s, want 401", code, body)
		}
		expectMessage(t, body)
	})

	t.Run("a malformed body is a 400 in the message envelope", func(t *testing.T) {
		code, body := send(r, http.MethodPost, "/auth/refresh", "", `{bad`)
		if code != http.StatusBadRequest {
			t.Fatalf("got %d %s, want 400", code, body)
		}
		expectMessage(t, body)
	})

	t.Run("a deleted account can no longer refresh", func(t *testing.T) {
		const gone = "refresh-deleted@tooken.test"
		doomed := signUpFresh(t, r, db, gone)
		if _, err := db.Exec(`DELETE FROM usr.users WHERE email = $1`, gone); err != nil {
			t.Fatal(err)
		}
		code, body := refresh(doomed.RefreshToken)
		if code != http.StatusUnauthorized {
			t.Fatalf("got %d %s, want 401", code, body)
		}
		expectMessage(t, body)
	})

	// The refresh token is still good: telling the client to sign in again
	// would throw away a session that a retry can save.
	t.Run("a database outage is a 500, not a 401", func(t *testing.T) {
		broken, err := sql.Open("postgres", "postgres://nobody@127.0.0.1:1/none?sslmode=disable")
		if err != nil {
			t.Fatal(err)
		}
		globals.DB = broken
		defer func() {
			globals.DB = db
			broken.Close()
		}()

		code, body := refresh(session.RefreshToken)
		if code != http.StatusInternalServerError {
			t.Fatalf("got %d %s, want 500", code, body)
		}
		expectMessage(t, body)
	})
}

func TestSignInCarriesTheRole(t *testing.T) {
	r, db := authRouter(t)

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
	signUpBody := `{"email":"` + email + `","password":"Sup3r-secret!","fullName":"Role Probe"}`
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
