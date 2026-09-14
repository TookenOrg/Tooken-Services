-- Rollback 000020 — the replay deadlock comes back.
--
--   Dropping the salt takes away the only key that can tie an orphaned
--   blk.token row back to its asset. After this, an asset whose deployment was
--   interrupted before it was bound becomes unpublishable again: the code
--   redeploys, token.go answers "Token already exists", and only a manual
--   UPDATE gets the asset out.
--
--   The recorded salts are lost, not recomputed. Re-running 000020 rebuilds
--   them for every token bound to an asset, so nothing is permanently
--   destroyed — but any salt that was filled in by hand for an unlinked token
--   is gone for good. Export them first if that applies:
--
--     SELECT id, token_name, salt FROM blk.token WHERE salt IS NOT NULL;

DROP INDEX IF EXISTS blk.token_salt_uk;

ALTER TABLE blk.token
    DROP COLUMN IF EXISTS salt;
