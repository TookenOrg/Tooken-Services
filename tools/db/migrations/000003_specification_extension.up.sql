-- 000003 — Extension de ass.real_estate_specification
--
-- Ajoute les diagnostics énergétiques (obligatoires à la vente en France comme
-- au Luxembourg) et resserre les garde-fous sur les mesures existantes.

ALTER TABLE ass.real_estate_specification
    ADD COLUMN IF NOT EXISTS energy_class CHAR(1),
    ADD COLUMN IF NOT EXISTS ges_class    CHAR(1);

ALTER TABLE ass.real_estate_specification
    ADD CONSTRAINT real_estate_specification_energy_class_ck
        CHECK (energy_class IS NULL OR energy_class BETWEEN 'A' AND 'G'),
    ADD CONSTRAINT real_estate_specification_ges_class_ck
        CHECK (ges_class IS NULL OR ges_class BETWEEN 'A' AND 'G'),
    ADD CONSTRAINT real_estate_specification_built_year_ck
        CHECK (built_year IS NULL OR built_year BETWEEN 1000 AND 2200),
    ADD CONSTRAINT real_estate_specification_measures_ck
        CHECK ((lot_size     IS NULL OR lot_size     >= 0)
           AND (surface_area IS NULL OR surface_area >= 0)
           AND (pool_size    IS NULL OR pool_size    >= 0)
           AND (terrace_size IS NULL OR terrace_size >= 0));
