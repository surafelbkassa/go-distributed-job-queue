package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	domain "github.com/surafelbkassa/go-distributed-job-queue/internal/Domain"
)

type RedisQueue struct {
	rdb *redis.Client
}

func NewRedisClient(rdb *redis.Client) *RedisQueue {
	return &RedisQueue{rdb: rdb}
}

var _ domain.Queue = (*RedisQueue)(nil)

func jobKey(id string) string { return "jobs:" + id }

func (q *RedisQueue) Enqueue(ctx context.Context, j domain.Job) error {
	pipe := q.rdb.TxPipeline()
	pipe.HSet(ctx, jobKey(j.ID), map[string]any{
		"id":           j.ID,
		"type":         j.Type,
		"payload":      string(j.Payload),
		"status":       string(domain.StatusPending),
		"attemps":      0,
		"max_attempts": j.MaxAttempts,
		"created_at":   time.Now().Unix(),
	})
	pipe.LPush(ctx, "queue:pending", j.ID)
	_, err := pipe.Exec(ctx)
	return err
}

func (q *RedisQueue) Dequeue(ctx context.Context, timeout time.Duration) (*domain.Job, error) {}

func (q *RedisQueue) Get(ctx context.Context, id string) (*domain.Job, error) {}
