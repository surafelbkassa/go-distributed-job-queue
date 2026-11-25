package domain

import (
	"context"
	"time"
)

type Outcome string

const (
	OutcomeRetry Outcome = "retry"
	OutcomeDead  Outcome = "dead"
	OutcomeStale Outcome = "stale"
)

type Queue interface {
	Enqueue(ctx context.Context, j Job) error
	Claim(ctx context.Context, lease time.Duration) (*Job, error) // nil, nil when the queue is empty
	Ack(ctx context.Context, id string) (owned bool, err error)
	Fail(ctx context.Context, id, errMsg string, backoff time.Duration) (Outcome, error)
	ReapExpired(ctx context.Context, limit int) ([]string, error)
	PromoteDue(ctx context.Context, limit int) (int, error)
	Get(ctx context.Context, id string) (*Job, error)
}
