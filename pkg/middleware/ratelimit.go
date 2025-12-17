package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter manages rate limiters per IP address
type IPRateLimiter struct {
	ips   map[string]*rate.Limiter
	mu    sync.RWMutex
	rate  rate.Limit
	burst int
}

// NewIPRateLimiter creates a new IP-based rate limiter
func NewIPRateLimiter(r rate.Limit, burst int) *IPRateLimiter {
	return &IPRateLimiter{
		ips:   make(map[string]*rate.Limiter),
		rate:  r,
		burst: burst,
	}
}

// GetLimiter returns the rate limiter for an IP
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.rate, i.burst)
		i.ips[ip] = limiter
	}

	return limiter
}

// RateLimiter is a middleware that limits requests per IP
// rate: requests per second, burst: max burst size
func RateLimiter(requestsPerSecond float64, burst int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(rate.Limit(requestsPerSecond), burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.GetLimiter(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Too many requests. Please slow down.",
			})
			return
		}
		c.Next()
	}
}
