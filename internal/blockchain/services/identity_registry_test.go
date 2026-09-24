package services

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/lib/pq"
)

// TestRegisterIdentityCountryGuard covers the bounds check that stands between an
// application int and the uint16 a country is stored as on-chain.
//
// Zero is a *valid* country in the IdentityRegistry, so a silent truncation does not
// fail: it registers a wrong nationality, permanently and invisibly. 70000 truncates
// to 4464, which is a perfectly plausible-looking country code.
//
// The guard is asserted on BOTH sides. A test that only exercises the refusing side
// cannot tell a correct guard from one hard-wired to reject everything — the lesson
// already paid for by isIdempotentIdentity (see identity_idempotency_integration_test.go).
//
// No chain and no database are needed: the rejecting side returns before touching
// either, and the accepting side is recognised by the fact that it fails *further
// down*, on a deliberately unreachable database, rather than on the country code.
func TestRegisterIdentityCountryGuard(t *testing.T) {
	logger.Init(true)

	// A handle that opens lazily and can never connect: reaching it proves the guard
	// let the call through, without requiring a real database.
	unreachable, err := sql.Open("postgres", "postgres://nobody@127.0.0.1:1/nowhere?sslmode=disable&connect_timeout=1")
	if err != nil {
		t.Fatalf("open unreachable handle: %v", err)
	}
	previousDB := globals.DB
	globals.DB = unreachable
	t.Cleanup(func() {
		globals.DB = previousDB
		unreachable.Close()
	})

	ctx := context.Background()
	wallet := common.HexToAddress("0x1111111111111111111111111111111111111111")
	identity := common.HexToAddress("0x2222222222222222222222222222222222222222")

	// ---- the side that refuses -------------------------------------------
	for _, tc := range []struct {
		name    string
		country int
	}{
		{"above uint16", 70000},
		{"just above uint16", 65536},
		{"negative", -1},
	} {
		t.Run("rejected: "+tc.name, func(t *testing.T) {
			_, err := registerIdentity(ctx, wallet, identity, tc.country)
			if err == nil {
				t.Fatalf("country %d must be rejected, got no error", tc.country)
			}
			if !strings.Contains(err.Error(), "country code") {
				t.Fatalf("country %d must be rejected for its value, got: %v", tc.country, err)
			}
		})
	}

	// ---- the side that lets through --------------------------------------
	for _, tc := range []struct {
		name    string
		country int
	}{
		{"France", 250},
		{"upper bound", 65535},
		{"zero is a valid country", 0},
	} {
		t.Run("accepted: "+tc.name, func(t *testing.T) {
			_, err := registerIdentity(ctx, wallet, identity, tc.country)
			// It must still fail — there is no database behind this handle — but it
			// must fail on the storage lookup, never on the country code.
			if err == nil {
				t.Fatalf("country %d reached a working backend, which this test does not provide", tc.country)
			}
			if strings.Contains(err.Error(), "country code") {
				t.Fatalf("country %d must be accepted, got: %v", tc.country, err)
			}
		})
	}
}
