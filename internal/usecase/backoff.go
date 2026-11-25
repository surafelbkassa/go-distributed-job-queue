package usecase

import (
	"math/rand"
	"time"
)

// Backoff doubles per attempt up to ceiling, then picks a random point in [d/2, d].
func Backoff(attempt int, base, ceiling time.Duration) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	d := base
	for i := 1; i < attempt && d < ceiling; i++ {
		d *= 2
	}
	if d > ceiling {
		d = ceiling
	}
	half := d / 2
	return half + time.Duration(rand.Int63n(int64(half)+1))
}
