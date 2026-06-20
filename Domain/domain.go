package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

type Job struct {
	ID        string
	Name      string
	Status    Status
	Payload   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewJob(name string, payload string) *Job {
	now := time.Now()
	return &Job{
		ID:        uuid.New().String(),
		Name:      name,
		Status:    StatusPending,
		Payload:   payload,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
