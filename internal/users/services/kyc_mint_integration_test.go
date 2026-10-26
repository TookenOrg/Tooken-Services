//go:build integration

package services

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	blkGlobals "github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	blkServices "github.com/TookenOrg/tooken-services/internal/blockchain/services"
	appGlobals "github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
)

const (
	scratchDatabaseName = "tooken_it_users_services_kyc_mint"
	baselineSQLPath     = "../../../tools/db/baseline.sql"
	hardhatAccount0Key  = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

type integrationEnv struct {
	DB     *sql.DB
	Client *ethclient.Client
}

func TestApprovedKYCEnablesMintAndNonKYCWalletReverts(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 900*time.Second)
	defer cancel()

	env := newIntegrationEnv(ctx, t)
	usersSvc := NewService()
	chainSvc := blkServices.NewService()

	token := deployTestToken(ctx, t, chainSvc)
	tokenAddr := common.HexToAddress(token.Address)
	tokenInstance, err := contracts.NewToken(tokenAddr, env.Client)
	if err != nil {
		t.Fatalf("bind token: %v", err)
	}
	irAddr, err := tokenInstance.IdentityRegistry(bindCall(ctx))
	if err != nil {
		t.Fatalf("token identity registry: %v", err)
	}
	identityRegistry, err := contracts.NewIdentityRegistry(irAddr, env.Client)
	if err != nil {
		t.Fatalf("bind identity registry: %v", err)
	}

	adminID := insertUser(t, env.DB, "kyc-mint-admin", "ADMIN")
	aliceID := insertUser(t, env.DB, "kyc-mint-alice", "USER")
	aliceVerificationID := insertApprovedKYC(t, env.DB, aliceID, adminID, "Alice Tooken", "FR")

	t.Run("Alice approved KYC syncs on-chain and can receive minted tokens", func(t *testing.T) {
		response, err := usersSvc.SyncApprovedKycVerificationOnChain(ctx, aliceVerificationID)
		if err != nil {
			t.Fatalf("SyncApprovedKycVerificationOnChain: %v", err)
		}
		if response.KycStatus != "verified" {
			t.Fatalf("kyc_status = %q, want verified", response.KycStatus)
		}
		if response.KycVerifiedAt.IsZero() {
			t.Fatal("kyc_verified_at must be filled")
		}
		if !common.IsHexAddress(response.WalletAddress) || !common.IsHexAddress(response.IdentityAddress) {
			t.Fatalf("sync must return wallet and identity addresses, got %+v", response)
		}
		if response.CountryCode != "FR" || response.CountryNumeric != 250 {
			t.Fatalf("country = %s/%d, want FR/250", response.CountryCode, response.CountryNumeric)
		}

		aliceWallet := common.HexToAddress(response.WalletAddress)
		verified, err := identityRegistry.IsVerified(bindCall(ctx), aliceWallet)
		if err != nil {
			t.Fatalf("IsVerified(alice): %v", err)
		}
		if !verified {
			t.Fatal("Alice must be verified on the token identity registry")
		}
		assertUserVerified(t, env.DB, aliceID)

		before, err := tokenInstance.BalanceOf(bindCall(ctx), aliceWallet)
		if err != nil {
			t.Fatalf("BalanceOf(alice before): %v", err)
		}
		txHash, err := mintShares(ctx, chainSvc, token.Address, aliceWallet.Hex(), 100)
		if err != nil {
			t.Fatalf("Mint(alice): %v", err)
		}
		if txHash == "" {
			t.Fatal("the mint must return its transaction hash")
		}
		after, err := tokenInstance.BalanceOf(bindCall(ctx), aliceWallet)
		if err != nil {
			t.Fatalf("BalanceOf(alice after): %v", err)
		}
		assertBalanceGain(t, before, after, 100)

		activeWalletsBefore := countRows(t, env.DB, `SELECT count(*) FROM blk.user_wallet WHERE user_id = $1 AND is_active`, aliceID)
		identitiesBefore := countRows(t, env.DB, `SELECT count(*) FROM blk.identity WHERE user_id = $1`, aliceID)
		registerTxBefore := countRows(t, env.DB, `SELECT count(*) FROM blk.eth_transaction WHERE tx_name = 'REGISTER_IDENTITY'`)

		secondResponse, err := usersSvc.SyncApprovedKycVerificationOnChain(ctx, aliceVerificationID)
		if err != nil {
			t.Fatalf("second SyncApprovedKycVerificationOnChain: %v", err)
		}
		if secondResponse.Transactions.Claim != nil || secondResponse.Transactions.RegisterIdentity != nil {
			t.Fatalf("idempotent sync must not add claim/register transactions, got %+v", secondResponse.Transactions)
		}
		if got := countRows(t, env.DB, `SELECT count(*) FROM blk.user_wallet WHERE user_id = $1 AND is_active`, aliceID); got != activeWalletsBefore {
			t.Fatalf("active wallet count changed: before=%d after=%d", activeWalletsBefore, got)
		}
		if got := countRows(t, env.DB, `SELECT count(*) FROM blk.identity WHERE user_id = $1`, aliceID); got != identitiesBefore {
			t.Fatalf("identity count changed: before=%d after=%d", identitiesBefore, got)
		}
		if got := countRows(t, env.DB, `SELECT count(*) FROM blk.eth_transaction WHERE tx_name = 'REGISTER_IDENTITY'`); got != registerTxBefore {
			t.Fatalf("register identity tx count changed: before=%d after=%d", registerTxBefore, got)
		}

		beforeSecondMint, err := tokenInstance.BalanceOf(bindCall(ctx), aliceWallet)
		if err != nil {
			t.Fatalf("BalanceOf(alice before second mint): %v", err)
		}
		if _, err := mintShares(ctx, chainSvc, token.Address, aliceWallet.Hex(), 25); err != nil {
			t.Fatalf("second Mint(alice): %v", err)
		}
		afterSecondMint, err := tokenInstance.BalanceOf(bindCall(ctx), aliceWallet)
		if err != nil {
			t.Fatalf("BalanceOf(alice after second mint): %v", err)
		}
		assertBalanceGain(t, beforeSecondMint, afterSecondMint, 25)
	})

	t.Run("Bob without verified KYC is rejected by the token", func(t *testing.T) {
		bobKey, err := crypto.GenerateKey()
		if err != nil {
			t.Fatalf("generate Bob key: %v", err)
		}
		bobWallet := crypto.PubkeyToAddress(bobKey.PublicKey)

		verified, err := identityRegistry.IsVerified(bindCall(ctx), bobWallet)
		if err != nil {
			t.Fatalf("IsVerified(bob): %v", err)
		}
		if verified {
			t.Fatal("Bob must not be verified before mint")
		}

		_, err = mintShares(ctx, chainSvc, token.Address, bobWallet.Hex(), 100)
		if err == nil {
			t.Fatal("Mint(bob) must fail")
		}
		if !strings.Contains(err.Error(), "Identity is not verified") {
			t.Fatalf("Bob mint error must contain the T-REX revert reason, got: %v", err)
		}

		balance, err := tokenInstance.BalanceOf(bindCall(ctx), bobWallet)
		if err != nil {
			t.Fatalf("BalanceOf(bob): %v", err)
		}
		if balance.Sign() != 0 {
			t.Fatalf("Bob balance = %s, want 0", balance)
		}
	})
}

