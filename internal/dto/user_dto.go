package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/hhung06/digimap-backend/internal/domain"
)

// ── Profile ───────────────────────────────────────────────────────────────────

// UserToResponse converts a domain.User to the shared UserResponse DTO (defined in auth_dto.go).
func UserToResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID: u.ID, Email: u.Email,
		FirstName: u.FirstName, LastName: u.LastName,
		Phone: u.Phone, AvatarURL: u.AvatarURL,
		IsActive: u.IsActive, IsSystemAdmin: u.IsSystemAdmin,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt, UpdatedAt: u.UpdatedAt,
	}
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatar_url"`
}

// ── Venue user management ─────────────────────────────────────────────────────

type VenueUserResponse struct {
	UserResponse
	Role string `json:"role"`
}

type ChangeRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=owner editor viewer"`
}

// ── Invitations ───────────────────────────────────────────────────────────────

type InvitationResponse struct {
	ID          uuid.UUID  `json:"id"`
	VenueID     uuid.UUID  `json:"venue_id"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	InvitedBy   uuid.UUID  `json:"invited_by"`
	AcceptedAt  *time.Time `json:"accepted_at,omitempty"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
	ExpiresAt   time.Time  `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type InviteUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=owner editor viewer"`
}

type AcceptInvitationRequest struct {
	Token string `json:"token" binding:"required"`
}

func InvitationToResponse(inv *domain.VenueInvitation) InvitationResponse {
	return InvitationResponse{
		ID: inv.ID, VenueID: inv.VenueID, Email: inv.Email,
		Role: inv.Role.String(), Status: string(inv.Status),
		InvitedBy: inv.InvitedBy,
		AcceptedAt: inv.AcceptedAt, CancelledAt: inv.CancelledAt,
		ExpiresAt: inv.ExpiresAt, CreatedAt: inv.CreatedAt,
	}
}
