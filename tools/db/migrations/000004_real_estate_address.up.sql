-- 000004 — ass.real_estate_address (1-1 avec ass.real_estate)
--
-- PK = FK : real_estate_id est à la fois clé primaire et clé étrangère.
-- Deux adresses pour un même bien deviennent impossibles, et ON DELETE CASCADE
-- interdit la ligne orpheline. C'est ce qui rend une table 1-1 sûre.
--
-- L'adresse est la donnée la plus sensible du modèle : la localisation exacte
-- d'un bien tokenisé. L'isoler permet de ne pas la joindre sur les endpoints
-- publics de liste et de n'exposer que city / country_code tant que
-- l'utilisateur n'est pas KYC.

CREATE TABLE IF NOT EXISTS ass.real_estate_address (
    real_estate_id    INTEGER      PRIMARY KEY
        REFERENCES ass.real_estate(id) ON DELETE CASCADE,

    street            TEXT         NOT NULL,
    street_complement TEXT,
    postal_code       TEXT         NOT NULL,
    city              TEXT         NOT NULL,
    region            TEXT,
    -- ISO 3166-1 alpha-2
    country_code      CHAR(2)      NOT NULL,

    -- NUMERIC(9,6) : ~0,11 m de résolution. Jamais de float sur des
    -- coordonnées que l'on compare ou que l'on indexe.
    latitude          NUMERIC(9,6),
    longitude         NUMERIC(9,6),

    created_at        TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT now(),

    CONSTRAINT real_estate_address_country_code_ck
        CHECK (country_code ~ '^[A-Z]{2}$'),
    CONSTRAINT real_estate_address_latitude_ck
        CHECK (latitude  IS NULL OR latitude  BETWEEN  -90 AND  90),
    CONSTRAINT real_estate_address_longitude_ck
        CHECK (longitude IS NULL OR longitude BETWEEN -180 AND 180),
    -- Une coordonnée partielle n'a aucun sens : les deux ou aucune.
    CONSTRAINT real_estate_address_coordinates_ck
        CHECK ((latitude IS NULL) = (longitude IS NULL))
);

-- Filtres de recherche les plus probables sur une place de marché.
CREATE INDEX IF NOT EXISTS real_estate_address_country_city_idx
    ON ass.real_estate_address (country_code, city);
