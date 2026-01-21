#!/bin/bash

# Load .env file from project root
if [ -f "../.env" ]; then
    export $(cat ../.env | grep -v '^#' | xargs)
fi

# Set defaults if not set
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-i18n}
DB_PASSWORD=${DB_PASSWORD:-change_this_password}
DB_NAME=${DB_NAME:-i18n}

echo "Dropping database..."
echo "Host: $DB_HOST"
echo "Port: $DB_PORT"
echo "Database: $DB_NAME"

# Set PGPASSWORD to avoid password prompt
export PGPASSWORD=$DB_PASSWORD

# Drop database
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -c "DROP DATABASE IF EXISTS $DB_NAME;" 2>/dev/null || echo "Database may not exist"

echo "Database dropped successfully!"