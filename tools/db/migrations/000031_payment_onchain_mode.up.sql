-- 000031: add ONCHAIN payment mode to iss.issuance_order_payments.
--
-- Manual bank transfers are declared by an ADMIN (recorded_by_user_id set,
-- tx_hash NULL). On-chain stablecoin payments and refunds are observed from
-- smart-contract events by the background watcher (recorded_by_user_id NULL,
-- tx_hash set and unique across the table).

ALTER TABLE iss.issuance_order_payments DROP CONSTRAINT IF EXISTS issuance_order_payments_provider_ck;
ALTER TABLE iss.issuance_order_payments ADD CONSTRAINT issuance_order_payments_provider_ck
    CHECK (provider IN ('MANUAL', 'ONCHAIN'));

ALTER TABLE iss.issuance_order_payments ALTER COLUMN recorded_by_user_id DROP NOT NULL;

ALTER TABLE iss.issuance_order_payments ADD COLUMN IF NOT EXISTS tx_hash varchar(66);

ALTER TABLE iss.issuance_order_payments DROP CONSTRAINT IF EXISTS issuance_order_payments_tx_hash_ck;
ALTER TABLE iss.issuance_order_payments ADD CONSTRAINT issuance_order_payments_tx_hash_ck
    CHECK (tx_hash IS NULL OR tx_hash ~ '^0x[0-9a-fA-F]{64}$');

-- PostgreSQL evaluates CHECK constraints in alphabetical order, so
-- issuance_order_payments_mode_ck runs before issuance_order_payments_provider_ck.
-- Writing the invariant as two implications keeps an unknown provider TRUE here
-- so that issuance_order_payments_provider_ck is the constraint that rejects it.
ALTER TABLE iss.issuance_order_payments DROP CONSTRAINT IF EXISTS issuance_order_payments_mode_ck;
ALTER TABLE iss.issuance_order_payments ADD CONSTRAINT issuance_order_payments_mode_ck
    CHECK (
        (provider <> 'MANUAL' OR (tx_hash IS NULL AND recorded_by_user_id IS NOT NULL))
        AND (provider <> 'ONCHAIN' OR (tx_hash IS NOT NULL AND recorded_by_user_id IS NULL))
    );

-- One on-chain transaction backs at most one payment or refund row; a payment
-- and its refund are always distinct transactions on-chain, and this partial
-- unique index makes replayed watcher events idempotent.
CREATE UNIQUE INDEX IF NOT EXISTS issuance_order_payments_tx_hash_uidx
    ON iss.issuance_order_payments (tx_hash) WHERE tx_hash IS NOT NULL;

COMMENT ON TABLE iss.issuance_order_payments IS
    'Money in and out for an order: either a bank transfer declared by an ADMIN (MANUAL) or a stablecoin transfer observed on-chain (ONCHAIN). That amount equals iss.issuance_orders.amount_due is checked by the application service, not here: a CHECK cannot read another table.';

COMMENT ON COLUMN iss.issuance_order_payments.tx_hash IS
    'On-chain transaction hash of an ONCHAIN payment or refund. NULL for MANUAL bank transfers, which carry no blockchain transaction.';

COMMENT ON CONSTRAINT issuance_order_payments_provider_ck ON iss.issuance_order_payments IS
    'A payment or refund is either declared off-chain by an ADMIN (MANUAL) or observed on-chain from a payment contract event (ONCHAIN).';

COMMENT ON CONSTRAINT issuance_order_payments_tx_hash_ck ON iss.issuance_order_payments IS
    'A transaction hash, when present, is 0x followed by 64 hexadecimal digits.';

COMMENT ON CONSTRAINT issuance_order_payments_mode_ck ON iss.issuance_order_payments IS
    'A MANUAL row carries the ADMIN who recorded it and no tx_hash; an ONCHAIN row carries its tx_hash and no human actor. Written as two implications so an unknown provider falls through alphabetically to issuance_order_payments_provider_ck.';
