package services

import (
	"errors"
	"strings"
	"testing"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/assets_managements/database"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

func strPtrT(v string) *string { return &v }
func intPtrT(v int) *int       { return &v }
func boolPtrT(v bool) *bool    { return &v }

func validAddress() server.RealEstateAddress {
	return server.RealEstateAddress{
		Street:      "12 rue de la Gare",
		PostalCode:  "L-1611",
		City:        "Luxembourg",
		CountryCode: "LU",
	}
}

func validRequest() server.RealEstateWriteRequest {
	return server.RealEstateWriteRequest{
		Title:        "Villa Belair",
		EstateTypeId: 1,
		Address:      validAddress(),
	}
}

// The payload is the only thing standing between a client and the registry:
// every value it can get wrong must be named, not silently coerced.
func TestToWriteDTORejectsInvalidPayloads(t *testing.T) {
	tests := []struct {
		name string
		req  func(*server.RealEstateWriteRequest)
	}{
		{"blank title", func(r *server.RealEstateWriteRequest) { r.Title = "   " }},
		{"missing street", func(r *server.RealEstateWriteRequest) { r.Address.Street = " " }},
		{"missing city", func(r *server.RealEstateWriteRequest) { r.Address.City = "" }},
		{"three letter country", func(r *server.RealEstateWriteRequest) { r.Address.CountryCode = "LUX" }},
		{"latitude without longitude", func(r *server.RealEstateWriteRequest) {
			r.Address.Latitude = strPtrT("49.6")
		}},
		{"latitude out of range", func(r *server.RealEstateWriteRequest) {
			r.Address.Latitude = strPtrT("91")
			r.Address.Longitude = strPtrT("6.1")
		}},
		{"latitude not a number", func(r *server.RealEstateWriteRequest) {
			r.Address.Latitude = strPtrT("north")
			r.Address.Longitude = strPtrT("6.1")
		}},
		{"fractional shares", func(r *server.RealEstateWriteRequest) {
			r.Configuration = &server.RealEstateConfigurationInput{TotalShares: "10.5", PricePerShare: "1"}
		}},
		{"zero shares", func(r *server.RealEstateWriteRequest) {
			r.Configuration = &server.RealEstateConfigurationInput{TotalShares: "0", PricePerShare: "1"}
		}},
		{"negative price", func(r *server.RealEstateWriteRequest) {
			r.Configuration = &server.RealEstateConfigurationInput{TotalShares: "10", PricePerShare: "-1"}
		}},
		{"price finer than the column", func(r *server.RealEstateWriteRequest) {
			r.Configuration = &server.RealEstateConfigurationInput{TotalShares: "10", PricePerShare: "1.123456789"}
		}},
		{"price not a number", func(r *server.RealEstateWriteRequest) {
			r.Configuration = &server.RealEstateConfigurationInput{TotalShares: "10", PricePerShare: "cheap"}
		}},
		{"fee above 100", func(r *server.RealEstateWriteRequest) {
			r.Configuration = &server.RealEstateConfigurationInput{
				TotalShares: "10", PricePerShare: "1", EntryFeeRate: strPtrT("101"),
			}
		}},
		{"min above max", func(r *server.RealEstateWriteRequest) {
			r.Configuration = &server.RealEstateConfigurationInput{
				TotalShares: "10", PricePerShare: "1",
				MinInvestment: strPtrT("500"), MaxInvestment: strPtrT("100"),
			}
		}},
		{"unknown currency", func(r *server.RealEstateWriteRequest) {
			r.Configuration = &server.RealEstateConfigurationInput{
				TotalShares: "10", PricePerShare: "1", CurrencyCode: strPtrT("EURO"),
			}
		}},
		{"negative surface", func(r *server.RealEstateWriteRequest) {
			r.Specification = &server.RealEstateSpecificationInput{SurfaceArea: strPtrT("-5")}
		}},
		{"negative bedrooms", func(r *server.RealEstateWriteRequest) {
			r.Specification = &server.RealEstateSpecificationInput{BedroomNumber: intPtrT(-1)}
		}},
		{"energy class out of scale", func(r *server.RealEstateWriteRequest) {
			r.Specification = &server.RealEstateSpecificationInput{EnergyClass: strPtrT("H")}
		}},
		{"media without url", func(r *server.RealEstateWriteRequest) {
			r.Media = &[]server.RealEstateMediaInput{{Url: "  "}}
		}},
		{"unknown media type", func(r *server.RealEstateWriteRequest) {
			r.Media = &[]server.RealEstateMediaInput{{Url: "a", MediaType: strPtrT("hologram")}}
		}},
		{"two covers", func(r *server.RealEstateWriteRequest) {
			r.Media = &[]server.RealEstateMediaInput{
				{Url: "a", IsCover: boolPtrT(true)},
				{Url: "b", IsCover: boolPtrT(true)},
			}
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validRequest()
			tt.req(&req)

			if _, err := toWriteDTO(req, writeCreate); !errors.Is(err, ErrInvalidRealEstate) {
				t.Errorf("want ErrInvalidRealEstate, got %v", err)
			}
		})
	}
}

// Normalisation is not cosmetic: the referential columns are CHAR(2), CHAR(3)
// and CHAR(1), and a lowercase code stored as-is would never match a lookup.
func TestToWriteDTONormalises(t *testing.T) {
	req := validRequest()
	req.Title = "  Villa Belair  "
	req.Description = strPtrT("   ")
	req.Address.CountryCode = "lu"
	req.Specification = &server.RealEstateSpecificationInput{EnergyClass: strPtrT(" c ")}
	req.Configuration = &server.RealEstateConfigurationInput{
		TotalShares: " 1500 ", PricePerShare: "199.99", CurrencyCode: strPtrT("eur"),
	}
	req.Media = &[]server.RealEstateMediaInput{{Url: "https://x/1.jpg"}, {Url: "https://x/2.jpg"}}

	in, err := toWriteDTO(req, writeCreate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if in.Title != "Villa Belair" {
		t.Errorf("title not trimmed: %q", in.Title)
	}
	// An empty description is no description: storing "" would create a second
	// way of saying "unknown", which every reader would then have to handle.
	if in.Description != nil {
		t.Errorf("blank description should be nil, got %q", *in.Description)
	}
	if in.Address.CountryCode != "LU" {
		t.Errorf("country code not upcased: %q", in.Address.CountryCode)
	}
	if *in.Specification.EnergyClass != "C" {
		t.Errorf("energy class not normalised: %q", *in.Specification.EnergyClass)
	}
	if in.SharesConfig.CurrencyCode != "EUR" {
		t.Errorf("currency not upcased: %q", in.SharesConfig.CurrencyCode)
	}
	if !in.SharesConfig.TotalShares.Equal(decimal.NewFromInt(1500)) {
		t.Errorf("total shares: %s", in.SharesConfig.TotalShares)
	}
	// The exact string must survive: 199.99 through a float64 comes back as
	// 199.99000000000000909.
	if in.SharesConfig.PricePerShare.String() != "199.99" {
		t.Errorf("price altered: %s", in.SharesConfig.PricePerShare)
	}
	// Positions default to the order the client sent, so the gallery does not
	// collapse into a pile of zeroes.
	if in.Media[0].Position != 0 || in.Media[1].Position != 1 {
		t.Errorf("positions not defaulted: %+v", in.Media)
	}
	if in.Media[0].MediaType != "image" {
		t.Errorf("media type not defaulted: %q", in.Media[0].MediaType)
	}
}

// A currency defaults to EUR, which is the database default too: the two must
// agree, otherwise the value depends on which layer wrote the row.
func TestToWriteDTODefaultsCurrency(t *testing.T) {
	req := validRequest()
	req.Configuration = &server.RealEstateConfigurationInput{TotalShares: "10", PricePerShare: "1"}

	in, err := toWriteDTO(req, writeCreate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if in.SharesConfig.CurrencyCode != "EUR" {
		t.Errorf("want EUR, got %q", in.SharesConfig.CurrencyCode)
	}
}

func nullDec(v string) decimal.NullDecimal {
	d, _ := decimal.NewFromString(v)
	return decimal.NullDecimal{Decimal: d, Valid: true}
}

func config(totalShares string) *database.RealEstateSharesConfigWriteDTO {
	d, _ := decimal.NewFromString(totalShares)
	return &database.RealEstateSharesConfigWriteDTO{TotalShares: d}
}

func TestCheckSharesConfigChange(t *testing.T) {
	tokenID := 42

	tests := []struct {
		name     string
		state    database.RealEstateGuardStateDTO
		in       *database.RealEstateSharesConfigWriteDTO
		conflict bool
	}{
		{
			name:  "untouched asset can be reconfigured freely",
			state: database.RealEstateGuardStateDTO{TotalShares: nullDec("1000")},
			in:    config("10"),
		},
		{
			name:     "tokenized asset freezes the total",
			state:    database.RealEstateGuardStateDTO{TokenId: &tokenID, TotalShares: nullDec("1000")},
			in:       config("1001"),
			conflict: true,
		},
		{
			name:  "tokenized asset accepts an unchanged total",
			state: database.RealEstateGuardStateDTO{TokenId: &tokenID, TotalShares: nullDec("1000")},
			in:    config("1000"),
		},
		{
			name:     "cannot shrink under the reserved shares",
			state:    database.RealEstateGuardStateDTO{ReservedShares: 400, TotalShares: nullDec("1000")},
			in:       config("399"),
			conflict: true,
		},
		{
			name:  "can shrink down to the reserved shares",
			state: database.RealEstateGuardStateDTO{ReservedShares: 400, TotalShares: nullDec("1000")},
			in:    config("400"),
		},
		{
			name:  "can always grow",
			state: database.RealEstateGuardStateDTO{ReservedShares: 400, TotalShares: nullDec("1000")},
			in:    config("5000"),
		},
		{
			name:     "cannot drop the configuration of a sold asset",
			state:    database.RealEstateGuardStateDTO{ReservedShares: 400, TotalShares: nullDec("1000")},
			in:       nil,
			conflict: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkSharesConfigChange(tt.state, tt.in)

			if tt.conflict && !errors.Is(err, ErrRealEstateConflict) {
				t.Errorf("want ErrRealEstateConflict, got %v", err)
			}
			if !tt.conflict && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// A create without an address is refused; a patch of an asset that never had
// one is not, otherwise an inherited record could never be corrected.
func TestAddressRequiredOnCreateOnly(t *testing.T) {
	req := validRequest()
	req.Address = server.RealEstateAddress{}

	if _, err := toWriteDTO(req, writeCreate); !errors.Is(err, ErrInvalidRealEstate) {
		t.Errorf("a create without an address must be refused, got %v", err)
	}

	in, err := toWriteDTO(req, writePatch)
	if err != nil {
		t.Fatalf("a patch without an address must be accepted, got %v", err)
	}
	if in.Address != nil {
		t.Errorf("no address should mean no row, got %+v", in.Address)
	}

	// A half-filled address is a mistake, not an omission: it stays refused
	// whatever the mode.
	req.Address.City = "Luxembourg"
	if _, err := toWriteDTO(req, writePatch); !errors.Is(err, ErrInvalidRealEstate) {
		t.Errorf("a partial address must be refused even on a patch, got %v", err)
	}
}

// applyPatch is where a field can silently disappear: every omission must be a
// no-op, and every supplied value must land.
func TestApplyPatch(t *testing.T) {
	current := validRequest()
	current.Description = strPtrT("Nice")
	current.Specification = &server.RealEstateSpecificationInput{SurfaceArea: strPtrT("128.50")}
	current.Configuration = &server.RealEstateConfigurationInput{
		TotalShares: "1500", PricePerShare: "199.99", CurrencyCode: strPtrT("EUR"),
	}
	current.Media = &[]server.RealEstateMediaInput{{Url: "a"}, {Url: "b"}}

	t.Run("an empty patch changes nothing", func(t *testing.T) {
		got := applyPatch(current, server.RealEstatePatchRequest{})

		if got.Title != current.Title || *got.Description != "Nice" {
			t.Errorf("scalars altered: %+v", got)
		}
		if got.Specification == nil || got.Configuration == nil || got.Media == nil {
			t.Errorf("sections dropped: %+v", got)
		}
		if got.Configuration.TotalShares != "1500" {
			t.Errorf("total shares altered: %q", got.Configuration.TotalShares)
		}
	})

	t.Run("a section patch merges field by field", func(t *testing.T) {
		got := applyPatch(current, server.RealEstatePatchRequest{
			Configuration: &server.RealEstateConfigurationPatch{PricePerShare: strPtrT("250")},
			Address:       &server.RealEstateAddressPatch{City: strPtrT("Esch")},
		})

		if got.Configuration.PricePerShare != "250" {
			t.Errorf("price not applied: %q", got.Configuration.PricePerShare)
		}
		if got.Configuration.TotalShares != "1500" {
			t.Errorf("total shares lost by a patch that never mentioned it: %q", got.Configuration.TotalShares)
		}
		if got.Address.City != "Esch" {
			t.Errorf("city not applied: %q", got.Address.City)
		}
		if got.Address.Street != current.Address.Street {
			t.Errorf("street lost by a patch that never mentioned it: %q", got.Address.Street)
		}
	})

	t.Run("a supplied gallery replaces the whole list", func(t *testing.T) {
		got := applyPatch(current, server.RealEstatePatchRequest{
			Media: &[]server.RealEstateMediaInput{{Url: "only"}},
		})

		if len(*got.Media) != 1 || (*got.Media)[0].Url != "only" {
			t.Errorf("gallery not replaced: %+v", *got.Media)
		}
	})

	t.Run("an empty string clears an optional value", func(t *testing.T) {
		got := applyPatch(current, server.RealEstatePatchRequest{Description: strPtrT("")})

		in, err := toWriteDTO(got, writePatch)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if in.Description != nil {
			t.Errorf("description not cleared: %q", *in.Description)
		}
	})
}

// A primary key collision comes from a sequence sitting behind its data, never
// from the payload. Reporting it as a conflict sent the caller looking for a
// duplicate value they never sent.
func TestTranslateWriteErrorTellsPrimaryKeyFromUserConflict(t *testing.T) {
	tests := []struct {
		name         string
		constraint   string
		wantSentinel error
	}{
		{"generated name", "real_estate_pkey", nil},
		{"handwritten name", "pk_real_estates", nil},
		{"a value the caller did send", "real_estate_shares_config_compartment_ref_uk", ErrRealEstateConflict},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := translateWriteError(&pq.Error{
				Code:       "23505",
				Constraint: tc.constraint,
				Detail:     "Key (id)=(1) already exists.",
			})

			if tc.wantSentinel != nil {
				if !errors.Is(err, tc.wantSentinel) {
					t.Fatalf("got %v, want a %v", err, tc.wantSentinel)
				}
				return
			}

			if errors.Is(err, ErrRealEstateConflict) || errors.Is(err, ErrInvalidRealEstate) {
				t.Fatalf("a sequence gap must not be reported to the caller as their mistake: %v", err)
			}
			if !strings.Contains(err.Error(), "000013") {
				t.Errorf("the message must point at the repair, got %q", err)
			}
		})
	}
}
