-- 000022 — the VALIDATE that only ever existed in a shell history
--
--   000021 added real_estate_active_requires_token_ck as NOT VALID and stopped
--   there, on purpose: legacy assets were still active with no token, and
--   refusing to delete them to force a migration through was the honest call.
--   Its comment presents the VALIDATE as an operational step, to be run once
--   those rows are cleared.
--
--   They were cleared, and the statement was run — by hand, in a terminal.
--   Which means it does not exist. Nothing in this repository performs it, so
--   every database rebuilt from tools/db/baseline.sql comes back NOT VALID:
--
--     $ psql -d fresh_db < tools/db/baseline.sql
--     $ SELECT convalidated FROM pg_constraint
--       WHERE conname = 'real_estate_active_requires_token_ck';
--      f
--
--   The baseline is not at fault, it is telling the truth: pg_get_constraintdef
--   only prints NOT VALID when convalidated is false. The invariant simply was
--   not validated anywhere the repository can see.
--
--   This is §20 all over again — what is not in the repository does not exist —
--   applied to a one-line statement that felt too small to deserve a file.
--   A VALIDATE is a schema change like any other.
--
-- WHAT THIS ACTUALLY CHANGES
--   Nothing for new writes. NOT VALID already rejected every INSERT and UPDATE
--   that would leave an asset active without a token; that has been true since
--   000021. What this adds is the other half: the rows already stored are read
--   once and certified to comply, so the invariant finally holds over the whole
--   table rather than over its future only.
--
--   A validated constraint is also usable by the planner as a proven fact,
--   which a NOT VALID one is not.
--
-- ⚠️  THIS MIGRATION IS MEANT TO BE ABLE TO FAIL
--   VALIDATE reads every row and aborts if a single one violates the rule.
--   That is the point: it is the proof that the cleanup really happened, on
--   this database and not just on someone's laptop. Before running it
--   elsewhere, look:
--
--     SELECT id, title FROM ass.real_estate
--     WHERE active AND token_id IS NULL AND deleted_at IS NULL;
--
--   If it returns rows, settle them first — republish each one, which deploys
--   the missing token, or unpublish it. Do not weaken the constraint to get
--   the migration through.
--
-- SAFE TO REPLAY
--   Validating an already validated constraint is a no-op, not an error, so
--   re-running this migration on a database that is already compliant does
--   nothing. The lock taken is SHARE UPDATE EXCLUSIVE: reads and writes carry
--   on during the scan.

ALTER TABLE ass.real_estate
    VALIDATE CONSTRAINT real_estate_active_requires_token_ck;

COMMENT ON CONSTRAINT real_estate_active_requires_token_ck ON ass.real_estate IS
    'No asset can be on sale without a token behind it. Added NOT VALID by '
    '000021 so the legacy rows would not block it, validated by 000022 once '
    'they were cleared — so it now holds for stored rows as well as new ones.';
