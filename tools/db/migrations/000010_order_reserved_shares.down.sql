-- Rollback 000010

DROP INDEX IF EXISTS iss.issuance_orders_asset_status_idx;

ALTER TABLE iss.issuance_order_statuses
    DROP COLUMN IF EXISTS counts_as_reserved;
