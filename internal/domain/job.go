package domain

import (
	"errors"
	"time"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

var ErrNotFound = errors.New("Job not found")

type Job struct {
	ID          string
	Type        string
	Name        string
	Status      Status
	Attempts    int
	MaxAttempts int
	Payload     []byte
	LastError   string
	CreatedAt   time.Time
}
