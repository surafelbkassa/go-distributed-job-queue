package queue

import (
	"github.com/surafelbkassa/go-distributed-job-queue/internal/job"
)

// Queue interface defines the operations for a job queue
type Queue interface {
	// Enqueue adds a job to the queue
	Enqueue(job *job.Job) error
	
	// Dequeue removes and returns a job from the queue
	Dequeue() (*job.Job, error)
	
	// UpdateJob updates a job's status and metadata
	UpdateJob(job *job.Job) error
	
	// GetJob retrieves a job by ID
	GetJob(jobID string) (*job.Job, error)
	
	// GetJobsByStatus retrieves jobs by status
	GetJobsByStatus(status job.Status, limit int) ([]*job.Job, error)
	
	// DeleteJob removes a job from the queue
	DeleteJob(jobID string) error
	
	// GetQueueSize returns the number of jobs in the queue
	GetQueueSize() (int64, error)
	
	// GetStats returns queue statistics
	GetStats() (*Stats, error)
	
	// Close closes the queue connection
	Close() error
}

// Stats represents queue statistics
type Stats struct {
	TotalJobs     int64 `json:"total_jobs"`
	PendingJobs   int64 `json:"pending_jobs"`
	RunningJobs   int64 `json:"running_jobs"`
	CompletedJobs int64 `json:"completed_jobs"`
	FailedJobs    int64 `json:"failed_jobs"`
}