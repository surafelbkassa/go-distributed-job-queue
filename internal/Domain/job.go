package domain

import "time"

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

type Job struct {
	ID          string
	Type        string
	Name        string
	Status      Status
	Attempts    int
	MaxAttempts int
	Payload     []byte
	CreatedAt   time.Time
}
