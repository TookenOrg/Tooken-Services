-- Rollback 000019 — the token name becomes a label again.
--
-- READ THIS BEFORE REPLAYING IT ANYWHERE THAT MATTERS
--   In every environment that already held token_name_unique, 000019 did not
--   create it: it was added by hand, long before this migration history
--   started. Rolling 000019 back there does not restore an earlier state — it
--   removes a protection that predates the migration entirely.
--
--   What goes with it is not cosmetic. token_name is the idempotency key of
--   CreateToken and the lookup key of GetTokenByName. Without the constraint,
--   a replayed deployment can leave two tokens sharing a name, and the next
--   lookup returns one of them arbitrarily — binding an asset to a contract
--   address that is not its own.
--
--   The rollback is written all the same. A migration that cannot be undone is
--   worse than one that states plainly what undoing it costs.

ALTER TABLE blk.token
    DROP CONSTRAINT IF EXISTS token_name_unique;

COMMENT ON COLUMN blk.token.token_name IS NULL;
