-- Rollback 000027 — the lifecycle goes back to seven states
--
--   The hard part is not dropping things. It is putting back exactly the orders
--   that were moved, and nothing else: the up merged PAYMENT_PENDING into
--   AWAITING_PAYMENT, so after it, status 2 holds both the orders that were
--   already there and the ones that arrived. The only thing that can still tell
--   them apart is the history the up wrote — which is why it is read here,
--   before the table that holds it is dropped.
--
-- WHAT THIS ROLLBACK DESTROYS
--   Every row of iss.issuance_order_status_history, including transitions
--   recorded by the application after the migration ran. That is unavoidable —
--   the table is what 000027 created — but it is a financial trail, so dump it
--   before running this if the platform has been live:
--
--     \copy (SELECT * FROM iss.issuance_order_status_history) TO 'history.csv' CSV HEADER
--
--   The delivery and reversal hashes go with it. An order that was delivered
--   keeps its status but loses the only proof of which transaction delivered it.

-- ── a) refuse rather than destroy ───────────────────────────────────────────
--
-- Statuses 8 to 14 do not exist before 000027. An order sitting in one of them
-- has nowhere to go, and inventing a destination would quietly rewrite its
-- history. Stopping here leaves the database usable; guessing would not.
DO $$
DECLARE
    stuck text;
BEGIN
    SELECT string_agg(format('order %s is in status %s', id, status_id), ', ' ORDER BY id)
    INTO stuck
    FROM iss.issuance_orders
    WHERE status_id BETWEEN 8 AND 14;

    IF stuck IS NOT NULL THEN
        RAISE EXCEPTION
            'cannot roll back 000027: % — these statuses do not exist before this migration. Move them to 1, 2, 4, 5, 6 or 7 deliberately, then re-run.',
            stuck;
    END IF;
END $$;

-- ── b) PAYMENT_PENDING comes back, and only its orders come back to it ──────

INSERT INTO iss.issuance_order_statuses (id, code, label, is_final, counts_as_reserved)
VALUES (3, 'PAYMENT_PENDING', 'Payment pending', false, true)
ON CONFLICT (id) DO UPDATE SET
    code               = EXCLUDED.code,
    label              = EXCLUDED.label,
    is_final           = EXCLUDED.is_final,
    counts_as_reserved = EXCLUDED.counts_as_reserved;

-- Named by the reason the up wrote, not by "every order in status 2": the
-- difference is the whole point of having recorded the move.
UPDATE iss.issuance_orders o
SET status_id = 3, updated_at = now()
FROM iss.issuance_order_status_history h
WHERE h.order_id = o.id
  AND h.from_status_id = 3
  AND h.to_status_id = 2
  AND h.reason = 'migration 000027: PAYMENT_PENDING withdrawn, merged into AWAITING_PAYMENT'
  AND o.status_id = 2;

-- ── c) the history table, read above, goes now ──────────────────────────────
--
-- DROP TABLE is not a row-level DELETE, so the append-only trigger does not
-- stand in the way — and it disappears with the table it guarded.
DROP TABLE IF EXISTS iss.issuance_order_status_history;
DROP FUNCTION IF EXISTS iss.issuance_order_status_history_append_only();

-- ── d) the seven original rows, as 000010 left them ─────────────────────────
--
-- CREATED goes back to reserving (T4 reversed): before 000027, a draft order
-- held shares out of the stock, and tokens_sold is computed from this flag.
UPDATE iss.issuance_order_statuses SET code = 'CREATED',   label = 'Order created',       is_final = false, counts_as_reserved = true  WHERE id = 1;
UPDATE iss.issuance_order_statuses SET code = 'RESERVED',  label = 'Tokens reserved',     is_final = false, counts_as_reserved = true  WHERE id = 2;
UPDATE iss.issuance_order_statuses SET code = 'COMPLETED', label = 'Tokens delivered',    is_final = true,  counts_as_reserved = true  WHERE id = 5;

DELETE FROM iss.issuance_order_statuses WHERE id BETWEEN 8 AND 14;

SELECT setval(
    pg_get_serial_sequence('iss.issuance_order_statuses', 'id'),
    (SELECT max(id) FROM iss.issuance_order_statuses),
    true
);

COMMENT ON TABLE iss.issuance_order_statuses IS NULL;

-- ── e) the hashes and the invariant that needed them ────────────────────────

ALTER TABLE iss.issuance_orders DROP CONSTRAINT IF EXISTS issuance_orders_delivered_hash_ck;

DROP INDEX IF EXISTS iss.issuance_orders_delivery_tx_hash_uidx;
DROP INDEX IF EXISTS iss.issuance_orders_reversal_tx_hash_uidx;

ALTER TABLE iss.issuance_orders DROP CONSTRAINT IF EXISTS issuance_orders_delivery_tx_hash_ck;
ALTER TABLE iss.issuance_orders DROP CONSTRAINT IF EXISTS issuance_orders_reversal_tx_hash_ck;

ALTER TABLE iss.issuance_orders DROP COLUMN IF EXISTS delivery_tx_hash;
ALTER TABLE iss.issuance_orders DROP COLUMN IF EXISTS reversal_tx_hash;
