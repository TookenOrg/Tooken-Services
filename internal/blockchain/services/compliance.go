package services

import (
	"context"
	"errors"
	"fmt"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/models"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func deployComplianceSuite(ctx context.Context) (complianceSuite models.ComplianceSuite, err error) {

	// 1 - Deploy modular compliance
	complianceAddr, _, complianceInstancePtr, err := deployModularCompliance(ctx)
	if err != nil {
		return
	}

	// 2 - Add transfer restrict module
	complianceModule, err := addTransferRestrictModuleInCompliance(ctx, complianceInstancePtr)
	if err != nil {
		return
	}

	complianceSuite = models.ComplianceSuite{
		Address:  complianceAddr,
		Instance: complianceInstancePtr,
		Modules: []models.ComplianceModule{
			complianceModule,
		},
	}

	return
}

func deployModularCompliance(ctx context.Context) (modularComplianceAddr common.Address, tx *types.Transaction, modularComplianceInstance *contracts.ModularCompliance, err error) {

	trexAuthAddr, err := retreiveTrexAuthorityAddress(ctx)
	if err != nil {
		return
	}

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Deploy modular compliance Proxy...")

	// modularComplianceAddr, not modularComplianceProxyAddr
	modularComplianceAddr, tx, _, err = contracts.DeployModularComplianceProxy(auth, globals.EthClient, trexAuthAddr)
	err = utils.WaitDeployedTransaction(ctx, tx, true)
	if err != nil {
		return
	}
	logger.LogInfo("📬 Modular compliance deployed with transaction : %s", tx.Hash().Hex())

	modularComplianceInstance, err = contracts.NewModularCompliance(modularComplianceAddr, globals.EthClient)
	if err != nil {
		return
	}

	return
}

func addTransferRestrictModuleInCompliance(ctx context.Context, complianceInstance *contracts.ModularCompliance) (module models.ComplianceModule, err error) {
	// 1 - Retreive Transfer Restrict module
	moduleAddr, found, err := database.FindModuleByName(globals.TransferRestrictionModuleName)
	if err != nil {
		return
	}
	if !found {
		return models.ComplianceModule{}, fmt.Errorf("Module %s not found", globals.TransferRestrictionModuleName)
	}

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Add transfer module in compliance...")
	tx, err := complianceInstance.AddModule(auth, moduleAddr)
	if err != nil {
		return
	}
	err = utils.WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return
	}
	logger.LogInfo("📬 Module added in compliance with transaction: %s", moduleAddr.Hex())

	return models.ComplianceModule{
		Address:     moduleAddr,
		Name:        globals.TransferRestrictionModuleName,
		Description: "Restrict transfer between investors",
		IsShared:    true,
	}, nil
}

func DeploySharedComplianceModules(ctx context.Context, moduleName string) (moduleAddr common.Address, err error) {

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	switch moduleName {
	case globals.TransferRestrictionModuleName:
		moduleAddr, err = DeployTransferRestrictModule(ctx, auth)
	default:
		err = errors.New("compliance module not recognized: " + moduleName)
		return
	}

	return
}

func DeployTransferRestrictModule(ctx context.Context, auth *bind.TransactOpts) (moduleAddr common.Address, err error) {

	logger.LogInfo("💌 Deploying Transfer Restriction Module...")
	moduleAddr, tx, _, err := contracts.DeployTransferRestrictModule(auth, globals.EthClient)
	if err != nil {
		return
	}
	err = utils.WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return
	}
	logger.LogInfo("📬 Transfer Restriction Module deployed at address: %s", moduleAddr.Hex())
	return
}

func bindTokenToCompliance(ctx context.Context, modularComplianceInstance *contracts.ModularCompliance, tokenAddress common.Address) (tx *types.Transaction, err error) {

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Binding token %s to modular compliance", tokenAddress)
	tx, err = modularComplianceInstance.BindToken(auth, tokenAddress)
	if err = utils.WaitDeployedTransaction(ctx, tx, false); err != nil {
		return
	}
	logger.LogInfo("📬 Binding Token completed at address [%s]", tx.Hash().Hex())

	return
}
