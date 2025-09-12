package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/surafelbkassa/go-distributed-job-queue/config"
	"github.com/surafelbkassa/go-distributed-job-queue/models"
	"github.com/surafelbkassa/go-distributed-job-queue/queue"
)

// Handlers contains the API handlers
type Handlers struct {
	queue  *queue.RedisQueue
	config *config.Config
}

// NewHandlers creates a new handlers instance
func NewHandlers(redisQueue *queue.RedisQueue, cfg *config.Config) *Handlers {
	return &Handlers{
		queue:  redisQueue,
		config: cfg,
	}
}

// JobRequest represents a job submission request
type JobRequest struct {
	Type       string                 `json:"type" binding:"required"`
	Payload    map[string]interface{} `json:"payload"`
	MaxRetries *int                   `json:"max_retries,omitempty"`
}

// JobResponse represents a job response
type JobResponse struct {
	Job     *models.Job `json:"job"`
	Message string      `json:"message,omitempty"`
}

// QueueStatsResponse represents queue statistics
type QueueStatsResponse struct {
	QueueLength int64  `json:"queue_length"`
	QueueName   string `json:"queue_name"`
}

// SubmitJob handles job submission
func (h *Handlers) SubmitJob(c *gin.Context) {
	var req JobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Set default max retries if not provided
	maxRetries := h.config.MaxRetries
	if req.MaxRetries != nil {
		maxRetries = *req.MaxRetries
	}

	// Create new job
	job := models.NewJob(req.Type, req.Payload, maxRetries)

	// Enqueue the job
	if err := h.queue.Enqueue(job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to enqueue job",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, JobResponse{
		Job:     job,
		Message: "Job submitted successfully",
	})
}

// GetJobStatus retrieves job status by ID
func (h *Handlers) GetJobStatus(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Job ID is required",
		})
		return
	}

	job, err := h.queue.GetJobStatus(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Job not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, JobResponse{
		Job: job,
	})
}

// GetQueueStats retrieves queue statistics
func (h *Handlers) GetQueueStats(c *gin.Context) {
	length, err := h.queue.GetQueueLength()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get queue statistics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, QueueStatsResponse{
		QueueLength: length,
		QueueName:   h.config.QueueName,
	})
}

// HealthCheck handles health check requests
func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"queue":  h.config.QueueName,
	})
}

// ListJobs handles listing jobs (simplified version - in production, implement pagination)
func (h *Handlers) ListJobs(c *gin.Context) {
	// This is a simplified implementation
	// In production, you'd implement proper pagination and filtering
	c.JSON(http.StatusOK, gin.H{
		"message": "Job listing not implemented in this demo",
		"hint":    "Use GET /jobs/{id} to get specific job status",
	})
}

// JobStatusSSE provides Server-Sent Events for real-time job status updates
func (h *Handlers) JobStatusSSE(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Job ID is required",
		})
		return
	}

	// Check if job exists
	_, err := h.queue.GetJobStatus(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Job not found",
			"details": err.Error(),
		})
		return
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Create a channel for this specific client
	clientChan := make(chan queue.StatusUpdate, 10)
	defer close(clientChan)

	// Get status updates from the queue
	statusChan := h.queue.GetStatusChannel()

	// Send initial job status
	if job, err := h.queue.GetJobStatus(jobID); err == nil {
		c.SSEvent("status", gin.H{
			"job_id": job.ID,
			"status": job.Status,
			"updated_at": job.UpdatedAt,
		})
		c.Writer.Flush()
	}

	// Listen for updates
	for {
		select {
		case <-c.Request.Context().Done():
			return
		case update := <-statusChan:
			if update.JobID == jobID {
				c.SSEvent("status", gin.H{
					"job_id": update.JobID,
					"status": update.Status,
				})
				c.Writer.Flush()

				// If job is completed or failed, end the stream
				if update.Status == models.StatusCompleted || update.Status == models.StatusFailed {
					return
				}
			}
		}
	}
}

// BulkSubmitJobs handles bulk job submission
func (h *Handlers) BulkSubmitJobs(c *gin.Context) {
	var requests []JobRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	if len(requests) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No jobs provided",
		})
		return
	}

	// Limit bulk submissions
	maxBulkSize := 100
	if len(requests) > maxBulkSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Too many jobs in bulk request",
			"limit": maxBulkSize,
		})
		return
	}

	var jobs []*models.Job
	var errors []string

	for i, req := range requests {
		// Set default max retries if not provided
		maxRetries := h.config.MaxRetries
		if req.MaxRetries != nil {
			maxRetries = *req.MaxRetries
		}

		// Create new job
		job := models.NewJob(req.Type, req.Payload, maxRetries)

		// Enqueue the job
		if err := h.queue.Enqueue(job); err != nil {
			errors = append(errors, "Job "+strconv.Itoa(i)+": "+err.Error())
		} else {
			jobs = append(jobs, job)
		}
	}

	response := gin.H{
		"submitted": len(jobs),
		"jobs":      jobs,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	status := http.StatusCreated
	if len(errors) > 0 && len(jobs) == 0 {
		status = http.StatusInternalServerError
	} else if len(errors) > 0 {
		status = http.StatusPartialContent
	}

	c.JSON(status, response)
}