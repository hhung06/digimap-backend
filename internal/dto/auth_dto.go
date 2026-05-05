package dto

import (
	"time"

	"github.com/google/uuid"
)

// ── Requests ──────────────────────────────────────────────────────────────────

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email     string `json:"email"      binding:"required,email"`
	Password  string `json:"password"   binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"  binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PasswordResetConfirmRequest struct {
	Token       string `json:"token"        binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ── Responses ─────────────────────────────────────────────────────────────────

type UserResponse struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	Phone         string     `json:"phone,omitempty"`
	AvatarURL     string     `json:"avatar_url,omitempty"`
	IsActive      bool       `json:"is_active"`
	IsSystemAdmin bool       `json:"is_system_admin"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type LoginResponse struct {
	AccessToken           string       `json:"access_token"`
	AccessExpiresIn       int          `json:"access_expires_in"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshExpiresIn      int          `json:"refresh_expires_in"`
	TokenType             string       `json:"token_type"`
	User                  UserResponse `json:"user"`
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	AccessExpiresIn  int    `json:"access_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
}
