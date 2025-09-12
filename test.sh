#!/bin/bash

# Simple test script to verify all functionality

echo "=== Go Distributed Job Queue Test ==="

# Run tests
echo "Running unit tests..."
go test -v ./...

# Build all binaries
echo "Building project..."
make build

# Check if binaries were created
echo "Checking binaries..."
if [ -f "bin/job-queue" ] && [ -f "bin/server" ] && [ -f "bin/worker" ] && [ -f "bin/client-example" ]; then
    echo "✅ All binaries built successfully"
else
    echo "❌ Some binaries are missing"
    exit 1
fi

# Test with Docker Compose (if Docker is available)
if command -v docker &> /dev/null && command -v docker-compose &> /dev/null; then
    echo "Testing with Docker Compose..."
    timeout 30s docker-compose up --build &
    sleep 20
    
    # Test API
    curl -f http://localhost:8080/api/v1/health &> /dev/null
    if [ $? -eq 0 ]; then
        echo "✅ Docker Compose setup works"
    else
        echo "❌ Docker Compose setup failed"
    fi
    
    docker-compose down &> /dev/null
else
    echo "Docker not available, skipping Docker tests"
fi

echo "=== Test Complete ==="