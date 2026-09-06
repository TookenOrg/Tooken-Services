-- 000009 — Rôles et statut KYC sur usr.users
--
-- usr.users ne porte aujourd'hui aucune notion de rôle : impossible de
-- protéger le CRUD de l'issue #28.
--
-- Deux dimensions ORTHOGONALES, volontairement séparées :
--   * role       -> ce que l'utilisateur a le DROIT de faire
--   * kyc_status -> ce qu'il a PROUVÉ de son identité
--
-- Les fondre ensemble (un rôle « investisseur non vérifié ») créerait une
-- machine à états où chaque nouveau rôle multiplierait les combinaisons.
-- Séparés, « investisseur éligible » se lit :
--     role = 'USER' AND kyc_status = 'verified'
--
-- kyc_status est une PROJECTION de la vérité on-chain
-- (IdentityRegistry.isVerified). Elle évite un appel RPC à chaque requête ;
-- la source de vérité reste la chaîne.

ALTER TABLE usr.users
    ADD COLUMN IF NOT EXISTS role            TEXT        NOT NULL DEFAULT 'USER',
    ADD COLUMN IF NOT EXISTS kyc_status      TEXT        NOT NULL DEFAULT 'none',
    ADD COLUMN IF NOT EXISTS kyc_verified_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS kyc_expires_at  TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS updated_at      TIMESTAMPTZ NOT NULL DEFAULT now();

ALTER TABLE usr.users
    ADD CONSTRAINT users_role_ck
        CHECK (role IN ('USER', 'MANAGER', 'ADMIN')),
    ADD CONSTRAINT users_kyc_status_ck
        CHECK (kyc_status IN ('none', 'pending', 'verified', 'rejected', 'expired')),
    -- Un KYC vérifié sans date de vérification est une incohérence silencieuse.
    ADD CONSTRAINT users_kyc_verified_at_ck
        CHECK (kyc_status <> 'verified' OR kyc_verified_at IS NOT NULL),
    ADD CONSTRAINT users_kyc_expiry_ck
        CHECK (kyc_expires_at IS NULL
               OR kyc_verified_at IS NULL
               OR kyc_expires_at > kyc_verified_at);

-- Un email identifie un compte : l'unicité doit être garantie par la base,
-- pas seulement par le SELECT applicatif de GetUserByEmail.
CREATE UNIQUE INDEX IF NOT EXISTS users_email_uk
    ON usr.users (lower(email));

CREATE INDEX IF NOT EXISTS users_role_kyc_idx
    ON usr.users (role, kyc_status);
