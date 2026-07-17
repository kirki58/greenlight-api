#!/bin/bash
set -e

APP_DB=$(cat /run/secrets/app_db)
APP_DB_USER=$(cat /run/secrets/app_db_user)
APP_DB_PASSWORD=$(cat /run/secrets/app_db_password)


# Execute the SQL commands using psql as the superuser
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_USER" <<-EOSQL
    CREATE USER "$APP_DB_USER" WITH PASSWORD '$APP_DB_PASSWORD';
    CREATE DATABASE "$APP_DB";
    GRANT ALL PRIVILEGES ON DATABASE "$APP_DB" TO "$APP_DB_USER";
    ALTER DATABASE "$APP_DB" OWNER TO "$APP_DB_USER";
EOSQL