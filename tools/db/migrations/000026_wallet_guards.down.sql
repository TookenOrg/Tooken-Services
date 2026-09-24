-- Rollback 000026 — the guards come off, in the order they were put on
--
--   Everything here is dropped, nothing is lost: no data was created by 000026,
--   only rules about what may be written. Rolling back cannot fail on existing
--   rows.
--
--   The custody column is the one exception — it holds a value per wallet, and
--   dropping it forgets which addresses the investors brought themselves.
--   Since that is exactly what decides whether a proof of control is required
--   before minting, a rollback silently makes every external wallet look
--   custodial. Dump the column before running this if any external wallet
--   exists:
--
--     SELECT id, user_id, wallet_address FROM blk.user_wallet
--     WHERE custody = 'external';
--
-- WHAT COMES BACK WITH IT
--   An ONCHAINID may again be deployed for a user who does not exist, and an
--   investor may again hold two active wallets. Both were possible until
--   2026-09-24, and both are the kind of defect that only shows up on-chain,
--   where nothing can be corrected.

DROP INDEX IF EXISTS blk.user_wallet_one_active_per_user_idx;

ALTER TABLE blk.user_wallet DROP CONSTRAINT IF EXISTS user_wallet_custody_ck;
ALTER TABLE blk.user_wallet DROP COLUMN IF EXISTS custody;

ALTER TABLE blk.identity DROP CONSTRAINT IF EXISTS identity_user_id_fkey;
ALTER TABLE blk.user_wallet DROP CONSTRAINT IF EXISTS user_wallet_user_id_fkey;
