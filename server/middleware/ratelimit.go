package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ipLimiter holds a rate limiter for a single IP address.
type ipLimiter struct {
	limiter *rate.Limiter
}

var (
	limiters sync.Map
	// Allow 10 requests per second with a burst of 20 per IP.
	rateLimit  = rate.Limit(10)
	burstLimit = 20
)

func getLimiter(ip string) *rate.Limiter {
	v, ok := limiters.Load(ip)
	if !ok {
		l := &ipLimiter{limiter: rate.NewLimiter(rateLimit, burstLimit)}
		v, _ = limiters.LoadOrStore(ip, l)
	}
	return v.(*ipLimiter).limiter
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
