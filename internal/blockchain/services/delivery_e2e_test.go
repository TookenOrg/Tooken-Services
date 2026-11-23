//go:build integration

// End-to-end delivery on a real chain and a real database.
//
// These tests live in the blockchain package because that is where the T-REX
// fixture lives — deploying a suite, registering an identity and signing a KYC
// claim is several hundred lines that cannot be reached from outside. The
// import therefore goes blockchain(test) -> delivery, which is the opposite of
// the production boundary: delivery never imports blockchain, and the
// dependency check still measures 0 both ways, since `go list -deps` ignores
// test files.
// Here the test plays the part of the handler, which is the one place allowed
// to know both sides.
//
// Requires a node and a local PostgreSQL whose user may CREATE DATABASE:
//
//	cd tools/hardhat && npx hardhat node
//	TEST_DATABASE_URL="postgres://postgres:pw@127.0.0.1:55444/tooken?sslmode=disable" \
//	  go test -tags integration -run TestDelivery ./internal/blockchain/services/...
package services

import (
	"context"
	"fmt"
	"math/big"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	chainglobals "github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	deliveryServices "github.com/TookenOrg/tooken-services/internal/delivery/services"
	"github.com/ethereum/go-ethereum/common"
)

// deliveryScene is one investor who paid for shares of one tokenised asset:
// everything the delivery reads, on the chain and in the database.
type deliveryScene struct {
	env      *chainDBEnv
	suite    *deployedSuite
	orderRef string
	wallet   common.Address
	shares   int
}

// seedAsset creates the issuer and the tokenised asset. ass.real_estate is
// UNIQUE on token_id, so one token backs one asset: two investors of the same
// fundraising share it, which is what a real one looks like.
func seedAsset(t *testing.T, env *chainDBEnv, suite *deployedSuite) int {
	t.Helper()
	return seedAssetOffering(t, env, suite, 1000)
}

// seedAssetOffering is seedAsset for an asset that offers totalShares shares.
// It is open for subscription (fundraising), like any asset that has orders.
func seedAssetOffering(t *testing.T, env *chainDBEnv, suite *deployedSuite, totalShares int) int {
	t.Helper()

	unique := time.Now().UnixNano()

	var tokenID int64
	if err := env.DB.QueryRow(
		`SELECT id FROM blk.token WHERE address = $1`, suite.TokenAddr.Hex()).Scan(&tokenID); err != nil {
		t.Fatalf("read token row: %v", err)
	}

	var issuerID int
	if err := env.DB.QueryRow(`
		INSERT INTO ass.issuer (name, legal_form) VALUES ($1, 'SA') RETURNING id`,
		fmt.Sprintf("Issuer %d", unique)).Scan(&issuerID); err != nil {
		t.Fatalf("insert issuer: %v", err)
	}

	var assetID int
	if err := env.DB.QueryRow(`
		INSERT INTO ass.real_estate (title, description, estate_type, issuer_id, token_id, status_id)
		VALUES ($1, 'e2e asset', 1, $2, $3, 4) RETURNING id`,
		fmt.Sprintf("Tooken Villa %d", unique), issuerID, tokenID,
	).Scan(&assetID); err != nil {
		t.Fatalf("insert real estate: %v", err)
	}

	if _, err := env.DB.Exec(`
		INSERT INTO ass.real_estate_shares_config (real_estate_id, total_shares, price_per_share)
		VALUES ($1, $2, 204.00)`, assetID, totalShares); err != nil {
		t.Fatalf("insert shares config: %v", err)
	}

	return assetID
}

