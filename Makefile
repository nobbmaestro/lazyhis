.PHONY: all clean build pre-commit lint fmt vendor test

PACKAGE_NAME := lazyhis
BUILD_DIR := build

VERSION := $(shell git describe --tags --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

all: build

pre-commit: vendor fmt lint test

build: 
	@echo "Building $(PACKAGE_NAME)..."
	go build -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o $(BUILD_DIR)/lazyhis main.go

clean:
	@echo "Cleaning $(BUILD_DIR) directory..."
	rm -rf $(BUILD_DIR)

vendor: 
	@echo "Vendor dependencies..."
	go mod tidy
	go mod vendor

fmt: 
	@echo "Formatting code..."
	gofmt -w .
	golines -w . --max-len=88

lint: 
	@echo "Running linters..."
	golangci-lint run

test: 
	@echo "Running tests..."
	go test -v ./...
