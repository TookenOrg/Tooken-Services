-- Rollback 000016 — naming an issuer becomes a convention again.
--
-- The rows written while the constraint was in force keep their issuer: this
-- only lifts the obligation for the next ones.

ALTER TABLE ass.real_estate
    ALTER COLUMN issuer_id DROP NOT NULL;

COMMENT ON COLUMN ass.real_estate.issuer_id IS NULL;
