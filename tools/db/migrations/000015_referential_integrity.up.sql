-- 000015 — Referential integrity for the asset detail tables and the order book
--
-- WHAT WAS MISSING
--   Four columns named a row in another table without any foreign key behind
--   them:
--
--     ass.real_estate_specification.real_estate_id  -> ass.real_estate(id)
--     ass.real_estate_shares_config.real_estate_id  -> ass.real_estate(id)
--     iss.issuance_orders.asset_id                  -> ass.real_estate(id)
--     iss.issuance_orders.user_id                   -> usr.users(id)
--
--   Only ass.real_estate_address, ass.real_estate_media and
--   rel.user_asset_favorite genuinely referenced ass.real_estate. Deleting an
--   asset therefore cascaded into those three and silently abandoned the rest.
--
--   This is not theoretical. Removing the 512 test assets from the demo
--   database was measured, before the fact, to leave 1 019 rows behind:
--   specifications, share configurations and orders still naming assets that
--   no longer existed. Nothing in the schema would have reported it, and the
--   list endpoint would have kept aggregating those orders into the progress
--   bar of assets that had been deleted.
--
-- WHY THE 1-1 TABLES ALSO NEED A UNIQUE CONSTRAINT
--   real_estate_id was NOT NULL but not unique, so an asset could carry two
--   specifications or two share configurations. Both tables are LEFT JOINed
--   directly by `realEstateBaseFrom` (internal/assets_managements/database/
--   real_estate_columns.go): a second row does not just add noise, it makes
--   the asset appear twice in the list. The write path already assumes a
--   single row — UpdateRealEstate deletes and reinserts exactly one — so the
--   constraint records an invariant the code already relies on.
--
--   ass.real_estate_address got this right from the start (000004) by making
--   real_estate_id the primary key. These two tables have their own surrogate
--   id, so the invariant is expressed as a unique index instead.
--
-- ON DELETE, AND WHY IT DIFFERS PER TABLE
--   CASCADE for the two detail tables: they describe the asset and have no
--   existence without it, which is what 000004 and 000005 already chose for
--   the address and the media.
--
--   RESTRICT for the orders: an order is a financial record. It must not
--   disappear because someone removed a listing, and it must keep naming its
--   subscriber. The application never hard-deletes an asset anyway — it sets
--   status_id = 7 and deleted_at (SoftDeleteRealEstate) — so RESTRICT changes
--   nothing in normal operation. It only refuses the accident, which is
--   precisely its job. Same reasoning as the existing issuer_id and token_id
--   keys, which are RESTRICT for the same reason.
--
-- ABOUT THE COLUMN WIDTHS
--   iss.issuance_orders.asset_id and user_id are BIGINT while
--   ass.real_estate.id and usr.users.id are INTEGER. PostgreSQL accepts a
--   foreign key across those two types: they share the same btree operator
--   family, so the equality used to check the constraint is defined. The
--   widths are left as they are on purpose — narrowing a column rewrites the
--   whole table under an ACCESS EXCLUSIVE lock, and the Go layer scans these
--   into int64. The mismatch is cosmetic; the missing key was not.
--
-- IF THIS MIGRATION FAILS
--   It means the database already holds rows the constraints reject: an
--   orphan, or an asset with two specifications. That is a finding, not an
--   obstacle to work around. The queries in the section 6 of this file list
--   the offending rows.


-- ---------------------------------------------------------------------
-- 1. One specification and one share configuration per asset
-- ---------------------------------------------------------------------
-- Created as unique indexes rather than table constraints, to match 000014
-- and because nothing references these columns.

CREATE UNIQUE INDEX IF NOT EXISTS real_estate_specification_real_estate_uk
    ON ass.real_estate_specification (real_estate_id);

CREATE UNIQUE INDEX IF NOT EXISTS real_estate_shares_config_real_estate_uk
    ON ass.real_estate_shares_config (real_estate_id);


-- ---------------------------------------------------------------------
-- 2. The detail tables belong to their asset
-- ---------------------------------------------------------------------
-- DROP IF EXISTS then ADD: ADD CONSTRAINT has no IF NOT EXISTS, and this pair
-- keeps the migration replayable without a DO block.
--
-- NOT VALID makes the ADD instant — it takes the lock, records the constraint
-- and skips the scan of existing rows. VALIDATE then checks them under
-- SHARE UPDATE EXCLUSIVE, which does not block reads or writes. The tables are
-- small today; the pattern costs one line and stops mattering the day they are
-- not. A constraint that has been validated behaves exactly like one created
-- in a single step.

