package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/surafelbkassa/go-distributed-job-queue/internal/domain"
)

type ProcessorConfig struct {
	Lease       time.Duration
	JobTimeout  time.Duration
	BackoffBase time.Duration
	BackoffMax  time.Duration
}

type Processor struct {
	q   domain.Queue
	reg *Registry
	cfg ProcessorConfig
	log *slog.Logger
}

func NewProcessor(q domain.Queue, reg *Registry, cfg ProcessorConfig, log *slog.Logger) *Processor {
	return &Processor{q: q, reg: reg, cfg: cfg, log: log}
}

// RunOne claims and processes one job. Returns false if the queue was empty.
func (p *Processor) RunOne(ctx context.Context) (bool, error) {
	job, err := p.q.Claim(ctx, p.cfg.Lease)
	if err != nil {
		return false, err
	}
	if job == nil {
		return false, nil
	}
	p.execute(job)
	return true, nil
}

func (p *Processor) execute(job *domain.Job) {
	// Fresh context on purpose: a shutdown signal must not cancel a job already running.
	ctx, cancel := context.WithTimeout(context.Background(), p.cfg.JobTimeout)
	defer cancel()

	var err error
	if h, ok := p.reg.Lookup(job.Type); !ok {
		err = fmt.Errorf("no handler registered for type %q", job.Type)
	} else {
		err = safeCall(ctx, h, job)
	}

	bg := context.Background()
	if err != nil {
		backoff := Backoff(job.Attempts, p.cfg.BackoffBase, p.cfg.BackoffMax)
		out, ferr := p.q.Fail(bg, job.ID, err.Error(), backoff)
		p.log.Warn("job failed", "job_id", job.ID, "attempt", job.Attempts,
			"outcome", out, "err", err, "fail_err", ferr)
		return
	}

	owned, aerr := p.q.Ack(bg, job.ID)
	switch {
	case aerr != nil:
		p.log.Error("ack failed", "job_id", job.ID, "err", aerr)
	case !owned:
		p.log.Warn("lease expired before ack; job may run twice", "job_id", job.ID)
	default:
		p.log.Info("job done", "job_id", job.ID, "attempt", job.Attempts)
	}
}

func safeCall(ctx context.Context, h HandlerFunc, j *domain.Job) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return h(ctx, j)
}
