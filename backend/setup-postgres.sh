#!/bin/bash
set -euo pipefail

brew services start postgresql
sleep 2

createdb abtalks 2>/dev/null || true
psql -d abtalks -c "CREATE EXTENSION IF NOT EXISTS pgcrypto;" 2>/dev/null || true

echo "PostgreSQL is ready for abtalks"
