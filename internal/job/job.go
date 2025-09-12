package job

import (
	"encoding/json"
	"time"
)

// Status represents the current state of a job
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusRetrying  Status = "retrying"
)

// Job represents a job in the queue
type Job struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Payload     map[string]interface{} `json:"payload"`
	Status      Status                 `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Attempts    int                    `json:"attempts"`
	MaxAttempts int                    `json:"max_attempts"`
	LastError   string                 `json:"last_error,omitempty"`
	Priority    int                    `json:"priority"`
	Delay       time.Duration          `json:"delay"`
}

// NewJob creates a new job with default values
func NewJob(jobType string, payload map[string]interface{}) *Job {
	now := time.Now()
	return &Job{
		ID:          generateJobID(),
		Type:        jobType,
		Payload:     payload,
		Status:      StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
		Attempts:    0,
		MaxAttempts: 3,
		Priority:    0,
		Delay:       0,
	}
}

// ToJSON serializes the job to JSON
func (j *Job) ToJSON() ([]byte, error) {
	return json.Marshal(j)
}

// FromJSON deserializes a job from JSON
func FromJSON(data []byte) (*Job, error) {
	var job Job
	err := json.Unmarshal(data, &job)
	return &job, err
}

// MarkRunning updates the job status to running
func (j *Job) MarkRunning() {
	now := time.Now()
	j.Status = StatusRunning
	j.UpdatedAt = now
	j.StartedAt = &now
}

// MarkCompleted updates the job status to completed
func (j *Job) MarkCompleted() {
	now := time.Now()
	j.Status = StatusCompleted
	j.UpdatedAt = now
	j.CompletedAt = &now
}

// MarkFailed updates the job status to failed
func (j *Job) MarkFailed(err error) {
	j.Status = StatusFailed
	j.UpdatedAt = time.Now()
	j.Attempts++
	if err != nil {
		j.LastError = err.Error()
	}
}

// MarkRetrying updates the job status to retrying
func (j *Job) MarkRetrying(err error) {
	j.Status = StatusRetrying
	j.UpdatedAt = time.Now()
	j.Attempts++
	if err != nil {
		j.LastError = err.Error()
	}
}

// CanRetry checks if the job can be retried
func (j *Job) CanRetry() bool {
	return j.Attempts < j.MaxAttempts
}

// ShouldDelay checks if the job should be delayed before processing
func (j *Job) ShouldDelay() bool {
	return j.Delay > 0 && time.Since(j.UpdatedAt) < j.Delay
}