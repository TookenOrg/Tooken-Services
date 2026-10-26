//go:build integration

// Database environment for the delivery transactions.
//
// These tests need a real PostgreSQL and no chain at all: what they check is
// what the schema enforces — the conditional UPDATE that arbitrates two
// deliveries, the CHECK that refuses a CLOSED order without a hash, the
// append-only timeline. A fake database would only prove that the fake agrees
// with itself.
//
// Each run builds its own database from tools/db/baseline.sql and drops it at
// the end, so these tests never touch the project database. Point
// TEST_DATABASE_URL at a local server whose user may CREATE DATABASE:
//
//	TEST_DATABASE_URL="postgres://postgres:pw@127.0.0.1:55444/tooken?sslmode=disable" \
//	  go test -tags integration ./internal/delivery/...
package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	_ "github.com/lib/pq"
)

const (
	deliveryScratchDatabase = "tooken_it_delivery_database"
	deliveryBaselinePath    = "../../../tools/db/baseline.sql"
)

// deliveryEnv is a private database wired into the global the production code
// reads, plus the few rows a delivery needs to exist at all.
type deliveryEnv struct {
	DB *sql.DB
}

func newDeliveryEnv(t *testing.T) *deliveryEnv {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	logger.Init(true)

	db := createDeliveryScratchDatabase(t, dsn)

	previous := globals.DB
	globals.DB = db
	t.Cleanup(func() {
		globals.DB = previous
		db.Close()
		dropDeliveryScratchDatabase(t, dsn)
	})

	return &deliveryEnv{DB: db}
}

func createDeliveryScratchDatabase(t *testing.T, dsn string) *sql.DB {
	t.Helper()

	adminDSN, scratchDSN, err := deliveryScratchDSNs(dsn)
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

	if _, err := admin.Exec(`DROP DATABASE IF EXISTS ` + deliveryScratchDatabase); err != nil {
		integrationUnavailable(t, "cannot drop %s (the test user needs CREATEDB): %v", deliveryScratchDatabase, err)
	}
	if _, err := admin.Exec(`CREATE DATABASE ` + deliveryScratchDatabase); err != nil {
		integrationUnavailable(t, "cannot create %s (the test user needs CREATEDB): %v", deliveryScratchDatabase, err)
	}

	baseline, err := os.ReadFile(deliveryBaselinePath)
	if err != nil {
		t.Fatalf("read %s: %v", deliveryBaselinePath, err)
	}

	db, err := sql.Open("postgres", scratchDSN)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(string(baseline)); err != nil {
		db.Close()
		t.Fatalf("apply %s: %v", deliveryBaselinePath, err)
	}

	return db
}

func dropDeliveryScratchDatabase(t *testing.T, dsn string) {
	t.Helper()

	adminDSN, _, err := deliveryScratchDSNs(dsn)
	if err != nil {
		return
	}
	admin, err := sql.Open("postgres", adminDSN)
	if err != nil {
		return
	}
	defer admin.Close()
	_, _ = admin.Exec(`DROP DATABASE IF EXISTS ` + deliveryScratchDatabase)
}

func deliveryScratchDSNs(dsn string) (adminDSN, scratchDSN string, err error) {
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
	scratch.Path = "/" + deliveryScratchDatabase

	return admin.String(), scratch.String(), nil
}

// orderFixture describes the order a test wants to start from.
type orderFixture struct {
	statusID     int
	kycStatus    string
	kycExpired   bool
	activeWallet bool
	tokenised    bool
	shares       int
}

func aPaidOrder() orderFixture {
	return orderFixture{
		statusID:     int(server.IssuanceOrderStatusIdPaymentConfirmed),
		kycStatus:    "verified",
		activeWallet: true,
		tokenised:    true,
		shares:       10,
	}
}

type seededOrder struct {
	OrderID      int64
	OrderRef     string
	UserID       int
	AssetID      int
	Wallet       string
	TokenAddress string
}

