package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/repository"
)

const ContextKeyVenueRole = "venueRole"

// VenueAccess checks that the authenticated user has at least minRole for the
// venue specified in the ":venue_id" URL parameter.
// System admins bypass the check and are granted owner-level access.
func VenueAccess(userRepo repository.UserRepository, minRole domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// System admins have full access to every venue
		if GetIsSystemAdmin(c) {
			c.Set(ContextKeyVenueRole, domain.RoleSystemAdmin)
			c.Next()
			return
		}

		// Routes use either ":id" (venue routes) or ":venueID"
		venueIDStr := c.Param("id")
		if venueIDStr == "" {
			venueIDStr = c.Param("venueID")
		}
		venueID, err := uuid.Parse(venueIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				dto.Fail(dto.CodeValidationError, "invalid venue_id"))
			return
		}

		userID := GetUserID(c)
		if userID == uuid.Nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				dto.Fail(dto.CodeAuthRequired, "authentication required"))
			return
		}

		vr, err := userRepo.GetVenueRole(c.Request.Context(), venueID, userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden,
				dto.Fail(dto.CodePermissionDenied, "access to this venue is not allowed"))
			return
		}

		if vr.Role < minRole {
			c.AbortWithStatusJSON(http.StatusForbidden,
				dto.Fail(dto.CodePermissionDenied, "insufficient permissions"))
			return
		}

		c.Set(ContextKeyVenueRole, vr.Role)
		c.Next()
	}
}

// GetVenueRole retrieves the caller's resolved role for the current venue from context.
func GetVenueRole(c *gin.Context) domain.Role {
	if v, exists := c.Get(ContextKeyVenueRole); exists {
		if r, ok := v.(domain.Role); ok {
			return r
		}
	}
	return domain.RoleViewer
}

// SystemAdminRequired aborts with 403 unless the caller is a system admin.
func SystemAdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !GetIsSystemAdmin(c) {
			c.AbortWithStatusJSON(http.StatusForbidden,
				dto.Fail(dto.CodePermissionDenied, "system admin access required"))
			return
		}
		c.Next()
	}
}

// RequireMinRole is a lightweight middleware to add after VenueAccess for
// routes that need a stricter minimum than the group default.
func RequireMinRole(minRole domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetIsSystemAdmin(c) {
			c.Next()
			return
		}
		if GetVenueRole(c) < minRole {
			c.AbortWithStatusJSON(http.StatusForbidden,
				dto.Fail(dto.CodePermissionDenied, "insufficient permissions"))
			return
		}
		c.Next()
	}
}
