package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
)

type ContractDTO struct {
	ID           int64     `json:"id"`
	TxHash       string    `json:"tx_hash"`
	Address      string    `json:"address"`
	ContractName string    `json:"contract_name"`
	CreatedAt    time.Time `json:"created_at"`
}

func GetAllContracts(ctx context.Context) (contractsDetails []server.ContractDetails, err error) {
	query := `
        SELECT id, tx_hash, address, contract_name, created_at 
		FROM blk.contract_role;`

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

func GetContractRoleByName(ctx context.Context, contractName string) (contract server.ContractDetails, err error) {

	query := `
        SELECT id, tx_hash, address, contract_name, created_at 
        FROM blk.contract_role
        WHERE contract_name = $1`

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

// InsertContractRole persists a singleton contract address (factory, claim issuer,
// compliance module, shared IRS, ...) keyed by contractName so it can be resolved
// later via GetContractRoleByName. Columns are inferred from the SELECT queries
// above (tx_hash, address, contract_name); id and created_at are expected to default.
func InsertContractRole(ctx context.Context, txHash, address, contractName string) (id int64, err error) {
	query := `
        INSERT INTO blk.contract_role
            (tx_hash, address, contract_name)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	err = globals.DB.QueryRow(query, txHash, address, contractName).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert contract role %s: %w", contractName, err)
	}

	return
}
