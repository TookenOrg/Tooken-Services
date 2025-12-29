package database

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

func FindModuleByName(ctx context.Context, complianceModuleName string) (moduleAddr common.Address, found bool, err error) {

	moduleDetails, err := GetContractRoleByName(ctx, complianceModuleName)
	if err != nil {
		return
	}

	found = true
	moduleAddr = common.HexToAddress(moduleDetails.Address)

	return
}
