package utils

import (
	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/ethereum/go-ethereum/common"
)

func FindContractByName(contracts []server.ContractDetails, name string) *common.Address {
	for _, contract := range contracts {
		if contract.Name == name {
			addr := common.HexToAddress(contract.Address)
			return &addr
		}
	}
	return nil
}
