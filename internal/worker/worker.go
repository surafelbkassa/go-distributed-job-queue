package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/surafelbkassa/go-distributed-job-queue/internal/job"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/queue"
)

// Worker represents a single worker that processes jobs
type Worker struct {
	id       int
	queue    queue.Queue
	registry *job.Registry
	ctx      context.Context
	cancel   context.CancelFunc
	wg       *sync.WaitGroup
}

// NewWorker creates a new worker
func NewWorker(id int, q queue.Queue, registry *job.Registry, wg *sync.WaitGroup) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Worker{
		id:       id,
		queue:    q,
		registry: registry,
		ctx:      ctx,
		cancel:   cancel,
		wg:       wg,
	}
}

// Start begins processing jobs
func (w *Worker) Start() {
	w.wg.Add(1)
	go w.run()
}

// Stop stops the worker
func (w *Worker) Stop() {
	w.cancel()
}

// run is the main worker loop
func (w *Worker) run() {
	defer w.wg.Done()
	
	log.Printf("Worker %d started", w.id)
	
	for {
		select {
		case <-w.ctx.Done():
			log.Printf("Worker %d stopped", w.id)
			return
		default:
			// Try to get a job from the queue
			j, err := w.queue.Dequeue()
			if err != nil {
				log.Printf("Worker %d: error dequeuing job: %v", w.id, err)
				time.Sleep(time.Second)
				continue
			}
			
			if j == nil {
				// No job available, wait a bit
				time.Sleep(100 * time.Millisecond)
				continue
			}
			
			// Check if job should be delayed
			if j.ShouldDelay() {
				// Re-enqueue the job
				if err := w.queue.Enqueue(j); err != nil {
					log.Printf("Worker %d: error re-enqueuing delayed job %s: %v", w.id, j.ID, err)
				}
				continue
			}
			
			// Process the job
			w.processJob(j)
		}
	}
}

// processJob processes a single job
func (w *Worker) processJob(j *job.Job) {
	log.Printf("Worker %d: processing job %s (type: %s, attempt: %d)", w.id, j.ID, j.Type, j.Attempts+1)
	
	// Mark job as running
	j.MarkRunning()
	if err := w.queue.UpdateJob(j); err != nil {
		log.Printf("Worker %d: error updating job status to running: %v", w.id, err)
	}
	
	// Get the handler for this job type
	handler, exists := w.registry.GetHandler(j.Type)
	if !exists {
		log.Printf("Worker %d: no handler found for job type %s", w.id, j.Type)
		j.MarkFailed(fmt.Errorf("no handler found for job type: %s", j.Type))
		w.queue.UpdateJob(j)
		return
	}
	
	// Execute the job
	err := handler(j)
	
	if err != nil {
		log.Printf("Worker %d: job %s failed: %v", w.id, j.ID, err)
		
		if j.CanRetry() {
			// Mark for retry
			j.MarkRetrying(err)
			j.Delay = time.Duration(j.Attempts) * time.Second // Exponential backoff
			
			// Re-enqueue the job
			if enqueueErr := w.queue.Enqueue(j); enqueueErr != nil {
				log.Printf("Worker %d: error re-enqueuing job for retry: %v", w.id, enqueueErr)
				j.MarkFailed(err)
				w.queue.UpdateJob(j)
			} else {
				log.Printf("Worker %d: job %s scheduled for retry (attempt %d/%d)", w.id, j.ID, j.Attempts, j.MaxAttempts)
				w.queue.UpdateJob(j)
			}
		} else {
			// Mark as failed
			j.MarkFailed(err)
			w.queue.UpdateJob(j)
		}
	} else {
		log.Printf("Worker %d: job %s completed successfully", w.id, j.ID)
		j.MarkCompleted()
		w.queue.UpdateJob(j)
	}
}