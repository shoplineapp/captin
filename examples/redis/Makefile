.PHONY: start start-docker

start:
	@go mod download && go mod vendor
	@go build -o build/app && ./build/app ./hooks.json

start-docker:
	@go mod download && go mod vendor
	@docker-compose build && docker-compose up

restart-captin:
	@docker-compose stop captin
	@docker-compose build && docker-compose up -d captin
	@docker-compose scale captin=3
	@docker-compose restart nginx

scale:
	@docker-compose scale captin=3
	@docker-compose restart nginx