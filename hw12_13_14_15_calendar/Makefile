BIN := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
BIN_MIGRATION := "./bin/migration"

DOCKER_IMG="calendar:develop"
DOCKER_IMG_SCHEDULER="scheduler:develop"
DOCKER_IMG_SENDER="sender:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run: build
	$(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile-cal .
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SCHEDULER) \
		-f build/Dockerfile-sch .
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SENDER) \
		-f build/Dockerfile-snd .

run-img:
	build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

up:
	docker-compose -f deployments/docker-compose.yaml up -d

down:
	docker-compose -f deployments/docker-compose.yaml down

generate:
	go generate ./...

build-migrations:
	go build -v -o $(BIN_MIGRATION) ./cmd/migration

migrate: build-migrations
	$(BIN_MIGRATION) --config=./configs/config.yaml

integration-tests:
	docker-compose -f deployments/docker-compose.test.yaml up --abort-on-container-exit --exit-code-from testruner testruner
	docker-compose -f deployments/docker-compose.test.yaml down

.PHONY: build run build-img run-img version test lint up down generate migrate build-migrations integration-tests