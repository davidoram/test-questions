package rate_limit

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConsume(t *testing.T) {
	bucket := NewTokenBucket(t.Context(), 1, time.Second)
	defer bucket.Close()

	require.NoError(t, bucket.Consume(context.Background()))
	require.Equal(t, 0, int(bucket.tokens))
	require.ErrorContains(t, bucket.Consume(context.Background()), "token unavailable")
}

func TestConsumeStopsWaitingForLockWhenContextIsCancelled(t *testing.T) {
	bucket := NewTokenBucket(t.Context(), 1, time.Second)
	defer bucket.Close()

	// Take a lock
	bucket.tokenLock.Lock()
	defer bucket.tokenLock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Fail can't aquirqe lock within deadline
	require.ErrorIs(t, bucket.Consume(ctx), context.DeadlineExceeded)
	require.Equal(t, 1, int(bucket.tokens))
}

func TestConsumeWithCancelledContext(t *testing.T) {
	bucket := NewTokenBucket(t.Context(), 1, time.Second)
	defer bucket.Close()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	require.ErrorIs(t, bucket.Consume(ctx), context.Canceled)
	require.Equal(t, 1, int(bucket.tokens))
}

func TestBucketRefilled(t *testing.T) {
	// Runs with a fake time
	synctest.Test(t, func(t *testing.T) {
		tokensPerPeriod := 10
		period := 10 * time.Millisecond
		bucket := NewTokenBucket(t.Context(), uint(tokensPerPeriod), period)
		defer bucket.Close()

		// Wait for more tokens to fill the bucket
		time.Sleep(period + time.Millisecond)

		// Consume all tokens
		for range tokensPerPeriod {
			require.NoError(t, bucket.Consume(t.Context()))
		}

		// Next hits rate limit
		require.ErrorContains(t, bucket.Consume(t.Context()), "token unavailable")

		// Wait for more tokens to fill the bucket
		time.Sleep(period + time.Millisecond)

		// Consume all tokens
		for range tokensPerPeriod {
			require.NoError(t, bucket.Consume(t.Context()))
		}
	})
}
