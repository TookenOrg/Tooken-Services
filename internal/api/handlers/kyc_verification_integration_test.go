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

	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/middleware"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/TookenOrg/tooken-services/internal/api/server"
)

func TestKycVerificationDecisionRoutes(t *testing.T) {
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

	if err := middleware.InitAuth("../../../api/openapi.yaml"); err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	g := r.Group(globals.BaseURL)
	g.Use(middleware.AutoAuthMiddleware())
	server.RegisterHandlers(g, NewHandler())

	prefix := fmt.Sprintf("m2-3-kyc-decision-%d", time.Now().UnixNano())
	t.Cleanup(func() {
		pattern := prefix + "%@tooken.test"
		_, _ = db.Exec(`
		    DELETE FROM usr.kyc_verification
		    WHERE user_id IN (SELECT id FROM usr.users WHERE email LIKE $1)
		       OR decided_by IN (SELECT id FROM usr.users WHERE email LIKE $1)
		       OR revoked_by IN (SELECT id FROM usr.users WHERE email LIKE $1)`, pattern)
		_, _ = db.Exec(`DELETE FROM usr.users WHERE email LIKE $1`, pattern)
	})

	type createdUser struct {
		id    int
		email string
		token string
	}
	newUser := func(t *testing.T, name, role string) createdUser {
		t.Helper()
		email := prefix + "-" + name + "@tooken.test"
		var id int
		if err := db.QueryRow(`
		    INSERT INTO usr.users (full_name, email, password, role)
		    VALUES ($1, $2, 'x', $3)
		    RETURNING id`, "KYC "+name, email, role).Scan(&id); err != nil {
			t.Fatal(err)
		}
		token, err := authUtils.GenerateJWT(id, email, role)
		if err != nil {
			t.Fatal(err)
		}
		return createdUser{id: id, email: email, token: token}
	}
	admin := newUser(t, "admin", authUtils.RoleAdmin)
	user := newUser(t, "user", authUtils.RoleUser)

	doJSON := func(t *testing.T, method, path, token, body string) (int, string) {
		t.Helper()
		var reader *bytes.Reader
		if body == "" {
			reader = bytes.NewReader(nil)
		} else {
			reader = bytes.NewReader([]byte(body))
		}
		req := httptest.NewRequest(method, globals.BaseURL+path, reader)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		if body != "" {
			req.Header.Set("Content-Type", "application/json")
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code, w.Body.String()
	}

	insertSubmitted := func(t *testing.T, owner createdUser, suffix string) int {
		t.Helper()
		var id int
		if err := db.QueryRow(`
		    INSERT INTO usr.kyc_verification (user_id, status, declared_full_name, declared_country_code)
		    VALUES ($1, 'submitted', $2, 'FR')
		    RETURNING id`, owner.id, "KYC "+suffix).Scan(&id); err != nil {
			t.Fatal(err)
		}
		return id
	}
	insertApproved := func(t *testing.T, owner createdUser, suffix string) int {
		t.Helper()
		var id int
		if err := db.QueryRow(`
		    INSERT INTO usr.kyc_verification (
		        user_id, status, declared_full_name, declared_country_code,
		        decided_at, decided_by, expires_at
		    )
		    VALUES ($1, 'approved', $2, 'FR', now(), $3, now() + interval '30 days')
		    RETURNING id`, owner.id, "KYC "+suffix, admin.id).Scan(&id); err != nil {
			t.Fatal(err)
		}
		return id
	}
	insertRejected := func(t *testing.T, owner createdUser, suffix string) int {
		t.Helper()
		var id int
		if err := db.QueryRow(`
		    INSERT INTO usr.kyc_verification (
		        user_id, status, declared_full_name, declared_country_code,
		        decided_at, decided_by, rejection_reason
		    )
		    VALUES ($1, 'rejected', $2, 'FR', now(), $3, 'Document illisible')
		    RETURNING id`, owner.id, "KYC "+suffix, admin.id).Scan(&id); err != nil {
			t.Fatal(err)
		}
		return id
	}
	newApplicant := func(t *testing.T, name string) createdUser {
		t.Helper()
		return newUser(t, name, authUtils.RoleUser)
	}

	t.Run("detail route is protected and returns the full KYC decision record", func(t *testing.T) {
		owner := newApplicant(t, "detail")
		kycID := insertApproved(t, owner, "detail")

		if code, body := doJSON(t, http.MethodGet, fmt.Sprintf("/users/kyc/verifications/%d", kycID), "", ""); code != http.StatusUnauthorized {
			t.Fatalf("anonymous: got %d %s, want 401", code, body)
		}
		if code, body := doJSON(t, http.MethodGet, fmt.Sprintf("/users/kyc/verifications/%d", kycID), user.token, ""); code != http.StatusForbidden {
			t.Fatalf("user: got %d %s, want 403", code, body)
		}

		code, body := doJSON(t, http.MethodGet, fmt.Sprintf("/users/kyc/verifications/%d", kycID), admin.token, "")
		if code != http.StatusOK {
			t.Fatalf("admin: got %d %s, want 200", code, body)
		}
		var got map[string]any
		if err := json.Unmarshal([]byte(body), &got); err != nil {
			t.Fatal(err)
		}
		if got["id"] != float64(kycID) || got["user_id"] != float64(owner.id) || got["email"] != owner.email {
			t.Fatalf("wrong identity fields: %v", got)
		}
		for _, key := range []string{"decided_at", "decided_by", "expires_at", "submitted_at"} {
			if _, ok := got[key]; !ok {
				t.Fatalf("missing %q in detail response: %s", key, body)
			}
		}
		raw := string(mustJSON(t, got))
		if strings.Contains(raw, "password") || strings.Contains(raw, "wallet") || strings.Contains(raw, "private") {
			t.Fatalf("detail leaks forbidden data: %s", raw)
		}

		if code, body := doJSON(t, http.MethodGet, "/users/kyc/verifications/999999999", admin.token, ""); code != http.StatusNotFound {
			t.Fatalf("missing id: got %d %s, want 404", code, body)
		}
	})

	t.Run("approve validates caller, body, status and projects approved user state", func(t *testing.T) {
		owner := newApplicant(t, "approve")
		kycID := insertSubmitted(t, owner, "approve")
		path := fmt.Sprintf("/kyc/verifications/%d/approve", kycID)
		future := time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)

		if code, body := doJSON(t, http.MethodPost, path, user.token, fmt.Sprintf(`{"expires_at":%q}`, future)); code != http.StatusForbidden {
			t.Fatalf("user approve: got %d %s, want 403", code, body)
		}
		if code, body := doJSON(t, http.MethodPost, path, admin.token, `{}`); code != http.StatusBadRequest {
			t.Fatalf("missing expires_at: got %d %s, want 400", code, body)
		}
		past := time.Now().UTC().Add(-24 * time.Hour).Format(time.RFC3339)
		if code, body := doJSON(t, http.MethodPost, path, admin.token, fmt.Sprintf(`{"expires_at":%q}`, past)); code != http.StatusBadRequest {
			t.Fatalf("past expires_at: got %d %s, want 400", code, body)
		}
		if code, body := doJSON(t, http.MethodPost, "/kyc/verifications/999999999/approve", admin.token, fmt.Sprintf(`{"expires_at":%q}`, future)); code != http.StatusNotFound {
			t.Fatalf("missing id: got %d %s, want 404", code, body)
		}

		code, body := doJSON(t, http.MethodPost, path, admin.token, fmt.Sprintf(`{"expires_at":%q}`, future))
		if code != http.StatusOK {
			t.Fatalf("approve: got %d %s, want 200", code, body)
		}
		var got map[string]any
		if err := json.Unmarshal([]byte(body), &got); err != nil {
			t.Fatal(err)
		}
		if got["status"] != "approved" || got["decided_by"] != float64(admin.id) || got["decided_at"] == nil || got["expires_at"] == nil {
			t.Fatalf("approve response missing decision fields: %v", got)
		}

		var status, country string
		var expiresAt sql.NullTime
		var verifiedAt sql.NullTime
		if err := db.QueryRow(`
		    SELECT kyc_status, country_code, kyc_expires_at, kyc_verified_at
		    FROM usr.users WHERE id = $1`, owner.id).Scan(&status, &country, &expiresAt, &verifiedAt); err != nil {
			t.Fatal(err)
		}
		if status != "approved" || country != "FR" || !expiresAt.Valid || verifiedAt.Valid {
			t.Fatalf("projection = status %q country %q expires %v verified %v, want approved/FR/expires/no verified",
				status, country, expiresAt, verifiedAt)
		}
		if code, body := doJSON(t, http.MethodPost, path, admin.token, fmt.Sprintf(`{"expires_at":%q}`, future)); code != http.StatusConflict {
			t.Fatalf("approve twice: got %d %s, want 409", code, body)
		}
	})

	t.Run("reject validates caller, reason, status and projects rejected user state", func(t *testing.T) {
		owner := newApplicant(t, "reject")
		kycID := insertSubmitted(t, owner, "reject")
		path := fmt.Sprintf("/kyc/verifications/%d/reject", kycID)

		if code, body := doJSON(t, http.MethodPost, path, user.token, `{"reason":"Document illisible"}`); code != http.StatusForbidden {
			t.Fatalf("user reject: got %d %s, want 403", code, body)
		}
		if code, body := doJSON(t, http.MethodPost, path, admin.token, `{}`); code != http.StatusBadRequest {
			t.Fatalf("missing reason: got %d %s, want 400", code, body)
		}
		if code, body := doJSON(t, http.MethodPost, path, admin.token, `{"reason":"   "}`); code != http.StatusBadRequest {
			t.Fatalf("blank reason: got %d %s, want 400", code, body)
		}
		if code, body := doJSON(t, http.MethodPost, "/kyc/verifications/999999999/reject", admin.token, `{"reason":"Document illisible"}`); code != http.StatusNotFound {
			t.Fatalf("missing id: got %d %s, want 404", code, body)
		}

		code, body := doJSON(t, http.MethodPost, path, admin.token, `{"reason":"  Document illisible  "}`)
		if code != http.StatusOK {
			t.Fatalf("reject: got %d %s, want 200", code, body)
		}
		var got map[string]any
		if err := json.Unmarshal([]byte(body), &got); err != nil {
			t.Fatal(err)
		}
		if got["status"] != "rejected" || got["decided_by"] != float64(admin.id) || got["rejection_reason"] != "Document illisible" {
			t.Fatalf("reject response missing decision fields: %v", got)
		}

		var status string
		var expiresAt sql.NullTime
		if err := db.QueryRow(`SELECT kyc_status, kyc_expires_at FROM usr.users WHERE id = $1`, owner.id).Scan(&status, &expiresAt); err != nil {
			t.Fatal(err)
		}
		if status != "rejected" || expiresAt.Valid {
			t.Fatalf("projection = status %q expires %v, want rejected/null", status, expiresAt)
		}
		if code, body := doJSON(t, http.MethodPost, path, admin.token, `{"reason":"Second decision"}`); code != http.StatusConflict {
			t.Fatalf("reject twice: got %d %s, want 409", code, body)
		}
	})

	t.Run("revoke validates caller, reason, status and projects revoked user state", func(t *testing.T) {
		owner := newApplicant(t, "revoke")
		kycID := insertApproved(t, owner, "revoke")
		path := fmt.Sprintf("/kyc/verifications/%d/revoke", kycID)

		if code, body := doJSON(t, http.MethodPost, path, user.token, `{"reason":"Document frauduleux"}`); code != http.StatusForbidden {
			t.Fatalf("user revoke: got %d %s, want 403", code, body)
		}
		if code, body := doJSON(t, http.MethodPost, path, admin.token, `{}`); code != http.StatusBadRequest {
			t.Fatalf("missing reason: got %d %s, want 400", code, body)
		}
		if code, body := doJSON(t, http.MethodPost, path, admin.token, `{"reason":"   "}`); code != http.StatusBadRequest {
			t.Fatalf("blank reason: got %d %s, want 400", code, body)
		}
		if code, body := doJSON(t, http.MethodPost, "/kyc/verifications/999999999/revoke", admin.token, `{"reason":"Document frauduleux"}`); code != http.StatusNotFound {
			t.Fatalf("missing id: got %d %s, want 404", code, body)
		}

		submittedOwner := newApplicant(t, "revoke-submitted")
		submittedID := insertSubmitted(t, submittedOwner, "revoke-submitted")
		if code, body := doJSON(t, http.MethodPost, fmt.Sprintf("/kyc/verifications/%d/revoke", submittedID), admin.token, `{"reason":"Trop tôt"}`); code != http.StatusConflict {
			t.Fatalf("revoke submitted: got %d %s, want 409", code, body)
		}
		rejectedOwner := newApplicant(t, "revoke-rejected")
		rejectedID := insertRejected(t, rejectedOwner, "revoke-rejected")
		if code, body := doJSON(t, http.MethodPost, fmt.Sprintf("/kyc/verifications/%d/revoke", rejectedID), admin.token, `{"reason":"Déjà rejeté"}`); code != http.StatusConflict {
			t.Fatalf("revoke rejected: got %d %s, want 409", code, body)
		}

		code, body := doJSON(t, http.MethodPost, path, admin.token, `{"reason":"  Document frauduleux  "}`)
		if code != http.StatusOK {
			t.Fatalf("revoke: got %d %s, want 200", code, body)
		}
		var got map[string]any
		if err := json.Unmarshal([]byte(body), &got); err != nil {
			t.Fatal(err)
		}
		if got["status"] != "revoked" || got["revoked_by"] != float64(admin.id) || got["revocation_reason"] != "Document frauduleux" {
			t.Fatalf("revoke response missing revocation fields: %v", got)
		}

		var status string
		if err := db.QueryRow(`SELECT kyc_status FROM usr.users WHERE id = $1`, owner.id).Scan(&status); err != nil {
			t.Fatal(err)
		}
		if status != "revoked" {
			t.Fatalf("projection status = %q, want revoked", status)
		}
		if code, body := doJSON(t, http.MethodPost, path, admin.token, `{"reason":"Second revocation"}`); code != http.StatusConflict {
			t.Fatalf("revoke twice: got %d %s, want 409", code, body)
		}
	})
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}
