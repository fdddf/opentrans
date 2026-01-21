#!/bin/bash

echo "Resetting database..."

# Get script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Drop database
$SCRIPT_DIR/db-drop.sh

# Create database
$SCRIPT_DIR/db-create.sh

# Run migrations
$SCRIPT_DIR/db-migrate.sh

echo "Database reset completed!"