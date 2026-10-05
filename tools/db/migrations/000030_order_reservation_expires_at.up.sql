-- 000030 — the order's deadline is the end of a reservation, and says so
--
--   iss.issuance_orders.expires_at was read as "when the order expires". On a
--   primary order the two are the same moment: an order expires only because
--   its reservation lapsed unpaid (EXPIRED). But the column keeps its value once
--   the order is paid, and a CLOSED order with expires_at in the past reads as
--   an expired order. The secondary market (M4) will also bring a genuine order
--   validity ("good till date") that has nothing to do with a reservation; the
--   name expires_at is left free for it.
--
--   A rename, not a new column: one fact, one place. RENAME COLUMN rewrites the
--   expression of issuance_orders_awaiting_payment_ck (I4) by itself, so the
--   invariant is untouched. The comment is restated because its wording moves.
--
--   Guarded so that a database that already carries the new name replays it as
--   a no-op.

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_schema = 'iss' AND table_name = 'issuance_orders'
                 AND column_name = 'expires_at') THEN
        ALTER TABLE iss.issuance_orders RENAME COLUMN expires_at TO reservation_expires_at;
    END IF;
END
$$;

COMMENT ON COLUMN iss.issuance_orders.reservation_expires_at IS
    'End of the share reservation. Until it passes, the shares are held out of the stock; past it, an unpaid order is EXPIRED. Kept once the order is paid, as part of its history. The duration is ORDER_RESERVATION_TTL (D28, M3-3).';
