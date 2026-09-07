-- Rollback 000009

DROP INDEX IF EXISTS usr.users_role_kyc_idx;
DROP INDEX IF EXISTS usr.users_email_uk;

ALTER TABLE usr.users
    DROP CONSTRAINT IF EXISTS users_role_ck,
    DROP CONSTRAINT IF EXISTS users_kyc_status_ck,
    DROP CONSTRAINT IF EXISTS users_kyc_verified_at_ck,
    DROP CONSTRAINT IF EXISTS users_kyc_expiry_ck;

ALTER TABLE usr.users
    DROP COLUMN IF EXISTS role,
    DROP COLUMN IF EXISTS kyc_status,
    DROP COLUMN IF EXISTS kyc_verified_at,
    DROP COLUMN IF EXISTS kyc_expires_at,
    DROP COLUMN IF EXISTS updated_at;