func newIntegrationEnv(ctx context.Context, t *testing.T) *integrationEnv {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	logger.Init(true)

	client := dialTestNode(ctx, t)
	db := createScratchDatabase(t, dsn)

	previousDB := appGlobals.DB
	appGlobals.DB = db
	previousClient := blkGlobals.EthClient
	blkGlobals.EthClient = client

	t.Setenv("PRIVATE_KEY", hardhatAccount0Key)

	t.Cleanup(func() {
		blkGlobals.EthClient = previousClient
		appGlobals.DB = previousDB
		db.Close()
		dropScratchDatabase(t, dsn)
	})

	return &integrationEnv{DB: db, Client: client}
}

func deployTestToken(ctx context.Context, t *testing.T, chainSvc *blkServices.Service) server.TokenInfos {
	t.Helper()

	if _, err := chainSvc.DeployAllImplementations(ctx); err != nil {
		t.Fatalf("DeployAllImplementations: %v", err)
	}
	if _, err := chainSvc.DeployIdentityFactory(ctx); err != nil {
		t.Fatalf("DeployIdentityFactory: %v", err)
	}
	if _, err := chainSvc.ConfigureAuthority(ctx); err != nil {
		t.Fatalf("ConfigureAuthority: %v", err)
	}
	if err := chainSvc.DeployAndInitTrexFactory(ctx); err != nil {
		t.Fatalf("DeployAndInitTrexFactory: %v", err)
	}
	if _, err := chainSvc.DeployTrexSuite(ctx); err != nil {
		t.Fatalf("DeployTrexSuite: %v", err)
	}

	_, token, err := chainSvc.CreateTokenWithSalt(ctx, server.CreateTokenRequest{
		TokenName: "Tooken Villa Belair #12",
		Symbol:    "TVB12",
		NbDecimal: 18,
		// Alice is minted 125 shares below: the per-investor cap must leave room for
		// them, since this test is about KYC, not about caps.
		TotalShares:          "1500",
		MaxSharesPerInvestor: "200",
	}, "re-12")
	if err != nil {
		t.Fatalf("CreateTokenWithSalt: %v", err)
	}
	if !common.IsHexAddress(token.Address) {
		t.Fatalf("token address = %q", token.Address)
	}
	return token
}

