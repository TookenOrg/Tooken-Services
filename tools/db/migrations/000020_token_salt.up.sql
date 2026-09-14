-- 000020 — A token's identity is its salt, not its name
--
-- WHAT IS BROKEN TODAY
--   Publishing an asset deploys a T-REX suite, inserts a blk.token row, then
--   binds the asset. If the chain succeeds and the write fails, the token
--   exists on-chain and in blk.token, while the asset stays a draft with
--   token_id IS NULL.
--
--   Replaying is then impossible. tokenForPublication sees no token on the
--   asset and deploys again; token.go:67 finds the existing row by name and
--   answers "Token already exists". Nothing links the orphan row back to the
--   asset, so the asset can never be published again without a manual UPDATE.
--
-- WHY THE SALT AND NOT THE NAME
--   The obvious repair is to look the orphan up by token_name — 000019 made it
--   unique, so the lookup is sound. It is still the wrong key.
--
--   token_name is built by defineTokenName from the asset *title*, and the
--   title is editable. If it changed between the failed attempt and the
--   replay, a lookup by name misses the orphan, the deployment starts over,
--   and it is the factory that refuses — because the salt, unlike the name,
--   did not change. The error then surfaces from the chain, opaque, instead of
--   being settled in the database.
--
--   defineSalt builds "re-<asset id>" from a primary key. It is stable by
--   construction, and it is the identity the T-REX factory itself uses to
--   derive the deployment address. That makes it the only honest key for
--   "has this asset already been deployed?".
--
-- WHAT IS BACKFILLED, AND WHAT IS DELIBERATELY NOT
--   Every token bound to an asset was deployed through CreateTokenWithSalt
--   with defineSalt(asset id), so its salt can be reconstructed exactly.
--   real_estate_token_id_uk guarantees one asset per token, so the backfill
--   cannot produce two identical salts.
--
--   Tokens bound to no asset are left NULL on purpose. They are of two kinds
--   and the schema cannot tell them apart:
--     - tokens from the manual POST /contract/token, whose salt was the name;
--     - orphans from the very failure described above, whose salt is "re-<id>".
--
--   Guessing would be worse than leaving the column empty: writing the name as
--   a salt onto an orphan would make a later publication deploy a second
--   contract for an asset that already has one. List them and decide each case
--   deliberately:
--
--     SELECT t.id, t.token_name, t.address
--     FROM blk.token t
--     LEFT JOIN ass.real_estate re ON re.token_id = t.id
--     WHERE re.id IS NULL;
--
-- WHY A PARTIAL UNIQUE INDEX
--   Two tokens sharing a salt would mean two contracts claiming the same
--   deployment slot, and adoption would pick one arbitrarily — the same damage
--   000019 prevents for the name. NULL is left out of the rule because it
--   means "unknown", not "none": several tokens may legitimately have no
--   recorded salt, and a plain UNIQUE would allow that anyway, but the partial
--   index states the intent rather than relying on how NULLs compare.
--
-- WHY A CHECK RATHER THAN NOT NULL
--   NOT NULL is out of reach: the tokens bound to no asset keep a NULL salt on
--   purpose, as explained above, and inventing one for them is the very thing
--   this migration refuses to do.
--
--   The empty string is a different matter. NULL means "not recorded"; '' means
--   nothing at all, and it can only come from a bug upstream. Letting it in
--   would be worse than it looks: tokenForPublication resolves on
--   "WHERE salt = $1", so an empty salt is a row the replay can never find —
--   precisely the orphan this migration exists to make adoptable — while
--   hiding the defect that produced it. The partial unique index would also
--   treat '' as a real value and reject the second such token.
--
--   modular_compliance_addr already carries that shape on this very table:
--   nullable, with a CHECK rejecting anything meaningless. This follows it.

ALTER TABLE blk.token
    ADD COLUMN IF NOT EXISTS salt TEXT;

UPDATE blk.token t
SET salt = 're-' || re.id
FROM ass.real_estate re
WHERE re.token_id = t.id
  AND t.salt IS NULL;

ALTER TABLE blk.token
    DROP CONSTRAINT IF EXISTS token_salt_not_empty_ck;
ALTER TABLE blk.token
    ADD CONSTRAINT token_salt_not_empty_ck
    CHECK (salt IS NULL OR salt <> '');

DROP INDEX IF EXISTS blk.token_salt_uk;
CREATE UNIQUE INDEX token_salt_uk
    ON blk.token (salt)
    WHERE salt IS NOT NULL;

COMMENT ON COLUMN blk.token.salt IS
    'Salt handed to the T-REX factory at deployment. Built by defineSalt as '
    '"re-<asset id>" for a tokenised asset, and equal to token_name for the '
    'manual POST /contract/token path. It is what tokenForPublication resolves '
    'on before deploying, so that a deployment interrupted before the asset '
    'was bound can be adopted on replay instead of being deployed twice. '
    'Derived from a primary key, it is stable even when the title — and '
    'therefore token_name — changes. NULL means "not recorded", which is the '
    'case for every token created before this migration that carries no asset.';
