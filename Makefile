.PHONY: all clean build install uninstall check lint fmt vendor test

PACKAGE_NAME := lazyhis
BUILD_DIR    := build
DST_DIR      ?= $(HOME)/.local/bin

VERSION := $(shell git describe --tags --dirty)
COMMIT  := $(shell git rev-parse --short HEAD)
DATE    := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

all: install

## CODE QUALITY & TESTS
check: vendor fmt lint test

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

## BUILD & INSTALLATION
info:
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Date:    $(DATE)"

build: info
	@echo "Building $(PACKAGE_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build \
		-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" \
		-o $(BUILD_DIR)/$(PACKAGE_NAME) main.go

install: build
	@echo "Installing $(PACKAGE_NAME)..."
	@mkdir -p $(DST_DIR)
	ln -sf $(PWD)/$(BUILD_DIR)/$(PACKAGE_NAME) $(DST_DIR)

## CLEAN
uninstall:
	@echo "Uninstalling $(PACKAGE_NAME)..."
	rm -f $(DST_DIR)/$(PACKAGE_NAME)

clean: uninstall
	@echo "Cleaning $(BUILD_DIR) directory..."
	rm -rf $(BUILD_DIR)
