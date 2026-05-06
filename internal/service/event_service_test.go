package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func newTestEventService(repo *mocks.EventRepository) service.EventService {
	return service.NewEventService(repo)
}

// ── List ──────────────────────────────────────────────────────────────────────

func TestEventService_List_Success(t *testing.T) {
	repo := &mocks.EventRepository{}
	svc := newTestEventService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	p := domain.Pagination{Page: 1, PageSize: 20}
	expected := []*domain.Event{{Title: "Summer Sale"}, {Title: "Opening Gala"}}

	repo.On("List", ctx, venueID, p).Return(expected, int64(2), nil)

	got, total, err := svc.List(ctx, venueID, p)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

// ── Get ───────────────────────────────────────────────────────────────────────

func TestEventService_Get_Success(t *testing.T) {
	repo := &mocks.EventRepository{}
	svc := newTestEventService(repo)

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.Event{ID: id, Title: "Summer Sale"}

	repo.On("FindByID", ctx, id).Return(expected, nil)

	got, err := svc.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	repo.AssertExpectations(t)
}

func TestEventService_Get_NotFound(t *testing.T) {
	repo := &mocks.EventRepository{}
	svc := newTestEventService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("FindByID", ctx, id).Return((*domain.Event)(nil), domain.NewNotFound("event not found"))

	_, err := svc.Get(ctx, id)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	repo.AssertExpectations(t)
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestEventService_Create_Success(t *testing.T) {
	repo := &mocks.EventRepository{}
	svc := newTestEventService(repo)

	ctx := context.Background()
	e := &domain.Event{Title: "Summer Sale"}
	tagIDs := []uuid.UUID{uuid.New()}
	locationIDs := []uuid.UUID{uuid.New()}

	repo.On("Create", ctx, e).Return(nil)
	repo.On("SetTags", ctx, e.ID, tagIDs).Return(nil)
	repo.On("SetLocations", ctx, e.ID, locationIDs).Return(nil)

	err := svc.Create(ctx, e, tagIDs, locationIDs)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestEventService_Create_NoTagsNoLocations(t *testing.T) {
	repo := &mocks.EventRepository{}
	svc := newTestEventService(repo)

	ctx := context.Background()
	e := &domain.Event{Title: "Summer Sale"}

	repo.On("Create", ctx, e).Return(nil)

	err := svc.Create(ctx, e, nil, nil)
	require.NoError(t, err)
	repo.AssertNotCalled(t, "SetTags")
	repo.AssertNotCalled(t, "SetLocations")
	repo.AssertExpectations(t)
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestEventService_Update_Success(t *testing.T) {
	repo := &mocks.EventRepository{}
	svc := newTestEventService(repo)

	ctx := context.Background()
	e := &domain.Event{ID: uuid.New(), Title: "Updated Sale"}
	tagIDs := []uuid.UUID{uuid.New()}
	locationIDs := []uuid.UUID{}

	repo.On("Update", ctx, e).Return(nil)
	repo.On("SetTags", ctx, e.ID, tagIDs).Return(nil)
	repo.On("SetLocations", ctx, e.ID, locationIDs).Return(nil)

	err := svc.Update(ctx, e, tagIDs, locationIDs)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestEventService_Delete_Success(t *testing.T) {
	repo := &mocks.EventRepository{}
	svc := newTestEventService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

// ── Images ────────────────────────────────────────────────────────────────────

func TestEventService_CreateImage_Success(t *testing.T) {
	repo := &mocks.EventRepository{}
	svc := newTestEventService(repo)

	ctx := context.Background()
	img := &domain.EventImage{EventID: uuid.New()}

	repo.On("CreateImage", ctx, img).Return(nil)

	err := svc.CreateImage(ctx, img)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestEventService_DeleteImage_Success(t *testing.T) {
	repo := &mocks.EventRepository{}
	svc := newTestEventService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("DeleteImage", ctx, id).Return(nil)

	err := svc.DeleteImage(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}
