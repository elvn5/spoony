#!/bin/bash
set -e

FRESH_INIT=false
if [ ! -f /var/lib/postgresql/data/PG_VERSION ]; then
    echo "[entrypoint] Initializing PostgreSQL..."
    su-exec postgres initdb -D /var/lib/postgresql/data --auth=trust --username=postgres
    FRESH_INIT=true
fi

# Start PG temporarily to create/fix the app user + database.
su-exec postgres pg_ctl start -D /var/lib/postgresql/data -w -t 30 \
    -o "-c listen_addresses='localhost'"

if [ "$FRESH_INIT" = true ]; then
    su-exec postgres psql -v ON_ERROR_STOP=1 --username=postgres <<-EOSQL
        CREATE USER spoony_user WITH PASSWORD 'spoony_password';
        CREATE DATABASE spoony OWNER spoony_user;
        GRANT ALL PRIVILEGES ON DATABASE spoony TO spoony_user;
EOSQL
    echo "[entrypoint] PostgreSQL initialized."
else
    # Volumes from deploys before the tma-boilerplate → spoony rename still
    # have the role/database under the old names. Rename them in place so
    # existing data survives the switch to spoony_user/spoony.
    OLD_ROLE_EXISTS=$(su-exec postgres psql --username=postgres -tAc "SELECT 1 FROM pg_roles WHERE rolname='tma_user'")
    if [ "$OLD_ROLE_EXISTS" = "1" ]; then
        echo "[entrypoint] Renaming legacy tma_user/tma_boilerplate to spoony_user/spoony..."
        su-exec postgres psql -v ON_ERROR_STOP=1 --username=postgres <<-EOSQL
            ALTER USER tma_user RENAME TO spoony_user;
            ALTER USER spoony_user WITH PASSWORD 'spoony_password';
            ALTER DATABASE tma_boilerplate RENAME TO spoony;
            ALTER DATABASE spoony OWNER TO spoony_user;
EOSQL
    fi
fi

su-exec postgres pg_ctl stop -D /var/lib/postgresql/data -w

echo "[entrypoint] Starting all services via supervisord..."
exec supervisord -c "${SUPERVISORD_CONF:-/etc/supervisord.conf}"
