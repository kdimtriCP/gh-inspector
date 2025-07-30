# Simple Makefile for gh-inspector

.PHONY: build test lint run

build:
	go build -o gh-inspector .

test:
	go test ./...

lint:
	golangci-lint run || true

run:
	go run main.go

