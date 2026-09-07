-- 000011 — Soft delete for ass.real_estate
--
-- A real estate asset is referenced by orders, by a token and by the on-chain
-- registry: deleting the row would either fail on a foreign key or orphan an
-- investor's history. Deletion is therefore a state, not a DELETE.
--
-- The state is carried by two columns on purpose:
--   * status_id = 7 (cancelled) is what the business reads, and what the
--     000006 trigger turns into active = FALSE, so a deleted asset leaves the
--     public listing without any query having to know about deletion;
--   * deleted_at answers "when", which a status alone cannot, and which an
--     audit of a financial registry needs.
--
-- Two representations can drift, so a CHECK forbids the drift instead of
-- trusting the application: a row cannot be dated without being cancelled.
-- The reverse stays allowed — an asset can be cancelled for business reasons
-- (funding failed) without ever having been deleted.

ALTER TABLE ass.real_estate
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

ALTER TABLE ass.real_estate
    DROP CONSTRAINT IF EXISTS real_estate_deleted_status_ck;

-- 7 = 'cancelled'. The id is hardcoded like the (3, 4, 5) of the 000006
-- trigger: a CHECK cannot query ass.real_estate_status.
ALTER TABLE ass.real_estate
    ADD CONSTRAINT real_estate_deleted_status_ck
        CHECK (deleted_at IS NULL OR status_id = 7);

ALTER TABLE ass.real_estate
    DROP CONSTRAINT IF EXISTS real_estate_deleted_not_tokenized_ck;

-- A tokenized asset cannot be deleted, even softly: the token exists on chain,
-- investors hold it, and nothing off-chain can undo that. The service layer
-- refuses it with a readable error; this constraint is the guarantee that a
-- script, a migration or a future endpoint cannot bypass the rule.
ALTER TABLE ass.real_estate
    ADD CONSTRAINT real_estate_deleted_not_tokenized_ck
        CHECK (deleted_at IS NULL OR token_id IS NULL);

-- Every read filters on deleted_at IS NULL. The index only carries the rows
-- that remain, so it stays small as deletions accumulate.
CREATE INDEX IF NOT EXISTS real_estate_not_deleted_idx
    ON ass.real_estate (id)
    WHERE deleted_at IS NULL;

COMMENT ON COLUMN ass.real_estate.deleted_at IS
    'Soft deletion timestamp. NULL means the asset exists. Non-NULL implies status_id = 7 (cancelled) and forbids tokenization.';
