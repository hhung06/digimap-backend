package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/email"
	"github.com/hhung06/digimap-backend/internal/repository"
)

const invitationTTL = 7 * 24 * time.Hour // 7 days

// UserService handles venue-scoped user management and invitation flows.
type UserService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	UpdateProfile(ctx context.Context, u *domain.User) error

	// Venue users
	ListVenueUsers(ctx context.Context, venueID uuid.UUID) ([]*domain.User, error)
	ChangeRole(ctx context.Context, venueID, userID uuid.UUID, role domain.Role) error
	RemoveFromVenue(ctx context.Context, venueID, userID uuid.UUID) error

	// Invitations
	ListInvitations(ctx context.Context, venueID uuid.UUID) ([]*domain.VenueInvitation, error)
	InviteUser(ctx context.Context, venueID, invitedBy uuid.UUID, email string, role domain.Role) (*domain.VenueInvitation, error)
	AcceptInvitation(ctx context.Context, token string, userID uuid.UUID) error
	CancelInvitation(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	userRepo repository.UserRepository
	mailer   email.Sender
}

// NewUserService creates a UserService.
func NewUserService(userRepo repository.UserRepository, mailer email.Sender) UserService {
	return &userService{userRepo: userRepo, mailer: mailer}
}

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

func (s *userService) UpdateProfile(ctx context.Context, u *domain.User) error {
	return s.userRepo.Update(ctx, u)
}

func (s *userService) ListVenueUsers(ctx context.Context, venueID uuid.UUID) ([]*domain.User, error) {
	return s.userRepo.ListVenueUsers(ctx, venueID)
}

func (s *userService) ChangeRole(ctx context.Context, venueID, userID uuid.UUID, role domain.Role) error {
	vr := &domain.VenueUserRole{
		VenueID: venueID,
		UserID:  userID,
		Role:    role,
	}
	return s.userRepo.UpsertVenueRole(ctx, vr)
}

func (s *userService) RemoveFromVenue(ctx context.Context, venueID, userID uuid.UUID) error {
	return s.userRepo.DeleteVenueRole(ctx, venueID, userID)
}

func (s *userService) ListInvitations(ctx context.Context, venueID uuid.UUID) ([]*domain.VenueInvitation, error) {
	return s.userRepo.ListInvitations(ctx, venueID)
}

func (s *userService) InviteUser(ctx context.Context, venueID, invitedBy uuid.UUID, emailAddr string, role domain.Role) (*domain.VenueInvitation, error) {
	token, err := generateInviteToken()
	if err != nil {
		return nil, fmt.Errorf("generate invite token: %w", err)
	}

	inv := &domain.VenueInvitation{
		VenueID:   venueID,
		InvitedBy: invitedBy,
		Email:     emailAddr,
		Role:      role,
		Token:     token,
		Status:    domain.InvitationStatusPending,
		ExpiresAt: time.Now().Add(invitationTTL),
	}

	if err := s.userRepo.CreateInvitation(ctx, inv); err != nil {
		return nil, err
	}

	// Best-effort email; failure does not roll back the invitation record.
	_ = s.mailer.SendInvitation(ctx, emailAddr, "", token)
	return inv, nil
}

func (s *userService) AcceptInvitation(ctx context.Context, token string, userID uuid.UUID) error {
	inv, err := s.userRepo.FindInvitationByToken(ctx, token)
	if err != nil {
		return err
	}

	if inv.Status != domain.InvitationStatusPending {
		return domain.NewConflict("invitation is no longer pending")
	}
	if time.Now().After(inv.ExpiresAt) {
		return domain.NewConflict("invitation has expired")
	}

	now := time.Now()
	if err := s.userRepo.UpdateInvitationStatus(ctx, inv.ID, domain.InvitationStatusAccepted, &now, nil); err != nil {
		return err
	}

	vr := &domain.VenueUserRole{
		VenueID: inv.VenueID,
		UserID:  userID,
		Role:    inv.Role,
	}
	return s.userRepo.UpsertVenueRole(ctx, vr)
}

func (s *userService) CancelInvitation(ctx context.Context, id uuid.UUID) error {
	inv, err := s.userRepo.FindInvitationByID(ctx, id)
	if err != nil {
		return err
	}
	if inv.Status != domain.InvitationStatusPending {
		return domain.NewConflict("only pending invitations can be cancelled")
	}

	now := time.Now()
	return s.userRepo.UpdateInvitationStatus(ctx, id, domain.InvitationStatusCancelled, nil, &now)
}

// generateInviteToken returns a 32-byte hex token (64 chars).
func generateInviteToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
