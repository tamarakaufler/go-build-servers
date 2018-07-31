# Database bootstrapping

MongoDB and Postgres bootstrapping is provided for running the respective database in a Docker container and set up for the use with an application. The examples here cater for the url-shorty-proxy server implementation (https://github.com/tamarakaufler/go-build-servers/tree/master/url-shorty-proxy).

## Preserving data

Container with a bind mount is used to preserve the database data.

A local directory is bind-mounted into the container at the container directory where the database stores its data:

	docker run --name url-shorty-postgres --network=url-shortener-bridge -v $(PWD)/psqldata:/var/lib/postgresql/data --rm -d -p 5432:5432 quay.io/tamarakaufler/url-shorty-postgres:$(IMAGE_TAG)

  docker run --name url-shorty-mongodb --network=url-shortener-bridge  -v $(PWD)/mgodata:/data/db --rm -d -p 27017:27017 quay.io/tamarakaufler/url-shorty-mongodb:$(IMAGE_TAG) --auth