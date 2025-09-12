package queue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/surafelbkassa/go-distributed-job-queue/config"
	"github.com/surafelbkassa/go-distributed-job-queue/models"
)

// JobProcessor defines the interface for processing jobs
type JobProcessor interface {
	Process(job *models.Job) error
}

// DefaultProcessor provides a default implementation for job processing
type DefaultProcessor struct{}

// Process processes a job based on its type
func (p *DefaultProcessor) Process(job *models.Job) error {
	log.Printf("Processing job %s of type %s", job.ID, job.Type)
	
	// Simulate job processing based on type
	switch job.Type {
	case "email":
		return p.processEmailJob(job)
	case "data_processing":
		return p.processDataJob(job)
	case "report_generation":
		return p.processReportJob(job)
	default:
		return p.processGenericJob(job)
	}
}

func (p *DefaultProcessor) processEmailJob(job *models.Job) error {
	// Simulate email sending
	time.Sleep(time.Duration(100+job.RetryCount*50) * time.Millisecond)
	
	recipient, ok := job.Payload["recipient"].(string)
	if !ok {
		return fmt.Errorf("recipient not specified")
	}
	
	log.Printf("Sending email to %s", recipient)
	return nil
}

func (p *DefaultProcessor) processDataJob(job *models.Job) error {
	// Simulate data processing
	time.Sleep(time.Duration(200+job.RetryCount*100) * time.Millisecond)
	
	dataSize, ok := job.Payload["data_size"].(float64)
	if !ok {
		return fmt.Errorf("data_size not specified")
	}
	
	log.Printf("Processing %v MB of data", dataSize)
	return nil
}

func (p *DefaultProcessor) processReportJob(job *models.Job) error {
	// Simulate report generation
	time.Sleep(time.Duration(500+job.RetryCount*200) * time.Millisecond)
	
	reportType, ok := job.Payload["report_type"].(string)
	if !ok {
		return fmt.Errorf("report_type not specified")
	}
	
	log.Printf("Generating %s report", reportType)
	return nil
}

func (p *DefaultProcessor) processGenericJob(job *models.Job) error {
	// Simulate generic job processing
	time.Sleep(time.Duration(150+job.RetryCount*75) * time.Millisecond)
	log.Printf("Processing generic job with payload: %+v", job.Payload)
	return nil
}

// WorkerPool manages a pool of workers for processing jobs
type WorkerPool struct {
	queue       *RedisQueue
	processor   JobProcessor
	config      *config.Config
	workers     []*Worker
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	stopOnce    sync.Once
}

// Worker represents a single worker goroutine
type Worker struct {
	id        int
	queue     *RedisQueue
	processor JobProcessor
	ctx       context.Context
	wg        *sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(queue *RedisQueue, processor JobProcessor, config *config.Config) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	
	if processor == nil {
		processor = &DefaultProcessor{}
	}
	
	return &WorkerPool{
		queue:     queue,
		processor: processor,
		config:    config,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	log.Printf("Starting worker pool with %d workers", wp.config.WorkerCount)
	
	wp.workers = make([]*Worker, wp.config.WorkerCount)
	
	for i := 0; i < wp.config.WorkerCount; i++ {
		worker := &Worker{
			id:        i + 1,
			queue:     wp.queue,
			processor: wp.processor,
			ctx:       wp.ctx,
			wg:        &wp.wg,
		}
		
		wp.workers[i] = worker
		wp.wg.Add(1)
		go worker.run()
	}
	
	log.Printf("Worker pool started successfully")
}

// Stop stops the worker pool gracefully
func (wp *WorkerPool) Stop() {
	wp.stopOnce.Do(func() {
		log.Println("Stopping worker pool...")
		wp.cancel()
		wp.wg.Wait()
		log.Println("Worker pool stopped")
	})
}

// run starts the worker's main processing loop
func (w *Worker) run() {
	defer w.wg.Done()
	
	log.Printf("Worker %d started", w.id)
	
	for {
		select {
		case <-w.ctx.Done():
			log.Printf("Worker %d stopping", w.id)
			return
		default:
			// Try to dequeue a job with timeout
			job, err := w.queue.Dequeue(5 * time.Second)
			if err != nil {
				log.Printf("Worker %d: error dequeuing job: %v", w.id, err)
				continue
			}
			
			if job == nil {
				// No job available, continue
				continue
			}
			
			w.processJob(job)
		}
	}
}

// processJob processes a single job
func (w *Worker) processJob(job *models.Job) {
	log.Printf("Worker %d: processing job %s (attempt %d/%d)", w.id, job.ID, job.RetryCount+1, job.MaxRetries+1)
	
	// Mark job as running
	job.MarkRunning()
	if err := w.queue.UpdateJobStatus(job); err != nil {
		log.Printf("Worker %d: failed to update job status to running: %v", w.id, err)
	}
	
	// Process the job
	err := w.processor.Process(job)
	
	if err != nil {
		log.Printf("Worker %d: job %s failed: %v", w.id, job.ID, err)
		
		// Mark job as failed
		job.MarkFailed(err)
		
		// Check if we can retry
		if job.CanRetry() {
			log.Printf("Worker %d: retrying job %s (retry %d/%d)", w.id, job.ID, job.RetryCount, job.MaxRetries)
			
			// Re-enqueue for retry with exponential backoff
			go func() {
				backoffDuration := time.Duration(job.RetryCount*job.RetryCount) * time.Second
				time.Sleep(backoffDuration)
				
				job.Status = models.StatusPending
				if err := w.queue.Enqueue(job); err != nil {
					log.Printf("Worker %d: failed to re-enqueue job %s: %v", w.id, job.ID, err)
				}
			}()
		} else {
			log.Printf("Worker %d: job %s failed permanently after %d retries", w.id, job.ID, job.RetryCount)
		}
		
		// Update job status
		if err := w.queue.UpdateJobStatus(job); err != nil {
			log.Printf("Worker %d: failed to update job status to failed: %v", w.id, err)
		}
	} else {
		log.Printf("Worker %d: job %s completed successfully", w.id, job.ID)
		
		// Mark job as completed
		job.MarkCompleted(map[string]interface{}{
			"processed_by": fmt.Sprintf("worker-%d", w.id),
			"processed_at": time.Now(),
		})
		
		// Update job status
		if err := w.queue.UpdateJobStatus(job); err != nil {
			log.Printf("Worker %d: failed to update job status to completed: %v", w.id, err)
		}
	}
}