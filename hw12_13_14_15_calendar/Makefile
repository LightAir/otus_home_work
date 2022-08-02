BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

MIGRATION := "./bin/migration"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

up-docker:
	docker-compose -f deployments/docker-compose.yaml up -d

down-docker:
	docker-compose -f deployments/docker-compose.yaml down

generate:
	go generate ./...

build-migrations:
	go build -v -o $(MIGRATION) ./cmd/migration

migrate: build-migrations
	$(MIGRATION) --config=./configs/config.yaml

.PHONY: build run build-img run-img version test lint up-docker down-docker generate migrate build-migrations
