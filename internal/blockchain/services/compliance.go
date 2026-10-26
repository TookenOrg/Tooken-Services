package services

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func DeploySharedComplianceModules(ctx context.Context, moduleName string) (moduleAddr common.Address, err error) {

	switch moduleName {
	case globals.TransferRestrictionModuleName:
		moduleAddr, err = DeployTransferRestrictModule(ctx)
	case globals.MaxBalanceModuleName:
		moduleAddr, err = DeployMaxBalanceModule(ctx)
	case globals.SupplyLimitModuleName:
		moduleAddr, err = DeploySupplyLimitModule(ctx)
	default:
		err = errors.New("compliance module not recognized: " + moduleName)
		return
	}

	return
}

// DeployTransferRestrictModule deploys, initializes and registers the shared module
// that restricts transfers to allowed users.
//
// It goes through deployComplianceModule like the capping modules: the module is a
// UUPS implementation too, and was deployed without Initialize until 2026-10-18, so
// any address could claim its ownership (B8). For this module ownership grants
// nothing (writes are onlyComplianceCall, upgradeTo refuses a non-proxy call), but
// initializing it keeps every module in the same, owned state.
func DeployTransferRestrictModule(ctx context.Context) (common.Address, error) {
	return deployComplianceModule(ctx, globals.TransferRestrictionModuleName, "TRANSFER_RESTRICT_MODULE",
		func(auth *bind.TransactOpts) (common.Address, *types.Transaction, initializableModule, error) {
			return contracts.DeployTransferRestrictModule(auth, globals.EthClient)
		})
}

// initializableModule is what the generated bindings of the 4.1.6 compliance
// modules have in common: an Initialize that sets the module owner.
type initializableModule interface {
	Initialize(opts *bind.TransactOpts) (*types.Transaction, error)
}

// DeployMaxBalanceModule deploys, initializes and registers the shared module
// that caps the balance one investor identity can hold in a token.
func DeployMaxBalanceModule(ctx context.Context) (common.Address, error) {
	return deployComplianceModule(ctx, globals.MaxBalanceModuleName, "MAX_BALANCE_MODULE",
		func(auth *bind.TransactOpts) (common.Address, *types.Transaction, initializableModule, error) {
			return contracts.DeployMaxBalanceModule(auth, globals.EthClient)
		})
}

// DeploySupplyLimitModule deploys, initializes and registers the shared module
// that caps the total supply of a token.
func DeploySupplyLimitModule(ctx context.Context) (common.Address, error) {
	return deployComplianceModule(ctx, globals.SupplyLimitModuleName, "SUPPLY_LIMIT_MODULE",
		func(auth *bind.TransactOpts) (common.Address, *types.Transaction, initializableModule, error) {
			return contracts.DeploySupplyLimitModule(auth, globals.EthClient)
		})
}

// deployComplianceModule deploys a module, initializes it, then registers it
// under name in blk.contract_role.
//
// Registration comes last, so a registered module is always an initialized one,
// and a failure is returned rather than logged: an address that is not recorded
// is never found again, and every following token would deploy a module of its
// own (B6). The address is returned along with that error, so the caller can
// tell which contract was left unused.
//
// Initialize is called because the 4.1.6 modules are UUPS implementations with
// no locked constructor: until it runs, anyone can claim ownership (B8).
func deployComplianceModule(
	ctx context.Context,
	name, txLabel string,
	deploy func(*bind.TransactOpts) (common.Address, *types.Transaction, initializableModule, error),
) (moduleAddr common.Address, err error) {

	logger.LogInfo("💌 Deploying %s...", name)
	var module initializableModule
	tx, err := utils.Transact(ctx, func(auth *bind.TransactOpts) (tx *types.Transaction, err error) {
		moduleAddr, tx, module, err = deploy(auth)
		return
	})
	if err != nil {
		return common.Address{}, fmt.Errorf("deploy %s: %w", name, err)
	}
	deployed, err := utils.WaitDeployedTransaction(ctx, tx, true)
	if err != nil {
		return common.Address{}, fmt.Errorf("deploy %s: %w", name, err)
	}
	logger.LogInfo("📬 %s deployed at address: %s", name, moduleAddr.Hex())

	if err = database.InsertEthTransaction(ctx, deployed.Tx.Hash().Hex(), "DEPLOY_"+txLabel, deployed.ToAddressHex(), deployed.BlockNumber.Int64(), big.Int{}); err != nil {
		return moduleAddr, err
	}

	initTx, err := utils.Transact(ctx, func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return module.Initialize(auth)
	})
	if err != nil {
		return moduleAddr, fmt.Errorf("initialize %s: %w", name, err)
	}
	initialized, err := utils.WaitDeployedTransaction(ctx, initTx, false)
	if err != nil {
		return moduleAddr, fmt.Errorf("initialize %s: %w", name, err)
	}
	if err = database.InsertEthTransaction(ctx, initialized.Tx.Hash().Hex(), "INITIALIZE_"+txLabel, initialized.ToAddressHex(), initialized.BlockNumber.Int64(), big.Int{}); err != nil {
		return moduleAddr, err
	}

	if _, err = database.InsertContractRole(ctx, deployed.Tx.Hash().Hex(), moduleAddr.Hex(), name); err != nil {
		return moduleAddr, err
	}
	logger.LogInfo("✅ %s registered", name)

	return moduleAddr, nil
}

// ensureSharedComplianceModule returns the shared module registered under name,
// deploying and registering it first if this environment has never had one.
//
// The two capping modules arrived after environments were installed, and
// nothing reruns DeployTrexSuite on an installed environment: the first token
// created since is the one that brings them in.
func ensureSharedComplianceModule(ctx context.Context, name string) (common.Address, error) {
	details, err := database.GetContractRoleByName(ctx, name)
	if err == nil {
		return common.HexToAddress(details.Address), nil
	}
	// A database that cannot answer is not a database saying "none": deploying
	// on that answer would leave one more module behind at every incident.
	if !errors.Is(err, database.ErrContractRoleNotFound) {
		return common.Address{}, err
	}

	logger.LogInfo("🚀 No %s registered yet, deploying the shared one", name)
	moduleAddr, err := DeploySharedComplianceModules(ctx, name)
	if err == nil {
		return moduleAddr, nil
	}
	if !database.IsContractRoleNameTaken(err) {
		return common.Address{}, err
	}

	// Another token creation deployed it at the same time and registered first
	// (B7). Its module is the shared one; this one stays on-chain, unused.
	details, err = database.GetContractRoleByName(ctx, name)
	if err != nil {
		return common.Address{}, err
	}
	logger.LogWarn("⚠️ %s was registered concurrently at %s; %s is left unused", name, details.Address, moduleAddr.Hex())

	return common.HexToAddress(details.Address), nil
}
