//go:build integration

// HTTP coverage of /users/kyc/verifications: an investor filing his KYC
// declaration (POST), and the back office listing the filings (GET).
//
// The decision routes on /users/kyc/verifications/{id} are covered in
// kyc_verification_integration_test.go; the on-chain synchronisation of an
// approval needs a chain and is covered with one elsewhere.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	_ "github.com/lib/pq"
)

const kycPath = "/users/kyc/verifications"

// kycFilings counts the filings of a user: a refused filing must leave it
// unchanged.
func kycFilings(t *testing.T, h *orderHarness, user int) int {
	t.Helper()
	return h.count(t, `SELECT count(*) FROM usr.kyc_verification WHERE user_id = $1`, user)
}

// seedKycFiling writes a filing in the given status straight in SQL, decided by
// decider, with every column its status demands. expiresIn is a PostgreSQL
// interval added to now() for expires_at, used only by an approval.
func seedKycFiling(t *testing.T, h *orderHarness, user, decider int, status, expiresIn string) int {
	t.Helper()
	var id int
	if err := h.db.QueryRow(`
	    INSERT INTO usr.kyc_verification (
	        user_id, status, declared_full_name, declared_country_code, submitted_at,
	        decided_at, decided_by, rejection_reason, expires_at,
	        revoked_at, revoked_by, revocation_reason)
	    SELECT $1, $2::text, 'Jeanne Martin', 'FR', now() - interval '60 days',
	        CASE WHEN $2::text <> 'submitted' THEN now() - interval '59 days' END,
	        CASE WHEN $2::text <> 'submitted' THEN $3::int END,
	        CASE WHEN $2::text = 'rejected' THEN 'Document illisible' END,
	        CASE WHEN $2::text IN ('approved', 'revoked') THEN now() + $4::interval END,
	        CASE WHEN $2::text = 'revoked' THEN now() - interval '1 day' END,
	        CASE WHEN $2::text = 'revoked' THEN $3::int END,
	        CASE WHEN $2::text = 'revoked' THEN 'Document frauduleux' END
	    RETURNING id`, user, status, decider, expiresIn).Scan(&id); err != nil {
		t.Fatalf("seeding a %s filing: %v", status, err)
	}
	return id
}

func kycDeclaration(fullName, country string) string {
	return fmt.Sprintf(`{"declared_full_name": %q, "declared_country_code": %q}`, fullName, country)
}

