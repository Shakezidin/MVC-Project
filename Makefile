.PHONY: build run test lint clean docker-up docker-down migrate seed swagger tidy

APP_NAME := bank-server
BUILD_DIR := bin
MAIN_PKG := ./cmd/server

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PKG)

run: build
	@echo "Running $(APP_NAME)..."
	./$(BUILD_DIR)/$(APP_NAME)

test:
	@echo "Running tests..."
	go test -v -race -cover ./...

test-unit:
	@echo "Running unit tests..."
	go test -v -short ./internal/...

test-integration:
	@echo "Running integration tests..."
	go test -v -tags=integration ./tests/...

lint:
	@echo "Running linter..."
	go vet ./...
	gofmt -l .

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean

tidy:
	go mod tidy

swagger:
	@echo "Generating Swagger docs..."
	swag init -g cmd/server/main.go -o docs

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f app

migrate:
	bash scripts/migrate.sh

seed:
	bash scripts/seed.sh

seed-go:
	go run ./cmd/seed

dev: docker-up
	@echo "Waiting for services..."
	sleep 5
	$(MAKE) run
