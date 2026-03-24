package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderRequestID = "X-Request-ID"

// RequestID injects a unique request ID (UUIDv7) into each request context and response header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(HeaderRequestID)
		if id == "" {
			id = newRequestID()
		}
		c.Set("requestID", id)
		c.Header(HeaderRequestID, id)
		c.Next()
	}
}

func newRequestID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String()
	}
	return id.String()
}
