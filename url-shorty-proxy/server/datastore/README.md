# Postgres schema

database: url-shortener
table: shorties
  columns:
    id             SERIAL PRIMARY KEY
    shorty         VARCHAR(50) NOT NULL
    url            VARCHAR(255) NOT NULL
    UNIQUE (shorty)

# Running Postgres locally in a docker container

## Local development - running both the database and the proxy in docker containers

Create a custom bridge and connect running containers to it.
Run databases in named containers and use their container name as their hostname.

a) custom network

    docker network create --driver bridge url-shortener-bridge

b) postgres database

    ba) start a container connected to the custom bridge

    docker run --network=url-shortener-bridge --rm -d -e POSTGRES_PASSWORD=mypass -p 5432:5432 -d --name url-shortener-postgres url-shortener-postgres

    (make runpostgres ... in Makefile under the bootstapping/postgres dir)

    OR

    bb) connect a running container to the custom bridge

    docker network connect url-shortener-bridge url-shortener-postgres

c) url shortener proxy

    ca) start a container connected to the custom bridge

    docker run --network=url-shortener-bridge --rm -d -e DB_HOST=url-shortener-postgres:5432 -e DB_USER=shorty -e DB_PASSWORD=shortypassword -p 50051:50051 --name url-shortener-proxy quay.io/tamarakaufler/url-shortener-proxy:v1alpha1

    (make runproxy ... in Makefile under server dir)

    OR

    cb) connect a running container to the custom bridge

    docker network connect url-shortener-bridge url-shortener-proxy

## Checking the database

a) Enter a running postgres container

    docker exec -it url-shortener-postgres bash

b) Access the database directly

    docker exec -it url-shortener-postgres psql -U postgres url_shortener

    docker exec -it url-shortener-postgres psql -U shorty_user url_shortener

## Local development - running Postgres in a container and the proxy locally

    a) Postgres

    docker run --rm -d -e POSTGRES_PASSWORD=mypass -p 5432:5432 -d --name url-shortener-postgres url-shortener-postgres

        Using -p 5432:5432 allows to access the database at localhost:5432/127.0.0.1:5432

    b) Proxy

    go run main.go

        No need to set up any env variables for access to the database. Default values are sufficient.


# Kubernetes

## Postgres deployment

https://gist.github.com/nathanborror/36ebcb42472775b1aa4a8edc135ee615