func TestSubmitKycVerification(t *testing.T) {
	h := newOrderHarness(t)
	admin := newTrackedUser(t, h, authUtils.RoleAdmin)

	submit := func(token, body string) (int, string) {
		return callAPI(h, http.MethodPost, kycPath, token, body)
	}

	t.Run("an anonymous caller is refused with 401", func(t *testing.T) {
		code, body := submit("", kycDeclaration("Jeanne Martin", "FR"))
		expectCode(t, code, body, http.StatusUnauthorized)
		expectAPIMessage(t, body, "Authorization header required")
	})

	t.Run("a first filing is stored as submitted, trimmed and upper-cased", func(t *testing.T) {
		investor := newTrackedUser(t, h, authUtils.RoleUser)

		code, body := submit(investor.token, kycDeclaration("  Jeanne Martin  ", " lu "))
		expectCode(t, code, body, http.StatusCreated)
		expectAPIMessage(t, body, "KYC verification submitted successfully.")

		var status, fullName, country, projected string
		if err := h.db.QueryRow(`
		    SELECT k.status, k.declared_full_name, k.declared_country_code, u.kyc_status
		    FROM usr.kyc_verification k JOIN usr.users u ON u.id = k.user_id
		    WHERE k.user_id = $1`, investor.id).Scan(&status, &fullName, &country, &projected); err != nil {
			t.Fatalf("reading the filing: %v", err)
		}
		if status != "submitted" || fullName != "Jeanne Martin" || country != "LU" {
			t.Fatalf("stored %q %q %q, want submitted / Jeanne Martin / LU", status, fullName, country)
		}
		// The user row is a projection of his latest filing: a pending review
		// is what the profile must show from now on.
		if projected != "pending" {
			t.Fatalf("usr.users.kyc_status = %q, want pending", projected)
		}
	})

	t.Run("an invalid declaration is refused with 400 and writes nothing", func(t *testing.T) {
		investor := newTrackedUser(t, h, authUtils.RoleUser)

		for _, c := range []struct {
			name   string
			body   string
			reason string
		}{
			{"a blank full name", kycDeclaration("   ", "FR"), "invalid full name"},
			{"no field at all", `{}`, "invalid full name"},
			{"an empty country code", kycDeclaration("Jeanne Martin", ""), "invalid country code"},
			{"a three-letter country code", kycDeclaration("Jeanne Martin", "FRA"), "invalid country code"},
			{"a country code with a digit", kycDeclaration("Jeanne Martin", "F1"), "invalid country code"},
			{"a body that is not JSON", `{"declared_full_name": "Jeanne`, "unexpected EOF"},
		} {
			t.Run(c.name, func(t *testing.T) {
				code, body := submit(investor.token, c.body)
				expectCode(t, code, body, http.StatusBadRequest)
				expectAPIMessage(t, body, c.reason)
				if n := kycFilings(t, h, investor.id); n != 0 {
					t.Fatalf("%d filing(s) written by a refused declaration", n)
				}
			})
		}
	})

	t.Run("a second filing is refused while the first awaits a decision", func(t *testing.T) {
		investor := newTrackedUser(t, h, authUtils.RoleUser)
		code, body := submit(investor.token, kycDeclaration("Jeanne Martin", "FR"))
		expectCode(t, code, body, http.StatusCreated)

		code, body = submit(investor.token, kycDeclaration("Jeanne Martin-Durand", "FR"))
		expectCode(t, code, body, http.StatusConflict)
		expectAPIMessage(t, body, "status 'submitted' already exists")
		if n := kycFilings(t, h, investor.id); n != 1 {
			t.Fatalf("%d filing(s), want the first one only", n)
		}
	})

	t.Run("a filing is refused while an approval is still valid", func(t *testing.T) {
		investor := newTrackedUser(t, h, authUtils.RoleUser)
		seedKycFiling(t, h, investor.id, admin.id, "approved", "30 days")

		code, body := submit(investor.token, kycDeclaration("Jeanne Martin", "FR"))
		expectCode(t, code, body, http.StatusConflict)
		expectAPIMessage(t, body, "a valid KYC verification already exists")
		if n := kycFilings(t, h, investor.id); n != 1 {
			t.Fatalf("%d filing(s), want the approval only", n)
		}
	})

	// A single refusal must never lock an investor out for good.
	t.Run("filing again is allowed once the previous filing no longer stands", func(t *testing.T) {
		for _, previous := range []struct {
			name      string
			status    string
			expiresIn string
		}{
			{"after a rejection", "rejected", "0"},
			{"after a revocation", "revoked", "30 days"},
			{"after the approval expired", "approved", "-1 day"},
		} {
			t.Run(previous.name, func(t *testing.T) {
				investor := newTrackedUser(t, h, authUtils.RoleUser)
				seedKycFiling(t, h, investor.id, admin.id, previous.status, previous.expiresIn)

				code, body := submit(investor.token, kycDeclaration("Jeanne Martin", "FR"))
				expectCode(t, code, body, http.StatusCreated)
				if n := h.count(t, `
				    SELECT count(*) FROM usr.kyc_verification
				    WHERE user_id = $1 AND status = 'submitted'`, investor.id); n != 1 {
					t.Fatalf("%d submitted filing(s), want 1", n)
				}
			})
		}
	})

	t.Run("a token whose account no longer exists files nothing", func(t *testing.T) {
		ghost := h.newUser(t, authUtils.RoleUser, "none", "NULL")
		if _, err := h.db.Exec(`DELETE FROM usr.users WHERE id = $1`, ghost.id); err != nil {
			t.Fatal(err)
		}

		// The token is still valid but points at nobody: the client must be
		// told to sign in again, not that the server failed.
		code, body := submit(ghost.token, kycDeclaration("Jeanne Martin", "FR"))
		expectCode(t, code, body, http.StatusUnauthorized)
		expectAPIMessage(t, body, "Sign in again")
		if n := kycFilings(t, h, ghost.id); n != 0 {
			t.Fatalf("%d filing(s) written for a deleted account", n)
		}
	})
}

