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
	l := &domain.Language{Code: "ja", Name: "Japanese", Enabled: true}

	repo.On("Create", ctx, l).Return(nil)

	err := svc.Create(ctx, l)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLanguageService_Create_SetsDefaultAfterInsert(t *testing.T) {
	repo := &mocks.LanguageRepository{}
	svc := service.NewLanguageService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	l := &domain.Language{VenueID: venueID, Code: "en", Name: "English", IsDefault: true, Enabled: true}

	repo.On("Create", ctx, l).Return(nil)
	repo.On("SetDefault", ctx, venueID, l.ID).Return(nil)

	err := svc.Create(ctx, l)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLanguageService_Create_RejectsDisabledDefault(t *testing.T) {
	repo := &mocks.LanguageRepository{}
	svc := service.NewLanguageService(repo)

	l := &domain.Language{Code: "en", Name: "English", IsDefault: true, Enabled: false}

	err := svc.Create(context.Background(), l)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
	repo.AssertNotCalled(t, "Create")
}

func TestLanguageService_Update(t *testing.T) {
	repo := &mocks.LanguageRepository{}
	svc := service.NewLanguageService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	l := &domain.Language{ID: uuid.New(), VenueID: venueID, Code: "en", Name: "English", IsDefault: true, Enabled: true}

	repo.On("Update", ctx, l).Return(nil)
	repo.On("SetDefault", ctx, venueID, l.ID).Return(nil)

	err := svc.Update(ctx, l)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLanguageService_Update_RejectsDisablingDefault(t *testing.T) {
	repo := &mocks.LanguageRepository{}
	svc := service.NewLanguageService(repo)

	l := &domain.Language{ID: uuid.New(), Code: "en", Name: "English", IsDefault: true, Enabled: false}

	err := svc.Update(context.Background(), l)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
	repo.AssertNotCalled(t, "Update")
}

func TestLanguageService_Update_NonDefaultDisableAllowed(t *testing.T) {
	repo := &mocks.LanguageRepository{}
	svc := service.NewLanguageService(repo)

	ctx := context.Background()
	l := &domain.Language{ID: uuid.New(), Code: "ja", Name: "Japanese", IsDefault: false, Enabled: false}

	repo.On("Update", ctx, l).Return(nil)

	err := svc.Update(ctx, l)
	require.NoError(t, err)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "SetDefault", mock.Anything, mock.Anything, mock.Anything)
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

func TestLanguageService_Reconcile_RejectsEmptyList(t *testing.T) {
	repo := &mocks.LanguageRepository{}
	svc := service.NewLanguageService(repo)

	_, err := svc.Reconcile(context.Background(), uuid.New(), []*domain.Language{})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
	repo.AssertNotCalled(t, "ReplaceAll")
}

func TestLanguageService_Reconcile_RejectsDuplicateCodes(t *testing.T) {
	repo := &mocks.LanguageRepository{}
	svc := service.NewLanguageService(repo)
	items := []*domain.Language{{Code: "en", Name: "English"}, {Code: "en", Name: "English (dup)"}}

	_, err := svc.Reconcile(context.Background(), uuid.New(), items)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
	repo.AssertNotCalled(t, "ReplaceAll")
}

func TestLanguageService_Reconcile_RejectsMultipleDefaults(t *testing.T) {
	repo := &mocks.LanguageRepository{}
	svc := service.NewLanguageService(repo)
	items := []*domain.Language{
		{Code: "en", Name: "English", IsDefault: true},
		{Code: "ja", Name: "Japanese", IsDefault: true},
	}

	_, err := svc.Reconcile(context.Background(), uuid.New(), items)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrValidation)
	repo.AssertNotCalled(t, "ReplaceAll")
}

func TestLanguageService_Reconcile_DefaultsToFirstItemWhenNoneMarked(t *testing.T) {
	repo := &mocks.LanguageRepository{}
	svc := service.NewLanguageService(repo)
	venueID := uuid.New()
	items := []*domain.Language{{Code: "en", Name: "English"}, {Code: "ja", Name: "Japanese"}}

	repo.On("ReplaceAll", mock.Anything, venueID, mock.MatchedBy(func(got []*domain.Language) bool {
		return len(got) == 2 && got[0].IsDefault && !got[1].IsDefault
	})).Return(items, nil)

	_, err := svc.Reconcile(context.Background(), venueID, items)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLanguageService_Reconcile_StampsVenueIDOnEachItem(t *testing.T) {
	repo := &mocks.LanguageRepository{}
	svc := service.NewLanguageService(repo)
	venueID := uuid.New()
	items := []*domain.Language{{Code: "en", Name: "English", IsDefault: true}}

	repo.On("ReplaceAll", mock.Anything, venueID, mock.MatchedBy(func(got []*domain.Language) bool {
		return got[0].VenueID == venueID
	})).Return(items, nil)

	_, err := svc.Reconcile(context.Background(), venueID, items)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}
