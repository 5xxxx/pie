package pie

import (
	"math/rand"
	"time"
)

// randomInt returns a uniformly distributed integer between min and max. When
// min is greater than or equal to max the function returns min to avoid
// panicking.
func randomInt(min, max int64) int64 {
	if min >= max {
		return min
	}
	return min + rand.Int63n(max-min+1)
}

// init seeds the pseudo-random number generator using a nanosecond timestamp so
// cache expiry jitter is unpredictable across process restarts.
func init() {
	rand.Seed(time.Now().UnixNano())
}
