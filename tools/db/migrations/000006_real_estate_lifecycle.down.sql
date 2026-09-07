-- Rollback 000006

DROP TRIGGER IF EXISTS real_estate_sync_active_trg ON ass.real_estate;
DROP FUNCTION IF EXISTS ass.real_estate_sync_active();

DROP INDEX IF EXISTS ass.real_estate_status_id_idx;

ALTER TABLE ass.real_estate
    DROP COLUMN IF EXISTS status_id,
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS published_at;

DROP TABLE IF EXISTS ass.real_estate_status;
