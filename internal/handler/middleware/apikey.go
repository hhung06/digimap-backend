package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hhung06/digimap-backend/internal/dto"
)

const ContextKeyVenueID = "venueID"

// VenueByPublicKeyLookup is the minimal interface needed by APIKeyAuth.
type VenueByPublicKeyLookup interface {
	FindByPublicKey(ctx context.Context, publicKey string) (id uuid.UUID, privateKey string, err error)
}

// APIKeyAuth validates an "ApiKey <public_key>:<private_key>" Authorization header.
// On success it sets the resolved venue ID in the Gin context.
func APIKeyAuth(venues VenueByPublicKeyLookup) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("Authorization")
		publicKey, privateKey, ok := parseAPIKey(raw)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				dto.Fail(dto.CodeAuthRequired, "API key required"))
			return
		}

		venueID, storedPrivate, err := venues.FindByPublicKey(c.Request.Context(), publicKey)
		if err != nil || storedPrivate != privateKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				dto.Fail(dto.CodeAuthRequired, "invalid API key"))
			return
		}

		c.Set(ContextKeyVenueID, venueID)
		c.Next()
	}
}

// GetVenueID extracts the venue UUID resolved by APIKeyAuth from the Gin context.
func GetVenueID(c *gin.Context) uuid.UUID {
	if v, exists := c.Get(ContextKeyVenueID); exists {
		if id, ok := v.(uuid.UUID); ok {
			return id
		}
	}
	return uuid.Nil
}

// parseAPIKey parses "ApiKey <public>:<private>" and returns the two key parts.
func parseAPIKey(header string) (publicKey, privateKey string, ok bool) {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "apikey") {
		return "", "", false
	}
	keys := strings.SplitN(strings.TrimSpace(parts[1]), ":", 2)
	if len(keys) != 2 || keys[0] == "" || keys[1] == "" {
		return "", "", false
	}
	return keys[0], keys[1], true
}
