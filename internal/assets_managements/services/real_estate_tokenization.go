package services

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const tokenDecimals = 18

// Visible name of the token based on the real estate name and ID
// We need to truncate the token name if it exceeds a certain length to comply with token name constraints.
func defineTokenName(realEstateName string, realEstateID int) string {
	const tookenCompanyName = "Tooken"
	const maxTokenNameOctets = 100
	const tokenNameFormat = "%s %s %s%s"
	const separatorsLength = 3
	const hashPrefix = "#"

	realEstateName = strings.TrimSpace(realEstateName)

	if realEstateName == "" {
		realEstateName = "Unnamed"
	}

	realEstateIDStr := strconv.Itoa(realEstateID)

	// Calculate the remaining size available for the real estate name to ensure the full token name does not exceed 100 octets.
	// 100 is the maximum allowed octets for the token name.
	// Subtract 3 for the spaces between the company name, real estate name, and the '#' character.
	remainingTitleSize := maxTokenNameOctets - len(tookenCompanyName) - len(realEstateIDStr) - separatorsLength

	cut := 0
	for i, r := range realEstateName {
		end := i + utf8.RuneLen(r)
		if end > remainingTitleSize {
			break
		}
		cut = end
	}
	realEstateName = realEstateName[:cut]

	return fmt.Sprintf(tokenNameFormat, tookenCompanyName, realEstateName, hashPrefix, realEstateIDStr)
}

func defineSymbol(realEstateID int) string {
	return fmt.Sprintf("TKN%d", realEstateID)
}

func defineSalt(realEstateID int) string {
	return fmt.Sprintf("re-%d", realEstateID)
}
