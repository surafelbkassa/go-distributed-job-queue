package repository

import (
	_ "embed"

	"github.com/redis/go-redis/v9"
)

//go:embed scripts/claim.lua
var claimSrc string

//go:embed scripts/ack.lua
var ackSrc string

//go:embed scripts/fail.lua
var failSrc string

//go:embed scripts/reap.lua
var reapSrc string

//go:embed scripts/promote.lua
var promoteSrc string

var (
	claimScript   = redis.NewScript(claimSrc)
	ackScript     = redis.NewScript(ackSrc)
	failScript    = redis.NewScript(failSrc)
	reapScript    = redis.NewScript(reapSrc)
	promoteScript = redis.NewScript(promoteSrc)
)

const (
	keyProcessing = "queue:processing"
	keyDelayed    = "queue:delayed"
	keyDead       = "queue:dead"
)
