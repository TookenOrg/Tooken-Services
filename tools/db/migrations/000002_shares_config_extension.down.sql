-- Rollback 000002

DROP INDEX IF EXISTS ass.real_estate_shares_config_compartment_ref_uk;

ALTER TABLE ass.real_estate_shares_config
    DROP CONSTRAINT IF EXISTS real_estate_shares_config_currency_ck,
    DROP CONSTRAINT IF EXISTS real_estate_shares_config_investment_range_ck,
    DROP CONSTRAINT IF EXISTS real_estate_shares_config_fees_ck;

ALTER TABLE ass.real_estate_shares_config
    DROP COLUMN IF EXISTS total_valuation,
    DROP COLUMN IF EXISTS currency_code,
    DROP COLUMN IF EXISTS compartment_ref,
    DROP COLUMN IF EXISTS min_investment,
    DROP COLUMN IF EXISTS max_investment,
    DROP COLUMN IF EXISTS entry_fee_rate,
    DROP COLUMN IF EXISTS mgmt_fee_rate,
    DROP COLUMN IF EXISTS exit_fee_rate,
    DROP COLUMN IF EXISTS updated_at;
