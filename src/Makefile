# Simple Makefile for a Go project

# Build the application
all: build test

# Create the OpenAPI Documentations
doc:
	@echo "Generating docs..."
	@swag init -g internal/server/routes.go

build:
	@echo "Building..."
	
	
	@go build -o main cmd/main.go

# Run the application
run:
	@go run cmd/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch: doc
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@cd .. && docker build -t goauth:latest .

docker-run:
	@echo "Running Docker container..."
	@cd .. && docker-compose up -d

docker-stop:
	@echo "Stopping Docker containers..."
	@cd .. && docker-compose down

docker-logs:
	@echo "Showing Docker logs..."
	@cd .. && docker-compose logs -f app

docker-rebuild:
	@echo "Rebuilding and restarting Docker containers..."
	@cd .. && docker-compose down && docker-compose up -d --build

# Linting
lint:
	@echo "Running linters..."
	@golangci-lint run

# Code coverage
coverage:
	@echo "Running tests with coverage..."
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed successfully"

# Help command
help:
	@echo "Available commands:"
	@echo "  make all           - Build and test the application"
	@echo "  make build         - Build the application binary"
	@echo "  make run           - Run the application"
	@echo "  make test          - Run tests"
	@echo "  make coverage      - Run tests with coverage report"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make doc           - Generate Swagger documentation"
	@echo "  make watch         - Run with live reload"
	@echo "  make lint          - Run linters"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run with Docker Compose"
	@echo "  make docker-stop   - Stop Docker containers"
	@echo "  make docker-logs   - Show Docker logs"
	@echo "  make docker-rebuild- Rebuild and restart containers"
	@echo "  make install-tools - Install development tools"

.PHONY: all build run test clean watch doc docker-build docker-run docker-stop docker-logs docker-rebuild lint coverage install-tools help
