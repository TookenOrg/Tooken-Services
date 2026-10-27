//go:build integration

package e2e

import (
	"net/http"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	chainglobals "github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	ordersvc "github.com/TookenOrg/tooken-services/internal/orders/services"
	paymentsvc "github.com/TookenOrg/tooken-services/internal/payments/services"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

// The routes an operator calls, once and in this order, to stand the
// platform up on a new chain.
const (
	routeImplementations = "/trex/init/implementations"
	routeIdentityFactory = "/trex/init/contract/factory/identity"
	routeAuthority       = "/trex/init/authority/configure"
	routeTrexFactory     = "/trex/init/contract/factory/trexFactory"
	routeSuite           = "/trex/init/suite/deploy"
)

var initRoutes = []string{routeImplementations, routeIdentityFactory, routeAuthority, routeTrexFactory, routeSuite}

// Order statuses, as the timeline records them.
const (
	statusCreated            = int(server.IssuanceOrderStatusIdCreated)
	statusAwaitingPayment    = int(server.IssuanceOrderStatusIdAwaitingPayment)
	statusPaymentConfirmed   = int(server.IssuanceOrderStatusIdPaymentConfirmed)
	statusDeliveryInProgress = int(server.IssuanceOrderStatusIdDeliveryInProgress)
	statusDelivered          = int(server.IssuanceOrderStatusIdDelivered)
	statusClosed             = int(server.IssuanceOrderStatusIdClosed)
)

// TestEndToEnd stands one platform up through its API, then plays each
// investor journey on it. Every journey lists its own asset and onboards its
// own investors: they share the platform, never their state.
func TestEndToEnd(t *testing.T) {
	p := bootPlatform(t)

	if !t.Run("the operator initialises the platform through the API", func(t *testing.T) {
		initialise(t, p)
	}) {
		t.Fatal("every journey needs an initialised platform")
	}

	t.Run("a paid order is delivered on chain", func(t *testing.T) {
		paidOrderIsDelivered(t, p)
	})
	t.Run("the chain refuses an investor struck off the whitelist", func(t *testing.T) {
		struckOffInvestorIsRefusedByTheChain(t, p)
	})
	t.Run("six investors race for the last shares", func(t *testing.T) {
		lastSharesRace(t, p)
	})
	t.Run("the platform refuses the orders it must not take", func(t *testing.T) {
		ordersRefused(t, p)
	})
	t.Run("the token enforces its limit when the database is wrong", func(t *testing.T) {
		tokenEnforcesItsLimit(t, p)
	})
}

// ---- initialisation -------------------------------------------------------

func initialise(t *testing.T, p *platform) {
	curious := p.signUp(t, "Curious User")

	t.Run("the identity factory refused before the implementations answers one error", func(t *testing.T) {
		before := p.platformNonce(t)

		r := p.call(t, http.MethodPost, routeIdentityFactory, p.admin.token, nil).
			expect(t, http.StatusInternalServerError, "identity factory on a blank platform")

		// decode refuses a second document: the handler used to append a
		// success message after the error.
		answer := decode[server.InitializeIdentityFactoryResponse](t, r)
		if answer.Data != nil {
			t.Fatalf("a failed deployment must not describe a contract: %s", r.body)
		}
		if !strings.Contains(answer.Message, "Failed to deploy Identity Factory") {
			t.Fatalf("the error does not say what failed: %q", answer.Message)
		}

		if n := p.count(t, `SELECT count(*) FROM blk.contract_role WHERE contract_name = $1`,
			chainglobals.IdentityFactoryName); n != 0 {
			t.Fatalf("a refused deployment recorded %d identity factory", n)
		}
		if after := p.platformNonce(t); after != before {
			t.Fatalf("a refused deployment moved the nonce from %d to %d", before, after)
		}
	})

	t.Run("only staff may initialise", func(t *testing.T) {
		before := p.platformNonce(t)

		for _, route := range initRoutes {
			p.call(t, http.MethodPost, route, "", nil).
				expect(t, http.StatusUnauthorized, "anonymous call to "+route)
			p.call(t, http.MethodPost, route, curious.token, nil).
				expect(t, http.StatusForbidden, "USER call to "+route)
		}

		// A refused call must not have started anything in the background.
		p.handler.WaitBackground()
		if n := p.count(t, `SELECT count(*) FROM blk.contract_implementation`); n != 0 {
			t.Fatalf("refused calls deployed %d implementations", n)
		}
		if after := p.platformNonce(t); after != before {
			t.Fatalf("refused calls moved the nonce from %d to %d", before, after)
		}
	})

	var identityFactory string

	t.Run("the five steps, in order", func(t *testing.T) {
		p.call(t, http.MethodPost, routeImplementations, p.admin.token, nil).
			expect(t, http.StatusAccepted, "implementations")
		p.handler.WaitBackground()
		if n := p.count(t, `SELECT count(*) FROM blk.contract_implementation`); n != 9 {
			t.Fatalf("%d implementations recorded, want 9", n)
		}

		r := p.call(t, http.MethodPost, routeIdentityFactory, p.admin.token, nil).
			expect(t, http.StatusOK, "identity factory")
		answer := decode[server.InitializeIdentityFactoryResponse](t, r)
		if answer.Data == nil || !common.IsHexAddress(answer.Data.Address) {
			t.Fatalf("the identity factory answered no address: %s", r.body)
		}
		identityFactory = answer.Data.Address

		p.call(t, http.MethodPost, routeAuthority, p.admin.token, nil).
			expect(t, http.StatusOK, "authority")

		p.call(t, http.MethodPost, routeTrexFactory, p.admin.token, nil).
			expect(t, http.StatusAccepted, "T-REX factory")
		p.handler.WaitBackground()

		p.call(t, http.MethodPost, routeSuite, p.admin.token, nil).
			expect(t, http.StatusAccepted, "T-REX suite")
		p.handler.WaitBackground()

		// The background steps answer before they work, so their outcome is
		// read where they leave it: one role per singleton the platform runs on.
		for _, role := range []string{
			chainglobals.IdentityFactoryName,
			chainglobals.TrexFactoryName,
			chainglobals.ClaimIssuerName,
			chainglobals.SharedIdentityRegistryStorageName,
			chainglobals.TransferRestrictionModuleName,
		} {
			if n := p.count(t, `SELECT count(*) FROM blk.contract_role WHERE contract_name = $1`, role); n != 1 {
				t.Fatalf("%d %s recorded, want 1", n, role)
			}
		}
		if n := p.count(t, `SELECT count(*) FROM blk.contract_role WHERE contract_name = $1 AND lower(address) = lower($2)`,
			chainglobals.IdentityFactoryName, identityFactory); n != 1 {
			t.Fatalf("the identity factory recorded is not the one answered, %s", identityFactory)
		}
	})

	t.Run("the identity factory is deployed once", func(t *testing.T) {
		if identityFactory == "" {
			t.Skip("the identity factory was not deployed")
		}
		before := p.platformNonce(t)

		r := p.call(t, http.MethodPost, routeIdentityFactory, p.admin.token, nil).
			expect(t, http.StatusOK, "identity factory called again")
		answer := decode[server.InitializeIdentityFactoryResponse](t, r)
		if answer.Data == nil || !strings.EqualFold(answer.Data.Address, identityFactory) {
			t.Fatalf("called again, the identity factory answered %s, want %s", r.body, identityFactory)
		}
		if after := p.platformNonce(t); after != before {
			t.Fatalf("calling it again redeployed: the nonce moved from %d to %d", before, after)
		}
	})
}

// ---- a paid order is delivered --------------------------------------------

func paidOrderIsDelivered(t *testing.T, p *platform) {
	assetID, token := p.listAsset(t, assetSpec{
		title: "Tooken Villa Belair", totalShares: 100, pricePerShare: "250.00",
		entryFeeRate: "2.5", maxPerInvestor: 20,
	})
	alice := p.onboardInvestor(t, "Alice Martin")

	key := uuid.NewString()
	r, err := p.placeOrder(alice, assetID, 10, key)
	if err != nil {
		t.Fatal(err)
	}
	placed := createdOrder(t, r.expect(t, http.StatusCreated, "Alice's order"))

	if placed.StatusCode != "AWAITING_PAYMENT" || placed.TokenQuantity != 10 {
		t.Fatalf("the order is %s for %d shares, want AWAITING_PAYMENT for 10", placed.StatusCode, placed.TokenQuantity)
	}
	if placed.UnitPrice != "250.00000000" || placed.GrossAmount != "2500.00" ||
		placed.FeeAmount != "62.50" || placed.AmountDue != "2562.50" {
		t.Fatalf("priced %s × 10 = %s + %s = %s, want 250.00000000 × 10 = 2500.00 + 62.50 = 2562.50",
			placed.UnitPrice, placed.GrossAmount, placed.FeeAmount, placed.AmountDue)
	}
	if len(placed.Fees) != 1 || placed.Fees[0].Code != "ENTRY" || placed.Fees[0].Rate != "2.5000" ||
		placed.Fees[0].BaseAmount != "2500.00" || placed.Fees[0].Amount != "62.50" {
		t.Fatalf("fees %+v, want one ENTRY fee of 62.50 at 2.5000 %% of 2500.00", placed.Fees)
	}
	if placed.DeliveryTxHash != nil {
		t.Fatalf("an unpaid order carries a delivery %s", *placed.DeliveryTxHash)
	}

	// A front end retrying after a lost answer sends the same key again.
	r, err = p.placeOrder(alice, assetID, 10, key)
	if err != nil {
		t.Fatal(err)
	}
	if replayed := createdOrder(t, r.expect(t, http.StatusOK, "Alice's order replayed")); replayed.OrderRef != placed.OrderRef {
		t.Fatalf("the replay answered %s, want the original %s", replayed.OrderRef, placed.OrderRef)
	}
	if n := p.count(t, `SELECT count(*) FROM iss.issuance_orders WHERE user_id = $1`, alice.id); n != 1 {
		t.Fatalf("Alice has %d orders after a replay, want 1", n)
	}

	r, err = p.confirmPayment(p.admin, placed.OrderRef)
	if err != nil {
		t.Fatal(err)
	}
	if paid := createdOrder(t, r.expect(t, http.StatusOK, "payment of Alice's order")); paid.StatusCode != "PAYMENT_CONFIRMED" {
		t.Fatalf("the payment answered %s, want PAYMENT_CONFIRMED", paid.StatusCode)
	}

	closed := p.awaitOrder(t, alice, placed.OrderRef, "CLOSED", "DELIVERY_FAILED")
	if closed.StatusCode != "CLOSED" {
		t.Fatalf("the order is %s, want CLOSED. Reason: %q", closed.StatusCode, p.lastReason(t, placed.OrderRef))
	}
	if closed.DeliveryTxHash == nil {
		t.Fatal("a delivered order must carry its transaction")
	}

	if got := p.balanceOf(t, token, alice.wallet); got.Cmp(shares(10)) != 0 {
		t.Fatalf("Alice holds %s, want %s", got, shares(10))
	}
	if got := p.totalSupply(t, token); got.Cmp(shares(10)) != 0 {
		t.Fatalf("the supply is %s, want %s", got, shares(10))
	}
	minted := p.mints(t, token, *closed.DeliveryTxHash)
	if len(minted) != 1 || minted[0].To != alice.wallet || minted[0].Value.Cmp(shares(10)) != 0 {
		t.Fatalf("the delivery minted %+v, want one mint of %s to %s", minted, shares(10), alice.wallet)
	}

	want := []int{statusCreated, statusAwaitingPayment, statusPaymentConfirmed,
		statusDeliveryInProgress, statusDelivered, statusClosed}
	if got := p.timeline(t, placed.OrderRef); !slices.Equal(got, want) {
		t.Fatalf("timeline %v, want %v", got, want)
	}
	if n := p.count(t, `SELECT count(*) FROM iss.issuance_order_payments`+orderRows, placed.OrderRef); n != 1 {
		t.Fatalf("%d payments recorded, want 1", n)
	}
	if sold := p.tokensSold(t, assetID); sold != 10 {
		t.Fatalf("the asset shows %d shares sold, want 10", sold)
	}

	t.Run("a closed order cannot be paid again", func(t *testing.T) {
		before := p.ledger(t)
		r, err := p.confirmPayment(p.admin, placed.OrderRef)
		if err != nil {
			t.Fatal(err)
		}
		p.refused(t, before, r, http.StatusConflict, paymentsvc.ErrOrderNotAwaitingPayment)
		if n := p.count(t, `SELECT count(*) FROM iss.issuance_order_payments`+orderRows, placed.OrderRef); n != 1 {
			t.Fatalf("%d payments recorded, want 1", n)
		}
	})

	t.Run("a payment confirmed twice at once mints once", func(t *testing.T) {
		second := p.order(t, alice, assetID, 5)

		start := make(chan struct{})
		var wg sync.WaitGroup
		replies := make([]reply, 2)
		errs := make([]error, 2)
		for i := range replies {
			wg.Go(func() {
				<-start
				replies[i], errs[i] = p.confirmPayment(p.admin, second.OrderRef)
			})
		}
		close(start)
		wg.Wait()

		accepted := 0
		for i, r := range replies {
			if errs[i] != nil {
				t.Fatal(errs[i])
			}
			switch r.code {
			case http.StatusOK:
				accepted++
			case http.StatusConflict:
			default:
				t.Fatalf("a concurrent confirmation answered %d: %s", r.code, r.body)
			}
		}
		if accepted == 0 {
			t.Fatal("neither confirmation was accepted")
		}

		if done := p.awaitOrder(t, alice, second.OrderRef, "CLOSED", "DELIVERY_FAILED"); done.StatusCode != "CLOSED" {
			t.Fatalf("the order is %s, want CLOSED. Reason: %q", done.StatusCode, p.lastReason(t, second.OrderRef))
		}
		if n := p.count(t, `SELECT count(*) FROM iss.issuance_order_payments`+orderRows, second.OrderRef); n != 1 {
			t.Fatalf("%d payments recorded, want 1", n)
		}
		if got := p.balanceOf(t, token, alice.wallet); got.Cmp(shares(15)) != 0 {
			t.Fatalf("Alice holds %s, want %s: the order was minted more than once", got, shares(15))
		}
		if got := p.totalSupply(t, token); got.Cmp(shares(15)) != 0 {
			t.Fatalf("the supply is %s, want %s", got, shares(15))
		}
	})
}

// orderRows narrows a per-order table to one order, by reference.
const orderRows = ` WHERE order_id = (SELECT id FROM iss.issuance_orders WHERE order_reference = $1)`

// ---- the chain has the last word ------------------------------------------

func struckOffInvestorIsRefusedByTheChain(t *testing.T, p *platform) {
	assetID, token := p.listAsset(t, assetSpec{
		title: "Tooken Résidence Garibaldi", totalShares: 40, pricePerShare: "100.00",
		entryFeeRate: "2.5", maxPerInvestor: 20,
	})
	bob := p.onboardInvestor(t, "Bob Leroy")

	// Simulates a compliance officer removing Bob from the whitelist on chain
	// while the database still says verified: nothing in the API can bring
	// the two apart, which is precisely why the chain must still say no.
	p.strikeOff(t, token, bob.wallet)

	placed := p.order(t, bob, assetID, 10)

	// Read after the strike-off, which is itself a transaction.
	nonce := p.platformNonce(t)
	supply := p.totalSupply(t, token)

	r, err := p.confirmPayment(p.admin, placed.OrderRef)
	if err != nil {
		t.Fatal(err)
	}
	r.expect(t, http.StatusOK, "payment of Bob's order")

	failed := p.awaitOrder(t, bob, placed.OrderRef, "CLOSED", "DELIVERY_FAILED")
	if failed.StatusCode != "DELIVERY_FAILED" {
		t.Fatalf("the order is %s, want DELIVERY_FAILED", failed.StatusCode)
	}
	if failed.DeliveryTxHash != nil {
		t.Fatalf("nothing was broadcast, yet the order carries %s", *failed.DeliveryTxHash)
	}
	if reason := p.lastReason(t, placed.OrderRef); !strings.Contains(reason, "Identity is not verified") {
		t.Fatalf("the timeline must name the refusal, got %q", reason)
	}
	if after := p.platformNonce(t); after != nonce {
		t.Fatalf("the nonce moved from %d to %d: a transaction went out", nonce, after)
	}
	if got := p.balanceOf(t, token, bob.wallet); got.Sign() != 0 {
		t.Fatalf("Bob holds %s, want nothing", got)
	}
	if got := p.totalSupply(t, token); got.Cmp(supply) != 0 {
		t.Fatalf("the supply moved from %s to %s", supply, got)
	}

	// The shares stay held: a paid order that failed waits for an operator to
	// retry or refund it, and selling them meanwhile would oversell the asset.
	if sold := p.tokensSold(t, assetID); sold != 10 {
		t.Fatalf("the asset shows %d shares sold, want the 10 still held", sold)
	}
}

// ---- the last shares ------------------------------------------------------

func lastSharesRace(t *testing.T, p *platform) {
	assetID, token := p.listAsset(t, assetSpec{
		title: "Tooken Les Terrasses", totalShares: 30, pricePerShare: "150.00",
		entryFeeRate: "2.5", maxPerInvestor: 10,
	})

	investors := p.onboardInvestors(t, "Ana Roux", "Ben Faure", "Chloé Blanc", "Dan Morel", "Emma Girard", "Fabien Henry")
	before := p.count(t, `SELECT count(*) FROM iss.issuance_orders`)

	// Six orders of 10 for 30 shares, all waiting on the asset before any of
	// them is let through.
	release := p.holdAsset(t, assetID)
	var wg sync.WaitGroup
	replies := make([]reply, len(investors))
	errs := make([]error, len(investors))
	for i, investor := range investors {
		wg.Go(func() {
			replies[i], errs[i] = p.placeOrder(investor, assetID, 10, uuid.NewString())
		})
	}
	lined := p.awaitLockWaiters(t, len(investors))
	release()
	wg.Wait()
	if !lined {
		t.Fatalf("the %d orders never all waited on the asset", len(investors))
	}

	var winners []account
	var orders []server.IssuanceOrder
	var losers []account
	for i, r := range replies {
		if errs[i] != nil {
			t.Fatal(errs[i])
		}
		switch r.code {
		case http.StatusCreated:
			winners = append(winners, investors[i])
			orders = append(orders, createdOrder(t, r))
		case http.StatusConflict:
			if got := r.message(t); !strings.Contains(got, ordersvc.ErrNotEnoughShares.Error()) {
				t.Fatalf("a losing order says %q, want %q", got, ordersvc.ErrNotEnoughShares)
			}
			losers = append(losers, investors[i])
		default:
			t.Fatalf("an order answered %d: %s", r.code, r.body)
		}
	}
	if len(winners) != 3 || len(losers) != 3 {
		t.Fatalf("%d orders accepted and %d refused, want 3 and 3", len(winners), len(losers))
	}
	if n := p.count(t, `SELECT count(*) FROM iss.issuance_orders`); n != before+3 {
		t.Fatalf("%d orders recorded, want %d", n-before, 3)
	}
	if sold := p.tokensSold(t, assetID); sold != 30 {
		t.Fatalf("%d shares held, want exactly 30", sold)
	}

	t.Run("the winners pay at once and the asset is funded", func(t *testing.T) {
		start := make(chan struct{})
		var wg sync.WaitGroup
		paid := make([]reply, len(orders))
		errs := make([]error, len(orders))
		for i, o := range orders {
			wg.Go(func() {
				<-start
				paid[i], errs[i] = p.confirmPayment(p.admin, o.OrderRef)
			})
		}
		close(start)
		wg.Wait()

		for i, o := range orders {
			if errs[i] != nil {
				t.Fatal(errs[i])
			}
			paid[i].expect(t, http.StatusOK, "payment of "+o.OrderRef)
		}
		for i, o := range orders {
			if done := p.awaitOrder(t, winners[i], o.OrderRef, "CLOSED", "DELIVERY_FAILED"); done.StatusCode != "CLOSED" {
				t.Fatalf("order %s is %s, want CLOSED. Reason: %q", o.OrderRef, done.StatusCode, p.lastReason(t, o.OrderRef))
			}
			if got := p.balanceOf(t, token, winners[i].wallet); got.Cmp(shares(10)) != 0 {
				t.Fatalf("%s holds %s, want %s", winners[i].email, got, shares(10))
			}
		}
		if got := p.totalSupply(t, token); got.Cmp(shares(30)) != 0 {
			t.Fatalf("the supply is %s, want %s", got, shares(30))
		}

		p.awaitAssetStatus(t, assetID, "funded")

		before := p.ledger(t)
		r, err := p.placeOrder(losers[0], assetID, 1, uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}
		p.refused(t, before, r, http.StatusConflict, ordersvc.ErrRealEstateNotAvailable)
	})
}

// ---- refusals -------------------------------------------------------------

func ordersRefused(t *testing.T, p *platform) {
	assetID, _ := p.listAsset(t, assetSpec{
		title: "Tooken Le Panoramique", totalShares: 30, pricePerShare: "200.00",
		entryFeeRate: "2.5", maxPerInvestor: 20,
	})
	draftID := p.createAsset(t, assetSpec{
		title: "Tooken Le Brouillon", totalShares: 10, pricePerShare: "100.00",
		entryFeeRate: "2.5", maxPerInvestor: 5,
	})
	publishedID := p.createAsset(t, assetSpec{
		title: "Tooken Bientôt", totalShares: 10, pricePerShare: "100.00",
		entryFeeRate: "2.5", maxPerInvestor: 5,
	})
	p.publish(t, publishedID)

	// One applicant at each stage of the KYC journey.
	var carla, david, noKyc, filed, approved, expired account
	together(t, "applicants", map[string]func(t *testing.T){
		"Carla, registered on chain": func(t *testing.T) { carla = p.onboardInvestor(t, "Carla Dupont") },
		"David, registered on chain": func(t *testing.T) { david = p.onboardInvestor(t, "David Moreau") },
		"Eve, no KYC":                func(t *testing.T) { noKyc = p.signUp(t, "Eve Lambert") },
		"Frank, KYC filed": func(t *testing.T) {
			filed = p.signUp(t, "Frank Mercier")
			p.submitKyc(t, filed, "Frank Mercier")
		},
		"Gina, KYC approved": func(t *testing.T) {
			approved = p.signUp(t, "Gina Rossi")
			p.approveKyc(t, p.submitKyc(t, approved, "Gina Rossi"))
		},
		"Hugo, KYC about to expire": func(t *testing.T) { expired = p.onboardInvestor(t, "Hugo Bernard") },
	})

	// Simulates a KYC whose validity ran out. The dates are moved on the user,
	// which is what the order placement reads; moving them on the
	// verification would let its trigger rewrite the status, and the refusal
	// would then come from the status rather than from the date.
	if _, err := p.db.Exec(`
		UPDATE usr.users
		SET kyc_verified_at = now() - interval '2 days',
		    kyc_expires_at  = now() - interval '1 day'
		WHERE id = $1`, expired.id); err != nil {
		t.Fatal(err)
	}

	// Carla takes 20 of the 30 shares, her limit; 10 remain.
	p.order(t, carla, assetID, 20)

	cases := []struct {
		name     string
		who      account
		asset    int
		quantity int
		want     int
		reason   error
	}{
		{"without a session", account{}, assetID, 1, http.StatusUnauthorized, nil},
		{"without a KYC", noKyc, assetID, 1, http.StatusForbidden, ordersvc.ErrInvestorNotEligible},
		{"with a KYC filed but not decided", filed, assetID, 1, http.StatusForbidden, ordersvc.ErrInvestorNotEligible},
		{"with a KYC approved but not registered on chain", approved, assetID, 1, http.StatusForbidden, ordersvc.ErrInvestorNotEligible},
		{"with an expired KYC", expired, assetID, 1, http.StatusForbidden, ordersvc.ErrInvestorNotEligible},
		{"on an unknown asset", david, 999_999, 1, http.StatusNotFound, ordersvc.ErrRealEstateNotFound},
		{"on a draft, which the public cannot see", david, draftID, 1, http.StatusNotFound, ordersvc.ErrRealEstateNotFound},
		{"on a published asset not yet open", david, publishedID, 1, http.StatusConflict, ordersvc.ErrRealEstateNotAvailable},
		{"for no share", david, assetID, 0, http.StatusBadRequest, ordersvc.ErrQuantityInvalid},
		{"for a negative quantity", david, assetID, -1, http.StatusBadRequest, ordersvc.ErrQuantityInvalid},
		{"above the limit per investor", carla, assetID, 1, http.StatusConflict, ordersvc.ErrStakeLimit},
		{"above the shares left, under the limit", david, assetID, 15, http.StatusConflict, ordersvc.ErrNotEnoughShares},
	}
	for _, c := range cases {
		t.Run("an order "+c.name, func(t *testing.T) {
			before := p.ledger(t)
			r, err := p.placeOrder(c.who, c.asset, c.quantity, uuid.NewString())
			if err != nil {
				t.Fatal(err)
			}
			p.refused(t, before, r, c.want, c.reason)
		})
	}

	t.Run("a key reused for another order", func(t *testing.T) {
		key := uuid.NewString()
		r, err := p.placeOrder(david, assetID, 5, key)
		if err != nil {
			t.Fatal(err)
		}
		r.expect(t, http.StatusCreated, "David's order")

		before := p.ledger(t)
		r, err = p.placeOrder(david, assetID, 6, key)
		if err != nil {
			t.Fatal(err)
		}
		p.refused(t, before, r, http.StatusUnprocessableEntity, ordersvc.ErrIdempotencyKeyReused)
	})

	t.Run("a payment confirmed by an investor", func(t *testing.T) {
		mine := p.order(t, david, assetID, 1)

		before := p.ledger(t)
		r, err := p.confirmPayment(david, mine.OrderRef)
		if err != nil {
			t.Fatal(err)
		}
		p.refused(t, before, r, http.StatusForbidden, nil)
		if still := p.fetchOrder(t, david, mine.OrderRef); still.StatusCode != "AWAITING_PAYMENT" {
			t.Fatalf("the order is %s, want it still AWAITING_PAYMENT", still.StatusCode)
		}
		if n := p.count(t, `SELECT count(*) FROM iss.issuance_order_payments`+orderRows, mine.OrderRef); n != 0 {
			t.Fatalf("%d payments recorded, want none", n)
		}
	})

	t.Run("a payment for an unknown order", func(t *testing.T) {
		before := p.ledger(t)
		r, err := p.confirmPayment(p.admin, "ORD-UNKNOWN")
		if err != nil {
			t.Fatal(err)
		}
		p.refused(t, before, r, http.StatusNotFound, ordersvc.ErrOrderNotFound)
	})
}

// ---- the token's own limit ------------------------------------------------

func tokenEnforcesItsLimit(t *testing.T, p *platform) {
	assetID, token := p.listAsset(t, assetSpec{
		title: "Tooken Domaine des Oliviers", totalShares: 30, pricePerShare: "120.00",
		entryFeeRate: "2.5", maxPerInvestor: 10,
	})

	// Simulates a limit raised in the database and never carried to the
	// token, whose compliance module still caps every holder at 10 shares.
	if _, err := p.db.Exec(`UPDATE ass.real_estate_shares_config SET max_shares_per_investor = 30 WHERE real_estate_id = $1`,
		assetID); err != nil {
		t.Fatal(err)
	}

	investors := p.onboardInvestors(t, "Inès Garnier", "Jules Petit")
	ines, jules := investors[0], investors[1]
	placed := p.order(t, ines, assetID, 20)

	nonce := p.platformNonce(t)
	r, err := p.confirmPayment(p.admin, placed.OrderRef)
	if err != nil {
		t.Fatal(err)
	}
	r.expect(t, http.StatusOK, "payment of Inès's order")

	failed := p.awaitOrder(t, ines, placed.OrderRef, "CLOSED", "DELIVERY_FAILED")
	if failed.StatusCode != "DELIVERY_FAILED" {
		t.Fatalf("the order is %s, want DELIVERY_FAILED", failed.StatusCode)
	}
	if reason := p.lastReason(t, placed.OrderRef); !strings.Contains(reason, "Compliance not followed") {
		t.Fatalf("the timeline must name the refusal, got %q", reason)
	}
	if after := p.platformNonce(t); after != nonce {
		t.Fatalf("the nonce moved from %d to %d: a transaction went out", nonce, after)
	}
	if got := p.balanceOf(t, token, ines.wallet); got.Sign() != 0 {
		t.Fatalf("Inès holds %s, want nothing", got)
	}

	t.Run("the asset still delivers an order within the limit", func(t *testing.T) {
		ok := p.order(t, jules, assetID, 10)

		r, err := p.confirmPayment(p.admin, ok.OrderRef)
		if err != nil {
			t.Fatal(err)
		}
		r.expect(t, http.StatusOK, "payment of Jules's order")

		if done := p.awaitOrder(t, jules, ok.OrderRef, "CLOSED", "DELIVERY_FAILED"); done.StatusCode != "CLOSED" {
			t.Fatalf("the order is %s, want CLOSED. Reason: %q", done.StatusCode, p.lastReason(t, ok.OrderRef))
		}
		if got := p.balanceOf(t, token, jules.wallet); got.Cmp(shares(10)) != 0 {
			t.Fatalf("Jules holds %s, want %s", got, shares(10))
		}
		if got := p.totalSupply(t, token); got.Cmp(shares(10)) != 0 {
			t.Fatalf("the supply is %s, want %s", got, shares(10))
		}

		// Every share is now reserved, twenty of them by the failed order
		// awaiting an operator, but only ten exist on chain. The asset is
		// funded by what the chain holds, never by what the database reserved.
		if sold := p.tokensSold(t, assetID); sold != 30 {
			t.Fatalf("%d shares reserved, want all 30", sold)
		}
		if status := p.asset(t, assetID).Status; status == nil || *status != "fundraising" {
			got := "without a status"
			if status != nil {
				got = *status
			}
			t.Fatalf("the asset is %q with 10 of its 30 shares on chain, want it still fundraising", got)
		}
	})
}
