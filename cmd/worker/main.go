package main

import (
	"context"
	"fmt"
	"time"

	infrastructure "github.com/surafelbkassa/go-distributed-job-queue/Infrastructure"
)

func main() {
	client := infrastructure.NewRedisClient()
	ctx := context.Background()
	queueName := "job_queue"
	fmt.Println("Worker live and waiting for jobs...")

	for {
		results, err := client.BLPop(ctx, 0, queueName).Result()
		if err != nil {
			fmt.Printf("Error fetching job: %v\n", err)
			continue
		}

		jobId := results[1]
		fmt.Printf("Processing job ID: %s\n", jobId)

		time.Sleep(2 * time.Second)
		fmt.Printf("✅Finished processing job ID: %s\n", jobId)
	}
}