func dialTestNode(ctx context.Context, t *testing.T) *ethclient.Client {
	t.Helper()

	rpc := os.Getenv("ETH_TEST_RPC")
	if rpc == "" {
		rpc = "ws://127.0.0.1:8545"
	}
	client, err := ethclient.DialContext(ctx, rpc)
	if err != nil {
		integrationUnavailable(t, "no local EVM node at %s: %v", rpc, err)
	}
	if _, err := client.ChainID(ctx); err != nil {
		integrationUnavailable(t, "local EVM node at %s not reachable: %v", rpc, err)
	}
	return client
}

func createScratchDatabase(t *testing.T, dsn string) *sql.DB {
	t.Helper()

	adminDSN, scratchDSN, err := scratchDSNs(dsn)
	if err != nil {
		integrationUnavailable(t, "TEST_DATABASE_URL is not a URL a database can be derived from: %v", err)
	}

	admin, err := sql.Open("postgres", adminDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	if err := admin.Ping(); err != nil {
		integrationUnavailable(t, "database server not reachable: %v", err)
	}

	if _, err := admin.Exec(`DROP DATABASE IF EXISTS ` + scratchDatabaseName); err != nil {
		integrationUnavailable(t, "cannot drop %s (the test user needs CREATEDB): %v", scratchDatabaseName, err)
	}
	if _, err := admin.Exec(`CREATE DATABASE ` + scratchDatabaseName); err != nil {
		integrationUnavailable(t, "cannot create %s (the test user needs CREATEDB): %v", scratchDatabaseName, err)
	}

	baseline, err := os.ReadFile(baselineSQLPath)
	if err != nil {
		t.Fatalf("read %s: %v", baselineSQLPath, err)
	}

	db, err := sql.Open("postgres", scratchDSN)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(string(baseline)); err != nil {
		db.Close()
		t.Fatalf("apply %s: %v", baselineSQLPath, err)
	}

	return db
}

func dropScratchDatabase(t *testing.T, dsn string) {
	t.Helper()

	adminDSN, _, err := scratchDSNs(dsn)
	if err != nil {
		return
	}
	admin, err := sql.Open("postgres", adminDSN)
	if err != nil {
		return
	}
	defer admin.Close()
	_, _ = admin.Exec(`DROP DATABASE IF EXISTS ` + scratchDatabaseName)
}

func scratchDSNs(dsn string) (adminDSN, scratchDSN string, err error) {
	parsed, err := url.Parse(dsn)
	if err != nil {
		return "", "", err
	}
	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		return "", "", fmt.Errorf("unsupported scheme %q", parsed.Scheme)
	}

	admin := *parsed
	admin.Path = "/postgres"

	scratch := *parsed
	scratch.Path = "/" + scratchDatabaseName

	return admin.String(), scratch.String(), nil
}

