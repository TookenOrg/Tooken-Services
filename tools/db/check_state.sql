-- Which migrations are actually applied?
--
-- Use this when migrations were applied by hand (copy/paste in a SQL client)
-- instead of through golang-migrate, to find out where the database really
-- stands before handing control over to the tool.
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
            SELECT 1 FROM pg_constraint WHERE conname = 'fk_issuance_orders_asset'))
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
        (15, EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_issuance_orders_asset'))
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
