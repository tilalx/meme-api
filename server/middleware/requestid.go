package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const requestIDHeader = "X-Request-ID"

// RequestID reads the incoming X-Request-ID header (or generates one) and
// echoes it back on the response so callers can trace requests end-to-end.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if id == "" {
			id = generateID()
		}
		c.Set("requestID", id)
		c.Header(requestIDHeader, id)
		c.Next()
	}
}

func generateID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