// seedDeliveryOrder builds the investor, the wallet and the order on an asset
// that already exists, and optionally takes the investor through the KYC
// journey on chain.
func seedDeliveryOrder(ctx context.Context, t *testing.T, env *chainDBEnv,
	suite *deployedSuite, assetID, shares int, verifiedOnChain bool) deliveryScene {

	t.Helper()

	wallet, identityAddr, identity := env.Fixture.NewInvestorIdentity(t)
	if verifiedOnChain {
		if _, err := registerIdentity(ctx, wallet, identityAddr, 250); err != nil {
			t.Fatalf("registerIdentity: %v", err)
		}
		env.Fixture.AddKYCClaim(t, identityAddr, identity)
	}

	unique := time.Now().UnixNano() + int64(shares)

	var userID int
	if err := env.DB.QueryRow(`
		INSERT INTO usr.users (full_name, email, password, kyc_status, kyc_verified_at)
		VALUES ($1, $2, 'not-a-real-hash', 'verified', now() - interval '1 hour')
		RETURNING id`,
		"Alice Investor", fmt.Sprintf("alice-%d@test.local", unique),
	).Scan(&userID); err != nil {
		t.Fatalf("insert user: %v", err)
	}

	if _, err := env.DB.Exec(`
		INSERT INTO blk.user_wallet (user_id, wallet_address, label, is_active)
		VALUES ($1, $2, 'test', true)`, userID, wallet.Hex()); err != nil {
		t.Fatalf("insert wallet: %v", err)
	}

	orderRef := fmt.Sprintf("ISS-E2E-%d", unique)
	if _, err := env.DB.Exec(`
		INSERT INTO iss.issuance_orders
			(user_id, asset_id, quantity, status_id, order_reference,
			 unit_price, currency_code, gross_amount, fee_amount, amount_due)
		VALUES ($1, $2, $3, $4, $5, 204.00, 'EUR', $6, 0.00, $6)`,
		userID, assetID, shares, int(server.IssuanceOrderStatusIdPaymentConfirmed),
		orderRef, fmt.Sprintf("%.2f", float64(shares)*204.00),
	); err != nil {
		t.Fatalf("insert order: %v", err)
	}

	return deliveryScene{env: env, suite: suite, orderRef: orderRef, wallet: wallet, shares: shares}
}

func (s deliveryScene) status(t *testing.T) (int, string) {
	t.Helper()

	var statusID int
	var hash *string
	if err := s.env.DB.QueryRow(
		`SELECT status_id, delivery_tx_hash FROM iss.issuance_orders WHERE order_reference = $1`,
		s.orderRef).Scan(&statusID, &hash); err != nil {
		t.Fatalf("read order %s: %v", s.orderRef, err)
	}
	if hash == nil {
		return statusID, ""
	}
	return statusID, *hash
}

func (s deliveryScene) lastReason(t *testing.T) string {
	t.Helper()

	var reason *string
	if err := s.env.DB.QueryRow(`
		SELECT h.reason
		FROM iss.issuance_order_status_history h
		JOIN iss.issuance_orders o ON o.id = h.order_id
		WHERE o.order_reference = $1
		ORDER BY h.created_at DESC, h.id DESC
		LIMIT 1`, s.orderRef).Scan(&reason); err != nil {
		t.Fatalf("read timeline of %s: %v", s.orderRef, err)
	}
	if reason == nil {
		return ""
	}
	return *reason
}

// platformNonce reads the nonce of the signing key. It is the sharpest proof
// that nothing was broadcast: a nonce that did not move means no transaction
// left, whatever the logs say.
func platformNonce(ctx context.Context, t *testing.T, from common.Address) uint64 {
	t.Helper()

	nonce, err := chainglobals.EthClient.PendingNonceAt(ctx, from)
	if err != nil {
		t.Fatalf("PendingNonceAt: %v", err)
	}
	return nonce
}

func shares18(n int) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(n)), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}

