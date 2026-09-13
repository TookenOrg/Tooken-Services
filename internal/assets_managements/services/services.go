package services

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/api/server"
)

type Service struct {
	deployer TokenDeployer
}

func NewService(tokenDeployer TokenDeployer) *Service {
	return &Service{
		deployer: tokenDeployer,
	}
}

// TokenDeployer is what the tokenisation needs from the blockchain module:
// deploying a T-REX suite under a given salt, and saying which blk.token row
// now carries it.
//
// It is declared here, on the calling side, rather than imported from the
// blockchain package. The point is testability: a fake deployer exercises the
// publication rule without a chain, a private key or a Hardhat node. It also
// states in one line the whole of what this package borrows from the chain —
// one method out of the twenty the blockchain service exposes.
type TokenDeployer interface {
	CreateTokenWithSalt(
		ctx context.Context,
		req server.CreateTokenRequest,
		salt string,
	) (int, server.TokenInfos, error)
}
