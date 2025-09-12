package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/surafelbkassa/go-distributed-job-queue/internal/api"
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

	// Create worker pool
	pool := worker.NewPool(cfg.WorkerConcurrency, q, registry)

	// Start worker pool
	log.Printf("Starting worker pool with %d workers", cfg.WorkerConcurrency)
	pool.Start()

	// Create and start server
	server := api.NewServer(cfg.ServerHost, cfg.ServerPort, q, registry, pool)

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down...")
		pool.Stop()
		os.Exit(0)
	}()

	// Start server
	log.Printf("Job Queue Server starting on %s:%d", cfg.ServerHost, cfg.ServerPort)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// registerJobHandlers registers example job handlers
func registerJobHandlers(registry *job.Registry) {
	// Email job handler
	registry.Register("email", func(j *job.Job) error {
		log.Printf("Processing email job: %s", j.ID)
		
		// Simulate email processing
		to, ok := j.Payload["to"].(string)
		if !ok {
			return fmt.Errorf("missing 'to' field in email job")
		}
		
		subject, ok := j.Payload["subject"].(string)
		if !ok {
			return fmt.Errorf("missing 'subject' field in email job")
		}
		
		log.Printf("Sending email to: %s, subject: %s", to, subject)
		
		// Simulate work
		time.Sleep(2 * time.Second)
		
		log.Printf("Email sent successfully to %s", to)
		return nil
	})

	// Image processing job handler
	registry.Register("image_resize", func(j *job.Job) error {
		log.Printf("Processing image resize job: %s", j.ID)
		
		imageURL, ok := j.Payload["image_url"].(string)
		if !ok {
			return fmt.Errorf("missing 'image_url' field in image resize job")
		}
		
		width, ok := j.Payload["width"].(float64)
		if !ok {
			return fmt.Errorf("missing 'width' field in image resize job")
		}
		
		height, ok := j.Payload["height"].(float64)
		if !ok {
			return fmt.Errorf("missing 'height' field in image resize job")
		}
		
		log.Printf("Resizing image %s to %dx%d", imageURL, int(width), int(height))
		
		// Simulate work
		time.Sleep(3 * time.Second)
		
		log.Printf("Image resized successfully: %s", imageURL)
		return nil
	})

	// Report generation job handler
	registry.Register("generate_report", func(j *job.Job) error {
		log.Printf("Processing report generation job: %s", j.ID)
		
		reportType, ok := j.Payload["type"].(string)
		if !ok {
			return fmt.Errorf("missing 'type' field in report generation job")
		}
		
		log.Printf("Generating %s report", reportType)
		
		// Simulate work
		time.Sleep(5 * time.Second)
		
		log.Printf("Report generated successfully: %s", reportType)
		return nil
	})

	// Data processing job handler
	registry.Register("data_processing", func(j *job.Job) error {
		log.Printf("Processing data processing job: %s", j.ID)
		
		dataSource, ok := j.Payload["source"].(string)
		if !ok {
			return fmt.Errorf("missing 'source' field in data processing job")
		}
		
		log.Printf("Processing data from source: %s", dataSource)
		
		// Simulate work
		time.Sleep(4 * time.Second)
		
		log.Printf("Data processing completed for source: %s", dataSource)
		return nil
	})

	log.Printf("Registered job handlers: %v", registry.GetRegisteredTypes())
}