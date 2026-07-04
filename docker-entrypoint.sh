#!/bin/bash
set -e

# Initialize PostgreSQL data directory on first run
if [ ! -f /var/lib/postgresql/data/PG_VERSION ]; then
    echo "[entrypoint] Initializing PostgreSQL..."
    su-exec postgres initdb -D /var/lib/postgresql/data --auth=trust --username=postgres

    # Start PG temporarily to create user + database
    su-exec postgres pg_ctl start -D /var/lib/postgresql/data -w -t 30 \
        -o "-c listen_addresses='localhost'"

    su-exec postgres psql -v ON_ERROR_STOP=1 --username=postgres <<-EOSQL
        CREATE USER spoony_user WITH PASSWORD 'spoony_password';
        CREATE DATABASE spoony OWNER spoony_user;
        GRANT ALL PRIVILEGES ON DATABASE spoony TO spoony_user;
EOSQL

    su-exec postgres pg_ctl stop -D /var/lib/postgresql/data -w
    echo "[entrypoint] PostgreSQL initialized."
fi

echo "[entrypoint] Starting all services via supervisord..."
exec supervisord -c "${SUPERVISORD_CONF:-/etc/supervisord.conf}"
