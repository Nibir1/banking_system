#!/bin/sh

# chmod +x start.sh - Run this so that it can execute - Execution permission

set -e

echo "run db migration"
source /app/app.env #Remove this code when building locally
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"


