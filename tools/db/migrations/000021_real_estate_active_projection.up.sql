-- 000021 — active becomes a projection, and the invariant enters the schema
--
-- ═══════════════════════════════════════════════════════════════════════════
-- PART ONE — the reference table and the trigger stopped agreeing
-- ═══════════════════════════════════════════════════════════════════════════
--
--   000006 created ass.real_estate_status with an is_public column, then wrote
--   the same rule a second time inside the trigger:
--
--     NEW.active := (NEW.status_id IN (3, 4, 5));
--
--   The two disagree on exactly one status. The reference table marks 'closed'
--   (6) as is_public = TRUE; the trigger makes it active = FALSE. The gap is
--   already documented and arbitrated in issuer_db.go:155-159, in favour of
--   re.active, because that is what the listing and the detail endpoints
--   filter on. So the behaviour is correct — it is the reference table that
--   says something the system does not do.
--
--   Two sources for one rule will eventually be read by two different people.
--   And the repository forbids this exact duplication elsewhere, in as many
--   words (real_estate_columns_test.go:73): what counts as a reserved share
--   must come from the reference table, never from a list of ids copied into
--   the code. The trigger is the only place that breaks its own rule.
--
--   So: the reference table is corrected to describe reality, and the trigger
--   starts reading it. Nothing changes at runtime — no Go code reads
--   ass.real_estate_status today.
--
-- WHY THE REVERSE DIRECTION IS REMOVED
--   The trigger also translated the other way: writing active = TRUE forced
--   status_id to 3. That branch is dead. Nothing writes active — every use in
--   Go is a read (WHERE re.active = true in real_estate_db.go:64 and :99), the
--   API exposes Active only on response models, and the single write in the
--   history is the one-shot backfill of 000006.
--
--   Keeping a second direction only kept a second possible truth. From here,
--   status_id is the input and active is the output; writing active has no
--   effect, because the trigger recomputes it.
--
-- WHY IT RECOMPUTES ON EVERY WRITE
--   The old guard only fired when status_id actually changed. That made the
--   column impossible to repair: UPDATE ... SET status_id = status_id changed
--   nothing, so a row whose active had drifted stayed wrong forever.
--   Recomputing unconditionally makes any write to a row re-derive its
--   projection, which is what the re-synchronisation below relies on.
--
--   The cost is one indexed lookup on a seven-row table per written row.
--
-- ⚠️  THE TRIGGER ONLY FIRES ON A WRITE
--   Changing is_public in the reference table later will NOT update assets
--   already stored — same class of problem as an expiry date that passes with
--   nobody writing. Any future migration that touches
--   ass.real_estate_status.is_public must be followed by the re-synchronisation
--   statement below, in that same migration.
--
-- ═══════════════════════════════════════════════════════════════════════════
-- PART TWO — the invariant that founds the platform was not in the schema
-- ═══════════════════════════════════════════════════════════════════════════
--
--   "No asset can be sold without a token" is the rule the whole product rests
--   on. It is stated in the design, implemented in Go, covered by integration
--   tests — and absent from the database. Before this migration, ass.real_estate
--   carried exactly two CHECK constraints, neither about active or token_id:
--
--     real_estate_deleted_not_tokenized_ck
--     real_estate_deleted_status_ck
--
--   A manual UPDATE, a future endpoint or a service bug therefore puts an
--   asset back on sale with no token behind it, silently. A rule of this
--   weight belongs where it cannot be bypassed.
--
-- WHY NOT VALID, AND WHY IT IS NOT A SOFT OPTION
--   NOT VALID is routinely misread as "disabled". It is not: the constraint is
--   enforced on every INSERT and every UPDATE from the moment it exists. It
--   only skips re-reading the rows already stored.
--
--   That is exactly what is needed here. Assets published before the
--   tokenising publication existed are active with token_id IS NULL. They are
--   known and deliberate; refusing to delete them to force a migration through
--   is the honest call. The future is locked immediately, and the past is
--   settled separately.
--
-- ⚠️  OPERATIONAL CONSEQUENCE — READ BEFORE PLAYING THIS
--   Because NOT VALID still checks UPDATEs, those legacy rows become
--   read-only: any edit to an active asset that carries no token will now be
--   rejected by the constraint. A manager editing such an asset gets an error.
--
--   This is not a side effect to discover in production. List them first:
--
--     SELECT id, title FROM ass.real_estate
--     WHERE active AND token_id IS NULL AND deleted_at IS NULL;
--
--   and clear them promptly — republish each one, which deploys the missing
--   token, or unpublish it. Once the query returns nothing:
--
--     ALTER TABLE ass.real_estate
--         VALIDATE CONSTRAINT real_estate_active_requires_token_ck;
--
--   and the invariant holds everywhere, past included.
--
-- ═══════════════════════════════════════════════════════════════════════════
-- ORDER MATTERS
-- ═══════════════════════════════════════════════════════════════════════════
--   The statements below are not interchangeable. The re-synchronisation
--   writes rows, so it must run BEFORE the constraint exists — otherwise it
--   would touch the legacy rows and be rejected by the very constraint this
--   migration adds, and the whole migration would roll back.
--
--   1. correct the reference table
--   2. rewrite the function to read it
--   3. re-synchronise the rows that drifted
--   4. only then, add the constraint

