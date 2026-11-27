-- Rollback of 000031.
--
-- Refuses to run when any ONCHAIN payment or refund row exists, so on-chain
-- transaction proofs are never silently dropped.

DO $$
DECLARE
    bad text;
BEGIN
    SELECT string_agg(id::text, ', ' ORDER BY id) INTO bad
    FROM iss.issuance_order_payments
    WHERE provider = 'ONCHAIN'
       OR tx_hash IS NOT NULL
       OR recorded_by_user_id IS NULL;

    IF bad IS NOT NULL THEN
        RAISE EXCEPTION
            '000031 down refused: iss.issuance_order_payments row(s) [%] use ONCHAIN mode (or carry tx_hash / NULL recorded_by_user_id). '
            'Rolling back would drop on-chain payment proofs; archive or remove those rows before rolling back.',
            bad;
    END IF;
END $$;

DROP INDEX IF EXISTS iss.issuance_order_payments_tx_hash_uidx;

ALTER TABLE iss.issuance_order_payments DROP CONSTRAINT IF EXISTS issuance_order_payments_mode_ck;
ALTER TABLE iss.issuance_order_payments DROP CONSTRAINT IF EXISTS issuance_order_payments_tx_hash_ck;
ALTER TABLE iss.issuance_order_payments DROP COLUMN IF EXISTS tx_hash;
ALTER TABLE iss.issuance_order_payments ALTER COLUMN recorded_by_user_id SET NOT NULL;

ALTER TABLE iss.issuance_order_payments DROP CONSTRAINT IF EXISTS issuance_order_payments_provider_ck;
ALTER TABLE iss.issuance_order_payments ADD CONSTRAINT issuance_order_payments_provider_ck
    CHECK (provider IN ('MANUAL'));

COMMENT ON CONSTRAINT issuance_order_payments_provider_ck ON iss.issuance_order_payments IS
    'D1: M3 records bank transfers declared by an ADMIN and nothing else. The day a real provider arrives it brings a webhook, a reconciliation and a state machine of its own — widening this CHECK is the smallest part of that work.';

COMMENT ON TABLE iss.issuance_order_payments IS
    'Money in and out for an order, declared by an ADMIN (D1). That amount equals iss.issuance_orders.amount_due is checked by M3-6, not here: a CHECK cannot read another table.';
