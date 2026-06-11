//go:build integration

package services

import (
	"context"
	"fmt"
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

const (
	kycTopic    = int64(7)    // claim topic required to hold the token (KYC)
	countryCode = uint16(250) // ISO-3166 numeric (France) for the investor
)

// keyHash returns the ERC-734 key id for an address: keccak256(abi.encode(addr)).
func keyHash(addr common.Address) [32]byte {
	return [32]byte(crypto.Keccak256Hash(common.LeftPadBytes(addr.Bytes(), 32)))
}

// TestTokenLifecycleIntegration reproduces the full ERC-3643 / T-REX flow on a local
// EVM (Hardhat), directly through the contract bindings (no DB), to prove that:
//
//  1. the T-REX infrastructure can be deployed and configured,
//  2. a security token can be issued by the shared factory,
//  3. a KYC'd investor identity can be registered, and
//  4. that investor wallet can then RECEIVE shares of the token (mint succeeds and
//     the balance reflects it) — i.e. the permissioning chain works end to end.
//
// Gated by the `integration` build tag. Requires a local node:
//
//	cd tools/hardhat && npm install && npx hardhat node
//	go test -tags integration -run TestTokenLifecycleIntegration -v ./internal/blockchain/...
func TestTokenLifecycleIntegration(t *testing.T) {
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

	// Deployer = platform manager (owner & agent of everything). Hardhat account #0.
	const deployerKeyHex = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	deployerKey, err := crypto.HexToECDSA(deployerKeyHex)
	if err != nil {
		t.Fatalf("parse deployer key: %v", err)
	}
	deployer := crypto.PubkeyToAddress(deployerKey.PublicKey)
	// generateSignatureAddClaim signs with PRIVATE_KEY; use the same deployer key.
	t.Setenv("PRIVATE_KEY", deployerKeyHex)

	// One shared transactor; leaving Nonce nil + GasLimit 0 lets go-ethereum
	// auto-manage the nonce and estimate gas for every sequential transaction.
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

	// ---- Step 1: deploy the 6 T-REX logic implementations ------------------
	t.Log("STEP 1 — deploying T-REX implementations (CTR, TIR, IRS, IR, MC, Token)")
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

	// ---- Step 2: ONCHAINID implementation + its ImplementationAuthority -----
	t.Log("STEP 2 — deploying ONCHAINID identity implementation + authority + factory")
	identityImpl, tx, _, err := contracts.DeployIdentity(auth, client, deployer, true)
	failOnErr(t, err, "deploy Identity impl")
	mine(tx, "Identity impl")
	idAuthority, tx, _, err := contracts.DeployImplementationAuthority(auth, client, identityImpl)
	failOnErr(t, err, "deploy identity ImplementationAuthority")
	mine(tx, "identity authority")
	idFactoryAddr, tx, idFactory, err := contracts.DeployIdFactory(auth, client, idAuthority)
	failOnErr(t, err, "deploy IdFactory")
	mine(tx, "IdFactory")

	// ---- Step 3: TREX ImplementationAuthority + register the version --------
	t.Log("STEP 3 — deploying TREX authority and registering the implementations version")
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

	// ---- Step 4: TREX factory, wired to the authority + identity factory ----
	t.Log("STEP 4 — deploying TREX factory and wiring it (addTokenFactory + setTREXFactory)")
	trexFactoryAddr, tx, trexFactory, err := contracts.DeployTREXFactory(auth, client, trexAuthorityAddr, idFactoryAddr)
	failOnErr(t, err, "deploy TREX factory")
	mine(tx, "TREX factory")
	tx, err = idFactory.AddTokenFactory(auth, trexFactoryAddr)
	failOnErr(t, err, "addTokenFactory")
	mine(tx, "addTokenFactory")
	tx, err = trexAuthority.SetTREXFactory(auth, trexFactoryAddr)
	failOnErr(t, err, "setTREXFactory")
	mine(tx, "setTREXFactory")

	// ---- Step 5: the single shared claim issuer (trusted for KYC) -----------
	t.Log("STEP 5 — deploying the shared ClaimIssuer and giving the deployer a CLAIM key")
	claimIssuerAddr, tx, claimIssuer, err := contracts.DeployClaimIssuer(auth, client, deployer)
	failOnErr(t, err, "deploy ClaimIssuer")
	mine(tx, "ClaimIssuer")
	tx, err = claimIssuer.AddKey(auth, keyHash(deployer), big.NewInt(3), big.NewInt(1)) // purpose 3 = CLAIM
	failOnErr(t, err, "ClaimIssuer.AddKey")
	mine(tx, "ClaimIssuer.AddKey")

	// ---- Step 6: deploy the token suite via the factory ---------------------
	t.Log("STEP 6 — deploying the token suite via deployTREXSuite (ONCHAINID=0, KYC topic 7)")
	salt := "E2ETokenLifecycle"
	tx, err = trexFactory.DeployTREXSuite(auth, salt,
		contracts.ITREXFactoryTokenDetails{
			Owner:              deployer,
			Name:               "Tooken Property #1",
			Symbol:             "TKP1",
			Decimals:           18,
			Irs:                common.Address{},
			ONCHAINID:          common.Address{},
			IrAgents:           []common.Address{deployer},
			TokenAgents:        []common.Address{deployer},
			ComplianceModules:  []common.Address{},
			ComplianceSettings: [][]byte{},
		},
		contracts.ITREXFactoryClaimDetails{
			ClaimTopics:  []*big.Int{big.NewInt(kycTopic)},
			Issuers:      []common.Address{claimIssuerAddr},
			IssuerClaims: [][]*big.Int{{big.NewInt(kycTopic)}},
		})
	failOnErr(t, err, "DeployTREXSuite")
	mine(tx, "DeployTREXSuite")

	tokenAddr, err := trexFactory.GetToken(co, salt)
	failOnErr(t, err, "GetToken")
	if tokenAddr == (common.Address{}) {
		t.Fatal("token not registered by the factory")
	}
	token, err := contracts.NewToken(tokenAddr, client)
	failOnErr(t, err, "bind token")
	irAddr, err := token.IdentityRegistry(co)
	failOnErr(t, err, "token.identityRegistry")
	ir, err := contracts.NewIdentityRegistry(irAddr, client)
	failOnErr(t, err, "bind IR")
	t.Logf("  token=%s  identityRegistry=%s", tokenAddr.Hex(), irAddr.Hex())

	// ---- Step 7: a KYC'd investor identity, registered in the IR ------------
	t.Log("STEP 7 — creating an investor identity, registering it and adding a KYC claim")
	investorKey, _ := crypto.GenerateKey()
	investor := crypto.PubkeyToAddress(investorKey.PublicKey)

	investorIdAddr, tx, investorID, err := contracts.DeployIdentity(auth, client, deployer, false)
	failOnErr(t, err, "deploy investor Identity")
	mine(tx, "investor Identity")
	// deployer needs a CLAIM key on the investor identity to be allowed to addClaim.
	tx, err = investorID.AddKey(auth, keyHash(deployer), big.NewInt(3), big.NewInt(1))
	failOnErr(t, err, "investorID.AddKey")
	mine(tx, "investorID.AddKey")
	// register the wallet -> identity link (deployer is an IR agent).
	tx, err = ir.RegisterIdentity(auth, investor, investorIdAddr, countryCode)
	failOnErr(t, err, "RegisterIdentity")
	mine(tx, "RegisterIdentity")
	// sign + add the KYC claim (issued by the shared claim issuer).
	sig, err := generateSignatureAddClaim(investorIdAddr, kycTopic)
	failOnErr(t, err, "generateSignatureAddClaim")
	tx, err = investorID.AddClaim(auth, big.NewInt(kycTopic), big.NewInt(1), claimIssuerAddr, sig.SignatureBytes, []byte(claimData), "")
	failOnErr(t, err, "AddClaim")
	mine(tx, "AddClaim")

	// ---- Step 8: unpause, then verify the investor is eligible --------------
	t.Log("STEP 8 — unpausing the token and checking the investor is verified")
	tx, err = token.Unpause(auth)
	failOnErr(t, err, "Unpause")
	mine(tx, "Unpause")

	verified, err := ir.IsVerified(co, investor)
	failOnErr(t, err, "IsVerified")
	if !verified {
		t.Fatal("investor is NOT verified — KYC claim / registration failed")
	}
	t.Logf("  investor %s is verified ✓", investor.Hex())

	// ---- Step 9: mint shares to the investor and assert the balance ---------
	t.Log("STEP 9 — minting shares to the investor wallet and asserting the balance")
	wei18 := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	amount := new(big.Int).Mul(big.NewInt(1000), wei18) // 1000 tokens

	before, err := token.BalanceOf(co, investor)
	failOnErr(t, err, "BalanceOf before")
	tx, err = token.Mint(auth, investor, amount)
	failOnErr(t, err, "Mint")
	mine(tx, "Mint")
	after, err := token.BalanceOf(co, investor)
	failOnErr(t, err, "BalanceOf after")

	gained := new(big.Int).Sub(after, before)
	if gained.Cmp(amount) != 0 {
		t.Fatalf("investor balance: expected +%s, got +%s", amount, gained)
	}
	t.Logf("  ✅ investor received %s base units (balance %s)", amount, after)

	// ---- Step 10: a NON-KYC investor must be rejected ----------------------
	t.Log("STEP 10 — a non-KYC investor must NOT be able to receive shares")
	investor2Key, _ := crypto.GenerateKey()
	investor2 := crypto.PubkeyToAddress(investor2Key.PublicKey)

	// Sanity: this wallet was never registered nor KYC'd.
	v2, err := ir.IsVerified(co, investor2)
	failOnErr(t, err, "IsVerified investor2")
	if v2 {
		t.Fatal("investor2 should NOT be verified")
	}

	// Attempt to mint to it — must fail, either at gas estimation (the binding
	// returns an error) or, if sent, on-chain (the transaction reverts).
	tx2, mintErr := token.Mint(auth, investor2, amount)
	if mintErr == nil {
		rcpt, werr := bind.WaitMined(ctx, client, tx2)
		switch {
		case werr != nil:
			mintErr = werr
		case rcpt.Status != types.ReceiptStatusSuccessful:
			mintErr = fmt.Errorf("transaction reverted (status 0)")
		}
	}
	if mintErr == nil {
		t.Fatal("expected mint to a non-KYC investor to FAIL, but it succeeded")
	}
	t.Logf("  mint to non-KYC investor correctly rejected: %v", mintErr)

	// The non-KYC wallet must still hold nothing.
	bal2, err := token.BalanceOf(co, investor2)
	failOnErr(t, err, "BalanceOf investor2")
	if bal2.Sign() != 0 {
		t.Fatalf("non-KYC investor balance should be 0, got %s", bal2)
	}
	t.Logf("  ✅ non-KYC investor balance is 0 as expected")
}

func failOnErr(t *testing.T, err error, what string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", what, err)
	}
}
