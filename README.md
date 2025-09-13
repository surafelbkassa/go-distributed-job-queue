# Go Distributed Job Queue

A distributed job queue system built with **Go** and **Redis**, designed for running background tasks at scale.  

## ✨ Features
- Enqueue and process jobs asynchronously
- Worker pool with configurable concurrency
- Job retries and failure handling
- REST API for submitting and monitoring jobs
- Real-time job status tracking (pending, running, completed, failed)
- Horizontal scalability with Redis as the central broker

## 🏗️ Architecture
- **Go (Goroutines & channels):** Worker pool + job handling
- **Redis:** Central queue for job persistence & distribution
- **REST API (Gin/Fiber):** Interface for job submission and monitoring
- **Optional Frontend (React):** Real-time job dashboard

## 🚀 Getting Started
1. Clone repo:  
   ```sh
   git clone https://github.com/<your-username>/go-distributed-job-queue.git
   cd go-distributed-job-queue
````

2. Start Redis (Docker recommended):

   ```sh
   docker run -d -p 6379:6379 redis
   ```
3. Run the app:

   ```sh
   go run cmd/server/main.go
   ```
4. Submit a job:

   ```sh
   curl -X POST http://localhost:8080/jobs -d '{"task":"send_email","payload":{"to":"user@example.com"}}'
   ```

## 🧪 Example Use Cases

* Sending bulk emails or notifications
* Image/video processing
* Data pipelines & ETL
* Scheduled background tasks

## 🛠️ Tech Stack

* Go
* Redis
* Gin (HTTP server)
* Docker (for local setup)

## 📌 Roadmap

* [ ] Add delayed/scheduled jobs
* [ ] Add distributed workers (multi-node support)
* [ ] Add metrics and observability
* [ ] Add dashboard UI for job monitoring(Optional)

---

## 📜 License

MIT


# ✅ Endpoints Checklist

### 🔐 Auth (JWT-based)

* [ ] `POST /auth/register` – Register new user (optional if internal use only)
* [ ] `POST /auth/login` – Login, receive access & refresh token
* [ ] `POST /auth/refresh` – Refresh JWT access token
* [ ] `POST /auth/logout` – Invalidate refresh token

---

### 📦 Jobs

* [ ] `POST /jobs` – Submit a new job (payload + metadata)
* [ ] `GET /jobs` – List all jobs (with filters: status, worker, time)
* [ ] `GET /jobs/:id` – Get details of a specific job
* [ ] `DELETE /jobs/:id` – Cancel a pending job (if not yet processed)
* [ ] `POST /jobs/:id/retry` – Retry a failed job manually

---

### ⚡ Workers

* [ ] `GET /workers` – List all active workers + their status
* [ ] `GET /workers/:id` – Worker details (uptime, jobs processed, failures)
* [ ] `POST /workers/register` – Worker registers itself (if not auto-discovery)
* [ ] `POST /workers/:id/stop` – Gracefully stop worker from taking new jobs

---

### 📊 Metrics & Monitoring

* [ ] `GET /metrics` – Global system metrics (jobs/sec, avg latency, success/fail %)
* [ ] `GET /metrics/jobs` – Aggregated job metrics (by status, type, worker)
* [ ] `GET /metrics/workers` – Worker performance metrics (processed count, failure rate)
* [ ] `GET /health` – Health check for API + Redis + DB

---

### Dashboard (Optional)

* [ ] `GET /dashboard/overview` – Aggregated data for dashboard cards (pending, running, failed, avg time)
* [ ] `GET /dashboard/timeline` – Recent jobs timeline for visualization
* [ ] `GET /dashboard/errors` – List of last N job errors
---


