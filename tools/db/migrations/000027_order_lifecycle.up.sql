-- 000027 — the order lifecycle: thirteen states, and a history that cannot lie
--
-- ═══════════════════════════════════════════════════════════════════════════
-- WHY THE REFERENTIAL IS WRITTEN BY EXPLICIT ID
-- ═══════════════════════════════════════════════════════════════════════════
--
--   000010 matched statuses by code, on the stated grounds that "ids may differ
--   between environments". That was defensible while nothing but Go read them.
--
--   It is not defensible any more. I5 below, and I4/I6 in 000028, are CHECK
--   constraints, and a CHECK cannot join: it can only name an integer. From the
--   moment a constraint says "status 5 requires a delivery hash", the id is part
--   of the contract, and an environment where 5 means something else is an
--   environment where the constraint protects the wrong state — silently.
--
--   So this migration does two things in order: it proves the ids it is about to
--   rely on, then it writes them. The guard is not ceremony. It is the only
--   moment where a mismatch can still be seen.
--
-- WHY A TABLE AND NOT AN ENUM
--   A status carries two flags that the application reads on every listing
--   (is_final, counts_as_reserved). An enum carries a name and nothing else, and
--   adding a value to one is a DDL that cannot be rolled back inside a
--   transaction. The table stays.

-- ── a) the guard: the ids must mean what the constraints will assume ────────
--
-- Accepts the old code or the new one for each id, because this migration is
-- the very thing that renames them. It fails loudly, naming the id, the code it
-- found and the code it expected — a migration that stops here has saved a
-- constraint from guarding the wrong row.
DO $$
DECLARE
    expected CONSTANT jsonb := jsonb_build_object(
        '1', jsonb_build_array('CREATED'),
        '2', jsonb_build_array('RESERVED', 'AWAITING_PAYMENT'),
        '4', jsonb_build_array('PAYMENT_CONFIRMED'),
        '5', jsonb_build_array('COMPLETED', 'CLOSED'),
        '6', jsonb_build_array('CANCELLED'),
        '7', jsonb_build_array('EXPIRED')
    );
    row_id    text;
    found     text;
BEGIN
    FOR row_id IN SELECT jsonb_object_keys(expected) LOOP
        SELECT code INTO found
        FROM iss.issuance_order_statuses
        WHERE id = row_id::integer;

        -- A missing row is fine: a database built from the migrations alone has
        -- an empty referential, and part (b) is about to fill it.
        IF found IS NOT NULL AND NOT (expected -> row_id) ? found THEN
            RAISE EXCEPTION
                'iss.issuance_order_statuses id % carries code %, expected one of % — the CHECK constraints of 000027/000028 name ids, so this database would guard the wrong status',
                row_id, found, expected -> row_id;
        END IF;
    END LOOP;
END $$;

-- ── b) the target referential ───────────────────────────────────────────────
--
-- Insert *or* update: a database built from the migrations alone holds the table
-- and none of its rows, while Aiven holds the seven rows of 000010. Both must
-- end with the same thirteen.
--
-- CREATED stops reserving (T4). It is the draft of an order, not a claim on the
-- stock; only AWAITING_PAYMENT holds shares back. The milestone says so in §3,
-- and real_estate_columns.go computes tokens_sold straight from this flag, so
-- the change takes effect with no Go to write.
INSERT INTO iss.issuance_order_statuses (id, code, label, is_final, counts_as_reserved) VALUES
    (1,  'CREATED',              'Order created',        false, false),
    (2,  'AWAITING_PAYMENT',     'Awaiting payment',     false, true),
    (4,  'PAYMENT_CONFIRMED',    'Payment confirmed',    false, true),
    (5,  'CLOSED',               'Order executed',       true,  true),
    (6,  'CANCELLED',            'Order cancelled',      true,  false),
    (7,  'EXPIRED',              'Reservation expired',  true,  false),
    (8,  'DELIVERY_IN_PROGRESS', 'Delivery in progress', false, true),
    (9,  'DELIVERY_FAILED',      'Delivery failed',      false, true),
    (10, 'DELIVERED',            'Tokens delivered',     false, true),
    (11, 'REVERSAL_IN_PROGRESS', 'Reversal in progress', false, true),
    (12, 'REVERSAL_FAILED',      'Reversal failed',      false, true),
    (13, 'REFUND_PENDING',       'Refund pending',       false, false),
    (14, 'REFUNDED',             'Refunded',             true,  false)
