package worker

import (
	"context"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/surafelbkassa/go-distributed-job-queue/internal/usecase"
)

type Pool struct {
	proc *usecase.Processor
	n    int
	log  *slog.Logger
}

func NewPool(proc *usecase.Processor, n int, log *slog.Logger) *Pool {
	return &Pool{proc: proc, n: n, log: log}
}

// Run blocks until ctx is cancelled AND every in-flight job has finished.
func (p *Pool) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < p.n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.loop(ctx)
		}()
	}
	wg.Wait()
}

func (p *Pool) loop(ctx context.Context) {
	for ctx.Err() == nil {
		ran, err := p.proc.RunOne(ctx)
		if err != nil {
			if ctx.Err() == nil {
				p.log.Error("claim", "err", err)
			}
			sleepJitter(ctx, 500*time.Millisecond, time.Second)
			continue
		}
		if !ran {
			sleepJitter(ctx, 100*time.Millisecond, 300*time.Millisecond)
		}
	}
}

func sleepJitter(ctx context.Context, lo, hi time.Duration) {
	d := lo + time.Duration(rand.Int63n(int64(hi-lo)+1))
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}
