-- Rollback 000018 — the optional details become mandatory again.
--
-- THIS MIGRATION CAN LEGITIMATELY FAIL
--   SET NOT NULL scans the table and refuses if a single NULL remains. Any
--   asset recorded while the columns were optional — which is the whole point
--   of 000018 — will stop it.
--
--   That failure is the honest outcome, and it is left as is. The alternative
--   would be to backfill with zero, which would state that a property was
--   measured and has no pool when in truth nobody ever looked. 000016 refused
--   the same shortcut for the issuer, and §11.20 refused it for the currency.
--
--   Find what blocks it with:
--
--     SELECT real_estate_id, lot_size, pool_size, terrace_size,
--            bedroom_number, bathroom_number, built_year
--     FROM ass.real_estate_specification
--     WHERE lot_size IS NULL OR pool_size IS NULL OR terrace_size IS NULL
--        OR bedroom_number IS NULL OR bathroom_number IS NULL
--        OR built_year IS NULL;
--
--   and decide each row deliberately before replaying this rollback.

ALTER TABLE ass.real_estate_specification
    ALTER COLUMN lot_size        SET NOT NULL,
    ALTER COLUMN pool_size       SET NOT NULL,
    ALTER COLUMN terrace_size    SET NOT NULL,
    ALTER COLUMN bedroom_number  SET NOT NULL,
    ALTER COLUMN bathroom_number SET NOT NULL,
    ALTER COLUMN built_year      SET NOT NULL;

COMMENT ON COLUMN ass.real_estate_specification.lot_size IS NULL;
COMMENT ON COLUMN ass.real_estate_specification.pool_size IS NULL;
COMMENT ON COLUMN ass.real_estate_specification.terrace_size IS NULL;
COMMENT ON COLUMN ass.real_estate_specification.surface_area IS NULL;
