package services

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// rate builds an entry fee rate; nil stands for an asset without entry fee.
func rate(s string) *decimal.Decimal {
	r := decimal.RequireFromString(s)
	return &r
}

// TestPriceOrder covers TICKET-M3-3 §5 (R1-R5) plus the rounding edge cases.
// Amounts are compared with StringFixed(2) because decimal.Equal ignores the
// scale and the API exposes exactly two decimals.
func TestPriceOrder(t *testing.T) {
	cases := []struct {
		name      string
		quantity  int64
		unitPrice string
		rate      *decimal.Decimal
		gross     string
		fee       string
		due       string
		withFee   bool
	}{
		{"R1 plain order with fee", 10, "200.00000000", rate("2.0000"), "2000.00", "40.00", "2040.00", true},
		{"R2 fee is computed on the rounded gross", 7, "199.99500000", rate("2.5"), "1399.97", "35.00", "1434.97", true},
		{"R3 asset without entry fee", 3, "100", nil, "300.00", "0.00", "300.00", false},
		{"R4 zero rate means no fee line", 3, "100", rate("0"), "300.00", "0.00", "300.00", false},
		{"R5 gross rounds half away from zero", 1, "0.12500000", rate("0"), "0.13", "0.00", "0.13", false},

		// 0.195 rounds to 0.20 first; 0.20 x 2.5 % = 0.005 -> 0.01.
		// Computing on the raw gross would give 0.004875 -> 0.00.
		{"fee on rounded gross flips the cent", 1, "0.195", rate("2.5"), "0.20", "0.01", "0.21", true},
		{"fee rounds half away from zero", 1, "1.00", rate("0.5"), "1.00", "0.01", "1.01", true},
		{"fee rounds down below half a cent", 1, "1.00", rate("0.4999"), "1.00", "0.00", "1.00", false},
		{"rate with four decimals", 1000, "1", rate("1.2345"), "1000.00", "12.35", "1012.35", true},
		{"positive rate rounding to zero fee writes no line", 1, "0.01", rate("2"), "0.01", "0.00", "0.01", false},
		{"rate above 100 %", 2, "50", rate("150"), "100.00", "150.00", "250.00", true},
		{"unit price with eight decimals", 3, "33.33333333", rate("1"), "100.00", "1.00", "101.00", true},
		{"gross rounding below half a cent", 1, "10.00499999", nil, "10.00", "0.00", "10.00", false},

		// Close to numeric(20,2): must stay exact, no float drift.
		{"large amounts stay exact", 1_000_000, "99999.99999999", rate("2"), "99999999999.99", "2000000000.00", "101999999999.99", true},

		// Defensive: the handler rejects these, priceOrder must not invent a fee.
		{"negative rate is ignored", 3, "100", rate("-1"), "300.00", "0.00", "300.00", false},
		{"zero quantity", 0, "100", rate("2"), "0.00", "0.00", "0.00", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			unitPrice := decimal.RequireFromString(c.unitPrice)
			var rateBefore string
			if c.rate != nil {
				rateBefore = c.rate.String()
			}

			got := priceOrder(c.quantity, unitPrice, c.rate)

			if got.GrossAmount.StringFixed(2) != c.gross {
				t.Errorf("GrossAmount = %s, want %s", got.GrossAmount.StringFixed(2), c.gross)
			}
			if got.FeeAmount.StringFixed(2) != c.fee {
				t.Errorf("FeeAmount = %s, want %s", got.FeeAmount.StringFixed(2), c.fee)
			}
			if got.AmountDue.StringFixed(2) != c.due {
				t.Errorf("AmountDue = %s, want %s", got.AmountDue.StringFixed(2), c.due)
			}
			if !got.UnitPrice.Equal(unitPrice) {
				t.Errorf("UnitPrice = %s, want it copied as is (%s)", got.UnitPrice, unitPrice)
			}
			if c.rate != nil && c.rate.String() != rateBefore {
				t.Errorf("entry fee rate was mutated: %s -> %s", rateBefore, c.rate)
			}

			// Invariants mirrored from the database CHECKs.
			if got.Fees == nil {
				t.Fatal("Fees is nil, want an empty slice (serialised as [] not null)")
			}
			if !got.AmountDue.Equal(got.GrossAmount.Add(got.FeeAmount)) {
				t.Errorf("AmountDue %s != GrossAmount %s + FeeAmount %s", got.AmountDue, got.GrossAmount, got.FeeAmount)
			}
			sum := decimal.Zero
			for _, f := range got.Fees {
				sum = sum.Add(f.Amount)
			}
			if !sum.Equal(got.FeeAmount) {
				t.Errorf("sum of fee lines %s != FeeAmount %s", sum, got.FeeAmount)
			}

			if !c.withFee {
				if len(got.Fees) != 0 {
					t.Errorf("Fees = %+v, want no fee line", got.Fees)
				}
				return
			}

			if len(got.Fees) != 1 {
				t.Fatalf("len(Fees) = %d, want 1", len(got.Fees))
			}
			f := got.Fees[0]
			if f.Code != "ENTRY" {
				t.Errorf("Code = %q, want ENTRY", f.Code)
			}
			if !f.Rate.Equal(*c.rate) {
				t.Errorf("Rate = %s, want %s", f.Rate, c.rate)
			}
			if !f.BaseAmount.Equal(got.GrossAmount) {
				t.Errorf("BaseAmount = %s, want the rounded gross %s", f.BaseAmount, got.GrossAmount)
			}
			// Same formula as the fees CHECK: round(base * rate / 100, 2).
			want := f.BaseAmount.Mul(f.Rate).Div(decimal.NewFromInt(100)).Round(2)
			if !f.Amount.Equal(want) {
				t.Errorf("Amount = %s, want round(base*rate/100, 2) = %s", f.Amount, want)
			}
		})
	}
}

func TestGenerateIssuanceOrderReference(t *testing.T) {
	const alphabet = "23456789ABCDEFGHJKLMNPQRSTVWXYZ"
	format := regexp.MustCompile(`^ISS-\d{8}-[` + alphabet + `]{6}$`)

	t.Run("format and date", func(t *testing.T) {
		at := time.Date(2026, time.October, 5, 23, 59, 0, 0, time.UTC)
		ref, err := generateIssuanceOrderReference(at)
		if err != nil {
			t.Fatal(err)
		}
		if !format.MatchString(ref) {
			t.Errorf("reference %q does not match %s", ref, format)
		}
		if !strings.HasPrefix(ref, "ISS-20261005-") {
			t.Errorf("reference %q does not carry the order date", ref)
		}
	})

	t.Run("no ambiguous characters", func(t *testing.T) {
		for i := 0; i < 500; i++ {
			ref, err := generateIssuanceOrderReference(time.Now())
			if err != nil {
				t.Fatal(err)
			}
			if suffix := ref[len(ref)-6:]; strings.ContainsAny(suffix, "01OIU") {
				t.Fatalf("suffix %q contains an ambiguous character", suffix)
			}
		}
	})

	// 31^6 ≈ 887 M combinations: 100 draws collide with probability ~6e-6.
	t.Run("suffix is random", func(t *testing.T) {
		at := time.Date(2026, time.October, 5, 0, 0, 0, 0, time.UTC)
		seen := make(map[string]bool)
		for i := 0; i < 100; i++ {
			ref, err := generateIssuanceOrderReference(at)
			if err != nil {
				t.Fatal(err)
			}
			if seen[ref] {
				t.Fatalf("duplicate reference %q", ref)
			}
			seen[ref] = true
		}
	})
}
