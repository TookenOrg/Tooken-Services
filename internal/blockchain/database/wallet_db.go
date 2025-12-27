package database

import (
	"context"
	"fmt"

	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/ethereum/go-ethereum/common"
)

func GetWalletByUserId(ctx context.Context, userId int) (walletPubKey *common.Address, err error) {
	// TODO MOCK
	walletPubKey = &common.MaxAddress
	return
}

func InsertWallet(userID int, walletAddress, label, privateKeyClear string) (id int64, err error) {
	encryptedKey, err := utils.EncryptAESGCM(privateKeyClear)
	if err != nil {
		return 0, fmt.Errorf("failed to encrypt private key: %w", err)
	}

	query := `
        INSERT INTO blk.user_wallet
            (user_id, wallet_address, label, private_key_clear, private_key_encrypted)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	err = globals.DB.QueryRow(query, userID, walletAddress, label, privateKeyClear, encryptedKey).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert wallet: %w", err)
	}

	fmt.Printf("Wallet inserted with id: %d\n", id)
	return
}
