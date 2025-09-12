# Simple Makefile for darts-counter
# Common tasks: build, test, lint

APP_NAME := darts-counter
BIN_DIR := bin
BUILD_OUT := $(BIN_DIR)/$(APP_NAME)
PKG := ./...

.PHONY: all build test lint clean

all: build

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

build: $(BIN_DIR)
	go build -o $(BUILD_OUT) ./cmd/server
	@echo "Built $(BUILD_OUT)"

# Run all Go tests
# Uses module-aware testing across all packages
# Frontend directory is ignored automatically by Go tooling

test:
	go test $(PKG) -cover

# Lint using golangci-lint (configured via .golangci.yml)
# Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.60.3
lint:
	golangci-lint run

clean:
	rm -rf $(BIN_DIR)
	@echo "Cleaned build artifacts"