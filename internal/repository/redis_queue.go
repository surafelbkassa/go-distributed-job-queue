package repository

import (
	"context"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	domain "github.com/surafelbkassa/go-distributed-job-queue/internal/domain"
)

const keyPending = "queue:pending"

type RedisQueue struct {
	rdb *redis.Client
}

func NewRedisQueue(rdb *redis.Client) *RedisQueue {
	return &RedisQueue{rdb: rdb}
}

var _ domain.Queue = (*RedisQueue)(nil)

func jobKey(id string) string { return "jobs:" + id }

func (q *RedisQueue) Enqueue(ctx context.Context, j domain.Job) error {
	now := time.Now().Unix()
	pipe := q.rdb.TxPipeline()
	pipe.HSet(ctx, jobKey(j.ID), map[string]any{
		"id":           j.ID,
		"type":         j.Type,
		"name":         j.Name,
		"payload":      string(j.Payload),
		"status":       string(domain.StatusPending),
		"attempts":     0,
		"max_attempts": j.MaxAttempts,
		"last_error":   "",
		"created_at":   now,
		"updated_at":   now,
	})
	pipe.LPush(ctx, keyPending, j.ID)
	_, err := pipe.Exec(ctx)
	return err
}

// func (q *RedisQueue) Dequeue(ctx context.Context, timeout time.Duration) (*domain.Job, error) {
// 	res, err := q.rdb.BRPop(ctx, timeout, "queue:pending").Result()
// 	if errors.Is(err, redis.Nil) {
// 		return nil, nil
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	return q.Get(ctx, res[1])

// }

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
func (q *RedisQueue) Claim(ctx context.Context, lease time.Duration) (*domain.Job, error) {
	id, err := claimScript.Run(ctx, q.rdb, []string{keyPending, keyProcessing}, int(lease.Seconds())).Text()
	if errors.Is(err, redis.Nil) {
		return nil, nil // queue empty, not an error
	}
	if err != nil {
		return nil, err
	}
	return q.Get(ctx, id)
}

func (q *RedisQueue) Ack(ctx context.Context, id string) (bool, error) {
	n, err := ackScript.Run(ctx, q.rdb, []string{keyProcessing}, id).Int()
	return n == 1, err
}

func (q *RedisQueue) Fail(ctx context.Context, id, errMsg string, backoff time.Duration) (domain.Outcome, error) {
	secs := int(math.Ceil(backoff.Seconds()))
	res, err := failScript.Run(ctx, q.rdb,
		[]string{keyProcessing, keyDelayed, keyDead}, id, errMsg, secs).Text()
	return domain.Outcome(res), err
}

func (q *RedisQueue) ReapExpired(ctx context.Context, limit int) ([]string, error) {
	return reapScript.Run(ctx, q.rdb, []string{keyProcessing}, limit).StringSlice()
}

func (q *RedisQueue) PromoteDue(ctx context.Context, limit int) (int, error) {
	return promoteScript.Run(ctx, q.rdb, []string{keyDelayed, keyPending}, limit).Int()
}
