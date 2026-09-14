-- 000018 — An asset is born incomplete, and the table must let it
--
-- WHAT WAS WRONG
--   ass.real_estate_specification demanded seven columns that the API declares
--   optional: lot_size, surface_area, bedroom_number, bathroom_number,
--   pool_size, terrace_size and built_year. All seven were NOT NULL with no
--   default, while the contract exposes every one of them as omitempty and the
--   DTO carries them as decimal.NullDecimal — a type that exists precisely to
--   represent absence.
--
--   The gap did not surface as a clean validation error. It surfaced as a 400
--   built from a PostgreSQL not_null_violation, translated by
--   real_estate_write.go:556 using pqErr.Column:
--
--       invalid real estate payload: lot_size cannot be empty in this database
--
--   Read literally, that message was accurate: the database, not the business
--   rule, was refusing. An apartment with no pool could not be recorded unless
--   its manager declared a pool size.
--
-- WHY THE TABLE IS WRONG AND NOT THE CONTRACT
--   Creation always produces a draft. A draft is incomplete by construction —
--   that is what it is for. Two moments therefore carry two different guards:
--
--     creation    the schema, which may only demand what is true at all times
--     publication publicationRequirements, which demands what an investor
--                 must be shown before an offer is put in front of them
--
--   issuer_id belongs to the first list and 000016 rightly moved it there: an
--   asset is issued by a vehicle from the moment it exists, visible or not.
--   A pool size belongs to neither. Holding a publication requirement in a
--   NOT NULL does not reinforce it, it forbids the draft.
--
--   The repository already settled this argument twice, in the same direction.
--   publicationRequirements explains why surface is not demanded even to
--   publish:
--
--       "demanding it would push a manager to invent a number to get past
--        the check"
--
--   and 000016 refused to invent an issuer:
--
--       "Attaching an asset to a vehicle picked by a migration would write a
--        plausible but false line into the register"
--
--   pool_size NOT NULL produced exactly what both texts reject: a false figure,
--   written to clear a barrier. A zero recorded for a property that was never
--   measured is not missing data, it is wrong data, and nothing downstream can
--   tell the two apart afterwards.
--
-- WHY surface_area KEEPS ITS CONSTRAINT
--   Of the seven, it is the only one true at all times: a building has a floor
--   area whether or not anyone has written it down yet. It is also the figure
--   every other number is read against — a price per share means little beside
--   an unknown surface.
--
--   Note the distinction with publicationRequirements, which deliberately does
--   NOT demand it. The two are consistent: the column must be filled for a row
--   to exist at all, and publication does not re-ask for what creation already
--   guaranteed.
--
-- WHY NO DEFAULT IS INVENTED
--   Not even zero, and least of all for pool_size and terrace_size, where zero
--   is a meaningful measurement: it says "measured, and there is none". NULL
--   says "not measured". Collapsing the two would destroy the only signal that
--   tells a manager what is left to fill in.
--
-- REPLAYING
--   DROP NOT NULL on a column that is already nullable is a no-op in
--   PostgreSQL, not an error. This migration is therefore safe to apply to a
--   database where the constraints were already lifted by hand, which is the
--   state some environments are in — and the reason it is written down here
--   rather than left as a manual step.

ALTER TABLE ass.real_estate_specification
    ALTER COLUMN lot_size        DROP NOT NULL,
    ALTER COLUMN pool_size       DROP NOT NULL,
    ALTER COLUMN terrace_size    DROP NOT NULL,
    ALTER COLUMN bedroom_number  DROP NOT NULL,
    ALTER COLUMN bathroom_number DROP NOT NULL,
    ALTER COLUMN built_year      DROP NOT NULL;

COMMENT ON COLUMN ass.real_estate_specification.lot_size IS
    'Land area. Optional: an apartment has none, and a draft may be recorded '
    'before the deed is read. NULL means unknown, 0 means measured and absent.';

COMMENT ON COLUMN ass.real_estate_specification.pool_size IS
    'Optional. NULL means unknown, 0 means there is no pool. Never defaulted: '
    'the distinction is the only thing telling a manager what is left to fill.';

COMMENT ON COLUMN ass.real_estate_specification.terrace_size IS
    'Optional, read exactly like pool_size.';

COMMENT ON COLUMN ass.real_estate_specification.surface_area IS
    'Floor area. The one specification required at all times: it is true of '
    'every building and every other figure is read against it. Publication '
    'does not re-ask for it — creation already guaranteed it.';
