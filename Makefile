.PHONY: all build test lint fmt clean

all: fmt lint test build

build:
	go build -o bin/cli-conn .

test:
	go test -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...
	gofumpt -w .

clean:
	rm -rf bin/ coverage.out
