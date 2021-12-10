.PHONY: test security run stop

SERVER_PORT = 8080
SERVER_URL = "0.0.0.0:$(SERVER_PORT)"
SERVER_READ_TIMEOUT = 60
JWT_SECRET_KEY = "GOWA_JWT_SECRET:base64(string):amNjx+OkIltCJU3aTYhO3A=="
JWT_SECRET_KEY_EXPIRE_MINUTES = 15
WHATSAPP_CLIENT_VERSION_MAJOR = 2
WHATSAPP_CLIENT_VERSION_MINOR = 2126
WHATSAPP_CLIENT_VERSION_BUILD = 11
WHATSAPP_CLIENT_SESSION_PATH = "./storage/temps"

BUILD_DIR = $(PWD)/bin/app

security:
	gosec -quiet ./...

test: security
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

swag:
	swag init

docker_build_image:
	docker build -t gowa .

docker_app: docker_build_image
	docker run -d \
        		--name gowa-c \
        		-p $(SERVER_PORT):$(SERVER_PORT) \
        		-e SERVER_URL=$(SERVER_URL) \
        		-e SERVER_READ_TIMEOUT=$(SERVER_READ_TIMEOUT) \
        		-e JWT_SECRET_KEY=$(JWT_SECRET_KEY) \
        		-e JWT_SECRET_KEY_EXPIRE_MINUTES=$(JWT_SECRET_KEY_EXPIRE_MINUTES) \
        		-e WHATSAPP_CLIENT_VERSION_MAJOR=$(WHATSAPP_CLIENT_VERSION_MAJOR) \
        		-e WHATSAPP_CLIENT_VERSION_MINOR=$(WHATSAPP_CLIENT_VERSION_MINOR) \
        		-e WHATSAPP_CLIENT_VERSION_BUILD=$(WHATSAPP_CLIENT_VERSION_BUILD) \
        		-e WHATSAPP_CLIENT_SESSION_PATH=$(WHATSAPP_CLIENT_SESSION_PATH) \
        		gowa

run: swag docker_app

stop:
	docker container stop gowa-c
	docker container rm gowa-c
	docker rmi gowa