-- 000023 — the platform stops holding investor private keys, anywhere
--
-- ═══════════════════════════════════════════════════════════════════════════
-- WHY THE COLUMNS GO AWAY INSTEAD OF BEING RENAMED
-- ═══════════════════════════════════════════════════════════════════════════
--
--   blk.user_wallet was built to store the private key of a wallet generated
--   for an investor: in clear (private_key_clear) and encrypted beside it
--   (private_key_encrypted). Storing both cancels the second: the whole point
--   of encrypting is lost when the original sits in the next column.
--
--   The first fix was to stop writing the clear one. The second was to merge
--   the two. Both kept a place to put a key, and a column that exists for a
--   key will eventually receive one. This migration removes the place.
--
-- WHY THE PLATFORM DOES NOT NEED THEM
--   Every operation Tooken performs on a holder's tokens is signed by the
--   token agent, never by the holder:
--
--     mint(_to, _amount)                   recipient is a parameter
--     burn(_userAddress, _amount)          victim is a parameter
--     forcedTransfer(_from, _to, _amount)  source AND destination are parameters
--
--   Only transfer(_to, _amount) is signed by the holder, because its source is
--   msg.sender — and that is the one case where the investor holds his own key.
--
--   Recovery is covered too. Each investor's ONCHAINID is deployed with the
--   platform as its management key (identity.go:99 passes auth.From), so the
--   platform can addKey a wallet the investor brings later and call
--   recoveryAddress to move the holding onto it.
--
--   Which is the sentence worth remembering: the wallet key is not the
--   ownership, the ONCHAINID is — and the platform controls that. A wallet
--   nobody can sign for is not a lost holding.
--
--   Keeping investor keys therefore added a serious risk and no capability at
--   all. The real secret to protect is PRIVATE_KEY, agent of every token and
--   management key of every identity — not a database column.
--
-- ═══════════════════════════════════════════════════════════════════════════
-- WHY THREE DROPS FOR WHAT LOOKS LIKE ONE COLUMN
-- ═══════════════════════════════════════════════════════════════════════════
--
--   This migration does not start from a known state. The column history was
--   partly played by hand, so the databases disagree:
--
--     local test database  → private_key_clear, private_key_encrypted
--     production / baseline → private_key
--
--   A migration that handled only one shape would fail on the other. The three
--   DROP ... IF EXISTS below converge every known state onto the same schema,
--   and make the migration replayable: on a database already converged they
--   emit three notices and change nothing.
--
--   This is also what repairs the divergence itself. The schema change existed
--   in one database and in tools/db/baseline.sql, in no migration — so a
--   database that was not rebuilt never received it. What is not in the
--   repository does not exist; this file is where it starts existing.
--
-- NOTHING TO REGENERATE
--   tools/db/baseline.sql still carries private_key. That is fine: a rebuild
--   replays the baseline and then this migration, which drops it. If the
--   baseline is regenerated later the column will already be gone and this
--   migration becomes a no-op. Both orders produce the same schema, which is
--   the property a migration should have.
--
-- ⚠️  THIS DESTROYS THE STORED KEYS, AND THAT IS THE POINT
--   Any key still present is lost when this runs. It is not recoverable from
--   the database and it is not meant to be. If a wallet holds tokens and its
--   key disappears, the holding moves with recoveryAddress; the key itself is
--   never needed again.

ALTER TABLE blk.user_wallet DROP COLUMN IF EXISTS private_key;
ALTER TABLE blk.user_wallet DROP COLUMN IF EXISTS private_key_encrypted;
ALTER TABLE blk.user_wallet DROP COLUMN IF EXISTS private_key_clear;

COMMENT ON TABLE blk.user_wallet IS
    'Delivery addresses for investors. Holds no private key and must never '
    'hold one: every token operation the platform performs is signed by the '
    'agent, and a wallet whose key is gone is recovered through the '
    'investor''s ONCHAINID, not through a stored secret.';

COMMENT ON COLUMN blk.user_wallet.wallet_address IS
    'The wallet address — the last 20 bytes of the keccak256 of the public '
    'key, not the public key itself. The public key is never stored: when it '
    'is needed it is recovered from a signature (identity_claim.go:142).';
