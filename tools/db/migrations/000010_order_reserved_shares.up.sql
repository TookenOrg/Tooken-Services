-- 000010 — Reserved shares accounting (debt 🔴: random tokens_sold)
--
-- `baseRealEstateQuery` currently computes:
--     FLOOR(random() * 2001)::int AS tokens_sold
-- that is, a DIFFERENT value on every call. The progress bar of an asset moves
-- between two page refreshes.
--
-- The real source already exists: iss.issuance_orders (quantity, status_id).
--
-- Decision: a share counts as taken FROM THE MOMENT THE ORDER IS CREATED, and
-- stays taken until the reservation is cancelled or expires.
--
-- Rather than encoding that rule in the query (where every new status would
-- force a SQL edit, in a file nobody reviews), it is carried by the reference
-- table: one flag per status. Adding a state to the flow then becomes an
-- INSERT, not a query change.
--
-- The flag CANNOT be derived from `is_final`: COMPLETED is final and holds the
-- share, CANCELLED is final and releases it. These are two distinct notions.

ALTER TABLE iss.issuance_order_statuses
    ADD COLUMN IF NOT EXISTS counts_as_reserved BOOLEAN NOT NULL DEFAULT TRUE;

COMMENT ON COLUMN iss.issuance_order_statuses.counts_as_reserved IS
    'Does an order in this status hold the ordered shares? '
    'FALSE only for the states that release them (cancellation, expiry, '
    'rejection, failure). Every new status must set this flag explicitly.';

-- Classification of the existing reference table (7 statuses):
--
--   1 CREATED           holds     -- the share is taken from creation
--   2 RESERVED          holds
--   3 PAYMENT_PENDING   holds
--   4 PAYMENT_CONFIRMED holds
--   5 COMPLETED         holds     -- shares delivered, permanently out of stock
--   6 CANCELLED         releases
--   7 EXPIRED           releases  -- the reservation falls back into stock
--
-- DEFAULT TRUE covers the first five; only the two that release are named.
-- Matching on `code` rather than on numeric ids keeps this migration correct
-- across environments, where the ids may differ.
--
-- Any cancelling status added later (REFUNDED, FAILED, VOIDED…) must set the
-- flag explicitly — the column comment above states that rule.
UPDATE iss.issuance_order_statuses
SET counts_as_reserved = FALSE
WHERE code IN ('CANCELLED', 'EXPIRED');


-- The tokens_sold computation aggregates on (asset_id, status_id): without this
-- index, every render of the asset list triggers a full scan of the orders
-- table.
CREATE INDEX IF NOT EXISTS issuance_orders_asset_status_idx
    ON iss.issuance_orders (asset_id, status_id);
