-- Read-only grants for the Copilot MCP database user.
--
-- Purpose
--   Give a least-privilege PostgreSQL user (default: copilot_readonly) SELECT-only
--   access to the Tooken schemas, so the Postgres MCP server can read the structure
--   and data of the database without being able to modify anything.
--
-- When / where to run (ONE time)
--   As an admin role (e.g. Aiven's `avnadmin`), from a network that can reach the DB
--   (a corporate firewall may block the PostgreSQL protocol — use an unrestricted
--   network such as a personal hotspot if needed). Examples:
--     psql "<admin connection URI>" -f tools/db/grants_readonly.sql
--   or paste this whole file into any SQL client (DBeaver, Aiven console, ...)
--   connected as the admin.
--
-- Prerequisites
--   The role must already exist (create it in the Aiven console → Users), e.g.:
--     copilot_readonly
--
-- Notes
--   * Replace `copilot_readonly` if your read-only role has a different name.
--   * Replace `defaultdb` if your database name differs (see your connection URI).
--   * Schema list = the application schemas (blk, usr, iss, ass, rel) + public.
--     Trim it if one of these schemas does not exist in your database.

-- Allow the role to connect to the database.
GRANT CONNECT ON DATABASE defaultdb TO copilot_readonly;

-- Read access to existing tables of the application schemas.
GRANT USAGE  ON SCHEMA blk, usr, iss, ass, rel, public TO copilot_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA blk, usr, iss, ass, rel, public TO copilot_readonly;

-- Read access to tables created in the future (by the role running this script).
ALTER DEFAULT PRIVILEGES IN SCHEMA blk, usr, iss, ass, rel, public
    GRANT SELECT ON TABLES TO copilot_readonly;