ON CONFLICT (id) DO UPDATE SET
    code               = EXCLUDED.code,
    label              = EXCLUDED.label,
    is_final           = EXCLUDED.is_final,
    counts_as_reserved = EXCLUDED.counts_as_reserved;

COMMENT ON TABLE iss.issuance_order_statuses IS
    'The thirteen states an issuance order can hold. Ids are part of the contract since 000027: the CHECK constraints I4, I5 and I6 name them and cannot join. Id 3 (PAYMENT_PENDING) was withdrawn and is never reused.';

-- ── c) the history, before anything needs to be written into it ─────────────
--
-- Current status on the order, full history beside it (T5). Deducing the state
-- from the last event would make every listing aggregate this table; storing it
-- twice is the trade, and the two are written in one transaction.
CREATE TABLE IF NOT EXISTS iss.issuance_order_status_history (
    id             bigint GENERATED ALWAYS AS IDENTITY,
    order_id       bigint      NOT NULL,
    from_status_id integer,
    to_status_id   integer     NOT NULL,
    actor_user_id  bigint,
    reason         text,
    tx_hash        varchar(66),
    created_at     timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE iss.issuance_order_status_history
    DROP CONSTRAINT IF EXISTS issuance_order_status_history_pkey;
ALTER TABLE iss.issuance_order_status_history
    ADD CONSTRAINT issuance_order_status_history_pkey PRIMARY KEY (id);

-- RESTRICT, not CASCADE: a financial trail does not leave with the row it
-- describes. Deleting an order that has a history must be refused, and the
-- refusal is the point.
ALTER TABLE iss.issuance_order_status_history
    DROP CONSTRAINT IF EXISTS issuance_order_status_history_order_fkey;
ALTER TABLE iss.issuance_order_status_history
    ADD CONSTRAINT issuance_order_status_history_order_fkey
    FOREIGN KEY (order_id) REFERENCES iss.issuance_orders(id) ON DELETE RESTRICT;

-- from_status_id carries NO foreign key, and that is the discovery this
-- migration made about itself.
--
--   An append-only history that points at a mutable referential forbids the
--   referential from ever shrinking. This very migration withdraws status 3,
--   and the rows it writes to record the withdrawal are rows that say "came
--   from 3" — so a RESTRICT foreign key here refuses the deletion the migration
--   exists to perform. Tried on PostgreSQL 16.14 on 2026-10-01: the DELETE
--   fails on issuance_order_status_history_from_status_fkey.
--
--   The relaxation is exactly as narrow as the problem. A status being retired
--   is only ever a *departure*: nothing lands on it any more, which is what
--   retiring it means. So to_status_id keeps its foreign key — a transition has
--   to arrive somewhere that exists — and from_status_id becomes what it
--   actually is, the archived name of a state that may since have been removed.
--
--   The alternative was to keep status 3 as a row flagged "retired". T3 rejected
--   it: a status that exists but must not be used is a trap for every future
--   reader. That decision stands; this is what it costs.

ALTER TABLE iss.issuance_order_status_history
    DROP CONSTRAINT IF EXISTS issuance_order_status_history_from_status_fkey;

ALTER TABLE iss.issuance_order_status_history
    DROP CONSTRAINT IF EXISTS issuance_order_status_history_to_status_fkey;
ALTER TABLE iss.issuance_order_status_history
    ADD CONSTRAINT issuance_order_status_history_to_status_fkey
    FOREIGN KEY (to_status_id) REFERENCES iss.issuance_order_statuses(id) ON DELETE RESTRICT;

COMMENT ON COLUMN iss.issuance_order_status_history.from_status_id IS
    'The status the order left. NULL on the first row of an order. Deliberately not a foreign key: a withdrawn status (3 PAYMENT_PENDING, removed by 000027) must stay readable here after it has left the referential.';

-- NULL actor means the system acted: an expiry, a confirmed mint, this very
-- migration. RESTRICT again — the account that cancelled an order is part of
-- the record.
ALTER TABLE iss.issuance_order_status_history
    DROP CONSTRAINT IF EXISTS issuance_order_status_history_actor_fkey;
ALTER TABLE iss.issuance_order_status_history
    ADD CONSTRAINT issuance_order_status_history_actor_fkey
    FOREIGN KEY (actor_user_id) REFERENCES usr.users(id) ON DELETE RESTRICT;

ALTER TABLE iss.issuance_order_status_history
    DROP CONSTRAINT IF EXISTS issuance_order_status_history_transition_ck;
ALTER TABLE iss.issuance_order_status_history
    ADD CONSTRAINT issuance_order_status_history_transition_ck
    CHECK (from_status_id IS NULL OR from_status_id <> to_status_id);

COMMENT ON CONSTRAINT issuance_order_status_history_transition_ck
    ON iss.issuance_order_status_history IS
    'I-history: a transition goes somewhere. A row from a status to itself is a write that happened for a reason nobody recorded, and it would show up on the timeline the front renders.';

ALTER TABLE iss.issuance_order_status_history
    DROP CONSTRAINT IF EXISTS issuance_order_status_history_tx_hash_ck;
ALTER TABLE iss.issuance_order_status_history
    ADD CONSTRAINT issuance_order_status_history_tx_hash_ck
    CHECK (tx_hash IS NULL OR tx_hash ~ '^0x[0-9a-fA-F]{64}$');

COMMENT ON CONSTRAINT issuance_order_status_history_tx_hash_ck
    ON iss.issuance_order_status_history IS
    'I2: a transaction hash is 0x and 64 hex digits. A truncated hash looks like a proof and finds nothing on a block explorer.';

-- The timeline the front will render, read in one direction only.
CREATE INDEX IF NOT EXISTS issuance_order_status_history_order_created_idx
    ON iss.issuance_order_status_history (order_id, created_at);

COMMENT ON TABLE iss.issuance_order_status_history IS
    'Append-only trail of every status change of an order. Written in the same transaction as iss.issuance_orders.status_id (T5).';

-- ── d) I1 — append-only ─────────────────────────────────────────────────────
--
-- No %ROWTYPE anywhere in this function, deliberately. %ROWTYPE is resolved when
-- the body is compiled, and tools/db/extract_baseline.sql emits every function
-- ahead of every table: a function declaring one cannot be created from the
-- baseline at all (PROGRESS.md §45.2). This one needs no row variable, which is
-- the simplest way to stay out of that trap.
CREATE OR REPLACE FUNCTION iss.issuance_order_status_history_append_only()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION
        'iss.issuance_order_status_history is append-only: % is refused. Correct an order by recording the transition that corrects it.',
        TG_OP;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION iss.issuance_order_status_history_append_only() IS
    'I1: refuses UPDATE and DELETE on the order history. An exceptional operation that genuinely must erase history disables the trigger for the length of one transaction — the cost belongs to the exception, not to the invariant.';

