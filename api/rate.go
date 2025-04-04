package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	max      int
	period   time.Duration
	requests map[string]int
	mu       sync.RWMutex
}

func newRateLimiter(max int, period time.Duration) *rateLimiter {
	limiter := &rateLimiter{
		max:      max,
		period:   period,
		requests: make(map[string]int),
	}
	limiter.startPruner()
	return limiter
}

func (r *rateLimiter) middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		r.mu.Lock()
		r.requests[c.ClientIP()]++
		count := r.requests[c.ClientIP()]
		r.mu.Unlock()

		if count > r.max {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}

func (r *rateLimiter) startPruner() {
	go func() {
		ticker := time.NewTicker(r.period)
		defer ticker.Stop()

		for range ticker.C {
			r.mu.Lock()
			r.requests = make(map[string]int)
			r.mu.Unlock()
		}
	}()
}
