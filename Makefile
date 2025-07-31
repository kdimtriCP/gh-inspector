# Simple Makefile for gh-inspector

.PHONY: build test test-verbose test-coverage lint run clean swag

build:
	go build -o gh-inspector .

test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run || true

run:
	go run main.go

clean:
	rm -f gh-inspector
	rm -f coverage.out coverage.html
	rm -rf .gh-inspector-cache
	rm -rf swagger

swag:
	swag init -g cmd/docs.go -o swagger --parseDependency --parseInternal

