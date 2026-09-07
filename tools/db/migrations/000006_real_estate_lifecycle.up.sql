-- 000006 — Cycle de vie du bien (dette 🟠 : active bool)
--
-- `active BOOLEAN` ne sait représenter que deux états. Le CRUD à venir
-- (issue #28) va nécessairement créer des biens non publiés, puis les faire
-- passer en revue, en levée, financés, clôturés. Un booléen ne peut pas
-- exprimer cela, et il empêche aussi de distinguer un 404 d'un bien inactif.
--
-- Le référentiel reprend la forme de iss.issuance_order_statuses
-- (id, code, label, is_final) déjà en place : cohérence du schéma.
--
-- `active` est CONSERVÉE et synchronisée par un trigger : le code de
-- production la lit encore. Elle sera retirée une fois la couche database
-- migrée (étape 3).

CREATE TABLE IF NOT EXISTS ass.real_estate_status (
    id         SMALLINT     PRIMARY KEY,
    code       TEXT         NOT NULL UNIQUE,
    label      TEXT         NOT NULL,
    is_final   BOOLEAN      NOT NULL DEFAULT FALSE,
    -- Visible par un utilisateur non authentifié.
    is_public  BOOLEAN      NOT NULL DEFAULT FALSE,
    position   SMALLINT     NOT NULL DEFAULT 0
);

INSERT INTO ass.real_estate_status (id, code, label, is_final, is_public, position) VALUES
    (1, 'draft',      'Brouillon',         FALSE, FALSE, 10),
    (2, 'in_review',  'En revue',          FALSE, FALSE, 20),
    (3, 'published',  'Publié',            FALSE, TRUE,  30),
    (4, 'fundraising','En levée',          FALSE, TRUE,  40),
    (5, 'funded',     'Financé',           FALSE, TRUE,  50),
    (6, 'closed',     'Clôturé',           TRUE,  TRUE,  60),
    (7, 'cancelled',  'Annulé',            TRUE,  FALSE, 70)
ON CONFLICT (id) DO NOTHING;

ALTER TABLE ass.real_estate
    ADD COLUMN IF NOT EXISTS status_id  SMALLINT
        REFERENCES ass.real_estate_status(id),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    ADD COLUMN IF NOT EXISTS published_at TIMESTAMPTZ;

-- Reprise : actif -> publié, inactif -> brouillon.
UPDATE ass.real_estate
SET status_id = CASE WHEN COALESCE(active, FALSE) THEN 3 ELSE 1 END
WHERE status_id IS NULL;

ALTER TABLE ass.real_estate
    ALTER COLUMN status_id SET DEFAULT 1,
    ALTER COLUMN status_id SET NOT NULL;

-- Tant que `active` existe, les deux représentations doivent rester cohérentes,
-- quelle que soit celle que le code écrit.
CREATE OR REPLACE FUNCTION ass.real_estate_sync_active()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' OR NEW.status_id IS DISTINCT FROM OLD.status_id THEN
        NEW.active := (NEW.status_id IN (3, 4, 5));
    ELSIF NEW.active IS DISTINCT FROM OLD.active THEN
        NEW.status_id := CASE WHEN NEW.active THEN 3 ELSE 1 END;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS real_estate_sync_active_trg ON ass.real_estate;
CREATE TRIGGER real_estate_sync_active_trg
    BEFORE INSERT OR UPDATE ON ass.real_estate
    FOR EACH ROW EXECUTE FUNCTION ass.real_estate_sync_active();

CREATE INDEX IF NOT EXISTS real_estate_status_id_idx
    ON ass.real_estate (status_id);
