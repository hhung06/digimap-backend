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

func TestLanguageService_List(t *testing.T) {
	repo := &mocks.LanguageRepository{}
	svc := service.NewLanguageService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	expected := []*domain.Language{{Code: "ja", Name: "Japanese"}, {Code: "en", Name: "English"}}

	repo.On("List", ctx, venueID).Return(expected, nil)

	langs, err := svc.List(ctx, venueID)
	require.NoError(t, err)
	assert.Len(t, langs, 2)
	repo.AssertExpectations(t)
}

func TestLanguageService_Get(t *testing.T) {
	repo := &mocks.LanguageRepository{}
	svc := service.NewLanguageService(repo)

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.Language{ID: id, Code: "ja", Name: "Japanese"}

	repo.On("FindByID", ctx, id).Return(expected, nil)

	l, err := svc.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "ja", l.Code)
	repo.AssertExpectations(t)
}

func TestLanguageService_Create(t *testing.T) {
	repo := &mocks.LanguageRepository{}
	svc := service.NewLanguageService(repo)

	ctx := context.Background()
	l := &domain.Language{Code: "ja", Name: "Japanese"}

	repo.On("Create", ctx, l).Return(nil)

	err := svc.Create(ctx, l)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLanguageService_Update(t *testing.T) {
	repo := &mocks.LanguageRepository{}
	svc := service.NewLanguageService(repo)

	ctx := context.Background()
	l := &domain.Language{ID: uuid.New(), Code: "en", Name: "English", IsDefault: true}

	repo.On("Update", ctx, l).Return(nil)

	err := svc.Update(ctx, l)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLanguageService_Delete(t *testing.T) {
	repo := &mocks.LanguageRepository{}
	svc := service.NewLanguageService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}
