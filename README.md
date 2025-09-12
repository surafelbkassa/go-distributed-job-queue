# Go Distributed Job Queue

A distributed job queue system built with Go, Redis, and REST APIs. Supports job scheduling, worker management, retries, and real-time monitoring, designed for scalable background task execution.

## ✨ Features

- **Asynchronous Job Processing**: Enqueue and process jobs asynchronously with configurable concurrency
- **Worker Pool Management**: Scalable worker pool with configurable concurrency using goroutines and channels
- **Job Retries & Failure Handling**: Automatic retry mechanism with exponential backoff for failed jobs
- **REST API**: Complete API for job submission, monitoring, and management
- **Real-time Job Status Tracking**: Track job status (pending, running, completed, failed, retrying)
- **Horizontal Scalability**: Redis-based central broker enables multiple worker instances
- **Priority Queue Support**: Jobs can be prioritized for faster processing
- **Statistics & Monitoring**: Real-time queue and worker statistics

## 🏗️ Architecture

- **Go (Goroutines & Channels)**: Worker pool and job handling with efficient concurrency
- **Redis**: Central queue for job persistence, distribution, and coordination
- **Gin Framework**: REST API for job submission and monitoring
- **Modular Design**: Separate server and worker processes for flexible deployment

## 📋 Prerequisites

- Go 1.19 or higher
- Redis server running locally or accessible via URL
- Git (for cloning the repository)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/surafelbkassa/go-distributed-job-queue.git
cd go-distributed-job-queue
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Start Redis

Make sure Redis is running. You can start it locally:

```bash
# Using Docker
docker run -d -p 6379:6379 redis:alpine

# Or install and start Redis locally
redis-server
```

### 4. Run the Application

#### Option A: All-in-One (Server + Workers)

```bash
# Build and run
go build -o bin/job-queue .
./bin/job-queue
```

#### Option B: Separate Server and Worker Processes

Terminal 1 (Server):
```bash
go build -o bin/server ./cmd/server
./bin/server
```

Terminal 2 (Worker):
```bash
go build -o bin/worker ./cmd/worker
./bin/worker
```

### 5. Test the API

The server will start on `http://localhost:8080` by default.

#### Create a Job

```bash
curl -X POST http://localhost:8080/api/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "payload": {
      "to": "user@example.com",
      "subject": "Welcome!",
      "body": "Welcome to our service!"
    }
  }'
```

#### Get Job Status

```bash
curl http://localhost:8080/api/v1/jobs/{job_id}
```

#### List Jobs

```bash
# All jobs
curl http://localhost:8080/api/v1/jobs

# Jobs by status
curl http://localhost:8080/api/v1/jobs?status=completed&limit=10
```

#### Get Statistics

```bash
curl http://localhost:8080/api/v1/stats
```

## 🔧 Configuration

Configure the application using environment variables:

```bash
# Redis configuration
export REDIS_URL="redis://localhost:6379"

# Server configuration
export SERVER_HOST="localhost"
export SERVER_PORT=8080

# Worker configuration
export WORKER_CONCURRENCY=5

# Job configuration
export DEFAULT_MAX_ATTEMPTS=3
export DEFAULT_TIMEOUT="30s"
```

## 📖 API Reference

### Job Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/jobs` | Create a new job |
| GET | `/api/v1/jobs` | List jobs (with optional status filter) |
| GET | `/api/v1/jobs/:id` | Get job details |
| DELETE | `/api/v1/jobs/:id` | Delete a job |

### Monitoring

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/stats` | Get queue and worker statistics |
| GET | `/api/v1/job-types` | Get registered job types |
| GET | `/api/v1/health` | Health check |

### Job Creation Request

```json
{
  "type": "email",
  "payload": {
    "to": "user@example.com",
    "subject": "Hello World"
  },
  "max_attempts": 3,
  "priority": 1
}
```

### Job Response

```json
{
  "id": "job_1672531200_abc123",
  "type": "email",
  "payload": {
    "to": "user@example.com",
    "subject": "Hello World"
  },
  "status": "pending",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z",
  "attempts": 0,
  "max_attempts": 3,
  "priority": 1
}
```

## 💼 Supported Job Types

The system comes with example job handlers:

- **email**: Email sending jobs
- **image_resize**: Image processing jobs
- **generate_report**: Report generation jobs
- **data_processing**: Data processing jobs

### Adding Custom Job Types

Register new job handlers in your main application:

```go
registry.Register("my_custom_job", func(j *job.Job) error {
    // Your job processing logic here
    log.Printf("Processing custom job: %s", j.ID)
    
    // Access job payload
    data := j.Payload["data"].(string)
    
    // Simulate work
    time.Sleep(time.Second)
    
    return nil
})
```

## 📊 Monitoring & Statistics

The system provides comprehensive monitoring through the `/api/v1/stats` endpoint:

```json
{
  "queue": {
    "total_jobs": 150,
    "pending_jobs": 5,
    "running_jobs": 3,
    "completed_jobs": 140,
    "failed_jobs": 2
  },
  "pool": {
    "worker_count": 5,
    "is_running": true,
    "concurrency": 5
  }
}
```

## 🔄 Job Lifecycle

1. **Pending**: Job is created and waiting in the queue
2. **Running**: Job is being processed by a worker
3. **Completed**: Job finished successfully
4. **Failed**: Job failed after all retry attempts
5. **Retrying**: Job failed but will be retried

## 🚀 Deployment

### Docker Deployment

Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o job-queue .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/job-queue .
CMD ["./job-queue"]
```

### Kubernetes Deployment

Deploy separate server and worker pods for better scalability:

- Server pods handle API requests
- Worker pods process jobs from the queue
- Scale workers independently based on queue size

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Related Projects

- [Redis](https://redis.io/) - In-memory data structure store
- [Gin](https://gin-gonic.com/) - HTTP web framework for Go
- [Go-Redis](https://github.com/go-redis/redis) - Redis client for Go
