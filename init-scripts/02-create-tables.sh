#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$SUBSCRIPTIONS_DB" <<-EOSQL
    CREATE TABLE IF NOT EXISTS subscriptions (
        id UUID PRIMARY KEY,
        service_name VARCHAR(100) NOT NULL,
        price INTEGER NOT NULL CHECK (price >= 0),
        user_id UUID NOT NULL,
        start_date VARCHAR(7) NOT NULL, -- Format: MM-YYYY
        end_date VARCHAR(7),             -- Format: MM-YYYY
        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

    CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
    CREATE INDEX IF NOT EXISTS idx_subscriptions_service_name ON subscriptions(service_name);

    GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO \$SUBSCRIPTIONS_USER;
EOSQL
