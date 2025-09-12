package api

import (
	"github.com/gin-gonic/gin"
	"github.com/surafelbkassa/go-distributed-job-queue/config"
	"github.com/surafelbkassa/go-distributed-job-queue/queue"
)

// SetupRoutes configures all API routes
func SetupRoutes(redisQueue *queue.RedisQueue, cfg *config.Config) *gin.Engine {
	// Create Gin router
	router := gin.Default()
	
	// Create handlers
	handlers := NewHandlers(redisQueue, cfg)
	
	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	
	// Health check
	router.GET("/health", handlers.HealthCheck)
	
	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Job management
		v1.POST("/jobs", handlers.SubmitJob)
		v1.POST("/jobs/bulk", handlers.BulkSubmitJobs)
		v1.GET("/jobs", handlers.ListJobs)
		v1.GET("/jobs/:id", handlers.GetJobStatus)
		v1.GET("/jobs/:id/events", handlers.JobStatusSSE)
		
		// Queue management
		v1.GET("/queue/stats", handlers.GetQueueStats)
	}
	
	return router
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
}