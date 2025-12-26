package services

import (
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/ethereum/go-ethereum/common"
)

func SetGlobals() {

	trexSuite, err := database.GetTrexSuite()
	if err != nil {
	}

	// IdentityRegistryInstance
	globals.IdentityRegistryAddress = common.HexToAddress(trexSuite.IdentityRegistry.Address)
	globals.IdentityRegistryInstance, _ = contracts.NewIdentityRegistry(globals.IdentityRegistryAddress, globals.EthClient)

	// TokenInstance
}
