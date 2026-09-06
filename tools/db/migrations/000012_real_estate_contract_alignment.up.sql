-- 000012 — Reconcile ass.real_estate with the API contract
--
-- The table predates the contract of étape 4 and still carries the shape of
-- the first prototype. Four divergences were reproduced against a copy of the
-- production DDL, and each of them is a failure the caller cannot act on:
--
--   1. description and imageurl are NOT NULL, while the OpenAPI write request
--      lists neither as required. Creating an asset without a description made
--      the endpoint answer 500. This is not theoretical: it is what the write
--      path does for every payload that omits them.
--   2. estate_type has no foreign key. `estate_type = 999999` was accepted and
--      stored. The service returns 400 for an unknown type only because it
--      relies on a foreign key violation, so on this table an unknown type was
--      silently written instead of being refused. A registry cannot reference
--      a type that does not exist.
--   3. title is varchar(50) while the contract advertises maxLength 255. A
--      60-character title — an ordinary one for a property listing — was
--      accepted by the API and rejected by the column, again as a 500.
--   4. created_at is `timestamp without time zone` while every other date on
--      the table (updated_at, published_at, deleted_at) carries its zone. Two
--      rows written from servers on different offsets are not comparable, and
--      the driver reads the naked value as UTC whatever it really was.
--
-- The lengths chosen below are mirrored in api/openapi.yaml and validated by
-- the service, so an oversized value is refused with a message naming the
-- field instead of reaching PostgreSQL.

-- 1 - Optional in the contract, therefore nullable in the table.
--
-- The existing rows keep their empty strings: '' and NULL both render as an
-- absent description through the read DTO, and rewriting history to make them
-- NULL would touch rows this migration has no reason to touch.
ALTER TABLE ass.real_estate
    ALTER COLUMN description DROP NOT NULL,
    ALTER COLUMN imageurl DROP NOT NULL;

-- 2 - The type must exist.
--
-- If a row already points at a missing type, the constraint cannot be created,
-- and papering over it with NOT VALID would leave the registry inconsistent
-- while pretending otherwise. The migration stops and names the offending
-- values instead, so they are fixed knowingly.
DO $$
DECLARE
    orphans TEXT;
BEGIN
    SELECT string_agg(DISTINCT re.estate_type::text, ', ')
    INTO orphans
    FROM ass.real_estate re
    LEFT JOIN ass.real_estate_type t ON t.id = re.estate_type
    WHERE re.estate_type IS NOT NULL AND t.id IS NULL;

    IF orphans IS NOT NULL THEN
        RAISE EXCEPTION
            'ass.real_estate references unknown estate types: %. Insert them into ass.real_estate_type or correct the rows, then replay this migration.',
            orphans;
    END IF;
END
$$;

ALTER TABLE ass.real_estate
    DROP CONSTRAINT IF EXISTS real_estate_estate_type_fkey;

-- RESTRICT rather than CASCADE: removing a type must not remove the assets
-- that use it.
ALTER TABLE ass.real_estate
    ADD CONSTRAINT real_estate_estate_type_fkey
        FOREIGN KEY (estate_type) REFERENCES ass.real_estate_type (id)
        ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS real_estate_estate_type_idx
    ON ass.real_estate (estate_type);

-- 3 - Lengths aligned with the contract.
--
-- Widening only: no existing value can be truncated by this migration.
ALTER TABLE ass.real_estate
    ALTER COLUMN title TYPE VARCHAR(255),
    ALTER COLUMN description TYPE VARCHAR(2000);

-- 4 - One kind of timestamp for the whole table.
--
-- The naked values were produced by now() on a server running in UTC, which is
-- what Aiven and the containers use, so they are reinterpreted as UTC. Stating
-- the zone explicitly rather than relying on the session TimeZone makes the
-- result identical wherever the migration runs.
ALTER TABLE ass.real_estate
    ALTER COLUMN created_at TYPE TIMESTAMPTZ
        USING created_at AT TIME ZONE 'UTC';

COMMENT ON COLUMN ass.real_estate.created_at IS
    'Creation instant. Timestamptz like every other date of the table: a registry compares instants, not wall clocks.';

COMMENT ON COLUMN ass.real_estate.imageurl IS
    'Legacy cover image. The gallery of ass.real_estate_media is the model going forward; this column is kept until the front stops reading it.';
