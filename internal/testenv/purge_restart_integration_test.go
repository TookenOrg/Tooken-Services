//go:build integration

package testenv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPurgeAndPreviewRestartScripts(t *testing.T) {
	db := ScratchDatabase(t, "tooken_it_purge_restart")

	baselinePath, err := BaselinePath()
	if err != nil {
		t.Fatal(err)
	}
	dbDir := filepath.Dir(baselinePath)

	previewBytes, err := os.ReadFile(filepath.Join(dbDir, "preview_restart.sql"))
	if err != nil {
		t.Fatalf("read preview_restart.sql: %v", err)
	}
	purgeBytes, err := os.ReadFile(filepath.Join(dbDir, "purge_restart.sql"))
	if err != nil {
		t.Fatalf("read purge_restart.sql: %v", err)
	}
	previewSQL := string(previewBytes)
	purgeDryRunSQL := string(purgeBytes)
	if !strings.Contains(purgeDryRunSQL, "\nROLLBACK;\n") {
		t.Fatal("purge_restart.sql must default to ROLLBACK")
	}
	purgeCommitSQL := strings.Replace(purgeDryRunSQL, "\nROLLBACK;\n", "\nCOMMIT;\n", 1)

	seedSQL := `
		INSERT INTO ass.real_estate_type (id, name) VALUES (1, 'Apartment') ON CONFLICT (id) DO NOTHING;
		INSERT INTO blk.chains (id, name, network, is_evm, native_symbol)
		VALUES (31337, 'Hardhat', 'local', true, 'ETH');

		INSERT INTO ass.issuer (name, legal_form, registration_number, country_code, status_id)
		VALUES ('Preserved SPV', 'SAS', '555666777', 'FR', 2);

		INSERT INTO usr.users (full_name, email, password, role) VALUES
			('Root Admin', 'admin@tooken.test', 'hash', 'ADMIN'),
			('Asset Manager', 'manager@tooken.test', 'hash', 'MANAGER'),
			('Retail Investor', 'investor@tooken.test', 'hash', 'USER');

		INSERT INTO blk.user_wallet (user_id, wallet_address) VALUES
			(1, '0x1111111111111111111111111111111111111111'),
			(3, '0x3333333333333333333333333333333333333333');

		INSERT INTO blk.identity (user_id, wallet_id, address, tx_hash)
		VALUES (3, 2, '0x4444444444444444444444444444444444444444', '0x' || repeat('a', 64));

		INSERT INTO usr.kyc_verification (
			user_id, status, declared_full_name, declared_country_code,
			submitted_at, decided_at, decided_by, expires_at
		) VALUES (
			3, 'approved', 'Retail Investor', 'FR',
			now() - interval '2 minutes', now() - interval '1 minute', 1, now() + interval '1 year'
		);
		UPDATE usr.users
		SET kyc_status = 'verified', kyc_verified_at = now(), kyc_expires_at = now() + interval '1 year', country_code = 'FR'
		WHERE id = 3;

		INSERT INTO blk.contract_implementation (tx_hash, address, contract_name, version, chain_id)
		VALUES ('0x' || repeat('b', 64), '0x5555555555555555555555555555555555555555', 'TokenImpl', 1, 31337);

		INSERT INTO blk.contract_role (tx_hash, address, contract_name)
		VALUES ('0x' || repeat('c', 64), '0x6666666666666666666666666666666666666666', 'TREXFactory');

		INSERT INTO blk.contract_instance (contract_implementation_id, parent_contract_id, address, contract_name, chain_id)
		VALUES (1, 1, '0x7777777777777777777777777777777777777777', 'TokenProxy', 31337);

		INSERT INTO blk.eth_transaction (tx_hash, tx_name, block_number, to_address, value_wei)
		VALUES ('0x' || repeat('d', 64), 'deployTREXSuite', 1, '0x6666666666666666666666666666666666666666', 0);

		INSERT INTO blk.token (symbol, token_name, address, nb_decimal, salt)
		VALUES ('TKN1', 'Tooken Asset 1', '0x8888888888888888888888888888888888888888', 0, 're-1');

		INSERT INTO blk.token_transaction (token_id, tx_hash, tx_name)
		VALUES (1, '0x' || repeat('e', 64), 'mint');

		INSERT INTO ass.real_estate (title, estate_type, issuer_id, status_id, active, token_id)
		VALUES ('Villa Purge Probe', 1, 1, 4, true, 1);

		INSERT INTO ass.real_estate_address (real_estate_id, street, city, postal_code, country_code)
		VALUES (1, '10 Rue de la Paix', 'Paris', '75002', 'FR');

		INSERT INTO ass.real_estate_specification (real_estate_id, surface_area)
		VALUES (1, 120.50);

		INSERT INTO ass.real_estate_shares_config (real_estate_id, total_shares, price_per_share, currency_code, entry_fee_rate)
		VALUES (1, 1000, 100.00, 'EUR', 2.0000);

		INSERT INTO ass.real_estate_media (real_estate_id, url, media_type, "position")
		VALUES (1, 'https://example.com/villa.jpg', 'image', 0);

		INSERT INTO rel.user_asset_favorite (user_id, real_estate_id)
		VALUES (3, 1);

		INSERT INTO iss.issuance_orders (
			user_id, asset_id, quantity, status_id, order_reference,
			unit_price, currency_code, gross_amount, fee_amount, amount_due
		) VALUES (
			3, 1, 10, 4, 'ORD-PURGE-001',
			100.00, 'EUR', 1000.00, 20.00, 1020.00
		);

		INSERT INTO iss.issuance_order_fees (order_id, fee_code, rate, base_amount, amount)
		VALUES (1, 'ENTRY', 2.0000, 1000.00, 20.00);

		INSERT INTO iss.issuance_order_payments (order_id, kind, provider, amount, currency_code, recorded_by_user_id)
		VALUES (1, 'PAYMENT', 'MANUAL', 1020.00, 'EUR', 1);

		INSERT INTO iss.issuance_order_status_history (order_id, from_status_id, to_status_id, actor_user_id, reason)
		VALUES (1, 2, 4, 1, 'payment confirmed');
	`
	if _, err := db.Exec(seedSQL); err != nil {
		t.Fatalf("seed database: %v", err)
	}

	count := func(table string) int {
		t.Helper()
		var n int
		if err := db.QueryRow(`SELECT count(*) FROM ` + table).Scan(&n); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		return n
	}

	if _, err := db.Exec(previewSQL); err != nil {
		t.Fatalf("preview_restart.sql failed: %v", err)
	}
	if got := count("usr.users"); got != 3 {
		t.Fatalf("preview changed usr.users count to %d, want 3", got)
	}

	if _, err := db.Exec(purgeDryRunSQL); err != nil {
		t.Fatalf("purge_restart.sql dry-run failed: %v", err)
	}
	if got := count("usr.users"); got != 3 {
		t.Fatalf("dry-run changed usr.users count to %d, want 3", got)
	}

	if _, err := db.Exec(purgeCommitSQL); err != nil {
		t.Fatalf("purge_restart.sql commit failed: %v", err)
	}

	purgedTables := []string{
		"usr.users",
		"usr.kyc_verification",
		"blk.user_wallet",
		"blk.identity",
		"ass.real_estate",
		"ass.real_estate_address",
		"ass.real_estate_media",
		"ass.real_estate_shares_config",
		"ass.real_estate_specification",
		"rel.user_asset_favorite",
		"iss.issuance_orders",
		"iss.issuance_order_fees",
		"iss.issuance_order_payments",
		"iss.issuance_order_status_history",
		"blk.token",
		"blk.token_transaction",
		"blk.contract_instance",
		"blk.contract_implementation",
		"blk.contract_role",
		"blk.eth_transaction",
	}
	for _, table := range purgedTables {
		if got := count(table); got != 0 {
			t.Errorf("after purge %s has %d row(s), want 0", table, got)
		}
	}

	if got := count("ass.issuer"); got != 1 {
		t.Errorf("after purge ass.issuer has %d row(s), want 1", got)
	}
	if got := count("blk.chains"); got != 1 {
		t.Errorf("after purge blk.chains has %d row(s), want 1", got)
	}
	if got := count("ass.real_estate_status"); got == 0 {
		t.Error("reference table ass.real_estate_status was emptied")
	}

	var nextUserID int
	if err := db.QueryRow(`
		INSERT INTO usr.users (full_name, email, password, role)
		VALUES ('Fresh Admin', 'fresh-admin@tooken.test', 'hash', 'ADMIN')
		RETURNING id
	`).Scan(&nextUserID); err != nil {
		t.Fatalf("insert fresh user after purge: %v", err)
	}
	if nextUserID <= 3 {
		t.Errorf("sequence was reset: new user got id %d, want > 3", nextUserID)
	}
}