func TestListKycVerifications(t *testing.T) {
	h := newOrderHarness(t)
	admin := newTrackedUser(t, h, authUtils.RoleAdmin)
	manager := newTrackedUser(t, h, authUtils.RoleManager)
	user := newTrackedUser(t, h, authUtils.RoleUser)

	// Jeanne was rejected once, then filed again; Paul holds an approval.
	jeanne := newTrackedUser(t, h, authUtils.RoleUser)
	jeanneRejected := seedKycFiling(t, h, jeanne.id, admin.id, "rejected", "0")
	jeanneSubmitted := seedKycFiling(t, h, jeanne.id, admin.id, "submitted", "0")
	paul := newTrackedUser(t, h, authUtils.RoleUser)
	paulApproved := seedKycFiling(t, h, paul.id, admin.id, "approved", "30 days")

	var jeanneEmail string
	if err := h.db.QueryRow(`SELECT email FROM usr.users WHERE id = $1`, jeanne.id).Scan(&jeanneEmail); err != nil {
		t.Fatal(err)
	}

	list := func(t *testing.T, token, query string) []server.KycVerification {
		t.Helper()
		code, body := callAPI(h, http.MethodGet, kycPath+query, token, "")
		expectCode(t, code, body, http.StatusOK)
		var filings []server.KycVerification
		if err := json.Unmarshal([]byte(body), &filings); err != nil {
			t.Fatalf("decoding %s: %v", body, err)
		}
		if filings == nil {
			t.Fatalf("want a JSON array, got %s", body)
		}
		return filings
	}
	ids := func(filings []server.KycVerification) map[int]server.KycVerification {
		byID := map[int]server.KycVerification{}
		for _, f := range filings {
			byID[f.Id] = f
		}
		return byID
	}

	t.Run("an anonymous caller is refused with 401", func(t *testing.T) {
		code, body := callAPI(h, http.MethodGet, kycPath, "", "")
		expectCode(t, code, body, http.StatusUnauthorized)
		expectAPIMessage(t, body, "Authorization header required")
	})

	t.Run("only an ADMIN may list the filings", func(t *testing.T) {
		for role, token := range map[string]string{"USER": user.token, "MANAGER": manager.token} {
			code, body := callAPI(h, http.MethodGet, fmt.Sprintf("%s?userId=%d", kycPath, jeanne.id), token, "")
			if code != http.StatusForbidden {
				t.Fatalf("%s: want 403, got %d: %s", role, code, body)
			}
			expectAPIMessage(t, body, "not allowed")
			if strings.Contains(body, jeanneEmail) || strings.Contains(body, "Jeanne Martin") {
				t.Fatalf("%s: a refused listing leaked filings: %s", role, body)
			}
		}
	})

	t.Run("filtering on a user returns all of his filings and only his", func(t *testing.T) {
		got := list(t, admin.token, fmt.Sprintf("?userId=%d", jeanne.id))
		byID := ids(got)
		if len(got) != 2 || byID[jeanneRejected].Status != "rejected" || byID[jeanneSubmitted].Status != "submitted" {
			t.Fatalf("got %+v, want Jeanne's rejected and submitted filings", got)
		}
		for _, f := range got {
			if f.UserId != jeanne.id || f.Email != jeanneEmail {
				t.Fatalf("filing %d belongs to user %d (%s), want only Jeanne's", f.Id, f.UserId, f.Email)
			}
		}
		if byID[jeanneRejected].RejectionReason == nil || *byID[jeanneRejected].RejectionReason != "Document illisible" {
			t.Fatalf("the rejection reason is missing from the listing: %+v", byID[jeanneRejected])
		}
	})

	t.Run("filtering on a user and a status combines both", func(t *testing.T) {
		got := list(t, admin.token, fmt.Sprintf("?userId=%d&status=submitted", jeanne.id))
		if len(got) != 1 || got[0].Id != jeanneSubmitted {
			t.Fatalf("got %+v, want only filing %d", got, jeanneSubmitted)
		}
		if got := list(t, admin.token, fmt.Sprintf("?userId=%d&status=approved", jeanne.id)); len(got) != 0 {
			t.Fatalf("got %+v, want no approved filing for Jeanne", got)
		}
	})

	t.Run("filtering on a status alone returns that status across users", func(t *testing.T) {
		got := list(t, admin.token, "?status=approved")
		if _, ok := ids(got)[paulApproved]; !ok {
			t.Fatalf("Paul's approval %d missing from %+v", paulApproved, got)
		}
		for _, f := range got {
			if f.Status != "approved" {
				t.Fatalf("filing %d is %q in a listing filtered on approved", f.Id, f.Status)
			}
		}
	})

	t.Run("without filter every filing is listed", func(t *testing.T) {
		byID := ids(list(t, admin.token, ""))
		for _, id := range []int{jeanneRejected, jeanneSubmitted, paulApproved} {
			if _, ok := byID[id]; !ok {
				t.Fatalf("filing %d missing from the unfiltered listing", id)
			}
		}
	})

	t.Run("a user who never filed gets an empty array, not null", func(t *testing.T) {
		code, body := callAPI(h, http.MethodGet, fmt.Sprintf("%s?userId=%d", kycPath, user.id), admin.token, "")
		expectCode(t, code, body, http.StatusOK)
		if body != "[]" {
			t.Fatalf("want [], got %s", body)
		}
	})

	// The generated wrapper does not check the enum of the filter: an unknown
	// status reached the query, matched nothing, and answered an empty list a
	// client could not tell from "nobody is in that state".
	t.Run("a status outside the filter values is refused with 400", func(t *testing.T) {
		for _, status := range []string{"bogus", "SUBMITTED", ""} {
			code, body := callAPI(h, http.MethodGet, kycPath+"?status="+status, admin.token, "")
			expectCode(t, code, body, http.StatusBadRequest)
			expectAPIMessage(t, body, "submitted", "approved")
		}
	})

	t.Run("a database failure answers 500, not an empty list", func(t *testing.T) {
		withClosedDatabase(t, func() {
			code, body := callAPI(h, http.MethodGet, kycPath, admin.token, "")
			expectCode(t, code, body, http.StatusInternalServerError)
			expectGenericFailure(t, body)
		})
	})
}

// Every KYC decision route is for ADMIN only. A refusal must say so in the
// {"message": ...} envelope: it used to serialise a Go error, which encodes
// as {}, and the client got a 403 with nothing to show.
func TestKycDecisionRoutesExplainTheirRefusal(t *testing.T) {
	h := newOrderHarness(t)
	manager := newTrackedUser(t, h, authUtils.RoleManager)
	id := fmt.Sprintf("%s/%d", kycPath, unknownID)

	for _, route := range []struct{ method, path, body string }{
		{http.MethodGet, kycPath, ""},
		{http.MethodGet, id, ""},
		{http.MethodPost, id + "/approve", `{"expires_at":"2099-01-01T00:00:00Z"}`},
		{http.MethodPost, id + "/reject", `{"reason":"Document illisible"}`},
		{http.MethodPost, id + "/revoke", `{"reason":"Document frauduleux"}`},
		{http.MethodPost, id + "/sync-onchain", ""},
	} {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			code, body := callAPI(h, route.method, route.path, manager.token, route.body)
			expectCode(t, code, body, http.StatusForbidden)
			expectAPIMessage(t, body, "not allowed")
		})
	}
}
