# Postgres schema

database: url-shortener
table: shorties
  columns:
    id             SERIAL PRIMARY KEY
    shorty         VARCHAR(50) NOT NULL
    url            VARCHAR(255) NOT NULL
    UNIQUE (shorty)

# Running Postgres locally in a docker container

## Local development

Create a custom bridge and connect running containers to it.
Run databases in named containers and use their container name as their hostname.

a) custom network

    docker network create --driver bridge url-shortener-bridge

b) postgres database

    ba) start a container connected to the custom bridge

    docker run --network=url-shortener-bridge --rm -d -e POSTGRES_PASSWORD=mypass -p 5432:5432 -d --name url-shortener-postgres url-shortener-postgres

    OR

    bb) connect a running container to the custom bridge

    docker network connect url-shortener-bridge url-shortener-postgres

c) url shortener proxy

    ca) start a container connected to the custom bridge

    docker run --network=url-shortener-bridge --rm -d -e DB_HOST=url-shortener-postgres:5432 -e DB_USER=shorty -e DB_PASSWORD=shortypassword -p 50051:50051 --name url-shortener-proxy quay.io/tamarakaufler/url-shortener-proxy:v1alpha1

    OR

    cb) connect a running container to the custom bridge

    docker network connect url-shortener-bridge url-shortener-proxy

## Checking the database

a) Enter a running postgres container

    docker exec -it url-shortener-postgres bash

b) Acess the dapabase directly

    docker exec -it url-shortener-postgres psql -U postgres url_shortener

    docker exec -it url-shortener-postgres psql -U shorty_user url_shortener

## Running Postgres in a container and the proxy locally

## Running both Postgres and the proxy in their respective containers


# Kubernetes

## Postgres deployment

https://gist.github.com/nathanborror/36ebcb42472775b1aa4a8edc135ee615

