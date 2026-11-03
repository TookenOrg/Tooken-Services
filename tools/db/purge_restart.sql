-- Destructive maintenance: back up the database and stop the API first,
-- including its background jobs. Run preview_restart.sql before this file.
-- Deletes ALL assets, orders, payments, fees, order history and favorites;
-- deletes account 2, its wallets and KYC records involving it.
-- This includes other users' KYC decisions made or revoked by account 2.
-- Clears ALL blockchain deployment records, identities and transaction logs.
-- Other accounts, their wallets, issuers and chain definitions are preserved.
-- Remaining KYC approvals require synchronization with the new platform.
-- This resets the DATABASE ONLY: old blockchain contracts, claims and balances
-- are NOT deleted or revoked. Follow restart.md for a fresh local node or a
-- new platform deployment on a public network before reopening the site.
-- Keep sequence values: reusing user IDs can revive old JWTs, and reusing
-- asset IDs can collide with existing deterministic token deployment salts.
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
DO $$ BEGIN IF EXISTS (
    SELECT 1
    FROM usr.users
    WHERE id = 2
        AND role = 'ADMIN'
)
AND NOT EXISTS (
    SELECT 1
    FROM usr.users
    WHERE id <> 2
        AND role = 'ADMIN'
) THEN RAISE EXCEPTION 'Account 2 is the last administrator; create another administrator before purging';
END IF;
END $$;
-- Explicit tables avoid TRUNCATE CASCADE silently including an unrelated table.
-- TRUNCATE also allows this one-off purge of the append-only order history
-- without disabling its protection against ordinary UPDATE/DELETE operations.
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
blk.token_transaction,
blk.token,
blk.contract_instance,
blk.contract_implementation,
blk.contract_role,
blk.eth_transaction CONTINUE IDENTITY RESTRICT;
-- DELETE, not TRUNCATE, lets the KYC trigger refresh affected user profiles.
DELETE FROM usr.kyc_verification;
DELETE FROM blk.user_wallet;
DELETE FROM usr.users;
-- Approval documents survive for other users, but verification on the old
-- registry cannot authorize investment on a new platform.
UPDATE usr.users
SET kyc_status = 'approved',
    kyc_verified_at = NULL,
    updated_at = now()
WHERE kyc_status = 'verified';
UPDATE usr.users
SET kyc_verified_at = NULL,
    updated_at = now()
WHERE kyc_verified_at IS NOT NULL;
SELECT 'remaining_real_estates' AS scope,
    count(*) AS rows
FROM ass.real_estate
UNION ALL
SELECT 'remaining_orders',
    count(*)
FROM iss.issuance_orders
UNION ALL
SELECT 'remaining_user_2',
    count(*)
FROM usr.users
UNION ALL
SELECT 'remaining_tokens',
    count(*)
FROM blk.token
UNION ALL
SELECT 'remaining_identities',
    count(*)
FROM blk.identity
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
SELECT 'remaining_onchain_verified_users',
    count(*)
FROM usr.users
WHERE kyc_status = 'verified'
UNION ALL
SELECT 'preserved_other_users',
    count(*)
FROM usr.users;
ROLLBACK;
# Change to COMMIT to commit the script