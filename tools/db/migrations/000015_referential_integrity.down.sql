-- Rollback 000015 — the asset detail tables and the order book go back to
-- naming rows nothing guarantees exist.
--
-- Dropping a foreign key never touches data, so this reverts cleanly. What it
-- cannot do is restore the orphans and the duplicates that the up migration
-- refused: if there were any, the up migration failed and nothing was applied.

DROP INDEX IF EXISTS iss.issuance_orders_user_idx;

ALTER TABLE iss.issuance_orders
    DROP CONSTRAINT IF EXISTS fk_issuance_orders_user;
ALTER TABLE iss.issuance_orders
    DROP CONSTRAINT IF EXISTS fk_issuance_orders_asset;

ALTER TABLE ass.real_estate_shares_config
    DROP CONSTRAINT IF EXISTS fk_real_estate_shares_config_real_estate;
ALTER TABLE ass.real_estate_specification
    DROP CONSTRAINT IF EXISTS fk_real_estate_specification_real_estate;

DROP INDEX IF EXISTS ass.real_estate_shares_config_real_estate_uk;
DROP INDEX IF EXISTS ass.real_estate_specification_real_estate_uk;
