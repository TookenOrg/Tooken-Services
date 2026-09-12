-- Where does this database really stand?
--
-- Two things are reported here, both read-only:
--
--   1. which migrations are actually applied -- use this when migrations were
--      applied by hand (copy/paste in a SQL client) instead of through
--      golang-migrate, to find out where the database stands before handing
--      control over to the tool;
--   2. the invariants the schema cannot express, listed at the end of the
--      file. They are enforced by the service layer; this is how we find out
--      whether one of them has been broken.
--
--   psql "$DATABASE_URL" -f tools/db/check_state.sql
--
-- Each migration leaves one uniquely named artifact behind; its presence is a
-- reliable proxy for "this migration ran".

WITH probes(version, label, applied) AS (
    VALUES
    (1, 'numeric_precision',      to_regclass('ass.real_estate_shares_config') IS NOT NULL AND EXISTS (
            SELECT 1 FROM pg_constraint WHERE conname = 'real_estate_shares_config_total_shares_ck')),
    (2, 'shares_config_extension', EXISTS (
            SELECT 1 FROM pg_constraint WHERE conname = 'real_estate_shares_config_currency_ck')),
    (3, 'specification_extension', EXISTS (
            SELECT 1 FROM pg_constraint WHERE conname = 'real_estate_specification_energy_class_ck')),
    (4, 'real_estate_address',     to_regclass('ass.real_estate_address') IS NOT NULL),
    (5, 'real_estate_media',       to_regclass('ass.real_estate_media')   IS NOT NULL),
    (6, 'real_estate_lifecycle',   to_regclass('ass.real_estate_status')  IS NOT NULL),
    (7, 'issuer',                  to_regclass('ass.issuer')              IS NOT NULL),
    (8, 'real_estate_token_link',  EXISTS (
            SELECT 1 FROM pg_constraint WHERE conname = 'token_nb_decimal_ck')),
    (9, 'user_role_kyc',           EXISTS (
            SELECT 1 FROM pg_constraint WHERE conname = 'users_role_ck')),
    (10, 'order_reserved_shares',  EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'iss' AND table_name = 'issuance_order_statuses'
              AND column_name = 'counts_as_reserved')),
    (11, 'real_estate_soft_delete', EXISTS (
            SELECT 1 FROM pg_constraint WHERE conname = 'real_estate_deleted_status_ck')),
    (12, 'real_estate_contract_alignment', to_regclass('ass.real_estate_estate_type_idx') IS NOT NULL),
    -- 000013 only calls setval(): it leaves no artifact, so it cannot be
    -- probed. It is reported as applied when 000014 is, since migrations run
    -- in order. When 000014 is absent the answer is 'unknown', which keeps it
    -- out of the contiguous count below and makes it replay -- harmless, as
    -- resynchronising a sequence to its own max id is idempotent.
    (13, 'resync_identity_sequences', to_regclass('ass.issuer_lei_uk') IS NOT NULL),
    (14, 'issuer_lei_unique',      to_regclass('ass.issuer_lei_uk') IS NOT NULL),
    (15, 'referential_integrity',  EXISTS (
            SELECT 1 FROM pg_constraint WHERE conname = 'fk_issuance_orders_asset')),
    (16, 'real_estate_issuer_required', EXISTS (
            SELECT 1 FROM information_schema.columns
            WHERE table_schema = 'ass' AND table_name = 'real_estate'
              AND column_name = 'issuer_id' AND is_nullable = 'NO'))
)
SELECT version,
       label,
       CASE
           WHEN version = 13 AND NOT applied THEN 'unknown (no artifact)'
           WHEN applied THEN 'applied'
           ELSE 'MISSING'
       END AS state
FROM probes
ORDER BY version;

-- Highest contiguous version reached: this is the number to pass to
-- `migrate force <n>`. It stops at the first gap on purpose — forcing past a
-- hole would silently skip a migration forever.
SELECT COALESCE(MAX(version), 0) AS force_this_version
FROM (
    WITH probes(version, applied) AS (
        VALUES
        (1, EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'real_estate_shares_config_total_shares_ck')),
        (2, EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'real_estate_shares_config_currency_ck')),
        (3, EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'real_estate_specification_energy_class_ck')),
        (4, to_regclass('ass.real_estate_address') IS NOT NULL),
        (5, to_regclass('ass.real_estate_media')   IS NOT NULL),
        (6, to_regclass('ass.real_estate_status')  IS NOT NULL),
        (7, to_regclass('ass.issuer')              IS NOT NULL),
        (8, EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'token_nb_decimal_ck')),
        (9, EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'users_role_ck')),
        (10, EXISTS (SELECT 1 FROM information_schema.columns
                     WHERE table_schema = 'iss' AND table_name = 'issuance_order_statuses'
                       AND column_name = 'counts_as_reserved')),
        (11, EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'real_estate_deleted_status_ck')),
        (12, to_regclass('ass.real_estate_estate_type_idx') IS NOT NULL),
        (13, to_regclass('ass.issuer_lei_uk') IS NOT NULL),
        (14, to_regclass('ass.issuer_lei_uk') IS NOT NULL),
        (15, EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_issuance_orders_asset')),
        (16, EXISTS (SELECT 1 FROM information_schema.columns
                     WHERE table_schema = 'ass' AND table_name = 'real_estate'
                       AND column_name = 'issuer_id' AND is_nullable = 'NO'))
    )
    SELECT version FROM probes p
    WHERE NOT EXISTS (SELECT 1 FROM probes q WHERE q.version <= p.version AND NOT q.applied)
) contiguous;

