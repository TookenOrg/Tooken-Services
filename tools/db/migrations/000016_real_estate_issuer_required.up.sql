-- 000016 — An asset is always issued by a vehicle
--
-- WHAT WAS MISSING
--   Antony settled the business rule at §11.21: an asset is ALWAYS attached to
--   an issuer. Since then the rule has been held by the service alone —
--   publicationRequirements refuses to publish without issuer_id, and
--   checkIssuerAttachment refuses to move an asset onto a dissolved vehicle.
--
--   But the column stayed nullable. An INSERT typed directly in psql, a
--   restored dump, or any future code path that forgets the check goes
--   straight through. A rule that only one caller enforces is a convention,
--   not an invariant: the schema is the only place that answers for every
--   writer at once.
--
--   The foreign key (real_estate_issuer_id_fkey, ON DELETE RESTRICT) already
--   guarantees that a NAMED issuer exists. What it cannot say is that one must
--   be named at all — that is exactly what NOT NULL adds, and nothing else.
--
-- WHY NO BACKFILL
--   §11.25 sized this as a separate project: "510 rows are in this state,
--   set the test data aside, backfill, then add the constraint". That was
--   measured before the demo database was cleaned. It no longer holds: the 13
--   assets now in the database ALL carry an issuer, so there is nothing to
--   backfill and the constraint applies as it stands.
--
--   No default issuer is invented for the rows that might still be missing one
--   elsewhere. Attaching an asset to a vehicle picked by a migration would
--   write a plausible but false line into the register — the same mistake
--   §11.20 corrected when the currency defaulted to EUR. If a database still
--   holds an asset without an issuer, this migration FAILS and says so, which
--   is the honest outcome. Find them with:
--
--     SELECT id, title, status_id FROM ass.real_estate WHERE issuer_id IS NULL;
--
--   and attach each one deliberately before replaying.
--
-- LOCKING
--   SET NOT NULL takes an ACCESS EXCLUSIVE lock and scans the table to prove
--   no NULL remains. On this table (13 rows, and ~500 at its largest) the scan
--   is instantaneous. The trick worth knowing if it ever grows: add a
--   validated CHECK (issuer_id IS NOT NULL) first, and PostgreSQL 12+ skips
--   the scan. Not worth the two extra steps here.

ALTER TABLE ass.real_estate
    ALTER COLUMN issuer_id SET NOT NULL;

COMMENT ON COLUMN ass.real_estate.issuer_id IS
    'The legal vehicle issuing the shares. Mandatory: a token represents a '
    'stake in an issuer, never in the walls themselves (decision A1). The API '
    'refuses a creation without it (400); this constraint answers for every '
    'other writer.';
