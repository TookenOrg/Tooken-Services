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

	// 0 - Validate input (T-REX Token.init enforces decimals <= 18 on-chain)
	if req.NbDecimal < 0 || req.NbDecimal > 18 {
		err = errors.New("nbDecimal must be between 0 and 18")
		return
	}

	// blk.token: token_name is varchar(100), symbol is varchar(50).
	if n := len(req.TokenName); n == 0 || n > 100 {
		err = errors.New("token name must be between 1 and 100 characters")
		return
	}
	if n := len(req.Symbol); n == 0 || n > 50 {
		err = errors.New("token symbol must be between 1 and 50 characters")
		return
	}

	// 1 - Idempotency
	if !isIdempotentToken(ctx, req) {
		err = errors.New("Token already exists")
		return
	}

	// 2 - Resolve the shared infrastructure (one factory, one claim issuer, one module)
	factoryDetails, err := database.GetContractRoleByName(ctx, globals.TrexFactoryName)
	if err != nil {
		return
	}
	factoryAddr := common.HexToAddress(factoryDetails.Address)

	claimIssuerDetails, err := database.GetContractRoleByName(ctx, globals.ClaimIssuerName)
	if err != nil {
		return
	}
	claimIssuerAddr := common.HexToAddress(claimIssuerDetails.Address)

	moduleAddr, found, err := database.FindModuleByName(ctx, globals.TransferRestrictionModuleName)
	if err != nil {
		return
	}
	if !found {
		err = fmt.Errorf("shared compliance module %s not found", globals.TransferRestrictionModuleName)
		return
	}

	// 3 - Resolve the shared IRS (shared investor whitelist). Zero address => first
	//     token: the factory deploys a fresh IRS that we persist for the next tokens.
	sharedIRS := resolveSharedIRS(ctx)

	// 4 - Build the suite parameters (owner & agents = the single platform manager)
	ethFrom := utils.GetEthFrom()
	tokenDetails := buildTokenDetails(ethFrom, req.TokenName, req.Symbol, req.NbDecimal, sharedIRS, []common.Address{moduleAddr})
	claimDetails := defineClaimSuiteDetails(claimIssuerAddr)

	// 5 - Deploy the token suite through the shared TREX factory
	deployment, suiteTxHash, err := deployTokenSuiteViaFactory(ctx, factoryAddr, req.TokenName, tokenDetails, claimDetails)
	if err != nil {
		return
	}
	tokenAddr := deployment.Token
	mcAddr := deployment.Mc

	onbehalfTransactions := []server.TxHashName{
		{OperationName: "DeployTREXSuite", TransactionHash: suiteTxHash},
	}

	// 6 - Persist the shared IRS the first time it is created
	if sharedIRS == (common.Address{}) {
		if _, perr := database.InsertContractRole(ctx, suiteTxHash, deployment.Irs.Hex(), globals.SharedIdentityRegistryStorageName); perr != nil {
			logger.LogWarn("could not persist shared IRS role: %s", perr.Error())
		}
	}

	// 7 - Unpause the token so transfers/mints become possible (ethFrom is a token agent)
	tokenInstance, err := createTokenInstance(tokenAddr)
	if err != nil {
		return
	}
	txUnpause, err := unpauseToken(ctx, tokenInstance)
	if err != nil {
		return
	}
	onbehalfTransactions = append(onbehalfTransactions, server.TxHashName{OperationName: "unpauseToken", TransactionHash: txUnpause.Hash().Hex()})

	// 8 - Persist the token
	if _, perr := database.InsertToken(ctx, req.Symbol, req.TokenName, tokenAddr.Hex(), req.NbDecimal, mcAddr.Hex()); perr != nil {
		logger.LogWarn("could not persist token: %s", perr.Error())
	}

	// 9 - Build response
	newToken = server.TokenInfos{
		Symbol:                req.Symbol,
		TokenName:             req.TokenName,
		Address:               tokenAddr.Hex(),
		NbDecimal:             int64(req.NbDecimal),
		OnbehalfTransactions:  &onbehalfTransactions,
		CreatedAt:             time.Now(),
		ModularComplianceAddr: swag.String(mcAddr.Hex()),
	}

	return
}

// resolveSharedIRS returns the shared IdentityRegistryStorage address persisted as a
// contract role, or the zero address if none has been recorded yet.
func resolveSharedIRS(ctx context.Context) common.Address {
	details, err := database.GetContractRoleByName(ctx, globals.SharedIdentityRegistryStorageName)
	if err != nil {
		return common.Address{}
	}
	return common.HexToAddress(details.Address)
}

// buildTokenDetails assembles the ITREXFactoryTokenDetails for a single token.
// irs is the shared IRS (zero lets the factory deploy a fresh one); ONCHAINID is left
// zero so the factory creates the token ONCHAINID through the IdFactory.
func buildTokenDetails(owner common.Address, name, symbol string, decimals int, irs common.Address, modules []common.Address) contracts.ITREXFactoryTokenDetails {
	return contracts.ITREXFactoryTokenDetails{
		Owner:              owner,
		Name:               name,
		Symbol:             symbol,
		Decimals:           uint8(decimals),
		Irs:                irs,
		ONCHAINID:          common.Address{},
		IrAgents:           []common.Address{owner},
		TokenAgents:        []common.Address{owner},
		ComplianceModules:  modules,
		ComplianceSettings: [][]byte{},
	}
}

// deployTokenSuiteViaFactory deploys a token suite through the shared TREX factory and
// waits for the TREXSuiteDeployed event to retrieve the deployed addresses.
func deployTokenSuiteViaFactory(ctx context.Context, factoryAddr common.Address, salt string, tokenDetails contracts.ITREXFactoryTokenDetails, claimDetails contracts.ITREXFactoryClaimDetails) (deployment *contracts.TREXFactoryTREXSuiteDeployed, txHash string, err error) {

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	factoryInstance, err := contracts.NewTREXFactory(factoryAddr, globals.EthClient)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Deploying token suite via TREX factory (salt=%s)...", salt)
	tx, err := factoryInstance.DeployTREXSuite(auth, salt, tokenDetails, claimDetails)
	if err != nil {
		return
	}

	deployment, err = utils.WaitTREXSuiteDeployment(ctx, factoryInstance, salt)
	if err != nil {
		return
	}
	txHash = tx.Hash().Hex()

	logger.LogInfo("📬 Token suite deployed: token=%s ir=%s irs=%s mc=%s", deployment.Token.Hex(), deployment.Ir.Hex(), deployment.Irs.Hex(), deployment.Mc.Hex())

	if derr := database.InsertEthTransaction(ctx, txHash, "DEPLOY_TREX_SUITE", factoryAddr.Hex(), int64(deployment.Raw.BlockNumber), big.Int{}); derr != nil {
		logger.LogWarn("could not persist DEPLOY_TREX_SUITE tx: %s", derr.Error())
	}

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

	err = database.InsertEthTransaction(ctx, deployedTxDetails.Tx.Hash().Hex(), "UNPAUSE_TOKEN", deployedTxDetails.Tx.To().Hex(), deployedTxDetails.BlockNumber.Int64(), big.Int{})
	if err != nil {
		return
	}
	return
}
