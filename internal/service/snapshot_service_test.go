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

func TestSnapshotService_List(t *testing.T) {
	repo := &mocks.SnapshotRepository{}
	storer := &mocks.StorerMock{}
	svc := service.NewSnapshotService(repo, storer, "local")

	ctx := context.Background()
	venueID := uuid.New()
	p := domain.Pagination{Page: 1, PageSize: 20}
	expected := []*domain.Snapshot{{VenueID: venueID, State: domain.SnapshotStateDraft}}

	repo.On("List", ctx, venueID, p).Return(expected, int64(1), nil)

	snapshots, total, err := svc.List(ctx, venueID, p)
	require.NoError(t, err)
	assert.Len(t, snapshots, 1)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}

func TestSnapshotService_Get_NotFound(t *testing.T) {
	repo := &mocks.SnapshotRepository{}
	storer := &mocks.StorerMock{}
	svc := service.NewSnapshotService(repo, storer, "local")

	ctx := context.Background()
	id := uuid.New()

	repo.On("FindByID", ctx, id).Return((*domain.Snapshot)(nil), domain.NewNotFound("snapshot not found"))

	_, err := svc.Get(ctx, id)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)
	repo.AssertExpectations(t)
}

func TestSnapshotService_LatestPublished(t *testing.T) {
	repo := &mocks.SnapshotRepository{}
	storer := &mocks.StorerMock{}
	svc := service.NewSnapshotService(repo, storer, "local")

	ctx := context.Background()
	venueID := uuid.New()
	expected := &domain.Snapshot{VenueID: venueID, State: domain.SnapshotStatePublic}

	repo.On("LatestPublished", ctx, venueID).Return(expected, nil)

	s, err := svc.LatestPublished(ctx, venueID)
	require.NoError(t, err)
	assert.Equal(t, domain.SnapshotStatePublic, s.State)
	repo.AssertExpectations(t)
}

func TestSnapshotService_CreateDraft_Success(t *testing.T) {
	repo := &mocks.SnapshotRepository{}
	storer := &mocks.StorerMock{}
	svc := service.NewSnapshotService(repo, storer, "local")

	ctx := context.Background()
	venueID := uuid.New()
	createdBy := uuid.New()
	bundle := []byte(`{"data":"map"}`)

	// S3 upload succeeds
	storer.On("PutObject", ctx, mock.AnythingOfType("string"), bundle).Return(nil)
	// DB insert succeeds
	repo.On("Create", ctx, mock.AnythingOfType("*domain.Snapshot")).Return(nil)
	// Version control: count below limit
	repo.On("CountDraftsByVenue", ctx, venueID).Return(int64(1), nil)

	s, err := svc.CreateDraft(ctx, venueID, createdBy, bundle)
	require.NoError(t, err)
	assert.Equal(t, venueID, s.VenueID)
	assert.Equal(t, domain.SnapshotStateDraft, s.State)
	assert.Equal(t, domain.SnapshotMethodManual, s.Method)
	assert.Equal(t, &createdBy, s.CreatedBy)
	storer.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestSnapshotService_CreateDraft_S3Failure(t *testing.T) {
	repo := &mocks.SnapshotRepository{}
	storer := &mocks.StorerMock{}
	svc := service.NewSnapshotService(repo, storer, "local")

	ctx := context.Background()
	venueID := uuid.New()
	createdBy := uuid.New()
	bundle := []byte(`{"data":"map"}`)

	storer.On("PutObject", ctx, mock.AnythingOfType("string"), bundle).Return(assert.AnError)

	_, err := svc.CreateDraft(ctx, venueID, createdBy, bundle)
	require.Error(t, err)
	// DB must not be touched
	repo.AssertNotCalled(t, "Create")
	storer.AssertExpectations(t)
}

func TestSnapshotService_CreateDraft_PrunesOldestWhenAtLimit(t *testing.T) {
	repo := &mocks.SnapshotRepository{}
	storer := &mocks.StorerMock{}
	svc := service.NewSnapshotService(repo, storer, "local")

	ctx := context.Background()
	venueID := uuid.New()
	createdBy := uuid.New()
	bundle := []byte(`{"data":"map"}`)

	storer.On("PutObject", ctx, mock.AnythingOfType("string"), bundle).Return(nil)
	repo.On("Create", ctx, mock.AnythingOfType("*domain.Snapshot")).Return(nil)
	// Count is at MaxSnapshotVersions — prune must be called
	repo.On("CountDraftsByVenue", ctx, venueID).Return(int64(domain.MaxSnapshotVersions), nil)
	repo.On("DeleteOldestDraft", ctx, venueID).Return(nil)

	_, err := svc.CreateDraft(ctx, venueID, createdBy, bundle)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestSnapshotService_Delete(t *testing.T) {
	repo := &mocks.SnapshotRepository{}
	storer := &mocks.StorerMock{}
	svc := service.NewSnapshotService(repo, storer, "local")

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}
