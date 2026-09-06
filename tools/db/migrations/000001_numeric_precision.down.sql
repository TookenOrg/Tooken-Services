-- Rollback 000001 — retour à la virgule flottante.
-- Attention : cette descente est destructrice au sens de la précision
-- (les décimales exactes sont irrémédiablement arrondies).

ALTER TABLE ass.real_estate_shares_config
    DROP CONSTRAINT IF EXISTS real_estate_shares_config_total_shares_ck,
    DROP CONSTRAINT IF EXISTS real_estate_shares_config_price_per_share_ck;

ALTER TABLE ass.real_estate_specification
    ALTER COLUMN lot_size     TYPE DOUBLE PRECISION USING lot_size::double precision,
    ALTER COLUMN surface_area TYPE DOUBLE PRECISION USING surface_area::double precision,
    ALTER COLUMN pool_size    TYPE DOUBLE PRECISION USING pool_size::double precision,
    ALTER COLUMN terrace_size TYPE DOUBLE PRECISION USING terrace_size::double precision;

ALTER TABLE ass.real_estate_shares_config
    ALTER COLUMN total_shares    TYPE INTEGER          USING total_shares::integer,
    ALTER COLUMN price_per_share TYPE DOUBLE PRECISION USING price_per_share::double precision,
    ALTER COLUMN yield           TYPE DOUBLE PRECISION USING yield::double precision;
