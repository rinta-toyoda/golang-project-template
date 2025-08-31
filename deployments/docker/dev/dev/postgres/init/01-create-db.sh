#!/bin/bash
set -e

DB="${POSTGRES_DB}"
USER="${POSTGRES_USER}"
PASS="${POSTGRES_PASSWORD}"

echo "Initialization start: USER='${USER}', DB='${DB}'"

psql -v ON_ERROR_STOP=1 --username "${USER}" <<-EOSQL
    DO \$\$
    BEGIN
        IF NOT EXISTS (
            SELECT FROM pg_database WHERE datname = '${DB}'
        ) THEN
            CREATE DATABASE "${DB}";
        END IF;
    END
    \$\$;
EOSQL

psql -v ON_ERROR_STOP=1 --username "${USER}" <<-EOSQL
    DO \$\$
    BEGIN
        IF NOT EXISTS (
            SELECT FROM pg_roles WHERE rolname = '${USER}'
        ) THEN
            CREATE ROLE "${USER}" WITH LOGIN PASSWORD '${PASS}';
        END IF;
    END
    \$\$;

    GRANT ALL PRIVILEGES ON DATABASE "${DB}" TO "${USER}";
EOSQL

echo "Database '${DB}' and User '${USER}' initialized."
