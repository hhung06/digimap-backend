package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func newTestVenueService(venueRepo *mocks.VenueRepository, levelRepo *mocks.LevelRepository) service.VenueService {
	return service.NewVenueService(venueRepo, levelRepo, nil)
}

// ── Get ───────────────────────────────────────────────────────────────────────

func TestVenueService_Get_Success(t *testing.T) {
	venueRepo := &mocks.VenueRepository{}
	svc := newTestVenueService(venueRepo, nil)

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.Venue{ID: id, Name: "Grand Mall"}

	venueRepo.On("FindByID", ctx, id).Return(expected, nil)

	got, err := svc.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	venueRepo.AssertExpectations(t)
}

func TestVenueService_Get_NotFound(t *testing.T) {
	venueRepo := &mocks.VenueRepository{}
	svc := newTestVenueService(venueRepo, nil)

	ctx := context.Background()
	id := uuid.New()

	venueRepo.On("FindByID", ctx, id).Return((*domain.Venue)(nil), domain.NewNotFound("venue not found"))

	_, err := svc.Get(ctx, id)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	venueRepo.AssertExpectations(t)
}

// ── GetByPublicKey ────────────────────────────────────────────────────────────

func TestVenueService_GetByPublicKey_Success(t *testing.T) {
	venueRepo := &mocks.VenueRepository{}
	svc := newTestVenueService(venueRepo, nil)

	ctx := context.Background()
	pubKey := "abc123"
	expected := &domain.Venue{PublicKey: pubKey}

	venueRepo.On("FindByPublicKey", ctx, pubKey).Return(expected, nil)

	got, err := svc.GetByPublicKey(ctx, pubKey)
	require.NoError(t, err)
	assert.Equal(t, pubKey, got.PublicKey)
	venueRepo.AssertExpectations(t)
}

// ── List ──────────────────────────────────────────────────────────────────────

func TestVenueService_List_Success(t *testing.T) {
	venueRepo := &mocks.VenueRepository{}
	svc := newTestVenueService(venueRepo, nil)

	ctx := context.Background()
	customerID := uuid.New()
	p := domain.Pagination{Page: 1, PageSize: 20}
	expected := []*domain.Venue{{Name: "Mall A"}, {Name: "Mall B"}}

	venueRepo.On("List", ctx, customerID, p).Return(expected, int64(2), nil)

	got, total, err := svc.List(ctx, customerID, p)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, int64(2), total)
	venueRepo.AssertExpectations(t)
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestVenueService_Create_Success(t *testing.T) {
	venueRepo := &mocks.VenueRepository{}
	svc := newTestVenueService(venueRepo, nil)

	ctx := context.Background()
	v := &domain.Venue{Name: "New Mall"}

	venueRepo.On("Create", ctx, v).Return(nil)

	err := svc.Create(ctx, v)
	require.NoError(t, err)
	assert.NotEmpty(t, v.PublicKey, "public key must be generated")
	assert.NotEmpty(t, v.PrivateKey, "private key must be generated")
	venueRepo.AssertExpectations(t)
}

