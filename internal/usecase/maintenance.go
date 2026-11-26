package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/surafelbkassa/go-distributed-job-queue/internal/domain"
)

type Maintenance struct {
	q   domain.Queue
	cfg ProcessorConfig
	log *slog.Logger
}

func NewMaintenance(q domain.Queue, cfg ProcessorConfig, log *slog.Logger) *Maintenance {
	return &Maintenance{q: q, cfg: cfg, log: log}
}

func (m *Maintenance) RunReaper(ctx context.Context, every time.Duration) {
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			m.reapOnce(ctx)
		}
	}
}

func (m *Maintenance) reapOnce(ctx context.Context) {
	ids, err := m.q.ReapExpired(ctx, 100)
	if err != nil {
		m.log.Error("reap", "err", err)
		return
	}
	for _, id := range ids {
		job, err := m.q.Get(ctx, id)
		if err != nil {
			continue
		}
		backoff := Backoff(job.Attempts, m.cfg.BackoffBase, m.cfg.BackoffMax)
		out, err := m.q.Fail(ctx, id, "lease expired", backoff)
		m.log.Warn("reaped expired lease", "job_id", id, "outcome", out, "err", err)
	}
}

func (m *Maintenance) RunPromoter(ctx context.Context, every time.Duration) {
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if n, err := m.q.PromoteDue(ctx, 100); err != nil {
				m.log.Error("promote", "err", err)
			} else if n > 0 {
				m.log.Info("promoted retries", "count", n)
			}
		}
	}
}
