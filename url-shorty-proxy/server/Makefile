IMAGE_TAG=v1alpha1
QUAY_PASS?=biggestsecret
#DB_HOST?=url-shorty-postgres:5432
DB_HOST?=url-shorty-postgres:27017
	
dev:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o url-shorty-proxy -a -installsuffix cgo main.go
	docker build -t quay.io/tamarakaufler/url-shorty-proxy:$(IMAGE_TAG) .
	
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o url-shorty-proxy -a -installsuffix cgo main.go
	docker build -t quay.io/tamarakaufler/url-shorty-proxy:$(IMAGE_TAG) .
	docker login quay.io -u tamarakaufler -p $(QUAY_PASS)
	docker push quay.io/tamarakaufler/url-shorty-proxy:$(IMAGE_TAG)

runproxy:
	docker run --net="url-shortener-bridge" \
		--name=url-shorty-proxy \
		--rm \
		-p 8888:8888 \
		-e DB_HOST=$(DB_HOST) \
	quay.io/tamarakaufler/url-shorty-proxy:$(IMAGE_TAG)

rmcontainer:
	docker ps | grep "shorty-proxy" | awk '{print $1}' | xargs docker rm -f
