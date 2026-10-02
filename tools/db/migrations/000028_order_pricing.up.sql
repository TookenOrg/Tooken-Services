-- 000028 — what the investor owes, frozen at the moment they ordered
--
-- ═══════════════════════════════════════════════════════════════════════════
-- WHY THE AMOUNTS ARE COPIED ONTO THE ORDER
-- ═══════════════════════════════════════════════════════════════════════════
--
--   Today an order is a quantity and nothing else. The price is read from
--   ass.real_estate_shares_config, which the manager may change at any time —
--   so an order placed at 200 EUR a share silently becomes an order at 240 the
--   moment the listing is edited, and the investor is asked for money they
--   never agreed to.
--
--   A price is part of an offer. Once the offer is accepted it stops being a
--   reference to a configuration and becomes a number on a contract. These
--   columns are that number.
--
-- WHY numeric(20,2) FOR THE AMOUNTS AND numeric(20,8) FOR THE PRICE (T6)
--   Nobody pays 39,998 EUR. The unit price keeps the precision of its source,
--   because 0,00000001 of a share is meaningful when a building is cut into a
--   million; what is *due* must be payable, so it is rounded to the cent and
--   checked there.
--
-- WHY THE COLUMNS STAY NULLABLE
--   InsertIssuranceOrder writes a CREATED order with no amounts at all, and it
--   must keep working until M3-3 rewrites it. NOT NULL would break it today for
--   a guarantee I3 and I4 give conditionally, which is the correct shape: an
--   order is allowed to be a draft, but it is not allowed to be half-priced.

-- ── a) the frozen amounts ───────────────────────────────────────────────────

ALTER TABLE iss.issuance_orders ADD COLUMN IF NOT EXISTS unit_price      numeric(20,8);
ALTER TABLE iss.issuance_orders ADD COLUMN IF NOT EXISTS currency_code   char(3);
ALTER TABLE iss.issuance_orders ADD COLUMN IF NOT EXISTS gross_amount    numeric(20,2);
ALTER TABLE iss.issuance_orders ADD COLUMN IF NOT EXISTS fee_amount      numeric(20,2);
ALTER TABLE iss.issuance_orders ADD COLUMN IF NOT EXISTS amount_due      numeric(20,2);
ALTER TABLE iss.issuance_orders ADD COLUMN IF NOT EXISTS expires_at      timestamptz;
ALTER TABLE iss.issuance_orders ADD COLUMN IF NOT EXISTS idempotency_key text;
ALTER TABLE iss.issuance_orders ADD COLUMN IF NOT EXISTS refund_reason   text;

COMMENT ON COLUMN iss.issuance_orders.unit_price IS
    'Copy of ass.real_estate_shares_config.price_per_share at the time of the order. Same type as the source, so no precision is lost on the way in.';
COMMENT ON COLUMN iss.issuance_orders.gross_amount IS
    'round(quantity * unit_price, 2). The rounding rule belongs to the database (I3): PostgreSQL rounds a half away from zero, and M3-3 has to agree with it.';
COMMENT ON COLUMN iss.issuance_orders.fee_amount IS
    'Sum of the fee lines charged to the investor. Kept beside the lines rather than derived from them, because the order is read far more often than its breakdown.';
COMMENT ON COLUMN iss.issuance_orders.amount_due IS
    'gross_amount + fee_amount. What the ADMIN compares a declared payment against (M3-6).';
COMMENT ON COLUMN iss.issuance_orders.expires_at IS
    'End of the reservation. Until it passes, the shares are held out of the stock; the duration itself is M3-3 (D28).';
COMMENT ON COLUMN iss.issuance_orders.idempotency_key IS
    'Key sent by the client. The unique index below is the lock: two simultaneous requests carrying the same key, one inserts and the other reads back (T9).';
COMMENT ON COLUMN iss.issuance_orders.refund_reason IS
    'PLATFORM or INVESTOR_REQUEST (D36). Who caused the refund decides who bears the cost, so it is recorded, not inferred.';

-- ── b) I3 — the five amounts are all there, or none of them is ──────────────
--
-- Half-priced is the dangerous state: an order with a unit_price and no
-- amount_due looks priced to a reader and is unpayable in practice. The
-- arithmetic is checked in the same breath, because an amount_due that does not
-- match its parts is a number somebody will act on.
--
-- num_nonnulls counts the five in one expression; spelling out five IS NULL and
-- five IS NOT NULL would say the same thing and hide the rule.
ALTER TABLE iss.issuance_orders DROP CONSTRAINT IF EXISTS issuance_orders_amounts_ck;
ALTER TABLE iss.issuance_orders ADD CONSTRAINT issuance_orders_amounts_ck
    CHECK (
        num_nonnulls(unit_price, currency_code, gross_amount, fee_amount, amount_due) IN (0, 5)
        AND (unit_price IS NULL OR (
                 gross_amount = round(quantity * unit_price, 2)
             AND fee_amount >= 0
             AND amount_due = gross_amount + fee_amount
        ))
    );