// The whole journey. A paid order ends up CLOSED, and the investor holds
// the shares on chain. This is the only test that proves the three layers —
// database, service, chain — agree on what a delivered order is.
func TestDeliveryE2ENominal(t *testing.T) {
	ctx := context.Background()
	env := newChainDBEnv(ctx, t)
	chainSvc := NewService()

	suite := env.Fixture.DeploySuite(t, "Delivery", "Tooken Delivery", "TKDL", common.Address{})
	env.SetSharedIRSRow(t, suite.IRSAddr)
	env.InsertTokenRow(t, suite.TokenAddr, "Tooken Delivery", "TKDL", 18)

	scene := seedDeliveryOrder(ctx, t, env, suite, seedAsset(t, env, suite), 10, true)

	deliveryServices.NewService(chainSvc).Deliver(ctx, scene.orderRef)

	status, hash := scene.status(t)
	if status != int(server.IssuanceOrderStatusIdClosed) {
		t.Fatalf("order is %d, want CLOSED (%d). Last reason: %s",
			status, int(server.IssuanceOrderStatusIdClosed), scene.lastReason(t))
	}
	if hash == "" {
		t.Fatal("a closed order must carry the mint that backs it")
	}

	balance, err := suite.Token.BalanceOf(env.Fixture.Call, scene.wallet)
	if err != nil {
		t.Fatalf("BalanceOf: %v", err)
	}
	if balance.Cmp(shares18(10)) != 0 {
		t.Fatalf("balance = %s, want %s", balance, shares18(10))
	}
}

func (s deliveryScene) assetStatus(t *testing.T, assetID int) int {
	t.Helper()

	var status int
	if err := s.env.DB.QueryRow(`SELECT status_id FROM ass.real_estate WHERE id = $1`, assetID).Scan(&status); err != nil {
		t.Fatalf("read asset %d: %v", assetID, err)
	}
	return status
}

// The asset is funded by the very delivery that issues its last share, and by
// no earlier one. The token is deployed with the decimals the platform keeps
// in blk.token, so the comparison is made in the unit the chain really uses:
// 100 × 10¹⁸ for an 18-decimal token, 100 for a token without decimals.
func TestDeliveryE2EFundsTheAssetOnTheLastShare(t *testing.T) {
	for _, decimals := range []int{18, 0} {
		t.Run(fmt.Sprintf("%d decimals", decimals), func(t *testing.T) {
			ctx := context.Background()
			env := newChainDBEnv(ctx, t)
			delivery := deliveryServices.NewService(NewService())

			salt := fmt.Sprintf("Funded%d", decimals)
			suite := env.Fixture.DeploySuiteDecimals(t, salt, "Tooken "+salt, "TKF", decimals)
			env.SetSharedIRSRow(t, suite.IRSAddr)
			env.InsertTokenRow(t, suite.TokenAddr, "Tooken "+salt, "TKF", decimals)

			asset := seedAssetOffering(t, env, suite, 100)
			first := seedDeliveryOrder(ctx, t, env, suite, asset, 90, true)
			last := seedDeliveryOrder(ctx, t, env, suite, asset, 10, true)

			delivery.Deliver(ctx, first.orderRef)
			if status, _ := first.status(t); status != int(server.IssuanceOrderStatusIdClosed) {
				t.Fatalf("first order is %d, want CLOSED. Last reason: %s", status, first.lastReason(t))
			}
			if got := first.assetStatus(t, asset); got != 4 {
				t.Fatalf("asset is %d after 90 of 100 shares, want 4 (fundraising)", got)
			}

			delivery.Deliver(ctx, last.orderRef)
			if status, _ := last.status(t); status != int(server.IssuanceOrderStatusIdClosed) {
				t.Fatalf("last order is %d, want CLOSED. Last reason: %s", status, last.lastReason(t))
			}
			if got := last.assetStatus(t, asset); got != 5 {
				t.Fatalf("asset is %d after 100 of 100 shares, want 5 (funded)", got)
			}
		})
	}
}

