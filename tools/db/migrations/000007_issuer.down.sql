-- Rollback 000007

DROP INDEX IF EXISTS ass.real_estate_issuer_id_idx;

ALTER TABLE ass.real_estate
    DROP COLUMN IF EXISTS issuer_id;

DROP TABLE IF EXISTS ass.issuer;
DROP TABLE IF EXISTS ass.issuer_status;
