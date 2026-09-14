-- 000019 — The token name is an identity, not a label
--
-- WHAT WAS MISSING
--   blk.token carries token_name_unique in production. No migration creates
--   it: the constraint was added by hand, and blk.token itself predates this
--   migration history. A database built from tools/db/migrations therefore
--   holds the table without the constraint that production relies on.
--
--   The gap stayed invisible for as long as no machine could build the schema
--   from the repository. It surfaced the wrong way round, as an integration
--   test that collided on its second run against a constraint the repository
--   never mentions.
--
-- WHY THE NAME IS AN IDENTITY
--   token_name is the idempotency key of CreateToken. token.go states it in
--   as many words: "Mint/Burn resolve it by address from this row and
--   idempotency keys on token_name". GetTokenByName resolves a token with
--   WHERE token_name = $1 and reads a single row out of it.
--
--   Two rows sharing a name would make that lookup return one of them,
--   arbitrarily. A replayed deployment would then bind an asset to the wrong
--   contract address, and a mint would be sent to a token nobody meant. The
--   damage is on-chain and irreversible, which is why this belongs in the
--   schema rather than in a code review.
--
-- WHY IT IS SAFE
--   The rule is already honoured by the only writer. defineTokenName builds
--   "Tooken <title> #<asset id>", defineSymbol builds "TKN<asset id>" and
--   defineSalt builds "re-<asset id>" — all keyed on a primary key, so two
--   assets cannot produce the same name. This migration introduces no rule; it
--   writes down one the code has always followed, so the schema answers for
--   every writer instead of a single caller. 000016 made the same argument for
--   issuer_id: a rule only one caller enforces is a convention, not an
--   invariant.
--
-- IF IT FAILS
--   Duplicates exist. The migration stops, and that is the honest outcome. No
--   row is renamed or deleted to force it through: a token name is what Mint
--   and Burn resolve on, so rewriting one would point an operation at a
--   different contract. Find them with:
--
--     SELECT token_name, count(*), array_agg(id)
--     FROM blk.token
--     GROUP BY token_name
--     HAVING count(*) > 1;
--
--   and decide each pair deliberately — which token is real, which asset it
--   belongs to — before replaying.
--
-- WHY DROP THEN ADD
--   ADD CONSTRAINT has no IF NOT EXISTS, and 000015 already settled the idiom:
--   the pair keeps the migration replayable without a DO block. NOT VALID is
--   not available here — it exists for foreign keys and checks, never for
--   UNIQUE, which has to build its index in one go.
--
--   Where the constraint already exists, this drops and rebuilds the index
--   under ACCESS EXCLUSIVE. blk.token holds one row per tokenised asset, a
--   handful today, so the rebuild is instantaneous; and migrate wraps the
--   whole migration in a transaction, so no other session ever observes the
--   table without its constraint. Should the table ever grow past the point
--   where that lock is acceptable, the replacement is CREATE UNIQUE INDEX
--   CONCURRENTLY followed by ADD CONSTRAINT ... USING INDEX — which cannot be
--   used from here, because CONCURRENTLY refuses to run inside a transaction.

ALTER TABLE blk.token
    DROP CONSTRAINT IF EXISTS token_name_unique;
ALTER TABLE blk.token
    ADD CONSTRAINT token_name_unique UNIQUE (token_name);

COMMENT ON COLUMN blk.token.token_name IS
    'Visible name of the token, built as "Tooken <title> #<asset id>" by '
    'defineTokenName. Unique because it is the idempotency key CreateToken '
    'resolves on: two tokens sharing a name would let a replayed deployment '
    'bind an asset to the wrong contract address.';
