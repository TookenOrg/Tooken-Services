package models

import (
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/ethereum/go-ethereum/common"
)

type ComplianceSuite struct {
	Address  common.Address
	Instance *contracts.ModularCompliance
	Modules  []ComplianceModule
}

type ComplianceModule struct {
	Address     common.Address
	Name        string
	Description string
	IsShared    bool
}
