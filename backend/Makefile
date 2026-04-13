# =============================================================================
# 太积堂充值与门店管理系统 - Makefile
# =============================================================================

.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } ' $(MAKEFILE_LIST)

# =============================================================================
# Environment
# =============================================================================

include .env
export

# =============================================================================
# Development
# =============================================================================

dev: compose-up ## Start local dev environment with docker-compose
	@echo "Starting Air with live reload..."
	@air

build: ## Build the application
	go build -o bin/server ./cmd/server

run: ## Run the application locally (requires mysql/redis)
	go run ./cmd/server

# =============================================================================
# Dependencies
# =============================================================================

deps: ## Download and tidy dependencies
	go mod tidy
	go mod verify

lint: golangci-lint ## Run linter
	@echo "Running golangci-lint..."
	golangci-lint run ./...

fix: ## Auto-fix code style issues
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .
	gci write .

# =============================================================================
# Testing
# =============================================================================

# Coverage threshold (percentage)
COVERAGE_MIN ?= 50

test: ## Run unit tests (skip integration tests)
	go test -v -race -short ./...

test-all: ## Run all tests including integration tests
	go test -v -race -short=false ./...

test-coverage: test ## Run tests with coverage report
	@echo "Coverage report:"
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Open coverage.html for detailed coverage"

test-threshold: ## Run tests with coverage threshold check
	@echo "Running tests with minimum coverage threshold: $(COVERAGE_MIN)%"
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./... 2>&1; \
	ACTUAL=$$(go tool cover -func=coverage.txt | grep total | awk -F'[ %]+' '{print $$3}'); \
	echo "Current coverage: $$ACTUAL%"; \
	if [ "$$(echo "$$ACTUAL < $(COVERAGE_MIN)" | bc)" -eq 1 ]; then \
		echo "ERROR: Coverage ($$ACTUAL%) is below threshold ($(COVERAGE_MIN)%)"; \
		exit 1; \
	fi

# =============================================================================
# Integration Tests
# =============================================================================

TEST_DSN ?= root:root123@tcp(localhost:3307)/test_db?parseTime=true&charset=utf8mb4

integration-test-up: ## Start integration test database
	cd tests/integration && docker-compose up -d

integration-test-down: ## Stop integration test database
	cd tests/integration && docker-compose down

integration-test: integration-test-up ## Run integration tests
	@echo "Waiting for database to be ready..."
	@sleep 10
	TEST_DSN=$(TEST_DSN) go test -v -tags=integration ./tests/integration/...
	@echo "Stopping test database..."
	cd tests/integration && docker-compose down

integration-test-clean: ## Clean integration test database
	docker exec sale-mysql-test mysql -uroot -proot123 -e "DROP DATABASE IF EXISTS test_db; CREATE DATABASE test_db;"

# =============================================================================
# Database Migration (GORM AutoMigrate)
# =============================================================================

migrate: ## Run GORM AutoMigrate (creates/updates all tables)
	go run ./migrations/

migrate-seed: ## Run migration with seed data (same as migrate, seed is included)
	go run ./migrations/

migrate-reset: ## Reset database: drop old tables and re-migrate
	@echo "Dropping all application tables..."
	@docker-compose exec -T mysql mysql -u$(DB_USER) -p$(DB_PASSWORD) $(DB_NAME) -e "\
		DROP TABLE IF EXISTS card_transactions, store_cards, c_recharges, recharge_applications, \
		recharge_operators, recharge_centers, users, casbin_rule;"
	go run ./migrations/

# =============================================================================
# Docker
# =============================================================================

compose-up: ## Start docker-compose services
	docker-compose up -d
	@echo "Waiting for MySQL to be ready..."
	@sleep 5
	@docker-compose exec -T mysql mysqladmin ping -h localhost -u root -p$(MYSQL_ROOT_PASSWORD) || true

compose-down: ## Stop docker-compose services
	docker-compose down

compose-logs: ## View docker-compose logs
	docker-compose logs -f

docker-build: ## Build Docker image
	docker build -t sale-backend:latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env sale-backend:latest

# =============================================================================
# CI/CD
# =============================================================================

pre-commit: lint test ## Run all pre-commit checks

ci-build: ## CI build step
	go mod download
	go build -o bin/server ./cmd/server

ci-test: ## CI test step
	go test -v -race ./...

ci-lint: ## CI lint step
	golangci-lint run ./...

# =============================================================================
# Tools
# =============================================================================

tools: ## Install required tools
	@echo "Installing tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/google/addlicense@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# =============================================================================
# Clean
# =============================================================================

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.txt coverage.html

clean-data: compose-down ## Clean docker volumes (WARNING: deletes data)
	docker-compose down -v
	rm -rf bin/

# =============================================================================
# Generate
# =============================================================================

swag: ## Generate Swagger docs
	swag init -g cmd/server/main.go -o docs/swagger

license: ## Add license header to files
	addlicense -c "太积堂" -y $$(date +%Y) .

.DEFAULT_GOAL := help
