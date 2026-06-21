package main

import (
	"fmt"

	infrastructure "github.com/surafelbkassa/go-distributed-job-queue/Infrastructure"
	repository "github.com/surafelbkassa/go-distributed-job-queue/Repository"

	usecases "github.com/surafelbkassa/go-distributed-job-queue/Usecases"
)

func main() {
	redisClient := infrastructure.NewRedisClient()
	jobRepo := repository.NewRedisJobRepository(redisClient, "jobs")
	jobUsecases := usecases.NewJobUsecase(jobRepo)
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
