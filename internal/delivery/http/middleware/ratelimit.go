package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// RateLimitMiddleware creates a gin middleware for rate limiting
// limit: request per second (e.g. 1/60 for 1 request per minute)
// burst: max burst size
func RateLimitMiddleware(limit rate.Limit, burst int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(limit, burst)

	// Perform maintenance to remove old entries to prevent memory leak (simplified version)
	// In production, we should use a proper cache with TTL or a separate cleanup goroutine
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			limiter.mu.Lock()
			// Reset map every 10 mins for simplicity in this demo
			// Ideally we check LastSeen for each IP
			limiter.ips = make(map[string]*rate.Limiter)
			limiter.mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.GetLimiter(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Terlalu banyak percobaan login. Silakan coba lagi nanti.",
			})
			return
		}
		c.Next()
	}
}
