package database

import (
	"context"
	"fmt"
	"math/big"

	"github.com/TookenOrg/tooken-services/internal/globals"
)

func InsertEthTransaction(ctx context.Context, txHash, txName, toAddress string, blockNumber int64, valueWei big.Int) (id int64, err error) {

	query := `
        INSERT INTO blk.eth_transaction
            (tx_hash, tx_name, block_number, to_address, value_wei)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	err = globals.DB.QueryRow(query, txHash, txName, blockNumber, toAddress, valueWei.String()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert transaction: %w", err)
	}

	fmt.Printf("Transaction inserted with id: %d\n", id)
	return
}
