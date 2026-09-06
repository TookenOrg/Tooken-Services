-- 000007 — ass.issuer : le véhicule juridique émetteur
--
-- Le token ne représente pas la brique et le mortier : il représente des parts
-- d'un véhicule qui détient le bien (décision A1). C'est ce qui permet à
-- l'ERC-3643 d'avoir un sens juridique — un transfert de token est un
-- transfert de titre, pas une mutation immobilière.
--
-- Cardinalité 1 issuer -> N biens (décision A2) : un véhicule de titrisation
-- luxembourgeois porte plusieurs compartiments, un par bien. Le compartiment
-- lui-même est identifié par shares_config.compartment_ref (migration 000002).

CREATE TABLE IF NOT EXISTS ass.issuer_status (
    id       SMALLINT PRIMARY KEY,
    code     TEXT     NOT NULL UNIQUE,
    label    TEXT     NOT NULL,
    is_final BOOLEAN  NOT NULL DEFAULT FALSE
);

INSERT INTO ass.issuer_status (id, code, label, is_final) VALUES
    (1, 'draft',      'Brouillon',    FALSE),
    (2, 'active',     'Actif',        FALSE),
    (3, 'suspended',  'Suspendu',     FALSE),
    (4, 'dissolved',  'Dissous',      TRUE)
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS ass.issuer (
    id                  SERIAL       PRIMARY KEY,

    name                TEXT         NOT NULL,
    legal_form          TEXT         NOT NULL,
    registration_number TEXT,
    -- Juridiction du véhicule (LU pour le POC). ISO 3166-1 alpha-2.
    country_code        CHAR(2)      NOT NULL DEFAULT 'LU',
    lei_code            CHAR(20),

    status_id           SMALLINT     NOT NULL DEFAULT 1
        REFERENCES ass.issuer_status(id),

    created_at          TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT now(),

    CONSTRAINT issuer_country_code_ck
        CHECK (country_code ~ '^[A-Z]{2}$'),
    CONSTRAINT issuer_lei_ck
        CHECK (lei_code IS NULL OR lei_code ~ '^[A-Z0-9]{20}$')
);

-- Un numéro d'immatriculation est unique dans sa juridiction.
CREATE UNIQUE INDEX IF NOT EXISTS issuer_registration_uk
    ON ass.issuer (country_code, registration_number)
    WHERE registration_number IS NOT NULL;

-- Nullable : un bien en brouillon n'a pas encore d'émetteur rattaché.
ALTER TABLE ass.real_estate
    ADD COLUMN IF NOT EXISTS issuer_id INTEGER
        REFERENCES ass.issuer(id) ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS real_estate_issuer_id_idx
    ON ass.real_estate (issuer_id);
