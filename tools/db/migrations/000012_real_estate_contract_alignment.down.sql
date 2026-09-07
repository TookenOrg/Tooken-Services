-- Reverse of 000012.
--
-- Two of the four changes cannot be undone faithfully, and the file says so
-- rather than pretending:
--   * restoring NOT NULL on description and imageurl requires the rows written
--     in the meantime to have a value, so they are backfilled with '' — which
--     is what the prototype stored anyway;
--   * narrowing title back to 50 characters would truncate. The migration
--     refuses instead of destroying data.

DO $$
DECLARE
    long_titles INT;
    long_descriptions INT;
BEGIN
    SELECT count(*) FILTER (WHERE length(title) > 50),
           count(*) FILTER (WHERE length(description) > 500)
    INTO long_titles, long_descriptions
    FROM ass.real_estate;

    IF long_titles > 0 OR long_descriptions > 0 THEN
        RAISE EXCEPTION
            '% row(s) carry a title longer than 50 characters and % a description longer than 500. Rolling back would truncate them; shorten them first if this is really wanted.',
            long_titles, long_descriptions;
    END IF;
END
$$;

ALTER TABLE ass.real_estate
    ALTER COLUMN created_at TYPE TIMESTAMP WITHOUT TIME ZONE
        USING created_at AT TIME ZONE 'UTC';

ALTER TABLE ass.real_estate
    ALTER COLUMN title TYPE VARCHAR(50),
    ALTER COLUMN description TYPE VARCHAR(500);

DROP INDEX IF EXISTS ass.real_estate_estate_type_idx;

ALTER TABLE ass.real_estate
    DROP CONSTRAINT IF EXISTS real_estate_estate_type_fkey;

UPDATE ass.real_estate
SET description = COALESCE(description, ''),
    imageurl = COALESCE(imageurl, '');

ALTER TABLE ass.real_estate
    ALTER COLUMN description SET NOT NULL,
    ALTER COLUMN imageurl SET NOT NULL;
