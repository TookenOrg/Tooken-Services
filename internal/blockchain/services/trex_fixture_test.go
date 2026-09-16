//go:build integration

// Shared T-REX fixture for the integration tests of this package.
//
// Deploying the T-REX infrastructure takes about fifty lines of sequential
// transactions, and that sequence was already copied in two test files. Extracting
// it here keeps the tests about what they assert rather than about what they set up.
//
// Requires a local node:
//
//	cd tools/hardhat && npm install && npx hardhat node
//	go test -tags integration ./internal/blockchain/...
package services

import (
	"context"
	"math/big"
	"os"
	"testing"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// trexFixture holds a deployed T-REX infrastructure: the six implementations, the
// ONCHAINID factory, the TREX factory and the single shared claim issuer.
//
// The deployer is the platform manager: owner and agent of everything, and the
// signer behind PRIVATE_KEY, exactly like the production `ethFrom`.
type trexFixture struct {
	Client        *ethclient.Client
	Auth          *bind.TransactOpts
	Call          *bind.CallOpts
	Deployer      common.Address
	Factory       *contracts.TREXFactory
	FactoryAddr   common.Address
	ClaimIssuer   common.Address
	Mine          func(tx *types.Transaction, what string)
	suiteCounter  int
	investorCount int
}

// newTREXFixture dials the local node and deploys the whole T-REX infrastructure.
// It skips the test — rather than failing it — when no node is reachable, so the
// default `go test ./...` stays green without a chain.
func newTREXFixture(ctx context.Context, t *testing.T) *trexFixture {
	t.Helper()

	rpc := os.Getenv("ETH_TEST_RPC")
	if rpc == "" {
		rpc = "http://127.0.0.1:8545"
	}
	client, err := ethclient.DialContext(ctx, rpc)
	if err != nil {
		t.Skipf("no local EVM node at %s: %v", rpc, err)
	}
	chainID, err := client.ChainID(ctx)
	if err != nil {
		t.Skipf("local EVM node at %s not reachable: %v", rpc, err)
	}

	deployerKey, err := crypto.HexToECDSA(hardhatAccount0Key)
	if err != nil {
		t.Fatalf("parse deployer key: %v", err)
	}
	deployer := crypto.PubkeyToAddress(deployerKey.PublicKey)
	// generateSignatureAddClaim and GenerateTransactOpts both sign with PRIVATE_KEY.
	t.Setenv("PRIVATE_KEY", hardhatAccount0Key)

	// Leaving Nonce nil and GasLimit 0 lets go-ethereum manage the nonce and estimate
	// gas for each sequential transaction.
	auth, err := bind.NewKeyedTransactorWithChainID(deployerKey, chainID)
	if err != nil {
		t.Fatalf("transactor: %v", err)
	}

	f := &trexFixture{
		Client:   client,
		Auth:     auth,
		Call:     &bind.CallOpts{Context: ctx},
		Deployer: deployer,
	}
	f.Mine = func(tx *types.Transaction, what string) {
		t.Helper()
		rcpt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			t.Fatalf("%s: wait mined: %v", what, err)
		}
		if rcpt.Status != types.ReceiptStatusSuccessful {
			t.Fatalf("%s: transaction reverted (%s)", what, tx.Hash().Hex())
		}
	}

	// ---- the six T-REX logic implementations -------------------------------
	ctrImpl, tx, _, err := contracts.DeployClaimTopicsRegistry(auth, client)
	failOnErr(t, err, "deploy CTR impl")
	f.Mine(tx, "CTR impl")
	tirImpl, tx, _, err := contracts.DeployTrustedIssuersRegistry(auth, client)
	failOnErr(t, err, "deploy TIR impl")
	f.Mine(tx, "TIR impl")
	irsImpl, tx, _, err := contracts.DeployIdentityRegistryStorage(auth, client)
	failOnErr(t, err, "deploy IRS impl")
	f.Mine(tx, "IRS impl")
	irImpl, tx, _, err := contracts.DeployIdentityRegistry(auth, client)
	failOnErr(t, err, "deploy IR impl")
	f.Mine(tx, "IR impl")
	mcImpl, tx, _, err := contracts.DeployModularCompliance(auth, client)
	failOnErr(t, err, "deploy MC impl")
	f.Mine(tx, "MC impl")
	tokenImpl, tx, _, err := contracts.DeployToken(auth, client)
	failOnErr(t, err, "deploy Token impl")
	f.Mine(tx, "Token impl")

	// ---- ONCHAINID implementation, authority and factory -------------------
	identityImpl, tx, _, err := contracts.DeployIdentity(auth, client, deployer, true)
	failOnErr(t, err, "deploy Identity impl")
	f.Mine(tx, "Identity impl")
	idAuthority, tx, _, err := contracts.DeployImplementationAuthority(auth, client, identityImpl)
	failOnErr(t, err, "deploy identity ImplementationAuthority")
	f.Mine(tx, "identity authority")
	idFactoryAddr, tx, idFactory, err := contracts.DeployIdFactory(auth, client, idAuthority)
	failOnErr(t, err, "deploy IdFactory")
	f.Mine(tx, "IdFactory")

	// ---- TREX ImplementationAuthority + registered version -----------------
	trexAuthorityAddr, tx, trexAuthority, err := contracts.DeployTREXImplementationAuthority(auth, client, true, common.Address{}, common.Address{})
	failOnErr(t, err, "deploy TREX authority")
	f.Mine(tx, "TREX authority")
	tx, err = trexAuthority.AddAndUseTREXVersion(auth,
		contracts.ITREXImplementationAuthorityVersion{Major: 1, Minor: 0, Patch: 0},
		contracts.ITREXImplementationAuthorityTREXContracts{
			TokenImplementation: tokenImpl,
			CtrImplementation:   ctrImpl,
			IrImplementation:    irImpl,
			IrsImplementation:   irsImpl,
			TirImplementation:   tirImpl,
			McImplementation:    mcImpl,
		})
	failOnErr(t, err, "AddAndUseTREXVersion")
	f.Mine(tx, "AddAndUseTREXVersion")

	// ---- TREX factory, wired to the authority and the identity factory -----
	factoryAddr, tx, factory, err := contracts.DeployTREXFactory(auth, client, trexAuthorityAddr, idFactoryAddr)
	failOnErr(t, err, "deploy TREX factory")
	f.Mine(tx, "TREX factory")
	tx, err = idFactory.AddTokenFactory(auth, factoryAddr)
	failOnErr(t, err, "addTokenFactory")
	f.Mine(tx, "addTokenFactory")
	tx, err = trexAuthority.SetTREXFactory(auth, factoryAddr)
	failOnErr(t, err, "setTREXFactory")
	f.Mine(tx, "setTREXFactory")
	f.Factory, f.FactoryAddr = factory, factoryAddr

	// ---- the single shared claim issuer, trusted for the KYC topic ---------
	claimIssuerAddr, tx, claimIssuer, err := contracts.DeployClaimIssuer(auth, client, deployer)
	failOnErr(t, err, "deploy ClaimIssuer")
	f.Mine(tx, "ClaimIssuer")
	tx, err = claimIssuer.AddKey(auth, keyHash(deployer), big.NewInt(3), big.NewInt(1)) // purpose 3 = CLAIM
	failOnErr(t, err, "ClaimIssuer.AddKey")
	f.Mine(tx, "ClaimIssuer.AddKey")
	f.ClaimIssuer = claimIssuerAddr

	return f
}

