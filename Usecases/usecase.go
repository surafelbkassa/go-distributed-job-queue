package usecases

import (
	"fmt"
	"time"

	domain "github.com/surafelbkassa/go-distributed-job-queue/Domain"
	repository "github.com/surafelbkassa/go-distributed-job-queue/Repository"
)

type JobUsecase struct {
	JobRepo repository.JobRepository
}

func NewJobUsecase(jobRepo repository.JobRepository) *JobUsecase {
	return &JobUsecase{JobRepo: jobRepo}
}

func (uc *JobUsecase) EnqueueJob(name, payload string) (*domain.Job, error) {
	job := domain.NewJob(name, payload)
	err := uc.JobRepo.EnqueueJob(job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (uc *JobUsecase) ProcessJob() error {

	job, err := uc.JobRepo.Dequeue()
	if err != nil {
		return err
	}
	err = uc.JobRepo.UpdateStatus(job.ID, domain.StatusInProgress)
	if err != nil {
		return err
	}
	fmt.Printf("Processing task %s payload %s at hand", job.Name, job.Payload)
	time.Sleep(2 * time.Second)
	err = uc.JobRepo.UpdateStatus(job.ID, domain.StatusCompleted)
	if err != nil {
		return err
	}
	return nil
}
