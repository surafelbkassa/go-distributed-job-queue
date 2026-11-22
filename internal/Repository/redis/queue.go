package repository

import (
	"context"
	"errors"
	"strconv"
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
		"attempts":     0,
		"max_attempts": j.MaxAttempts,
		"created_at":   time.Now().Unix(),
	})
	pipe.LPush(ctx, "queue:pending", j.ID)
	_, err := pipe.Exec(ctx)
	return err
}

func (q *RedisQueue) Dequeue(ctx context.Context, timeout time.Duration) (*domain.Job, error) {
	res, err := q.rdb.BRPop(ctx, timeout, "queue:pending").Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return q.Get(ctx, res[1])

}

func (q *RedisQueue) Get(ctx context.Context, id string) (*domain.Job, error) {
	m, err := q.rdb.HGetAll(ctx, jobKey(id)).Result()
	if err != nil {
		return nil, err
	}
	if len(m) == 0 {
		return nil, domain.ErrNotFound
	}
	attempts, _ := strconv.Atoi(m["attempts"])
	maxAttempts, _ := strconv.Atoi(m["max_attempts"])
	created, _ := strconv.ParseInt(m["created_at"], 10, 64)
	return &domain.Job{
		ID: m["id"], Type: m["type"], Payload: []byte(m["payload"]),
		Status: domain.Status(m["status"]), Attempts: attempts,
		MaxAttempts: maxAttempts, LastError: m["last_error"],
		CreatedAt: time.Unix(created, 0),
	}, nil

}
