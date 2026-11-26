package worker

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/domain"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/usecase"
)

func RegisterDemo(reg *usecase.Registry, rdb *redis.Client) {
	// Idempotent handler: checks a record keyed by job ID before doing the effect.
	reg.Register("email", func(ctx context.Context, j *domain.Job) error {
		rdb.HIncrBy(ctx, "executions", j.ID, 1) // counts runs, used by the chaos test
		if rdb.Exists(ctx, "sent:"+j.ID).Val() == 1 {
			return nil // an earlier attempt already did the work
		}
		select {
		case <-time.After(300 * time.Millisecond): // pretend to call the provider
		case <-ctx.Done():
			return ctx.Err()
		}
		return rdb.Set(ctx, "sent:"+j.ID, 1, 24*time.Hour).Err()
	})

	// Always fails: shows retries with growing gaps, then the dead-letter queue.
	reg.Register("always_fail", func(ctx context.Context, j *domain.Job) error {
		return errors.New("simulated failure")
	})
}
