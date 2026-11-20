jobqueue/
  cmd/
    api/main.go          # wires config, redis, usecases, gin router
    worker/main.go       # wires config, redis, pool, reaper, promoter
  internal/
    domain/
      job.go             # Job, Status, errors
      queue.go           # Queue interface (the port)
    usecase/
      enqueue.go         # validate, build Job, call Queue.Enqueue
      get_job.go
      process.go         # execute a claimed job: run handler, ack or fail
      maintenance.go     # reaper and promoter logic
    repository/
      redis/
        queue.go         # implements domain.Queue using go-redis
        scripts.go       # go:embed of the lua files
        scripts/         # enqueue.lua claim.lua ack.lua fail.lua reap.lua promote.lua
    api/
      router.go handlers.go dto.go
    worker/
      pool.go            # N goroutines claiming and executing
      registry.go        # map of job type to handler func
      backoff.go
    config/config.go     # env vars: REDIS_ADDR, WORKERS, LEASE, MAX_ATTEMPTS...
  scripts/chaos.sh       # the kill-a-worker proof
  docs/                  # architecture.md, delivery-semantics.md
  Dockerfile docker-compose.yml README.md