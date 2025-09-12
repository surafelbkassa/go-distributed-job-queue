# Distributed Job Queue

A distributed job queue system built with Go, Redis, and REST APIs. Supports job scheduling, worker management, retries, and real-time monitoring, designed for scalable background task execution.

## Features

- **Asynchronous Job Processing**: Enqueue and process jobs asynchronously
- **Worker Pool**: Configurable concurrency with goroutine-based workers
- **Retry Mechanism**: Automatic job retries with exponential backoff
- **REST API**: Complete API for job submission and monitoring
- **Real-time Status Tracking**: Live job status updates (pending, running, completed, failed)
- **Horizontal Scalability**: Redis as central broker for distributed processing
- **Graceful Shutdown**: Proper cleanup and shutdown handling

## Architecture

- **Go**: Goroutines & channels for worker pool and job handling
- **Redis**: Central queue for job persistence and distribution
- **Gin**: REST API framework for job submission and monitoring
- **Docker**: Containerized deployment with Docker Compose

## Quick Start

### Prerequisites

- Go 1.19+
- Redis (or use Docker Compose)

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/surafelbkassa/go-distributed-job-queue.git
cd go-distributed-job-queue
```

2. Start the services:
```bash
docker-compose up -d
```

This will start both Redis and the job queue application.

### Manual Setup

1. Start Redis:
```bash
redis-server
```

2. Clone and build the application:
```bash
git clone https://github.com/surafelbkassa/go-distributed-job-queue.git
cd go-distributed-job-queue
go mod tidy
go build -o job-queue-server .
```

3. Run the application:
```bash
./job-queue-server
```

## Configuration

Configure the application using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `REDIS_URL` | `redis://localhost:6379` | Redis connection URL |
| `REDIS_DB` | `0` | Redis database number |
| `SERVER_PORT` | `8080` | HTTP server port |
| `WORKER_COUNT` | `5` | Number of worker goroutines |
| `MAX_RETRIES` | `3` | Default maximum retries for jobs |
| `QUEUE_NAME` | `jobs` | Redis queue name |
| `STATUS_PREFIX` | `job:status:` | Redis key prefix for job status |

## API Endpoints

### Submit a Job

```bash
POST /api/v1/jobs
Content-Type: application/json

{
  "type": "email",
  "payload": {
    "recipient": "user@example.com",
    "subject": "Test Email"
  },
  "max_retries": 3
}
```

### Get Job Status

```bash
GET /api/v1/jobs/{job_id}
```

### Real-time Job Status (Server-Sent Events)

```bash
GET /api/v1/jobs/{job_id}/events
```

### Bulk Job Submission

```bash
POST /api/v1/jobs/bulk
Content-Type: application/json

[
  {
    "type": "email",
    "payload": {
      "recipient": "user1@example.com"
    }
  },
  {
    "type": "data_processing",
    "payload": {
      "data_size": 100.5
    }
  }
]
```

### Queue Statistics

```bash
GET /api/v1/queue/stats
```

### Health Check

```bash
GET /health
```

## Job Types

The system supports different job types with custom processing logic:

- **`email`**: Email sending jobs
- **`data_processing`**: Data processing tasks
- **`report_generation`**: Report generation jobs
- **Custom types**: Any custom job type with generic processing

## Job Status Lifecycle

1. **`pending`**: Job submitted and waiting in queue
2. **`running`**: Job picked up by a worker and being processed
3. **`completed`**: Job finished successfully
4. **`failed`**: Job failed after all retry attempts

## Example Usage

### Submit an Email Job

```bash
curl -X POST http://localhost:8080/api/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "payload": {
      "recipient": "user@example.com",
      "subject": "Welcome Email",
      "body": "Welcome to our service!"
    },
    "max_retries": 3
  }'
```

### Check Job Status

```bash
curl http://localhost:8080/api/v1/jobs/{job_id}
```

### Monitor Queue

```bash
curl http://localhost:8080/api/v1/queue/stats
```

## Development

### Running Tests

```bash
go test -v ./...
```

### Building

```bash
go build -o job-queue-server .
```

### Running Locally

```bash
# Start Redis
redis-server

# Run the application
./job-queue-server
```

## Production Considerations

1. **Redis Persistence**: Configure Redis with appropriate persistence settings
2. **Monitoring**: Add metrics collection and monitoring
3. **Logging**: Implement structured logging
4. **Authentication**: Add API authentication and authorization
5. **Rate Limiting**: Implement rate limiting for API endpoints
6. **Scaling**: Deploy multiple instances behind a load balancer

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
