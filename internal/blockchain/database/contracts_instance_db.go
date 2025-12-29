package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
)

type ContractInstanceDTO struct {
	ID                       int64     `json:"id"`
	Address                  string    `json:"address"`
	ContractName             string    `json:"contract_name"`
	ContractImplementationId int64     `json:"contract_implementation_id"`
	ParentContractId         int64     `json:"parent_contract_id"`
	CreatedAt                time.Time `json:"created_at"`
}

func GetContractInstanceByName(ctx context.Context, contractName string) (contract server.ContractDetails, err error) {

	query := `
        SELECT id, address, contract_name, contract_implementation_id, parent_contract_id, created_at 
        FROM blk.contract_instance
        WHERE contract_name = $1`

	var contractDB ContractInstanceDTO
	err = globals.DB.QueryRow(query, contractName).Scan(
		&contractDB.ID,
		&contractDB.Address,
		&contractDB.ContractName,
		&contractDB.ContractImplementationId,
		&contractDB.ParentContractId,
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
