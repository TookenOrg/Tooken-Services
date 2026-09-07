package services

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestBuildTokenDetails(t *testing.T) {
	owner := common.HexToAddress("0xAbC0000000000000000000000000000000000001")
	irs := common.HexToAddress("0xAbC0000000000000000000000000000000000002")
	module := common.HexToAddress("0xAbC0000000000000000000000000000000000003")

	td := buildTokenDetails(owner, "MyToken", "MTK", 6, irs, []common.Address{module})

	// ERC-3643 conformance: the token ONCHAINID must be the zero address so the
	// TREX factory creates it via the IdFactory (never an EOA).
	if td.ONCHAINID != (common.Address{}) {
		t.Errorf("ONCHAINID must be zero, got %s", td.ONCHAINID.Hex())
	}
	if td.Owner != owner {
		t.Errorf("owner mismatch: got %s", td.Owner.Hex())
	}
	// The shared IRS must be passed through so the investor whitelist is shared.
	if td.Irs != irs {
		t.Errorf("irs must be passed through, got %s", td.Irs.Hex())
	}
	if td.Decimals != 6 {
		t.Errorf("decimals: want 6, got %d", td.Decimals)
	}
	if td.Name != "MyToken" || td.Symbol != "MTK" {
		t.Errorf("name/symbol mismatch: %q / %q", td.Name, td.Symbol)
	}
	if len(td.IrAgents) != 1 || td.IrAgents[0] != owner {
		t.Errorf("irAgents should be [owner]")
	}
	if len(td.TokenAgents) != 1 || td.TokenAgents[0] != owner {
		t.Errorf("tokenAgents should be [owner]")
	}
	if len(td.ComplianceModules) != 1 || td.ComplianceModules[0] != module {
		t.Errorf("compliance modules mismatch")
	}
}

func TestDefineClaimSuiteDetails(t *testing.T) {
	issuer := common.HexToAddress("0xAbC0000000000000000000000000000000000009")
	cd := defineClaimSuiteDetails(issuer)

	if len(cd.ClaimTopics) != 1 || cd.ClaimTopics[0].Cmp(big.NewInt(7)) != 0 {
		t.Errorf("expected a single claim topic 7 (KYC)")
	}
	if len(cd.Issuers) != 1 || cd.Issuers[0] != issuer {
		t.Errorf("expected the shared claim issuer as trusted issuer")
	}
	if len(cd.IssuerClaims) != 1 || len(cd.IssuerClaims[0]) != 1 || cd.IssuerClaims[0][0].Cmp(big.NewInt(7)) != 0 {
		t.Errorf("issuerClaims must map the issuer to topic 7")
	}
}

func TestDefineAuthorityVersion(t *testing.T) {
	v := defineAuthorityVersion()
	if v.Major != 1 || v.Minor != 0 || v.Patch != 0 {
		t.Errorf("want version 1.0.0, got %d.%d.%d", v.Major, v.Minor, v.Patch)
	}
}

func TestControlInputMintBurn(t *testing.T) {
	valid := "0x1111111111111111111111111111111111111111"
	tests := []struct {
		name     string
		token    string
		to       string
		amount   float64
		decimals int64
		wantWei  string // base-10 string; empty means the input must be refused
	}{
		{"valid", valid, valid, 1.0, 18, "1000000000000000000"},
		{"bad token address", "not-an-address", valid, 1.0, 18, ""},
		{"bad recipient address", valid, "0xZZZ", 1.0, 18, ""},
		{"negative amount", valid, valid, -1.0, 18, ""},

		// An amount the token cannot represent must be refused, not rounded:
		// the mint is irreversible and the ledger has to stay exact.
		{"finer than the token decimals", valid, valid, 1.005, 2, ""},
		{"would mint nothing at all", valid, valid, 0.9, 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWei, err := controlInputMintBurn(tt.token, tt.to, tt.amount, tt.decimals)

			if tt.wantWei == "" {
				if err == nil {
					t.Fatalf("expected a refusal, got %v wei", gotWei)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected refusal: %v", err)
			}
			if gotWei.String() != tt.wantWei {
				t.Errorf("want %s wei, got %s", tt.wantWei, gotWei)
			}
		})
	}
}
