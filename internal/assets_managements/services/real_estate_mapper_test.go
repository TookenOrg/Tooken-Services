package services

import (
	"encoding/json"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/assets_managements/database"
	"github.com/shopspring/decimal"
)

func TestRealEstateResponsesIncludeInvestmentTerms(t *testing.T) {
	exact := func(value string) decimal.NullDecimal {
		return decimal.NullDecimal{Decimal: decimal.RequireFromString(value), Valid: true}
	}
	dto := database.RealEstateDTO{
		SharesConfig: &database.RealEstateSharesConfigDTO{
			MinInvestment:  exact("199.99000001"),
			MaxInvestment:  exact("9007199254740993"),
			EntryFeeRate:   exact("1.125"),
			MgmtFeeRate:    exact("0.75"),
			ExitFeeRate:    exact("0"),
			CompartmentRef: strPtr("re-12"),
		},
	}
	for name, response := range map[string]any{
		"summary": toServerRealEstateSummary(dto),
		"detail":  toServerRealEstate(dto, false),
	} {
		t.Run(name, func(t *testing.T) {
			raw, err := json.Marshal(response)
			if err != nil {
				t.Fatal(err)
			}
			var body struct {
				Configuration map[string]json.RawMessage `json:"configuration"`
			}
			if err := json.Unmarshal(raw, &body); err != nil {
				t.Fatal(err)
			}
			for field, want := range map[string]string{
				"min_investment":  "199.99000001",
				"max_investment":  "9007199254740993",
				"entry_fee_rate":  "1.125",
				"mgmt_fee_rate":   "0.75",
				"exit_fee_rate":   "0",
				"compartment_ref": "re-12",
			} {
				var got string
				if err := json.Unmarshal(body.Configuration[field], &got); err != nil {
					t.Fatalf("%s is not a decimal string: %s (%v)", field, body.Configuration[field], err)
				}
				if got != want {
					t.Errorf("%s = %q, want %q", field, got, want)
				}
			}
		})
	}
}

func TestRealEstateResponsesOmitUnsetInvestmentTerms(t *testing.T) {
	for name, dto := range map[string]database.RealEstateDTO{
		"no configuration": {},
		"nullable terms":   {SharesConfig: &database.RealEstateSharesConfigDTO{}},
	} {
		t.Run(name, func(t *testing.T) {
			raw, err := json.Marshal(toServerRealEstateSummary(dto))
			if err != nil {
				t.Fatal(err)
			}
			var body struct {
				Configuration map[string]json.RawMessage `json:"configuration"`
			}
			if err := json.Unmarshal(raw, &body); err != nil {
				t.Fatal(err)
			}
			for _, field := range []string{"min_investment", "max_investment", "entry_fee_rate", "mgmt_fee_rate", "exit_fee_rate", "compartment_ref"} {
				if _, exists := body.Configuration[field]; exists {
					t.Errorf("unset %s was included", field)
				}
			}
		})
	}
}
