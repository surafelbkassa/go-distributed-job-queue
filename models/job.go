package models

import (
	"encoding/json"
	"time"
)

// JobStatus represents the current status of a job
type JobStatus string

const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

// Job represents a job in the queue
type Job struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Payload     map[string]interface{} `json:"payload"`
	Status      JobStatus              `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	RetryCount  int                    `json:"retry_count"`
	MaxRetries  int                    `json:"max_retries"`
	Error       string                 `json:"error,omitempty"`
	Result      interface{}            `json:"result,omitempty"`
}

// NewJob creates a new job with default values
func NewJob(jobType string, payload map[string]interface{}, maxRetries int) *Job {
	now := time.Now()
	return &Job{
		ID:         generateJobID(),
		Type:       jobType,
		Payload:    payload,
		Status:     StatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
		RetryCount: 0,
		MaxRetries: maxRetries,
	}
}

// ToJSON converts job to JSON string
func (j *Job) ToJSON() (string, error) {
	data, err := json.Marshal(j)
	return string(data), err
}

// FromJSON creates job from JSON string
func FromJSON(data string) (*Job, error) {
	var job Job
	err := json.Unmarshal([]byte(data), &job)
	return &job, err
}

// CanRetry checks if the job can be retried
func (j *Job) CanRetry() bool {
	return j.RetryCount < j.MaxRetries
}

// MarkRunning updates job status to running
func (j *Job) MarkRunning() {
	j.Status = StatusRunning
	now := time.Now()
	j.StartedAt = &now
	j.UpdatedAt = now
}

// MarkCompleted updates job status to completed
func (j *Job) MarkCompleted(result interface{}) {
	j.Status = StatusCompleted
	j.Result = result
	now := time.Now()
	j.CompletedAt = &now
	j.UpdatedAt = now
}

// MarkFailed updates job status to failed
func (j *Job) MarkFailed(err error) {
	j.Status = StatusFailed
	if err != nil {
		j.Error = err.Error()
	} else {
		j.Error = "Unknown error"
	}
	j.UpdatedAt = time.Now()
	j.RetryCount++
}

// generateJobID generates a unique job ID
func generateJobID() string {
	// Simple implementation - in production, use UUID or similar
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}