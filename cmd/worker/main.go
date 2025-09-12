package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/surafelbkassa/go-distributed-job-queue/internal/job"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/queue"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/worker"
	"github.com/surafelbkassa/go-distributed-job-queue/pkg/config"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Redis queue
	q, err := queue.NewRedisQueue(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to initialize Redis queue: %v", err)
	}
	defer q.Close()

	// Create job registry
	registry := job.NewRegistry()
	registerJobHandlers(registry)

	// Create and start worker pool
	pool := worker.NewPool(cfg.WorkerConcurrency, q, registry)
	
	log.Printf("Starting worker pool with %d workers", cfg.WorkerConcurrency)
	pool.Start()

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	log.Println("Worker started. Press Ctrl+C to exit.")
	<-c
	
	log.Println("Shutting down worker...")
	pool.Stop()
	log.Println("Worker stopped")
}

// registerJobHandlers registers the same job handlers as the server
func registerJobHandlers(registry *job.Registry) {
	// You can register the same handlers or different ones for different worker instances
	// This allows for specialization of workers for specific job types
	
	// Email job handler
	registry.Register("email", func(j *job.Job) error {
		log.Printf("Worker processing email job: %s", j.ID)
		// Email processing logic here
		return nil
	})

	// Image processing job handler
	registry.Register("image_resize", func(j *job.Job) error {
		log.Printf("Worker processing image resize job: %s", j.ID)
		// Image processing logic here
		return nil
	})

	// Report generation job handler
	registry.Register("generate_report", func(j *job.Job) error {
		log.Printf("Worker processing report generation job: %s", j.ID)
		// Report generation logic here
		return nil
	})

	log.Printf("Worker registered job handlers: %v", registry.GetRegisteredTypes())
}