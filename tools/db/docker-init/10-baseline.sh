#!/bin/sh
# Builds the schema of a fresh local database, the way the integration tests
# build theirs: tools/db/baseline.sql is the current schema, reference data
# included. The postgres image runs this once, when its data volume is empty.
#
# The golang-migrate cursor is then set to the latest migration, as
# `migrate force` would, so that a migration added later applies on top of
# this database with a plain `migrate up` instead of replaying from 000001.
set -eu

db_dir=/tooken/db

latest=$(ls "$db_dir"/migrations/*.up.sql | sed 's#.*/\([0-9]*\)_.*#\1#' | sort | tail -n 1)

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
    --single-transaction --quiet -f "$db_dir/baseline.sql"

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
    --single-transaction --quiet -v version="$latest" <<'SQL'
CREATE TABLE IF NOT EXISTS public.schema_migrations (
    version bigint  NOT NULL PRIMARY KEY,
    dirty   boolean NOT NULL
);
DELETE FROM public.schema_migrations;
INSERT INTO public.schema_migrations (version, dirty) VALUES (:'version'::bigint, false);
SQL

echo "tooken: schema built from baseline.sql, migration cursor set to $latest"
