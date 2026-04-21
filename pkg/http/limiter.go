package http

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

// DelayRateLimiter implements RateLimiter with configurable base delay and jitter.
// It ensures that the time between consecutive requests is at least baseDelay + random jitter.
type DelayRateLimiter struct {
	baseDelay time.Duration
	jitter    time.Duration
	mu        sync.Mutex
	lastTime  time.Time
}

// NewDelayRateLimiter returns a DelayRateLimiter with the given base delay and jitter.
func NewDelayRateLimiter(baseDelay, jitter time.Duration) *DelayRateLimiter {
	return &DelayRateLimiter{
		baseDelay: baseDelay,
		jitter:    jitter,
	}
}

// Wait implements RateLimiter by enforcing a minimum delay between requests.
// It accounts for the actual time elapsed since the last request and sleeps for the remaining duration.
func (d *DelayRateLimiter) Wait(ctx context.Context) error {
	d.mu.Lock()
	now := time.Now()
	elapsed := now.Sub(d.lastTime)
	d.mu.Unlock()

	totalDelay := d.baseDelay
	if d.jitter > 0 {
		randomJitter := time.Duration(rand.Int63n(int64(d.jitter)))
		totalDelay += randomJitter
	}

	// Calculate how long we still need to wait
	remaining := totalDelay - elapsed
	if remaining <= 0 {
		// Enough time has passed, just update the timestamp
		d.mu.Lock()
		d.lastTime = time.Now()
		d.mu.Unlock()
		return nil
	}

	// Sleep for the remaining duration
	select {
	case <-time.After(remaining):
		d.mu.Lock()
		d.lastTime = time.Now()
		d.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
