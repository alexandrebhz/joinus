.PHONY: help build run test docker-up docker-down migrate backend-build frontend-build

help:
	@echo "Available commands:"
	@echo "  make build          - Build both backend and frontend"
	@echo "  make backend-build  - Build the backend"
	@echo "  make frontend-build - Build the frontend"
	@echo "  make run            - Run backend locally"
	@echo "  make test           - Run tests"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make migrate       - Run database migrations"

build: backend-build frontend-build

backend-build:
	cd backend && go build -o ../bin/crawler ./cmd/crawler

frontend-build:
	cd frontend && npm run build

run:
	cd backend && go run ./cmd/crawler

test:
	cd backend && go test -v ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-rebuild:
	docker-compose up -d --build

migrate:
	psql $(DATABASE_URL) -f backend/migrations/001_initial_schema.sql

.DEFAULT_GOAL := help
