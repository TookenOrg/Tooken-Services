-- Read-only preview of purge_restart.sql, including blockchain database records.
-- No account details are exposed. Actual chain state is not read or changed.
SELECT current_database() AS database_name,
    current_user AS database_user;
SELECT 'real_estates' AS scope,
    count(*) AS rows
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
SELECT 'favorites',
    count(*)
FROM rel.user_asset_favorite
UNION ALL
SELECT 'user_2',
    count(*)
FROM usr.users
UNION ALL
SELECT 'user_2_wallets',
    count(*)
FROM blk.user_wallet
UNION ALL
SELECT 'user_2_identities',
    count(*)
FROM blk.identity
UNION ALL
SELECT 'all_identities',
    count(*)
FROM blk.identity
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
SELECT 'kyc_onchain_statuses_to_reset',
    count(*)
FROM usr.users
WHERE kyc_status = 'verified'
UNION ALL
SELECT 'related_kyc',
    count(*)
FROM usr.kyc_verification
UNION ALL
SELECT 'other_users_kyc_decided_by_user_2',
    count(*)
FROM usr.kyc_verification;