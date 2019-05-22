.PHONY: start start-docker

start:
	@go build -o build/app && ./build/app ./hooks.json

start-docker:
	@docker-compose build && docker-compose up

scale:
	@docker-compose scale captin=3