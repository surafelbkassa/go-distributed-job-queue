package main

import (
	"errors"
	"testing"

	"github.com/surafelbkassa/go-distributed-job-queue/config"
	"github.com/surafelbkassa/go-distributed-job-queue/models"
	"github.com/surafelbkassa/go-distributed-job-queue/queue"
)

func TestJobCreation(t *testing.T) {
	payload := map[string]interface{}{
		"recipient": "test@example.com",
		"subject":   "Test Email",
	}
	
	job := models.NewJob("email", payload, 3)
	
	if job.Type != "email" {
		t.Errorf("Expected job type 'email', got '%s'", job.Type)
	}
	
	if job.Status != models.StatusPending {
		t.Errorf("Expected job status 'pending', got '%s'", job.Status)
	}
	
	if job.MaxRetries != 3 {
		t.Errorf("Expected max retries 3, got %d", job.MaxRetries)
	}
	
	if job.RetryCount != 0 {
		t.Errorf("Expected retry count 0, got %d", job.RetryCount)
	}
}

func TestJobStatusTransitions(t *testing.T) {
	job := models.NewJob("test", map[string]interface{}{}, 3)
	
	// Test marking as running
	job.MarkRunning()
	if job.Status != models.StatusRunning {
		t.Errorf("Expected job status 'running', got '%s'", job.Status)
	}
	if job.StartedAt == nil {
		t.Error("Expected StartedAt to be set when marking as running")
	}
	
	// Test marking as completed
	result := map[string]interface{}{"result": "success"}
	job.MarkCompleted(result)
	if job.Status != models.StatusCompleted {
		t.Errorf("Expected job status 'completed', got '%s'", job.Status)
	}
	if job.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set when marking as completed")
	}
}

func TestJobRetryLogic(t *testing.T) {
	job := models.NewJob("test", map[string]interface{}{}, 2)
	
	if !job.CanRetry() {
		t.Error("Expected job to be retryable initially")
	}
	
	// Simulate first failure
	job.MarkFailed(errors.New("test error"))
	if job.RetryCount != 1 {
		t.Errorf("Expected retry count 1, got %d", job.RetryCount)
	}
	if !job.CanRetry() {
		t.Error("Expected job to still be retryable after first failure")
	}
	
	// Simulate second failure
	job.MarkFailed(errors.New("test error 2"))
	if job.RetryCount != 2 {
		t.Errorf("Expected retry count 2, got %d", job.RetryCount)
	}
	if job.CanRetry() {
		t.Error("Expected job to not be retryable after reaching max retries")
	}
}

func TestConfigLoading(t *testing.T) {
	cfg := config.LoadConfig()
	
	if cfg.WorkerCount <= 0 {
		t.Error("Expected positive worker count")
	}
	
	if cfg.MaxRetries < 0 {
		t.Error("Expected non-negative max retries")
	}
	
	if cfg.QueueName == "" {
		t.Error("Expected non-empty queue name")
	}
}

func TestDefaultProcessor(t *testing.T) {
	processor := &queue.DefaultProcessor{}
	
	// Test email job
	emailJob := models.NewJob("email", map[string]interface{}{
		"recipient": "test@example.com",
	}, 3)
	
	err := processor.Process(emailJob)
	if err != nil {
		t.Errorf("Expected email job to process successfully, got error: %v", err)
	}
	
	// Test data processing job
	dataJob := models.NewJob("data_processing", map[string]interface{}{
		"data_size": 100.5,
	}, 3)
	
	err = processor.Process(dataJob)
	if err != nil {
		t.Errorf("Expected data processing job to process successfully, got error: %v", err)
	}
	
	// Test generic job
	genericJob := models.NewJob("custom_type", map[string]interface{}{
		"custom_field": "value",
	}, 3)
	
	err = processor.Process(genericJob)
	if err != nil {
		t.Errorf("Expected generic job to process successfully, got error: %v", err)
	}
}

func TestJobJSONSerialization(t *testing.T) {
	originalJob := models.NewJob("test", map[string]interface{}{
		"key": "value",
	}, 3)
	
	// Serialize to JSON
	jsonData, err := originalJob.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize job to JSON: %v", err)
	}
	
	// Deserialize from JSON
	restoredJob, err := models.FromJSON(jsonData)
	if err != nil {
		t.Fatalf("Failed to deserialize job from JSON: %v", err)
	}
	
	// Compare key fields
	if originalJob.ID != restoredJob.ID {
		t.Errorf("Job ID mismatch: expected %s, got %s", originalJob.ID, restoredJob.ID)
	}
	
	if originalJob.Type != restoredJob.Type {
		t.Errorf("Job type mismatch: expected %s, got %s", originalJob.Type, restoredJob.Type)
	}
	
	if originalJob.Status != restoredJob.Status {
		t.Errorf("Job status mismatch: expected %s, got %s", originalJob.Status, restoredJob.Status)
	}
	
	if originalJob.MaxRetries != restoredJob.MaxRetries {
		t.Errorf("Max retries mismatch: expected %d, got %d", originalJob.MaxRetries, restoredJob.MaxRetries)
	}
}