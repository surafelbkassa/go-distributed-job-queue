package main

import (
	"github.com/gin-gonic/gin"
	infrastructure "github.com/surafelbkassa/go-distributed-job-queue/Infrastructure"
	repository "github.com/surafelbkassa/go-distributed-job-queue/Repository"
	usecases "github.com/surafelbkassa/go-distributed-job-queue/Usecases"
)

func main() {
	r := gin.Default()

	redisClient := infrastructure.NewRedisClient()
	repo := repository.NewRedisJobRepository(redisClient, "jobs")
	usecase := usecases.NewJobUsecase(repo)
	r.POST("/enqueue", func(c *gin.Context) {
		var body struct {
			Name    string `json:"name"`
			Payload string `json:"payload"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		job, err := usecase.EnqueueJob(body.Name, body.Payload)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Job enqueued", "job": job})

	})
}