COMMENT ON CONSTRAINT issuance_orders_amounts_ck ON iss.issuance_orders IS
    'I3: the five frozen amounts are all NULL (a draft) or all set (a priced order), and when set, gross = round(quantity * unit_price, 2), fee >= 0, due = gross + fee. A half-priced order reads as priced and cannot be paid.';

-- ── c) I4 — a reservation has an end, and a price ───────────────────────────
--
-- AWAITING_PAYMENT holds shares out of the stock. Without expires_at they are
-- held for ever, and the listing never sells out; without the amounts, nobody
-- can say what to pay for them.
--
-- expires_at > created_at rather than > now(): a CHECK is re-evaluated on every
-- UPDATE of the row, so now() would make an order unmodifiable the second it
-- expires — and M3-4 has to write EXPIRED onto exactly those rows.
--
-- Posed VALID: the table is empty since the purge of 2026-09-28 (T10 revised).
--
-- Which makes this migration inapplicable to any database that still holds a
-- reserved order — including the orders 000027 has just moved here out of
-- PAYMENT_PENDING, since they arrive with neither a date nor a price. That is
-- the price of T10 revised, and the guard is what makes it say so instead of
-- failing on "is violated by some row".
DO $$
DECLARE
    offenders text;
BEGIN
    SELECT string_agg(COALESCE(order_reference, 'id ' || id), ', ' ORDER BY id)
    INTO offenders
    FROM iss.issuance_orders
    WHERE status_id = 2
      AND (expires_at IS NULL OR expires_at <= created_at OR amount_due IS NULL);

    IF offenders IS NOT NULL THEN
        RAISE EXCEPTION
            'cannot apply I4: order(s) % reserve shares with no end date or no price. Orders merged out of PAYMENT_PENDING by 000027 land here. Give each an expires_at and the five frozen amounts, or move them to 7 EXPIRED, then re-run.',
            offenders;
    END IF;
END $$;

ALTER TABLE iss.issuance_orders DROP CONSTRAINT IF EXISTS issuance_orders_awaiting_payment_ck;
ALTER TABLE iss.issuance_orders ADD CONSTRAINT issuance_orders_awaiting_payment_ck
    CHECK (
        status_id <> 2
        OR (expires_at IS NOT NULL AND expires_at > created_at AND amount_due IS NOT NULL)
    );

COMMENT ON CONSTRAINT issuance_orders_awaiting_payment_ck ON iss.issuance_orders IS
    'I4: status 2 AWAITING_PAYMENT reserves shares, so it needs an end date and a price. Compared to created_at and not now(), because a CHECK re-fires on every UPDATE and M3-4 must still be able to write EXPIRED onto a row whose date has passed.';

-- ── d) I6 and I7 — a refund says why ────────────────────────────────────────

ALTER TABLE iss.issuance_orders DROP CONSTRAINT IF EXISTS issuance_orders_refund_reason_ck;
ALTER TABLE iss.issuance_orders ADD CONSTRAINT issuance_orders_refund_reason_ck
    CHECK (refund_reason IS NULL OR refund_reason IN ('PLATFORM', 'INVESTOR_REQUEST'));

COMMENT ON CONSTRAINT issuance_orders_refund_reason_ck ON iss.issuance_orders IS
    'I7: a refund reason is PLATFORM or INVESTOR_REQUEST (D36). Free text here would be a field nobody could group by when the question is how much the platform refunded of its own doing.';

ALTER TABLE iss.issuance_orders DROP CONSTRAINT IF EXISTS issuance_orders_refund_has_reason_ck;
ALTER TABLE iss.issuance_orders ADD CONSTRAINT issuance_orders_refund_has_reason_ck
    CHECK (status_id NOT IN (13, 14) OR refund_reason IS NOT NULL);

COMMENT ON CONSTRAINT issuance_orders_refund_has_reason_ck ON iss.issuance_orders IS
    'I6: statuses 13 REFUND_PENDING and 14 REFUNDED send money back. Why it went back is the question asked months later, and it has no answer unless it is written at the time.';

-- ── e) I11 — one key, one order, per investor ───────────────────────────────
--
-- Scoped to the investor and not global: two clients picking the same UUID is
-- improbable, two clients picking the same "retry-1" is not. A collision across
-- users would refuse an order that has nothing to do with the first.
CREATE UNIQUE INDEX IF NOT EXISTS issuance_orders_idempotency_uidx
    ON iss.issuance_orders (user_id, idempotency_key) WHERE idempotency_key IS NOT NULL;

