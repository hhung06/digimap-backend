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

func TestLevelTypeService_Create(t *testing.T) {
	repo := &mocks.LevelTypeRepository{}
	svc := service.NewLevelTypeService(repo)
	ctx := context.Background()

	venueID := uuid.New()
	repo.On("Create", ctx, mock.AnythingOfType("*domain.LevelType")).Return(nil)

	result, err := svc.Create(ctx, venueID, "Floor", "floor-icon")

	require.NoError(t, err)
	assert.Equal(t, "Floor", result.Name)
	assert.Equal(t, venueID, result.VenueID)
	repo.AssertExpectations(t)
}

func TestLevelTypeService_Delete_NotFound(t *testing.T) {
	repo := &mocks.LevelTypeRepository{}
	svc := service.NewLevelTypeService(repo)
	ctx := context.Background()

	id := uuid.New()
	repo.On("FindByID", ctx, id).Return((*domain.LevelType)(nil), domain.NewNotFound("level type not found"))

	err := svc.Delete(ctx, id)

	require.Error(t, err)
	var appErr *domain.AppError
	assert.ErrorAs(t, err, &appErr)
}
