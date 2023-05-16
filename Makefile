.PHONY: build test

build:
	go build -o bin/go-together ./cmd

test:
	go test -v ./...
