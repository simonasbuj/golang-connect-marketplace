PATH := $(PATH):$(shell go env GOPATH)/bin
export PATH

# Load .env.example and .env
ifneq (,$(wildcard .env.example))
    include .env.example
    export $(shell sed 's/=.*//' .env.example)
endif

ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif


run:
	go run ./cmd/api

# pre-commit
lint:
	golangci-lint run --verbose --max-issues-per-linter=0 --max-same-issues=0

lint-fix:
	golangci-lint run --verbose --fix

.PHONY: test
test:
	go test -v ./...

# database
start-postgres:
	docker-compose -p go-marketplace -f infra/docker/postgres-docker-compose.yml up -d --build

stop-postgres:
	docker-compose -p go-marketplace -f infra/docker/postgres-docker-compose.yml stop