DROP TRIGGER IF EXISTS issuance_order_status_history_append_only_trg
    ON iss.issuance_order_status_history;
CREATE TRIGGER issuance_order_status_history_append_only_trg
    BEFORE UPDATE OR DELETE ON iss.issuance_order_status_history
    FOR EACH ROW EXECUTE FUNCTION iss.issuance_order_status_history_append_only();

-- ── e) withdrawing PAYMENT_PENDING without losing anything ──────────────────
--
-- Id 3 said "waiting for the money", which is what id 2 now says. Keeping it as
-- a row flagged "do not use" would leave a trap for every future reader; moving
-- its orders costs one UPDATE.
--
-- The history row is written FIRST and is not decoration: it is the only thing
-- that lets the rollback put these exact orders back in 3, rather than every
-- order that happens to sit in 2 — including those that were already there.
INSERT INTO iss.issuance_order_status_history (order_id, from_status_id, to_status_id, actor_user_id, reason)
SELECT id, 3, 2, NULL, 'migration 000027: PAYMENT_PENDING withdrawn, merged into AWAITING_PAYMENT'
FROM iss.issuance_orders
WHERE status_id = 3;

UPDATE iss.issuance_orders SET status_id = 2, updated_at = now() WHERE status_id = 3;

DELETE FROM iss.issuance_order_statuses WHERE id = 3;

-- The sequence has to clear the ids this migration wrote by hand, or the next
-- status created without an explicit id collides with CREATED.
SELECT setval(
    pg_get_serial_sequence('iss.issuance_order_statuses', 'id'),
    (SELECT max(id) FROM iss.issuance_order_statuses),
    true
);

-- ── f) the delivery and reversal hashes ─────────────────────────────────────
--
-- They live in this migration rather than in 000028 because I5 just below needs
-- them, and I5 is a statement about the lifecycle, not about the money.
ALTER TABLE iss.issuance_orders ADD COLUMN IF NOT EXISTS delivery_tx_hash varchar(66);
ALTER TABLE iss.issuance_orders ADD COLUMN IF NOT EXISTS reversal_tx_hash varchar(66);

