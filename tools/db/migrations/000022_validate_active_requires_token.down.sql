-- Rollback 000022 — the stored rows stop being covered.
--
--   PostgreSQL has no "un-validate". Once a constraint is validated, the only
--   way back to the previous state is to drop it and add it again as NOT
--   VALID. That looks heavier than the migration it reverses, and it is not a
--   workaround: it restores exactly the definition 000021 left behind, down to
--   the NOT VALID flag.
--
--   Both statements run inside the same transaction — golang-migrate wraps
--   each migration — so the table is never left without the constraint.
--
-- WHAT IT COSTS
--   New writes stay protected: NOT VALID still rejects every INSERT and UPDATE
--   that would put an asset on sale with no token. What is given up is the
--   certification of the rows already stored. They are compliant today, so
--   nothing breaks immediately; the database simply stops being able to prove
--   it, and the planner loses a fact it could rely on.
--
--   Rolling this back is therefore only useful to undo the migration itself,
--   never to make room for a non-compliant row. If a row needs to violate the
--   rule, the rule is wrong and that is a design discussion, not a rollback.
--
-- ⚠️  RE-RUNNING 000022 AFTERWARDS RESCANS THE TABLE
--   Going down then up is not free on a large table: the second VALIDATE reads
--   every row again. At this scale it is instant, and it will stay the honest
--   check rather than a cached answer.

ALTER TABLE ass.real_estate
    DROP CONSTRAINT IF EXISTS real_estate_active_requires_token_ck;

ALTER TABLE ass.real_estate
    ADD CONSTRAINT real_estate_active_requires_token_ck
    CHECK (NOT active OR token_id IS NOT NULL) NOT VALID;
