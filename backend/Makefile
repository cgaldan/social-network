.PHONY: help build run test clean docker-build docker-run migrate migrate-down seed seed-reset seed-scale

SCALE ?= 1000

help:
	@echo "Available commands:"
	@echo "  make build                - Build the application"
	@echo "  make run                  - Run the application"
	@echo "  make test                 - Run tests"
	@echo "  make clean                - Clean build artifacts"
	@echo "  make docker-build         - Build Docker image"
	@echo "  make docker-run           - Run Docker container"
	@echo "  make deps                 - Install dependencies"
	@echo "  make seed                 - Apply curated seed dataset (no reset)"
	@echo "  make seed-reset           - Wipe DB then apply curated seed dataset"
	@echo "  make seed-scale [SCALE=N] - Wipe DB, apply curated seed, then add N procedural users (default: 1000)"

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

seed:
	@echo "Seeding database with curated dataset..."
	@go run ./cmd/seed

seed-reset:
	@echo "Resetting and seeding database with curated dataset..."
	@go run ./cmd/seed -reset

seed-scale:
	@echo "Resetting and seeding with curated dataset + $(SCALE) procedural users..."
	@go run ./cmd/seed -reset -scale $(SCALE)