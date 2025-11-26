.PHONY: help build run test docker-up docker-down docker-logs docker-rebuild dev migrate-up migrate-down lint

help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make test          - Run tests"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make docker-logs   - View API container logs"
	@echo "  make docker-rebuild - Rebuild and start containers"
	@echo "  make dev           - Start dev services (postgres, minio)"
	@echo "  make lint          - Run linter"

build:
	go build -o bin/server cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f api

docker-rebuild:
	docker-compose up -d --build

dev:
	docker-compose -f docker-compose.dev.yml up -d

lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.DEFAULT_GOAL := help


