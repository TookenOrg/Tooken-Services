-- Reset the id sequences of every emptied table, so the next insert gets id 1.
--
-- Purpose
--   TRUNCATE without RESTART IDENTITY leaves the sequences where the deleted rows
--   left them: after the 2026-09-28 purge, ass.real_estate was empty yet its next
--   id was 514. That is harmless, but it makes a fresh environment read like a used
--   one. This script brings those counters back to 1.
--
-- Safety
--   A sequence is reset only when its owning table is EMPTY. Reference tables that
--   still hold rows (ass.issuer, the status and type tables, blk.chains) keep their
--   counter, because restarting it would hand out an id that already exists and the
--   next insert would fail on the primary key.
--
--   The owning table is resolved through pg_depend, so the rule holds for both
--   flavours used in this schema: GENERATED ALWAYS AS IDENTITY columns and plain
--   DEFAULT nextval(...) columns.
--
-- How to run it
--   Paste the two statements below into Aiven PG Studio, one after the other.
--   ALTER SEQUENCE ... RESTART is transactional, so wrapping the DO block in
--   BEGIN / ROLLBACK is a genuine dry run: the counters go back to where they were.
--
--   Verified on PostgreSQL 16 against tools/db/baseline.sql on 2026-09-28.


-- Statement 1 — the reset itself.
DO $$
DECLARE
    seq        record;
    row_count  bigint;
    reset_list text := '';
    kept_list  text := '';
BEGIN
    FOR seq IN
        SELECT s.relname                       AS sequence_name,
               sn.nspname                      AS sequence_schema,
               tn.nspname || '.' || t.relname  AS table_name
        FROM pg_class s
        JOIN pg_namespace sn ON sn.oid = s.relnamespace
        JOIN pg_depend d ON d.objid = s.oid AND d.classid = 'pg_class'::regclass
        JOIN pg_class t ON t.oid = d.refobjid
        JOIN pg_namespace tn ON tn.oid = t.relnamespace
        WHERE s.relkind = 'S'
          AND d.deptype IN ('a', 'i')  -- 'a' = OWNED BY, 'i' = identity column
          AND sn.nspname IN ('ass', 'blk', 'iss', 'rel', 'usr')
        ORDER BY 3, 1
    LOOP
        EXECUTE format('SELECT count(*) FROM %s', seq.table_name) INTO row_count;

        IF row_count = 0 THEN
            EXECUTE format('ALTER SEQUENCE %I.%I RESTART WITH 1',
                           seq.sequence_schema, seq.sequence_name);
            reset_list := reset_list || '  ' || seq.table_name || E'\n';
        ELSE
            kept_list := kept_list || '  ' || seq.table_name
                                   || ' (' || row_count || E' rows)\n';
        END IF;
    END LOOP;

    RAISE NOTICE E'\nReset to 1:\n%\nLeft untouched, table not empty:\n%',
                 reset_list, kept_list;
END
$$;


-- Statement 2 — the proof: the id the next insert will get in each table.
--
-- pg_sequence_last_value answers NULL for a sequence that has never been read, and
-- a restarted sequence is back in that state — which is why a NULL here means "the
-- next id is the start value", that is 1 for every sequence of this schema.
--
-- pg_sequences has no is_called column, so reading last_value from it cannot tell a
-- fresh sequence from a used one. That is the reason this query goes through
-- pg_sequence_last_value instead.
SELECT tn.nspname || '.' || t.relname AS table_name,
       COALESCE(pg_sequence_last_value(s.oid) + 1, ps.start_value) AS next_id,
       (xpath('/row/c/text()',
              query_to_xml(format('SELECT count(*) AS c FROM %I.%I',
                                  tn.nspname, t.relname),
                           false, true, '')))[1]::text::bigint AS rows_now
FROM pg_class s
JOIN pg_namespace sn ON sn.oid = s.relnamespace
JOIN pg_depend d ON d.objid = s.oid AND d.classid = 'pg_class'::regclass
JOIN pg_class t ON t.oid = d.refobjid
JOIN pg_namespace tn ON tn.oid = t.relnamespace
JOIN pg_sequences ps ON ps.schemaname = sn.nspname AND ps.sequencename = s.relname
WHERE s.relkind = 'S'
  AND d.deptype IN ('a', 'i')
  AND sn.nspname IN ('ass', 'blk', 'iss', 'rel', 'usr')
ORDER BY 1;
