package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pandaman404/finance-tracker-go/internal/shared"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	r        rate.Limit
	b        int
}

func newIPLimiter(r rate.Limit, b int) *ipLimiter {
	return &ipLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

func (l *ipLimiter) get(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, exists := l.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(l.r, l.b)
		l.limiters[ip] = limiter
	}

	return limiter
}

// RateLimiter limits requests per IP.
// r = requests per second, b = burst size.
// Example: RateLimiter(5, 10) → 5 req/s, burst up to 10.
func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	limiter := newIPLimiter(r, b)

	return func(c *gin.Context) {
		if !limiter.get(c.ClientIP()).Allow() {
			c.AbortWithStatusJSON(429, gin.H{"error": shared.ErrTooManyRequests.Error()})
			return
		}
		c.Next()
	}
}
