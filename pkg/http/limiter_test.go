package http

import (
	"context"
	"testing"
	"time"
)

func TestDelayRateLimiterEnforcesBaseDelay(t *testing.T) {
	baseDelay := 100 * time.Millisecond
	limiter := NewDelayRateLimiter(baseDelay, 0)

	ctx := context.Background()

	// First call initiates the limiter
	err := limiter.Wait(ctx)
	if err != nil {
		t.Fatalf("first wait failed: %v", err)
	}

	// Second call should enforce the delay
	start := time.Now()
	err = limiter.Wait(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("second wait failed: %v", err)
	}
	if elapsed < baseDelay {
		t.Errorf("expected to wait at least %v, but waited %v", baseDelay, elapsed)
	}
}

func TestDelayRateLimiterEnforcesDelayBetweenCalls(t *testing.T) {
	baseDelay := 100 * time.Millisecond
	limiter := NewDelayRateLimiter(baseDelay, 0)

	ctx := context.Background()

	// First call
	err := limiter.Wait(ctx)
	if err != nil {
		t.Fatalf("first wait failed: %v", err)
	}

	// Second call should enforce delay from first call
	start := time.Now()
	err = limiter.Wait(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("second wait failed: %v", err)
	}
	if elapsed < baseDelay {
		t.Errorf("expected to wait at least %v between calls, but waited %v", baseDelay, elapsed)
	}
}

func TestDelayRateLimiterAccountsForElapsedTime(t *testing.T) {
	baseDelay := 200 * time.Millisecond
	limiter := NewDelayRateLimiter(baseDelay, 0)

	ctx := context.Background()

	// First call
	err := limiter.Wait(ctx)
	if err != nil {
		t.Fatalf("first wait failed: %v", err)
	}

	// Wait half the base delay manually
	time.Sleep(baseDelay / 2)

	// Second call should only wait for the remaining duration
	start := time.Now()
	err = limiter.Wait(ctx)
	remainingElapsed := time.Since(start)

	if err != nil {
		t.Fatalf("second wait failed: %v", err)
	}

	// Should wait approximately the remaining half (with some tolerance for timing)
	expectedMax := baseDelay/2 + 50*time.Millisecond
	if remainingElapsed > expectedMax {
		t.Errorf("expected to wait at most ~%v for remaining duration, but waited %v", baseDelay/2, remainingElapsed)
	}
}

func TestDelayRateLimiterWithJitter(t *testing.T) {
	baseDelay := 50 * time.Millisecond
	jitter := 50 * time.Millisecond

	ctx := context.Background()

	// Run multiple times to observe jitter distribution
	minWait := time.Duration(1<<63 - 1)
	maxWait := time.Duration(0)

	for i := 0; i < 10; i++ {
		limiter := NewDelayRateLimiter(baseDelay, jitter)

		// First call initiates the limiter
		err := limiter.Wait(ctx)
		if err != nil {
			t.Fatalf("first wait failed: %v", err)
		}

		// Second call experiences the delay
		start := time.Now()
		err = limiter.Wait(ctx)
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("second wait failed: %v", err)
		}

		if elapsed < minWait {
			minWait = elapsed
		}
		if elapsed > maxWait {
			maxWait = elapsed
		}
	}

	// The minimum wait should be at least baseDelay
	if minWait < baseDelay {
		t.Errorf("minimum wait %v is less than base delay %v", minWait, baseDelay)
	}

	// The maximum wait should be at most baseDelay + jitter (with some tolerance)
	maxPossible := baseDelay + jitter + 20*time.Millisecond
	if maxWait > maxPossible {
		t.Errorf("maximum wait %v exceeds expected max %v", maxWait, maxPossible)
	}

	// Ideally we see variation due to jitter (but this is probabilistic, so we allow it to be flaky)
	variation := maxWait - minWait
	if variation > 10*time.Millisecond {
		t.Logf("observed jitter variation: %v", variation)
	}
}

func TestDelayRateLimiterRespectsCancellation(t *testing.T) {
	baseDelay := 5 * time.Second
	limiter := NewDelayRateLimiter(baseDelay, 0)

	ctx := context.Background()

	// First call to initialize the limiter
	err := limiter.Wait(ctx)
	if err != nil {
		t.Fatalf("first wait failed: %v", err)
	}

	// Second call with a timeout - should be cancelled before the wait completes
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	err = limiter.Wait(ctx)
	elapsed := time.Since(start)

	if err != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
	if elapsed > 200*time.Millisecond {
		t.Errorf("context cancellation should have been respected quickly, but waited %v", elapsed)
	}
}

func TestDelayRateLimiterZeroDelay(t *testing.T) {
	limiter := NewDelayRateLimiter(0, 0)

	ctx := context.Background()
	start := time.Now()
	err := limiter.Wait(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed > 50*time.Millisecond {
		t.Errorf("zero delay should not cause significant wait, but waited %v", elapsed)
	}
}

func TestDelayRateLimiterConcurrentSafety(t *testing.T) {
	baseDelay := 50 * time.Millisecond
	limiter := NewDelayRateLimiter(baseDelay, 0)

	ctx := context.Background()
	done := make(chan error, 5)

	// Launch multiple goroutines calling Wait concurrently
	for i := 0; i < 5; i++ {
		go func() {
			err := limiter.Wait(ctx)
			done <- err
		}()
	}

	// Collect results
	for i := 0; i < 5; i++ {
		err := <-done
		if err != nil {
			t.Errorf("concurrent wait failed: %v", err)
		}
	}
}
