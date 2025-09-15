# Simple Makefile for darts-counter
# Common tasks: build, test, lint

APP_NAME := darts-counter
BIN_DIR := bin
BUILD_OUT := $(BIN_DIR)/$(APP_NAME)
PKG := ./...

.PHONY: all build test lint fmt clean

all: build

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

build: $(BIN_DIR)
	GOTOOLCHAIN=go1.25.0+auto go build -o $(BUILD_OUT) ./cmd/server
	@echo "Built $(BUILD_OUT)"

# Run all Go tests
# Uses module-aware testing across all packages
# Frontend directory is ignored automatically by Go tooling

test:
	GOTOOLCHAIN=go1.25.0+auto go test $(PKG) -cover

# Lint using golangci-lint (configured via .golangci-lint.yml)
# Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.62.0
lint:
	GOTOOLCHAIN=go1.25.0+auto go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.0 run -c .golangci-lint.yml

# Format code with standard gofmt (matches linter)
fmt:
	gofmt -s -w .

clean:
	rm -rf $(BIN_DIR)
	@echo "Cleaned build artifacts"