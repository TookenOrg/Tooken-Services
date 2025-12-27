package database

import (
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/ethereum/go-ethereum/common"
)

func GetIdentityRegistryInstance() (irInstance *contracts.IdentityRegistry, err error) {

	// TODO: MOCK
	// https://github.com/TookenOrg/tooken-services/issues/13
	trexSuiteIr := common.HexToAddress("0x3bD5E6f6253e177e8c3E6ED24F51bEA358343BD0")

	irInstance, err = contracts.NewIdentityRegistry(trexSuiteIr, globals.EthClient)
	if err != nil {
		return
	}

	return
}
