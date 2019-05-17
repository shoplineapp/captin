.PHONY: build test

start: build run-example

build:
	@go build -o build/captin cmd/captin/main.go

test:
	@go test ./test/...

run-example:
	@./build/captin ./example/config.json