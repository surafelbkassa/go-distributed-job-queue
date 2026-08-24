package main

import (
	"net/http"

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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		job, err := usecase.EnqueueJob(body.Name, body.Payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Job enqueued", "job": job})

	})
	r.GET("/job/:id", func(c *gin.Context) {
		id := c.Param("id")
		job, err := usecase.GetJob(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"job": job})
	})
	r.Run(":8080")
}
