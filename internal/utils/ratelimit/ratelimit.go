package ratelimit

import (
	"sync"
	"time"
)

var rateLimitStore = struct {
	m map[string][]int64
	sync.Mutex
}{
	m: make(map[string][]int64),
}

// RateLimit returns (limited, retryAfterSeconds)
// key: unique identifier (e.g., "otp:phonenumber")
// max: max requests
// windowSeconds: time window in seconds
func RateLimit(key string, maxReqs int, windowSeconds int64) (bool, int64) {
	now := time.Now().Unix()
	windowStart := now - windowSeconds

	rateLimitStore.Lock()
	defer rateLimitStore.Unlock()

	timestamps := rateLimitStore.m[key]
	// Remove old timestamps
	validTimestamps := make([]int64, 0, len(timestamps))
	for _, ts := range timestamps {
		if ts >= windowStart {
			validTimestamps = append(validTimestamps, ts)
		}
	}
	if len(validTimestamps) >= maxReqs {
		// Calculate retry after
		retryAfter := windowSeconds - (now - validTimestamps[0])
		if retryAfter < 1 {
			retryAfter = 1
		}
		return true, retryAfter
	}
	// Add current timestamp
	validTimestamps = append(validTimestamps, now)
	rateLimitStore.m[key] = validTimestamps
	return false, 0
}
