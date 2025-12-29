package services

import (
	"context"

	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func registerIdentity(ctx context.Context, userWallet, identityAddress common.Address, countryCode int) (tx *types.Transaction, err error) {

	irInstance, err := database.GetIdentityRegistryInstance()

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Registration of new identity in identity registry...")
	tx, err = irInstance.RegisterIdentity(auth, userWallet, identityAddress, uint16(countryCode))
	if err != nil {
		return
	}
	if err = utils.WaitDeployedTransaction(ctx, tx, false); err != nil {
		return
	}

	logger.LogInfo("📬 Identity registred on transaction: %s", tx.Hash().Hex())

	// TODO Save DB

	return
}

func addAgentOnIdentityRegistry(ctx context.Context, tokenAddr common.Address) (tx *types.Transaction, err error) {
	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Add agent on identity registry...")
	tx, err = globals.IdentityRegistryInstance.AddAgent(auth, tokenAddr)
	if err != nil {
		return
	}
	if err = utils.WaitDeployedTransaction(ctx, tx, false); err != nil {
		return
	}

	logger.LogInfo("📬 Agent added on identity registry: %s", tx.Hash().Hex())
	return
}
