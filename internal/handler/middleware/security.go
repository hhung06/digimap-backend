package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hhung06/digimap-backend/config"
)

// SecurityHeaders returns a Gin middleware that sets security-related HTTP response headers.
// HSTS is only applied in production to avoid breaking local HTTPS dev flows.
func SecurityHeaders(cfg config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent the response from being framed (clickjacking protection)
		c.Header("X-Frame-Options", "DENY")

		// Stop browsers from MIME-sniffing the content type
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable the browser's built-in XSS filter (legacy browsers)
		c.Header("X-XSS-Protection", "1; mode=block")

		// Control what information is sent in the Referer header
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Restrict which browser features this page can use
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Swagger UI requires inline styles, data URIs, and blob workers — relax CSP for its path only
		if strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
			c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; worker-src blob:")
		} else {
			c.Header("Content-Security-Policy", "default-src 'self'; object-src 'none'; base-uri 'self'")
		}

		// HSTS: only send over HTTPS, only in production
		if cfg.IsProduction() {
			c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}

		c.Next()
	}
}
