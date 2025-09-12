# Go Distributed Job Queue Makefile

.PHONY: build run server worker clean test docker-build docker-up docker-down example

# Build all binaries
build:
	@echo "Building binaries..."
	@mkdir -p bin
	@go build -o bin/job-queue .
	@go build -o bin/server ./cmd/server
	@go build -o bin/worker ./cmd/worker
	@go build -o bin/client-example ./examples/client_example.go
	@echo "Build complete!"

# Run the all-in-one application
run: build
	@echo "Starting job queue (server + workers)..."
	@./bin/job-queue

# Run only the server
server: build
	@echo "Starting server only..."
	@./bin/server

# Run only workers
worker: build
	@echo "Starting worker only..."
	@./bin/worker

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@go clean

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t job-queue .

# Start with Docker Compose
docker-up:
	@echo "Starting with Docker Compose..."
	@docker-compose up --build

# Stop Docker Compose
docker-down:
	@echo "Stopping Docker Compose..."
	@docker-compose down

# Run client example
example: build
	@echo "Running client example..."
	@./bin/client-example

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@go vet ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Install
install:
	@echo "Installing job-queue..."
	@go install .

# Help
help:
	@echo "Available commands:"
	@echo "  build       - Build all binaries"
	@echo "  run         - Run all-in-one application"
	@echo "  server      - Run server only"
	@echo "  worker      - Run worker only"
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run tests"
	@echo "  docker-build- Build Docker image"
	@echo "  docker-up   - Start with Docker Compose"
	@echo "  docker-down - Stop Docker Compose"
	@echo "  example     - Run client example"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  deps        - Download dependencies"
	@echo "  install     - Install binary"
	@echo "  help        - Show this help"