COMMENT ON INDEX iss.issuance_orders_idempotency_uidx IS
    'I11: an idempotency key is used once per investor. The index is the lock itself (T9) — of two concurrent requests carrying the same key, one inserts and the other is rejected and reads the first back.';

-- ── f) the fee lines — iss.issuance_order_fees ──────────────────────────────
--
-- fee_amount on the order answers "how much", these rows answer "of what". The
-- first is read on every listing, the second only on a detail page, which is why
-- they are not the same table.
CREATE TABLE IF NOT EXISTS iss.issuance_order_fees (
    id          bigint GENERATED ALWAYS AS IDENTITY,
    order_id    bigint         NOT NULL,
    fee_code    text           NOT NULL,
    rate        numeric(7,4)   NOT NULL,
    base_amount numeric(20,2)  NOT NULL,
    amount      numeric(20,2)  NOT NULL,
    created_at  timestamptz    NOT NULL DEFAULT now()
);

ALTER TABLE iss.issuance_order_fees DROP CONSTRAINT IF EXISTS issuance_order_fees_pkey;
ALTER TABLE iss.issuance_order_fees
    ADD CONSTRAINT issuance_order_fees_pkey PRIMARY KEY (id);

ALTER TABLE iss.issuance_order_fees DROP CONSTRAINT IF EXISTS issuance_order_fees_order_fkey;
ALTER TABLE iss.issuance_order_fees
    ADD CONSTRAINT issuance_order_fees_order_fkey
    FOREIGN KEY (order_id) REFERENCES iss.issuance_orders(id) ON DELETE RESTRICT;

-- CHECK rather than a reference table (T7): 000010 chose a table for statuses
-- because adding one needs no code. Adding a fee type always needs code — the
-- formula that computes it — so the ALTER travels in the same pull request.
-- EXIT comes with M4 (D37).
ALTER TABLE iss.issuance_order_fees DROP CONSTRAINT IF EXISTS issuance_order_fees_code_ck;
ALTER TABLE iss.issuance_order_fees ADD CONSTRAINT issuance_order_fees_code_ck
    CHECK (fee_code IN ('ENTRY'));

COMMENT ON CONSTRAINT issuance_order_fees_code_ck ON iss.issuance_order_fees IS
    'T7: only entry fees exist in M3. A fee type nobody can compute is a line that shows on an invoice and reconciles with nothing.';

-- The rate is a percentage, like its source entry_fee_rate: 2.5 means 2,5 %.
-- Dividing by 100 here, once, is what keeps M3-3 from having to guess.
ALTER TABLE iss.issuance_order_fees DROP CONSTRAINT IF EXISTS issuance_order_fees_amount_ck;
ALTER TABLE iss.issuance_order_fees ADD CONSTRAINT issuance_order_fees_amount_ck
    CHECK (
        rate >= 0 AND rate <= 100
        AND base_amount >= 0
        AND amount = round(base_amount * rate / 100, 2)
    );

COMMENT ON CONSTRAINT issuance_order_fees_amount_ck ON iss.issuance_order_fees IS
    'I8: the rate is a percentage between 0 and 100, like ass.real_estate_shares_config.entry_fee_rate, and the amount is round(base * rate / 100, 2). The database arbitrates the rounding so that five writers cannot each pick their own.';

CREATE UNIQUE INDEX IF NOT EXISTS issuance_order_fees_order_code_uidx
    ON iss.issuance_order_fees (order_id, fee_code);

COMMENT ON INDEX iss.issuance_order_fees_order_code_uidx IS
    'An order pays its entry fee once. A retried write that inserted a second line would double fee_amount with nothing to show it happened.';

COMMENT ON TABLE iss.issuance_order_fees IS
    'What makes up iss.issuance_orders.fee_amount. The consistency between the two is a service rule (M3-3), not a constraint: a CHECK cannot sum another table.';

-- ── g) money in and out — iss.issuance_order_payments ───────────────────────
--
-- One table, a kind column (T8): a payment and a refund have the same shape,
-- the same provider and the same declarer, and M3-11 reads what came in and what
-- went out in the same query.
CREATE TABLE IF NOT EXISTS iss.issuance_order_payments (
    id                  bigint GENERATED ALWAYS AS IDENTITY,
    order_id            bigint        NOT NULL,
    kind                text          NOT NULL,
    provider            text          NOT NULL,
    amount              numeric(20,2) NOT NULL,
    currency_code       char(3)       NOT NULL,
    external_reference  text,
    recorded_by_user_id bigint        NOT NULL,
    idempotency_key     text,
    created_at          timestamptz   NOT NULL DEFAULT now()
);

