package services

import (
	"context"
	"errors"
	"math/big"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func (s *Service) CreateIdentity(ctx context.Context, identityReq server.CreateIdentityRequest) (newIdentityResponse server.IdentityInfos, err error) {

	// 1 - Idempotency
	if !isIdempotentIdentity(ctx, identityReq) {
		return server.IdentityInfos{}, errors.New("User already has an identity")
	}

	// 2 - Generate wallet
	walletPubKey, walletId, err := GenerateNewWallet(identityReq.UserId)
	if err != nil {
		return
	}

	// 3 - Deploy Identity Proxy
	proxyAddr, txProxy, err := deployIdentityProxy(ctx)
	if err != nil {
		return
	}

	// 4 - Register identity + wallet into Identity registry
	registerIdentityTxHashPtr, err := registerIdentity(ctx, walletPubKey, proxyAddr, identityReq.CountryCode)
	if err != nil {
		return
	}

	err = database.InsertIdentity(ctx, identityReq.UserId, walletId, proxyAddr.Hex(), txProxy.Hash().Hex())
	if err != nil {
		return
	}

	// log to remove after insert db
	logger.LogDebug("%s, %s", txProxy.Hash().Hex(), registerIdentityTxHashPtr.Hash().Hex())

	newIdentityResponse.CountryCode = identityReq.CountryCode
	newIdentityResponse.UserId = identityReq.UserId
	newIdentityResponse.WalletAddress = walletPubKey.Hex()
	newIdentityResponse.IdentityAddress = proxyAddr.Hex()
	newIdentityResponse.TransactionHash = txProxy.Hash().Hex()

	logger.LogInfo("New Identity Created for userId [%d]", identityReq.UserId)

	return
}

func isIdempotentIdentity(ctx context.Context, identityReq server.CreateIdentityRequest) bool {
	existingWalletPtr, err := database.GetWalletByUserId(ctx, identityReq.UserId)
	if err != nil {
		return false
	}
	return existingWalletPtr == nil
}

func GenerateNewWallet(userId int) (publicKey common.Address, walletId int64, err error) {
	logger.LogDebug("Generating new wallet address for user [%d]", userId)

	privateKeyecdsa, err := crypto.GenerateKey()
	if err != nil {
		return
	}
	publicKey = crypto.PubkeyToAddress(privateKeyecdsa.PublicKey)

	logger.LogDebug("New wallet created: public key [%s]", publicKey)

	walletId, err = database.InsertWallet(userId, publicKey.Hex(), "Main Wallet")

	return
}

func deployIdentityProxy(ctx context.Context) (proxyAddr common.Address, tx *types.Transaction, err error) {

	// Get Impl Identity Authority Address
	implIdentityAuthorityDetails, err := database.GetImplementationContractByName(ctx, globals.ImplIdentityAuthorityName)
	if err != nil {
		return
	}

	identityAuthorityAddress := common.HexToAddress(implIdentityAuthorityDetails.Address)

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Deploying identity Proxy...")
	proxyAddr, tx, _, err = contracts.DeployIdentityProxy(auth, globals.EthClient, identityAuthorityAddress, (auth.From))
	if err != nil {
		return
	}

	deployedTxDetails, err := utils.WaitDeployedTransaction(ctx, tx, true)
	if err != nil {
		return
	}

	logger.LogInfo("📬 Identity Proxy deployed at address: %s", proxyAddr.Hex())

	err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "IDENTITY_PROXY", deployedTxDetails.ToAddressHex(), deployedTxDetails.BlockNumber.Int64(), big.Int{})
	if err != nil {
		return
	}

	return
}

func getIdentityInstance(identityAddr common.Address) (instance *contracts.Identity, err error) {
	return contracts.NewIdentity(identityAddr, globals.EthClient)
}
