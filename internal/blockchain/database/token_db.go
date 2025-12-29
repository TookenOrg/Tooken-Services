package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
)

func GetTokenByName(ctx context.Context, tokenName string, empryResultAllowed bool) (token *server.TokenInfos, err error) {
	query := `
        SELECT id, symbol, token_name,  address, nb_decimal, modular_compliance_addr, created_at
		FROM blk.token
        WHERE token_name = $1
    `

	// On crée une instance vide
	tokenRow := &server.TokenInfos{}

	// QueryRow().Scan renvoie sql.ErrNoRows si aucune ligne trouvée
	err = globals.DB.QueryRow(query, tokenName).Scan(
		&tokenRow.Id,
		&tokenRow.Symbol,
		&tokenRow.TokenName,
		&tokenRow.Address,
		&tokenRow.NbDecimal,
		&tokenRow.ModularComplianceAddr,
		&tokenRow.CreatedAt,
	)
	if err != nil {
		if empryResultAllowed && err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to select token for tokenName %s: %w", tokenName, err)
	}

	logger.LogDebug("Token found for tokenName %s: %s", tokenName, tokenRow.Address)

	token = tokenRow
	return
}
