//go:build integration

// Coverage of the identity idempotency guard against a real PostgreSQL.
//
// This guard broke twice in one evening, in opposite directions, and neither
// break was visible to the compiler or to any existing test:
//
//   - first, isIdempotentIdentity returned "a wallet exists" where it meant "the
//     way is clear". Paired with a mocked GetWalletByUserId that never returned
//     nil, the two mistakes cancelled out: the guard let everything through and
//     nothing looked wrong;
//   - then the comparison was fixed while the mock was still in place, and the
//     endpoint started refusing every single user — including, and especially,
//     the new investors it exists for.
//
// The lesson is not about the comparison. A guard has two sides, and a test that
// only exercises the side that refuses cannot tell a correct guard from one that
// always says no. Both sides are asserted below, and each fails on its own when
// the guard is broken in the matching direction.
//
// No chain is needed here, only a database:
//
//	TEST_DATABASE_URL="postgres://…" go test -tags integration ./internal/blockchain/services/...
package services

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	_ "github.com/lib/pq"
)

const (
	idempotencyProbeUser  = 900011
	idempotencyProbeLabel = "identity-idempotency-probe"
)

func TestIsIdempotentIdentity(t *testing.T) {
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

	clean := func() {
		if _, err := db.Exec(
			`DELETE FROM blk.user_wallet WHERE label = $1`, idempotencyProbeLabel); err != nil {
			t.Fatal(err)
		}
	}
	clean()
	t.Cleanup(func() {
		clean()
		db.Close()
	})

	ctx := context.Background()
	req := server.CreateIdentityRequest{UserId: idempotencyProbeUser, CountryCode: 442}

	// The side that lets through. Remove it and a guard hard-coded to refuse
	// would still look perfectly healthy — which is exactly what happened.
	t.Run("an investor with no wallet may get an identity", func(t *testing.T) {
		if !isIdempotentIdentity(ctx, req) {
			t.Fatal("a user without a wallet must be allowed to create an identity")
		}
	})

	// The side that bites. This is the whole point of the guard: a second
	// identity would leave the investor with two wallets, and the registry with
	// two answers to "where do we deliver".
	t.Run("an investor who already has a wallet is refused", func(t *testing.T) {
		if _, err := db.Exec(`
			INSERT INTO blk.user_wallet (user_id, wallet_address, label)
			VALUES ($1, $2, $3)`,
			idempotencyProbeUser,
			"0x5555555555555555555555555555555555555555",
			idempotencyProbeLabel); err != nil {
			t.Fatal(err)
		}

		if isIdempotentIdentity(ctx, req) {
			t.Fatal("a user who already has a wallet must be refused")
		}
	})

	// A wallet that was retired must not keep its owner out. The guard reads
	// through GetWalletByUserId, which filters on is_active; this test fails if
	// that filter is ever dropped from under it.
	t.Run("a retired wallet does not block a new identity", func(t *testing.T) {
		if _, err := db.Exec(
			`UPDATE blk.user_wallet SET is_active = false WHERE label = $1`,
			idempotencyProbeLabel); err != nil {
			t.Fatal(err)
		}

		if !isIdempotentIdentity(ctx, req) {
			t.Fatal("a deactivated wallet must not prevent a new identity")
		}
	})
}
