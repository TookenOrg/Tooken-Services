package services

import (
	"context"

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
		logger.LogWarn("User already has an identity")
	}

	// 2 - Generate wallet
	wallet, err := generateNewWallet(identityReq.UserId)
	if err != nil {
		return
	}

	// 3 - Deploy Identity Proxy
	proxyAddr, txProxy, err := deployIdentityProxy(ctx)
	if err != nil {
		return
	}

	// 4 - Register identity + wallet into Identity registry
	registerIdentityTxHashPtr, err := registerIdentity(ctx, wallet, proxyAddr, identityReq.CountryCode)
	if err != nil {
		return
	}

	// 5 - TODO Save DB
	logger.LogInfo("New Identity Created for userId [%d]", identityReq.UserId)

	// log to remove after insert db
	logger.LogDebug("%s, %s", txProxy.Hash().Hex(), registerIdentityTxHashPtr.Hash().Hex())

	newIdentityResponse.CountryCode = identityReq.CountryCode
	newIdentityResponse.UserId = identityReq.UserId
	newIdentityResponse.WalletAddress = wallet.Hex()
	newIdentityResponse.IdentityAddress = proxyAddr.Hex()
	newIdentityResponse.TransactionHash = txProxy.Hash().Hex()

	return
}

func isIdempotentIdentity(ctx context.Context, identityReq server.CreateIdentityRequest) bool {
	existingWalletPtr, err := database.GetWalletByUserId(ctx, identityReq.UserId)
	if err != nil {
		return false
	}
	return existingWalletPtr != nil
}

func generateNewWallet(userId int) (publicKey common.Address, err error) {
	logger.LogDebug("Generating new wallet address for user [%d]", userId)

	privateKeyecdsa, err := crypto.GenerateKey()
	if err != nil {
		return
	}
	// privateKey = fmt.Sprintf("%x", crypto.FromECDSA(privateKeyecdsa))
	publicKey = crypto.PubkeyToAddress(privateKeyecdsa.PublicKey)

	logger.LogDebug("New wallet created: public key [%s]", publicKey)

	// TODO
	// Save pubKey + private Key in DB

	return
}

func deployIdentityProxy(ctx context.Context) (proxyAddr common.Address, tx *types.Transaction, err error) {

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Deploying identity Proxy...")
	proxyAddr, tx, _, err = contracts.DeployIdentityProxy(auth, globals.EthClient, globals.ImplIdentityAuthorityAddress, (auth.From))
	if err != nil {
		return
	}

	if err = utils.WaitDeployedTransaction(ctx, tx, true); err != nil {
		return
	}

	logger.LogInfo("📬 Identity Proxy deployed at address: %s", proxyAddr.Hex())
	return
}

func getIdentityInstance(identityAddr common.Address) (instance *contracts.Identity, err error) {
	return contracts.NewIdentity(identityAddr, globals.EthClient)
}
