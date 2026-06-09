#!/usr/bin/env bash
set -euo pipefail

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-bankuser}"
DB_PASSWORD="${DB_PASSWORD:-bankpass}"
DB_NAME="${DB_NAME:-bankdb}"

export PGPASSWORD="$DB_PASSWORD"
SEED_FILE="$(cd "$(dirname "$0")" && pwd)/seed.sql"

echo "Seeding database ${DB_NAME}..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$SEED_FILE"
echo "Seed complete."
