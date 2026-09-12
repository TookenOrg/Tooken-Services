-- Rollback 000017 — the role name becomes a convention again.
--
-- Restores UNIQUE (address) and lifts UNIQUE (contract_name). From here on a
-- duplicate role name is accepted again, and GetContractRoleByName silently
-- resolves to whichever row was created last.
--
-- TWO THINGS ARE NOT RESTORED, ON PURPOSE:
--
--   UNIQUE (tx_hash) is not put back. It cannot be: SHARED_IRS and TREX_SUITE
--   share the transaction that deployed the suite, which is the normal state
--   of affairs once a factory has run. Re-adding it would abort this rollback
--   on any database where the anchor exists. The constraint contradicted the
--   architecture, so undoing 000017 does not resurrect it.
--
--   blk.users is not recreated. It held no rows, no foreign key referenced it
--   and no code addressed it; there is nothing to restore and no definition
--   worth guessing.
--
-- Re-adding UNIQUE (address) succeeds today only because TREX_FACTORY and
-- TREX_SUITE spell the same address in different cases. If the Go-side
-- canonicalisation has since rewritten either row to its EIP-55 form, the two
-- become the same string and this rollback fails - correctly, since the
-- constraint was false to begin with.

ALTER TABLE blk.contract_role
    DROP CONSTRAINT IF EXISTS contract_role_contract_name_key;

ALTER TABLE blk.contract_role
    ADD CONSTRAINT contract_address_key UNIQUE (address);

COMMENT ON TABLE blk.contract_role IS NULL;
COMMENT ON COLUMN blk.contract_role.contract_name IS NULL;
COMMENT ON COLUMN blk.contract_role.address IS NULL;
COMMENT ON COLUMN blk.contract_role.tx_hash IS NULL;
