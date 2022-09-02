.PHONY: build test

start: build run-example

start-api: build-api run-example-api

start-docker:
	@docker-compose build && docker-compose up

build:
	@go build -o build/captin cmd/captin/main.go

test:
	@go test -parallel 4 -race $(shell go list ./test/... | grep -v mocks)

run-example:
	@./build/captin ./example/config.json

build-api:
	@go build -o build/captin cmd/captin/api.go

run-example-api:
	@./build/captin ./example/config.json