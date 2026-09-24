-- Rollback 000024 — the history goes, the projection stays behind
--
--   Dropping this table destroys every KYC decision ever recorded: who
--   approved, when, on what declared data, and why anything was refused. None
--   of it is recoverable from usr.users, which only ever held the current
--   state — that asymmetry is the whole reason the table was created.
--
--   Take a dump of usr.kyc_verification before running this if the decisions
--   still matter. An AML audit asks about the past, and the past is here.
--
-- WHAT IS LEFT BEHIND
--   usr.users.kyc_status and its dates keep whatever value the trigger last
--   wrote, and nothing maintains them any more. They become a frozen snapshot
--   of the moment the rollback happened — plausible, readable, and no longer
--   backed by anything. Run the rollback of 000025 first: it removes the
--   trigger that reads this table.

DROP INDEX IF EXISTS usr.kyc_verification_user_submitted_idx;
DROP INDEX IF EXISTS usr.kyc_verification_one_open_per_user_idx;
DROP TABLE IF EXISTS usr.kyc_verification;
