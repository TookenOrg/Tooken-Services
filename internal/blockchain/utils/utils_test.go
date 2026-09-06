package utils

import (
	"math"
	"math/big"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/ethereum/go-ethereum/common"
)

func mustBig(t *testing.T, s string) *big.Int {
	t.Helper()
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		t.Fatalf("invalid big int %q", s)
	}
	return n
}

func TestConvertFloatToWei(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		decimals int64
		want     string // base-10 string; ignored when wantErr is true
		wantErr  bool
	}{
		{"one token, 18 decimals", 1, 18, "1000000000000000000", false},
		{"fractional, 2 decimals", 1.5, 2, "150", false},
		{"zero", 0, 18, "0", false},
		{"2.5 with 6 decimals", 2.5, 6, "2500000", false},
		{"integer with 0 decimals", 7, 0, "7", false},
		{"negative is rejected", -1, 18, "", true},

		// The amount ends up in an immutable ledger: every case below used to
		// return a plausible but wrong value instead of an error.
		{"finer than the token decimals", 1.005, 2, "", true},
		{"half a share on a whole-share token", 1.5, 0, "", true},
		{"less than one share is not a free mint", 0.9, 0, "", true},
		{"beyond 2^53 a unit is no longer representable", 1e16 + 1, 18, "", true},
		{"more significant digits than a float64 carries", 1.234567890123456789, 18, "", true},
		{"negative decimals would divide instead of scaling", 100, -2, "", true},

		// The boundary itself stays valid: refusing too much would be its own
		// kind of bug.
		{"exactly 15 significant digits", 999999999999999, 18, "999999999999999000000000000000000", false},
		{"smallest unit of an 18 decimals token", 0.000000000000000001, 18, "1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertFloatToWei(tt.amount, tt.decimals)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error for a negative amount")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got == nil {
				t.Fatalf("expected %s, got nil", tt.want)
			}
			if got.Cmp(mustBig(t, tt.want)) != 0 {
				t.Fatalf("want %s, got %s", tt.want, got.String())
			}
		})
	}
}

// NaN and Inf used to panic inside decimal.NewFromFloat. Mint and burn run in
// a goroutine, so that panic would not have been a failed request: it would
// have taken the whole server down.
func TestConvertFloatToWeiRejectsNonFiniteAmounts(t *testing.T) {
	for name, amount := range map[string]float64{
		"NaN":  math.NaN(),
		"+Inf": math.Inf(1),
		"-Inf": math.Inf(-1),
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panicked instead of returning an error: %v", r)
				}
			}()

			if _, err := ConvertFloatToWei(amount, 18); err == nil {
				t.Fatalf("expected an error for %s", name)
			}
		})
	}
}

func TestFindContractByName(t *testing.T) {
	addr := "0x1111111111111111111111111111111111111111"
	list := []server.ContractDetails{
		{Name: "A", Address: "0x2222222222222222222222222222222222222222"},
		{Name: "TARGET", Address: addr},
	}

	got := FindContractByName(list, "TARGET")
	if got == nil {
		t.Fatal("expected to find TARGET")
	}
	if *got != common.HexToAddress(addr) {
		t.Fatalf("want %s, got %s", addr, got.Hex())
	}

	if FindContractByName(list, "MISSING") != nil {
		t.Fatal("expected nil for a missing name")
	}
}