// deployedSuite is what a single deployTREXSuite call produces.
type deployedSuite struct {
	Token     *contracts.Token
	TokenAddr common.Address
	IR        *contracts.IdentityRegistry
	IRAddr    common.Address
	IRS       *contracts.IdentityRegistryStorage
	IRSAddr   common.Address
}

// DeploySuite deploys one token suite through the factory. Pass the zero address as
// irs to let the factory create a fresh storage; pass an existing one to share the
// investor whitelist. Salts are unique per fixture: the factory reverts on reuse.
func (f *trexFixture) DeploySuite(t *testing.T, name, symbol string, irs common.Address) *deployedSuite {
	t.Helper()
	f.suiteCounter++
	salt := name

	tx, err := f.Factory.DeployTREXSuite(f.Auth, salt,
		contracts.ITREXFactoryTokenDetails{
			Owner:              f.Deployer,
			Name:               name,
			Symbol:             symbol,
			Decimals:           18,
			Irs:                irs,
			ONCHAINID:          common.Address{},
			IrAgents:           []common.Address{f.Deployer},
			TokenAgents:        []common.Address{f.Deployer},
			ComplianceModules:  []common.Address{},
			ComplianceSettings: [][]byte{},
		},
		contracts.ITREXFactoryClaimDetails{
			ClaimTopics:  []*big.Int{big.NewInt(kycTopic)},
			Issuers:      []common.Address{f.ClaimIssuer},
			IssuerClaims: [][]*big.Int{{big.NewInt(kycTopic)}},
		})
	failOnErr(t, err, "DeployTREXSuite "+salt)
	f.Mine(tx, "DeployTREXSuite "+salt)

	tokenAddr, err := f.Factory.GetToken(f.Call, salt)
	failOnErr(t, err, "GetToken "+salt)
	if tokenAddr == (common.Address{}) {
		t.Fatalf("token %s not registered by the factory", salt)
	}
	token, err := contracts.NewToken(tokenAddr, f.Client)
	failOnErr(t, err, "bind token "+salt)
	irAddr, err := token.IdentityRegistry(f.Call)
	failOnErr(t, err, "token.identityRegistry "+salt)
	ir, err := contracts.NewIdentityRegistry(irAddr, f.Client)
	failOnErr(t, err, "bind IR "+salt)
	irsAddr, err := ir.IdentityStorage(f.Call)
	failOnErr(t, err, "ir.identityStorage "+salt)
	irsInstance, err := contracts.NewIdentityRegistryStorage(irsAddr, f.Client)
	failOnErr(t, err, "bind IRS "+salt)

	// A token is paused on creation; nothing can be minted until it is unpaused.
	tx, err = token.Unpause(f.Auth)
	failOnErr(t, err, "Unpause "+salt)
	f.Mine(tx, "Unpause "+salt)

	return &deployedSuite{
		Token: token, TokenAddr: tokenAddr,
		IR: ir, IRAddr: irAddr,
		IRS: irsInstance, IRSAddr: irsAddr,
	}
}

