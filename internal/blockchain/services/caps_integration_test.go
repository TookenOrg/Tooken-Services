//go:build integration

// Integration tests of the T-REX caps: the per-investor cap enforced by
// MaxBalanceModule and the supply cap enforced by SupplyLimitModule.
//
// Tokens are created through CreateTokenWithSalt, against a node and a private
// database, so the modules are deployed, initialized and registered by the very code
// production runs. The fixture is only used to build what production refuses to: a
// suite with its modules in the wrong order, or without their settings.
//
// Requires a node and a PostgreSQL user allowed to create databases:
//
//	cd tools/hardhat && npx hardhat node
//	TEST_DATABASE_URL="postgres://…" go test -tags integration -run TestCapsIntegration ./internal/blockchain/services/
package services

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	"github.com/TookenOrg/tooken-services/internal/blockchain/database"
	chainglobals "github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	appglobals "github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// shares turns a whole number of shares into base units at 18 decimals.
func shares(n int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(n), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}

// sendAndWait reports a transaction refused by the node, at estimation or on-chain, as
// an error. Bindings estimate gas before sending, so a revert usually surfaces there,
// with its reason; a transaction that is sent anyway can still revert when mined.
func sendAndWait(ctx context.Context, client bind.DeployBackend, tx *types.Transaction, err error) error {
	if err != nil {
		return err
	}
	rcpt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return err
	}
	if rcpt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("transaction %s reverted", tx.Hash().Hex())
	}
	return nil
}

func TestCapsIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	env := newChainDBEnv(ctx, t)
	f := env.Fixture
	svc := NewService()

	// What an environment installed before the capping modules existed looks like to
	// CreateTokenWithSalt: a factory and a claim issuer, and none of the capping modules.
	env.SetContractRole(t, chainglobals.TrexFactoryName, f.FactoryAddr)

	cappedModules := []string{chainglobals.MaxBalanceModuleName, chainglobals.SupplyLimitModuleName}

	countRows := func(t *testing.T, query string, arg string) int {
		t.Helper()
		var n int
		if err := env.DB.QueryRow(query, arg).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}
	roleRows := func(t *testing.T, name string) int {
		t.Helper()
		return countRows(t, `SELECT count(*) FROM blk.contract_role WHERE contract_name = $1`, name)
	}
	txRows := func(t *testing.T, txName string) int {
		t.Helper()
		return countRows(t, `SELECT count(*) FROM blk.eth_transaction WHERE tx_name = $1`, txName)
	}
	moduleAddr := func(t *testing.T, name string) common.Address {
		t.Helper()
		role, err := database.GetContractRoleByName(ctx, name)
		if err != nil {
			t.Fatalf("%s must be registered: %v", name, err)
		}
		return common.HexToAddress(role.Address)
	}
	createToken := func(t *testing.T, name, symbol, salt, total, perInvestor string) server.TokenInfos {
		t.Helper()
		_, token, err := svc.CreateTokenWithSalt(ctx, server.CreateTokenRequest{
			TokenName: name, Symbol: symbol, NbDecimal: 18,
			TotalShares: total, MaxSharesPerInvestor: perInvestor,
		}, salt)
		if err != nil {
			t.Fatalf("CreateTokenWithSalt %s: %v", salt, err)
		}
		return token
	}
	balanceOf := func(t *testing.T, token string, wallet common.Address) *big.Int {
		t.Helper()
		instance, err := contracts.NewToken(common.HexToAddress(token), f.Client)
		failOnErr(t, err, "bind token")
		balance, err := instance.BalanceOf(f.Call, wallet)
		failOnErr(t, err, "BalanceOf")
		return balance
	}
	// newInvestor walks the production onboarding: an ONCHAINID, registered in the
	// shared whitelist, holding a KYC claim.
	newInvestor := func(t *testing.T) (common.Address, common.Address) {
		t.Helper()
		wallet, identityAddr, identity := f.NewInvestorIdentity(t)
		if _, err := registerIdentity(ctx, wallet, identityAddr, int(countryCode)); err != nil {
			t.Fatalf("registerIdentity: %v", err)
		}
		f.AddKYCClaim(t, identityAddr, identity)
		return wallet, identityAddr
	}

	for _, name := range cappedModules {
		if n := roleRows(t, name); n != 0 {
			t.Fatalf("the environment must start without %s, found %d rows", name, n)
		}
	}

	// Token A: 100 shares, at most 10 per investor.
	tokenA := createToken(t, "Tooken Caps A", "TKCA", "caps-a", "100", "10")
	mcA := common.HexToAddress(*tokenA.ModularComplianceAddr)

	// ---- the first token brings the modules in, the next ones reuse them ----------

	t.Run("the first token deploys, initializes and registers each capping module", func(t *testing.T) {
		for _, name := range cappedModules {
			if n := roleRows(t, name); n != 1 {
				t.Fatalf("%s: want one contract_role row, found %d", name, n)
			}
			if n := txRows(t, "DEPLOY_"+name); n != 1 {
				t.Fatalf("%s: want one DEPLOY_ row in eth_transaction, found %d", name, n)
			}
			if n := txRows(t, "INITIALIZE_"+name); n != 1 {
				t.Fatalf("%s: want one INITIALIZE_ row in eth_transaction, found %d", name, n)
			}
		}

		// The transfer restriction module goes through the same path, so the same rows.
		if n := txRows(t, "INITIALIZE_TRANSFER_RESTRICT_MODULE"); n != 1 {
			t.Fatalf("TRANSFER_RESTRICTION_MODULE: want one INITIALIZE_ row in eth_transaction, found %d", n)
		}

		// Initialized by the platform, so nobody else can claim the module.
		trm, err := contracts.NewTransferRestrictModule(moduleAddr(t, chainglobals.TransferRestrictionModuleName), f.Client)
		failOnErr(t, err, "bind TransferRestrictModule")
		mb, err := contracts.NewMaxBalanceModule(moduleAddr(t, chainglobals.MaxBalanceModuleName), f.Client)
		failOnErr(t, err, "bind MaxBalanceModule")
		sl, err := contracts.NewSupplyLimitModule(moduleAddr(t, chainglobals.SupplyLimitModuleName), f.Client)
		failOnErr(t, err, "bind SupplyLimitModule")
		for name, owner := range map[string]func(*bind.CallOpts) (common.Address, error){
			"MaxBalanceModule":       mb.Owner,
			"SupplyLimitModule":      sl.Owner,
			"TransferRestrictModule": trm.Owner,
		} {
			got, err := owner(f.Call)
			failOnErr(t, err, name+".owner")
			if got != f.Deployer {
				t.Fatalf("%s owner = %s, want the platform %s", name, got.Hex(), f.Deployer.Hex())
			}
		}
	})

	// Token B: 20 shares, at most 15 per investor — small enough to reach the supply
	// cap with two investors, and a cap above token A's to tell them apart.
	tokenB := createToken(t, "Tooken Caps B", "TKCB", "caps-b", "20", "15")
	mcB := common.HexToAddress(*tokenB.ModularComplianceAddr)

	t.Run("a second token reuses the registered modules", func(t *testing.T) {
		for _, name := range cappedModules {
			if n := roleRows(t, name); n != 1 {
				t.Fatalf("%s: want still one contract_role row, found %d", name, n)
			}
			if n := txRows(t, "DEPLOY_"+name); n != 1 {
				t.Fatalf("%s: a second token must not deploy again, found %d DEPLOY_ rows", name, n)
			}
			got, err := ensureSharedComplianceModule(ctx, name)
			failOnErr(t, err, "ensureSharedComplianceModule "+name)
			if got != moduleAddr(t, name) {
				t.Fatalf("%s: ensure answered %s instead of the registered module", name, got.Hex())
			}
		}
	})

	t.Run("a database that cannot answer deploys nothing", func(t *testing.T) {
		broken, err := sql.Open("postgres", "postgres://nobody@127.0.0.1:1/none?sslmode=disable")
		failOnErr(t, err, "open")
		broken.Close()

		nonceBefore, err := f.Client.PendingNonceAt(ctx, f.Deployer)
		failOnErr(t, err, "nonce before")

		previous := appglobals.DB
		appglobals.DB = broken
		_, ensureErr := ensureSharedComplianceModule(ctx, chainglobals.MaxBalanceModuleName)
		appglobals.DB = previous

		if ensureErr == nil {
			t.Fatal("a failed lookup must be reported, not read as « no module yet »")
		}
		nonceAfter, err := f.Client.PendingNonceAt(ctx, f.Deployer)
		failOnErr(t, err, "nonce after")
		if nonceAfter != nonceBefore {
			t.Fatalf("no transaction may be sent on a failed lookup: nonce went from %d to %d", nonceBefore, nonceAfter)
		}
	})

	t.Run("a second registration of a module is recognised as taken", func(t *testing.T) {
		// The constraint name lives in two places — the schema and the Go constant —
		// and only the real schema can tell whether they still agree.
		_, err := database.InsertContractRole(ctx, probeTxHash("caps-race"), common.Address{1}.Hex(), chainglobals.MaxBalanceModuleName)
		if err == nil {
			t.Fatal("a second row for the same role must be refused")
		}
		if !database.IsContractRoleNameTaken(err) {
			t.Fatalf("the refusal must be recognised as « role already registered », got: %v", err)
		}
	})

	// ---- the caps are set on-chain, in base units ---------------------------------

	investor1, identity1 := newInvestor(t)
	mb, err := contracts.NewMaxBalanceModule(moduleAddr(t, chainglobals.MaxBalanceModuleName), f.Client)
	failOnErr(t, err, "bind MaxBalanceModule")
	sl, err := contracts.NewSupplyLimitModule(moduleAddr(t, chainglobals.SupplyLimitModuleName), f.Client)
	failOnErr(t, err, "bind SupplyLimitModule")

	t.Run("the caps read back on-chain in base units", func(t *testing.T) {
		for _, c := range []struct {
			mc   common.Address
			want *big.Int
		}{{mcA, shares(100)}, {mcB, shares(20)}} {
			got, err := sl.GetSupplyLimit(f.Call, c.mc)
			failOnErr(t, err, "getSupplyLimit")
			if got.Cmp(c.want) != 0 {
				t.Fatalf("supply limit of %s = %s, want %s", c.mc.Hex(), got, c.want)
			}
		}

		// MaxBalanceModule exposes no getter for the cap: its own check is the way
		// to read it — exactly the cap passes, one base unit more does not.
		for _, c := range []struct {
			mc  common.Address
			cap *big.Int
		}{{mcA, shares(10)}, {mcB, shares(15)}} {
			ok, err := mb.ModuleCheck(f.Call, common.Address{}, investor1, c.cap, c.mc)
			failOnErr(t, err, "moduleCheck at the cap")
			if !ok {
				t.Fatalf("%s: a mint of exactly the cap (%s) must pass", c.mc.Hex(), c.cap)
			}
			above := new(big.Int).Add(c.cap, big.NewInt(1))
			ok, err = mb.ModuleCheck(f.Call, common.Address{}, investor1, above, c.mc)
			failOnErr(t, err, "moduleCheck above the cap")
			if ok {
				t.Fatalf("%s: a mint of the cap plus one base unit must not pass", c.mc.Hex())
			}
		}
	})

	// ---- the per-investor cap -----------------------------------------------------

	t.Run("a mint up to the per-investor cap succeeds", func(t *testing.T) {
		if _, err := mintSharesForTest(ctx, svc, tokenA.Address, investor1.Hex(), 10); err != nil {
			t.Fatalf("minting exactly the cap must succeed: %v", err)
		}
		if got := balanceOf(t, tokenA.Address, investor1); got.Cmp(shares(10)) != 0 {
			t.Fatalf("balance = %s, want %s", got, shares(10))
		}
	})

	t.Run("one share above the cap reverts and changes nothing", func(t *testing.T) {
		_, err := mintSharesForTest(ctx, svc, tokenA.Address, investor1.Hex(), 1)
		if err == nil {
			t.Fatal("a mint above the per-investor cap must be refused")
		}
		if !strings.Contains(err.Error(), "Compliance not followed") {
			t.Fatalf("the refusal must come from the compliance, got: %v", err)
		}
		if got := balanceOf(t, tokenA.Address, investor1); got.Cmp(shares(10)) != 0 {
			t.Fatalf("balance after the refusal = %s, want it unchanged at %s", got, shares(10))
		}
		held, err := mb.GetIDBalance(f.Call, mcA, identity1)
		failOnErr(t, err, "getIDBalance")
		if held.Cmp(shares(10)) != 0 {
			t.Fatalf("the module counts %s for the identity, want %s", held, shares(10))
		}
	})

	// ---- the cap belongs to the identity, not to the wallet -----------------------

	investor2, _ := newInvestor(t)

	t.Run("a second wallet of the same identity shares its cap", func(t *testing.T) {
		key, err := crypto.GenerateKey()
		failOnErr(t, err, "generate key")
		secondWallet := crypto.PubkeyToAddress(key.PublicKey)
		if _, err := registerIdentity(ctx, secondWallet, identity1, int(countryCode)); err != nil {
			t.Fatalf("registering a second wallet for the same identity: %v", err)
		}

		if _, err := mintSharesForTest(ctx, svc, tokenA.Address, secondWallet.Hex(), 1); err == nil {
			t.Fatal("the identity already holds its cap: its second wallet must be refused")
		}
		if got := balanceOf(t, tokenA.Address, secondWallet); got.Sign() != 0 {
			t.Fatalf("second wallet balance = %s, want 0", got)
		}

		// The control: on the same token, another identity is still free to buy.
		// Without it, a token refusing everybody would pass the assertion above.
		if _, err := mintSharesForTest(ctx, svc, tokenA.Address, investor2.Hex(), 1); err != nil {
			t.Fatalf("another identity must not be bound by the first one's cap: %v", err)
		}
	})

	// ---- one cap per token ---------------------------------------------------------

	t.Run("the same investor has an independent cap on another token", func(t *testing.T) {
		// Full on token A, yet token B lets them buy up to B's own cap.
		if _, err := mintSharesForTest(ctx, svc, tokenB.Address, investor1.Hex(), 15); err != nil {
			t.Fatalf("token B's cap must not be consumed by token A: %v", err)
		}
		if got := balanceOf(t, tokenB.Address, investor1); got.Cmp(shares(15)) != 0 {
			t.Fatalf("balance on token B = %s, want %s", got, shares(15))
		}
	})

	// ---- the supply cap ------------------------------------------------------------

	t.Run("mints add up to the total, and one share more reverts", func(t *testing.T) {
		// 15 already minted on token B: investor2 takes the last 5.
		if _, err := mintSharesForTest(ctx, svc, tokenB.Address, investor2.Hex(), 5); err != nil {
			t.Fatalf("minting up to the total supply must succeed: %v", err)
		}

		// investor2 holds 5 out of an allowed 15: only the supply cap can refuse now.
		ok, err := mb.ModuleCheck(f.Call, common.Address{}, investor2, shares(1), mcB)
		failOnErr(t, err, "moduleCheck")
		if !ok {
			t.Fatal("precondition: the per-investor cap must not be the one refusing")
		}
		ok, err = sl.ModuleCheck(f.Call, common.Address{}, investor2, shares(1), mcB)
		failOnErr(t, err, "supply moduleCheck")
		if ok {
			t.Fatal("SupplyLimitModule must refuse a mint beyond the total supply")
		}

		if _, err := mintSharesForTest(ctx, svc, tokenB.Address, investor2.Hex(), 1); err == nil {
			t.Fatal("a mint beyond the total supply must be refused")
		}
		token, err := contracts.NewToken(common.HexToAddress(tokenB.Address), f.Client)
		failOnErr(t, err, "bind token B")
		supply, err := token.TotalSupply(f.Call)
		failOnErr(t, err, "totalSupply")
		if supply.Cmp(shares(20)) != 0 {
			t.Fatalf("total supply = %s, want exactly %s", supply, shares(20))
		}
	})

	// ---- what production must never build ------------------------------------------

	sharedIRS := moduleAddr(t, chainglobals.SharedIdentityRegistryStorageName)
	mbAddr := moduleAddr(t, chainglobals.MaxBalanceModuleName)
	slAddr := moduleAddr(t, chainglobals.SupplyLimitModuleName)
	trmAddr := moduleAddr(t, chainglobals.TransferRestrictionModuleName)
	big1000 := shares(1000)
	setMaxBalance, setSupplyLimit, err := encodeCapSettings(big1000, big1000)
	failOnErr(t, err, "encodeCapSettings")

	// deployRaw sends a suite deployment through the factory and reports a revert as
	// an error instead of failing the test.
	deployRaw := func(salt string, modules []common.Address, settings [][]byte) error {
		details := buildTokenDetails(f.Deployer, "Tooken "+salt, strings.ToUpper(salt), 18, sharedIRS, modules, settings)
		tx, err := f.Factory.DeployTREXSuite(f.Auth, salt, details, defineClaimSuiteDetails(f.ClaimIssuer))
		return sendAndWait(ctx, f.Client, tx, err)
	}

	t.Run("the modules out of order make the deployment revert", func(t *testing.T) {
		// The factory sends settings[0], setMaxBalance, to modules[0]: here the
		// transfer restriction module, which has no such function and no fallback.
		err := deployRaw("caps-b1", []common.Address{trmAddr, mbAddr, slAddr}, [][]byte{setMaxBalance, setSupplyLimit})
		if err == nil {
			t.Fatal("a setting sent to the wrong module must make the deployment revert")
		}
		if token, err := f.Factory.GetToken(f.Call, "caps-b1"); err != nil || token != (common.Address{}) {
			t.Fatalf("no token may exist under the salt, got %s (%v)", token.Hex(), err)
		}

		// The control: the same call in production order goes through, so the
		// revert above is the order and nothing else.
		if err := deployRaw("caps-b1-ok", []common.Address{mbAddr, slAddr, trmAddr}, [][]byte{setMaxBalance, setSupplyLimit}); err != nil {
			t.Fatalf("the same deployment in production order must succeed: %v", err)
		}
	})

	t.Run("a capping module left without its setting blocks every mint", func(t *testing.T) {
		mintOne := func(t *testing.T, suite *deployedSuite) error {
			t.Helper()
			tx, err := suite.Token.Mint(f.Auth, investor2, shares(1))
			return sendAndWait(ctx, f.Client, tx, err)
		}

		for _, c := range []struct {
			name     string
			salt     string
			modules  []common.Address
			settings [][]byte
			wantMint bool
		}{
			{"supply limit unset", "caps-b2-sl", []common.Address{mbAddr, slAddr}, [][]byte{setMaxBalance}, false},
			{"max balance unset", "caps-b2-mb", []common.Address{slAddr, mbAddr}, [][]byte{setSupplyLimit}, false},
			// The control: both set, the same mint goes through.
			{"both set", "caps-b2-ok", []common.Address{mbAddr, slAddr}, [][]byte{setMaxBalance, setSupplyLimit}, true},
		} {
			t.Run(c.name, func(t *testing.T) {
				suite := f.DeploySuiteWith(t, c.salt, "Tooken "+c.salt, strings.ToUpper(c.salt), sharedIRS, c.modules, c.settings)
				err := mintOne(t, suite)
				if c.wantMint && err != nil {
					t.Fatalf("with both caps set, a one-share mint must succeed: %v", err)
				}
				if !c.wantMint && err == nil {
					t.Fatal("a cap left at zero must refuse a one-share mint")
				}
			})
		}
	})

	// ---- the setters only answer the token's compliance -----------------------------

	t.Run("the platform key cannot set a cap directly", func(t *testing.T) {
		mbTx, err := contracts.NewMaxBalanceModuleTransactor(mbAddr, f.Client)
		failOnErr(t, err, "bind MaxBalanceModule transactor")
		slTx, err := contracts.NewSupplyLimitModuleTransactor(slAddr, f.Client)
		failOnErr(t, err, "bind SupplyLimitModule transactor")

		for name, send := range map[string]func() (*types.Transaction, error){
			"setMaxBalance":  func() (*types.Transaction, error) { return mbTx.SetMaxBalance(f.Auth, big1000) },
			"setSupplyLimit": func() (*types.Transaction, error) { return slTx.SetSupplyLimit(f.Auth, big1000) },
		} {
			tx, err := send()
			err = sendAndWait(ctx, f.Client, tx, err)
			if err == nil {
				t.Fatalf("%s from the platform key must revert", name)
			}
			if !strings.Contains(err.Error(), "only bound compliance can call") {
				t.Fatalf("%s: the revert must be the compliance guard, got: %v", name, err)
			}
		}
	})
}
