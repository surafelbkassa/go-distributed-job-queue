package main

import (
	"fmt"

	infrastructure "github.com/surafelbkassa/go-distributed-job-queue/internal/infrastructure"
	repository "github.com/surafelbkassa/go-distributed-job-queue/internal/repository"

	"github.com/surafelbkassa/go-distributed-job-queue/internal/usecase"
)

func main() {
	redisClient := infrastructure.NewRedisClient()
	jobRepo := repository.NewRedisJobRepository(redisClient, "jobs")
	jobUsecases := usecase.NewJobUsecase(jobRepo)
	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			fmt.Printf("Worker %d started", workerID)
			for {
				err := jobUsecases.ProcessJob()
				if err != nil {
					fmt.Printf("Error processing job: %v\n", err)
				}

			}
		}(i)
	}
	select {}
}
