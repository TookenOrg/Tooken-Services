package database

import (
	"database/sql"
	"errors"
	"testing"
)

// TestCheckSharesAvailable covers the decision of step 5; the SQL that feeds
// it is covered by the integration tests.
func TestCheckSharesAvailable(t *testing.T) {
	noCap := sql.NullInt64{}
	cap50 := sql.NullInt64{Int64: 50, Valid: true}

	cases := []struct {
		name     string
		total    int64
		cap      sql.NullInt64
		onAsset  int64
		byUser   int64
		quantity int
		want     error
		message  string
	}{
		{"enough shares, no cap", 1500, noCap, 0, 0, 1000, nil, ""},
		{"exactly the last shares", 1500, noCap, 1495, 0, 5, nil, ""},
		{"one share too many", 1500, noCap, 1495, 0, 6, ErrOrderNotEnoughShares, "not enough shares left: only 5 shares left, 6 requested"},
		{"sold out", 1500, noCap, 1500, 0, 1, ErrOrderNotEnoughShares, "not enough shares left: only 0 shares left, 1 requested"},
		{"oversold never shows a negative stock", 1500, noCap, 1510, 0, 1, ErrOrderNotEnoughShares, "not enough shares left: only 0 shares left, 1 requested"},
		{"exactly at the cap", 1500, cap50, 45, 45, 5, nil, ""},
		{"above the cap", 1500, cap50, 45, 45, 10, ErrOrderStakeLimit, "investor stake limit reached: you already hold 45 of 50 shares allowed on this asset, 10 requested"},
		{"cap counts only this investor", 1500, cap50, 900, 0, 50, nil, ""},
		{"stock is checked before the cap", 1500, cap50, 1495, 45, 10, ErrOrderNotEnoughShares, "not enough shares left: only 5 shares left, 10 requested"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			asset := lockedRealEstateDTO{TotalShares: c.total, MaxSharesPerInvestor: c.cap}

			err := checkSharesAvailable(asset, c.onAsset, c.byUser, c.quantity)

			if c.want == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if !errors.Is(err, c.want) {
				t.Fatalf("error = %v, want %v", err, c.want)
			}
			if err.Error() != c.message {
				t.Errorf("message = %q, want %q", err.Error(), c.message)
			}
		})
	}
}
