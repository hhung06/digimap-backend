package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func newTestUserService(userRepo *mocks.UserRepository, mailer *mocks.EmailSender) service.UserService {
	return service.NewUserService(userRepo, mailer)
}

// ── GetProfile ────────────────────────────────────────────────────────────────

func TestUserService_GetProfile_Success(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	svc := newTestUserService(userRepo, nil)

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.User{ID: id, Email: "alice@example.com"}

	userRepo.On("FindByID", ctx, id).Return(expected, nil)

	got, err := svc.GetProfile(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	userRepo.AssertExpectations(t)
}

func TestUserService_GetProfile_NotFound(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	svc := newTestUserService(userRepo, nil)

	ctx := context.Background()
	id := uuid.New()

	userRepo.On("FindByID", ctx, id).Return((*domain.User)(nil), domain.NewNotFound("user not found"))

	_, err := svc.GetProfile(ctx, id)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	userRepo.AssertExpectations(t)
}

// ── UpdateProfile ─────────────────────────────────────────────────────────────

func TestUserService_UpdateProfile_Success(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	svc := newTestUserService(userRepo, nil)

	ctx := context.Background()
	u := &domain.User{ID: uuid.New(), Email: "alice@example.com"}

	userRepo.On("Update", ctx, u).Return(nil)

	err := svc.UpdateProfile(ctx, u)
	require.NoError(t, err)
	userRepo.AssertExpectations(t)
}

// ── ListVenueUsers ────────────────────────────────────────────────────────────

func TestUserService_ListVenueUsers_Success(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	svc := newTestUserService(userRepo, nil)

	ctx := context.Background()
	venueID := uuid.New()
	expected := []*domain.User{{Email: "alice@example.com"}, {Email: "bob@example.com"}}

	userRepo.On("ListVenueUsers", ctx, venueID).Return(expected, nil)

	got, err := svc.ListVenueUsers(ctx, venueID)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	userRepo.AssertExpectations(t)
}

// ── ChangeRole ────────────────────────────────────────────────────────────────

func TestUserService_ChangeRole_Success(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	svc := newTestUserService(userRepo, nil)

	ctx := context.Background()
	venueID, userID := uuid.New(), uuid.New()

	userRepo.On("UpsertVenueRole", ctx, mockAny).Return(nil)

	err := svc.ChangeRole(ctx, venueID, userID, domain.RoleEditor)
	require.NoError(t, err)
	userRepo.AssertExpectations(t)
}

// ── RemoveFromVenue ───────────────────────────────────────────────────────────

func TestUserService_RemoveFromVenue_Success(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	svc := newTestUserService(userRepo, nil)

	ctx := context.Background()
	venueID, userID := uuid.New(), uuid.New()

	userRepo.On("DeleteVenueRole", ctx, venueID, userID).Return(nil)

	err := svc.RemoveFromVenue(ctx, venueID, userID)
	require.NoError(t, err)
	userRepo.AssertExpectations(t)
}

// ── ListInvitations ───────────────────────────────────────────────────────────

func TestUserService_ListInvitations_Success(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	svc := newTestUserService(userRepo, nil)

	ctx := context.Background()
	venueID := uuid.New()
	expected := []*domain.VenueInvitation{{Email: "alice@example.com"}}

	userRepo.On("ListInvitations", ctx, venueID).Return(expected, nil)

	got, err := svc.ListInvitations(ctx, venueID)
	require.NoError(t, err)
	assert.Len(t, got, 1)
	userRepo.AssertExpectations(t)
}

// ── InviteUser ────────────────────────────────────────────────────────────────

func TestUserService_InviteUser_Success(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	mailer := &mocks.EmailSender{}
	svc := newTestUserService(userRepo, mailer)

	ctx := context.Background()
	venueID, invitedBy := uuid.New(), uuid.New()

	userRepo.On("CreateInvitation", ctx, mockAny).Return(nil)
	// Email is best-effort; failure not surfaced — but it should be called
	mailer.On("SendInvitation", ctx, "alice@example.com", "", mockAny).Return(nil)

	inv, err := svc.InviteUser(ctx, venueID, invitedBy, "alice@example.com", domain.RoleViewer)
	require.NoError(t, err)
	assert.NotEmpty(t, inv.Token, "invitation token must be generated")
	assert.Equal(t, domain.InvitationStatusPending, inv.Status)
	userRepo.AssertExpectations(t)
	mailer.AssertExpectations(t)
}

func TestUserService_InviteUser_RepoError(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	mailer := &mocks.EmailSender{}
	svc := newTestUserService(userRepo, mailer)

	ctx := context.Background()
	venueID, invitedBy := uuid.New(), uuid.New()

	userRepo.On("CreateInvitation", ctx, mockAny).Return(errors.New("db error"))

	_, err := svc.InviteUser(ctx, venueID, invitedBy, "alice@example.com", domain.RoleViewer)
	require.Error(t, err)
	// Mailer must NOT be called if repo fails
	mailer.AssertNotCalled(t, "SendInvitation")
	userRepo.AssertExpectations(t)
}

// ── AcceptInvitation ──────────────────────────────────────────────────────────

func TestUserService_AcceptInvitation_Success(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	svc := newTestUserService(userRepo, nil)

	ctx := context.Background()
	token := "valid-token"
	userID := uuid.New()
	inv := &domain.VenueInvitation{
		ID:        uuid.New(),
		VenueID:   uuid.New(),
		Status:    domain.InvitationStatusPending,
		ExpiresAt: time.Now().Add(time.Hour),
		Role:      domain.RoleViewer,
	}

	userRepo.On("FindInvitationByToken", ctx, token).Return(inv, nil)
	userRepo.On("UpdateInvitationStatus", ctx, inv.ID, domain.InvitationStatusAccepted, mockAny, (*time.Time)(nil)).Return(nil)
	userRepo.On("UpsertVenueRole", ctx, mockAny).Return(nil)

	err := svc.AcceptInvitation(ctx, token, userID)
	require.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestUserService_AcceptInvitation_AlreadyAccepted(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	svc := newTestUserService(userRepo, nil)

	ctx := context.Background()
	token := "used-token"
	inv := &domain.VenueInvitation{
		Status:    domain.InvitationStatusAccepted,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	userRepo.On("FindInvitationByToken", ctx, token).Return(inv, nil)

	err := svc.AcceptInvitation(ctx, token, uuid.New())
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrConflict))
	userRepo.AssertExpectations(t)
}

func TestUserService_AcceptInvitation_Expired(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	svc := newTestUserService(userRepo, nil)

	ctx := context.Background()
	token := "expired-token"
	inv := &domain.VenueInvitation{
		Status:    domain.InvitationStatusPending,
		ExpiresAt: time.Now().Add(-time.Hour),
	}

	userRepo.On("FindInvitationByToken", ctx, token).Return(inv, nil)

	err := svc.AcceptInvitation(ctx, token, uuid.New())
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrConflict))
	userRepo.AssertExpectations(t)
}

// ── CancelInvitation ──────────────────────────────────────────────────────────

func TestUserService_CancelInvitation_Success(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	svc := newTestUserService(userRepo, nil)

	ctx := context.Background()
	id := uuid.New()
	inv := &domain.VenueInvitation{ID: id, Status: domain.InvitationStatusPending}

	userRepo.On("FindInvitationByID", ctx, id).Return(inv, nil)
	userRepo.On("UpdateInvitationStatus", ctx, id, domain.InvitationStatusCancelled, (*time.Time)(nil), mockAny).Return(nil)

	err := svc.CancelInvitation(ctx, id)
	require.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestUserService_CancelInvitation_NotPending(t *testing.T) {
	userRepo := &mocks.UserRepository{}
	svc := newTestUserService(userRepo, nil)

	ctx := context.Background()
	id := uuid.New()
	inv := &domain.VenueInvitation{ID: id, Status: domain.InvitationStatusAccepted}

	userRepo.On("FindInvitationByID", ctx, id).Return(inv, nil)

	err := svc.CancelInvitation(ctx, id)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrConflict))
	userRepo.AssertExpectations(t)
}