ALTER TABLE ass.real_estate_specification
    DROP CONSTRAINT IF EXISTS fk_real_estate_specification_real_estate;
ALTER TABLE ass.real_estate_specification
    ADD CONSTRAINT fk_real_estate_specification_real_estate
    FOREIGN KEY (real_estate_id) REFERENCES ass.real_estate(id)
    ON DELETE CASCADE NOT VALID;
ALTER TABLE ass.real_estate_specification
    VALIDATE CONSTRAINT fk_real_estate_specification_real_estate;

ALTER TABLE ass.real_estate_shares_config
    DROP CONSTRAINT IF EXISTS fk_real_estate_shares_config_real_estate;
ALTER TABLE ass.real_estate_shares_config
    ADD CONSTRAINT fk_real_estate_shares_config_real_estate
    FOREIGN KEY (real_estate_id) REFERENCES ass.real_estate(id)
    ON DELETE CASCADE NOT VALID;
ALTER TABLE ass.real_estate_shares_config
    VALIDATE CONSTRAINT fk_real_estate_shares_config_real_estate;


-- ---------------------------------------------------------------------
-- 3. An order names an asset and a subscriber that exist
-- ---------------------------------------------------------------------

ALTER TABLE iss.issuance_orders
    DROP CONSTRAINT IF EXISTS fk_issuance_orders_asset;
ALTER TABLE iss.issuance_orders
    ADD CONSTRAINT fk_issuance_orders_asset
    FOREIGN KEY (asset_id) REFERENCES ass.real_estate(id)
    ON DELETE RESTRICT NOT VALID;
ALTER TABLE iss.issuance_orders
    VALIDATE CONSTRAINT fk_issuance_orders_asset;

ALTER TABLE iss.issuance_orders
    DROP CONSTRAINT IF EXISTS fk_issuance_orders_user;
ALTER TABLE iss.issuance_orders
    ADD CONSTRAINT fk_issuance_orders_user
    FOREIGN KEY (user_id) REFERENCES usr.users(id)
    ON DELETE RESTRICT NOT VALID;
ALTER TABLE iss.issuance_orders
    VALIDATE CONSTRAINT fk_issuance_orders_user;


-- ---------------------------------------------------------------------
-- 4. Index behind the new user key
-- ---------------------------------------------------------------------
-- asset_id is already covered by issuance_orders_asset_status_idx (000010).
-- user_id was not indexed at all: RESTRICT has to scan the orders on every
-- attempt to delete a user, and "the orders of a subscriber" is a query the
-- product needs regardless of the constraint.

CREATE INDEX IF NOT EXISTS issuance_orders_user_idx
    ON iss.issuance_orders (user_id);


-- ---------------------------------------------------------------------
-- 5. What the schema now states
-- ---------------------------------------------------------------------

COMMENT ON CONSTRAINT fk_issuance_orders_asset ON iss.issuance_orders IS
    'RESTRICT, not CASCADE: an order is a financial record and must survive '
    'the removal of the listing. Assets are retired with status_id = 7 and '
    'deleted_at, never with a DELETE.';


-- ---------------------------------------------------------------------
-- 6. Diagnosing a failure -- run these by hand, they are not part of the
--    migration
-- ---------------------------------------------------------------------
-- Orphans:
--   SELECT 'specification' AS tbl, s.id FROM ass.real_estate_specification s
--   WHERE NOT EXISTS (SELECT 1 FROM ass.real_estate r WHERE r.id = s.real_estate_id)
--   UNION ALL
--   SELECT 'shares_config', c.id FROM ass.real_estate_shares_config c
--   WHERE NOT EXISTS (SELECT 1 FROM ass.real_estate r WHERE r.id = c.real_estate_id)
--   UNION ALL
--   SELECT 'order/asset', o.id FROM iss.issuance_orders o
--   WHERE NOT EXISTS (SELECT 1 FROM ass.real_estate r WHERE r.id = o.asset_id)
--   UNION ALL
--   SELECT 'order/user', o.id FROM iss.issuance_orders o
--   WHERE NOT EXISTS (SELECT 1 FROM usr.users u WHERE u.id = o.user_id);
--
-- Assets carrying more than one detail row:
--   SELECT real_estate_id, count(*) FROM ass.real_estate_specification
--   GROUP BY real_estate_id HAVING count(*) > 1;
--   SELECT real_estate_id, count(*) FROM ass.real_estate_shares_config
--   GROUP BY real_estate_id HAVING count(*) > 1;
