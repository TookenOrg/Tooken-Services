package utils

import (
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
