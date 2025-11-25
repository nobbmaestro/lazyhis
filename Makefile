.PHONY: all clean build install uninstall install-bin install-man uninstall-bin uninstall-man check lint fmt vendor test

PACKAGE_NAME := lazyhis
BUILD_DIR    := build
DST_DIR      ?= $(HOME)/.local/bin
DST_MAN_DIR  ?= $(HOME)/.local/share/man/man1

VERSION := $(shell git describe --tags --dirty)
COMMIT  := $(shell git rev-parse --short HEAD)
DATE    := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

all: build

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

install-bin: build
	@echo "Installing $(PACKAGE_NAME)..."
	@mkdir -p $(DST_DIR)
	install -m 755 $(PWD)/$(BUILD_DIR)/$(PACKAGE_NAME) $(DST_DIR)

install-man: build
	@echo "Installing man pages..."
	@mkdir -p $(DST_MAN_DIR)
	$(BUILD_DIR)/$(PACKAGE_NAME) gen man --dst $(DST_MAN_DIR)

install: install-bin install-man

## CLEAN
uninstall-bin:
	@echo "Uninstalling $(PACKAGE_NAME)..."
	rm -f $(DST_DIR)/$(PACKAGE_NAME)

uninstall-man:
	@echo "Uninstalling man pages..."
	rm -f $(DST_MAN_DIR)/$(PACKAGE_NAME)*.1

uninstall: uninstall-bin uninstall-man

clean:
	@echo "Cleaning $(BUILD_DIR) directory..."
	rm -rf $(BUILD_DIR)
