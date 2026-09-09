# Database migrations

Until now the PostgreSQL schema lived **outside the repository**: there was no
way to recreate the database, to review a schema change in a PR, or to
guarantee that dev, CI and prod were aligned. This directory fixes that.

Tool: [`golang-migrate`](https://github.com/golang-migrate/migrate).

## Installation

```bash
brew install golang-migrate
# or
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Usage

`MIGRATE_URL` below is only a local alias for `DATABASE_URL`, the connection
string the application itself uses (`main.go`, `os.Getenv("DATABASE_URL")`).

The repository ships **no `.env`** — it is gitignored and only injected by the
VS Code `launch.json` when debugging. Take the connection string from the Aiven
console (service → *Overview* → *Service URI*); `?sslmode=require` is mandatory
there and is already part of the URI Aiven displays.

```bash
export MIGRATE_URL='postgres://user:pass@host:port/db?sslmode=require'

migrate -path tools/db/migrations -database "$MIGRATE_URL" up        # apply everything
migrate -path tools/db/migrations -database "$MIGRATE_URL" up 1      # apply a single one
migrate -path tools/db/migrations -database "$MIGRATE_URL" down 1    # revert the last one
migrate -path tools/db/migrations -database "$MIGRATE_URL" version   # current version
```

Quote the URL with single quotes so the shell does not interpret `?` and `&`,
and prefix the `export` with a space to keep the password out of your shell
history.

## Where the current version is stored

`golang-migrate` creates a `public.schema_migrations` table on the first run. It
holds a single row with two columns:

```sql
SELECT * FROM public.schema_migrations;

 version | dirty
---------+-------
       9 | f
```

`version` is the last migration applied — a cursor, not a list. `dirty` is set
when a migration failed halfway through.

The state therefore lives **inside the database it describes**: nothing to
synchronise, and each database (local, CI, prod) carries its own cursor. That is
why `up` always applies exactly what is missing, wherever it is run.

## Recovering from a failed migration

While `dirty` is true, `migrate` refuses to go any further, to avoid stacking
changes on top of an unknown state. Fix the SQL, then:

```bash
migrate -path tools/db/migrations -database "$MIGRATE_URL" force <previous_version>
```

⚠️ `force` runs **no SQL at all** — it only rewrites the cursor and clears
`dirty`. Use it only after manually undoing whatever the failed migration left
half-applied, otherwise the database ends up in a state nothing can describe.

## Adopting a database that was migrated by hand

If the SQL files were applied manually (copy/paste in a SQL client),
`schema_migrations` does not exist and `migrate up` would try to replay
everything from scratch. Hand control over to the tool in three steps:

```bash
# 1. Find out where the database really stands
psql "$DATABASE_URL" -f tools/db/check_state.sql

# 2. Tell golang-migrate about it (this writes the cursor, runs no SQL)
migrate -path tools/db/migrations -database "$MIGRATE_URL" force 9

# 3. Confirm
migrate -path tools/db/migrations -database "$MIGRATE_URL" version
```

`force` creates `schema_migrations` if needed. From then on the database is
tracked normally and later migrations apply with a plain `up`.

Pass the number reported as `force_this_version` by the check script: it is the
highest *contiguous* version reached, and stops at the first gap on purpose —
forcing past a hole would skip that migration forever.

If the `migrate` CLI cannot reach the database (corporate VPN, IP allowlist),
the same cursor can be written from any SQL client — this is exactly what
`force 9` does, no more, no less:

```sql
CREATE TABLE IF NOT EXISTS public.schema_migrations (
    version bigint  NOT NULL PRIMARY KEY,
    dirty   boolean NOT NULL
);
DELETE FROM public.schema_migrations;
INSERT INTO public.schema_migrations (version, dirty) VALUES (9, false);
```

## Creating a migration

```bash
migrate create -ext sql -dir tools/db/migrations -seq explicit_name
```

Every migration must ship a `.down.sql` that genuinely restores the previous
state — this is verified (see below).

> **A migration that has been committed and deployed is never modified.** Add a
> new one instead. Editing an applied migration silently desynchronises every
> database that already ran it.

## Contents

| Version | Purpose |
|---|---|
| 000001 | Exact numeric precision: `float` → `NUMERIC` on every amount and measurement |
| 000002 | `shares_config` extension: currency, derived valuation, investment bounds, fees, compartment |
| 000003 | `specification` extension: energy/GES class + guard rails on measurements |
| 000004 | New `ass.real_estate_address` table (1‑1) |
| 000005 | New `ass.real_estate_media` table (1‑N) + backfill from `imageurl` |
| 000006 | Lifecycle: `ass.real_estate_status` + `status_id`, `active` kept in sync by a trigger |
| 000007 | New `ass.issuer` table + `real_estate.issuer_id` |
| 000008 | `real_estate.token_id` → `blk.token(id)`, `contract_address` kept in sync by a trigger, `nb_decimal BETWEEN 0 AND 18` |
| 000009 | `usr.users`: `role` and `kyc_status` |
| 000010 | Reserved shares: `counts_as_reserved` per order status, replacing a random `tokens_sold` |
| 000011 | Soft delete on `ass.real_estate` (`deleted_at` + cancelled status) |
| 000012 | `ass.real_estate` realigned with the API contract |
| 000013 | Identity sequences resynchronised after seeding with explicit ids |
| 000014 | `ass.issuer.lei_code` unique (partial index, the LEI stays optional) |
| 000015 | Referential integrity: the missing foreign keys on the detail tables and the order book, one detail row per asset |

## Compatibility with the code running in production

Two columns are **deliberately kept** even though their replacement already
exists: `ass.real_estate.active` (→ `status_id`) and
`ass.real_estate.contract_address` (→ `token_id`).

They are kept consistent by triggers, in both directions. Dropping a column that
production code still reads is the surest way to break the service during a
deployment: they will be removed by a later migration, once the `database` layer
has been migrated.

## Checking the state of a database

`tools/db/check_state.sql` reports, migration by migration, what is actually
present in a given database:

```bash
psql "$DATABASE_URL" -f tools/db/check_state.sql
```

It relies on the fact that every migration leaves a uniquely named artifact
behind (a constraint or a table), so it works even on a database golang-migrate
has never touched. It also flags known leftovers from earlier drafts, and prints
the `schema_migrations` cursor when there is one.
