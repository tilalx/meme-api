package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ipLimiter holds a rate limiter and the last time the IP was seen.
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu         sync.Mutex
	limiters   = make(map[string]*ipLimiter)
	// Allow 10 requests per second with a burst of 20 per IP.
	rateLimit  = rate.Limit(10)
	burstLimit = 20
)

func init() {
	// Evict entries that haven't been seen in 10 minutes, every 5 minutes.
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			for ip, l := range limiters {
				if time.Since(l.lastSeen) > 10*time.Minute {
					delete(limiters, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	l, ok := limiters[ip]
	if !ok {
		l = &ipLimiter{limiter: rate.NewLimiter(rateLimit, burstLimit)}
		limiters[ip] = l
	}
	l.lastSeen = time.Now()
	return l.limiter
}

// RateLimit returns a gin middleware that limits each client IP to 10 req/s (burst 20).
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "Rate limit exceeded. Please slow down.",
			})
			return
		}

		c.Next()
	}
}
