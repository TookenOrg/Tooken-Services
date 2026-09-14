-- Rollback 000023 — the place to store a key comes back, the keys do not.
--
--   This restores one column, private_key_encrypted, and it comes back empty.
--   The keys dropped by 000023 are gone from the database and cannot be
--   recomputed from anything: a private key is not derivable from the address
--   it produced. Rolling back therefore gives back the shape, never the
--   contents.
--
--   Only one column is recreated, not the three the history went through.
--   private_key_clear is deliberately not restored: storing a key in clear
--   beside its encrypted copy was the original defect, and a rollback should
--   not reopen it. A rollback restores a working schema, not every mistake
--   that was made on the way to the current one.
--
-- WHAT COMES BACK WITH IT
--   An empty column to write a key into. That is the whole risk 000023 was
--   removing — code that once wrote there would start writing there again.
--   Before running this, be sure something actually needs to store a key, and
--   that it is not simply a way to avoid recoveryAddress.
--
-- ⚠️  IF WALLETS WERE CREATED AFTER 000023
--   They carry no key and never will. The column is nullable, so they stay
--   valid; but any code assuming "a custodial wallet has a key" would be wrong
--   about them. The absence is permanent, not a gap to backfill.

ALTER TABLE blk.user_wallet
    ADD COLUMN IF NOT EXISTS private_key_encrypted bytea;

COMMENT ON TABLE blk.user_wallet IS NULL;
COMMENT ON COLUMN blk.user_wallet.wallet_address IS NULL;
