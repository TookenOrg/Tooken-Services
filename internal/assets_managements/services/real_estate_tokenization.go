package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/TookenOrg/tooken-services/internal/api/server"
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

// tokenForPublication returns the blk.token id the asset must be bound to,
// deploying a T-REX suite only when it does not already carry one.
//
// An asset that already carries a token is handed back untouched. Republishing
// must not spend a second salt: the binding it holds is the one investors own,
// and the factory would refuse the deployment anyway.
//
// Deploying here, before the asset is written as published, is what makes a
// chain failure harmless: the asset stays a draft, and no investor is ever
// shown an offer backed by nothing.
func (s *Service) tokenForPublication(ctx context.Context, realEstateID int, realEstateName string, existingTokenID *int) (int, error) {
	if existingTokenID != nil {
		return *existingTokenID, nil
	}

	salt := defineSalt(realEstateID)
	token, err := s.deployer.GetTokenBySalt(ctx, salt, true)
	if err != nil {
		return 0, err
	}
	if token != nil {
		return token.Id, nil
	}

	req := server.CreateTokenRequest{
		TokenName: defineTokenName(realEstateName, realEstateID),
		Symbol:    defineSymbol(realEstateID),
		NbDecimal: tokenDecimals,
	}

	tokenID, _, err := s.deployer.CreateTokenWithSalt(ctx, req, salt)
	if err != nil {
		return 0, err
	}

	return tokenID, nil
}
