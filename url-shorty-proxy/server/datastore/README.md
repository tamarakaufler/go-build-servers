# Postgres schema

database: url-shortener
table: shorty
  columns:
    id      INTEGER PRIMARY
    short   CHAR VARYING(50) NOT NULL
    full    CHAR VARYING(255) NOT NULL
    UNIQUE (short)

# Running Postgres locally in a docker container

## Local development

Create a custom bridge and connect running containers to it.
Run databases in named containers and use their container name as their hostname.

a) custom network

    docker network create --driver bridge url-shortener-bridge

b) postgres database

    docker run --network=url-shortener-bridge --rm -d -e POSTGRES_PASSWORD=mypass -p 5432:5432 -d --name url-shortener-postgres url-shortener-postgres

c)
    ca) start a container connected to the custom bridge

    docker run --network=url-shortener-bridge --rm -d -e DB_HOST=url-shortener-postgres:5432 -e DB_USER=shorty -e DB_PASSWORD=shortypassword -p 50051:50051 --name url-shortener-proxy quay.io/tamarakaufler/url-shortener-proxy:v1alpha1

    OR

    ca) connect a running container to the custom bridge

    - docker network connect url-shortener-bridge url-shortener-postgres
    - docker network connect url-shortener-bridge url-shortener-proxy

## Checking the database

a) Enter a running postgres container

    docker exec -it url-shortener-postgres bash

b) Acess the dapabase directly

    docker exec -it url-shortener-postgres psql -U postgres url_shortener

    docker exec -it url-shortener-postgres psql -U shorty_user url_shortener

# Postgres deployment in Kubernetes

https://gist.github.com/nathanborror/36ebcb42472775b1aa4a8edc135ee615

