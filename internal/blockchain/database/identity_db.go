package database

import (
	"context"
	"fmt"

	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/ethereum/go-ethereum/common"
)

func GetIdentityAddrByUserId(ctx context.Context, userId int) (identityAddr *common.Address, err error) {
	return &common.MaxAddress, nil
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
		return fmt.Errorf("failed to insert wallet: %w", err)
	}

	fmt.Printf("Identity inserted with id: %d\n", identityID)
	return
}
