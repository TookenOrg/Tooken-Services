-- Rollback 000021 — two sources of truth again, and no invariant.
--
--   This restores the 000006 shape exactly: the status ids copied into the
--   trigger, the reverse direction that lets a write to active move status_id,
--   'closed' marked public again while the trigger keeps it inactive, and no
--   constraint tying active to token_id.
--
--   What that costs, plainly: an asset can go back on sale with no token
--   behind it, and nothing in the database will object. The protection removed
--   here is the one the product rests on.
--
-- WHY NO ROW HAS TO BE REWRITTEN
--   The two rules project identically. Once 000021 set is_public = FALSE for
--   'closed', reading is_public and testing status_id IN (3, 4, 5) give the
--   same answer for all seven statuses — that correction is precisely what
--   made them agree. So the stored active values are already correct under the
--   restored rule, and rolling back touches no asset.
--
--   The constraint is dropped first all the same: it must not outlive the
--   guarantee the new trigger provided.

ALTER TABLE ass.real_estate
    DROP CONSTRAINT IF EXISTS real_estate_active_requires_token_ck;

CREATE OR REPLACE FUNCTION ass.real_estate_sync_active()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' OR NEW.status_id IS DISTINCT FROM OLD.status_id THEN
        NEW.active := (NEW.status_id IN (3, 4, 5));
    ELSIF NEW.active IS DISTINCT FROM OLD.active THEN
        NEW.status_id := CASE WHEN NEW.active THEN 3 ELSE 1 END;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS real_estate_sync_active_trg ON ass.real_estate;
CREATE TRIGGER real_estate_sync_active_trg
    BEFORE INSERT OR UPDATE ON ass.real_estate
    FOR EACH ROW EXECUTE FUNCTION ass.real_estate_sync_active();

-- Put the disagreement back: 'closed' is public in the reference table and
-- inactive in the trigger, exactly as 000006 left it.
UPDATE ass.real_estate_status
SET is_public = TRUE
WHERE id = 6;

COMMENT ON COLUMN ass.real_estate.active IS NULL;
COMMENT ON COLUMN ass.real_estate_status.is_public IS NULL;
