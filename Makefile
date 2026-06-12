.PHONY: help build run test test-unit test-integration lint fmt clean docker-up docker-down docker-logs migrate seed seed-go swagger tidy deps \
mcp-build mcp-run mcp-tidy mcp-clean clean-temp

# Default target
help:
	@echo "Available targets:"
	@echo "  help              - Show this help message"
	@echo "  build             - Build binary"
	@echo "  run               - Build and run"
	@echo "  test              - Run all tests"
	@echo "  test-unit         - Unit tests only"
	@echo "  test-integration  - Integration tests"
	@echo "  lint              - Run go vet"
	@echo "  fmt               - Format code with gofmt"
	@echo "  tidy              - Tidy go modules"
	@echo "  deps              - Install dependencies"
	@echo "  swagger           - Regenerate Swagger docs"
	@echo "  clean             - Clean build files"
	@echo "  clean-temp        - Clean temporary build files"
	@echo "  docker-up         - Start Docker services"
	@echo "  docker-down       - Stop Docker services"
	@echo "  docker-logs       - Show Docker logs"
	@echo "  migrate           - Run SQL migrations"
	@echo "  seed              - Seed database with SQL"
	@echo "  seed-go           - Seed database with Go program"
	@echo "  mcp-build         - Build MCP server"
	@echo "  mcp-run           - Build and run MCP server"
	@echo "  mcp-tidy          - Tidy MCP server modules"
	@echo "  mcp-clean         - Clean MCP build files"

# Configuration
APP_NAME := bank-server
BUILD_DIR := bin
MAIN_PKG := ./cmd/server

MCP_APP_NAME := bank-mcp-server
MCP_BUILD_DIR := mcp_bin
MCP_MAIN_PKG := ./cmd/mcp

# Check if gcp-key.json exists
ifeq ($(OS),Windows_NT)
    GCP_KEY_EXISTS := $(if $(wildcard $(CURDIR)\gcp-key.json),yes,no)
    GCP_KEY := $(CURDIR)\gcp-key.json
else
    GCP_KEY_EXISTS := $(if $(wildcard $(CURDIR)/gcp-key.json),yes,no)
    GCP_KEY := $(CURDIR)/gcp-key.json
endif

# OS-specific commands
ifeq ($(OS),Windows_NT)
    RM := if exist
    RMDIR := rmdir /s /q
    MKDIR := mkdir
    EXE := .exe
    ENV_SET := set
    PATHSEP := \
else
    RM := rm -f
    RMDIR := rm -rf
    MKDIR := mkdir -p
    EXE :=
    ENV_SET := export
    PATHSEP := /
endif

# Default environment variables
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= bankdb

# Clean temporary files
clean-temp:
	@echo "Cleaning temporary files..."
ifeq ($(OS),Windows_NT)
	@if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR)
	@if exist $(MCP_BUILD_DIR) rmdir /s /q $(MCP_BUILD_DIR)
else
	@rm -rf $(BUILD_DIR) $(MCP_BUILD_DIR)
endif

# Build main server
build: clean-temp
	@echo "Building $(APP_NAME)..."
	@$(MKDIR) $(BUILD_DIR)
	go build -o $(BUILD_DIR)$(PATHSEP)$(APP_NAME)$(EXE) $(MAIN_PKG)

# Run main server
run: clean-temp
	@echo "Running $(APP_NAME)..."
ifeq ($(OS),Windows_NT)
	@if exist "$(GCP_KEY)" (set GOOGLE_APPLICATION_CREDENTIALS=$(GCP_KEY) && go run $(MAIN_PKG)) else (go run $(MAIN_PKG))
else
	@if [ -f "$(GCP_KEY)" ]; then GOOGLE_APPLICATION_CREDENTIALS=$(GCP_KEY) go run $(MAIN_PKG); else go run $(MAIN_PKG); fi
endif

# Build MCP server
mcp-build: clean-temp
	@echo "Building MCP server..."
	@$(MKDIR) $(MCP_BUILD_DIR)
	go build -o $(MCP_BUILD_DIR)$(PATHSEP)$(MCP_APP_NAME)$(EXE) $(MCP_MAIN_PKG)

# Run MCP server
mcp-run: clean-temp
	@echo "Starting MCP server..."
ifeq ($(OS),Windows_NT)
	@if exist "$(GCP_KEY)" (set GOOGLE_APPLICATION_CREDENTIALS=$(GCP_KEY) && go run $(MCP_MAIN_PKG)) else (go run $(MCP_MAIN_PKG))
else
	@if [ -f "$(GCP_KEY)" ]; then GOOGLE_APPLICATION_CREDENTIALS=$(GCP_KEY) go run $(MCP_MAIN_PKG); else go run $(MCP_MAIN_PKG); fi
endif

# Tidy MCP server modules
mcp-tidy:
	@echo "Tidying MCP server modules..."
	go mod tidy

# Clean MCP build files
mcp-clean:
	@echo "Cleaning MCP build files..."
	-$(RMDIR) $(MCP_BUILD_DIR)

# Run all tests
test:
	@echo "Running all tests..."
	go test -v -race -cover ./...

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test -v -short ./internal/...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v ./tests/integration/...

# Lint code
lint:
	@echo "Linting code..."
	go vet ./...

# Format code
fmt:
	@echo "Formatting code..."
	gofmt -w .

# Tidy modules
tidy:
	@echo "Tidying Go modules..."
	go mod tidy

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download

# Generate swagger docs
swagger:
	@echo "Generating swagger docs..."
	swag init -g cmd/server/main.go -o docs

# Clean build files
clean: clean-temp
	@echo "Cleaning project..."
	go clean

# Docker targets
docker-up:
	@echo "Starting docker compose..."
	docker compose up -d

docker-down:
	@echo "Stopping docker compose..."
	docker compose down

docker-logs:
	@echo "Showing docker logs..."
	docker compose logs -f app

# Database targets
migrate:
	@echo "Running migrations..."

seed:
	@echo "Seeding database..."

seed-go:
	@echo "Seeding database with Go program..."
	go run ./cmd/seed
