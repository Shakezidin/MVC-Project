#!/usr/bin/env bash
set -euo pipefail

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-bankuser}"
DB_PASSWORD="${DB_PASSWORD:-bankpass}"
DB_NAME="${DB_NAME:-bankdb}"

export PGPASSWORD="$DB_PASSWORD"
MIGRATIONS_DIR="$(cd "$(dirname "$0")/../migrations" && pwd)"

echo "Running migrations against ${DB_HOST}:${DB_PORT}/${DB_NAME}..."

for file in "$MIGRATIONS_DIR"/*.up.sql; do
    echo "Applying $(basename "$file")..."
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file"
done

echo "Migrations complete."
