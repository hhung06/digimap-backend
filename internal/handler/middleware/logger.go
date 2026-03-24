package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	applog "github.com/hhung06/digimap-backend/log"
)

// Logger returns a Gin middleware that logs each request with method, path, status, and latency.
func Logger(logger applog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		requestID, _ := c.Get("requestID")

		fields := applog.Fields{
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"status":     status,
			"latency_ms": latency.Milliseconds(),
			"ip":         c.ClientIP(),
			"request_id": requestID,
		}

		entry := logger.WithFields(fields)
		switch {
		case status >= 500:
			entry.Error("server error")
		case status >= 400:
			entry.Warn("client error")
		default:
			entry.Info("request")
		}
	}
}
