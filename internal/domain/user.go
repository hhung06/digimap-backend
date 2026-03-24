package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a user's permission level within a venue.
// Higher numeric values have more permissions (owner > editor > viewer).
type Role int

const (
	RoleViewer      Role = 1
	RoleEditor      Role = 2
	RoleOwner       Role = 3
	RoleSystemAdmin Role = 4 // not stored in venue_user_roles; derived from users.is_system_admin
)

// RoleFromString converts the DB enum string to a Role.
func RoleFromString(s string) Role {
	switch s {
	case "owner":
		return RoleOwner
	case "editor":
		return RoleEditor
	default:
		return RoleViewer
	}
}

// String returns the DB enum value for a Role.
func (r Role) String() string {
	switch r {
	case RoleOwner:
		return "owner"
	case RoleEditor:
		return "editor"
	default:
		return "viewer"
	}
}

// User represents an application user.
type User struct {
	ID            uuid.UUID
	Email         string
	PasswordHash  string
	FirstName     string
	LastName      string
	Phone         string
	AvatarURL     string
	IsActive      bool
	IsSystemAdmin bool
	LastLoginAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

// FullName returns the user's display name.
func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Email
	}
	return u.FirstName + " " + u.LastName
}

// VenueUserRole records a user's role within a specific venue.
type VenueUserRole struct {
	ID        uuid.UUID
	VenueID   uuid.UUID
	UserID    uuid.UUID
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// InvitationStatus represents the state of a venue invitation.
type InvitationStatus string

const (
	InvitationStatusPending   InvitationStatus = "pending"
	InvitationStatusAccepted  InvitationStatus = "accepted"
	InvitationStatusCancelled InvitationStatus = "cancelled"
)

// VenueInvitation is a pending invite for an email address to join a venue.
type VenueInvitation struct {
	ID          uuid.UUID
	VenueID     uuid.UUID
	Email       string
	Role        Role
	Token       string
	InvitedBy   uuid.UUID
	Status      InvitationStatus
	AcceptedAt  *time.Time
	CancelledAt *time.Time
	ExpiresAt   time.Time
	CreatedAt   time.Time
	DeletedAt   *time.Time
}

// RefreshToken is the stored record of an issued refresh token.
// Only the SHA-256 hash of the raw token is persisted.
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// ResetPasswordToken is a one-time token for password reset flows.
type ResetPasswordToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}
