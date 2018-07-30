IMAGE_TAG?=v1alpha1
QUAY_PASS?=biggestsecret
MONGODB_ADMIN_USER?=admin
MONGODB_ADMIN_PASS?=admin
SHORTY_DB_USER?=shorty_user
SHORTY_DB_PASS?=shortypass
SHORTY_DB?=url_shortener

net:
	docker network create --driver bridge url-shortener-bridge

build:
	docker build -t quay.io/tamarakaufler/url-shorty-mongodb:$(IMAGE_TAG) .

runmongodb:
	docker run --name url-shorty-mongodb --network=url-shortener-bridge  -v $(PWD)/dbdata:/data/db --rm -d -p 27017:27017 quay.io/tamarakaufler/url-shorty-mongodb:$(IMAGE_TAG) --auth

rmdbcontainer:
	docker ps | grep "mongodb" | awk '{print $1}' | xargs docker rm -f

access:
	docker exec -it url-shorty-mongodb mongo --username $(SHORTY_DB_USER) --password $(SHORTY_DB_PASS)  $(SHORTY_DB)

adminaccess:
	docker exec -it url-shorty-mongodb mongo --username $(MONGODB_ADMIN_USER) --password $(MONGODB_ADMIN_PASS) admin