ALTER TABLE iss.issuance_order_payments DROP CONSTRAINT IF EXISTS issuance_order_payments_pkey;
ALTER TABLE iss.issuance_order_payments
    ADD CONSTRAINT issuance_order_payments_pkey PRIMARY KEY (id);

ALTER TABLE iss.issuance_order_payments DROP CONSTRAINT IF EXISTS issuance_order_payments_order_fkey;
ALTER TABLE iss.issuance_order_payments
    ADD CONSTRAINT issuance_order_payments_order_fkey
    FOREIGN KEY (order_id) REFERENCES iss.issuance_orders(id) ON DELETE RESTRICT;

-- NOT NULL and RESTRICT: a declared payment is an act, and the account that
-- performed it is the only thing that makes it auditable.
ALTER TABLE iss.issuance_order_payments DROP CONSTRAINT IF EXISTS issuance_order_payments_recorded_by_fkey;
ALTER TABLE iss.issuance_order_payments
    ADD CONSTRAINT issuance_order_payments_recorded_by_fkey
    FOREIGN KEY (recorded_by_user_id) REFERENCES usr.users(id) ON DELETE RESTRICT;

ALTER TABLE iss.issuance_order_payments DROP CONSTRAINT IF EXISTS issuance_order_payments_kind_ck;
ALTER TABLE iss.issuance_order_payments ADD CONSTRAINT issuance_order_payments_kind_ck
    CHECK (kind IN ('PAYMENT', 'REFUND'));

COMMENT ON CONSTRAINT issuance_order_payments_kind_ck ON iss.issuance_order_payments IS
    'T8: money comes in or goes out, there is no third direction. The sign is carried by the kind and not by the amount, so no query can forget to look at it.';

ALTER TABLE iss.issuance_order_payments DROP CONSTRAINT IF EXISTS issuance_order_payments_provider_ck;
ALTER TABLE iss.issuance_order_payments ADD CONSTRAINT issuance_order_payments_provider_ck
    CHECK (provider IN ('MANUAL'));

COMMENT ON CONSTRAINT issuance_order_payments_provider_ck ON iss.issuance_order_payments IS
    'D1: M3 records bank transfers declared by an ADMIN and nothing else. The day a real provider arrives it brings a webhook, a reconciliation and a state machine of its own — widening this CHECK is the smallest part of that work.';

ALTER TABLE iss.issuance_order_payments DROP CONSTRAINT IF EXISTS issuance_order_payments_amount_ck;
ALTER TABLE iss.issuance_order_payments ADD CONSTRAINT issuance_order_payments_amount_ck
    CHECK (amount > 0);

COMMENT ON CONSTRAINT issuance_order_payments_amount_ck ON iss.issuance_order_payments IS
    'I9: a movement of zero is not a payment, and a negative one is a refund written the wrong way round — the kind column is where the direction lives.';

-- At most one of each per order. A second PAYMENT would be either a duplicate
-- declaration or a partial payment, and M3 does not do partial payments.
CREATE UNIQUE INDEX IF NOT EXISTS issuance_order_payments_one_payment_uidx
    ON iss.issuance_order_payments (order_id) WHERE kind = 'PAYMENT';
CREATE UNIQUE INDEX IF NOT EXISTS issuance_order_payments_one_refund_uidx
    ON iss.issuance_order_payments (order_id) WHERE kind = 'REFUND';

COMMENT ON INDEX iss.issuance_order_payments_one_payment_uidx IS
    'I9: one payment per order. M3 has no partial payments, so a second row is a double declaration — the kind of mistake that is only noticed when the books are closed.';

CREATE UNIQUE INDEX IF NOT EXISTS issuance_order_payments_idempotency_uidx
    ON iss.issuance_order_payments (recorded_by_user_id, idempotency_key) WHERE idempotency_key IS NOT NULL;

COMMENT ON INDEX iss.issuance_order_payments_idempotency_uidx IS
    'I11: scoped to the ADMIN who declares, exactly as the order key is scoped to the investor. A retried declaration must not record the money twice.';

COMMENT ON TABLE iss.issuance_order_payments IS
    'Money in and out for an order, declared by an ADMIN (D1). That amount equals iss.issuance_orders.amount_due is checked by M3-6, not here: a CHECK cannot read another table.';

-- ── h) iss.issuance_allocations goes ────────────────────────────────────────
--
-- Read by no Go, written by no Go, empty on Aiven, and carrying an expires_at
-- that competes with the one added above. Two places holding the same date is
-- the guarantee that one day they disagree, and no reader will know which to
-- believe. Partial delivery, the feature it was presumably drafted for, is not
-- in M3 (T11).
DROP TABLE IF EXISTS iss.issuance_allocations;
DROP SEQUENCE IF EXISTS iss.issuance_allocations_id_seq;
