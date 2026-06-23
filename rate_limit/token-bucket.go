package rate_limit

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrTokenUnavailable = errors.New("token unavailable")

type TokenBucket struct {
	tokens    uint64     // Total Tokens in the bucket
	tokenLock sync.Mutex // Lock used when adding/removing tokens from bucket
	closeChan chan bool  // closing channel
	closed    bool
}

// NewTokenBucket creates a TokenBucket initally with with count tokens.
// After each period the TokenBucket is reset to having that same number of tokens.
// Caller must cancel the context to clean up resources
func NewTokenBucket(ctx context.Context, count uint, period time.Duration) *TokenBucket {
	bucket := TokenBucket{tokens: uint64(count), closeChan: make(chan bool)}
	ticker := time.NewTicker(period)
	go func(ctx context.Context) {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				bucket.fill(uint64(count))
			case <-bucket.closeChan:
				return
			}
		}
	}(ctx)
	return &bucket
}

func (bucket *TokenBucket) Close() {
	bucket.closeChan <- true
	close(bucket.closeChan)
	bucket.closed = true
}

func (bucket *TokenBucket) Consume(ctx context.Context) error {
	if bucket.closed {
		return fmt.Errorf("bucket closed")
	}
	const lockRetryInterval = time.Millisecond

	retry := time.NewTicker(lockRetryInterval)
	defer retry.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("acquiring token lock: %w", ctx.Err())
		default:
		}

		if bucket.tokenLock.TryLock() {
			break
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("acquiring token lock: %w", ctx.Err())
		case <-retry.C:
		}
	}
	defer bucket.tokenLock.Unlock()

	if bucket.tokens == 0 {
		return ErrTokenUnavailable
	}

	bucket.tokens--
	return nil
}

func (bucket *TokenBucket) fill(newTokens uint64) {
	bucket.tokenLock.Lock()
	defer bucket.tokenLock.Unlock()

	bucket.tokens = newTokens
}
