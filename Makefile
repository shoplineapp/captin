.PHONY: build test

build:
	@go build -o build/captin cmd/captin/main.go

test:
	@go test ./test/...