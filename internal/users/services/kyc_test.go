package services

import (
	"errors"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/users/database"
)

func TestCanAcceptNewKycFiling(t *testing.T) {
	now := time.Date(2026, 9, 28, 16, 0, 0, 0, time.UTC)
	future := now.Add(365 * 24 * time.Hour)
	past := now.Add(-24 * time.Hour)

	cases := []struct {
		name    string
		latest  *server.KycVerification
		wantErr error
	}{
		{
			name:    "no previous filing is accepted",
			latest:  nil,
			wantErr: nil,
		},
		{
			name:    "a filing awaiting a decision blocks a new one",
			latest:  &server.KycVerification{Status: database.StatusSubmitted},
			wantErr: ErrMessageKycAlreadyExists,
		},
		{
			name:    "a valid approval blocks a new filing",
			latest:  &server.KycVerification{Status: database.StatusApproved, ExpiresAt: &future},
			wantErr: ErrMessageKycStillValid,
		},
		{
			name:    "an expired approval lets the user file again",
			latest:  &server.KycVerification{Status: database.StatusApproved, ExpiresAt: &past},
			wantErr: nil,
		},
		{
			name:    "an approval expiring exactly now lets the user file again",
			latest:  &server.KycVerification{Status: database.StatusApproved, ExpiresAt: &now},
			wantErr: nil,
		},
		{
			name:    "an approval without expiry is treated as still valid",
			latest:  &server.KycVerification{Status: database.StatusApproved},
			wantErr: ErrMessageKycStillValid,
		},
		{
			name:    "a rejection lets the user file again",
			latest:  &server.KycVerification{Status: database.StatusRejected},
			wantErr: nil,
		},
		{
			name:    "a revocation lets the user file again",
			latest:  &server.KycVerification{Status: database.StatusRevoked, ExpiresAt: &future},
			wantErr: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := canAcceptNewKycFiling(tc.latest, now)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("canAcceptNewKycFiling() error = %v, want %v", err, tc.wantErr)
			}
		})
	}
}
