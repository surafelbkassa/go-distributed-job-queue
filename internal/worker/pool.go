package worker

import (
	"log"
	"sync"

	"github.com/surafelbkassa/go-distributed-job-queue/internal/job"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/queue"
)

// Pool represents a pool of workers
type Pool struct {
	workers     []*Worker
	concurrency int
	queue       queue.Queue
	registry    *job.Registry
	wg          sync.WaitGroup
	isRunning   bool
	mu          sync.RWMutex
}

// NewPool creates a new worker pool
func NewPool(concurrency int, q queue.Queue, registry *job.Registry) *Pool {
	return &Pool{
		concurrency: concurrency,
		queue:       q,
		registry:    registry,
		workers:     make([]*Worker, 0, concurrency),
	}
}

// Start starts all workers in the pool
func (p *Pool) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.isRunning {
		log.Println("Worker pool is already running")
		return
	}
	
	log.Printf("Starting worker pool with %d workers", p.concurrency)
	
	// Create and start workers
	for i := 0; i < p.concurrency; i++ {
		worker := NewWorker(i+1, p.queue, p.registry, &p.wg)
		p.workers = append(p.workers, worker)
		worker.Start()
	}
	
	p.isRunning = true
	log.Printf("Worker pool started with %d workers", len(p.workers))
}

// Stop stops all workers in the pool
func (p *Pool) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if !p.isRunning {
		log.Println("Worker pool is not running")
		return
	}
	
	log.Printf("Stopping worker pool with %d workers", len(p.workers))
	
	// Stop all workers
	for _, worker := range p.workers {
		worker.Stop()
	}
	
	// Wait for all workers to finish
	p.wg.Wait()
	
	// Clear workers slice
	p.workers = p.workers[:0]
	p.isRunning = false
	
	log.Println("Worker pool stopped")
}

// IsRunning returns whether the pool is currently running
func (p *Pool) IsRunning() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isRunning
}

// GetWorkerCount returns the number of workers in the pool
func (p *Pool) GetWorkerCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.workers)
}

// Resize changes the number of workers in the pool
func (p *Pool) Resize(newConcurrency int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if !p.isRunning {
		p.concurrency = newConcurrency
		return
	}
	
	currentCount := len(p.workers)
	
	if newConcurrency > currentCount {
		// Add more workers
		for i := currentCount; i < newConcurrency; i++ {
			worker := NewWorker(i+1, p.queue, p.registry, &p.wg)
			p.workers = append(p.workers, worker)
			worker.Start()
		}
		log.Printf("Worker pool resized from %d to %d workers", currentCount, newConcurrency)
	} else if newConcurrency < currentCount {
		// Remove workers
		workersToStop := p.workers[newConcurrency:]
		p.workers = p.workers[:newConcurrency]
		
		// Stop the excess workers
		for _, worker := range workersToStop {
			worker.Stop()
		}
		
		log.Printf("Worker pool resized from %d to %d workers", currentCount, newConcurrency)
	}
	
	p.concurrency = newConcurrency
}

// GetStats returns pool statistics
func (p *Pool) GetStats() *PoolStats {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	return &PoolStats{
		WorkerCount: len(p.workers),
		IsRunning:   p.isRunning,
		Concurrency: p.concurrency,
	}
}

// PoolStats represents worker pool statistics
type PoolStats struct {
	WorkerCount int  `json:"worker_count"`
	IsRunning   bool `json:"is_running"`
	Concurrency int  `json:"concurrency"`
}