// NewInvestorIdentity deploys an ONCHAINID for a fresh wallet and grants the platform
// a CLAIM key on it, which is what lets AddKYCClaim sign a claim for that identity.
// It deliberately does NOT register the wallet anywhere: that is what the tests assert.
func (f *trexFixture) NewInvestorIdentity(t *testing.T) (wallet, identityAddr common.Address, identity *contracts.Identity) {
	t.Helper()
	f.investorCount++

	key, err := crypto.GenerateKey()
	failOnErr(t, err, "generate investor key")
	wallet = crypto.PubkeyToAddress(key.PublicKey)

	identityAddr, tx, identity, err := contracts.DeployIdentity(f.Auth, f.Client, f.Deployer, false)
	failOnErr(t, err, "deploy investor Identity")
	f.Mine(tx, "investor Identity")

	tx, err = identity.AddKey(f.Auth, keyHash(f.Deployer), big.NewInt(3), big.NewInt(1)) // purpose 3 = CLAIM
	failOnErr(t, err, "investorID.AddKey")
	f.Mine(tx, "investorID.AddKey")

	return wallet, identityAddr, identity
}

// AddKYCClaim signs a KYC claim with the platform key and stores it on the investor
// ONCHAINID. The claim lives on the identity, not on any token, which is why a single
// KYC serves every token of the platform.
func (f *trexFixture) AddKYCClaim(t *testing.T, identityAddr common.Address, identity *contracts.Identity) {
	t.Helper()

	sig, err := generateSignatureAddClaim(identityAddr, kycTopic)
	failOnErr(t, err, "generateSignatureAddClaim")
	tx, err := identity.AddClaim(f.Auth, big.NewInt(kycTopic), big.NewInt(1), f.ClaimIssuer, sig.SignatureBytes, []byte(claimData), "")
	failOnErr(t, err, "AddClaim")
	f.Mine(tx, "AddClaim")
}

// MakePlatformIRSAgent runs the one-off setup that lets the platform write into the
// shared IdentityRegistryStorage: the factory owns it, so ownership is recovered, the
// agent added, and ownership handed straight back.
//
// The last step is not optional. bindIdentityRegistry is onlyOwner on the storage, so
// a factory that no longer owns it can no longer deploy a suite reusing the shared
// IRS — token creation would break to fix investor registration.
func (f *trexFixture) MakePlatformIRSAgent(t *testing.T, suite *deployedSuite) {
	t.Helper()

	tx, err := f.Factory.RecoverContractOwnership(f.Auth, suite.IRSAddr, f.Deployer)
	failOnErr(t, err, "RecoverContractOwnership(IRS)")
	f.Mine(tx, "RecoverContractOwnership")

	tx, err = suite.IRS.AddAgent(f.Auth, f.Deployer)
	failOnErr(t, err, "IRS.AddAgent")
	f.Mine(tx, "IRS.AddAgent")

	tx, err = suite.IRS.TransferOwnership(f.Auth, f.FactoryAddr)
	failOnErr(t, err, "IRS.TransferOwnership back to the factory")
	f.Mine(tx, "IRS.TransferOwnership")
}
