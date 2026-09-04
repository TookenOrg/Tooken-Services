//go:build integration

package services

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TestSharedIRSReuseIntegration proves the core promise of the CreateToken refactor:
// a SECOND token deployed through the factory while REUSING the shared
// IdentityRegistryStorage (SHARED_IRS) shares one investor whitelist with the first
// token. An investor KYC'd once (via the IR the factory bound to the shared IRS) is
// verified and can receive shares on BOTH tokens.
//
// It exercises the real production helpers buildTokenDetails + defineClaimSuiteDetails
// (the same construction CreateToken uses) and confirms that the factory retains
// ownership of the shared IRS and binds EVERY token's IR to it: both IRs become agents
// of the shared storage, so any token can register investors into the one shared
// whitelist and all tokens verify against it.
//
// Gated by the `integration` build tag. Requires a local node:
//
//	cd tools/hardhat && npm install && npx hardhat node
//	go test -tags integration -run TestSharedIRSReuseIntegration -v ./internal/blockchain/...
func TestSharedIRSReuseIntegration(t *testing.T) {
	rpc := os.Getenv("ETH_TEST_RPC")
	if rpc == "" {
		rpc = "http://127.0.0.1:8545"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpc)
	if err != nil {
		t.Skipf("no local EVM node at %s: %v", rpc, err)
	}
	chainID, err := client.ChainID(ctx)
	if err != nil {
		t.Skipf("local EVM node at %s not reachable: %v", rpc, err)
	}

	const deployerKeyHex = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	deployerKey, err := crypto.HexToECDSA(deployerKeyHex)
	if err != nil {
		t.Fatalf("parse deployer key: %v", err)
	}
	deployer := crypto.PubkeyToAddress(deployerKey.PublicKey)
	// generateSignatureAddClaim signs with PRIVATE_KEY; use the same deployer key.
	t.Setenv("PRIVATE_KEY", deployerKeyHex)

	auth, err := bind.NewKeyedTransactorWithChainID(deployerKey, chainID)
	if err != nil {
		t.Fatalf("transactor: %v", err)
	}
	co := &bind.CallOpts{Context: ctx}

	mine := func(tx *types.Transaction, what string) {
		t.Helper()
		rcpt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			t.Fatalf("%s: wait mined: %v", what, err)
		}
		if rcpt.Status != types.ReceiptStatusSuccessful {
			t.Fatalf("%s: transaction reverted (%s)", what, tx.Hash().Hex())
		}
	}

	// ---- T-REX infrastructure (implementations, ONCHAINID, authority, factory) ----
	t.Log("deploying T-REX infrastructure")
	ctrImpl, tx, _, err := contracts.DeployClaimTopicsRegistry(auth, client)
	failOnErr(t, err, "deploy CTR impl")
	mine(tx, "CTR impl")
	tirImpl, tx, _, err := contracts.DeployTrustedIssuersRegistry(auth, client)
	failOnErr(t, err, "deploy TIR impl")
	mine(tx, "TIR impl")
	irsImpl, tx, _, err := contracts.DeployIdentityRegistryStorage(auth, client)
	failOnErr(t, err, "deploy IRS impl")
	mine(tx, "IRS impl")
	irImpl, tx, _, err := contracts.DeployIdentityRegistry(auth, client)
	failOnErr(t, err, "deploy IR impl")
	mine(tx, "IR impl")
	mcImpl, tx, _, err := contracts.DeployModularCompliance(auth, client)
	failOnErr(t, err, "deploy MC impl")
	mine(tx, "MC impl")
	tokenImpl, tx, _, err := contracts.DeployToken(auth, client)
	failOnErr(t, err, "deploy Token impl")
	mine(tx, "Token impl")

	identityImpl, tx, _, err := contracts.DeployIdentity(auth, client, deployer, true)
	failOnErr(t, err, "deploy Identity impl")
	mine(tx, "Identity impl")
	idAuthority, tx, _, err := contracts.DeployImplementationAuthority(auth, client, identityImpl)
	failOnErr(t, err, "deploy identity ImplementationAuthority")
	mine(tx, "identity authority")
	idFactoryAddr, tx, idFactory, err := contracts.DeployIdFactory(auth, client, idAuthority)
	failOnErr(t, err, "deploy IdFactory")
	mine(tx, "IdFactory")

	trexAuthorityAddr, tx, trexAuthority, err := contracts.DeployTREXImplementationAuthority(auth, client, true, common.Address{}, common.Address{})
	failOnErr(t, err, "deploy TREX authority")
	mine(tx, "TREX authority")
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
	mine(tx, "AddAndUseTREXVersion")

	trexFactoryAddr, tx, trexFactory, err := contracts.DeployTREXFactory(auth, client, trexAuthorityAddr, idFactoryAddr)
	failOnErr(t, err, "deploy TREX factory")
	mine(tx, "TREX factory")
	tx, err = idFactory.AddTokenFactory(auth, trexFactoryAddr)
	failOnErr(t, err, "addTokenFactory")
	mine(tx, "addTokenFactory")
	tx, err = trexAuthority.SetTREXFactory(auth, trexFactoryAddr)
	failOnErr(t, err, "setTREXFactory")
	mine(tx, "setTREXFactory")

	claimIssuerAddr, tx, claimIssuer, err := contracts.DeployClaimIssuer(auth, client, deployer)
	failOnErr(t, err, "deploy ClaimIssuer")
	mine(tx, "ClaimIssuer")
	tx, err = claimIssuer.AddKey(auth, keyHash(deployer), big.NewInt(3), big.NewInt(1)) // purpose 3 = CLAIM
	failOnErr(t, err, "ClaimIssuer.AddKey")
	mine(tx, "ClaimIssuer.AddKey")

	// deploySuite deploys one token suite through the factory, reusing the SAME code
	// path as production (buildTokenDetails + defineClaimSuiteDetails). Passing a
	// non-zero irs makes the factory reuse an existing IdentityRegistryStorage.
	deploySuite := func(salt, name, symbol string, irs common.Address) (token *contracts.Token, ir *contracts.IdentityRegistry, irAddr, irsAddr common.Address) {
		td := buildTokenDetails(deployer, name, symbol, 18, irs, []common.Address{})
		cd := defineClaimSuiteDetails(claimIssuerAddr)

		tx, err := trexFactory.DeployTREXSuite(auth, salt, td, cd)
		failOnErr(t, err, "DeployTREXSuite "+salt)
		mine(tx, "DeployTREXSuite "+salt)

		tokenAddr, err := trexFactory.GetToken(co, salt)
		failOnErr(t, err, "GetToken "+salt)
		if tokenAddr == (common.Address{}) {
			t.Fatalf("token %s not registered by the factory", salt)
		}
		token, err = contracts.NewToken(tokenAddr, client)
		failOnErr(t, err, "bind token "+salt)
		irAddr, err = token.IdentityRegistry(co)
		failOnErr(t, err, "token.identityRegistry "+salt)
		ir, err = contracts.NewIdentityRegistry(irAddr, client)
		failOnErr(t, err, "bind IR "+salt)
		irsAddr, err = ir.IdentityStorage(co)
		failOnErr(t, err, "ir.identityStorage "+salt)
		return
	}

	// ---- Token A: fresh IRS (factory deploys + binds A's IR to it) ----------
	t.Log("deploying token A with a fresh IRS")
	tokenA, irA, irAAddr, sharedIRS := deploySuite("SharedTokenA", "Tooken Property A", "TKPA", common.Address{})
	if sharedIRS == (common.Address{}) {
		t.Fatal("token A must have a non-zero IdentityRegistryStorage")
	}

	// ---- Token B: REUSE the shared IRS (this is the CreateToken reuse path) --
	t.Log("deploying token B reusing token A's IRS (shared whitelist)")
	tokenB, irB, irBAddr, irsB := deploySuite("SharedTokenB", "Tooken Property B", "TKPB", sharedIRS)
	if irsB != sharedIRS {
		t.Fatalf("token B must reuse the shared IRS: got %s, want %s", irsB.Hex(), sharedIRS.Hex())
	}
	t.Logf("  shared IRS = %s (token A IR=%s, token B IR=%s)", sharedIRS.Hex(), irAAddr.Hex(), irBAddr.Hex())

	// The factory retains ownership of the shared IRS and binds every token's IR to it,
	// so BOTH token A's and token B's IRs must be agents of the shared storage. This is
	// what makes the shared whitelist writable from any token and readable by all.
	irsInstance, err := contracts.NewIdentityRegistryStorage(sharedIRS, client)
	failOnErr(t, err, "bind shared IRS")
	linked, err := irsInstance.LinkedIdentityRegistries(co)
	failOnErr(t, err, "LinkedIdentityRegistries")
	contains := func(list []common.Address, a common.Address) bool {
		for _, x := range list {
			if x == a {
				return true
			}
		}
		return false
	}
	if !contains(linked, irAAddr) {
		t.Fatalf("token A IR %s must be bound to the shared IRS, linked=%v", irAAddr.Hex(), linked)
	}
	if !contains(linked, irBAddr) {
		t.Fatalf("token B IR %s must also be bound to the shared IRS on reuse, linked=%v", irBAddr.Hex(), linked)
	}
	t.Logf("  both token IRs are bound to the shared IRS: %v", linked)

	// ---- Unpause both tokens ------------------------------------------------
	tx, err = tokenA.Unpause(auth)
	failOnErr(t, err, "unpause token A")
	mine(tx, "unpause token A")
	tx, err = tokenB.Unpause(auth)
	failOnErr(t, err, "unpause token B")
	mine(tx, "unpause token B")

	// ---- KYC an investor ONCE, through the bound IR (token A) ---------------
	t.Log("registering + KYC'ing an investor once, via the bound IR")
	investorKey, _ := crypto.GenerateKey()
	investor := crypto.PubkeyToAddress(investorKey.PublicKey)

	investorIdAddr, tx, investorID, err := contracts.DeployIdentity(auth, client, deployer, false)
	failOnErr(t, err, "deploy investor Identity")
	mine(tx, "investor Identity")
	tx, err = investorID.AddKey(auth, keyHash(deployer), big.NewInt(3), big.NewInt(1))
	failOnErr(t, err, "investorID.AddKey")
	mine(tx, "investorID.AddKey")
	tx, err = irA.RegisterIdentity(auth, investor, investorIdAddr, countryCode)
	failOnErr(t, err, "RegisterIdentity via token A IR")
	mine(tx, "RegisterIdentity")
	sig, err := generateSignatureAddClaim(investorIdAddr, kycTopic)
	failOnErr(t, err, "generateSignatureAddClaim")
	tx, err = investorID.AddClaim(auth, big.NewInt(kycTopic), big.NewInt(1), claimIssuerAddr, sig.SignatureBytes, []byte(claimData), "")
	failOnErr(t, err, "AddClaim")
	mine(tx, "AddClaim")

	// ---- The shared whitelist must make the investor verified on BOTH tokens -
	vA, err := irA.IsVerified(co, investor)
	failOnErr(t, err, "IsVerified token A")
	if !vA {
		t.Fatal("investor must be verified on token A")
	}
	vB, err := irB.IsVerified(co, investor)
	failOnErr(t, err, "IsVerified token B")
	if !vB {
		t.Fatal("investor must be verified on token B via the SHARED IRS — reuse is broken")
	}
	t.Log("  investor verified on BOTH token A and token B ✓")

	// ---- Mint on BOTH tokens must succeed for the single KYC'd investor -----
	wei18 := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	amount := new(big.Int).Mul(big.NewInt(500), wei18)

	for _, tc := range []struct {
		name  string
		token *contracts.Token
	}{{"token A", tokenA}, {"token B", tokenB}} {
		tx, err = tc.token.Mint(auth, investor, amount)
		failOnErr(t, err, "Mint "+tc.name)
		mine(tx, "Mint "+tc.name)
		bal, err := tc.token.BalanceOf(co, investor)
		failOnErr(t, err, "BalanceOf "+tc.name)
		if bal.Cmp(amount) != 0 {
			t.Fatalf("%s: expected balance %s, got %s", tc.name, amount, bal)
		}
		t.Logf("  ✅ %s minted %s to the shared-whitelist investor", tc.name, amount)
	}
}
