-- 000005 — ass.real_estate_media (1-N)
--
-- ass.real_estate.imageurl ne porte qu'une seule image. Une fiche de bien
-- immobilier en exige une galerie.
--
-- imageurl n'est PAS supprimée ici : la migration 000009 la recalculera depuis
-- la couverture. Supprimer une colonne encore lue par le code en production est
-- la meilleure façon de casser le service pendant le déploiement.

CREATE TABLE IF NOT EXISTS ass.real_estate_media (
    id             SERIAL       PRIMARY KEY,
    real_estate_id INTEGER      NOT NULL
        REFERENCES ass.real_estate(id) ON DELETE CASCADE,

    url            TEXT         NOT NULL,
    alt_text       TEXT,
    media_type     TEXT         NOT NULL DEFAULT 'image',
    position       INTEGER      NOT NULL DEFAULT 0,
    is_cover       BOOLEAN      NOT NULL DEFAULT FALSE,

    created_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),

    CONSTRAINT real_estate_media_type_ck
        CHECK (media_type IN ('image', 'video', 'floorplan', 'virtual_tour')),
    CONSTRAINT real_estate_media_position_ck
        CHECK (position >= 0)
);

-- Une seule couverture par bien, garantie par la base plutôt que par le code.
CREATE UNIQUE INDEX IF NOT EXISTS real_estate_media_single_cover_uk
    ON ass.real_estate_media (real_estate_id)
    WHERE is_cover;

CREATE INDEX IF NOT EXISTS real_estate_media_real_estate_position_idx
    ON ass.real_estate_media (real_estate_id, position);

-- Reprise de l'existant : l'image unique actuelle devient la couverture.
INSERT INTO ass.real_estate_media (real_estate_id, url, position, is_cover)
SELECT re.id, re.imageurl, 0, TRUE
FROM ass.real_estate re
WHERE re.imageurl IS NOT NULL
  AND re.imageurl <> ''
  AND NOT EXISTS (
      SELECT 1 FROM ass.real_estate_media m WHERE m.real_estate_id = re.id
  );
