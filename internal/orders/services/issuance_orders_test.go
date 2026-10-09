package services

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/orders/database"
)

func TestGenerateIssuanceOrderReference(t *testing.T) {
	const alphabet = "23456789ABCDEFGHJKLMNPQRSTVWXYZ"
	format := regexp.MustCompile(`^ISS-\d{8}-[` + alphabet + `]{6}$`)

	t.Run("format and date", func(t *testing.T) {
		at := time.Date(2026, time.October, 5, 23, 59, 0, 0, time.UTC)
		ref, err := generateIssuanceOrderReference(at)
		if err != nil {
			t.Fatal(err)
		}
		if !format.MatchString(ref) {
			t.Errorf("reference %q does not match %s", ref, format)
		}
		if !strings.HasPrefix(ref, "ISS-20261005-") {
			t.Errorf("reference %q does not carry the order date", ref)
		}
	})

	t.Run("no ambiguous characters", func(t *testing.T) {
		for i := 0; i < 500; i++ {
			ref, err := generateIssuanceOrderReference(time.Now())
			if err != nil {
				t.Fatal(err)
			}
			if suffix := ref[len(ref)-6:]; strings.ContainsAny(suffix, "01OIU") {
				t.Fatalf("suffix %q contains an ambiguous character", suffix)
			}
		}
	})

	// 31^6 ≈ 887 M combinations: 100 draws collide with probability ~6e-6.
	t.Run("suffix is random", func(t *testing.T) {
		at := time.Date(2026, time.October, 5, 0, 0, 0, 0, time.UTC)
		seen := make(map[string]bool)
		for i := 0; i < 100; i++ {
			ref, err := generateIssuanceOrderReference(at)
			if err != nil {
				t.Fatal(err)
			}
			if seen[ref] {
				t.Fatalf("duplicate reference %q", ref)
			}
			seen[ref] = true
		}
	})
}

// TestReservationTTL covers the fallbacks around ORDER_RESERVATION_TTL (U1).
func TestReservationTTL(t *testing.T) {
	tests := []struct {
		name  string
		unset bool
		env   string
		want  time.Duration
	}{
		{name: "unset falls back to the default", unset: true, want: 15 * time.Minute},
		{name: "empty falls back to the default", env: "", want: 15 * time.Minute},
		{name: "minutes", env: "30m", want: 30 * time.Minute},
		{name: "seconds", env: "90s", want: 90 * time.Second},
		{name: "composite duration", env: "1h30m", want: 90 * time.Minute},
		{name: "surrounding spaces are tolerated", env: " 10m ", want: 10 * time.Minute},
		{name: "a unitless number falls back", env: "900", want: 15 * time.Minute},
		{name: "plain words fall back", env: "quinze minutes", want: 15 * time.Minute},
		{name: "zero falls back", env: "0s", want: 15 * time.Minute},
		{name: "negative falls back", env: "-5m", want: 15 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ORDER_RESERVATION_TTL", tt.env)
			if tt.unset {
				os.Unsetenv("ORDER_RESERVATION_TTL")
			}

			if got := reservationTTL(); got != tt.want {
				t.Fatalf("reservationTTL() = %s, want %s", got, tt.want)
			}
		})
	}
}

// TestPlacementRefusalsReachTheHandler pins the contract the handler relies
// on: a refusal wrapped by the database, with its detail, still matches the
// service sentinel the handler switches on.
func TestPlacementRefusalsReachTheHandler(t *testing.T) {
	tests := []struct {
		db      error
		service error
	}{
		{database.ErrOrderAssetNotFound, ErrRealEstateNotFound},
		{database.ErrOrderAssetNotOpen, ErrRealEstateNotAvailable},
		{database.ErrOrderNotEnoughShares, ErrNotEnoughShares},
		{database.ErrOrderStakeLimit, ErrStakeLimit},
		{database.ErrOrderInvestorNotEligible, ErrInvestorNotEligible},
		{database.ErrOrderIdempotencyKeyReuse, ErrIdempotencyKeyReused},
	}

	for _, tt := range tests {
		t.Run(tt.db.Error(), func(t *testing.T) {
			wrapped := fmt.Errorf("%w: only 5 shares left, 10 requested", tt.db)
			if !errors.Is(wrapped, tt.service) {
				t.Fatalf("errors.Is(%q, %q) = false", wrapped, tt.service)
			}
		})
	}
}
