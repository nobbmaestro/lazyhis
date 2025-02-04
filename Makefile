.PHONY: all clean build

PACKAGE_NAME := lazyhis
BUILD_DIR := build

VERSION := $(shell git describe --tags --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

all: build

build:
	@echo "Building..."
	go build -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o $(BUILD_DIR)/lazyhis main.go

clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)
