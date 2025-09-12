package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/surafelbkassa/go-distributed-job-queue/api"
	"github.com/surafelbkassa/go-distributed-job-queue/config"
	"github.com/surafelbkassa/go-distributed-job-queue/queue"
)

func main() {
	log.Println("Starting Distributed Job Queue Server...")

	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("Configuration loaded: Workers=%d, Port=%s, Redis=%s", 
		cfg.WorkerCount, cfg.ServerPort, cfg.RedisURL)

	// Initialize Redis queue
	redisQueue, err := queue.NewRedisQueue(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis queue: %v", err)
	}
	defer redisQueue.Close()

	// Create worker pool
	workerPool := queue.NewWorkerPool(redisQueue, nil, cfg)
	
	// Start worker pool
	workerPool.Start()
	defer workerPool.Stop()

	// Setup API routes
	router := api.SetupRoutes(redisQueue, cfg)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Print API endpoints
	printAPIEndpoints(cfg.ServerPort)

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

// printAPIEndpoints prints available API endpoints
func printAPIEndpoints(port string) {
	baseURL := "http://localhost:" + port
	
	log.Print("\n=== API Endpoints ===\n")
	log.Printf("Health Check:          GET    %s/health", baseURL)
	log.Printf("Submit Job:            POST   %s/api/v1/jobs", baseURL)
	log.Printf("Bulk Submit Jobs:      POST   %s/api/v1/jobs/bulk", baseURL)
	log.Printf("Get Job Status:        GET    %s/api/v1/jobs/{id}", baseURL)
	log.Printf("Job Status Events:     GET    %s/api/v1/jobs/{id}/events", baseURL)
	log.Printf("Queue Statistics:      GET    %s/api/v1/queue/stats", baseURL)
	log.Println("=====================")
	
	log.Println("Example job submission:")
	log.Printf(`curl -X POST %s/api/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "payload": {
      "recipient": "user@example.com",
      "subject": "Test Email"
    },
    "max_retries": 3
  }'`, baseURL)
}