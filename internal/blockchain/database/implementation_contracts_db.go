package database

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/api/server"
)

func GetAllContractImplementations(ctx context.Context) ([]server.ContractDetails, error) {
	// Mock données
	return []server.ContractDetails{
		{
			Name:    "TREXImplementationAuthority",
			Address: "0xAuthorityAddress",
		},
	}, nil
}
