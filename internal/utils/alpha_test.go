package utils

import "testing"

func TestCountryAlpha2ToNumeric(t *testing.T) {
	tests := []struct {
		name  string
		alpha string
		want  int
	}{
		{name: "France", alpha: "FR", want: 250},
		{name: "Luxembourg", alpha: "LU", want: 442},
		{name: "Belgium", alpha: "BE", want: 56},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CountryAlpha2ToNumeric(tt.alpha)
			if err != nil {
				t.Fatalf("CountryAlpha2ToNumeric(%q) error = %v", tt.alpha, err)
			}
			if got != tt.want {
				t.Fatalf("CountryAlpha2ToNumeric(%q) = %d, want %d", tt.alpha, got, tt.want)
			}
		})
	}
}

func TestCountryAlpha2ToNumericRejectsInvalidCountries(t *testing.T) {
	tests := []struct {
		name  string
		alpha string
	}{
		{name: "empty", alpha: ""},
		{name: "unknown alpha 2", alpha: "ZZ"},
		{name: "non country region", alpha: "EU"},
		{name: "too long", alpha: "FRA"},
		{name: "numeric", alpha: "12"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CountryAlpha2ToNumeric(tt.alpha)
			if err == nil {
				t.Fatalf("CountryAlpha2ToNumeric(%q) = %d, want error", tt.alpha, got)
			}
			if got != 0 {
				t.Fatalf("CountryAlpha2ToNumeric(%q) returned %d with error, want 0", tt.alpha, got)
			}
		})
	}
}
