package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hhung06/digimap-backend/config"
)

// CORS returns a Gin middleware configured from AppConfig.
func CORS(cfg config.AppConfig) gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()

	origins := cfg.CORSAllowedOrigins
	if len(origins) == 1 && origins[0] == "*" {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = origins
	}

	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{
		"Origin", "Content-Type", "Accept", "Authorization",
		"Cache-Control", "X-Requested-With", "X-Request-ID", "X-Public-Key",
		"X-Lang-Code",
	}
	corsConfig.ExposeHeaders = []string{"X-Request-ID"}
	corsConfig.AllowCredentials = !corsConfig.AllowAllOrigins

	return cors.New(corsConfig)
}
