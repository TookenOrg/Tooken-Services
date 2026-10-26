package services

import (
	"context"
	"errors"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/delivery/database"
)

// fakeMinter stands in for the blockchain module. It counts what Deliver asked
// of it, which is how these tests check the one thing that matters: a mint is
// attempted at most once per order, and never after a refusal.
type fakeMinter struct {
	sendCalls    int
	receiptCalls int

	txHash  string
	sent    bool
	sendErr error

	found      bool
	succeeded  bool
	receiptErr error

	panicOnSend bool
}

// SendMintShares always calls record, like the real one: the hash reaches the
// database before the network, and a record that fails stops the send.
func (f *fakeMinter) SendMintShares(_ context.Context, _, _ string, _ int,
	record func(txHash string) error) (string, bool, error) {

	f.sendCalls++

	if f.panicOnSend {
		panic("minter is not wired")
	}

	if record != nil {
		if err := record(f.txHash); err != nil {
			return f.txHash, false, err
		}
	}

	return f.txHash, f.sent, f.sendErr
}

func (f *fakeMinter) MintReceipt(_ context.Context, _, _ string, _ int) (bool, bool, error) {
	f.receiptCalls++
	return f.found, f.succeeded, f.receiptErr
}

// deliveryProbe collects what Deliver did to the order, so each test can state
// its expectation as "this happened, that did not".
type deliveryProbe struct {
	t *testing.T

	target    database.DeliveryTargetDTO
	startErr  error
	recordErr error

	recordOrderID int64
	recordHash    string

	failCalls      int
	failHash       *string
	failReached    bool
	completeCalls  int
	completeTxHash string
}

func (p *deliveryProbe) start(context.Context, string) (database.DeliveryTargetDTO, error) {
	return p.target, p.startErr
}

func (p *deliveryProbe) record(_ context.Context, orderID int64, txHash string) error {
	p.recordOrderID = orderID
	p.recordHash = txHash
	return p.recordErr
}

func (p *deliveryProbe) fail(_ context.Context, _, _ string, txHash *string, reachedNetwork bool) error {
	p.failCalls++
	p.failHash = txHash
	p.failReached = reachedNetwork
	return nil
}

func (p *deliveryProbe) complete(_ context.Context, _, txHash string) error {
	p.completeCalls++
	p.completeTxHash = txHash
	return nil
}

// run wires the probe into a Service and delivers one order.
func (p *deliveryProbe) run(minter TokenMinter, orderRef string) {
	p.t.Helper()

	svc := NewService(minter)
	svc.startDelivery = p.start
	svc.recordTxHash = p.record
	svc.failDelivery = p.fail
	svc.completeDelivery = p.complete

	svc.Deliver(context.Background(), orderRef)
}

func aTarget() database.DeliveryTargetDTO {
	return database.DeliveryTargetDTO{
		OrderID:      1,
		UserID:       7,
		Shares:       10,
		Wallet:       "0x1111111111111111111111111111111111111111",
		TokenAddress: "0x2222222222222222222222222222222222222222",
	}
}

// A delivery refused before any on-chain action must never reach
// the chain. StartDelivery has already written DELIVERY_FAILED by then, so
// Deliver only has to stop — and stopping is what must be proven, because a
// mint sent for a failed order is a mint nobody expects.
func TestDeliverDoesNotMintWhenTheOrderWasRefused(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{"refused before any on-chain action", database.ErrDeliveryRefused},
		{"order is not awaiting delivery", database.ErrOrderNotDeliverable},
		{"the database could not be reached", errors.New("connection refused")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			minter := &fakeMinter{}
			probe := &deliveryProbe{t: t, startErr: tc.err}

			probe.run(minter, "ISS-TEST-0001")

			if minter.sendCalls != 0 {
				t.Fatalf("SendMintShares was called %d time(s); a refused order must never be minted", minter.sendCalls)
			}
			if probe.failCalls != 0 {
				t.Fatalf("FailDelivery was called %d time(s); StartDelivery already settled this order", probe.failCalls)
			}
		})
	}
}

// When the hash cannot be stored, nothing may be broadcast. This is the
// rule the whole delivery rests on — "no hash in the database" has to keep
// meaning "nothing was ever sent".
func TestDeliverTreatsAFailedRecordAsNotSent(t *testing.T) {
	minter := &fakeMinter{txHash: "0xabc", sent: true}
	probe := &deliveryProbe{t: t, target: aTarget(), recordErr: errors.New("database is down")}

	probe.run(minter, "ISS-TEST-0002")

	if minter.receiptCalls != 0 {
		t.Fatal("a mint that was never sent must not be waited on")
	}
	if probe.failCalls != 1 {
		t.Fatalf("the order must be marked as failed, FailDelivery called %d time(s)", probe.failCalls)
	}
	if probe.failReached {
		t.Fatal("reachedNetwork must be false: the send never happened")
	}
}

// Two deliveries race on the same order. The one that loses the
// hash must stop without touching the order — marking it failed would
// contradict the other one, which is busy minting.
func TestDeliverStopsQuietlyWhenAnotherDeliveryOwnsTheMint(t *testing.T) {
	minter := &fakeMinter{txHash: "0xabc", sent: true}
	probe := &deliveryProbe{t: t, target: aTarget(), recordErr: database.ErrDeliveryAlreadyRecorded}

	probe.run(minter, "ISS-TEST-0003")

	if probe.failCalls != 0 {
		t.Fatalf("FailDelivery was called %d time(s); the other delivery owns this mint", probe.failCalls)
	}
	if probe.completeCalls != 0 {
		t.Fatalf("CompleteDelivery was called %d time(s); this delivery minted nothing", probe.completeCalls)
	}
}

