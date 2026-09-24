-- 000025 — usr.users.kyc_* becomes a projection of the decision history
--
-- ═══════════════════════════════════════════════════════════════════════════
-- WHAT THIS TRIGGER MUST NEVER DO
-- ═══════════════════════════════════════════════════════════════════════════
--
--   It must never write 'verified'. That word is reserved for a fact the chain
--   has confirmed: the investor is in the IdentityRegistry and holds a claim
--   signed by the trusted issuer. An administrator clicking "approve" proves
--   that a human decided — not that the registration succeeded.
--
--   Between the two there is a deployment, a registration transaction and a
--   claim, any of which can fail. Projecting 'verified' at decision time would
--   write a statement nothing can take back, and every later read — an order,
--   a mint, a listing — would believe it.
--
--   So (D13): 'approved' is what a human decided, 'verified' is what the chain
--   proved. The projection stops at 'approved'; M2-4 writes 'verified' together
--   with kyc_verified_at, after confirmation.
--
-- WHY THE LATEST ROW, NOT THE ROW BEING WRITTEN
--   Correcting an old verification must not resurrect it. The trigger always
--   recomputes from the most recent submission of that user, whichever row the
--   write touched — the same reason 000021 recomputes unconditionally: a
--   projection that only follows the last write cannot be repaired.
--
-- WHY country_code IS PROJECTED TOO
--   It is read on every listing and every order, and it has no business being
--   fetched from the chain for that. But it is a *copy*: the source is the
--   approved verification, and the on-chain value is another copy still. The
--   day one of them changes, they must be resynchronised — that is M2-6.

-- ── a) the country, projected from the approved verification ────────────────
--
-- Nullable on purpose. An account that never submitted a KYC has no country,
-- and inventing one would be worse than having none: it would travel to the
-- IdentityRegistry as a plausible, wrong nationality. Same argument as 000016
-- and 000018.
ALTER TABLE usr.users ADD COLUMN IF NOT EXISTS country_code text;

ALTER TABLE usr.users DROP CONSTRAINT IF EXISTS users_country_code_ck;
ALTER TABLE usr.users ADD CONSTRAINT users_country_code_ck
    CHECK (country_code IS NULL OR country_code ~ '^[A-Z]{2}$');

COMMENT ON COLUMN usr.users.country_code IS
    'Projection of the country of the latest approved KYC. Never written by hand, never modifiable by the user (D11). The on-chain copy lives in the shared IdentityRegistryStorage.';

-- ── b) the states the projection needs to be able to express ────────────────
--
-- 'approved' and 'revoked' are added. Widening a CHECK cannot invalidate any
-- existing row, so this is safe on a live table.
--
--   approved : a human decided, the chain has not confirmed yet (D13)
--   revoked  : granted, then withdrawn — which is not the same story as
--              'rejected', and unlike it, requires an on-chain action
ALTER TABLE usr.users DROP CONSTRAINT IF EXISTS users_kyc_status_ck;
ALTER TABLE usr.users ADD CONSTRAINT users_kyc_status_ck
    CHECK (kyc_status IN ('none', 'pending', 'approved', 'verified', 'rejected', 'expired', 'revoked'));

-- ── c) the trigger ──────────────────────────────────────────────────────────

CREATE OR REPLACE FUNCTION usr.kyc_verification_sync_user()
RETURNS TRIGGER AS $$
DECLARE
    latest usr.kyc_verification%ROWTYPE;
    target_user integer;
BEGIN
    -- On DELETE, OLD carries the user; on INSERT and UPDATE, NEW does.
    target_user := COALESCE(NEW.user_id, OLD.user_id);

    SELECT * INTO latest
    FROM usr.kyc_verification
    WHERE user_id = target_user
    ORDER BY submitted_at DESC, id DESC
    LIMIT 1;

    IF NOT FOUND THEN
        -- Every verification of this user is gone: the profile goes back to
        -- never having submitted anything, rather than keeping a state no row
        -- supports any more.
        UPDATE usr.users
        SET kyc_status     = 'none',
            kyc_expires_at = NULL,
            country_code   = NULL,
            updated_at     = now()
        WHERE id = target_user;
        RETURN NULL;
    END IF;

    UPDATE usr.users
    SET kyc_status = CASE latest.status
                         WHEN 'submitted' THEN 'pending'
                         WHEN 'approved'  THEN 'approved'
                         WHEN 'rejected'  THEN 'rejected'
                         WHEN 'revoked'   THEN 'revoked'
                     END,
        -- Only an approval carries an expiry and a country. Projecting them
        -- from a rejected or revoked row would leave a valid-looking date next
        -- to a status that grants nothing.
        kyc_expires_at = CASE WHEN latest.status = 'approved' THEN latest.expires_at ELSE NULL END,
        country_code   = CASE WHEN latest.status = 'approved' THEN latest.declared_country_code ELSE country_code END,
        updated_at     = now()
    WHERE id = target_user;

    -- kyc_verified_at is deliberately untouched: it dates the on-chain
    -- confirmation, not the decision (D13). The existing CHECK
    -- (kyc_status <> 'verified' OR kyc_verified_at IS NOT NULL) therefore stays
    -- satisfied, because M2-4 writes both in a single UPDATE.

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION usr.kyc_verification_sync_user() IS
    'Projects the latest KYC decision of a user onto usr.users. Never writes verified: that word belongs to the chain (D13).';

DROP TRIGGER IF EXISTS kyc_verification_sync_user_trg ON usr.kyc_verification;
CREATE TRIGGER kyc_verification_sync_user_trg
    AFTER INSERT OR UPDATE OR DELETE ON usr.kyc_verification
    FOR EACH ROW EXECUTE FUNCTION usr.kyc_verification_sync_user();

-- ── d) initial synchronisation ──────────────────────────────────────────────
--
-- The table is empty at this point (000024 just created it), so this changes
-- nothing today. It is written anyway because a migration must be able to run
-- on a database where rows were inserted between the two — and because a
-- projection that is only correct for writes made after its trigger existed is
-- a projection nobody can trust.
UPDATE usr.users u
SET kyc_status = CASE latest.status
                     WHEN 'submitted' THEN 'pending'
                     WHEN 'approved'  THEN 'approved'
                     WHEN 'rejected'  THEN 'rejected'
                     WHEN 'revoked'   THEN 'revoked'
                 END,
    kyc_expires_at = CASE WHEN latest.status = 'approved' THEN latest.expires_at ELSE NULL END,
    country_code   = CASE WHEN latest.status = 'approved' THEN latest.declared_country_code ELSE u.country_code END,
    updated_at     = now()
FROM (
    SELECT DISTINCT ON (user_id) user_id, status, expires_at, declared_country_code
    FROM usr.kyc_verification
    ORDER BY user_id, submitted_at DESC, id DESC
) AS latest
WHERE u.id = latest.user_id;
