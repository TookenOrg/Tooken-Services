//go:build integration

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
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/middleware"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func TestCurrentUserProfileAndKYCProjection(t *testing.T) {
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
	t.Cleanup(func() { db.Close() })

	prefix := fmt.Sprintf("m2-2-profile-%d", time.Now().UnixNano())
	t.Cleanup(func() {
		pattern := prefix + "%@tooken.test"
		_, _ = db.Exec(`
		    DELETE FROM usr.kyc_verification
		    WHERE user_id IN (SELECT id FROM usr.users WHERE email LIKE $1)
		       OR decided_by IN (SELECT id FROM usr.users WHERE email LIKE $1)
		       OR revoked_by IN (SELECT id FROM usr.users WHERE email LIKE $1)`, pattern)
		_, _ = db.Exec(`
		    DELETE FROM blk.user_wallet
		    WHERE user_id IN (SELECT id FROM usr.users WHERE email LIKE $1)`, pattern)
		_, _ = db.Exec(`DELETE FROM usr.users WHERE email LIKE $1`, pattern)
	})

	if err := middleware.InitAuth("../../../api/openapi.yaml"); err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	g := r.Group(globals.BaseURL)
	g.Use(middleware.AutoAuthMiddleware())
	server.RegisterHandlers(g, NewHandler())

	do := func(method, url, token string) (int, string) {
		req := httptest.NewRequest(method, globals.BaseURL+url, bytes.NewReader(nil))
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code, w.Body.String()
	}
	newUser := func(name, role string) (int, string) {
		email := prefix + "-" + name + "@tooken.test"
		var id int
		if err := db.QueryRow(`
		    INSERT INTO usr.users (full_name, email, password, role)
		    VALUES ($1, $2, 'x', $3)
		    RETURNING id`, "M2 "+name, email, role).Scan(&id); err != nil {
			t.Fatal(err)
		}
		token, err := authUtils.GenerateJWT(id, email, role)
		if err != nil {
			t.Fatal(err)
		}
		return id, token
	}
	readProfile := func(t *testing.T, token string) map[string]any {
		t.Helper()
		code, body := do(http.MethodGet, "/users/me", token)
		if code != http.StatusOK {
			t.Fatalf("GET /users/me: got %d %s, want 200", code, body)
		}
		var got map[string]any
		if err := json.Unmarshal([]byte(body), &got); err != nil {
			t.Fatal(err)
		}
		return got
	}

	managerID, _ := newUser("manager", authUtils.RoleManager)
	userID, userToken := newUser("user", authUtils.RoleUser)

	t.Run("anonymous is refused by the OpenAPI-driven auth middleware", func(t *testing.T) {
		code, body := do(http.MethodGet, "/users/me", "")
		if code != http.StatusUnauthorized {
			t.Fatalf("got %d %s, want 401", code, body)
		}
	})

	t.Run("a missing user from a valid token is a 404", func(t *testing.T) {
		token, err := authUtils.GenerateJWT(999_999_991, prefix+"-missing@tooken.test", authUtils.RoleUser)
		if err != nil {
			t.Fatal(err)
		}
		code, body := do(http.MethodGet, "/users/me", token)
		if code != http.StatusNotFound {
			t.Fatalf("got %d %s, want 404", code, body)
		}
	})

	t.Run("a user with no KYC sees only profile fields and status none", func(t *testing.T) {
		body := readProfile(t, userToken)
		if body["id"] != float64(userID) {
			t.Errorf("id = %v, want %d", body["id"], userID)
		}
		if body["kyc_status"] != "none" {
			t.Errorf("kyc_status = %v, want none", body["kyc_status"])
		}
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(raw), "password") || strings.Contains(string(raw), "wallet") {
			t.Fatalf("profile leaks forbidden data: %s", raw)
		}
	})

	t.Run("submitted KYC is projected as pending", func(t *testing.T) {
		if _, err := db.Exec(`
		    INSERT INTO usr.kyc_verification (user_id, status, declared_full_name, declared_country_code)
		    VALUES ($1, 'submitted', 'M2 User', 'FR')`, userID); err != nil {
			t.Fatal(err)
		}

		var projected string
		if err := db.QueryRow(`SELECT kyc_status FROM usr.users WHERE id = $1`, userID).Scan(&projected); err != nil {
			t.Fatal(err)
		}
		if projected != "pending" {
			t.Fatalf("projected status = %q, want pending", projected)
		}
		if got := readProfile(t, userToken)["kyc_status"]; got != "pending" {
			t.Fatalf("API kyc_status = %v, want pending", got)
		}
	})

	t.Run("approved KYC projects country and does not pretend to be on-chain verified", func(t *testing.T) {
		decidedAt := time.Now().UTC()
		expiresAt := decidedAt.Add(48 * time.Hour)
		if _, err := db.Exec(`
		    UPDATE usr.kyc_verification
		    SET status = 'approved',
		        decided_at = $2,
		        decided_by = $3,
		        expires_at = $4,
		        declared_country_code = 'FR'
		    WHERE user_id = $1 AND status = 'submitted'`,
			userID, decidedAt, managerID, expiresAt); err != nil {
			t.Fatal(err)
		}

		var status string
		var country sql.NullString
		var verifiedAt sql.NullTime
		if err := db.QueryRow(`
		    SELECT kyc_status, country_code, kyc_verified_at
		    FROM usr.users WHERE id = $1`, userID).Scan(&status, &country, &verifiedAt); err != nil {
			t.Fatal(err)
		}
		if status != "approved" {
			t.Fatalf("projected status = %q, want approved", status)
		}
		if !country.Valid || country.String != "FR" {
			t.Fatalf("country_code = %v, want FR", country)
		}
		if verifiedAt.Valid {
			t.Fatalf("approved KYC should not set kyc_verified_at: %v", verifiedAt.Time)
		}

		body := readProfile(t, userToken)
		if body["kyc_status"] != "approved" || body["country_code"] != "FR" {
			t.Fatalf("API projection = %v, want approved/FR", body)
		}
	})

	t.Run("verified KYC is returned as expired once its validity date has passed", func(t *testing.T) {
		if _, err := db.Exec(`
		    UPDATE usr.users
		    SET kyc_status = 'verified',
		        kyc_verified_at = now() - interval '2 days',
		        kyc_expires_at = now() + interval '1 day'
		    WHERE id = $1`, userID); err != nil {
			t.Fatal(err)
		}
		if got := readProfile(t, userToken)["kyc_status"]; got != "verified" {
			t.Fatalf("future verified KYC = %v, want verified", got)
		}

		if _, err := db.Exec(`
		    UPDATE usr.users
		    SET kyc_status = 'verified',
		        kyc_verified_at = now() - interval '2 days',
		        kyc_expires_at = now() - interval '1 day'
		    WHERE id = $1`, userID); err != nil {
			t.Fatal(err)
		}
		if got := readProfile(t, userToken)["kyc_status"]; got != "expired" {
			t.Fatalf("expired verified KYC = %v, want expired", got)
		}
	})
}

