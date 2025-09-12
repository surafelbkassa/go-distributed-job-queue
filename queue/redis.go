package queue

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/surafelbkassa/go-distributed-job-queue/config"
	"github.com/surafelbkassa/go-distributed-job-queue/models"
)

// RedisQueue implements a Redis-backed job queue
type RedisQueue struct {
	client       *redis.Client
	config       *config.Config
	ctx          context.Context
	pubsub       *redis.PubSub
	statusChan   chan StatusUpdate
}

// StatusUpdate represents a job status update
type StatusUpdate struct {
	JobID  string           `json:"job_id"`
	Status models.JobStatus `json:"status"`
}

// NewRedisQueue creates a new Redis queue instance
func NewRedisQueue(cfg *config.Config) (*RedisQueue, error) {
	// Parse Redis URL
	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}
	
	opts.DB = cfg.RedisDB
	
	client := redis.NewClient(opts)
	
	ctx := context.Background()
	
	// Test connection
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	
	// Subscribe to status updates
	pubsub := client.Subscribe(ctx, "job:status:*")
	
	rq := &RedisQueue{
		client:     client,
		config:     cfg,
		ctx:        ctx,
		pubsub:     pubsub,
		statusChan: make(chan StatusUpdate, 100),
	}
	
	// Start listening for status updates
	go rq.listenForStatusUpdates()
	
	return rq, nil
}

// Enqueue adds a job to the queue
func (rq *RedisQueue) Enqueue(job *models.Job) error {
	// Convert job to JSON
	jobData, err := job.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize job: %w", err)
	}
	
	// Add to Redis list (queue)
	err = rq.client.LPush(rq.ctx, rq.config.QueueName, jobData).Err()
	if err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}
	
	// Store job status
	err = rq.storeJobStatus(job)
	if err != nil {
		return fmt.Errorf("failed to store job status: %w", err)
	}
	
	log.Printf("Job %s enqueued successfully", job.ID)
	return nil
}

// Dequeue removes and returns a job from the queue
func (rq *RedisQueue) Dequeue(timeout time.Duration) (*models.Job, error) {
	// Block pop from Redis list
	result, err := rq.client.BRPop(rq.ctx, timeout, rq.config.QueueName).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // No job available
		}
		return nil, fmt.Errorf("failed to dequeue job: %w", err)
	}
	
	if len(result) < 2 {
		return nil, fmt.Errorf("invalid dequeue result")
	}
	
	// Parse job data
	job, err := models.FromJSON(result[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse job data: %w", err)
	}
	
	return job, nil
}

// GetJobStatus retrieves the status of a job by ID
func (rq *RedisQueue) GetJobStatus(jobID string) (*models.Job, error) {
	key := rq.config.StatusPrefix + jobID
	data, err := rq.client.Get(rq.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("failed to get job status: %w", err)
	}
	
	job, err := models.FromJSON(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse job data: %w", err)
	}
	
	return job, nil
}

// UpdateJobStatus updates the status of a job
func (rq *RedisQueue) UpdateJobStatus(job *models.Job) error {
	err := rq.storeJobStatus(job)
	if err != nil {
		return err
	}
	
	// Publish status update
	channel := "job:status:" + job.ID
	err = rq.client.Publish(rq.ctx, channel, fmt.Sprintf(`{"job_id":"%s","status":"%s"}`, job.ID, job.Status)).Err()
	if err != nil {
		log.Printf("Failed to publish status update for job %s: %v", job.ID, err)
	}
	
	return nil
}

// GetQueueLength returns the number of jobs in the queue
func (rq *RedisQueue) GetQueueLength() (int64, error) {
	length, err := rq.client.LLen(rq.ctx, rq.config.QueueName).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue length: %w", err)
	}
	return length, nil
}

// GetStatusChannel returns the channel for status updates
func (rq *RedisQueue) GetStatusChannel() <-chan StatusUpdate {
	return rq.statusChan
}

// Close closes the Redis connection
func (rq *RedisQueue) Close() error {
	if rq.pubsub != nil {
		rq.pubsub.Close()
	}
	close(rq.statusChan)
	return rq.client.Close()
}

// storeJobStatus stores job status in Redis
func (rq *RedisQueue) storeJobStatus(job *models.Job) error {
	key := rq.config.StatusPrefix + job.ID
	jobData, err := job.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize job: %w", err)
	}
	
	// Store with TTL of 24 hours
	err = rq.client.Set(rq.ctx, key, jobData, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store job status: %w", err)
	}
	
	return nil
}

// listenForStatusUpdates listens for status updates from Redis pub/sub
func (rq *RedisQueue) listenForStatusUpdates() {
	ch := rq.pubsub.Channel()
	for msg := range ch {
		if strings.HasPrefix(msg.Channel, "job:status:") {
			jobID := strings.TrimPrefix(msg.Channel, "job:status:")
			// Extract status from message payload
			// This is a simplified implementation
			status := extractStatusFromPayload(msg.Payload)
			
			select {
			case rq.statusChan <- StatusUpdate{JobID: jobID, Status: status}:
			default:
				// Channel full, drop update
				log.Printf("Status update channel full, dropping update for job %s", jobID)
			}
		}
	}
}

// extractStatusFromPayload extracts status from JSON payload
func extractStatusFromPayload(payload string) models.JobStatus {
	// Simple extraction - in production, use proper JSON parsing
	if strings.Contains(payload, "pending") {
		return models.StatusPending
	} else if strings.Contains(payload, "running") {
		return models.StatusRunning
	} else if strings.Contains(payload, "completed") {
		return models.StatusCompleted
	} else if strings.Contains(payload, "failed") {
		return models.StatusFailed
	}
	return models.StatusPending
}