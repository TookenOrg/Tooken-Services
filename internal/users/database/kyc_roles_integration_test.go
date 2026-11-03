//go:build integration

package database

import (
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/testenv"
	"github.com/TookenOrg/tooken-services/pkg/logger"
)

func TestKycDirectoryReadsCurrentAccountRole(t *testing.T) {
	db := testenv.ScratchDatabase(t, "tooken_it_kyc_roles")
	oldDB := globals.DB
	globals.DB = db
	t.Cleanup(func() { globals.DB = oldDB })
	logger.Init(false)

	var userID int
	if err := db.QueryRow(`INSERT INTO usr.users (full_name, email, password, role)
		VALUES ('Account', 'account@example.test', 'x', 'USER') RETURNING id`).Scan(&userID); err != nil {
		t.Fatal(err)
	}
	filingID, _, err := InsertKycVerification(t.Context(), userID, &server.KycVerificationRequest{
		DeclaredFullName: "Account", DeclaredCountryCode: "LU",
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, role := range []string{"USER", "MANAGER", "ADMIN"} {
		if _, err := db.Exec(`UPDATE usr.users SET role = $1 WHERE id = $2`, role, userID); err != nil {
			t.Fatal(err)
		}
		rows, err := GetKycVerificationByUserIDAndStatus(t.Context(), &userID, nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(rows) != 1 {
			t.Fatalf("filings = %d, want 1", len(rows))
		}
		detail, err := GetKycVerificationByID(t.Context(), filingID)
		if err != nil {
			t.Fatal(err)
		}
		latest, err := GetLatestKycVerificationByUserID(t.Context(), userID)
		if err != nil {
			t.Fatal(err)
		}
		for _, filing := range []server.KycVerification{rows[0], detail, latest} {
			if filing.UserRole == nil || string(*filing.UserRole) != role {
				t.Fatalf("filing %d does not expose current role %s: %v", filing.Id, role, filing.UserRole)
			}
		}

		for _, mutate := range []func() (server.KycVerification, error){
			func() (server.KycVerification, error) {
				return ApproveKycVerification(t.Context(), filingID, time.Now().Add(time.Hour), userID)
			},
			func() (server.KycVerification, error) {
				return RevokeKycVerification(t.Context(), filingID, "Review required", userID)
			},
			func() (server.KycVerification, error) {
				return RejectKycVerification(t.Context(), filingID, "Review incomplete", userID)
			},
		} {
			filing, err := mutate()
			if err != nil {
				t.Fatal(err)
			}
			if filing.UserRole == nil {
				t.Fatalf("updated filing does not expose current account role: %v", filing.UserRole)
			}
			if string(*filing.UserRole) != role {
				t.Fatalf("updated filing %s role = %q, want %s", filing.Status, *filing.UserRole, role)
			}
		}
	}
}
