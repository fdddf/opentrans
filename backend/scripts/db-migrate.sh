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

echo "Running database migrations..."
echo "Host: $DB_HOST"
echo "Port: $DB_PORT"
echo "Database: $DB_NAME"
echo "User: $DB_USER"

# Set PGPASSWORD to avoid password prompt
export PGPASSWORD=$DB_PASSWORD

# Construct DATABASE_URL for Go application
export DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo "DATABASE_URL: postgres://${DB_USER}:****@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Build and run migration
go build -o opentrans main.go
./opentrans migrate

echo "Migration completed!"