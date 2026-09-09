-- 000014 — Un LEI ne désigne qu'une entité
--
-- Question laissée ouverte en §11.19 : `registration_number` et `lei_code`
-- doivent-ils être uniques ?
--
--   * `registration_number` l'est déjà, et seulement dans sa juridiction
--     (index partiel `issuer_registration_uk`, migration 000007) : deux pays
--     peuvent parfaitement attribuer le même numéro.
--   * `lei_code` ne l'était pas. Or le LEI est **mondialement** unique par
--     construction (ISO 17442) : deux lignes le partageant ne décrivent pas
--     deux émetteurs, elles décrivent deux fois le même — un doublon que rien
--     ne signalait, et qui rattache des biens à des véhicules qu'un rapport
--     réglementaire compterait ensuite comme distincts.
--
-- L'index est partiel : le LEI reste facultatif — toutes les structures n'en
-- ont pas — et plusieurs NULL ne sont pas un doublon.
--
-- Si la base porte déjà des doublons, la création échoue et le dit : mieux
-- vaut refuser la migration que valider en silence des lignes contradictoires.

CREATE UNIQUE INDEX IF NOT EXISTS issuer_lei_uk
    ON ass.issuer (lei_code)
    WHERE lei_code IS NOT NULL;
