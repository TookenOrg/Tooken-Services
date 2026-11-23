-- Destructive maintenance: back up the database and stop the API first,
-- including its background jobs. Run preview_restart.sql before this file.
-- Deletes ALL user accounts (including administrators), user wallets, KYC
-- verifications, favorites, assets, orders, payments, fees and order history.
-- Clears ALL blockchain deployment records, identities, tokens and transaction
-- logs. Issuers, chain definitions and reference status/type catalogs are
-- preserved.
-- This resets the DATABASE ONLY: old blockchain contracts, claims and balances
-- are NOT deleted or revoked. After committing:
--   1. Rotate JWT_SECRET so every pre-reset access or refresh token stops
--      working immediately.
--   2. Sign up a new user and grant it the initial ADMIN role via SQL (see
--      README.md).
--   3. Run tools/db/deploy_fresh_platform.sh against a fresh local node or a
--      new platform deployment on a public network before reopening traffic.
-- Keep sequence values: reusing user IDs can revive old JWTs if JWT_SECRET was
-- not rotated, and reusing asset IDs can collide with deterministic token salts.
-- Default is a dry run. Replace ONLY the final ROLLBACK with COMMIT to apply.
BEGIN;
SET LOCAL lock_timeout = '10s';
SET LOCAL statement_timeout = '2min';
LOCK TABLE usr.users,
usr.kyc_verification,
blk.user_wallet,
blk.identity,
ass.real_estate,
ass.real_estate_address,
ass.real_estate_media,
ass.real_estate_shares_config,
ass.real_estate_specification,
rel.user_asset_favorite,
iss.issuance_orders,
iss.issuance_order_fees,
iss.issuance_order_payments,
iss.issuance_order_status_history,
blk.token,
blk.token_transaction,
blk.contract_instance,
blk.contract_implementation,
blk.contract_role,
blk.eth_transaction IN ACCESS EXCLUSIVE MODE;
LOCK TABLE ass.issuer,
blk.chains IN ACCESS SHARE MODE;
CREATE TEMP TABLE _purge_expected ON COMMIT DROP AS
SELECT (
        SELECT count(*)
        FROM ass.issuer
    ) AS issuers,
    (
        SELECT count(*)
        FROM blk.chains
    ) AS chains;
-- Explicit tables avoid TRUNCATE CASCADE silently including an unrelated table,
-- while RESTRICT aborts if a new table still references any purged table.
-- TRUNCATE also clears the append-only order history and KYC rows cleanly
-- without firing row-level UPDATE/DELETE triggers.
TRUNCATE TABLE iss.issuance_order_status_history,
iss.issuance_order_payments,
iss.issuance_order_fees,
iss.issuance_orders,
rel.user_asset_favorite,
ass.real_estate_address,
ass.real_estate_media,
ass.real_estate_shares_config,
ass.real_estate_specification,
ass.real_estate,
blk.identity,
blk.user_wallet,
usr.kyc_verification,
usr.users,
blk.token_transaction,
blk.token,
blk.contract_instance,
blk.contract_implementation,
blk.contract_role,
blk.eth_transaction CONTINUE IDENTITY RESTRICT;
DO $$
DECLARE
    expected RECORD;
    actual_issuers bigint;
    actual_chains bigint;
BEGIN
    SELECT * INTO STRICT expected FROM _purge_expected;
    IF EXISTS (SELECT 1 FROM usr.users) THEN
        RAISE EXCEPTION 'usr.users is not empty after purge';
    END IF;
    IF EXISTS (SELECT 1 FROM usr.kyc_verification) THEN
        RAISE EXCEPTION 'usr.kyc_verification is not empty after purge';
    END IF;
    IF EXISTS (SELECT 1 FROM blk.user_wallet) THEN
        RAISE EXCEPTION 'blk.user_wallet is not empty after purge';
    END IF;
    SELECT count(*) INTO actual_issuers FROM ass.issuer;
    IF actual_issuers <> expected.issuers THEN
        RAISE EXCEPTION 'Preserved issuer count changed from % to %', expected.issuers, actual_issuers;
    END IF;
    SELECT count(*) INTO actual_chains FROM blk.chains;
    IF actual_chains <> expected.chains THEN
        RAISE EXCEPTION 'Preserved chain count changed from % to %', expected.chains, actual_chains;
    END IF;
END $$;
SELECT 'remaining_users' AS scope,
    count(*) AS rows
FROM usr.users
UNION ALL
SELECT 'remaining_user_wallets',
    count(*)
FROM blk.user_wallet
UNION ALL
SELECT 'remaining_kyc_verifications',
    count(*)
FROM usr.kyc_verification
UNION ALL
SELECT 'remaining_identities',
    count(*)
FROM blk.identity
UNION ALL
SELECT 'remaining_favorites',
    count(*)
FROM rel.user_asset_favorite
UNION ALL
SELECT 'remaining_real_estates',
    count(*)
FROM ass.real_estate
UNION ALL
SELECT 'remaining_orders',
    count(*)
FROM iss.issuance_orders
UNION ALL
SELECT 'remaining_tokens',
    count(*)
FROM blk.token
UNION ALL
SELECT 'remaining_contract_roles',
    count(*)
FROM blk.contract_role
UNION ALL
SELECT 'remaining_contract_implementations',
    count(*)
FROM blk.contract_implementation
UNION ALL
SELECT 'remaining_contract_instances',
    count(*)
FROM blk.contract_instance
UNION ALL
SELECT 'remaining_chain_transaction_logs',
    count(*)
FROM blk.eth_transaction
UNION ALL
SELECT 'preserved_issuers',
    count(*)
FROM ass.issuer
UNION ALL
SELECT 'preserved_chains',
    count(*)
FROM blk.chains;
ROLLBACK;
-- Replace ROLLBACK above with COMMIT to apply the purge.