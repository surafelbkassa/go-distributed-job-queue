package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/job"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/queue"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/worker"
)

// Handlers contains all API handlers
type Handlers struct {
	queue    queue.Queue
	registry *job.Registry
	pool     *worker.Pool
}

// NewHandlers creates new API handlers
func NewHandlers(q queue.Queue, registry *job.Registry, pool *worker.Pool) *Handlers {
	return &Handlers{
		queue:    q,
		registry: registry,
		pool:     pool,
	}
}

// CreateJobRequest represents a job creation request
type CreateJobRequest struct {
	Type        string                 `json:"type" binding:"required"`
	Payload     map[string]interface{} `json:"payload"`
	MaxAttempts int                    `json:"max_attempts,omitempty"`
	Priority    int                    `json:"priority,omitempty"`
}

// CreateJobResponse represents a job creation response
type CreateJobResponse struct {
	JobID   string      `json:"job_id"`
	Status  job.Status  `json:"status"`
	Message string      `json:"message"`
	Job     *job.Job    `json:"job"`
}

// CreateJob handles POST /jobs
func (h *Handlers) CreateJob(c *gin.Context) {
	var req CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Validate job type
	if _, exists := h.registry.GetHandler(req.Type); !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid job type",
			"details": "No handler registered for job type: " + req.Type,
		})
		return
	}

	// Create job
	j := job.NewJob(req.Type, req.Payload)
	if req.MaxAttempts > 0 {
		j.MaxAttempts = req.MaxAttempts
	}
	if req.Priority > 0 {
		j.Priority = req.Priority
	}

	// Enqueue job
	if err := h.queue.Enqueue(j); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to enqueue job",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, CreateJobResponse{
		JobID:   j.ID,
		Status:  j.Status,
		Message: "Job created successfully",
		Job:     j,
	})
}

// GetJob handles GET /jobs/:id
func (h *Handlers) GetJob(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Job ID is required",
		})
		return
	}

	j, err := h.queue.GetJob(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Job not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"job": j,
	})
}

// ListJobs handles GET /jobs
func (h *Handlers) ListJobs(c *gin.Context) {
	status := c.Query("status")
	limitStr := c.DefaultQuery("limit", "50")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit parameter",
		})
		return
	}

	var jobs []*job.Job
	if status != "" {
		jobs, err = h.queue.GetJobsByStatus(job.Status(status), limit)
	} else {
		// Get all jobs (this is a simplified implementation)
		// In a real application, you might want to paginate through all statuses
		allJobs := make([]*job.Job, 0)
		statuses := []job.Status{
			job.StatusPending,
			job.StatusRunning,
			job.StatusCompleted,
			job.StatusFailed,
			job.StatusRetrying,
		}
		
		for _, s := range statuses {
			statusJobs, err := h.queue.GetJobsByStatus(s, limit)
			if err == nil {
				allJobs = append(allJobs, statusJobs...)
			}
		}
		jobs = allJobs
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get jobs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs":  jobs,
		"count": len(jobs),
	})
}

// DeleteJob handles DELETE /jobs/:id
func (h *Handlers) DeleteJob(c *gin.Context) {
	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Job ID is required",
		})
		return
	}

	if err := h.queue.DeleteJob(jobID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete job",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Job deleted successfully",
	})
}

// GetStats handles GET /stats
func (h *Handlers) GetStats(c *gin.Context) {
	queueStats, err := h.queue.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get queue stats",
			"details": err.Error(),
		})
		return
	}

	poolStats := h.pool.GetStats()

	c.JSON(http.StatusOK, gin.H{
		"queue": queueStats,
		"pool":  poolStats,
	})
}

// GetJobTypes handles GET /job-types
func (h *Handlers) GetJobTypes(c *gin.Context) {
	types := h.registry.GetRegisteredTypes()
	c.JSON(http.StatusOK, gin.H{
		"job_types": types,
		"count":     len(types),
	})
}

// HealthCheck handles GET /health
func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "job-queue",
	})
}