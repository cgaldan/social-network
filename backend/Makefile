.PHONY: help build run test clean docker-build docker-run deps \
	migrate-up migrate-down migrate-down-one migrate-version

help:
	@echo "Available commands:"
	@echo "  make build              - Build the application"
	@echo "  make run                - Run the application"
	@echo "  make test               - Run tests"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make docker-build       - Build Docker image"
	@echo "  make docker-run         - Run Docker container"
	@echo "  make deps               - Install dependencies"
	@echo "  make migrate-up         - Apply all pending migrations"
	@echo "  make migrate-down       - Roll back ALL migrations (destructive)"
	@echo "  make migrate-down-one   - Roll back the most recent migration"
	@echo "  make migrate-version    - Print current schema version"

build:
	@echo "Building..."
	@go build -o bin/forum-backend ./cmd/server

run:
	@echo "Running..."
	@go run ./cmd/server/main.go

test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

docker-build:
	@echo "Building Docker image..."
	@docker build -t forum-backend:latest .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8000:8000 --name forum-backend forum-backend:latest

deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

migrate-up:
	@go run ./cmd/migrate up

migrate-down:
	@go run ./cmd/migrate down

migrate-down-one:
	@go run ./cmd/migrate down 1

migrate-version:
	@go run ./cmd/migrate version