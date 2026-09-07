-- Rollback 000008

ALTER TABLE blk.token
    DROP CONSTRAINT IF EXISTS token_nb_decimal_ck;

DROP TRIGGER IF EXISTS real_estate_sync_contract_address_trg ON ass.real_estate;
DROP FUNCTION IF EXISTS ass.real_estate_sync_contract_address();

DROP INDEX IF EXISTS ass.real_estate_token_id_uk;

ALTER TABLE ass.real_estate
    DROP COLUMN IF EXISTS token_id;
