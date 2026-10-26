//go:build integration

package database

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/server"
)

const (
	statusPaymentConfirmed = int(server.IssuanceOrderStatusIdPaymentConfirmed)
	statusInProgress       = int(server.IssuanceOrderStatusIdDeliveryInProgress)
	statusFailed           = int(server.IssuanceOrderStatusIdDeliveryFailed)
	statusDelivered        = int(server.IssuanceOrderStatusIdDelivered)
	statusClosed           = int(server.IssuanceOrderStatusIdClosed)
)

// The database half: a deliverable order hands back everything the mint
// needs, and moves to DELIVERY_IN_PROGRESS with the system — not the investor —
// as the actor.
func TestStartDeliveryHandsBackTheMintTarget(t *testing.T) {
	ctx := context.Background()
	env := newDeliveryEnv(t)
	order := env.seed(t, aPaidOrder())

	target, err := StartDelivery(ctx, order.OrderRef)
	if err != nil {
		t.Fatalf("StartDelivery: %v", err)
	}

	if target.OrderID != order.OrderID || target.UserID != order.UserID || target.Shares != 10 {
		t.Fatalf("target identifies the wrong order: %+v", target)
	}
	if target.Wallet != order.Wallet {
		t.Fatalf("wallet = %q, want %q", target.Wallet, order.Wallet)
	}
	// Read through the ass.real_estate trigger, not from the token row: that
	// join is the part a unit test cannot check.
	if target.TokenAddress != order.TokenAddress {
		t.Fatalf("token address = %q, want %q", target.TokenAddress, order.TokenAddress)
	}

	if got, hash := env.status(t, order.OrderRef); got != statusInProgress || hash != nil {
		t.Fatalf("order is %d with hash %v, want %d with no hash", got, hash, statusInProgress)
	}

	timeline := env.timeline(t, order.OrderID)
	if len(timeline) != 1 {
		t.Fatalf("want one transition, got %d: %+v", len(timeline), timeline)
	}
	if timeline[0].From != statusPaymentConfirmed || timeline[0].To != statusInProgress {
		t.Fatalf("transition is %d -> %d, want %d -> %d",
			timeline[0].From, timeline[0].To, statusPaymentConfirmed, statusInProgress)
	}
	if !timeline[0].ActorIsNull {
		t.Fatal("the actor must be NULL: the delivery is automatic, no human asked for it")
	}
}

// An invalid KYC, a missing wallet and the asset that was never tokenised.
// Each one is refused before anything is signed, and the order goes straight
// to DELIVERY_FAILED: one transition, not two.
func TestStartDeliveryRefusesWhatCannotBeMinted(t *testing.T) {
	cases := []struct {
		name    string
		fixture func() orderFixture
		reason  string
	}{
		{
			name: "the investor KYC is not verified",
			fixture: func() orderFixture {
				f := aPaidOrder()
				f.kycStatus = "pending"
				return f
			},
			reason: "KYC",
		},
		{
			name: "the investor KYC expired between payment and delivery",
			fixture: func() orderFixture {
				f := aPaidOrder()
				f.kycExpired = true
				return f
			},
			reason: "KYC",
		},
		{
			name: "the investor has no active wallet",
			fixture: func() orderFixture {
				f := aPaidOrder()
				f.activeWallet = false
				return f
			},
			reason: "wallet",
		},
		{
			name: "the asset carries no token",
			fixture: func() orderFixture {
				f := aPaidOrder()
				f.tokenised = false
				return f
			},
			reason: "token",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			env := newDeliveryEnv(t)
			order := env.seed(t, tc.fixture())

			_, err := StartDelivery(ctx, order.OrderRef)
			if !errors.Is(err, ErrDeliveryRefused) {
				t.Fatalf("want ErrDeliveryRefused, got %v", err)
			}

			// The refusal has to survive: it is written inside the very
			// transaction that reports it, and a rollback would erase it.
			got, hash := env.status(t, order.OrderRef)
			if got != statusFailed {
				t.Fatalf("order is %d, want %d — the failure was rolled back", got, statusFailed)
			}
			if hash != nil {
				t.Fatalf("nothing was signed, yet the order carries hash %v", *hash)
			}

			timeline := env.timeline(t, order.OrderID)
			if len(timeline) != 1 {
				t.Fatalf("want one transition, got %d: %+v", len(timeline), timeline)
			}
			if timeline[0].From != statusPaymentConfirmed || timeline[0].To != statusFailed {
				t.Fatalf("transition is %d -> %d, want %d -> %d",
					timeline[0].From, timeline[0].To, statusPaymentConfirmed, statusFailed)
			}
			if timeline[0].Reason == "" {
				t.Fatal("a failure without a reason tells the operator nothing")
			}
		})
	}
}

