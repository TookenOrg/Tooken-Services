-- Reverts 000011. Cancelled assets keep their status: only the deletion
-- timestamp and its guards disappear.

DROP INDEX IF EXISTS ass.real_estate_not_deleted_idx;

ALTER TABLE ass.real_estate
    DROP CONSTRAINT IF EXISTS real_estate_deleted_not_tokenized_ck,
    DROP CONSTRAINT IF EXISTS real_estate_deleted_status_ck;

ALTER TABLE ass.real_estate
    DROP COLUMN IF EXISTS deleted_at;
