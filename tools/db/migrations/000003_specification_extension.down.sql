-- Rollback 000003

ALTER TABLE ass.real_estate_specification
    DROP CONSTRAINT IF EXISTS real_estate_specification_energy_class_ck,
    DROP CONSTRAINT IF EXISTS real_estate_specification_ges_class_ck,
    DROP CONSTRAINT IF EXISTS real_estate_specification_built_year_ck,
    DROP CONSTRAINT IF EXISTS real_estate_specification_measures_ck;

ALTER TABLE ass.real_estate_specification
    DROP COLUMN IF EXISTS energy_class,
    DROP COLUMN IF EXISTS ges_class;

-- Une version antérieure de cette migration ajoutait une colonne `attributes`
-- JSONB et son index GIN, retirés depuis. Le nettoyage est conservé ici pour
-- que les bases ayant joué cette version-là reviennent au même état que les
-- autres après un `down`.
DROP INDEX IF EXISTS ass.real_estate_specification_attributes_gin;
ALTER TABLE ass.real_estate_specification
    DROP COLUMN IF EXISTS attributes;
