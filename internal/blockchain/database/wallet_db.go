package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/ethereum/go-ethereum/common"
)

func GetWalletByUserId(ctx context.Context, userId int) (walletPubKey *common.Address, err error) {
	var hex string
	err = globals.DB.QueryRowContext(ctx, `
    SELECT wallet_address FROM blk.user_wallet
    WHERE user_id = $1 AND is_active
`, userId).Scan(&hex)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get wallet by user id: %w", err)
	}
	addr := common.HexToAddress(hex)
	return &addr, nil
}

func InsertWallet(userID int, walletAddress, label string) (id int64, err error) {
	query := `
        INSERT INTO blk.user_wallet
            (user_id, wallet_address, label)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	err = globals.DB.QueryRow(query, userID, walletAddress, label).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert wallet: %w", err)
	}

	fmt.Printf("Wallet inserted with id: %d\n", id)
	return
}
