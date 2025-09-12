package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// JobRequest represents a job creation request
type JobRequest struct {
	Type        string                 `json:"type"`
	Payload     map[string]interface{} `json:"payload"`
	MaxAttempts int                    `json:"max_attempts,omitempty"`
	Priority    int                    `json:"priority,omitempty"`
}

// JobResponse represents the API response for job creation
type JobResponse struct {
	JobID   string                 `json:"job_id"`
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Job     map[string]interface{} `json:"job"`
}

const baseURL = "http://localhost:8080/api/v1"

func main() {
	fmt.Println("Job Queue Client Example")
	fmt.Println("========================")

	// Wait for server to start
	fmt.Println("Waiting for server to be ready...")
	waitForServer()

	// Create different types of jobs
	jobs := []JobRequest{
		{
			Type: "email",
			Payload: map[string]interface{}{
				"to":      "alice@example.com",
				"subject": "Welcome!",
				"body":    "Welcome to our platform!",
			},
			Priority: 2,
		},
		{
			Type: "image_resize",
			Payload: map[string]interface{}{
				"image_url": "https://example.com/photo1.jpg",
				"width":     1200.0,
				"height":    800.0,
			},
		},
		{
			Type: "generate_report",
			Payload: map[string]interface{}{
				"type":   "analytics",
				"period": "weekly",
			},
		},
		{
			Type: "data_processing",
			Payload: map[string]interface{}{
				"source": "user_activity_logs",
				"format": "csv",
			},
		},
	}

	// Submit jobs
	jobIDs := make([]string, 0)
	for i, jobReq := range jobs {
		fmt.Printf("Creating job %d: %s\n", i+1, jobReq.Type)
		jobID, err := createJob(jobReq)
		if err != nil {
			log.Printf("Failed to create job: %v", err)
			continue
		}
		jobIDs = append(jobIDs, jobID)
		fmt.Printf("Created job: %s\n", jobID)
	}

	// Wait for jobs to process
	fmt.Println("\nWaiting for jobs to process...")
	time.Sleep(8 * time.Second)

	// Check job statuses
	fmt.Println("\nChecking job statuses:")
	for _, jobID := range jobIDs {
		status, err := getJobStatus(jobID)
		if err != nil {
			log.Printf("Failed to get job status: %v", err)
			continue
		}
		fmt.Printf("Job %s: %s\n", jobID, status)
	}

	// Get statistics
	fmt.Println("\nGetting queue statistics:")
	stats, err := getStats()
	if err != nil {
		log.Printf("Failed to get stats: %v", err)
	} else {
		fmt.Printf("Queue Stats: %+v\n", stats)
	}
}

func waitForServer() {
	maxAttempts := 30
	for i := 0; i < maxAttempts; i++ {
		resp, err := http.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			fmt.Println("Server is ready!")
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
	log.Fatal("Server not ready after 30 seconds")
}

func createJob(jobReq JobRequest) (string, error) {
	jobData, err := json.Marshal(jobReq)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(baseURL+"/jobs", "application/json", bytes.NewBuffer(jobData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create job: status %d", resp.StatusCode)
	}

	var jobResp JobResponse
	if err := json.NewDecoder(resp.Body).Decode(&jobResp); err != nil {
		return "", err
	}

	return jobResp.JobID, nil
}

func getJobStatus(jobID string) (string, error) {
	resp, err := http.Get(baseURL + "/jobs/" + jobID)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get job status: status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	job, ok := result["job"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	status, ok := job["status"].(string)
	if !ok {
		return "", fmt.Errorf("status not found")
	}

	return status, nil
}

func getStats() (map[string]interface{}, error) {
	resp, err := http.Get(baseURL + "/stats")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get stats: status %d", resp.StatusCode)
	}

	var stats map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}

	return stats, nil
}