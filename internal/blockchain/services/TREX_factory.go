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
)

// DeployAndInitTrexFactory deploys and configures the TREX Factory contract
func (s *Service) DeployAndInitTrexFactory(ctx context.Context) (err error) {

	allContractsImplementations, err := database.GetAllContractImplementations(ctx)
	if err != nil {
		return
	}

	authorityAddrPtr := utils.FindContractByName(allContractsImplementations, globals.ImplTrexAuthorityName)
	if authorityAddrPtr == nil {
		return errors.New("TREX Authority implementation address not found")
	}
	authorityAddr := *authorityAddrPtr

	identityFactoryAddrPtr, err := database.GetContractByName(ctx, globals.IdentityFactoryName)
	if err != nil {
		return
	}
	identityFactoryAddr := common.HexToAddress(identityFactoryAddrPtr.Address)

	// 1 - Deploy TREX Factory
	factoryAddr, err := deployTrexFactory(ctx, authorityAddr, identityFactoryAddr)
	if err != nil {
		return
	}

	// 2 - Configure TREX Factory
	/*onBehalfTransactions*/
	_, err = configureTrexFactory(ctx, authorityAddr, factoryAddr, identityFactoryAddr)
	if err != nil {
		return
	}

	// save Db factory address + onbehalf txs

	return
}

func deployTrexFactory(ctx context.Context, authorityAddr, identityFactoryAddr common.Address) (factoryAddr common.Address, err error) {

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	logger.LogInfo("💌 Deploying Trex Factory...")
	factoryAddr, tx, _, err := contracts.DeployTREXFactory(auth, globals.EthClient, authorityAddr, identityFactoryAddr)
	if err != nil {
		return
	}
	if err = utils.WaitDeployedTransaction(ctx, tx, true); err != nil {
		return
	}

	logger.LogInfo("📬 Trex Factory deployed at address: %s", factoryAddr.Hex())
	return
}

func configureTrexFactory(ctx context.Context, authorityAddr, factoryAddr, identityFactoryAddr common.Address) (txDetails []server.TxHashName, err error) {

	// 1 - Add TREX Factory in Identity Factory
	addTokenFactoryTx, err := addTokenFactory(ctx, identityFactoryAddr, factoryAddr)
	if err != nil {
		return
	}
	txDetails = append(txDetails, addTokenFactoryTx)

	// 2 - Set TREX Factory as active in Authority
	setTREXFactoryTx, err := setTREXFactory(ctx, authorityAddr, factoryAddr)
	if err != nil {
		return
	}
	txDetails = append(txDetails, setTREXFactoryTx)
	return
}

func addTokenFactory(ctx context.Context, identityFactoryAddr, factoryAddr common.Address) (txDetails server.TxHashName, err error) {

	identityFactoryInstance, err := contracts.NewIdFactory(identityFactoryAddr, globals.EthClient)
	if err != nil {
		return
	}

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}
	logger.LogInfo("💌 Send transaction to add token factory...")
	tx, err := identityFactoryInstance.AddTokenFactory(auth, factoryAddr)
	if err != nil {
		return
	}
	if err = utils.WaitDeployedTransaction(ctx, tx, false); err != nil {
		return
	}
	logger.LogInfo("📬 Token factory added to identity factory")

	txDetails = server.TxHashName{
		TransactionHash: tx.Hash().Hex(),
		OperationName:   "AddTokenFactory",
	}

	return
}

func setTREXFactory(ctx context.Context, authorityAddr, factoryAddr common.Address) (txDetails server.TxHashName, err error) {

	authorityInstance, err := contracts.NewTREXImplementationAuthorityTransactor(authorityAddr, globals.EthClient)
	if err != nil {
		return
	}

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}
	logger.LogInfo("💌 Send transaction to set TREX factory in Authority...")
	tx, err := authorityInstance.SetTREXFactory(auth, factoryAddr)
	if err != nil {
		return
	}
	if err = utils.WaitDeployedTransaction(ctx, tx, false); err != nil {
		return
	}
	logger.LogInfo("📬 TREX factory set")

	txDetails = server.TxHashName{
		TransactionHash: tx.Hash().Hex(),
		OperationName:   "SetTREXFactory",
	}

	return
}

