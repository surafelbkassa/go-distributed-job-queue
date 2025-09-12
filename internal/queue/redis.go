package queue

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/job"
)

const (
	defaultJobQueue    = "job_queue"
	defaultJobHash     = "job_data"
	defaultStatsHash   = "job_stats"
	defaultTimeout     = 30 * time.Second
)

// RedisQueue implements the Queue interface using Redis
type RedisQueue struct {
	client      *redis.Client
	ctx         context.Context
	queueName   string
	jobHashName string
	statsHash   string
}

// NewRedisQueue creates a new Redis-based queue
func NewRedisQueue(redisURL string) (*RedisQueue, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)
	ctx := context.Background()

	// Test connection
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisQueue{
		client:      client,
		ctx:         ctx,
		queueName:   defaultJobQueue,
		jobHashName: defaultJobHash,
		statsHash:   defaultStatsHash,
	}, nil
}

// Enqueue adds a job to the queue
func (r *RedisQueue) Enqueue(j *job.Job) error {
	// Serialize job to JSON
	jobData, err := j.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize job: %w", err)
	}

	pipe := r.client.TxPipeline()

	// Store job data in hash
	pipe.HSet(r.ctx, r.jobHashName, j.ID, jobData)

	// Add job ID to queue based on priority
	if j.Priority > 0 {
		// Use sorted set for priority queue
		pipe.ZAdd(r.ctx, r.queueName+"_priority", &redis.Z{
			Score:  float64(j.Priority),
			Member: j.ID,
		})
	} else {
		// Use list for regular queue (FIFO)
		pipe.LPush(r.ctx, r.queueName, j.ID)
	}

	// Update stats
	pipe.HIncrBy(r.ctx, r.statsHash, "total_jobs", 1)
	pipe.HIncrBy(r.ctx, r.statsHash, string(j.Status), 1)

	_, err = pipe.Exec(r.ctx)
	return err
}

// Dequeue removes and returns a job from the queue
func (r *RedisQueue) Dequeue() (*job.Job, error) {
	// Try priority queue first
	jobID, err := r.client.ZPopMax(r.ctx, r.queueName+"_priority").Result()
	if err == nil && len(jobID) > 0 {
		return r.getJobByID(jobID[0].Member.(string))
	}

	// Try regular queue
	jobIDResult, err := r.client.BRPop(r.ctx, defaultTimeout, r.queueName).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // No jobs available
		}
		return nil, fmt.Errorf("failed to dequeue job: %w", err)
	}

	if len(jobIDResult) < 2 {
		return nil, fmt.Errorf("invalid dequeue result")
	}

	return r.getJobByID(jobIDResult[1])
}

// getJobByID retrieves a job by ID from Redis
func (r *RedisQueue) getJobByID(jobID string) (*job.Job, error) {
	jobData, err := r.client.HGet(r.ctx, r.jobHashName, jobID).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("job not found: %s", jobID)
		}
		return nil, fmt.Errorf("failed to get job data: %w", err)
	}

	return job.FromJSON([]byte(jobData))
}

// UpdateJob updates a job's status and metadata
func (r *RedisQueue) UpdateJob(j *job.Job) error {
	// Get old job for stats update
	oldJob, err := r.GetJob(j.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing job: %w", err)
	}

	// Serialize updated job
	jobData, err := j.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize job: %w", err)
	}

	pipe := r.client.TxPipeline()

	// Update job data
	pipe.HSet(r.ctx, r.jobHashName, j.ID, jobData)

	// Update stats
	if oldJob.Status != j.Status {
		pipe.HIncrBy(r.ctx, r.statsHash, string(oldJob.Status), -1)
		pipe.HIncrBy(r.ctx, r.statsHash, string(j.Status), 1)
	}

	_, err = pipe.Exec(r.ctx)
	return err
}

// GetJob retrieves a job by ID
func (r *RedisQueue) GetJob(jobID string) (*job.Job, error) {
	return r.getJobByID(jobID)
}

// GetJobsByStatus retrieves jobs by status
func (r *RedisQueue) GetJobsByStatus(status job.Status, limit int) ([]*job.Job, error) {
	// Get all job IDs and data
	allJobs, err := r.client.HGetAll(r.ctx, r.jobHashName).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs: %w", err)
	}

	var jobs []*job.Job
	count := 0

	for _, jobData := range allJobs {
		if limit > 0 && count >= limit {
			break
		}

		j, err := job.FromJSON([]byte(jobData))
		if err != nil {
			continue // Skip invalid jobs
		}

		if j.Status == status {
			jobs = append(jobs, j)
			count++
		}
	}

	return jobs, nil
}

// DeleteJob removes a job from the queue
func (r *RedisQueue) DeleteJob(jobID string) error {
	// Get job for stats update
	j, err := r.GetJob(jobID)
	if err != nil {
		return err
	}

	pipe := r.client.TxPipeline()

	// Remove from job hash
	pipe.HDel(r.ctx, r.jobHashName, jobID)

	// Remove from queues
	pipe.LRem(r.ctx, r.queueName, 0, jobID)
	pipe.ZRem(r.ctx, r.queueName+"_priority", jobID)

	// Update stats
	pipe.HIncrBy(r.ctx, r.statsHash, "total_jobs", -1)
	pipe.HIncrBy(r.ctx, r.statsHash, string(j.Status), -1)

	_, err = pipe.Exec(r.ctx)
	return err
}

// GetQueueSize returns the number of jobs in the queue
func (r *RedisQueue) GetQueueSize() (int64, error) {
	totalStr, err := r.client.HGet(r.ctx, r.statsHash, "total_jobs").Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	return strconv.ParseInt(totalStr, 10, 64)
}

// GetStats returns queue statistics
func (r *RedisQueue) GetStats() (*Stats, error) {
	statsData, err := r.client.HGetAll(r.ctx, r.statsHash).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	stats := &Stats{}

	if val, exists := statsData["total_jobs"]; exists {
		stats.TotalJobs, _ = strconv.ParseInt(val, 10, 64)
	}
	if val, exists := statsData[string(job.StatusPending)]; exists {
		stats.PendingJobs, _ = strconv.ParseInt(val, 10, 64)
	}
	if val, exists := statsData[string(job.StatusRunning)]; exists {
		stats.RunningJobs, _ = strconv.ParseInt(val, 10, 64)
	}
	if val, exists := statsData[string(job.StatusCompleted)]; exists {
		stats.CompletedJobs, _ = strconv.ParseInt(val, 10, 64)
	}
	if val, exists := statsData[string(job.StatusFailed)]; exists {
		stats.FailedJobs, _ = strconv.ParseInt(val, 10, 64)
	}

	return stats, nil
}

// Close closes the queue connection
func (r *RedisQueue) Close() error {
	return r.client.Close()
}