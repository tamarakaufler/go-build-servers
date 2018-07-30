#!/usr/bin/env bash
set -o errexit

# Application Database User
URL_SHORTENER_DB=${URL_SHORTENER_DB:-url_shortener}
SHORTY_DB_USER=${SHORTY_DB_USER:-shorty_user}
SHORTY_DB_PASS=${SHORTY_DB_PASS:-shortypass}

# Admin User
MONGODB_ADMIN_USER=${MONGODB_ADMIN_USER:-admin}
MONGODB_ADMIN_PASS=${MONGODB_ADMIN_PASS:-admin}

# Debug ----------------------------------------------------
echo "==> MONGODB_ADMIN_USER ... $MONGODB_ADMIN_USER"
echo "==> URL_SHORTENER_DB ... $URL_SHORTENER_DB"
echo "==> SHORTY_DB_USER ... $SHORTY_DB_USER"
echo "==> SHORTY_DB_PASS ... [$SHORTY_DB_PASS]"
echo "==> MONGODB_ADMIN_PASS = [$MONGODB_ADMIN_PASS]"
# ----------------------------------------------------------

if [[ $MONGODB_ADMIN_PASS == "" ]] ; then
  echo "=> You must specifity MONGODB_ADMIN_PASS environment variable"
  exit 1
fi

main() {
  init_user_and_db
}

# Perform initialization in the already-started Mongodb
#   - create the database
#   - create admin user
#   - set up user for the url-shortener-proxy database:
#         this user needs to be able to create a table,
#         to insert/update and delete records
init_user_and_db() {
    # Create the admin user
    echo "========================================================================"
    echo "=> Creating admin user $MONGODB_ADMIN_USER with a password $MONGODB_ADMIN_PASS in MongoDB"
    mongo admin --eval "db.createUser({user: '$MONGODB_ADMIN_USER', pwd: '$MONGODB_ADMIN_PASS', roles:[{role:'userAdminAnyDatabase',db:'admin'}]});"
    echo "$MONGODB_ADMIN_USER user created!"
    echo "========================================================================"

    if [[ "$URL_SHORTENER_DB" != "admin" ]]; then
        echo "=> Creating an ${URL_SHORTENER_DB} user with a password in MongoDB"
        mongo admin <<-EOF
        use $URL_SHORTENER_DB
        db.createUser({user: '$SHORTY_DB_USER', pwd: '$SHORTY_DB_PASS', roles:[{role:'dbOwner', db:'$URL_SHORTENER_DB'}]})
EOF
        echo "$URL_SHORTENER_DB database created!"
        echo "$SHORTY_DB_USER user for $URL_SHORTENER_DB created!"
        echo "========================================================================"
    else
        echo "=> Skipping database and user creation."
    fi
}

# Executes the main routine with environment variables
# passed through the command line.
main "$@"