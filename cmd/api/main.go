package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	infrastructure "github.com/surafelbkassa/go-distributed-job-queue/Infrastructure"
)

func main() {
	r := gin.Default()

	redisClient := infrastructure.NewRedisClient()
	ctx := context.Background()
	r.POST("/enqueue", func(c *gin.Context) {
		jobID := c.Query("jobID")

		if jobID == "" {
			c.JSON(400, gin.H{"error": "jobID is required"})
			return
		}

		err := redisClient.LPush(ctx, "job_queue", jobID).Err()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue job"})
		}

		fmt.Println("Job sent to redis", jobID)

		c.JSON(http.StatusOK, gin.H{"message": "Job enqueued successfully", "jobID": jobID})
	})
	fmt.Println("API server running on 8080...")
	r.Run(":8080")
}
