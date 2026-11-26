package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/surafelbkassa/go-distributed-job-queue/internal/config"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/infrastructure"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/repository"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/usecase"
	"github.com/surafelbkassa/go-distributed-job-queue/internal/worker"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()
	if cfg.JobTimeout >= cfg.Lease {
		log.Error("JOB_TIMEOUT must be shorter than LEASE")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	rdb := infrastructure.NewRedisClient(cfg.RedisAddr)
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Error("redis unreachable", "err", err)
		os.Exit(1)
	}

	q := repository.NewRedisQueue(rdb)
	reg := usecase.NewRegistry()
	worker.RegisterDemo(reg, rdb)

	pcfg := usecase.ProcessorConfig{
		Lease: cfg.Lease, JobTimeout: cfg.JobTimeout,
		BackoffBase: cfg.BackoffBase, BackoffMax: cfg.BackoffMax,
	}
	proc := usecase.NewProcessor(q, reg, pcfg, log)
	maint := usecase.NewMaintenance(q, pcfg, log)

	go maint.RunReaper(ctx, 5*time.Second)
	go maint.RunPromoter(ctx, time.Second)

	pool := worker.NewPool(proc, cfg.Workers, log)
	done := make(chan struct{})
	go func() { pool.Run(ctx); close(done) }()

	<-ctx.Done()
	log.Info("shutting down: draining in-flight jobs")
	select {
	case <-done:
		log.Info("drained cleanly")
	case <-time.After(cfg.JobTimeout + 5*time.Second):
		log.Warn("drain timeout: leases will recover the remaining jobs")
	}
}
