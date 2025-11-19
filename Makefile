# Makefile for a Go project

## ------------------------------------------------------------------------------------
## Configuration Variables
## ------------------------------------------------------------------------------------

# Name of the final binary. By default, it uses the current directory's name.
BINARY_NAME ?= $(shell basename $(CURDIR))

# Go command
GO ?= go

# Flags for the Go compiler. -v for verbose output.
GOFLAGS ?= -v

# Output directory for the binary
OUTPUT_DIR ?= ./tmp

# Environment file for local development.
ENV_FILE ?= local.env

# Inject version information at build time.
# Usage: make build VERSION=1.0.0
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.BuildDate=$(BUILD_DATE)'

# Include variables from the env file if it exists, and export them to the environment.
# The leading dash `-` prevents errors if the file doesn't exist.
-include $(ENV_FILE)
export


## ------------------------------------------------------------------------------------
## Main Targets
## ------------------------------------------------------------------------------------

# Default goal when running `make` without arguments.
.DEFAULT_GOAL := help

# .PHONY prevents conflicts with files that have the same name as targets.
.PHONY: all build run watch clean test test-coverage lint

all: build ## Builds the binary. An alias for 'build'.

build-indexer: tidy ## Compiles the source code and creates the binary in $(OUTPUT_DIR).
	@echo "==> Compiling binary..."
	@mkdir -p $(OUTPUT_DIR)
	$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME) cmd/indexer/main.go

run-indexer: build-indexer ## Builds and runs the binary.
	@echo "==> Running the application..."
	@$(OUTPUT_DIR)/$(BINARY_NAME) -path=data

build-server: tidy ## Compiles the source code and creates the binary in $(OUTPUT_DIR).
	@echo "==> Compiling binary..."
	@mkdir -p $(OUTPUT_DIR)
	$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME) cmd/server/main.go

run-server: build-server ## Builds and runs the binary.
	@echo "==> Running the application..."
	@$(OUTPUT_DIR)/$(BINARY_NAME)

build-crawler: tidy ## Compiles the source code and creates the binary in $(OUTPUT_DIR).
	@echo "==> Compiling binary..."
	@mkdir -p $(OUTPUT_DIR)
	$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(BINARY_NAME) cmd/crawler/main.go

run-crawler: build-crawler ## Builds and runs the binary.
	@echo "==> Running the application..."
	@$(OUTPUT_DIR)/$(BINARY_NAME)

watch: build-server ## Runs the application in development mode with live-reloading using Air.
	@echo "==> Starting in watch mode with Air (loading $(ENV_FILE))..."
	@air

clean: ## Removes compiled binaries and temporary files.
	@echo "==> Cleaning project..."
	@rm -rf $(OUTPUT_DIR)
	@rm -f coverage.*
	@rm -rf tmp

test: ## Runs all project tests.
	@echo "==> Running tests..."
	$(GO) test $(GOFLAGS) ./...

test-coverage: ## Runs tests and generates an HTML coverage report.
	@echo "==> Running tests with coverage..."
	$(GO) test -cover -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "==> Coverage report generated in coverage.html"

lint: ## Runs the linter (golangci-lint) to analyze the code.
	@echo "==> Linting code with golangci-lint..."
	@golangci-lint run

## ------------------------------------------------------------------------------------
## Utility Targets
## ------------------------------------------------------------------------------------

.PHONY: deps tidy lint-install air-install air-init help

deps: ## Downloads the Go module dependencies.
	@echo "==> Downloading module dependencies..."
	$(GO) mod download

tidy: ## Tidies and cleans the Go module dependencies.
	@echo "==> Tidying module dependencies..."
	$(GO) mod tidy

lint-install: ## Installs golangci-lint.
	@echo "==> Installing golangci-lint..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

air-install: ## Installs Air for live-reloading.
	@echo "==> Installing Air..."
	$(GO) install github.com/air-verse/air@latest

air-init: ## Creates a default .air.toml configuration file if it doesn't exist.
	@echo "==> Creating default .air.toml configuration file..."
	@[ -f .air.toml ] || printf 'root = "."\ntmp_dir = "tmp"\n\n[build]\n  cmd = "go build -o ./tmp/$(BINARY_NAME) ."\n  bin = "./tmp/$(BINARY_NAME)"\n  full_bin = "./tmp/$(BINARY_NAME)"\n  include_ext = ["go", "tpl", "tmpl", "html"]\n  exclude_dir = ["assets", "tmp", "vendor", "testdata"]\n  log = "air_errors.log"\n  delay = 1000 # ms\n\n[log]\n  time = true\n\n[color]\n  main = "yellow"\n  watcher = "cyan"\n  build = "blue"\n  runner = "green"\n' > .air.toml

help: ## Shows this help message.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)