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

	tokenRow := &server.TokenInfos{}

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

func GetTokenByAddress(ctx context.Context, address string) (token *server.TokenInfos, err error) {
	query := `
        SELECT id, symbol, token_name,  address, nb_decimal, modular_compliance_addr, created_at
		FROM blk.token
        WHERE address = $1
    `

	tokenRow := &server.TokenInfos{}

	err = globals.DB.QueryRow(query, address).Scan(
		&tokenRow.Id,
		&tokenRow.Symbol,
		&tokenRow.TokenName,
		&tokenRow.Address,
		&tokenRow.NbDecimal,
		&tokenRow.ModularComplianceAddr,
		&tokenRow.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to select token for address %s: %w", address, err)
	}

	logger.LogDebug("Token found for address %s", address)

	token = tokenRow
	return
}

// InsertToken persists a newly created token. Columns are inferred from the SELECT
// queries above (symbol, token_name, address, nb_decimal, modular_compliance_addr);
// id and created_at are expected to default.
func InsertToken(ctx context.Context, symbol, tokenName, address string, nbDecimal int, modularComplianceAddr string) (id int, err error) {
	query := `
        INSERT INTO blk.token
            (symbol, token_name, address, nb_decimal, modular_compliance_addr)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	// blk.token.modular_compliance_addr is nullable with a CHECK that rejects any
	// value that is neither NULL nor a 0x-address — insert NULL instead of "".
	mcAddr := sql.NullString{String: modularComplianceAddr, Valid: modularComplianceAddr != ""}

	err = globals.DB.QueryRow(query, symbol, tokenName, address, nbDecimal, mcAddr).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert token %s: %w", tokenName, err)
	}

	return
}
