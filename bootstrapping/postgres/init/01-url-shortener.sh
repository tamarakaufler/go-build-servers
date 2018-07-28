#!/usr/bin/env bash

# Credits: based on https://medium.com/@beld_pro/quick-tip-creating-a-postgresql-container-with-default-user-and-password-8bb2adb82342

# This script is used to initialize postgres, after it started running,
# to provide the database(s) and table(s) expected by a connecting
# application.

# postgres is used by a URL shortener proxy,
# which expects:

#   - database called url_shortener
#   - within it a table called shorties

#   * db user with approriate privileges to the database

set -o errexit

URL_SHORTENER_DB=${URL_SHORTENER_DB:-url_shortener}
SHORTY_DB_TABLE=${SHORTY_DB_TABLE:-shorties}
SHORTY_DB_USER=${SHORTY_DB_USER:-shorty_user}
SHORTY_DB_PASSWORD=${SHORTY_DB_PASSWORD:-shortypass}
POSTGRES_USER=${POSTGRES_USER:-postgres}

# By default POSTGRES_PASSWORD is an empty string. For security reasons it is advisable
# to set set it up when we start running the container:
#
#   docker run --rm -e POSTGRES_PASSWORD=mypass -p 5432:5432 -d --name url-shortener-postgres  quay.io/tamarakaufler/url-shortener-postgres:$(IMAGE_TAG)

#   psql -h localhost -p 5432 -U postgres

#       Unlike in MySQL, psql does not provide a flag for providing password.
#       The password is provided interactively.
#       The PostgreSQL image sets up trust authentication locally, so password is not required
#       when connecting from localhost (inside the same container). Ie. psql in this script, 
#       that runs after Postgres starts, does not need the authentication. 

POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-}

# Debug ----------------------------------------------------
echo "==> POSTGRES_USER ... $POSTGRES_USER"
echo "==> POSTGRES_DB ... $POSTGRES_DB"
echo "==> URL_SHORTENER_DB ... $URL_SHORTENER_DB"
echo "==> SHORTY_DB_USER ... $SHORTY_DB_USER"
echo "==> SHORTY_DB_PASSWORD ... [$SHORTY_DB_PASSWORD]"
echo "==> SHORTY_DB_TABLE ... $SHORTY_DB_TABLE"
echo "==> POSTGRES_PASSWORD = [$POSTGRES_PASSWORD]"
# ----------------------------------------------------------

# What environment variables need to be set up.
#   Environment variable defaults are set up in this case, 
#   however we want to ensure the defaults are not accidentally
#   removed from this file causing a problem.
readonly REQUIRED_ENV_VARS=(
  "URL_SHORTENER_DB"
  "SHORTY_DB_USER"
  "SHORTY_DB_PASSWORD"
  "SHORTY_DB_TABLE")

# Main execution:
# - verifies all environment variables are set
# - runs SQL code to create user and database
# - runs SQL code to create table
main() {
  check_env_vars_set
  init_user_and_db

  # Comment out if wanting to use the url-shorty-proxy to use gorm AutoMigrate feature:
  #   the gorm AutoMigrate feature creates extra columns (xxx_unrecognized, xxx_sizecache)
  #   based on the proto message, which are required for proto messages transactions
  #   to work with the table
  # init_db_tables
}

# ----------------------------------------------------------
# HELPER FUNCTIONS

# Check if all of the required environment
# variables are set
check_env_vars_set() {
  for required_env_var in ${REQUIRED_ENV_VARS[@]}; do
    if [[ -z "${!required_env_var}" ]]; then
      echo "Error:
    Environment variable '$required_env_var' not set.
    Make sure you have the following environment variables set:
      ${REQUIRED_ENV_VARS[@]}
Aborting."
      exit 1
    fi
  done
}

# Perform initialization in the already-started PostgreSQL
#   - create the database
#   - set up user for the url-shortener-proxy database:
#         this user needs to be able to create a table,
#         to insert/update and delete records
init_user_and_db() {
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
     CREATE DATABASE $URL_SHORTENER_DB;
     CREATE USER $SHORTY_DB_USER WITH PASSWORD '$SHORTY_DB_PASSWORD';
     GRANT ALL PRIVILEGES ON DATABASE $URL_SHORTENER_DB TO $SHORTY_DB_USER;
EOSQL
}

#   - create database tables
init_db_tables() {
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" "$URL_SHORTENER_DB" <<-EOSQL
    CREATE TABLE $SHORTY_DB_TABLE(
    id             SERIAL PRIMARY KEY,
    shorty         VARCHAR(50) NOT NULL,
    url            VARCHAR(255) NOT NULL,
    UNIQUE (shorty)
);
EOSQL
}

# Executes the main routine with environment variables
# passed through the command line. Added for completeness 
# as not used here.
main "$@"