// DeployTrexSuite deploys a TREX Suite (Claim Issuer, Compliance Modules, Token)
func (s *Service) DeployTrexSuite(ctx context.Context) (deployedContracts []server.ContractDetails, err error) {

	// 1 - Deploy Compliance Modules
	moduleAddr, err := DeploySharedComplianceModules(ctx, globals.TransferRestrictionModuleName)
	if err != nil {
		return
	}
	deployedContracts = append(deployedContracts, server.ContractDetails{
		Address: moduleAddr.Hex(),
		Name:    globals.TransferRestrictionModuleName,
	})

	// 2 - Deploy claim issuer
	claimIssuerAddr, claimIssuerInstance, err := DeployClaimIssuer(ctx)
	if err != nil {
		return
	}
	deployedContracts = append(deployedContracts, server.ContractDetails{
		Address: claimIssuerAddr.Hex(),
		Name:    globals.ClaimIssuerName,
	})

	// 3 - Add management key to claim issuer
	err = AddManagementKeyToClaimIssuer(ctx, claimIssuerInstance)
	if err != nil {
		return
	}

	// 4 - Define tokenDetails
	tokenDetails := defineTokenSuiteDetails([]common.Address{moduleAddr})

	// 5 - Define claimDetails
	claimDetails := defineClaimSuiteDetails(claimIssuerAddr)

	// 6 - Deploy TREX Suite
	trexSuiteAddr, err := deployTrexSuite(ctx, tokenDetails, claimDetails)
	if err != nil {
		return
	}
	deployedContracts = append(deployedContracts, server.ContractDetails{
		Address: trexSuiteAddr.Hex(),
		Name:    globals.TrexSuiteName,
	})

	return
}

func defineTokenSuiteDetails(moduleAddrs []common.Address) contracts.ITREXFactoryTokenDetails {

	ethFrom := utils.GetEthFrom()

	return contracts.ITREXFactoryTokenDetails{
		Owner:              ethFrom,
		Name:               "TREXSuiteV1",
		Symbol:             "TOOK",
		Decimals:           18,
		Irs:                common.Address{},
		ONCHAINID:          ethFrom,
		IrAgents:           []common.Address{ethFrom},
		TokenAgents:        []common.Address{ethFrom},
		ComplianceModules:  moduleAddrs,
		ComplianceSettings: [][]byte{},
	}
}

func defineClaimSuiteDetails(claimIssuerAddr common.Address) contracts.ITREXFactoryClaimDetails {
	return contracts.ITREXFactoryClaimDetails{
		ClaimTopics:  []*big.Int{big.NewInt(7)},
		Issuers:      []common.Address{claimIssuerAddr},
		IssuerClaims: [][]*big.Int{{big.NewInt(7)}},
	}
}

func deployTrexSuite(ctx context.Context, tokenDetails contracts.ITREXFactoryTokenDetails, claimDetails contracts.ITREXFactoryClaimDetails) (deploySuiteAddr common.Address, err error) {

	auth, err := utils.GenerateTransactOpts(ctx)
	if err != nil {
		return
	}

	trexFactoryDetails, err := database.GetContractByName(ctx, globals.TrexFactoryName)
	if err != nil {
		return
	}

	trexFactoryAddr := common.HexToAddress(trexFactoryDetails.Address)

	trexFactoryInstance, err := contracts.NewTREXFactory(trexFactoryAddr, globals.EthClient)
	if err != nil {
		return
	}

	salt := "TookenSuiteV1Salt"

	logger.LogInfo("💌 Deploying TREX Suite...")
	txDeploySuite, err := trexFactoryInstance.DeployTREXSuite(
		auth,
		salt,
		tokenDetails,
		claimDetails,
	)
	if err != nil {
		return
	}
	deploymentDetails, err := utils.WaitTREXSuiteDeployment(ctx, trexFactoryInstance, salt)
	if err != nil {
		return
	}

	deploySuiteAddr = deploymentDetails.Raw.Address
	logger.LogInfo("📬 TREX Suite deployment transaction mined: %s", txDeploySuite.Hash().Hex())

	// TODO: save deploymentDetails in DB

	return
}
