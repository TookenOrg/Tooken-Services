-- Rollback 000028 — the money columns come off
--
--   Everything 000028 added holds value: what an investor agreed to pay, what
--   they were charged in fees, what was actually received. Dropping the columns
--   and the tables does not archive any of it.
--
--   So this rollback refuses to run while there is money recorded, and says how
--   to proceed on purpose. An order that keeps its status and loses its price is
--   an order nobody can honour or contest.

-- ── a) refuse rather than destroy ───────────────────────────────────────────
DO $$
DECLARE
    fees     bigint;
    payments bigint;
BEGIN
    SELECT count(*) INTO fees     FROM iss.issuance_order_fees;
    SELECT count(*) INTO payments FROM iss.issuance_order_payments;

    IF fees > 0 OR payments > 0 THEN
        RAISE EXCEPTION
            'cannot roll back 000028: % fee line(s) and % payment(s) would be destroyed. Archive them first, then empty the tables deliberately: \copy (SELECT * FROM iss.issuance_order_payments) TO ''payments.csv'' CSV HEADER',
            fees, payments;
    END IF;
END $$;

-- The frozen amounts go with the columns. There is no table to move them to —
-- they only ever existed here — so this is the warning, not a safeguard:
--
--   \copy (SELECT id, order_reference, unit_price, currency_code, gross_amount,
--                 fee_amount, amount_due, expires_at FROM iss.issuance_orders
--          WHERE amount_due IS NOT NULL) TO 'priced_orders.csv' CSV HEADER

-- ── b) iss.issuance_allocations comes back, as baseline.sql had it ──────────
--
-- Rebuilt to the letter, including the two foreign keys that say the same thing
-- under two names: fk_issuance_allocations_order and the one PostgreSQL named
-- itself. They were both there before 000028, so they are both here — a rollback
-- that tidies up is a rollback whose result does not match the dump it claims to
-- restore.
CREATE SEQUENCE IF NOT EXISTS iss.issuance_allocations_id_seq
    AS bigint INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START WITH 1 NO CYCLE;

CREATE TABLE IF NOT EXISTS iss.issuance_allocations (
    id                bigint      DEFAULT nextval('iss.issuance_allocations_id_seq'::regclass) NOT NULL,
    issuance_order_id bigint      NOT NULL,
    asset_id          bigint      NOT NULL,
    quantity          integer     NOT NULL,
    expires_at        timestamptz NOT NULL,
    created_at        timestamptz DEFAULT now() NOT NULL
);

ALTER TABLE iss.issuance_allocations DROP CONSTRAINT IF EXISTS issuance_allocations_pkey;
ALTER TABLE iss.issuance_allocations
    ADD CONSTRAINT issuance_allocations_pkey PRIMARY KEY (id);

ALTER TABLE iss.issuance_allocations DROP CONSTRAINT IF EXISTS issuance_allocations_quantity_check;
ALTER TABLE iss.issuance_allocations
    ADD CONSTRAINT issuance_allocations_quantity_check CHECK ((quantity > 0));

ALTER TABLE iss.issuance_allocations DROP CONSTRAINT IF EXISTS fk_issuance_allocations_order;
ALTER TABLE iss.issuance_allocations
    ADD CONSTRAINT fk_issuance_allocations_order
    FOREIGN KEY (issuance_order_id) REFERENCES iss.issuance_orders(id) ON DELETE CASCADE;

ALTER TABLE iss.issuance_allocations DROP CONSTRAINT IF EXISTS issuance_allocations_issuance_order_id_fkey;
ALTER TABLE iss.issuance_allocations
    ADD CONSTRAINT issuance_allocations_issuance_order_id_fkey
    FOREIGN KEY (issuance_order_id) REFERENCES iss.issuance_orders(id) ON DELETE CASCADE;

ALTER SEQUENCE iss.issuance_allocations_id_seq OWNED BY iss.issuance_allocations.id;

-- ── c) the two tables 000028 created ────────────────────────────────────────
--
-- Their indexes and constraints go with them; both were empty, part (a) made
-- sure of it.
DROP TABLE IF EXISTS iss.issuance_order_payments;
DROP TABLE IF EXISTS iss.issuance_order_fees;

-- ── d) the columns on the order ─────────────────────────────────────────────

DROP INDEX IF EXISTS iss.issuance_orders_idempotency_uidx;

ALTER TABLE iss.issuance_orders DROP CONSTRAINT IF EXISTS issuance_orders_amounts_ck;
ALTER TABLE iss.issuance_orders DROP CONSTRAINT IF EXISTS issuance_orders_awaiting_payment_ck;
ALTER TABLE iss.issuance_orders DROP CONSTRAINT IF EXISTS issuance_orders_refund_reason_ck;
ALTER TABLE iss.issuance_orders DROP CONSTRAINT IF EXISTS issuance_orders_refund_has_reason_ck;

ALTER TABLE iss.issuance_orders DROP COLUMN IF EXISTS unit_price;
ALTER TABLE iss.issuance_orders DROP COLUMN IF EXISTS currency_code;
ALTER TABLE iss.issuance_orders DROP COLUMN IF EXISTS gross_amount;
ALTER TABLE iss.issuance_orders DROP COLUMN IF EXISTS fee_amount;
ALTER TABLE iss.issuance_orders DROP COLUMN IF EXISTS amount_due;
ALTER TABLE iss.issuance_orders DROP COLUMN IF EXISTS expires_at;
ALTER TABLE iss.issuance_orders DROP COLUMN IF EXISTS idempotency_key;
ALTER TABLE iss.issuance_orders DROP COLUMN IF EXISTS refund_reason;
