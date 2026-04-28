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

func TestAssetService_Create(t *testing.T) {
	repo := &mocks.AssetRepository{}
	svc := service.NewAssetService(repo)
	ctx := context.Background()

	venueID := uuid.New()
	userID := uuid.New()
	req := service.CreateAssetInput{
		VenueID:     venueID,
		Name:        "logo.png",
		Key:         "venue/logo.png",
		ContentType: "image/png",
		SizeBytes:   1024,
		URL:         "https://cdn.example.com/logo.png",
		CreatedBy:   &userID,
	}

	repo.On("Create", ctx, mock.AnythingOfType("*domain.Asset")).Return(nil)

	result, err := svc.Create(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, "logo.png", result.Name)
	assert.Equal(t, venueID, result.VenueID)
	repo.AssertExpectations(t)
}

func TestAssetService_Delete_NotFound(t *testing.T) {
	repo := &mocks.AssetRepository{}
	svc := service.NewAssetService(repo)
	ctx := context.Background()

	id := uuid.New()
	repo.On("FindByID", ctx, id).Return((*domain.Asset)(nil), domain.NewNotFound("asset not found"))

	err := svc.Delete(ctx, id)

	require.Error(t, err)
	var appErr *domain.AppError
	assert.ErrorAs(t, err, &appErr)
}