func TestInvestorDatabaseGuards(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	globals.DB = db
	t.Cleanup(func() { db.Close() })

	prefix := fmt.Sprintf("m2-2-guards-%d", time.Now().UnixNano())
	t.Cleanup(func() {
		pattern := prefix + "%@tooken.test"
		_, _ = db.Exec(`
		    DELETE FROM usr.kyc_verification
		    WHERE user_id IN (SELECT id FROM usr.users WHERE email LIKE $1)
		       OR decided_by IN (SELECT id FROM usr.users WHERE email LIKE $1)
		       OR revoked_by IN (SELECT id FROM usr.users WHERE email LIKE $1)`, pattern)
		_, _ = db.Exec(`
		    DELETE FROM blk.user_wallet
		    WHERE user_id IN (SELECT id FROM usr.users WHERE email LIKE $1)`, pattern)
		_, _ = db.Exec(`DELETE FROM usr.users WHERE email LIKE $1`, pattern)
	})

	newUser := func(name, role string) int {
		var id int
		if err := db.QueryRow(`
		    INSERT INTO usr.users (full_name, email, password, role)
		    VALUES ($1, $2, 'x', $3)
		    RETURNING id`, "M2 "+name, prefix+"-"+name+"@tooken.test", role).Scan(&id); err != nil {
			t.Fatal(err)
		}
		return id
	}

	managerID := newUser("manager", authUtils.RoleManager)

	t.Run("rejected KYC without a rejection reason is refused", func(t *testing.T) {
		userID := newUser("rejected", authUtils.RoleUser)
		_, err := db.Exec(`
		    INSERT INTO usr.kyc_verification (user_id, status, decided_at, decided_by)
		    VALUES ($1, 'rejected', now(), $2)`, userID, managerID)
		if err == nil {
			t.Fatal("rejected KYC without rejection_reason was accepted")
		}
	})

	t.Run("only one submitted KYC can remain open per user", func(t *testing.T) {
		userID := newUser("submitted", authUtils.RoleUser)
		if _, err := db.Exec(`
		    INSERT INTO usr.kyc_verification (user_id, status, declared_country_code)
		    VALUES ($1, 'submitted', 'LU')`, userID); err != nil {
			t.Fatal(err)
		}
		_, err := db.Exec(`
		    INSERT INTO usr.kyc_verification (user_id, status, declared_country_code)
		    VALUES ($1, 'submitted', 'FR')`, userID)
		if err == nil {
			t.Fatal("second submitted KYC was accepted")
		}
	})

	t.Run("only one active wallet can exist per user", func(t *testing.T) {
		userID := newUser("wallet", authUtils.RoleUser)
		firstWallet := fmt.Sprintf("0x%040x", userID*10+1)
		secondWallet := fmt.Sprintf("0x%040x", userID*10+2)
		if _, err := db.Exec(`
		    INSERT INTO blk.user_wallet (user_id, wallet_address, label)
		    VALUES ($1, $2, 'Main')`, userID, firstWallet); err != nil {
			t.Fatal(err)
		}
		_, err := db.Exec(`
		    INSERT INTO blk.user_wallet (user_id, wallet_address, label)
		    VALUES ($1, $2, 'Second')`, userID, secondWallet)
		if err == nil {
			t.Fatal("second active wallet was accepted")
		}
	})
}
