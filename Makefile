ifneq (,$(wildcard .versions))
include .versions
export
endif

BINARY := server
GO_IMAGE := golang:$(GO_VERSION)-alpine$(ALPINE_VERSION)

.PHONY: all go-tidy go-get swag-init docker-build docker-down docker-up

all: docker-down swag-init docker-build docker-up


go-tidy:
	docker run --rm \
		-v $(PWD):/app:ro \
		-v $(PWD)/go.mod:/app/go.mod \
		-v $(PWD)/go.sum:/app/go.sum \
		-w /app \
		$(GO_IMAGE) \
		go mod tidy

# for e.g. make go-get PKG=github.com/golang-jwt/jwt/v5
go-get:
	docker run --rm \
		-v $(PWD):/app:ro \
		-v $(PWD)/go.mod:/app/go.mod \
		-v $(PWD)/go.sum:/app/go.sum \
		-w /app \
		$(GO_IMAGE) \
		sh -c 'go get $(PKG) && go mod tidy'

swag-init:
	swag init

docker-build:
	docker compose --env-file .versions build

docker-down:
	docker compose --env-file .versions down

docker-up:
	docker compose --env-file .versions up -d
