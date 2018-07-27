# Postgres refresher

## psql commands

\d __database__

\c __database__
\d __table__

## Schema

    CREATE TABLE $SHORTY_DB_TABLE(
    id             SERIAL PRIMARY KEY,
    shorty         VARCHAR(50) NOT NULL,
    url            VARCHAR(255) NOT NULL,
    UNIQUE (shorty)
    );

## SQL
INSERT INTO shorties VALUES(DEFAULT,'ggl.c', 'https://google.com');

INSERT INTO shorties(shorty,url) VALUES('ggl.u', 'https://google.co.uk');

## Postgres running in a Docker container

### Access the postgres interactive terminal
docker exec -it url-shortener-postgres psql -U shorty_user url_shortener

### Run SQL command locally
docker exec -it url-shortener-postgres psql -U shorty_user url_shortener --command "SELECT * FROM shorties"

    - url-shortener-postgres    : container name
    - url_shortener             : database
    - psql ... --command "...." : running an SQL command outside the database env