#!/bin/sh
set -eu

: "${API_DB_USER:?API_DB_USER is required}"
: "${API_DB_PASSWORD:?API_DB_PASSWORD is required}"
: "${MIGRATION_DB_USER:?MIGRATION_DB_USER is required}"
: "${MIGRATION_DB_PASSWORD:?MIGRATION_DB_PASSWORD is required}"

psql --set=ON_ERROR_STOP=1 \
  --username "$POSTGRES_USER" \
  --dbname "$POSTGRES_DB" \
  --set=api_user="$API_DB_USER" \
  --set=api_password="$API_DB_PASSWORD" \
  --set=migration_user="$MIGRATION_DB_USER" \
  --set=migration_password="$MIGRATION_DB_PASSWORD" \
  --set=database_name="$POSTGRES_DB" <<'SQL'
CREATE ROLE :"api_user" LOGIN PASSWORD :'api_password';
CREATE ROLE :"migration_user" LOGIN PASSWORD :'migration_password';
GRANT CONNECT ON DATABASE :"database_name" TO :"api_user";
GRANT CONNECT ON DATABASE :"database_name" TO :"migration_user";
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
ALTER SCHEMA public OWNER TO :"migration_user";
GRANT USAGE, CREATE ON SCHEMA public TO :"migration_user";
GRANT USAGE ON SCHEMA public TO :"api_user";
SQL
