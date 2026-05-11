package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
)

func TestSnapshotService_Publish(t *testing.T) {
	repo := &mocks.SnapshotRepository{}
	svc := newTestSnapshotSvc(repo)

	snapID := uuid.New()
	venueID := uuid.New()
	existing := &domain.Snapshot{
		ID:      snapID,
		VenueID: venueID,
		State:   domain.SnapshotStateDraft,
	}

	repo.On("FindByID", mock.Anything, snapID).Return(existing, nil)
	repo.On("UpdateState", mock.Anything, snapID, domain.SnapshotStatePublic, mock.AnythingOfType("*time.Time")).Return(nil)

	snap, err := svc.Publish(context.Background(), snapID)

	require.NoError(t, err)
	assert.Equal(t, domain.SnapshotStatePublic, snap.State)
	assert.NotNil(t, snap.PublishAt)
	repo.AssertExpectations(t)
}

func TestSnapshotService_Publish_NotFound(t *testing.T) {
	repo := &mocks.SnapshotRepository{}
	svc := newTestSnapshotSvc(repo)

	snapID := uuid.New()
	notFound := domain.NewNotFound("snapshot not found")

	repo.On("FindByID", mock.Anything, snapID).Return((*domain.Snapshot)(nil), notFound)

	snap, err := svc.Publish(context.Background(), snapID)

	assert.Nil(t, snap)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	repo.AssertExpectations(t)
}

func TestSnapshotService_Revert(t *testing.T) {
	repo := &mocks.SnapshotRepository{}
	svc := newTestSnapshotSvc(repo)
	ctx := context.Background()

	venueID := uuid.New()
	targetID := uuid.New()
	target := &domain.Snapshot{
		ID:      targetID,
		VenueID: venueID,
		State:   domain.SnapshotStateDraft,
	}

	repo.On("FindByID", ctx, targetID).Return(target, nil)
	repo.On("UnpublishVenue", ctx, venueID).Return(nil)
	repo.On("UpdateState", ctx, targetID, domain.SnapshotStatePublic, mock.AnythingOfType("*time.Time")).Return(nil)

	result, err := svc.Revert(ctx, targetID)

	require.NoError(t, err)
	assert.Equal(t, domain.SnapshotStatePublic, result.State)
	repo.AssertExpectations(t)
}

// TestSnapshotService_AutoPublish_ReturnsErrorWhenNoDraft verifies that AutoPublish
// returns an error when there is no existing draft for the venue (new behavior).
func TestSnapshotService_AutoPublish_ReturnsErrorWhenNoDraft(t *testing.T) {
	repo := &mocks.SnapshotRepository{}
	svc := newTestSnapshotSvc(repo)
	ctx := context.Background()
	venueID := uuid.New()
	userID := uuid.New()

	repo.On("LatestDraft", ctx, venueID).Return((*domain.Snapshot)(nil), nil)

	_, err := svc.AutoPublish(ctx, venueID, userID)

	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	repo.AssertExpectations(t)
}

// TestSnapshotService_AutoPublish_PublishesExistingDraft verifies that AutoPublish
// republishes the latest existing draft snapshot.
func TestSnapshotService_AutoPublish_PublishesExistingDraft(t *testing.T) {
	repo := &mocks.SnapshotRepository{}
	svc := newTestSnapshotSvc(repo)
	ctx := context.Background()
	venueID := uuid.New()
	snapID := uuid.New()
	userID := uuid.New()

	draft := &domain.Snapshot{
		ID:      snapID,
		VenueID: venueID,
		State:   domain.SnapshotStateDraft,
	}

	repo.On("LatestDraft", ctx, venueID).Return(draft, nil)
	repo.On("UnpublishVenue", ctx, venueID).Return(nil)
	repo.On("UpdateState", ctx, snapID, domain.SnapshotStatePublic, mock.AnythingOfType("*time.Time")).Return(nil)

	result, err := svc.AutoPublish(ctx, venueID, userID)

	require.NoError(t, err)
	assert.Equal(t, domain.SnapshotStatePublic, result.State)
	assert.Equal(t, snapID, result.ID)
	repo.AssertExpectations(t)
}