// Two last deliveries finishing together: both orders close, the asset becomes
// funded once, and nobody sees an error.
func TestDeliveryE2ETwoLastDeliveriesFundTheAssetOnce(t *testing.T) {
	ctx := context.Background()
	env := newChainDBEnv(ctx, t)
	delivery := deliveryServices.NewService(NewService())

	suite := env.Fixture.DeploySuite(t, "FundedTwice", "Tooken FundedTwice", "TKFT", common.Address{})
	env.SetSharedIRSRow(t, suite.IRSAddr)
	env.InsertTokenRow(t, suite.TokenAddr, "Tooken FundedTwice", "TKFT", 18)

	asset := seedAssetOffering(t, env, suite, 100)
	early := seedDeliveryOrder(ctx, t, env, suite, asset, 90, true)
	delivery.Deliver(ctx, early.orderRef)

	a := seedDeliveryOrder(ctx, t, env, suite, asset, 5, true)
	b := seedDeliveryOrder(ctx, t, env, suite, asset, 5, true)

	var wg sync.WaitGroup
	for _, scene := range []deliveryScene{a, b} {
		wg.Add(1)
		go func(ref string) {
			defer wg.Done()
			delivery.Deliver(ctx, ref)
		}(scene.orderRef)
	}
	wg.Wait()

	for _, scene := range []deliveryScene{a, b} {
		if status, _ := scene.status(t); status != int(server.IssuanceOrderStatusIdClosed) {
			t.Fatalf("order %s is %d, want CLOSED. Last reason: %s", scene.orderRef, status, scene.lastReason(t))
		}
	}
	if got := a.assetStatus(t, asset); got != 5 {
		t.Fatalf("asset is %d, want 5 (funded)", got)
	}
}

// The chain says no at gas estimation, before anything is signed. The
// order must fail *and* the platform nonce must not move — a consumed nonce
// would mean a transaction went out after all.
//
// The other chain refusal, a mint above the per-investor cap, is proven at the
// mint level by caps_integration_test.go. Raising it here would mean
// replaying the whole module installation for a refusal that travels the same
// code path as this one.
func TestDeliveryE2EChainRefusal(t *testing.T) {
	ctx := context.Background()
	env := newChainDBEnv(ctx, t)
	chainSvc := NewService()

	suite := env.Fixture.DeploySuite(t, "Refusal", "Tooken Refusal", "TKRF", common.Address{})
	env.SetSharedIRSRow(t, suite.IRSAddr)
	env.InsertTokenRow(t, suite.TokenAddr, "Tooken Refusal", "TKRF", 18)

	// An investor who never went through the KYC journey on chain.
	scene := seedDeliveryOrder(ctx, t, env, suite, seedAsset(t, env, suite), 10, false)

	signer := env.Fixture.Auth.From
	before := platformNonce(ctx, t, signer)
	deliveryServices.NewService(chainSvc).Deliver(ctx, scene.orderRef)
	after := platformNonce(ctx, t, signer)

	status, hash := scene.status(t)
	if status != int(server.IssuanceOrderStatusIdDeliveryFailed) {
		t.Fatalf("order is %d, want DELIVERY_FAILED (%d)",
			status, int(server.IssuanceOrderStatusIdDeliveryFailed))
	}
	if hash != "" {
		t.Fatalf("nothing was broadcast, yet the order carries %s", hash)
	}
	if after != before {
		t.Fatalf("the platform nonce moved from %d to %d: a transaction went out", before, after)
	}
	if reason := scene.lastReason(t); !strings.Contains(reason, "Identity is not verified") {
		t.Fatalf("the timeline must name the refusal, got %q", reason)
	}

	balance, err := suite.Token.BalanceOf(env.Fixture.Call, scene.wallet)
	if err != nil {
		t.Fatalf("BalanceOf: %v", err)
	}
	if balance.Sign() != 0 {
		t.Fatalf("the wallet must hold nothing, got %s", balance)
	}
}

