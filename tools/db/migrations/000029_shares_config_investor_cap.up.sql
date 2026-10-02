-- 000029 — how many shares one investor may hold of one property
--
--   ERC-3643 can enforce this on-chain: MaxBalanceModule refuses a transfer that
--   would take a holder past a ceiling. But the module has to be bound and set
--   inside the same deployTREXSuite call that creates the token (trap 4.6 of the
--   milestone — it refuses a token that already has balances), so the number has
--   to exist in the database before the first deployment, not after.
--
--   This migration puts the column there. Requiring it at tokenisation, and
--   passing it to the module, is M3-7.
--
-- WHY NULLABLE
--   A property is drafted long before it is tokenised, and the manager does not
--   have to answer this question on the first screen. NOT NULL would force a
--   number into every draft and, worse, force a default — and a default here is
--   a regulatory limit chosen by nobody.
--
-- WHY numeric(20,0) AND NOT integer (T12)
--   total_shares is numeric(20,0). Comparing the two is the whole purpose of the
--   column, and comparing across types is how a cast quietly rounds a limit.

ALTER TABLE ass.real_estate_shares_config
    ADD COLUMN IF NOT EXISTS max_shares_per_investor numeric(20,0);

COMMENT ON COLUMN ass.real_estate_shares_config.max_shares_per_investor IS
    'Ceiling of shares a single investor may hold of this property. NULL while the property is a draft; required before tokenisation and passed to MaxBalanceModule in M3-7 (D35). max_investment, an amount, stays unused (D27).';

ALTER TABLE ass.real_estate_shares_config
    DROP CONSTRAINT IF EXISTS real_estate_shares_config_max_shares_ck;
ALTER TABLE ass.real_estate_shares_config
    ADD CONSTRAINT real_estate_shares_config_max_shares_ck
    CHECK (max_shares_per_investor IS NULL
           OR (max_shares_per_investor > 0 AND max_shares_per_investor <= total_shares));

COMMENT ON CONSTRAINT real_estate_shares_config_max_shares_ck ON ass.real_estate_shares_config IS
    'I10: NULL, or in ]0 ; total_shares]. Zero would mean nobody may buy anything — a listing that cannot be invested in, written by accident. Above total_shares it is not a ceiling, just a number that reads like one.';