// An order that is not awaiting delivery must be left exactly as it is. This is
// what makes a second Deliver harmless.
func TestStartDeliveryWritesNothingOnAnOrderThatMovedOn(t *testing.T) {
	for _, status := range []int{statusInProgress, statusFailed, statusDelivered, statusClosed} {
		t.Run("", func(t *testing.T) {
			ctx := context.Background()
			env := newDeliveryEnv(t)

			f := aPaidOrder()
			f.statusID = status
			order := env.seed(t, f)

			_, err := StartDelivery(ctx, order.OrderRef)
			if !errors.Is(err, ErrOrderNotDeliverable) {
				t.Fatalf("want ErrOrderNotDeliverable, got %v", err)
			}

			if got, _ := env.status(t, order.OrderRef); got != status {
				t.Fatalf("order moved from %d to %d", status, got)
			}
			if timeline := env.timeline(t, order.OrderID); len(timeline) != 0 {
				t.Fatalf("nothing must be written, got %+v", timeline)
			}
		})
	}
}

// The half no unit test can reach: the conditional UPDATE is the arbiter
// between two deliveries of the same order. Whatever the interleaving, exactly
// one of them may record a hash — and therefore exactly one may send a mint.
func TestRecordDeliveryTxHashLetsOnlyOneDeliveryThrough(t *testing.T) {
	ctx := context.Background()
	env := newDeliveryEnv(t)
	order := env.seed(t, aPaidOrder())

	if _, err := StartDelivery(ctx, order.OrderRef); err != nil {
		t.Fatalf("StartDelivery: %v", err)
	}

	const racers = 8
	var wg sync.WaitGroup
	results := make([]error, racers)

	for i := 0; i < racers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// A different hash each time, as two signed transactions would have.
			results[i] = RecordDeliveryTxHash(ctx, order.OrderID, nthTxHash(i))
		}(i)
	}
	wg.Wait()

	var winners int
	for i, err := range results {
		switch {
		case err == nil:
			winners++
		case errors.Is(err, ErrDeliveryAlreadyRecorded):
		default:
			t.Fatalf("racer %d failed for the wrong reason: %v", i, err)
		}
	}

	if winners != 1 {
		t.Fatalf("%d deliveries recorded a hash; exactly one must, or the order is minted twice", winners)
	}

	if _, hash := env.status(t, order.OrderRef); hash == nil {
		t.Fatal("the winner's hash must be on the order")
	}
}

// The same gate, seen from the other side: an order that is not
// DELIVERY_IN_PROGRESS never gets a hash, so a stale delivery cannot revive it.
func TestRecordDeliveryTxHashRefusesAnOrderThatIsNotBeingDelivered(t *testing.T) {
	ctx := context.Background()
	env := newDeliveryEnv(t)
	order := env.seed(t, aPaidOrder())

	err := RecordDeliveryTxHash(ctx, order.OrderID, aTxHash)
	if !errors.Is(err, ErrDeliveryAlreadyRecorded) {
		t.Fatalf("want ErrDeliveryAlreadyRecorded, got %v", err)
	}

	if _, hash := env.status(t, order.OrderRef); hash != nil {
		t.Fatalf("the order must carry no hash, got %v", *hash)
	}
}

// The send was refused and the node never saw the transaction. The
// attempt stays on the timeline, but the column goes back to NULL — it answers
// "what was broadcast?", and the answer is nothing.
func TestFailDeliveryTakesTheHashBackWhenNothingWasBroadcast(t *testing.T) {
	ctx := context.Background()
	env := newDeliveryEnv(t)
	order := env.seed(t, aPaidOrder())

	if _, err := StartDelivery(ctx, order.OrderRef); err != nil {
		t.Fatalf("StartDelivery: %v", err)
	}
	if err := RecordDeliveryTxHash(ctx, order.OrderID, aTxHash); err != nil {
		t.Fatalf("RecordDeliveryTxHash: %v", err)
	}

	hash := aTxHash
	if err := FailDelivery(ctx, order.OrderRef, "the send was refused", &hash, false); err != nil {
		t.Fatalf("FailDelivery: %v", err)
	}

	got, column := env.status(t, order.OrderRef)
	if got != statusFailed {
		t.Fatalf("order is %d, want %d", got, statusFailed)
	}
	if column != nil {
		t.Fatalf("delivery_tx_hash = %v, want NULL: nothing reached the network", *column)
	}

	timeline := env.timeline(t, order.OrderID)
	last := timeline[len(timeline)-1]
	if last.TxHash == nil || *last.TxHash != aTxHash {
		t.Fatalf("the timeline must keep the attempt, got %v", last.TxHash)
	}
}