-- 1. The reference table starts describing what the system actually does.
UPDATE ass.real_estate_status
SET is_public = FALSE
WHERE id = 6;

-- 2. One direction, one source. status_id decides, active follows.
CREATE OR REPLACE FUNCTION ass.real_estate_sync_active()
RETURNS TRIGGER AS $$
DECLARE
    projected BOOLEAN;
BEGIN
    SELECT s.is_public INTO projected
    FROM ass.real_estate_status s
    WHERE s.id = NEW.status_id;

    -- A foreign key already rejects an unknown status, but it is checked after
    -- this trigger has run. Raising here turns a constraint violation into a
    -- sentence, and keeps the projection honest if that key is ever dropped.
    IF NOT FOUND THEN
        RAISE EXCEPTION 'unknown ass.real_estate_status id: %', NEW.status_id;
    END IF;

    NEW.active := projected;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- The trigger itself is unchanged; recreated so that a database built from
-- this history is identical whichever migration last touched it.
DROP TRIGGER IF EXISTS real_estate_sync_active_trg ON ass.real_estate;
CREATE TRIGGER real_estate_sync_active_trg
    BEFORE INSERT OR UPDATE ON ass.real_estate
    FOR EACH ROW EXECUTE FUNCTION ass.real_estate_sync_active();

-- 3. Repair whatever drifted, and only that. Writing status_id onto itself is
--    enough: the function now recomputes on every write.
UPDATE ass.real_estate re
SET status_id = re.status_id
FROM ass.real_estate_status s
WHERE s.id = re.status_id
  AND re.active IS DISTINCT FROM s.is_public;

-- 4. The rule the platform rests on, where it cannot be bypassed.
ALTER TABLE ass.real_estate
    DROP CONSTRAINT IF EXISTS real_estate_active_requires_token_ck;
ALTER TABLE ass.real_estate
    ADD CONSTRAINT real_estate_active_requires_token_ck
    CHECK (NOT active OR token_id IS NOT NULL) NOT VALID;

COMMENT ON COLUMN ass.real_estate.active IS
    'Projection of status_id through ass.real_estate_status.is_public, '
    'maintained by real_estate_sync_active_trg. Read-only in practice: writing '
    'it has no effect, the trigger recomputes it from status_id on every '
    'write. It exists because the listing and the detail queries filter on it; '
    'status_id is the value to set.';

COMMENT ON COLUMN ass.real_estate_status.is_public IS
    'Whether an asset in this status is served to the public — and, since '
    '000021, the single source of ass.real_estate.active. Changing it here '
    'does not update assets already stored: any migration that does must '
    'follow with UPDATE ass.real_estate SET status_id = status_id on the '
    'affected rows.';
