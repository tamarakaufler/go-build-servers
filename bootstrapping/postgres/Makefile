IMAGE_TAG?=v1alpha1
QUAY_PASS?=biggestsecret

net:
	docker network create --driver bridge url-shortener-bridge		

build:
	docker build -t quay.io/tamarakaufler/url-shorty-postgres:$(IMAGE_TAG) .

runpostgres:
	docker run --name url-shorty-postgres --network=url-shortener-bridge -v $(PWD)/psqldata:/var/lib/postgresql/data --rm -d -p 5432:5432 quay.io/tamarakaufler/url-shorty-postgres:$(IMAGE_TAG)

rmdbcontainer:
	docker ps | grep "postgres" | awk '{print $1}' | xargs docker rm -f

access:
	docker exec -it url-shorty-postgres psql -U postgres url_shortener