// The counterpart: a transaction that was mined and reverted exists on the
// chain, so its hash stays on the order.
func TestFailDeliveryKeepsTheHashWhenTheTransactionExists(t *testing.T) {
	ctx := context.Background()
	env := newDeliveryEnv(t)
	order := env.seed(t, aPaidOrder())

	if _, err := StartDelivery(ctx, order.OrderRef); err != nil {
		t.Fatalf("StartDelivery: %v", err)
	}
	if err := RecordDeliveryTxHash(ctx, order.OrderID, aTxHash); err != nil {
		t.Fatalf("RecordDeliveryTxHash: %v", err)
	}

	hash := aTxHash
	if err := FailDelivery(ctx, order.OrderRef, "the mint reverted", &hash, true); err != nil {
		t.Fatalf("FailDelivery: %v", err)
	}

	got, column := env.status(t, order.OrderRef)
	if got != statusFailed {
		t.Fatalf("order is %d, want %d", got, statusFailed)
	}
	if column == nil || *column != aTxHash {
		t.Fatalf("delivery_tx_hash = %v, want the hash: the transaction is on the chain", column)
	}
}

// The end of the nominal path: DELIVERED and CLOSED land in the same
// transaction, both carrying the hash. The schema CHECKs that a CLOSED order
// has one, so a bug here is a constraint violation, not a silent wrong state.
func TestCompleteDeliveryClosesTheOrderInOneTransaction(t *testing.T) {
	ctx := context.Background()
	env := newDeliveryEnv(t)
	order := env.seed(t, aPaidOrder())

	if _, err := StartDelivery(ctx, order.OrderRef); err != nil {
		t.Fatalf("StartDelivery: %v", err)
	}
	if err := RecordDeliveryTxHash(ctx, order.OrderID, aTxHash); err != nil {
		t.Fatalf("RecordDeliveryTxHash: %v", err)
	}

	if err := CompleteDelivery(ctx, order.OrderRef, aTxHash); err != nil {
		t.Fatalf("CompleteDelivery: %v", err)
	}

	got, column := env.status(t, order.OrderRef)
	if got != statusClosed {
		t.Fatalf("order is %d, want %d", got, statusClosed)
	}
	if column == nil || *column != aTxHash {
		t.Fatalf("delivery_tx_hash = %v, want the mint hash", column)
	}

	timeline := env.timeline(t, order.OrderID)
	if len(timeline) != 3 {
		t.Fatalf("want 4->8, 8->10, 10->5, got %d transitions: %+v", len(timeline), timeline)
	}
	if timeline[1].From != statusInProgress || timeline[1].To != statusDelivered {
		t.Fatalf("second transition is %d -> %d, want %d -> %d",
			timeline[1].From, timeline[1].To, statusInProgress, statusDelivered)
	}
	if timeline[2].From != statusDelivered || timeline[2].To != statusClosed {
		t.Fatalf("third transition is %d -> %d, want %d -> %d",
			timeline[2].From, timeline[2].To, statusDelivered, statusClosed)
	}
	for _, tr := range timeline[1:] {
		if tr.TxHash == nil || *tr.TxHash != aTxHash {
			t.Fatalf("every delivery transition must name the mint, got %v", tr.TxHash)
		}
	}
}

// CompleteDelivery is told which mint it is closing. A hash that is not the
// one on the order means the caller is a step behind, and closing the order on
// that basis would claim shares nobody proved were minted.
func TestCompleteDeliveryRefusesAHashThatIsNotTheOneOnTheOrder(t *testing.T) {
	ctx := context.Background()
	env := newDeliveryEnv(t)
	order := env.seed(t, aPaidOrder())

	if _, err := StartDelivery(ctx, order.OrderRef); err != nil {
		t.Fatalf("StartDelivery: %v", err)
	}
	if err := RecordDeliveryTxHash(ctx, order.OrderID, aTxHash); err != nil {
		t.Fatalf("RecordDeliveryTxHash: %v", err)
	}

	other := nthTxHash(99)
	if err := CompleteDelivery(ctx, order.OrderRef, other); !errors.Is(err, ErrOrderNotDeliverable) {
		t.Fatalf("want ErrOrderNotDeliverable, got %v", err)
	}

	if got, _ := env.status(t, order.OrderRef); got != statusInProgress {
		t.Fatalf("order is %d, want %d: it must not have moved", got, statusInProgress)
	}
}
