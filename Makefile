.PHONY: run test build

run:
	go run ./cmd/api

test:
	go test ./...

build:
	go build -o bin/api ./cmd/api
