package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/hhung06/digimap-backend/config"
	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/email"
	"github.com/hhung06/digimap-backend/internal/repository"
)

// Claims is the JWT payload for access tokens.
type Claims struct {
	UserID        string `json:"sub"`
	Email         string `json:"email"`
	IsSystemAdmin bool   `json:"is_system_admin"`
	jwt.RegisteredClaims
}

// TokenPair holds both tokens issued at login or refresh.
type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresIn  int // seconds
	RefreshExpiresIn int // seconds
}

// AuthService handles all authentication and session management.
type AuthService interface {
	Login(ctx context.Context, email, password string) (*domain.User, TokenPair, error)
	Register(ctx context.Context, req RegisterRequest) (*domain.User, error)
	RefreshToken(ctx context.Context, rawRefreshToken string) (TokenPair, error)
	Logout(ctx context.Context, rawRefreshToken string) error
	RequestPasswordReset(ctx context.Context, email string) error
	ConfirmPasswordReset(ctx context.Context, token, newPassword string) error
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	// ValidateClaims parses and validates a JWT access token string.
	ValidateClaims(tokenString string) (*Claims, error)
}

// RegisterRequest contains the fields needed to create a new user.
type RegisterRequest struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type authService struct {
	users  repository.UserRepository
	tokens repository.TokenRepository
	mailer email.Sender
	cfg    config.JWTConfig
}

// NewAuthService creates an AuthService.
func NewAuthService(
	users repository.UserRepository,
	tokens repository.TokenRepository,
	mailer email.Sender,
	cfg config.JWTConfig,
) AuthService {
	return &authService{users: users, tokens: tokens, mailer: mailer, cfg: cfg}
}

// ── Login ─────────────────────────────────────────────────────────────────────

func (s *authService) Login(ctx context.Context, emailAddr, password string) (*domain.User, TokenPair, error) {
	u, err := s.users.FindByEmail(ctx, emailAddr)
	if err != nil {
		// Return generic message to prevent email enumeration
		return nil, TokenPair{}, domain.NewUnauthorized("invalid email or password")
	}

	if !u.IsActive {
		return nil, TokenPair{}, domain.NewUnauthorized("account is disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, TokenPair{}, domain.NewUnauthorized("invalid email or password")
	}

	pair, err := s.issuePair(ctx, u)
	if err != nil {
		return nil, TokenPair{}, err
	}

	_ = s.users.UpdateLastLogin(ctx, u.ID)
	return u, pair, nil
}

// ── Register ──────────────────────────────────────────────────────────────────

func (s *authService) Register(ctx context.Context, req RegisterRequest) (*domain.User, error) {
	if len(req.Password) < 8 {
		return nil, domain.NewValidation(map[string]string{"password": "must be at least 8 characters"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     true,
	}

	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

// ── Refresh ───────────────────────────────────────────────────────────────────

func (s *authService) RefreshToken(ctx context.Context, rawRefreshToken string) (TokenPair, error) {
	hash := hashToken(rawRefreshToken)
	stored, err := s.tokens.FindRefreshToken(ctx, hash)
	if err != nil {
		return TokenPair{}, err
	}

	u, err := s.users.FindByID(ctx, stored.UserID)
	if err != nil {
		return TokenPair{}, err
	}
	if !u.IsActive {
		return TokenPair{}, domain.NewUnauthorized("account is disabled")
	}

	// Revoke the used token (rotation) before issuing a new pair, so a stolen
	// refresh token can be replayed at most once.
	if err := s.tokens.RevokeRefreshToken(ctx, stored.ID); err != nil {
		return TokenPair{}, fmt.Errorf("revoke old refresh token: %w", err)
	}

	return s.issuePair(ctx, u)
}

// ── Logout ────────────────────────────────────────────────────────────────────

func (s *authService) Logout(ctx context.Context, rawRefreshToken string) error {
	hash := hashToken(rawRefreshToken)
	stored, err := s.tokens.FindRefreshToken(ctx, hash)
	if err != nil {
		return err
	}
	return s.tokens.RevokeRefreshToken(ctx, stored.ID)
}

// ── Password reset ────────────────────────────────────────────────────────────

func (s *authService) RequestPasswordReset(ctx context.Context, emailAddr string) error {
	u, err := s.users.FindByEmail(ctx, emailAddr)
	if err != nil {
		// Silently succeed to prevent email enumeration
		return nil
	}

	rawToken, err := generateSecureToken()
	if err != nil {
		return fmt.Errorf("generate reset token: %w", err)
	}

	rec := &domain.ResetPasswordToken{
		UserID:    u.ID,
		TokenHash: hashToken(rawToken),
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	if err := s.tokens.CreateResetToken(ctx, rec); err != nil {
		return fmt.Errorf("store reset token: %w", err)
	}

	return s.mailer.SendPasswordReset(ctx, emailAddr, rawToken)
}

func (s *authService) ConfirmPasswordReset(ctx context.Context, token, newPassword string) error {
	if len(newPassword) < 8 {
		return domain.NewValidation(map[string]string{"password": "must be at least 8 characters"})
	}

	hash := hashToken(token)
	rec, err := s.tokens.FindResetToken(ctx, hash)
	if err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	u, err := s.users.FindByID(ctx, rec.UserID)
	if err != nil {
		return err
	}
	u.PasswordHash = string(passwordHash)
	if err := s.users.Update(ctx, u); err != nil {
		return err
	}

	_ = s.tokens.MarkResetTokenUsed(ctx, rec.ID)
	_ = s.tokens.RevokeAllUserRefreshTokens(ctx, u.ID)
	return nil
}

// ── Change password ───────────────────────────────────────────────────────────

func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	if len(newPassword) < 8 {
		return domain.NewValidation(map[string]string{"password": "must be at least 8 characters"})
	}

	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(oldPassword)); err != nil {
		return domain.NewUnauthorized("current password is incorrect")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	return s.users.UpdatePassword(ctx, userID, string(hash))
}

// ── JWT helpers ───────────────────────────────────────────────────────────────

func (s *authService) ValidateClaims(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, domain.NewUnauthorized("invalid or expired token")
	}
	return claims, nil
}

func (s *authService) issueAccessToken(u *domain.User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:        u.ID.String(),
		Email:         u.Email,
		IsSystemAdmin: u.IsSystemAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   u.ID.String(),
			Issuer:    s.cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.AccessExpiry)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.SecretKey))
}

func (s *authService) issuePair(ctx context.Context, u *domain.User) (TokenPair, error) {
	accessToken, err := s.issueAccessToken(u)
	if err != nil {
		return TokenPair{}, fmt.Errorf("sign access token: %w", err)
	}

	rawRefresh, err := generateSecureToken()
	if err != nil {
		return TokenPair{}, fmt.Errorf("generate refresh token: %w", err)
	}

	rec := &domain.RefreshToken{
		UserID:    u.ID,
		TokenHash: hashToken(rawRefresh),
		ExpiresAt: time.Now().Add(s.cfg.RefreshExpiry),
	}
	if err := s.tokens.CreateRefreshToken(ctx, rec); err != nil {
		return TokenPair{}, fmt.Errorf("store refresh token: %w", err)
	}

	return TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     rawRefresh,
		AccessExpiresIn:  int(s.cfg.AccessExpiry.Seconds()),
		RefreshExpiresIn: int(s.cfg.RefreshExpiry.Seconds()),
	}, nil
}

// ── Utility ───────────────────────────────────────────────────────────────────

// hashToken returns the hex-encoded SHA-256 hash of a token string.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// generateSecureToken returns a 32-byte cryptographically random hex string.
func generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
