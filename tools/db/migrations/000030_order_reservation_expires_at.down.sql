-- Rollback 000030 — the column takes back its former name
--
--   Lossless: a rename moves no data, and issuance_orders_awaiting_payment_ck
--   follows the column either way. Code reading reservation_expires_at stops
--   working, so roll the application back first.

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_schema = 'iss' AND table_name = 'issuance_orders'
                 AND column_name = 'reservation_expires_at') THEN
        ALTER TABLE iss.issuance_orders RENAME COLUMN reservation_expires_at TO expires_at;
    END IF;
END
$$;

COMMENT ON COLUMN iss.issuance_orders.expires_at IS
    'End of the reservation. Until it passes, the shares are held out of the stock; the duration itself is M3-3 (D28).';
