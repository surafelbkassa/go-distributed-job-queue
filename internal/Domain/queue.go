package domain

import (
	"context"
	"time"
)

type Queue interface {
	Enqueue(ctx context.Context, j Job) error
	Dequeue(ctx context.Context, timeout time.Duration) (*Job, error)
	Get(ctx context.Context, id string) (*Job, error)
}
