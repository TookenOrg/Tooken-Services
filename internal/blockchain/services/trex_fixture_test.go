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
const (
	kycTopic    = int64(7)    // claim topic required to hold the token (KYC)
	countryCode = uint16(250) // ISO-3166 numeric (France) for the investor
)

// keyHash returns the ERC-734 key id for an address: keccak256(abi.encode(addr)).
func keyHash(addr common.Address) [32]byte {
	return [32]byte(crypto.Keccak256Hash(common.LeftPadBytes(addr.Bytes(), 32)))
}

func failOnErr(t *testing.T, err error, what string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", what, err)
	}
}

type trexFixture struct {
	Ctx                  context.Context
	Client               *ethclient.Client
	Auth                 *bind.TransactOpts
	Call                 *bind.CallOpts
	Deployer             common.Address
	Factory              *contracts.TREXFactory
	FactoryAddr          common.Address
	ClaimIssuer          common.Address
	IdentityAuthorityRef common.Address // ONCHAINID ImplementationAuthority, what deployIdentityProxy needs
	Mine                 func(tx *types.Transaction, what string)
	suiteCounter         int
	investorCount        int
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
		Ctx:      ctx,
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
	f.IdentityAuthorityRef = idAuthority

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
// investor whitelist. A salt can only be spent once per factory.
//
// It goes through buildTokenDetails and defineClaimSuiteDetails — the very functions
// CreateToken uses — so a change in the way production assembles a suite is reflected
// here instead of quietly diverging from it.
func (f *trexFixture) DeploySuite(t *testing.T, salt, name, symbol string, irs common.Address) *deployedSuite {
	t.Helper()
	f.suiteCounter++

	tokenDetails := buildTokenDetails(f.Deployer, name, symbol, 18, irs, []common.Address{})
	claimDetails := defineClaimSuiteDetails(f.ClaimIssuer)

	tx, err := f.Factory.DeployTREXSuite(f.Auth, salt, tokenDetails, claimDetails)
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

	// A token is paused on creation. Minting is NOT gated by the pause — verified by
	// mutation: removing this call leaves every test green — but transfers are, so a
	// paused token is not a realistic starting state. Unpausing here keeps the fixture
	// honest for whatever a future test does, and changes no permissioning.
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

// NewStandaloneIRS deploys an IdentityRegistryStorage outside any suite: nothing is
// bound to it, so it has no registry to write through. It exists to exercise the case
// the factory can never produce — a whitelist with no door.
func (f *trexFixture) NewStandaloneIRS(t *testing.T) common.Address {
	t.Helper()

	addr, tx, irs, err := contracts.DeployIdentityRegistryStorage(f.Auth, f.Client)
	failOnErr(t, err, "deploy standalone IRS")
	f.Mine(tx, "standalone IRS")

	tx, err = irs.Init(f.Auth)
	failOnErr(t, err, "standalone IRS.init")
	f.Mine(tx, "standalone IRS.init")

	return addr
}

// MakePlatformIRSAgent used to run a one-off ownership round-trip so the platform
// could write straight into the shared IdentityRegistryStorage. It is gone: T-REX
// already wires the authorisation chain, since bindIdentityRegistry makes every
// IdentityRegistry an agent of the storage and buildTokenDetails makes the platform
// an agent of every IdentityRegistry. Registration needs no setup at all.
