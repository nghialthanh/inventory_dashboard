#!/bin/sh
set -eu

DATA_DIR=/var/lib/postgresql/data
PRIMARY_HOST="${PRIMARY_HOST:-postgres}"
PRIMARY_PORT="${PRIMARY_PORT:-5432}"

echo "waiting for primary ${PRIMARY_HOST}:${PRIMARY_PORT}"
until pg_isready -h "${PRIMARY_HOST}" -p "${PRIMARY_PORT}" -U inventory; do
  sleep 1
done

# Fresh replica data directory each start (volume may be empty or stale).
rm -rf "${DATA_DIR:?}/"*
export PGPASSWORD=replicator
pg_basebackup \
  -h "${PRIMARY_HOST}" \
  -p "${PRIMARY_PORT}" \
  -U replicator \
  -D "${DATA_DIR}" \
  -Fp -Xs -P -R

chown -R postgres:postgres "${DATA_DIR}"
exec docker-entrypoint.sh postgres -c hot_standby=on
