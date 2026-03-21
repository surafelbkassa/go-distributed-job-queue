FROM golang:1.21-alpine

# Create a folder inside the container to hold our app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the API and the Worker into executable binary files
RUN go build -o api_bin ./cmd/api/main.go
RUN go build -o worker_bin ./cmd/worker/main.go

# This file doesn't "run" anything yet. 
# Our docker-compose.yml will tell it which binary to start.