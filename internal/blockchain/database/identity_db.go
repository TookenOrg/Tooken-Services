package database

import (
	"context"
	"fmt"

	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
)

func GetIdentityAddrByUserId(ctx context.Context, userId int) (identityAddr *common.Address, err error) {
	query := `
        SELECT address
		FROM blk.identity
        WHERE user_id = $1
    `

	var identityAddress string
	err = globals.DB.QueryRow(query, userId).Scan(&identityAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to select identity address for userId %d: %w", userId, err)
	}

	logger.LogDebug("Identity found for userId %d: %s", userId, identityAddress)

	identityAddrCommon := common.HexToAddress(identityAddress)
	identityAddr = &identityAddrCommon
	return
}

func InsertIdentity(ctx context.Context, userId int, walletId int64, identityAddress, identityTxHash string) (err error) {
	query := `
        INSERT INTO blk.identity
            (user_id, wallet_id, address, tx_hash)
        VALUES ($1, $2, $3, $4)
		RETURNING id
    `

	var identityID int64
	err = globals.DB.QueryRow(query, userId, walletId, identityAddress, identityTxHash).Scan(&identityID)
	if err != nil {
		return fmt.Errorf("failed to insert Identity: %w", err)
	}

	logger.LogDebug("Identity inserted with id: %d\n", identityID)
	return
}
