package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func TestBeaconService_List(t *testing.T) {
	repo := &mocks.BeaconRepository{}
	svc := service.NewBeaconService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	p := domain.Pagination{Page: 1, PageSize: 20}
	expected := []*domain.Beacon{{Name: "Beacon A"}, {Name: "Beacon B"}}

	repo.On("List", ctx, venueID, p).Return(expected, 2, nil)

	beacons, total, err := svc.List(ctx, venueID, p)
	require.NoError(t, err)
	assert.Len(t, beacons, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

func TestBeaconService_Get_Found(t *testing.T) {
	repo := &mocks.BeaconRepository{}
	svc := service.NewBeaconService(repo)

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.Beacon{ID: id, Name: "Beacon X"}

	repo.On("FindByID", ctx, id).Return(expected, nil)

	b, err := svc.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, b.ID)
	repo.AssertExpectations(t)
}

func TestBeaconService_Get_NotFound(t *testing.T) {
	repo := &mocks.BeaconRepository{}
	svc := service.NewBeaconService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("FindByID", ctx, id).Return((*domain.Beacon)(nil), domain.NewNotFound("beacon not found"))

	_, err := svc.Get(ctx, id)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestBeaconService_Create(t *testing.T) {
	repo := &mocks.BeaconRepository{}
	svc := service.NewBeaconService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	b := &domain.Beacon{VenueID: &venueID, Name: "New Beacon", HwID: "AA:BB:CC:DD"}

	repo.On("Create", ctx, b).Return(nil)

	err := svc.Create(ctx, b)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBeaconService_Update(t *testing.T) {
	repo := &mocks.BeaconRepository{}
	svc := service.NewBeaconService(repo)

	ctx := context.Background()
	b := &domain.Beacon{ID: uuid.New(), Name: "Updated"}

	repo.On("Update", ctx, b).Return(nil)

	err := svc.Update(ctx, b)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBeaconService_Delete(t *testing.T) {
	repo := &mocks.BeaconRepository{}
	svc := service.NewBeaconService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}
