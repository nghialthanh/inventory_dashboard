#!/bin/bash
set -euo pipefail

# Allow streaming replication from other Compose services.
echo "host replication replicator 0.0.0.0/0 md5" >> "${PGDATA}/pg_hba.conf"
echo "host replication replicator ::0/0 md5" >> "${PGDATA}/pg_hba.conf"

psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" --dbname "${POSTGRES_DB}" <<'SQL'
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'replicator') THEN
    CREATE ROLE replicator WITH REPLICATION LOGIN PASSWORD 'replicator';
  END IF;
END
$$;
SQL
