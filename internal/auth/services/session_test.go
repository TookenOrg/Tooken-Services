package services

import (
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
)

func TestProjectSessionUserKycStatus(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)
	past := now.Add(-time.Minute)
	future := now.Add(time.Hour)

	approved := server.UserKycStatusApproved
	verified := server.UserKycStatusVerified

	expiredUser := server.User{KycStatus: &approved, KycExpiresAt: &past}
	projectSessionUserKycStatus(&expiredUser, now)
	if expiredUser.KycStatus == nil || *expiredUser.KycStatus != server.UserKycStatusExpired {
		t.Fatalf("expired user kyc_status = %v, want %q", expiredUser.KycStatus, server.UserKycStatusExpired)
	}

	activeUser := server.User{KycStatus: &verified, KycExpiresAt: &future}
	projectSessionUserKycStatus(&activeUser, now)
	if activeUser.KycStatus == nil || *activeUser.KycStatus != server.UserKycStatusVerified {
		t.Fatalf("active user kyc_status = %v, want %q", activeUser.KycStatus, server.UserKycStatusVerified)
	}

	noKycUser := server.User{}
	projectSessionUserKycStatus(&noKycUser, now)
	if noKycUser.KycStatus != nil {
		t.Fatalf("nil kyc_status became %v, want nil", noKycUser.KycStatus)
	}
}
