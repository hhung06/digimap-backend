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

func newTestConnectionService(repo *mocks.ConnectionRepository) service.ConnectionService {
	return service.NewConnectionService(repo)
}

func TestConnectionService_List_Success(t *testing.T) {
	repo := &mocks.ConnectionRepository{}
	svc := newTestConnectionService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	p := domain.Pagination{Page: 1, PageSize: 20}
	expected := []*domain.Connection{{ID: uuid.New()}, {ID: uuid.New()}}

	repo.On("List", ctx, venueID, p).Return(expected, 2, nil)

	got, total, err := svc.List(ctx, venueID, p)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, 2, total)
	repo.AssertExpectations(t)
}

func TestConnectionService_Get_Success(t *testing.T) {
	repo := &mocks.ConnectionRepository{}
	svc := newTestConnectionService(repo)

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.Connection{ID: id}

	repo.On("FindByID", ctx, id).Return(expected, nil)

	got, err := svc.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	repo.AssertExpectations(t)
}

func TestConnectionService_Get_NotFound(t *testing.T) {
	repo := &mocks.ConnectionRepository{}
	svc := newTestConnectionService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("FindByID", ctx, id).Return((*domain.Connection)(nil), domain.NewNotFound("connection not found"))

	_, err := svc.Get(ctx, id)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	repo.AssertExpectations(t)
}

func TestConnectionService_Create_Success(t *testing.T) {
	repo := &mocks.ConnectionRepository{}
	svc := newTestConnectionService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	c := &domain.Connection{VenueID: &venueID}

	repo.On("Create", ctx, c).Return(nil)

	err := svc.Create(ctx, c)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestConnectionService_Update_Success(t *testing.T) {
	repo := &mocks.ConnectionRepository{}
	svc := newTestConnectionService(repo)

	ctx := context.Background()
	c := &domain.Connection{ID: uuid.New()}

	repo.On("Update", ctx, c).Return(nil)

	err := svc.Update(ctx, c)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestConnectionService_Delete_Success(t *testing.T) {
	repo := &mocks.ConnectionRepository{}
	svc := newTestConnectionService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestConnectionService_AddLevel_Success(t *testing.T) {
	repo := &mocks.ConnectionRepository{}
	svc := newTestConnectionService(repo)

	ctx := context.Background()
	lvlID := uuid.New()
	cl := &domain.ConnectionLevel{ConnectionID: uuid.New(), LevelID: &lvlID}

	repo.On("AddLevel", ctx, cl).Return(nil)

	err := svc.AddLevel(ctx, cl)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestConnectionService_RemoveLevel_Success(t *testing.T) {
	repo := &mocks.ConnectionRepository{}
	svc := newTestConnectionService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("RemoveLevel", ctx, id).Return(nil)

	err := svc.RemoveLevel(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}
