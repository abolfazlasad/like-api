ifneq (,$(wildcard .versions))
include .versions
export
endif

BINARY := server
GO_IMAGE := golang:$(GO_VERSION)-alpine$(ALPINE_VERSION)

.PHONY: \
	all \
	test \
	go-tidy \
	go-get \
	swag-init \
	docker-build \
	docker-down \
	docker-up \
	integration-test \
	integration-test-build \
	integration-test-up \
	integration-test-down

all: docker-down swag-init docker-build docker-up

test:
	richgo test -v -cover -race -short -count=1 ./...

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

integration-test-build:
	docker compose --env-file .versions \
		-f docker-compose.yml \
		-f docker-compose.test.yml \
		build
integration-test-up:
	docker compose --env-file .versions \
		-f docker-compose.yml \
		-f docker-compose.test.yml \
		up  --abort-on-container-exit --exit-code-from tests --attach tests;
integration-test-down:
	docker compose --env-file .versions \
		-f docker-compose.yml \
		-f docker-compose.test.yml \
		down --remove-orphans
integration-test: integration-test-build integration-test-up integration-test-down