// Two investors paid at the same second. Both deliveries read
// PendingNonceAt, and without the mutex they would claim the same nonce — one
// transaction rejected, its hash already recorded, its order stuck forever.
func TestDeliveryE2ETwoOrdersDeliveredAtOnce(t *testing.T) {
	ctx := context.Background()
	env := newChainDBEnv(ctx, t)
	chainSvc := NewService()

	suite := env.Fixture.DeploySuite(t, "Concurrent", "Tooken Concurrent", "TKCC", common.Address{})
	env.SetSharedIRSRow(t, suite.IRSAddr)
	env.InsertTokenRow(t, suite.TokenAddr, "Tooken Concurrent", "TKCC", 18)

	asset := seedAsset(t, env, suite)
	first := seedDeliveryOrder(ctx, t, env, suite, asset, 10, true)
	second := seedDeliveryOrder(ctx, t, env, suite, asset, 7, true)

	// One service instance, as production has: the mutex only protects what
	// shares it.
	delivery := deliveryServices.NewService(chainSvc)

	var wg sync.WaitGroup
	for _, scene := range []deliveryScene{first, second} {
		wg.Add(1)
		go func(ref string) {
			defer wg.Done()
			delivery.Deliver(ctx, ref)
		}(scene.orderRef)
	}
	wg.Wait()

	for _, scene := range []deliveryScene{first, second} {
		status, _ := scene.status(t)
		if status != int(server.IssuanceOrderStatusIdClosed) {
			t.Fatalf("order %s is %d, want CLOSED. Last reason: %s",
				scene.orderRef, status, scene.lastReason(t))
		}
		balance, err := suite.Token.BalanceOf(env.Fixture.Call, scene.wallet)
		if err != nil {
			t.Fatalf("BalanceOf: %v", err)
		}
		if balance.Cmp(shares18(scene.shares)) != 0 {
			t.Fatalf("balance of %s = %s, want %s", scene.orderRef, balance, shares18(scene.shares))
		}
	}
}

// The rule the whole delivery exists to protect. Two deliveries of the *same*
// order, started at the same instant, must mint once. If this fails, an
// investor holds shares nobody sold them, on an immutable ledger.
func TestDeliveryE2EAnOrderIsNeverMintedTwice(t *testing.T) {
	ctx := context.Background()
	env := newChainDBEnv(ctx, t)
	chainSvc := NewService()

	suite := env.Fixture.DeploySuite(t, "DoubleMint", "Tooken DoubleMint", "TKDM", common.Address{})
	env.SetSharedIRSRow(t, suite.IRSAddr)
	env.InsertTokenRow(t, suite.TokenAddr, "Tooken DoubleMint", "TKDM", 18)

	scene := seedDeliveryOrder(ctx, t, env, suite, seedAsset(t, env, suite), 10, true)
	delivery := deliveryServices.NewService(chainSvc)

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			delivery.Deliver(ctx, scene.orderRef)
		}()
	}
	wg.Wait()

	// And once more afterwards, the way a second confirm-payment would.
	delivery.Deliver(ctx, scene.orderRef)

	balance, err := suite.Token.BalanceOf(env.Fixture.Call, scene.wallet)
	if err != nil {
		t.Fatalf("BalanceOf: %v", err)
	}
	if balance.Cmp(shares18(10)) != 0 {
		t.Fatalf("balance = %s, want %s — the order was minted more than once", balance, shares18(10))
	}

	if status, _ := scene.status(t); status != int(server.IssuanceOrderStatusIdClosed) {
		t.Fatalf("order is %d, want CLOSED", status)
	}
}

// The boundary, asserted rather than trusted. The delivery module must
// never reach for the blockchain package, and vice versa.
func TestDeliveryBoundaryHolds(t *testing.T) {
	const mod = "github.com/TookenOrg/tooken-services/"
	cases := []struct{ from, mustNotImport string }{
		{mod + "internal/delivery/...", "internal/blockchain"},
		{mod + "internal/blockchain/...", "internal/delivery"},
	}

	for _, tc := range cases {
		out, err := exec.Command("go", "list", "-deps", tc.from).Output()
		if err != nil {
			t.Fatalf("go list -deps %s: %v", tc.from, err)
		}
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, tc.mustNotImport) {
				t.Errorf("%s depends on %s", tc.from, line)
			}
		}
	}
}