func insertUser(t *testing.T, db *sql.DB, emailPrefix, role string) int {
	t.Helper()

	var id int
	email := fmt.Sprintf("%s-%d@tooken.test", emailPrefix, time.Now().UnixNano())
	if err := db.QueryRow(`
		INSERT INTO usr.users (full_name, email, password, role)
		VALUES ($1, $2, 'not-a-real-hash', $3)
		RETURNING id`,
		emailPrefix, email, role,
	).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}

func insertApprovedKYC(t *testing.T, db *sql.DB, userID, decidedBy int, fullName, country string) int {
	t.Helper()

	var id int
	if err := db.QueryRow(`
		INSERT INTO usr.kyc_verification (
			user_id,
			status,
			declared_full_name,
			declared_country_code,
			submitted_at,
			decided_at,
			decided_by,
			expires_at
		)
		VALUES ($1, 'approved', $2, $3, NOW(), NOW(), $4, NOW() + INTERVAL '30 days')
		RETURNING id`,
		userID, fullName, country, decidedBy,
	).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}

func assertUserVerified(t *testing.T, db *sql.DB, userID int) {
	t.Helper()

	var status string
	var verifiedAt sql.NullTime
	if err := db.QueryRow(`
		SELECT kyc_status, kyc_verified_at
		FROM usr.users
		WHERE id = $1`, userID).Scan(&status, &verifiedAt); err != nil {
		t.Fatal(err)
	}
	if status != "verified" {
		t.Fatalf("usr.users.kyc_status = %q, want verified", status)
	}
	if !verifiedAt.Valid {
		t.Fatal("usr.users.kyc_verified_at is NULL")
	}
}

func countRows(t *testing.T, db *sql.DB, query string, args ...any) int {
	t.Helper()

	var count int
	if err := db.QueryRow(query, args...).Scan(&count); err != nil {
		t.Fatal(err)
	}
	return count
}

func assertBalanceGain(t *testing.T, before, after *big.Int, humanAmount int64) {
	t.Helper()

	expected := new(big.Int).Mul(big.NewInt(humanAmount), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	if gained := new(big.Int).Sub(after, before); gained.Cmp(expected) != 0 {
		t.Fatalf("balance gain = %s, want %s", gained, expected)
	}
}

func bindCall(ctx context.Context) *bind.CallOpts {
	return &bind.CallOpts{Context: ctx}
}

// mintShares mints whole shares and waits for the chain, the way the delivery
// does: sign, record, send, then read the receipt. It stands in for the
// Service.Mint that was removed, so these tests keep asking the same question
// through the path the product actually takes.
func mintShares(ctx context.Context, svc *blkServices.Service, tokenAddr, wallet string, shares int) (string, error) {

	txHash, sent, err := svc.SendMintShares(ctx, tokenAddr, wallet, shares, func(string) error { return nil })
	if err != nil {
		return txHash, err
	}
	if !sent {
		return txHash, fmt.Errorf("mint of %d shares to %s never reached the network", shares, wallet)
	}

	found, succeeded, err := svc.MintReceipt(ctx, txHash, wallet, shares)
	if err != nil {
		return txHash, err
	}
	if !found {
		return txHash, fmt.Errorf("mint %s was not mined in time", txHash)
	}
	if !succeeded {
		return txHash, fmt.Errorf("mint %s was mined but did not create the shares", txHash)
	}

	return txHash, nil
}

// integrationUnavailable skips the test when no integration run was asked for,
// and fails it when one was: with TEST_DATABASE_URL or ETH_TEST_RPC set, a
// missing database or node must not let the suite pass without running.
func integrationUnavailable(t *testing.T, format string, args ...any) {
	t.Helper()
	if os.Getenv("TEST_DATABASE_URL") != "" || os.Getenv("ETH_TEST_RPC") != "" {
		t.Fatalf(format, args...)
	}
	t.Skipf(format, args...)
}
