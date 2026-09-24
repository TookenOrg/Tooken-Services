package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
)

const contractVersion = 1

func InsertContractImplementation(ctx context.Context, txHash, address, contractName string) error {
	query := `
        INSERT INTO blk.contract_implementation (tx_hash, address, contract_name, version, created_at)
        VALUES ($1, $2, $3, $4, NOW())`

	_, err := globals.DB.ExecContext(ctx, query, txHash, address, contractName, contractVersion)
	if err != nil {
		return fmt.Errorf("insert failed: %w", err)
	}
	return nil
}
func GetAllContractImplementations(ctx context.Context) (contractsDetails []server.ContractDetails, err error) {
	query := `
        SELECT id, tx_hash, address, contract_name, created_at 
		FROM blk.contract_implementation;`

	rows, err := globals.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var contractsDB []ContractDTO
	for rows.Next() {
		var contractDB ContractDTO
		err := rows.Scan(
			&contractDB.ID,
			&contractDB.TxHash,
			&contractDB.Address,
			&contractDB.ContractName,
			&contractDB.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		contractsDB = append(contractsDB, contractDB)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	contractsDetails = make([]server.ContractDetails, len(contractsDB))
	for i, c := range contractsDB {
		contractsDetails[i] = server.ContractDetails{
			Address: c.Address,
			Name:    c.ContractName,
		}
	}

	return
}

func GetImplementationContractByName(ctx context.Context, contractName string) (contract server.ContractDetails, err error) {

	query := `
        SELECT id, tx_hash, address, contract_name, created_at 
        FROM blk.contract_implementation 
        WHERE contract_name = $1
        ORDER BY created_at DESC NULLS LAST
        LIMIT 1`

	var contractDB ContractDTO
	err = globals.DB.QueryRow(query, contractName).Scan(
		&contractDB.ID,
		&contractDB.TxHash,
		&contractDB.Address,
		&contractDB.ContractName,
		&contractDB.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return server.ContractDetails{}, fmt.Errorf("contract '%s' not found", contractName)
	}
	if err != nil {
		return server.ContractDetails{}, fmt.Errorf("query failed: %w", err)
	}

	contract.Address = contractDB.Address
	contract.Name = contractName

	return contract, nil

}
