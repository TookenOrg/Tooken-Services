-- Rollback 000029 — the investor ceiling goes
--
--   The column holds one number per property, and nothing else references it
--   yet: M3-7 is what will read it. Dropping it forgets the limits already
--   entered, so dump them if any property has been configured:
--
--     \copy (SELECT real_estate_id, max_shares_per_investor
--            FROM ass.real_estate_shares_config
--            WHERE max_shares_per_investor IS NOT NULL) TO 'caps.csv' CSV HEADER
--
--   What comes back with the rollback: nothing stops a single investor from
--   holding the whole building. That was the state before 2026-10-01, and it is
--   a concentration risk the chain cannot catch either once the module is gone.

ALTER TABLE ass.real_estate_shares_config
    DROP CONSTRAINT IF EXISTS real_estate_shares_config_max_shares_ck;

ALTER TABLE ass.real_estate_shares_config
    DROP COLUMN IF EXISTS max_shares_per_investor;
