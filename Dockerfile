# Build Stage
FROM golang:1.12.5-alpine3.9 AS build-env
RUN mkdir /app
WORKDIR /app
RUN apk add git
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /app
RUN cd /app && go build -o captin cmd/captin/api.go

# Final Stage
FROM alpine
WORKDIR /app
COPY --from=build-env /app/captin /app/
COPY ./example/config.json config.json
ENTRYPOINT ./captin ./config.json