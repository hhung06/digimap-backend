package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hhung06/digimap-backend/internal/dto"
	applog "github.com/hhung06/digimap-backend/log"
)

// Recovery returns a Gin middleware that catches panics and returns a 500 response.
func Recovery(logger applog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("panic recovered: %v", err)
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					dto.Fail(dto.CodeInternalError, "internal server error"),
				)
			}
		}()
		c.Next()
	}
}
