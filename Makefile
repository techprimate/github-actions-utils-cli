# ============================================================================
# GITHUB ACTIONS UTILS CLI MAKEFILE
# ============================================================================
# This Makefile provides automation for building, testing, and developing
# the GitHub Actions Utils CLI. Run 'make help' to see all available commands.
# ============================================================================

# ============================================================================
# SETUP & INSTALLATION
# ============================================================================

## Initialize project for development (installs all dependencies)
#
# This command sets up your development environment by:
# - Installing system dependencies via Homebrew (if available)
# - Installing Go module dependencies
# - Preparing the project for development
#
# Run this once when you first clone the repository.
.PHONY: init
init:
	@if [ "$$(uname)" = "Darwin" ]; then \
		echo "Darwin detected."; \
		$(MAKE) init-darwin; \
	elif [ "$$(uname)" = "Linux" ]; then \
		echo "Linux detected."; \
		$(MAKE) init-linux; \
	else \
		echo "Not running on Darwin or Linux."; \
		exit 1; \
	fi
	$(MAKE) install

.PHONY: init-darwin
init-darwin:
	@if ! command -v brew >/dev/null 2>&1; then \
		echo "Homebrew not detected. Skipping system dependency installation."; \
		exit 1; \
	fi
	echo "Homebrew detected. Installing system dependencies..."; \
	brew bundle

.PHONY: init-linux
init-linux:
	@if ! command -v dprint >/dev/null 2>&1; then \
		echo "dprint not detected. Installing it using: curl -fsSL https://dprint.dev/install.sh | sh"; \
		exit 1; \
	fi

## Install and tidy Go module dependencies
#
# Downloads and installs all Go module dependencies and removes
# unused modules. Safe to run multiple times.
.PHONY: install
install:
	go mod tidy

# ============================================================================
# BUILDING
# ============================================================================

## Build production CLI binary optimized for deployment
#
# Creates an optimized binary in dist/github-actions-utils-cli suitable for
# distribution. This is the binary users will download and use.
.PHONY: build
build:
	mkdir -p dist
	go build -o dist/github-actions-utils-cli ./cmd/cli

## Build Linux binaries for Docker image
#
# Builds static Linux binaries for amd64 architecture.
# These binaries are used in the Docker multi-platform build.
.PHONY: build-linux
build-linux:
	mkdir -p dist
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags "-s -w -extldflags '-static'" \
		-a -installsuffix cgo \
		-o dist/github-actions-utils-cli-linux-amd64 \
		./cmd/cli

# ============================================================================
# DOCKER
# ============================================================================

## Build Docker image locally
#
# Builds a Docker image for local testing.
# Requires Linux binaries to be built first (make build-linux).
# Image is tagged as github-actions-utils-cli:latest
.PHONY: docker-build
docker-build: build-linux
	docker buildx build \
		--platform linux/amd64 \
		-t github-actions-utils-cli:latest \
		--load \
		.

## Run Docker container interactively
#
# Runs the Docker image with an interactive shell.
# Useful for testing the CLI within the container environment.
.PHONY: docker-run
docker-run:
	docker run --rm -it github-actions-utils-cli:latest --help

## Test Docker image
#
# Builds and tests the Docker image by running version check.
# Verifies that the image works correctly.
.PHONY: docker-test
docker-test: docker-build
	docker run --rm github-actions-utils-cli:latest --version

# ============================================================================
# DEVELOPMENT & RUNNING
# ============================================================================

## Build and run the CLI
#
# Compiles the CLI binary and then runs it directly.
# Useful for quick local testing of CLI changes.
#
# This will build the CLI, then execute the resulting binary.
# Pass additional arguments after '--' if desired:
#   make run -- [arguments...]
.PHONY: run
run: build
	@echo "Building and running the GitHub Actions Utils CLI..."
	@echo ""
	./dist/github-actions-utils-cli $(ARGS)

## Run MCP server for agent integration
#
# Starts the MCP (Model Context Protocol) server that exposes CLI
# functionality as tools over stdin/stdout. This allows AI agents
# to interact with GitHub Actions programmatically.
#
# The server will run until manually stopped with Ctrl+C or EOF.
# Connect an MCP client to test the integration.
.PHONY: mcp
mcp: build
	@echo "Starting MCP server..."
	@echo "Press Ctrl+C to stop the server."
	@echo ""
	./dist/github-actions-utils-cli mcp

# ============================================================================
# TESTING & QUALITY ASSURANCE
# ============================================================================

## Run all tests in the project
#
# Executes all unit tests, integration tests, and benchmarks.
# Tests are run with Go's built-in testing framework.
#
# Use 'go test -v ./...' for verbose output.
# Use 'go test -race ./...' to check for race conditions.
.PHONY: test
test:
	go test ./...

## Run comprehensive static analysis and security checks
#
# Performs multiple code quality checks:
# - go vet: Examines Go source code for suspicious constructs
# - staticcheck: Advanced static analysis for bugs and performance issues
# - govulncheck: Scans for known security vulnerabilities
#
# Fix any issues reported before committing code.
.PHONY: analyze
analyze:
	go vet ./...
	go tool staticcheck ./...
	go tool govulncheck ./...

## Format code and organize imports
#
# Automatically formats all code in the project:
# - go mod tidy: Cleans up module dependencies
# - go fmt: Formats Go source code to standard style
# - dprint fmt: Formats other files (JSON, YAML, etc.) using dprint
#
# Run this before committing to ensure consistent code style.
.PHONY: format
format:
	go mod tidy
	go fmt ./...
	dprint fmt

# ============================================================================
# MAINTENANCE
# ============================================================================

## Update all dependencies to latest compatible versions
#
# Updates all Go module dependencies to their latest minor/patch versions
# while respecting semantic versioning constraints. After updating:
# - Dependencies are updated to latest compatible versions
# - Code is automatically formatted
# - Module files are tidied
#
# Review changes carefully before committing dependency updates.
.PHONY: upgrade-deps
upgrade-deps:
	go get -u ./...
	$(MAKE) format

# ============================================================================
# HELP & DOCUMENTATION
# ============================================================================

## Show this help message with all available commands
#
# Displays a formatted list of all available make targets with descriptions.
# Commands are organized by topic for easy navigation.
.PHONY: help
help:
	@echo "=============================================="
	@echo "🔧 GITHUB ACTIONS UTILS CLI - COMMANDS"
	@echo "=============================================="
	@echo ""
	@awk 'BEGIN { desc = ""; target = "" } \
	/^## / { desc = substr($$0, 4) } \
	/^\.PHONY: / && desc != "" { \
		target = $$2; \
		printf "\033[36m%-20s\033[0m %s\n", target, desc; \
		desc = ""; target = "" \
	}' $(MAKEFILE_LIST)
	@echo ""
	@echo "💡 Use 'make <command>' to run any command above."
	@echo "📖 For detailed information, see comments in the Makefile."
	@echo ""