func TestVenueService_Create_RepoError(t *testing.T) {
	venueRepo := &mocks.VenueRepository{}
	svc := newTestVenueService(venueRepo, nil)

	ctx := context.Background()
	v := &domain.Venue{Name: "New Mall"}
	repoErr := errors.New("db error")

	venueRepo.On("Create", ctx, v).Return(repoErr)

	err := svc.Create(ctx, v)
	require.Error(t, err)
	venueRepo.AssertExpectations(t)
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestVenueService_Update_Success(t *testing.T) {
	venueRepo := &mocks.VenueRepository{}
	svc := newTestVenueService(venueRepo, nil)

	ctx := context.Background()
	id := uuid.New()
	customerID := uuid.New()
	existing := &domain.Venue{ID: id, PublicKey: "pub", PrivateKey: "priv", CustomerID: customerID}
	updated := &domain.Venue{ID: id, Name: "Updated Mall"}

	venueRepo.On("FindByID", ctx, id).Return(existing, nil)
	venueRepo.On("Update", ctx, updated).Return(nil)

	err := svc.Update(ctx, updated)
	require.NoError(t, err)
	assert.Equal(t, "pub", updated.PublicKey, "public key must be preserved")
	assert.Equal(t, "priv", updated.PrivateKey, "private key must be preserved")
	assert.Equal(t, customerID, updated.CustomerID, "customerID must be preserved")
	venueRepo.AssertExpectations(t)
}

func TestVenueService_Update_NotFound(t *testing.T) {
	venueRepo := &mocks.VenueRepository{}
	svc := newTestVenueService(venueRepo, nil)

	ctx := context.Background()
	id := uuid.New()

	venueRepo.On("FindByID", ctx, id).Return((*domain.Venue)(nil), domain.NewNotFound("venue not found"))

	err := svc.Update(ctx, &domain.Venue{ID: id})
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	venueRepo.AssertExpectations(t)
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestVenueService_Delete_Success(t *testing.T) {
	venueRepo := &mocks.VenueRepository{}
	svc := newTestVenueService(venueRepo, nil)

	ctx := context.Background()
	id := uuid.New()

	venueRepo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	venueRepo.AssertExpectations(t)
}

// ── RegenerateKeys ────────────────────────────────────────────────────────────

func TestVenueService_RegenerateKeys_Success(t *testing.T) {
	venueRepo := &mocks.VenueRepository{}
	svc := newTestVenueService(venueRepo, nil)

	ctx := context.Background()
	id := uuid.New()
	updated := &domain.Venue{ID: id, PublicKey: "newpub", PrivateKey: "newpriv"}

	venueRepo.On("UpdateKeys", ctx, id, mockAny, mockAny).Return(nil)
	venueRepo.On("FindByID", ctx, id).Return(updated, nil)

	got, err := svc.RegenerateKeys(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	venueRepo.AssertExpectations(t)
}

func TestVenueService_RegenerateKeys_UpdateError(t *testing.T) {
	venueRepo := &mocks.VenueRepository{}
	svc := newTestVenueService(venueRepo, nil)

	ctx := context.Background()
	id := uuid.New()

	venueRepo.On("UpdateKeys", ctx, id, mockAny, mockAny).Return(errors.New("db error"))

	_, err := svc.RegenerateKeys(ctx, id)
	require.Error(t, err)
	venueRepo.AssertExpectations(t)
}

// ── Clone ─────────────────────────────────────────────────────────────────────

func TestVenueService_Clone_SourceNotFound(t *testing.T) {
	venueRepo := &mocks.VenueRepository{}
	svc := newTestVenueService(venueRepo, &mocks.LevelRepository{})

	ctx := context.Background()
	sourceID := uuid.New()

	venueRepo.On("FindByID", ctx, sourceID).Return((*domain.Venue)(nil), domain.NewNotFound("venue not found"))

	_, err := svc.Clone(ctx, sourceID)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	venueRepo.AssertExpectations(t)
}

// ── PBT: key invariants ───────────────────────────────────────────────────────

func TestVenueService_Create_KeyInvariants(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		venueRepo := &mocks.VenueRepository{}
		svc := newTestVenueService(venueRepo, nil)

		ctx := context.Background()
		v := &domain.Venue{Name: rapid.StringN(1, 50, 50).Draw(rt, "name")}

		venueRepo.On("Create", ctx, v).Return(nil)

		err := svc.Create(ctx, v)
		require.NoError(rt, err)
		// public key: base64url of 32 bytes = 43 chars (no padding)
		assert.Equal(rt, 43, len(v.PublicKey), "public key must always be 43 chars")
		// private key: hex of sha256 = 64 chars
		assert.Equal(rt, 64, len(v.PrivateKey), "private key must always be 64 chars")
	})
}
