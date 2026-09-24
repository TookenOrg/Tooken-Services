-- 000024 — a KYC is a dated decision, not a column
--
-- ═══════════════════════════════════════════════════════════════════════════
-- WHY A TABLE, WHEN THREE COLUMNS ALREADY EXIST
-- ═══════════════════════════════════════════════════════════════════════════
--
--   usr.users already carries kyc_status, kyc_verified_at and kyc_expires_at.
--   They describe a state. A KYC is not a state: it is a decision, taken by a
--   named person, on a given day, about data frozen at that instant.
--
--   With columns alone:
--     - a rejection disappears the moment the investor submits again;
--     - a renewal overwrites the one before it;
--     - "who approved this, when, on what evidence?" has no answer at all.
--
--   That last question is precisely what an AML regime requires to be
--   answerable, years later, about a decision nobody remembers taking.
--
--   So the history becomes the truth, and the three columns become a
--   projection of it — maintained by a trigger in 000025, never written by
--   hand. The same pattern as ass.real_estate.active (000021): one source,
--   one derived value, and no way for the two to disagree.
--
-- WHY THE DECLARED DATA IS DUPLICATED HERE
--   declared_full_name and declared_country_code repeat what usr.users holds.
--   That duplication is the point (D11): what was verified is a *snapshot*. If
--   the investor changes name tomorrow, this row must keep attesting what was
--   actually checked. A verification that silently follows the profile proves
--   nothing.

CREATE TABLE IF NOT EXISTS usr.kyc_verification (
    id                bigint GENERATED ALWAYS AS IDENTITY,
    user_id           integer     NOT NULL,
    status            text        NOT NULL,

    -- What was declared AND verified at that instant (D11).
    declared_full_name    text,
    declared_country_code text,

    submitted_at      timestamptz NOT NULL DEFAULT now(),
    decided_at        timestamptz,
    decided_by        integer,
    rejection_reason  text,
    expires_at        timestamptz,

    revoked_at        timestamptz,
    revoked_by        integer,
    revocation_reason text,

    created_at        timestamptz NOT NULL DEFAULT now(),
    updated_at        timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT kyc_verification_pkey PRIMARY KEY (id),

    CONSTRAINT kyc_verification_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES usr.users(id) ON DELETE RESTRICT,
    -- The decider and the revoker are platform operators, and deleting one must
    -- not erase the fact that they decided.
    CONSTRAINT kyc_verification_decided_by_fkey
        FOREIGN KEY (decided_by) REFERENCES usr.users(id) ON DELETE RESTRICT,
    CONSTRAINT kyc_verification_revoked_by_fkey
        FOREIGN KEY (revoked_by) REFERENCES usr.users(id) ON DELETE RESTRICT,

    -- 'expired' is deliberately absent. An expiry happens by the mere passing
    -- of time: nobody writes it, no trigger fires. Storing it would require a
    -- periodic sweep whose failure would leave expired KYCs looking valid —
    -- the worst possible direction for that kind of bug. It is derived instead:
    --   status = 'approved' AND (expires_at IS NULL OR expires_at > now())
    CONSTRAINT kyc_verification_status_ck
        CHECK (status IN ('submitted', 'approved', 'rejected', 'revoked')),

    -- A decision is dated and signed, or it did not happen.
    CONSTRAINT kyc_verification_decided_ck
        CHECK (status = 'submitted'
               OR (decided_at IS NOT NULL AND decided_by IS NOT NULL)),

    -- A refusal says why. Without a reason the investor cannot correct anything,
    -- and the platform cannot justify the refusal later.
    CONSTRAINT kyc_verification_rejection_reason_ck
        CHECK (status <> 'rejected' OR rejection_reason IS NOT NULL),

    -- A revocation is dated, signed and motivated: it withdraws a right that was
    -- granted, which is a heavier act than refusing to grant it.
    CONSTRAINT kyc_verification_revocation_ck
        CHECK (status <> 'revoked'
               OR (revoked_at IS NOT NULL AND revoked_by IS NOT NULL
                   AND revocation_reason IS NOT NULL)),

    -- Only what was granted can be revoked.
    CONSTRAINT kyc_verification_revoked_after_decision_ck
        CHECK (revoked_at IS NULL OR decided_at IS NOT NULL),

    -- A validity that ends before it starts is not a validity.
    CONSTRAINT kyc_verification_expiry_ck
        CHECK (expires_at IS NULL OR decided_at IS NULL OR expires_at > decided_at),

    -- An approval must carry a country: it is what will be written into the
    -- IdentityRegistry in M2-4, where a missing country becomes 0 — a valid,
    -- plausible and wrong nationality, recorded on-chain for good.
    CONSTRAINT kyc_verification_approved_country_ck
        CHECK (status <> 'approved' OR declared_country_code IS NOT NULL),

    CONSTRAINT kyc_verification_country_code_ck
        CHECK (declared_country_code IS NULL OR declared_country_code ~ '^[A-Z]{2}$')
);

-- One open request at a time. Two concurrent submissions would produce two
-- decisions on the same investor, and nothing would say which one counts.
CREATE UNIQUE INDEX IF NOT EXISTS kyc_verification_one_open_per_user_idx
    ON usr.kyc_verification (user_id) WHERE status = 'submitted';

-- Both the projection trigger and the back-office look for "the latest row of
-- this user", and they will do it on every decision.
CREATE INDEX IF NOT EXISTS kyc_verification_user_submitted_idx
    ON usr.kyc_verification (user_id, submitted_at DESC);

COMMENT ON TABLE usr.kyc_verification IS
    'History of KYC decisions. The source of truth: usr.users.kyc_* is a projection of this table, maintained by trigger (000025).';
COMMENT ON COLUMN usr.kyc_verification.declared_full_name IS
    'Name as declared and verified at that moment. Deliberately duplicated from usr.users: a past verification must keep attesting what was actually checked (D11).';
COMMENT ON COLUMN usr.kyc_verification.declared_country_code IS
    'ISO 3166-1 alpha-2, as declared and verified. Converted to numeric only when writing to the IdentityRegistry (M2-4).';
COMMENT ON COLUMN usr.kyc_verification.status IS
    'submitted | approved | rejected | revoked. Never "expired": an expiry is derived from expires_at, because nothing writes it.';
COMMENT ON COLUMN usr.kyc_verification.expires_at IS
    'End of validity. An approved row past this date is no longer valid even though its status still says approved — always evaluate both.';
