package repository

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	domain "github.com/surafelbkassa/go-distributed-job-queue/Domain"
)

type JobRepository interface {
	EnqueueJob(job *domain.Job) error
	Dequeue() (*domain.Job, error)
	UpdateStatus(id string, status domain.Status) error
}
type RedisJobRepository struct {
	client    *redis.Client
	queuename string
}

func (r *RedisJobRepository) EnqueueJob(Job *domain.Job) error {

	ctx := context.Background()
	data, err := json.Marshal(Job)
	if err != nil {
		return err
	}
	err = r.client.LPush(ctx, r.queuename, Job.ID).Err()
	if err != nil {
		return err
	}
	return r.client.Set(ctx, "job:"+Job.ID, data, 0).Err()
}

func (r *RedisJobRepository) Dequeue() (*domain.Job, error) {
	ctx := context.Background()
	results, err := r.client.BLPop(ctx, 0, r.queuename).Result()
	if err != nil {
		return nil, err
	}
	JobID := results[1]
	data, err := r.client.Get(ctx, "job:"+JobID).Bytes()
	if err != nil {
		return nil, err
	}
	var job domain.Job
	err = json.Unmarshal(data, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}
func (r *RedisJobRepository) UpdateStatus(id string, status domain.Status) error {
	ctx := context.Background()
	data, err := r.client.Get(ctx, "job:"+id).Bytes()
	if err != nil {
		return err
	}
	var job domain.Job
	err = json.Unmarshal(data, &job)
	job.Status = status
	updated, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, "job:"+id, updated, 0).Err()
}
