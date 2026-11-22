FROM golang:1.25-alpine AS build

# Create a folder inside the container to hold our app
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/api_bin ./cmd/api \&& CGO_ENABLED=0 go build -o /put/worker_bin ./cmd/worker

# Build the API and the Worker into executable binary files
# RUN go build -o api_bin ./cmd/api/main.go
# RUN go build -o worker_bin ./cmd/worker/main.go

# This file doesn't "run" anything yet. 
# Our docker-compose.yml will tell it which binary to start.
FROM gcr.io/distroless/static
COPY --from=build /out/app /app
ENTRYPOINT ["/src"]