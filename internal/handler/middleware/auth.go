package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/service"
)

// Context keys used to pass auth data through the Gin request chain.
type contextKey string

const (
	ContextKeyUserID        contextKey = "userID"
	ContextKeyEmail         contextKey = "email"
	ContextKeyIsSystemAdmin contextKey = "isSystemAdmin"
)

// AuthRequired validates the Bearer JWT and injects user identity into the Gin context.
// Returns 401 if the token is missing or invalid.
func AuthRequired(authSvc service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := extractBearer(c.GetHeader("Authorization"))
		if raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				dto.Fail(dto.CodeAuthRequired, "authentication required"))
			return
		}

		claims, err := authSvc.ValidateClaims(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				dto.Fail(dto.CodeAuthRequired, "invalid or expired token"))
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				dto.Fail(dto.CodeAuthRequired, "invalid token subject"))
			return
		}

		c.Set(string(ContextKeyUserID), userID)
		c.Set(string(ContextKeyEmail), claims.Email)
		c.Set(string(ContextKeyIsSystemAdmin), claims.IsSystemAdmin)
		c.Next()
	}
}

// GetUserID extracts the authenticated user's UUID from the Gin context.
func GetUserID(c *gin.Context) uuid.UUID {
	if v, exists := c.Get(string(ContextKeyUserID)); exists {
		if id, ok := v.(uuid.UUID); ok {
			return id
		}
	}
	return uuid.Nil
}

// GetIsSystemAdmin returns true if the current user is a system administrator.
func GetIsSystemAdmin(c *gin.Context) bool {
	if v, exists := c.Get(string(ContextKeyIsSystemAdmin)); exists {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// extractBearer parses "Bearer <token>" from an Authorization header value.
func extractBearer(header string) string {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
		return strings.TrimSpace(parts[1])
	}
	return ""
}
