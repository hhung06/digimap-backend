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

func TestThemeService_Create(t *testing.T) {
	repo := &mocks.ThemeRepository{}
	svc := service.NewThemeService(repo)
	ctx := context.Background()

	venueID := uuid.New()
	repo.On("Create", ctx, mock.AnythingOfType("*domain.Theme")).Return(nil)

	result, err := svc.Create(ctx, venueID, "Dark", "#000000", "#ffffff")

	require.NoError(t, err)
	assert.Equal(t, "Dark", result.Name)
	assert.Equal(t, venueID, result.VenueID)
	assert.Equal(t, "#000000", result.PrimaryColor)
	assert.Equal(t, "#ffffff", result.SecondaryColor)
	repo.AssertExpectations(t)
}

func TestThemeService_Delete_NotFound(t *testing.T) {
	repo := &mocks.ThemeRepository{}
	svc := service.NewThemeService(repo)
	ctx := context.Background()

	id := uuid.New()
	repo.On("FindByID", ctx, id).Return((*domain.Theme)(nil), domain.NewNotFound("theme not found"))

	err := svc.Delete(ctx, id)

	require.Error(t, err)
	var appErr *domain.AppError
	assert.ErrorAs(t, err, &appErr)
}
