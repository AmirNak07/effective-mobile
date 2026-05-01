DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'subscriptions') THEN
      CREATE USER subscriptions WITH PASSWORD 'subscriptions_password' NOCREATEDB NOCREATEROLE;
   END IF;
END
$$;

SELECT 'CREATE DATABASE subscriptions OWNER admin'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'subscriptions')\gexec

\c subscriptions

ALTER SCHEMA public OWNER TO admin;

REVOKE ALL ON SCHEMA public FROM PUBLIC;

GRANT USAGE ON SCHEMA public TO subscriptions;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO subscriptions;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO subscriptions;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO subscriptions;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO subscriptions;

REVOKE CREATE ON SCHEMA public FROM subscriptions;

DO $$
DECLARE
    tbl text;
BEGIN
    FOR tbl IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
        EXECUTE format('GRANT SELECT, INSERT, UPDATE, DELETE ON %I TO subscriptions', tbl);
    END LOOP;
END;
$$;