IMAGE_TAG?=v1alpha1
QUAY_PASS?=biggestsecret

net:
	docker network create --driver bridge url-shortener-bridge		

build:
	docker build -t quay.io/tamarakaufler/url-shortener-postgres:$(IMAGE_TAG) .

localrun:
	docker run  --name url-shortener-postgres  --network=url-shortener-bridge --rm -d -p 5432:5432 quay.io/tamarakaufler/url-shortener-postgres:$(IMAGE_TAG)

rmdbcontainer:
	docker ps | grep "postgres" | awk '{print $1}' | xargs docker rm -f

access:
	docker exec -it url-shortener-postgres psql -U shorty_user url_shortener