COMMENT ON COLUMN iss.issuance_orders.delivery_tx_hash IS
    'Hash of the mint that delivered the shares. Written before waiting for the receipt (M3-8), so a crash mid-wait leaves the proof behind.';
COMMENT ON COLUMN iss.issuance_orders.reversal_tx_hash IS
    'Hash of the burn that took the shares back (M3-9).';

ALTER TABLE iss.issuance_orders DROP CONSTRAINT IF EXISTS issuance_orders_delivery_tx_hash_ck;
ALTER TABLE iss.issuance_orders ADD CONSTRAINT issuance_orders_delivery_tx_hash_ck
    CHECK (delivery_tx_hash IS NULL OR delivery_tx_hash ~ '^0x[0-9a-fA-F]{64}$');

COMMENT ON CONSTRAINT issuance_orders_delivery_tx_hash_ck ON iss.issuance_orders IS
    'I2: a transaction hash is 0x and 64 hex digits.';

ALTER TABLE iss.issuance_orders DROP CONSTRAINT IF EXISTS issuance_orders_reversal_tx_hash_ck;
ALTER TABLE iss.issuance_orders ADD CONSTRAINT issuance_orders_reversal_tx_hash_ck
    CHECK (reversal_tx_hash IS NULL OR reversal_tx_hash ~ '^0x[0-9a-fA-F]{64}$');

COMMENT ON CONSTRAINT issuance_orders_reversal_tx_hash_ck ON iss.issuance_orders IS
    'I2: a transaction hash is 0x and 64 hex digits.';

-- One transaction delivers one order. The same hash on two orders would mean a
-- single mint was counted twice, and M3-8 relaunches a failed delivery: this
-- index is what stops a relaunch from recording the first attempt twice.
CREATE UNIQUE INDEX IF NOT EXISTS issuance_orders_delivery_tx_hash_uidx
    ON iss.issuance_orders (delivery_tx_hash) WHERE delivery_tx_hash IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS issuance_orders_reversal_tx_hash_uidx
    ON iss.issuance_orders (reversal_tx_hash) WHERE reversal_tx_hash IS NOT NULL;

-- ── g) I5 — a delivered order knows which transaction delivered it ──────────
--
-- DELIVERED, CLOSED, REVERSAL_IN_PROGRESS and REVERSAL_FAILED all assert that
-- shares reached a wallet. Without the hash, nothing on the chain backs the
-- claim, and a reversal would not know what to burn.
--
-- Posed VALID: iss.issuance_orders is empty since the purge of 2026-09-28
-- (T10 revised). NOT VALID would have exempted nothing from future writes
-- anyway — it only skips the initial scan.
--
-- The guard below exists because "is violated by some row" is what PostgreSQL
-- says otherwise, and it names neither the row nor the way out. Any database
-- still holding a COMPLETED order is in that situation by construction: the
-- hash column did not exist before this file, so no such order can ever have
-- carried one.
DO $$
DECLARE
    offenders text;
BEGIN
    SELECT string_agg(COALESCE(order_reference, 'id ' || id), ', ' ORDER BY id)
    INTO offenders
    FROM iss.issuance_orders
    WHERE status_id IN (5, 10, 11, 12) AND delivery_tx_hash IS NULL;

    IF offenders IS NOT NULL THEN
        RAISE EXCEPTION
            'cannot apply I5: order(s) % claim the shares were delivered and carry no transaction hash. They predate this migration, so none ever could. Record the real hash on each, or move them to 6 CANCELLED if the delivery never happened, then re-run.',
            offenders;
    END IF;
END $$;

ALTER TABLE iss.issuance_orders DROP CONSTRAINT IF EXISTS issuance_orders_delivered_hash_ck;
ALTER TABLE iss.issuance_orders ADD CONSTRAINT issuance_orders_delivered_hash_ck
    CHECK (status_id NOT IN (5, 10, 11, 12) OR delivery_tx_hash IS NOT NULL);

COMMENT ON CONSTRAINT issuance_orders_delivered_hash_ck ON iss.issuance_orders IS
    'I5: statuses 10 DELIVERED, 5 CLOSED, 11 REVERSAL_IN_PROGRESS and 12 REVERSAL_FAILED all claim the shares were minted. The hash is what backs the claim — and what a reversal reads to know what to burn.';
