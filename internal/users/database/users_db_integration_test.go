//go:build integration

package database

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/globals"
	_ "github.com/lib/pq"
)

func TestMarkUserKycVerified(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	previousDB := globals.DB
	globals.DB = db
	t.Cleanup(func() {
		globals.DB = previousDB
		db.Close()
	})

	email := fmt.Sprintf("mark-kyc-verified-%d@tooken.test", time.Now().UnixNano())
	var userID int
	if err := db.QueryRow(`
	    INSERT INTO usr.users (full_name, email, password, role, kyc_status)
	    VALUES ('Mark KYC Verified', $1, 'x', 'USER', 'approved')
	    RETURNING id`, email).Scan(&userID); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = db.Exec(`DELETE FROM usr.users WHERE email = $1`, email)
	})

	verifiedAt, err := MarkUserKycVerified(t.Context(), userID)
	if err != nil {
		t.Fatalf("MarkUserKycVerified: %v", err)
	}
	if verifiedAt.IsZero() {
		t.Fatal("MarkUserKycVerified returned a zero verifiedAt")
	}

	var status string
	var storedVerifiedAt sql.NullTime
	if err := db.QueryRow(`
	    SELECT kyc_status, kyc_verified_at
	    FROM usr.users
	    WHERE id = $1`, userID).Scan(&status, &storedVerifiedAt); err != nil {
		t.Fatal(err)
	}
	if status != StatusVerified {
		t.Fatalf("kyc_status = %q, want %q", status, StatusVerified)
	}
	if !storedVerifiedAt.Valid {
		t.Fatal("kyc_verified_at is NULL")
	}
	if !storedVerifiedAt.Time.Equal(verifiedAt) {
		t.Fatalf("stored verifiedAt = %v, returned %v", storedVerifiedAt.Time, verifiedAt)
	}
}
