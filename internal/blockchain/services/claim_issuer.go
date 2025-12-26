package services

import (
	"context"
	"math/big"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func DeployClaimIssuer(ctx context.Context) (claimIssuerAddr common.Address, claimIssuerInstance *contracts.ClaimIssuer, err error) {

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Deploying Claim Issuer...")
	claimIssuerAddr, tx, claimIssuerInstance, err := contracts.DeployClaimIssuer(auth, globals.EthClient, auth.From)
	if err != nil {
		return
	}

	err = utils.WaitDeployedTransaction(ctx, tx, true)
	if err != nil {
		return
	}
	logger.LogInfo("📬 Claim Issuer deployed at address: %s", claimIssuerAddr.Hex())

	// TODO : save in dB

	return
}

func AddManagementKeyToClaimIssuer(ctx context.Context, claimIssuerInstance *contracts.ClaimIssuer) (err error) {

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	key := crypto.Keccak256Hash(auth.From.Bytes())

	logger.LogInfo("💌 Send transaction to add management key to Claim Issuer...")
	tx, err := claimIssuerInstance.AddKey(
		auth,
		key, big.NewInt(3),
		big.NewInt(1),
	)
	if err != nil {
		return
	}

	err = utils.WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return
	}
	logger.LogInfo("📬 Management key added to Claim Issuer")

	// TODO : save in dB

	return
}
