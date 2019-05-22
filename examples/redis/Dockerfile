# Build Stage
FROM golang:1.12.5-alpine3.9 AS build-env
RUN mkdir /app
WORKDIR /app
RUN apk add git
COPY vendor ./vendor

COPY . /app
RUN cd /app && go build -mod vendor -o captin

# Final Stage
FROM alpine
WORKDIR /app
COPY --from=build-env /app/captin /app/
COPY ./hooks.json hooks.json
ENTRYPOINT ./captin ./hooks.json