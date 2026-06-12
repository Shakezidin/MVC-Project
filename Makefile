SHELL := cmd.exe

.PHONY: build run test lint clean docker-up docker-down migrate seed seed-go swagger tidy \
mcp_run mcp_build mcp_dev mcp_tidy mcp_clean

# ======================================================
# MAIN APP
# ======================================================

APP_NAME := bank-server
BUILD_DIR := bin
MAIN_PKG := ./cmd/server

# ======================================================
# MCP SERVER
# ======================================================

MCP_APP_NAME := bank-mcp-server
MCP_BUILD_DIR := mcp_bin
MCP_MAIN_FILE := ./cmd/mcp

# ======================================================
# GCP CONFIG
# ======================================================

GCP_KEY := $(CURDIR)\gcp-key.json

# ======================================================
# DATABASE CONFIG
# ======================================================

DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= bankdb

export PGPASSWORD := $(DB_PASSWORD)

PSQL := "C:\Program Files\PostgreSQL\18\bin\psql.exe"

# ======================================================
# CLEAN TEMP FILES
# ======================================================

clean-temp:
	@if exist bin rmdir /s /q bin
	@if exist mcp_bin rmdir /s /q mcp_bin
	@if exist -p rmdir /s /q -p

# ======================================================
# BUILD MAIN SERVER
# ======================================================

build: clean-temp
	@echo Building $(APP_NAME)...
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	go build -o $(BUILD_DIR)\$(APP_NAME).exe $(MAIN_PKG)

run: build
	@echo Running $(APP_NAME)...
	set GOOGLE_APPLICATION_CREDENTIALS=$(GCP_KEY) && $(BUILD_DIR)\$(APP_NAME).exe

# ======================================================
# MCP SERVER
# ======================================================

mcp_build: clean-temp
	@echo Building MCP Server...
	@if not exist $(MCP_BUILD_DIR) mkdir $(MCP_BUILD_DIR)
	go build -o $(MCP_BUILD_DIR)\$(MCP_APP_NAME).exe $(MCP_MAIN_FILE)

mcp_run: mcp_build
	@echo Starting MCP Server...
	set GOOGLE_APPLICATION_CREDENTIALS=$(GCP_KEY) && $(MCP_BUILD_DIR)\$(MCP_APP_NAME).exe

mcp_dev:
	@echo Running MCP Server in development mode...
	set GOOGLE_APPLICATION_CREDENTIALS=$(GCP_KEY) && air

mcp_tidy:
	@echo Tidying Go modules...
	go mod tidy

mcp_clean:
	@echo Cleaning MCP build files...
	@if exist $(MCP_BUILD_DIR) rmdir /s /q $(MCP_BUILD_DIR)

# ======================================================
# DATABASE MIGRATIONS
# ======================================================

migrate:
	@echo Running migrations...
	@for %%f in (migrations\*.up.sql) do $(PSQL) -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f "%%f"

# ======================================================
# SEED DATA
# ======================================================

seed:
	@echo Seeding database...
	$(PSQL) -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f migrations\seed.sql
	@echo Seeding completed.

seed-go:
	go run ./cmd/seed

# ======================================================
# TESTING
# ======================================================

test:
	go test -v -race -cover ./...

lint:
	go vet ./...
	gofmt -w .

# ======================================================
# CLEAN
# ======================================================

clean:
	@echo Cleaning project...
	@if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR)
	@if exist $(MCP_BUILD_DIR) rmdir /s /q $(MCP_BUILD_DIR)
	go clean

tidy:
	go mod tidy

# ======================================================
# SWAGGER
# ======================================================

swagger:
	swag init -g cmd/server/main.go -o docs

# ======================================================
# DOCKER
# ======================================================

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f app