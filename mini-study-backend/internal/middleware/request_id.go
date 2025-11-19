package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDKey = "X-Request-ID"

// RequestID ensures every request has a unique identifier for tracing.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(requestIDKey)
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Set(requestIDKey, rid)
		c.Writer.Header().Set(requestIDKey, rid)
		c.Next()
	}
}

// GetRequestID extracts the request id from context.
func GetRequestID(c *gin.Context) string {
	if value, ok := c.Get(requestIDKey); ok {
		if rid, ok := value.(string); ok {
			return rid
		}
	}
	return c.Writer.Header().Get(requestIDKey)
}
