-- 000002 — Extension de ass.real_estate_shares_config
--
-- Ajoute la devise, la valorisation dérivée, les bornes d'investissement,
-- la grille de frais et la référence de compartiment luxembourgeois.
--
-- total_valuation est une colonne GÉNÉRÉE : PostgreSQL la calcule et interdit
-- toute écriture. C'est la garantie structurelle qu'elle ne peut jamais
-- diverger de total_shares * price_per_share (règle : stocker deux des trois
-- grandeurs liées, dériver la troisième).

ALTER TABLE ass.real_estate_shares_config
    ADD COLUMN IF NOT EXISTS currency_code    CHAR(3)      NOT NULL DEFAULT 'EUR',
    ADD COLUMN IF NOT EXISTS compartment_ref  TEXT,
    ADD COLUMN IF NOT EXISTS min_investment   NUMERIC(20,8),
    ADD COLUMN IF NOT EXISTS max_investment   NUMERIC(20,8),
    ADD COLUMN IF NOT EXISTS entry_fee_rate   NUMERIC(7,4),
    ADD COLUMN IF NOT EXISTS mgmt_fee_rate    NUMERIC(7,4),
    ADD COLUMN IF NOT EXISTS exit_fee_rate    NUMERIC(7,4),
    ADD COLUMN IF NOT EXISTS updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now();

-- Colonne dérivée : ni ADD COLUMN IF NOT EXISTS ni GENERATED ne sont
-- combinables avec un DEFAULT, on la déclare séparément.
ALTER TABLE ass.real_estate_shares_config
    ADD COLUMN IF NOT EXISTS total_valuation NUMERIC(40,8)
        GENERATED ALWAYS AS (total_shares * price_per_share) STORED;

ALTER TABLE ass.real_estate_shares_config
    ADD CONSTRAINT real_estate_shares_config_currency_ck
        CHECK (currency_code ~ '^[A-Z]{3}$'),
    ADD CONSTRAINT real_estate_shares_config_investment_range_ck
        CHECK (min_investment IS NULL
               OR max_investment IS NULL
               OR min_investment <= max_investment),
    ADD CONSTRAINT real_estate_shares_config_fees_ck
        CHECK ((entry_fee_rate IS NULL OR entry_fee_rate BETWEEN 0 AND 100)
           AND (mgmt_fee_rate  IS NULL OR mgmt_fee_rate  BETWEEN 0 AND 100)
           AND (exit_fee_rate  IS NULL OR exit_fee_rate  BETWEEN 0 AND 100));

-- Une référence de compartiment identifie un patrimoine cloisonné : elle est
-- unique par nature. UNIQUE tolère plusieurs NULL en PostgreSQL, ce qui laisse
-- les biens non encore tokenisés sans contrainte.
CREATE UNIQUE INDEX IF NOT EXISTS real_estate_shares_config_compartment_ref_uk
    ON ass.real_estate_shares_config (compartment_ref)
    WHERE compartment_ref IS NOT NULL;
