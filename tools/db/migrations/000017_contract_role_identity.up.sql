-- 000017 — A role is identified by its name, not by its address
--
-- WHAT THIS FIXES
--   blk.contract_role is the platform directory: it answers "which address
--   currently plays role X" (TREX_FACTORY, CLAIM_ISSUER, SHARED_IRS...). Every
--   single reader asks it that way — GetContractRoleByName does
--   WHERE contract_name = $1 (contract_role_db.go:69), and there is no other
--   lookup. The name IS the business key.
--
--   Yet the table protected everything EXCEPT the name:
--
--     contract_address_key   UNIQUE (address)   <- wrong key
--     contract_tx_hash_key   UNIQUE (tx_hash)   <- wrong key, already dropped
--     (nothing on contract_name)                <- the key that matters
--
--   Both of those guards are not merely useless here, they are false. They
--   assert a one-to-one mapping that a factory architecture breaks by design:
--
--   * UNIQUE (tx_hash) — ONE transaction deploys SIX contracts. The call to
--     TREXFactory.DeployTREXSuite creates the token, IR, IRS, TIR, CTR and MC
--     in a single transaction, and the code stores that same suite hash for
--     the role it registers (TREX_factory.go:297 passes txDeploySuite.Hash()).
--     Registering a second role from the same suite was therefore impossible.
--     This is what blocked SHARED_IRS from ever being inserted: its only
--     legitimate hash was already held by the TREX_SUITE row. Antony dropped
--     this constraint by hand to unblock the anchor; this migration records
--     that decision so every other environment gets the same schema.
--
--   * UNIQUE (address) — ONE deployed contract can legitimately hold SEVERAL
--     roles. A T-REX "suite" has no address of its own: it has six. So
--     deploySuiteAddr = deploymentDetails.Raw.Address (TREX_factory.go:288)
--     resolves to the emitter of the event, i.e. the factory. That is why
--     TREX_FACTORY and TREX_SUITE legitimately carry the same address.
--
--   Worse, that constraint never actually held. Addresses are stored as
--   varchar and PostgreSQL compares them byte for byte, while the two rows
--   differ only in case:
--
--     TREX_FACTORY  0x7f501340e10577e7b69eb8b03c0e1b9e40bffc0f   (lowercase)
--     TREX_SUITE    0x7F501340e10577e7b69Eb8b03C0e1B9e40bFFC0f   (EIP-55)
--
--   The CHECK ^0x[0-9a-fA-F]{40}$ accepts both, so the duplicate walked in
--   through the front door. A UNIQUE that any change of case defeats is worse
--   than no UNIQUE at all: it buys confidence it does not deliver.
--
-- WHAT IT ADDS INSTEAD
--   UNIQUE (contract_name) — the guard the readers actually depend on.
--   Without it, GetContractRoleByName falls back on
--   ORDER BY created_at DESC LIMIT 1, which silently means "the most recent
--   row wearing that name". A duplicate SHARED_IRS would not raise anything:
--   it would quietly re-point the shared investor whitelist, and every token
--   registered afterwards would verify its holders against a different KYC
--   list. Nothing would fail — transfers would simply start being rejected
--   for investors everyone believed were onboarded. That is the exact class
--   of bug the schema must answer for, because it never shows up in a log.
--
-- WHAT IT DELIBERATELY DOES NOT DO
--   No data is inserted here. The SHARED_IRS row names a contract deployed on
--   one specific chain; hardcoding a Sepolia address into a migration would
--   write a fact about one environment into the history of all of them. The
--   code already creates that anchor where it is missing (token.go:94-98), so
--   a fresh environment bootstraps itself correctly.
--
--   The lowercase TREX_FACTORY address is left exactly as it is. EIP-55
--   checksumming needs keccak256, which SQL cannot compute, and the value is
--   harmless: common.HexToAddress parses either case. Normalisation belongs on
--   the Go side, where the readers now canonicalise what they are given before
--   querying.
--
-- LOCKING
--   ADD CONSTRAINT ... UNIQUE builds an index under ACCESS EXCLUSIVE. Seven
--   rows: instantaneous. On a large table the same result is reached without
--   holding the table by doing CREATE UNIQUE INDEX CONCURRENTLY first, then
--   ADD CONSTRAINT ... USING INDEX. Not worth it here.

ALTER TABLE blk.contract_role
    DROP CONSTRAINT IF EXISTS contract_tx_hash_key;

ALTER TABLE blk.contract_role
    DROP CONSTRAINT IF EXISTS contract_address_key;

ALTER TABLE blk.contract_role
    ADD CONSTRAINT contract_role_contract_name_key UNIQUE (contract_name);

COMMENT ON TABLE blk.contract_role IS
    'Directory of the platform singletons: which address plays which role '
    '(TREX_FACTORY, CLAIM_ISSUER, TRANSFER_RESTRICTION_MODULE, SHARED_IRS). '
    'A designation, not an inventory - contract_instance records what was '
    'deployed. contract_name is the business key and the only lookup path.';

COMMENT ON COLUMN blk.contract_role.contract_name IS
    'The role played, not the contract type. Unique: every reader resolves a '
    'role through this column alone. Mind the distinction with the '
    'implementation key - IDENTITY_REGISTRY_STORAGE names the deployed logic '
    'in contract_implementation, while SHARED_IRS names the role of being the '
    'one investor whitelist shared by every token.';

COMMENT ON COLUMN blk.contract_role.address IS
    'Not unique: a T-REX suite has no address of its own, so the same contract '
    'legitimately holds several roles. Canonical form is EIP-55 checksummed, '
    'the form common.Address.Hex() produces; readers must canonicalise any '
    'address they receive before comparing, since this is a case-sensitive '
    'varchar.';

COMMENT ON COLUMN blk.contract_role.tx_hash IS
    'Not unique: one factory transaction deploys six contracts, so several '
    'roles legitimately share the transaction that created them.';

-- blk.users was an empty leftover from an earlier split of the schema: zero
-- rows against the 25 carried by usr.users, no foreign key pointing at it, and
-- no Go code reading or writing it. Two tables named users in one database is
-- an invitation to query the wrong one. Antony dropped it by hand; recorded
-- here so no environment recreates it.
DROP TABLE IF EXISTS blk.users;
