package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/job"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/queue"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/worker"
)

// Server represents the HTTP server
type Server struct {
	router   *gin.Engine
	handlers *Handlers
	host     string
	port     int
}

// NewServer creates a new HTTP server
func NewServer(host string, port int, q queue.Queue, registry *job.Registry, pool *worker.Pool) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	handlers := NewHandlers(q, registry, pool)

	server := &Server{
		router:   router,
		handlers: handlers,
		host:     host,
		port:     port,
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")
	{
		// Job management
		api.POST("/jobs", s.handlers.CreateJob)
		api.GET("/jobs", s.handlers.ListJobs)
		api.GET("/jobs/:id", s.handlers.GetJob)
		api.DELETE("/jobs/:id", s.handlers.DeleteJob)

		// Statistics and monitoring
		api.GET("/stats", s.handlers.GetStats)
		api.GET("/job-types", s.handlers.GetJobTypes)
		api.GET("/health", s.handlers.HealthCheck)
	}

	// Root endpoints
	s.router.GET("/health", s.handlers.HealthCheck)
	s.router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "go-distributed-job-queue",
			"version": "1.0.0",
			"status":  "running",
		})
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	log.Printf("Starting HTTP server on %s", addr)
	return s.router.Run(addr)
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}