// A mint whose fate is unknown is undecided, not failed. Failing it is
// precisely what would let a second mint start, while the first one may still
// land a minute later.
func TestDeliverLeavesAnUndecidedMintAlone(t *testing.T) {
	cases := []struct {
		name    string
		found   bool
		receipt error
	}{
		{"not mined before the wait ran out", false, nil},
		{"the node could not be reached", false, errors.New("connection reset")},
		// MintReceipt answers this way when the receipt says success but the
		// token cannot be read back: the shares may well exist.
		{"mined, but the proof could not be read", true, errors.New("token lookup failed")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			minter := &fakeMinter{txHash: "0xabc", sent: true, found: tc.found, receiptErr: tc.receipt}
			probe := &deliveryProbe{t: t, target: aTarget()}

			probe.run(minter, "ISS-TEST-0004")

			if probe.failCalls != 0 {
				t.Fatalf("FailDelivery was called %d time(s); an undecided mint has not failed", probe.failCalls)
			}
			if probe.completeCalls != 0 {
				t.Fatalf("CompleteDelivery was called %d time(s); nothing proves the shares exist", probe.completeCalls)
			}
		})
	}
}

// The send was refused and the node never saw the transaction. The hash
// goes on the timeline, but reachedNetwork is false so the order column is
// cleared: it must keep answering "what was broadcast?" truthfully.
func TestDeliverKeepsTheHashOnTheTimelineWhenNothingWasBroadcast(t *testing.T) {
	minter := &fakeMinter{
		txHash:  "0xabc",
		sent:    false,
		sendErr: errors.New("replacement transaction underpriced"),
	}
	probe := &deliveryProbe{t: t, target: aTarget()}

	probe.run(minter, "ISS-TEST-0005")

	if probe.failCalls != 1 {
		t.Fatalf("the order must be marked as failed, FailDelivery called %d time(s)", probe.failCalls)
	}
	if probe.failHash == nil || *probe.failHash != "0xabc" {
		t.Fatalf("the timeline must keep the attempt, got %v", probe.failHash)
	}
	if probe.failReached {
		t.Fatal("reachedNetwork must be false: the node never saw the transaction")
	}
}

// The nominal path. The mint is mined and carries its Transfer, so the
// order is completed with the very hash that was broadcast.
func TestDeliverCompletesTheOrderOnASuccessfulMint(t *testing.T) {
	minter := &fakeMinter{txHash: "0xabc", sent: true, found: true, succeeded: true}
	probe := &deliveryProbe{t: t, target: aTarget()}

	probe.run(minter, "ISS-TEST-0006")

	if probe.completeCalls != 1 {
		t.Fatalf("the order must be completed once, got %d call(s)", probe.completeCalls)
	}
	if probe.completeTxHash != "0xabc" {
		t.Fatalf("the order must be completed with the broadcast hash, got %q", probe.completeTxHash)
	}
	if probe.failCalls != 0 {
		t.Fatal("a successful mint must not mark the order as failed")
	}
	if minter.sendCalls != 1 {
		t.Fatalf("exactly one mint must be sent, got %d", minter.sendCalls)
	}
	if probe.recordOrderID != aTarget().OrderID || probe.recordHash != "0xabc" {
		t.Fatalf("the hash must be recorded on this order: got order %d, hash %q", probe.recordOrderID, probe.recordHash)
	}
}

// Deliver runs in a goroutine of its own: a panic there would take the whole
// server down, not just this delivery. It must be contained, and the order
// left as it is, since nobody knows how far the mint went.
func TestDeliverContainsAPanic(t *testing.T) {
	minter := &fakeMinter{panicOnSend: true}
	probe := &deliveryProbe{t: t, target: aTarget()}

	defer func() {
		if p := recover(); p != nil {
			t.Fatalf("Deliver let a panic escape: %v", p)
		}
	}()

	probe.run(minter, "ISS-TEST-0008")

	if probe.failCalls != 0 || probe.completeCalls != 0 {
		t.Fatalf("after a panic the order must be left alone, got %d fail and %d complete", probe.failCalls, probe.completeCalls)
	}
}

// A mint that was mined but created nothing is a failure the chain decided,
// and the transaction does exist: its hash must stay on the order.
func TestDeliverFailsAMinedMintThatCreatedNothing(t *testing.T) {
	minter := &fakeMinter{txHash: "0xabc", sent: true, found: true, succeeded: false}
	probe := &deliveryProbe{t: t, target: aTarget()}

	probe.run(minter, "ISS-TEST-0007")

	if probe.failCalls != 1 {
		t.Fatalf("a mined mint that created nothing must fail the order, got %d call(s)", probe.failCalls)
	}
	if !probe.failReached {
		t.Fatal("reachedNetwork must be true: the transaction is on the chain")
	}
	if probe.completeCalls != 0 {
		t.Fatal("an order whose mint created nothing must not be closed")
	}
}
