#!/bin/bash
set -e

echo "Running database migrations..."
for f in /docker-entrypoint-initdb.d/migrations/*.up.sql; do
    echo "Applying $(basename "$f")..."
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f "$f"
done

if [ -f /docker-entrypoint-initdb.d/seed.sql ]; then
    echo "Seeding database..."
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/seed.sql
fi

echo "Database initialization complete."
