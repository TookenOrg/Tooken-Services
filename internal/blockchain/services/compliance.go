package services

import (
	"context"
	"errors"
	"math/big"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

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
	deployedTxDetails, err := utils.WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return
	}
	logger.LogInfo("📬 Transfer Restriction Module deployed at address: %s", moduleAddr.Hex())

	err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "DEPLOY_TRANSFER_RESTRICT_MODULE", deployedTxDetails.Tx.To().Hex(), deployedTxDetails.BlockNumber.Int64(), big.Int{})
	if err != nil {
		return
	}

	if _, perr := database.InsertContractRole(ctx, deployedTxDetails.Tx.Hash().Hex(), moduleAddr.Hex(), globals.TransferRestrictionModuleName); perr != nil {
		logger.LogWarn("could not persist TRANSFER_RESTRICTION_MODULE role: %s", perr.Error())
	}
	return
}
