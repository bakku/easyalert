#!/bin/bash
set -e

POSTGRES="psql --username ${POSTGRES_USER}"

echo "Creating database: ${DEV_DB}"

$POSTGRES <<EOSQL
CREATE DATABASE ${DEV_DB} OWNER ${POSTGRES_USER};
EOSQL

echo "Creating database: ${TEST_DB}"

$POSTGRES <<EOSQL
CREATE DATABASE ${TEST_DB} OWNER ${POSTGRES_USER};
EOSQL
