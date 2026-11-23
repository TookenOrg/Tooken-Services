-- Read-only preview of purge_restart.sql: counts every user, wallet, KYC,
-- asset, order and blockchain record that will be wiped, alongside preserved
-- issuers and chain definitions. No personal account details are exposed.
-- Actual chain state is not read or changed.
SELECT current_database() AS database_name,
    current_user AS database_user;
SELECT 'users' AS scope,
    count(*) AS rows
FROM usr.users
UNION ALL
SELECT 'admin_users',
    count(*)
FROM usr.users
WHERE role = 'ADMIN'
UNION ALL
SELECT 'user_wallets',
    count(*)
FROM blk.user_wallet
UNION ALL
SELECT 'kyc_verifications',
    count(*)
FROM usr.kyc_verification
UNION ALL
SELECT 'identities',
    count(*)
FROM blk.identity
UNION ALL
SELECT 'favorites',
    count(*)
FROM rel.user_asset_favorite
UNION ALL
SELECT 'real_estates',
    count(*)
FROM ass.real_estate
UNION ALL
SELECT 'orders',
    count(*)
FROM iss.issuance_orders
UNION ALL
SELECT 'order_payments',
    count(*)
FROM iss.issuance_order_payments
UNION ALL
SELECT 'order_fees',
    count(*)
FROM iss.issuance_order_fees
UNION ALL
SELECT 'order_history',
    count(*)
FROM iss.issuance_order_status_history
UNION ALL
SELECT 'tokens',
    count(*)
FROM blk.token
UNION ALL
SELECT 'token_transactions',
    count(*)
FROM blk.token_transaction
UNION ALL
SELECT 'contract_implementations',
    count(*)
FROM blk.contract_implementation
UNION ALL
SELECT 'contract_instances',
    count(*)
FROM blk.contract_instance
UNION ALL
SELECT 'contract_roles',
    count(*)
FROM blk.contract_role
UNION ALL
SELECT 'chain_transaction_logs',
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