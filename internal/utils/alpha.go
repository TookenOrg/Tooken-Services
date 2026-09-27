package utils

import (
	"fmt"

	"golang.org/x/text/language"
)

var ErrInvalidCountryCode = fmt.Errorf("invalid country code")

func IsAlphaCountryCode(s string) bool {
	if len(s) != 2 {
		return false
	}
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}

func CountryAlpha2ToNumeric(alpha2 string) (int, error) {

	isAlpha2 := IsAlphaCountryCode(alpha2)
	if !isAlpha2 {
		return 0, ErrInvalidCountryCode
	}

	region, err := language.ParseRegion(alpha2)
	if err != nil {
		return 0, err
	}
	if !region.IsCountry() {
		return 0, ErrInvalidCountryCode
	}
	countryCode := region.M49()
	if countryCode == 0 {
		return 0, ErrInvalidCountryCode
	}
	return countryCode, nil
}