// seed builds the whole chain of rows a delivery reads: an investor, a wallet,
// a token, an asset pointing at it, and the order itself. Production reads all
// of them, so a test that invents ids proves nothing about the joins.
func (e *deliveryEnv) seed(t *testing.T, f orderFixture) seededOrder {
	t.Helper()

	unique := time.Now().UnixNano()

	var userID int
	// users_kyc_verified_at_ck: a verified user must carry the date it was
	// verified on. The schema refuses the shortcut.
	kycVerifiedAt := "NULL"
	if f.kycStatus == "verified" {
		kycVerifiedAt = "now() - interval '1 hour'"
	}
	kycExpiry := "NULL"
	if f.kycExpired {
		// users_kyc_expiry_ck: the expiry must come after the verification,
		// so an expired KYC is one verified long ago.
		kycVerifiedAt = "now() - interval '2 years'"
		kycExpiry = "now() - interval '1 day'"
	}
	if err := e.DB.QueryRow(fmt.Sprintf(`
		INSERT INTO usr.users (full_name, email, password, kyc_status, kyc_verified_at, kyc_expires_at)
		VALUES ($1, $2, 'not-a-real-hash', $3, %s, %s)
		RETURNING id`, kycVerifiedAt, kycExpiry),
		"Alice Investor",
		fmt.Sprintf("alice-%d@test.local", unique),
		f.kycStatus,
	).Scan(&userID); err != nil {
		t.Fatalf("insert user: %v", err)
	}

	wallet := fmt.Sprintf("0x%040x", unique)
	if f.activeWallet {
		if _, err := e.DB.Exec(`
			INSERT INTO blk.user_wallet (user_id, wallet_address, label, is_active)
			VALUES ($1, $2, 'test', true)`, userID, wallet); err != nil {
			t.Fatalf("insert wallet: %v", err)
		}
	}

	var tokenID *int64
	tokenAddress := fmt.Sprintf("0x%040x", unique+1)
	if f.tokenised {
		var id int64
		if err := e.DB.QueryRow(`
			INSERT INTO blk.token (symbol, token_name, salt, address, nb_decimal)
			VALUES ($1, $2, $3, $4, 18)
			RETURNING id`,
			"TKT", "Tooken Test", fmt.Sprintf("re-%d", unique), tokenAddress,
		).Scan(&id); err != nil {
			t.Fatalf("insert token: %v", err)
		}
		tokenID = &id
	}

	var assetID int
	// ass.real_estate needs an issuer and an estate type; contract_address is
	// filled by a trigger from token_id, which is the path production relies on.
	var issuerID int
	if err := e.DB.QueryRow(`
		INSERT INTO ass.issuer (name, legal_form)
		VALUES ($1, 'SA')
		RETURNING id`, fmt.Sprintf("Tooken Issuer %d", unique)).Scan(&issuerID); err != nil {
		t.Fatalf("insert issuer: %v", err)
	}

	if err := e.DB.QueryRow(`
		INSERT INTO ass.real_estate (title, description, estate_type, issuer_id, token_id)
		VALUES ($1, 'a test asset', 1, $2, $3)
		RETURNING id`,
		fmt.Sprintf("Tooken Villa %d", unique), issuerID, tokenID,
	).Scan(&assetID); err != nil {
		t.Fatalf("insert real estate: %v", err)
	}

	if _, err := e.DB.Exec(`
		INSERT INTO ass.real_estate_shares_config (real_estate_id, total_shares, price_per_share)
		VALUES ($1, 1500, 204.00)`, assetID); err != nil {
		t.Fatalf("insert shares config: %v", err)
	}

	orderRef := fmt.Sprintf("ISS-%d", unique)

	// issuance_orders_delivered_hash_ck: statuses 5, 10, 11 and 12 all claim
	// the shares were minted, so they cannot exist without a hash.
	var seedHash any
	switch f.statusID {
	case int(server.IssuanceOrderStatusIdDelivered), int(server.IssuanceOrderStatusIdClosed):
		seedHash = aTxHash
	}

	var orderID int64
	if err := e.DB.QueryRow(`
		INSERT INTO iss.issuance_orders
			(user_id, asset_id, quantity, status_id, order_reference,
			 unit_price, currency_code, gross_amount, fee_amount, amount_due, delivery_tx_hash)
		VALUES ($1, $2, $3, $4, $5, 204.00, 'EUR', $6, 0.00, $6, $7)
		RETURNING id`,
		userID, assetID, f.shares, f.statusID, orderRef,
		fmt.Sprintf("%.2f", float64(f.shares)*204.00), seedHash,
	).Scan(&orderID); err != nil {
		t.Fatalf("insert order: %v", err)
	}

	return seededOrder{
		OrderID:      orderID,
		OrderRef:     orderRef,
		UserID:       userID,
		AssetID:      assetID,
		Wallet:       wallet,
		TokenAddress: tokenAddress,
	}
}

// status reads back the order status and its delivery hash.
func (e *deliveryEnv) status(t *testing.T, orderRef string) (int, *string) {
	t.Helper()

	var statusID int
	var hash sql.NullString
	if err := e.DB.QueryRow(`
		SELECT status_id, delivery_tx_hash FROM iss.issuance_orders WHERE order_reference = $1`,
		orderRef).Scan(&statusID, &hash); err != nil {
		t.Fatalf("read order %s: %v", orderRef, err)
	}

	if !hash.Valid {
		return statusID, nil
	}
	return statusID, &hash.String
}

// timeline reads the transitions of an order, oldest first.
func (e *deliveryEnv) timeline(t *testing.T, orderID int64) []transition {
	t.Helper()

	rows, err := e.DB.Query(`
		SELECT from_status_id, to_status_id, actor_user_id, reason, tx_hash
		FROM iss.issuance_order_status_history
		WHERE order_id = $1
		ORDER BY created_at, id`, orderID)
	if err != nil {
		t.Fatalf("read timeline of order %d: %v", orderID, err)
	}
	defer rows.Close()

	var out []transition
	for rows.Next() {
		var tr transition
		var from sql.NullInt64
		var actor sql.NullInt64
		var reason sql.NullString
		var hash sql.NullString
		if err := rows.Scan(&from, &tr.To, &actor, &reason, &hash); err != nil {
			t.Fatal(err)
		}
		tr.From = int(from.Int64)
		tr.ActorIsNull = !actor.Valid
		tr.Reason = reason.String
		if hash.Valid {
			tr.TxHash = &hash.String
		}
		out = append(out, tr)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}

	return out
}

type transition struct {
	From        int
	To          int
	ActorIsNull bool
	Reason      string
	TxHash      *string
}

const aTxHash = "0xab12cd34ef56ab12cd34ef56ab12cd34ef56ab12cd34ef56ab12cd34ef56ab12"

// nthTxHash gives distinct, schema-valid hashes: 0x and 64 hex digits
// (issuance_orders_delivery_tx_hash_ck).
func nthTxHash(n int) string {
	return fmt.Sprintf("0x%064x", n+1)
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
