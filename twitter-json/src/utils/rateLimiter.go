package utils

import (
	"sync/atomic"
	"time"
)

// RateLimiter rate limiter type
type RateLimiter struct {
	limit  uint64
	count  uint64
	ticker *time.Ticker
	lockCh chan struct{}
}

// NewRateLimiter create rate limiter
func NewRateLimiter(d time.Duration, limit uint64) *RateLimiter {
	ticker := time.NewTicker(d)
	r := &RateLimiter{
		limit:  limit,
		count:  0,
		ticker: ticker,
		lockCh: make(chan struct{}),
	}

	go func() {
		for range ticker.C {
			if r.count > r.limit {
				select {
				case <-r.lockCh:
				default:
					r.resetCount()
				}
			}

			if r.count > 0 {
				r.resetCount()
			}
		}
	}()

	return r
}

// Limit whether limit
func (r *RateLimiter) Limit() bool {
	r.addCount(1)

	if r.getCount() > r.limit {
		var s struct{}
		r.lockCh <- s
	}

	return true
}

func (r *RateLimiter) addCount(interval uint64) {
	atomic.AddUint64(&r.count, interval)
}

func (r *RateLimiter) getCount() uint64 {
	return atomic.LoadUint64(&r.count)
}

func (r *RateLimiter) resetCount() {
	atomic.StoreUint64(&r.count, 1)
}
