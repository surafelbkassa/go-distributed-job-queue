package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
	domain "github.com/surafelbkassa/go-distributed-job-queue/Domain"
)

type UserRepository interface {
	Save(user *domain.User) error
}

func EnqueueJob(client *redis.Client, queueName string, jobID string) error {
	ctx := context.Background()
	return client.LPush(ctx, queueName, jobID).Err()
}
