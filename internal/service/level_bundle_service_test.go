package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func TestLevelBundleService_ListBySnapshot(t *testing.T) {
	repo := &mocks.LevelBundleRepository{}
	snapRepo := &mocks.SnapshotRepository{}
	storer := &mocks.StorerMock{}
	svc := service.NewLevelBundleService(repo, snapRepo, storer, "local")

	ctx := context.Background()
	snapshotID := uuid.New()
	expected := []*domain.LevelBundle{{SnapshotID: snapshotID, State: domain.LevelBundleStateDraft}}

	repo.On("ListBySnapshot", ctx, snapshotID).Return(expected, nil)

	bundles, err := svc.ListBySnapshot(ctx, snapshotID)
	require.NoError(t, err)
	assert.Len(t, bundles, 1)
	repo.AssertExpectations(t)
}

func TestLevelBundleService_Create_InheritsParentState(t *testing.T) {
	repo := &mocks.LevelBundleRepository{}
	snapRepo := &mocks.SnapshotRepository{}
	storer := &mocks.StorerMock{}
	svc := service.NewLevelBundleService(repo, snapRepo, storer, "local")

	ctx := context.Background()
	snapshotID := uuid.New()
	venueID := uuid.New()
	levelID := uuid.New()
	bundle := []byte(`{"tiles":[]}`)
	parentSnap := &domain.Snapshot{ID: snapshotID, VenueID: venueID, State: domain.SnapshotStateDraft}

	snapRepo.On("FindByID", ctx, snapshotID).Return(parentSnap, nil)
	storer.On("PutObject", ctx, mock.AnythingOfType("string"), bundle).Return(nil)
	repo.On("Create", ctx, mock.AnythingOfType("*domain.LevelBundle")).Return(nil)

	b, err := svc.Create(ctx, snapshotID, venueID, levelID, bundle)
	require.NoError(t, err)
	assert.Equal(t, snapshotID, b.SnapshotID)
	assert.Equal(t, venueID, b.VenueID)
	assert.Equal(t, levelID, b.LevelID)
	// State must match parent snapshot's state
	assert.Equal(t, domain.LevelBundleStateDraft, b.State)
	snapRepo.AssertExpectations(t)
	storer.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestLevelBundleService_Create_S3Failure(t *testing.T) {
	repo := &mocks.LevelBundleRepository{}
	snapRepo := &mocks.SnapshotRepository{}
	storer := &mocks.StorerMock{}
	svc := service.NewLevelBundleService(repo, snapRepo, storer, "local")

	ctx := context.Background()
	snapshotID := uuid.New()
	venueID := uuid.New()
	levelID := uuid.New()
	bundle := []byte(`{"tiles":[]}`)
	parentSnap := &domain.Snapshot{ID: snapshotID, VenueID: venueID, State: domain.LevelBundleStateDraft}

	snapRepo.On("FindByID", ctx, snapshotID).Return(parentSnap, nil)
	storer.On("PutObject", ctx, mock.AnythingOfType("string"), bundle).Return(assert.AnError)

	_, err := svc.Create(ctx, snapshotID, venueID, levelID, bundle)
	require.Error(t, err)
	repo.AssertNotCalled(t, "Create")
}

func TestLevelBundleService_Create_SnapshotNotFound(t *testing.T) {
	repo := &mocks.LevelBundleRepository{}
	snapRepo := &mocks.SnapshotRepository{}
	storer := &mocks.StorerMock{}
	svc := service.NewLevelBundleService(repo, snapRepo, storer, "local")

	ctx := context.Background()
	snapshotID := uuid.New()

	snapRepo.On("FindByID", ctx, snapshotID).Return((*domain.Snapshot)(nil), domain.NewNotFound("snapshot not found"))

	_, err := svc.Create(ctx, snapshotID, uuid.New(), uuid.New(), []byte(`{}`))
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestLevelBundleService_Delete(t *testing.T) {
	repo := &mocks.LevelBundleRepository{}
	snapRepo := &mocks.SnapshotRepository{}
	storer := &mocks.StorerMock{}
	svc := service.NewLevelBundleService(repo, snapRepo, storer, "local")

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}