-- Leftovers from an earlier draft of migration 000003, which used to add an
-- `attributes JSONB` column before it was dropped from the design. If this
-- returns a row, run:
--   ALTER TABLE ass.real_estate_specification DROP COLUMN attributes;
--   DROP INDEX IF EXISTS ass.real_estate_specification_attributes_idx;
SELECT 'stale column ass.real_estate_specification.attributes' AS leftover
FROM information_schema.columns
WHERE table_schema = 'ass'
  AND table_name   = 'real_estate_specification'
  AND column_name  = 'attributes';

-- Does golang-migrate already track this database?
DO $$
DECLARE v record;
BEGIN
    IF to_regclass('public.schema_migrations') IS NULL THEN
        RAISE NOTICE 'schema_migrations absent - golang-migrate has never run here';
    ELSE
        FOR v IN SELECT version, dirty FROM public.schema_migrations LOOP
            RAISE NOTICE 'schema_migrations: version=% dirty=%', v.version, v.dirty;
        END LOOP;
    END IF;
END $$;


-- ---------------------------------------------------------------------------
-- Invariants the schema cannot express
-- ---------------------------------------------------------------------------
--
-- Invariant: a visible asset must stand behind an active issuer.
--
-- The rule is enforced in the service layer -- an issuer carrying published
-- assets is refused any status other than 'active' (issuer.go,
-- explainIssuerWriteRefusal). The database knows nothing about it, and cannot:
-- a CHECK constraint only ever sees one row of one table, and this one spans
-- ass.real_estate and ass.issuer.
--
-- A trigger could close the gap, and is deliberately not used. This schema
-- already carries one on ass.real_estate.active, and that trigger is precisely
-- why `active` drifted from real_estate_status.is_public (see the note above
-- GetIssuerGuardState in issuer_db.go). Adding a third piece of invisible
-- logic to guard against a case never observed would buy a guarantee at the
-- price of the problem it guards against.
--
-- So the invariant stays in the code, and this block makes it observable.
-- Silence (beyond the "holds" notice) means it holds. Any VIOLATION line means
-- investors can see an offer whose vehicle is suspended or dissolved; fix it by
-- unpublishing the asset or by putting the issuer back to 'active'.
--
-- "Visible" is read as re.active, not as real_estate_status.is_public: active
-- is what the listing and detail queries filter on, and it is the same column
-- the service counts when it refuses the status change. Checking the rule with
-- anything else would report violations the rule never claimed to prevent.
--
-- Like every probe above, this must not fail on a half-migrated database: the
-- whole point of the file is to be runnable when the state is unknown. Hence
-- the guard -- a missing table here means the invariant does not exist yet.
DO $$
DECLARE
    v         record;
    violations int := 0;
BEGIN
    IF to_regclass('ass.real_estate') IS NULL
       OR to_regclass('ass.issuer') IS NULL
       OR to_regclass('ass.issuer_status') IS NULL
       OR NOT EXISTS (SELECT 1 FROM information_schema.columns
                      WHERE table_schema = 'ass' AND table_name = 'real_estate'
                        AND column_name = 'deleted_at')
    THEN
        RAISE NOTICE 'issuer invariant: skipped - schema too old to carry it';
        RETURN;
    END IF;

    FOR v IN
        SELECT re.id AS real_estate_id,
               re.title,
               i.id   AS issuer_id,
               i.name AS issuer_name,
               ist.code AS issuer_status
        FROM ass.real_estate re
        JOIN ass.issuer i          ON i.id   = re.issuer_id
        JOIN ass.issuer_status ist ON ist.id = i.status_id
        WHERE re.deleted_at IS NULL
          AND re.active
          AND i.status_id <> 2      -- 2 = 'active', seeded by migration 000007
        ORDER BY re.id
    LOOP
        violations := violations + 1;
        RAISE NOTICE 'VIOLATION: real estate % (%) is visible under issuer % (%) in status ''%''',
            v.real_estate_id, v.title, v.issuer_id, v.issuer_name, v.issuer_status;
    END LOOP;

    IF violations = 0 THEN
        RAISE NOTICE 'issuer invariant: holds - every visible asset stands behind an active issuer';
    ELSE
        RAISE WARNING 'issuer invariant: % visible asset(s) stand behind a non-active issuer', violations;
    END IF;
END $$;
