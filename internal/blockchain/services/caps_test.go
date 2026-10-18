package services

import (
	"bytes"
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// B-u1, B-u3: a valid number of shares is scaled exactly, however large it is.
func TestSharesToBaseUnitsScalesExactly(t *testing.T) {
	cases := []struct {
		shares   string
		decimals int
		want     string
	}{
		{"50", 18, "50000000000000000000"},
		{"1500", 18, "1500000000000000000000"},
		{"50", 6, "50000000"},
		{"50", 0, "50"},
		// Leading zeros are still a whole number: the value, not the spelling, counts.
		{"007", 2, "700"},
		// 20 digits: far beyond int64 and beyond what a float64 holds exactly.
		{"99999999999999999999", 18, "99999999999999999999000000000000000000"},
	}
	for _, c := range cases {
		got, err := sharesToBaseUnits(c.shares, c.decimals)
		if err != nil {
			t.Errorf("sharesToBaseUnits(%q, %d): unexpected error %v", c.shares, c.decimals, err)
			continue
		}
		if got.String() != c.want {
			t.Errorf("sharesToBaseUnits(%q, %d) = %s, want %s", c.shares, c.decimals, got, c.want)
		}
	}
}

// B-u2: anything that is not a strictly positive whole number is refused. A cap of
// zero blocks every mint (B2), and a share cannot be split.
func TestSharesToBaseUnitsRefusesWhatIsNotAWholePositiveNumber(t *testing.T) {
	for _, shares := range []string{"", "0", "000", "-3", "+3", "1.5", "1.0", "1e3", "1_000", " 5", "5 ", "abc", "0x10"} {
		if got, err := sharesToBaseUnits(shares, 18); err == nil {
			t.Errorf("sharesToBaseUnits(%q) must be refused, got %s", shares, got)
		}
	}
	if got, err := sharesToBaseUnits("5", -1); err == nil {
		t.Errorf("negative decimals must be refused, got %s", got)
	}
}

// B-u4: the two settings are calls to the right functions, carrying the value given.
// A wrong selector would not fail here at encoding time: the factory would replay it
// on the module and the whole suite deployment would revert.
func TestEncodeCapSettings(t *testing.T) {
	perInvestor, _ := new(big.Int).SetString("50000000000000000000", 10)
	supply, _ := new(big.Int).SetString("1500000000000000000000", 10)

	maxBalance, supplyLimit, err := encodeCapSettings(perInvestor, supply)
	if err != nil {
		t.Fatalf("encodeCapSettings: %v", err)
	}

	cases := []struct {
		name      string
		call      []byte
		signature string
		want      *big.Int
	}{
		{"maxBalance", maxBalance, "setMaxBalance(uint256)", perInvestor},
		{"supplyLimit", supplyLimit, "setSupplyLimit(uint256)", supply},
	}
	for _, c := range cases {
		selector := crypto.Keccak256([]byte(c.signature))[:4]
		if len(c.call) != 4+32 {
			t.Fatalf("%s: a single uint256 argument encodes to 36 bytes, got %d", c.name, len(c.call))
		}
		if !bytes.Equal(c.call[:4], selector) {
			t.Fatalf("%s: selector %x, want %x (%s)", c.name, c.call[:4], selector, c.signature)
		}
		if got := new(big.Int).SetBytes(c.call[4:]); got.Cmp(c.want) != 0 {
			t.Fatalf("%s: argument %s, want %s", c.name, got, c.want)
		}
	}

	// And the selectors are the ones of the generated bindings, so they do exist on
	// the deployed modules — not merely a signature typed the same way twice.
	mbABI, err := contracts.MaxBalanceModuleMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}
	slABI, err := contracts.SupplyLimitModuleMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(maxBalance[:4], mbABI.Methods["setMaxBalance"].ID) {
		t.Fatal("setMaxBalance selector does not match the MaxBalanceModule binding")
	}
	if !bytes.Equal(supplyLimit[:4], slABI.Methods["setSupplyLimit"].ID) {
		t.Fatal("setSupplyLimit selector does not match the SupplyLimitModule binding")
	}
}

// B-u5: buildTokenDetails hands modules and settings to the factory untouched. The
// factory pairs settings[i] with modules[i], so a reordering or a dropped entry here
// would send a setter to the wrong module (B1).
func TestBuildTokenDetailsPassesSettingsThrough(t *testing.T) {
	owner := common.HexToAddress("0xAbC0000000000000000000000000000000000001")
	maxBalanceModule := common.HexToAddress("0xAbC0000000000000000000000000000000000004")
	supplyLimitModule := common.HexToAddress("0xAbC0000000000000000000000000000000000005")
	transferRestrictModule := common.HexToAddress("0xAbC0000000000000000000000000000000000006")

	maxBalance, supplyLimit, err := encodeCapSettings(big.NewInt(50), big.NewInt(1500))
	if err != nil {
		t.Fatal(err)
	}
	modules := []common.Address{maxBalanceModule, supplyLimitModule, transferRestrictModule}
	settings := [][]byte{maxBalance, supplyLimit}

	td := buildTokenDetails(owner, "Tooken Villa Belair #12", "TVB12", 18, common.Address{}, modules, settings)

	if len(td.ComplianceSettings) > len(td.ComplianceModules) {
		t.Fatalf("%d settings for %d modules: the factory would index past the modules", len(td.ComplianceSettings), len(td.ComplianceModules))
	}
	for i, m := range modules {
		if td.ComplianceModules[i] != m {
			t.Fatalf("module %d: got %s, want %s", i, td.ComplianceModules[i].Hex(), m.Hex())
		}
	}
	if len(td.ComplianceSettings) != 2 ||
		!bytes.Equal(td.ComplianceSettings[0], maxBalance) ||
		!bytes.Equal(td.ComplianceSettings[1], supplyLimit) {
		t.Fatal("settings must be passed through in order: setMaxBalance, then setSupplyLimit")
	}
}

// The share counts are checked before anything is read or sent. Nothing here sets a
// database or a node: should a check move after the first lookup, the test panics on
// the nil connection instead of answering an error, and that is the point. A bad
// input must never cost a module deployment (each one is paid gas).
func TestCreateTokenRefusesBadSharesBeforeAnyIO(t *testing.T) {
	base := server.CreateTokenRequest{TokenName: "Tooken Villa Belair #12", Symbol: "TVB12", NbDecimal: 18}

	cases := []struct {
		name        string
		total, cap  string
		wantMessage string
	}{
		{"no total shares", "", "50", "whole number"},
		{"fractional total shares", "1500.5", "50", "whole number"},
		{"zero total shares", "0", "50", "greater than zero"},
		{"no per-investor cap", "1500", "", "whole number"},
		{"zero per-investor cap", "1500", "0", "greater than zero"},
		{"negative per-investor cap", "1500", "-1", "whole number"},
		{"cap above the total", "1500", "1501", "must not exceed totalShares"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := base
			req.TotalShares, req.MaxSharesPerInvestor = c.total, c.cap

			_, _, err := NewService().CreateTokenWithSalt(context.Background(), req, "re-12")
			if err == nil {
				t.Fatal("must be refused")
			}
			if !strings.Contains(err.Error(), c.wantMessage) {
				t.Fatalf("error %q must mention %q", err, c.wantMessage)
			}
		})
	}
}
