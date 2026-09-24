-- Rollback 000025 — the projection stops being maintained
--
--   The trigger goes, the columns stay. usr.users.kyc_status, kyc_expires_at
--   and country_code keep the last value they were given and then drift: a KYC
--   decided after this rollback will not appear there at all.
--
--   That is the dangerous part. The columns still read exactly the same, and
--   every listing and every order keeps filtering on them — on a snapshot that
--   no longer follows anything. If this rollback is run on a live database,
--   the decision history in usr.kyc_verification remains correct, and it is the
--   only thing that does.
--
-- WHY country_code IS NOT DROPPED
--   Dropping it would destroy the projected country of every investor, and
--   rebuilding it requires the verification history to still exist. The column
--   is kept, empty of maintenance rather than empty of data.
--
-- WHY THE STATUS CHECK IS NOT NARROWED BACK
--   Restoring the original list would reject any row already sitting at
--   'approved' or 'revoked', and the rollback would fail on exactly the
--   databases that used the feature. The wider CHECK accepts everything the
--   narrower one did.

DROP TRIGGER IF EXISTS kyc_verification_sync_user_trg ON usr.kyc_verification;
DROP FUNCTION IF EXISTS usr.kyc_verification_sync_user();
