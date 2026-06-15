package service_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/cdn"
	"github.com/hhung06/digimap-backend/internal/platform/storage"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func newTestThemeService(repo *mocks.ThemeRepository, venueRepo *mocks.VenueRepository) service.ThemeService {
	return service.NewThemeService(repo, venueRepo, storage.NewLogStorer(), cdn.NewLogInvalidator(), "test")
}

func TestThemeService_Create(t *testing.T) {
	repo := &mocks.ThemeRepository{}
	venueRepo := &mocks.VenueRepository{}
	svc := newTestThemeService(repo, venueRepo)
	ctx := context.Background()

	venueID := uuid.New()
	data := json.RawMessage(`{"primary_color":"#000"}`)
	repo.On("Create", ctx, mock.AnythingOfType("*domain.Theme")).Return(nil)

	result, err := svc.Create(ctx, venueID, "Dark", data)

	require.NoError(t, err)
	assert.Equal(t, "Dark", result.Name)
	require.NotNil(t, result.VenueID)
	assert.Equal(t, venueID, *result.VenueID)
	assert.Equal(t, domain.ThemeScopeCustom, result.Scope)
	repo.AssertExpectations(t)
}

func TestThemeService_CreateGlobal_SetsStoragePath(t *testing.T) {
	repo := &mocks.ThemeRepository{}
	venueRepo := &mocks.VenueRepository{}
	svc := newTestThemeService(repo, venueRepo)
	ctx := context.Background()

	data := json.RawMessage(`{"color":"#fff"}`)
	repo.On("Create", ctx, mock.AnythingOfType("*domain.Theme")).Return(nil)

	result, err := svc.CreateGlobal(ctx, "Light", data)

	require.NoError(t, err)
	assert.Equal(t, domain.ThemeScopeGlobal, result.Scope)
	assert.Equal(t, "test/venue_themes/global/Light.json", result.StoragePath)
	repo.AssertExpectations(t)
}

func TestThemeService_Delete_NotFound(t *testing.T) {
	repo := &mocks.ThemeRepository{}
	venueRepo := &mocks.VenueRepository{}
	svc := newTestThemeService(repo, venueRepo)
	ctx := context.Background()

	id := uuid.New()
	repo.On("FindByID", ctx, id).Return((*domain.Theme)(nil), domain.NewNotFound("theme not found"))

	err := svc.Delete(ctx, id)

	require.Error(t, err)
	var appErr *domain.AppError
	assert.ErrorAs(t, err, &appErr)
}

func TestThemeService_Delete_InUse(t *testing.T) {
	repo := &mocks.ThemeRepository{}
	venueRepo := &mocks.VenueRepository{}
	svc := newTestThemeService(repo, venueRepo)
	ctx := context.Background()

	id := uuid.New()
	repo.On("FindByID", ctx, id).Return(&domain.Theme{ID: id, Scope: domain.ThemeScopeGlobal}, nil)
	repo.On("IsUsedByVenues", ctx, id).Return(true, nil)

	err := svc.Delete(ctx, id)

	require.Error(t, err)
	var appErr *domain.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.ErrorIs(t, appErr.Err, domain.ErrConflict)
}

func TestThemeService_SetVenueTheme_CustomWrongVenue(t *testing.T) {
	repo := &mocks.ThemeRepository{}
	venueRepo := &mocks.VenueRepository{}
	svc := newTestThemeService(repo, venueRepo)
	ctx := context.Background()

	venueID := uuid.New()
	otherVenueID := uuid.New()
	themeID := uuid.New()
	repo.On("FindByID", ctx, themeID).Return(&domain.Theme{
		ID:      themeID,
		Scope:   domain.ThemeScopeCustom,
		VenueID: &otherVenueID,
	}, nil)

	err := svc.SetVenueTheme(ctx, venueID, &themeID)

	require.Error(t, err)
	var appErr *domain.AppError
	assert.ErrorAs(t, err, &appErr)
}
