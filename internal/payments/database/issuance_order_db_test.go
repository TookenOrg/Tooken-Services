package database

import (
	"errors"
	"fmt"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/locks"
)

// TestCheckOrderForPaymentConfirmation covers the decision taken under the
// lock. The SQL that feeds it (reservation_expired) is covered by the
// integration tests.
func TestCheckOrderForPaymentConfirmation(t *testing.T) {
	t.Run("awaiting payment, reservation running: accepted", func(t *testing.T) {
		if err := checkOrderForPaymentConfirmation(locks.LockedIssuanceOrderDTO{StatusId: 2}); err != nil {
			t.Fatalf("error = %v, want nil", err)
		}
	})

	t.Run("awaiting payment, reservation lapsed: expired", func(t *testing.T) {
		err := checkOrderForPaymentConfirmation(locks.LockedIssuanceOrderDTO{StatusId: 2, ReservationExpired: true})
		if !errors.Is(err, ErrOrderReservationExpired) {
			t.Fatalf("error = %v, want %v", err, ErrOrderReservationExpired)
		}
	})

	// The status is judged before the expiry: reservation_expires_at is never
	// cleared, so a paid or cancelled order keeps a lapsed date forever. Every
	// status other than 2 must say "not awaiting payment", lapsed or not.
	// Status 4 included: the "already paid" exit happens before this function.
	for _, status := range []int{1, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14} {
		for _, expired := range []bool{false, true} {
			t.Run(fmt.Sprintf("status %d, expired %v: not awaiting payment", status, expired), func(t *testing.T) {
				err := checkOrderForPaymentConfirmation(locks.LockedIssuanceOrderDTO{StatusId: status, ReservationExpired: expired})
				if !errors.Is(err, ErrOrderNotAwaitingPayment) {
					t.Fatalf("error = %v, want %v", err, ErrOrderNotAwaitingPayment)
				}
			})
		}
	}
}

// TestSentinels pins the identities the handler relies on: the 404 of this
// package must be the one of locks, otherwise errors.Is fails and a missing
// order answers 500.
func TestSentinels(t *testing.T) {
	if !errors.Is(ErrOrderNotFound, locks.ErrOrderNotFound) {
		t.Fatal("ErrOrderNotFound is not an alias of locks.ErrOrderNotFound")
	}
	if errors.Is(ErrOrderNotAwaitingPayment, ErrOrderReservationExpired) {
		t.Fatal("the two 409 sentinels must stay distinct")
	}
}
