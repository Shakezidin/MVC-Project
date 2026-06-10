SHELL := bash

.PHONY: build run test lint clean docker-up docker-down migrate seed seed-go swagger tidy

APP_NAME := bank-server
BUILD_DIR := bin
MAIN_PKG := ./cmd/server

# DB CONFIG (can override via environment variables)
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= bankdb

# IMPORTANT: set password for psql
export PGPASSWORD := $(DB_PASSWORD)

# If psql is not in PATH, uncomment and fix this:
# PSQL := "/c/Program Files/PostgreSQL/18/bin/psql"
PSQL := psql


# -----------------------
# BUILD
# -----------------------
build: 
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PKG)

run: build
	@echo "Running $(APP_NAME)..."
	./$(BUILD_DIR)/$(APP_NAME)


# -----------------------
# DATABASE MIGRATIONS
# -----------------------
migrate:
	@echo "Running migrations..."
	@for %%f in (migrations\*.up.sql) do \
		"C:\Program Files\PostgreSQL\18\bin\psql.exe" -h localhost -p 5432 -U postgres -d bankdb -f "%%f"


# -----------------------
# SEED DATA
# -----------------------
seed:
	@echo "Seeding database..."
	$(PSQL) -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f migrations/seed.sql
	@echo "Seeding completed."


seed-go:
	go run ./cmd/seed


# -----------------------
# TESTING
# -----------------------
test:
	go test -v -race -cover ./...

lint:
	go vet ./...
	gofmt -l .


# -----------------------
# CLEAN
# -----------------------
clean:
	rm -rf $(BUILD_DIR)
	go clean


tidy:
	go mod tidy


# -----------------------
# SWAGGER
# -----------------------
swagger:
	swag init -g cmd/server/main.go -o docs


# -----------------------
# DOCKER
# -----------------------
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f app


MCP_APP_NAME := bank-mcp-server
MCP_BUILD_DIR := mcp_bin
MCP_MAIN_FILE := ./cmd/mcp/.

mcp_run:
	@echo "Starting MCP Server..."
	@echo $(MCP_MAIN_FILE)
	go run $(MCP_MAIN_FILE)/main.go

mcp_build:
	@echo "Building MCP Server..."
	@mkdir -p $(MCP_BUILD_DIR)
	go build -o $(MCP_BUILD_DIR)/$(MCP_APP_NAME) $(MCP_MAIN_FILE)

mcp_dev:
	@echo "Running MCP Server in development mode..."
	air

mcp_tidy:
	@echo "Tidying go modules..."
	go mod tidy

mcp_clean:
	@echo "Cleaning build files..."
	rm -rf $(MCP_BUILD_DIR)