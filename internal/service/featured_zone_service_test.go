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

func TestFeaturedZoneService_List(t *testing.T) {
	repo := &mocks.FeaturedZoneRepository{}
	svc := service.NewFeaturedZoneService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	expected := []*domain.FeaturedZone{{Name: "Zone A", IsActive: true}, {Name: "Zone B", IsActive: false}}

	repo.On("List", ctx, venueID).Return(expected, nil)

	zones, err := svc.List(ctx, venueID)
	require.NoError(t, err)
	assert.Len(t, zones, 2)
	repo.AssertExpectations(t)
}

func TestFeaturedZoneService_ListActive(t *testing.T) {
	repo := &mocks.FeaturedZoneRepository{}
	svc := service.NewFeaturedZoneService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	expected := []*domain.FeaturedZone{{Name: "Zone A", IsActive: true}}

	repo.On("ListActive", ctx, venueID).Return(expected, nil)

	zones, err := svc.ListActive(ctx, venueID)
	require.NoError(t, err)
	assert.Len(t, zones, 1)
	assert.True(t, zones[0].IsActive)
	repo.AssertExpectations(t)
}

func TestFeaturedZoneService_Get(t *testing.T) {
	repo := &mocks.FeaturedZoneRepository{}
	svc := service.NewFeaturedZoneService(repo)

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.FeaturedZone{ID: id, Name: "Zone A"}

	repo.On("FindByID", ctx, id).Return(expected, nil)

	z, err := svc.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Zone A", z.Name)
	repo.AssertExpectations(t)
}

func TestFeaturedZoneService_Create(t *testing.T) {
	repo := &mocks.FeaturedZoneRepository{}
	svc := service.NewFeaturedZoneService(repo)

	ctx := context.Background()
	z := &domain.FeaturedZone{Name: "New Zone", IsActive: true}

	repo.On("Create", ctx, z).Return(nil)

	err := svc.Create(ctx, z)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestFeaturedZoneService_Delete(t *testing.T) {
	repo := &mocks.FeaturedZoneRepository{}
	svc := service.NewFeaturedZoneService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}
