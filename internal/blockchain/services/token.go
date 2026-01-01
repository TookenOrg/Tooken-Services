package services

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	"github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-openapi/swag"
)

func (s *Service) CreateToken(ctx context.Context, req server.CreateTokenRequest) (newToken server.TokenInfos, err error) {

	// 1 - Idempotency
	if !isIdempotentToken(ctx, req) {
		err = errors.New("Token already exists")
		return
	}

	var onbehalfTransactions []server.TxHashName

	// 2 - Get TrexAuthority implementation
	trexAuthorityImplementationAddr, err := retreiveTrexAuthorityAddress(ctx)
	if err != nil {
		return
	}
	// 3 - Deploy Compliance Suite
	complianceSuite, err := deployComplianceSuite(ctx)
	if err != nil {
		return
	}

	// 4 - Deploy Token Proxy and retreive the generated Token.sol
	tokenAddr, txDeploytokenProxy, _, err := deployTokenProxy(ctx, trexAuthorityImplementationAddr, complianceSuite, req.TokenName, req.TokenName, req.NbDecimal)
	if err != nil {
		return
	}
	onbehalfTransactions = append(onbehalfTransactions, server.TxHashName{OperationName: "deployTokenProxy", TransactionHash: txDeploytokenProxy.Hash().Hex()})

	// 5 - Bind Token in compliance
	txBindToken, err := bindTokenToCompliance(ctx, complianceSuite.Instance, tokenAddr)
	if err != nil {
		return
	}
	onbehalfTransactions = append(onbehalfTransactions, server.TxHashName{OperationName: "bindTokenToCompliance", TransactionHash: txBindToken.Hash().Hex()})

	// 6 - Create Instance of Token with Token Address (not the proxy)
	tokenInstance, err := createTokenInstance(tokenAddr)
	if err != nil {
		return
	}

	// 7 - Add Agent on Token
	txAddAgentToken, err := addAgentOnToken(ctx, tokenInstance)
	if err != nil {
		return
	}
	onbehalfTransactions = append(onbehalfTransactions, server.TxHashName{OperationName: "addAgentOnToken", TransactionHash: txAddAgentToken.Hash().Hex()})

	// 8 - Unpause Token
	txUnpause, err := unpauseToken(ctx, tokenInstance)
	if err != nil {
		return
	}
	onbehalfTransactions = append(onbehalfTransactions, server.TxHashName{OperationName: "unpauseToken", TransactionHash: txUnpause.Hash().Hex()})

	// 9 - Add Agent on Identity Registry
	txIrAddAgent, err := addAgentOnIdentityRegistry(ctx, tokenAddr)
	if err != nil {
		return
	}
	onbehalfTransactions = append(onbehalfTransactions, server.TxHashName{OperationName: "AddAgentToIdentityRegister", TransactionHash: txIrAddAgent.Hash().Hex()})

	// 10 - Build TokenInfo for return
	newToken = server.TokenInfos{
		Symbol:                req.TokenName,
		Address:               tokenAddr.Hex(),
		NbDecimal:             int64(req.NbDecimal),
		OnbehalfTransactions:  &onbehalfTransactions,
		CreatedAt:             time.Now(),
		ModularComplianceAddr: swag.String(complianceSuite.Address.Hex()),
	}

	// 11 - DB

	return
}

func isIdempotentToken(ctx context.Context, tokenReq server.CreateTokenRequest) bool {
	existingTokenPtr, err := database.GetTokenByName(ctx, tokenReq.TokenName, true)
	if err != nil {
		logger.LogError("%s", err.Error())
		return false
	}
	return existingTokenPtr == nil
}

func createTokenInstance(tokenAddr common.Address) (instance contracts.Token, err error) {
	instancePtr, err := contracts.NewToken(tokenAddr, globals.EthClient)
	if instancePtr == nil {
		return contracts.Token{}, fmt.Errorf("Error when creating newToken instance: %+v", err)
	}
	return *instancePtr, nil
}

func addAgentOnToken(ctx context.Context, tokenInstance contracts.Token) (tx *types.Transaction, err error) {
	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Add agent on token...")
	tx, err = tokenInstance.AddAgent(auth, auth.From)
	if err != nil {
		return
	}

	deployedTxDetails, err := utils.WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return
	}
	logger.LogInfo("📬 Agent added on transaction: %s", tx.Hash().Hex())

	_, err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "ADD_AGENT_TOKEN", deployedTxDetails.Tx.To().Hex(), deployedTxDetails.BlockNumber.Int64(), big.Int{})
	if err != nil {
		return
	}

	return
}

func unpauseToken(ctx context.Context, tokenInstance contracts.Token) (tx *types.Transaction, err error) {
	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Unpause token...")
	tx, err = tokenInstance.Unpause(auth)
	if err != nil {
		return
	}

	deployedTxDetails, err := utils.WaitDeployedTransaction(ctx, tx, false)
	if err != nil {
		return
	}
	logger.LogInfo("📬 Token unpaused on transaction: %s", tx.Hash().Hex())

	_, err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "UNPAUSE_TOKEN", deployedTxDetails.Tx.To().Hex(), deployedTxDetails.BlockNumber.Int64(), big.Int{})
	if err != nil {
		return
	}
	return
}
