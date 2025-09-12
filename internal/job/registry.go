package job

import (
	"crypto/rand"
	"fmt"
	"time"
)

// generateJobID generates a unique job ID
func generateJobID() string {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	return fmt.Sprintf("job_%d_%x", timestamp, randomBytes)
}

// Handler represents a function that processes a job
type Handler func(*Job) error

// Registry holds job handlers
type Registry struct {
	handlers map[string]Handler
}

// NewRegistry creates a new job handler registry
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]Handler),
	}
}

// Register registers a handler for a job type
func (r *Registry) Register(jobType string, handler Handler) {
	r.handlers[jobType] = handler
}

// GetHandler returns the handler for a job type
func (r *Registry) GetHandler(jobType string) (Handler, bool) {
	handler, exists := r.handlers[jobType]
	return handler, exists
}

// GetRegisteredTypes returns all registered job types
func (r *Registry) GetRegisteredTypes() []string {
	types := make([]string, 0, len(r.handlers))
	for jobType := range r.handlers {
		types = append(types, jobType)
	}
	return types
}