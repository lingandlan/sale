# =============================================================================
# Marketplace API - Makefile
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

test: ## Run unit tests
	go test -v -race -coverprofile=coverage.txt ./...

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

integration-test: compose-up ## Run integration tests (requires docker-compose)
	go clean -testcache
	go test -v ./tests/integration/...

# =============================================================================
# Database Migration
# =============================================================================

migrate-create: ## Create a new migration: make migrate-create name=create_users
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	goose -dir migrations mysql "$(MYSQL_DSN)" create $(name) sql

migrate-up: ## Run migrations
	goose -dir migrations mysql "$(MYSQL_DSN)" up

migrate-down: ## Rollback last migration
	goose -dir migrations mysql "$(MYSQL_DSN)" down

migrate-status: ## Show migration status
	goose -dir migrations mysql "$(MYSQL_DSN)" status

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
	docker build -t marketplace-api:latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env marketplace-api:latest

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
	go install github.com/pressly/goose/v3/cmd/goose@latest
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
	addlicense -c "Marketplace" -y $$(date +%Y) .

.DEFAULT_GOAL := help
