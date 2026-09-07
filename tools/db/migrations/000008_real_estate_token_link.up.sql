-- 000008 — Lien réel entre le bien et son contrat
--
-- ass.real_estate.contract_address est une colonne texte libre : rien ne
-- garantit que l'adresse corresponde à un token connu, rien n'empêche deux
-- biens de pointer le même contrat, et une chaîne vide est indistinguable
-- d'un bien non encore tokenisé.
--
-- token_id NULL devient un état métier explicite : « pas encore tokenisé ».
--
-- contract_address est CONSERVÉE et synchronisée par trigger, pour la même
-- raison que `active` en 000006 : le code de production la lit encore.

ALTER TABLE ass.real_estate
    ADD COLUMN IF NOT EXISTS token_id INTEGER
        REFERENCES blk.token(id) ON DELETE RESTRICT;

-- Reprise : on relie ce qui peut l'être, en normalisant la casse
-- (les adresses Ethereum circulent en EIP-55 comme en minuscules).
UPDATE ass.real_estate re
SET token_id = t.id
FROM blk.token t
WHERE re.token_id IS NULL
  AND re.contract_address IS NOT NULL
  AND re.contract_address <> ''
  AND lower(re.contract_address) = lower(t.address);

-- La chaîne vide était le marqueur de fait de « non tokenisé ». Le marqueur
-- explicite est désormais token_id IS NULL : on normalise pour qu'il n'existe
-- qu'une seule représentation de cet état.
UPDATE ass.real_estate
SET contract_address = NULL
WHERE contract_address = '';

-- Un token ne peut représenter qu'un seul bien.
CREATE UNIQUE INDEX IF NOT EXISTS real_estate_token_id_uk
    ON ass.real_estate (token_id)
    WHERE token_id IS NOT NULL;

CREATE OR REPLACE FUNCTION ass.real_estate_sync_contract_address()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.token_id IS NULL THEN
        NEW.contract_address := NULL;
    ELSE
        SELECT t.address INTO NEW.contract_address
        FROM blk.token t WHERE t.id = NEW.token_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS real_estate_sync_contract_address_trg ON ass.real_estate;
CREATE TRIGGER real_estate_sync_contract_address_trg
    BEFORE INSERT OR UPDATE OF token_id ON ass.real_estate
    FOR EACH ROW EXECUTE FUNCTION ass.real_estate_sync_contract_address();

-- decimals est IMMUABLE une fois le token déployé : T-REX n'expose aucun
-- SetDecimals (vérifié dans internal/blockchain/contracts/bindings/Token.go).
-- Le choix est donc définitif, token par token.
--
-- La borne haute 18 est celle imposée on-chain par Token.init, et déjà
-- appliquée côté service (services/token.go). On la réplique ici pour que la
-- base reste cohérente même si une écriture contourne le service.
--
-- La valeur reste LIBRE entre 0 et 18 : un bien dont les parts doivent rester
-- indivisibles se déploie à 0, un bien destiné à la revente fractionnée à 18.
-- Conséquence à connaître : à 0, céder 0,1 part sera impossible pour toujours.
ALTER TABLE blk.token
    ADD CONSTRAINT token_nb_decimal_ck
        CHECK (nb_decimal BETWEEN 0 AND 18) NOT VALID;
