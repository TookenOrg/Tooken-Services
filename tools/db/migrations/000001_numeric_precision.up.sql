-- 000001 — Précision numérique exacte (dette 🔴 : montants en float)
--
-- Les montants et mesures sont aujourd'hui en virgule flottante. Ils sont
-- confrontés à des valeurs on-chain exprimées en wei (entiers exacts sur 256
-- bits) : tout arrondi flottant est un écart de comptabilité.
--
-- Conventions retenues :
--   * montants          -> NUMERIC(20,8)
--   * nombre de parts   -> NUMERIC(20,0)  (une part est indivisible à l'émission ;
--                          la fractionnabilité vit on-chain via decimals = 18)
--   * mesures / surfaces-> NUMERIC(12,2)
--   * pourcentages      -> NUMERIC(7,4)   (ex. 4,2500 %)
--
-- USING ...::numeric est explicite : PostgreSQL refuse la conversion implicite
-- double precision -> numeric lors d'un ALTER TYPE.

ALTER TABLE ass.real_estate_shares_config
    ALTER COLUMN total_shares    TYPE NUMERIC(20,0) USING total_shares::numeric,
    ALTER COLUMN price_per_share TYPE NUMERIC(20,8) USING price_per_share::numeric,
    ALTER COLUMN yield           TYPE NUMERIC(7,4)  USING yield::numeric;

ALTER TABLE ass.real_estate_specification
    ALTER COLUMN lot_size     TYPE NUMERIC(12,2) USING lot_size::numeric,
    ALTER COLUMN surface_area TYPE NUMERIC(12,2) USING surface_area::numeric,
    ALTER COLUMN pool_size    TYPE NUMERIC(12,2) USING pool_size::numeric,
    ALTER COLUMN terrace_size TYPE NUMERIC(12,2) USING terrace_size::numeric;

-- Garde-fous métier : aucune de ces grandeurs ne peut être négative.
ALTER TABLE ass.real_estate_shares_config
    ADD CONSTRAINT real_estate_shares_config_total_shares_ck
        CHECK (total_shares IS NULL OR total_shares > 0),
    ADD CONSTRAINT real_estate_shares_config_price_per_share_ck
        CHECK (price_per_share IS NULL OR price_per_share >= 0);
