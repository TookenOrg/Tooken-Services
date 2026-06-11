package services

import (
	"context"
	"math/big"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
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

	deployedTxDetails, err := utils.WaitDeployedTransaction(ctx, tx, true)
	if err != nil {
		return
	}
	logger.LogInfo("📬 Claim Issuer deployed at address: %s", claimIssuerAddr.Hex())

	err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "CLAIM_ISSUER", deployedTxDetails.Tx.To().Hex(), deployedTxDetails.BlockNumber.Int64(), big.Int{})
	if err != nil {
		return
	}

	return
}

func AddClaimSignerKeyToClaimIssuer(ctx context.Context, claimIssuerInstance *contracts.ClaimIssuer) (err error) {

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	// ONCHAINID keys are keccak256(abi.encode(address)) == keccak256(left-pad-32(address)).
	key := crypto.Keccak256Hash(common.LeftPadBytes(auth.From.Bytes(), 32))

	logger.LogInfo("💌 Send transaction to add claim-signer key to Claim Issuer...")
	tx, err := claimIssuerInstance.AddKey(
		auth,
		key, big.NewInt(3),
		big.NewInt(1),
	)
	if err != nil {
		return
	}

	deployedTxDetails, err := utils.WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return
	}
	logger.LogInfo("📬 Claim-signer key added to Claim Issuer")

	err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "ADD_MANAGEMENT_KEY", deployedTxDetails.Tx.To().Hex(), deployedTxDetails.BlockNumber.Int64(), big.Int{})
	if err != nil {
		return
	}

	return
}
