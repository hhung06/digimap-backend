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

func TestProductPlazaService_Create(t *testing.T) {
	repo := &mocks.ProductPlazaRepository{}
	svc := service.NewProductPlazaService(repo)
	ctx := context.Background()

	venueID := uuid.New()
	repo.On("Create", ctx, mock.AnythingOfType("*domain.ProductPlaza")).Return(nil)

	result, err := svc.Create(ctx, service.CreateProductPlazaInput{
		VenueID:     venueID,
		Name:        "Main Plaza",
		Description: "Primary shopping zone",
	})

	require.NoError(t, err)
	assert.Equal(t, "Main Plaza", result.Name)
	assert.Equal(t, venueID, result.VenueID)
	assert.Equal(t, "Primary shopping zone", result.Description)
	repo.AssertExpectations(t)
}

func TestProductPlazaService_Delete_NotFound(t *testing.T) {
	repo := &mocks.ProductPlazaRepository{}
	svc := service.NewProductPlazaService(repo)
	ctx := context.Background()

	id := uuid.New()
	repo.On("FindByID", ctx, id).Return((*domain.ProductPlaza)(nil), domain.NewNotFound("product plaza not found"))

	err := svc.Delete(ctx, id)

	require.Error(t, err)
	var appErr *domain.AppError
	assert.ErrorAs(t, err, &appErr)
}
