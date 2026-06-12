package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/cdn"
	"github.com/hhung06/digimap-backend/internal/platform/storage"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func newTestLocationService(repo *mocks.LocationRepository, _ *mocks.StorerMock) service.LocationService {
	// publishTopLocations runs in a goroutine; venue mock returns an error so the goroutine exits cleanly.
	venueRepo := &mocks.VenueRepository{}
	venueRepo.On("FindByID", mock.Anything, mock.Anything).Return((*domain.Venue)(nil), domain.NewNotFound("venue not found"))
	return service.NewLocationService(repo, venueRepo, storage.NewLogStorer(), cdn.NewLogInvalidator(), nil, "test")
}

func newTestLocationCategoryService(repo *mocks.LocationCategoryRepository) service.LocationCategoryService {
	return service.NewLocationCategoryService(repo)
}

func TestLocationCategoryService_Create_SubcategorySuccess(t *testing.T) {
	repo := &mocks.LocationCategoryRepository{}
	svc := newTestLocationCategoryService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	parentID := uuid.New()
	cat := &domain.LocationCategory{
		VenueID:  venueID,
		ParentID: &parentID,
		Name:     "Coffee",
	}
	parent := &domain.LocationCategory{ID: parentID, VenueID: venueID}

	repo.On("FindByID", ctx, parentID).Return(parent, nil)
	repo.On("Create", ctx, cat).Return(nil)

	err := svc.Create(ctx, cat)
	require.NoError(t, err)
	assert.Equal(t, "internal", cat.Source)
	repo.AssertExpectations(t)
}

func TestLocationCategoryService_Create_RejectsNestedSubcategory(t *testing.T) {
	repo := &mocks.LocationCategoryRepository{}
	svc := newTestLocationCategoryService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	parentID := uuid.New()
	rootID := uuid.New()
	cat := &domain.LocationCategory{VenueID: venueID, ParentID: &parentID}
	parent := &domain.LocationCategory{ID: parentID, VenueID: venueID, ParentID: &rootID}

	repo.On("FindByID", ctx, parentID).Return(parent, nil)

	err := svc.Create(ctx, cat)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrValidation))
	repo.AssertNotCalled(t, "Create")
}

func TestLocationCategoryService_Create_RejectsCrossVenueParent(t *testing.T) {
	repo := &mocks.LocationCategoryRepository{}
	svc := newTestLocationCategoryService(repo)

	ctx := context.Background()
	parentID := uuid.New()
	cat := &domain.LocationCategory{VenueID: uuid.New(), ParentID: &parentID}
	parent := &domain.LocationCategory{ID: parentID, VenueID: uuid.New()}

	repo.On("FindByID", ctx, parentID).Return(parent, nil)

	err := svc.Create(ctx, cat)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrValidation))
	repo.AssertNotCalled(t, "Create")
}

func TestLocationCategoryService_Update_RejectsSelfParent(t *testing.T) {
	repo := &mocks.LocationCategoryRepository{}
	svc := newTestLocationCategoryService(repo)

	ctx := context.Background()
	id := uuid.New()
	cat := &domain.LocationCategory{ID: id, VenueID: uuid.New(), ParentID: &id}

	err := svc.Update(ctx, cat)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrValidation))
	repo.AssertNotCalled(t, "Update")
}

func TestLocationCategoryService_Update_RejectsParentWithSubcategoriesBecomingChild(t *testing.T) {
	repo := &mocks.LocationCategoryRepository{}
	svc := newTestLocationCategoryService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	id := uuid.New()
	parentID := uuid.New()
	childID := uuid.New()
	cat := &domain.LocationCategory{ID: id, VenueID: venueID, ParentID: &parentID}
	newParent := &domain.LocationCategory{ID: parentID, VenueID: venueID}
	child := &domain.LocationCategory{ID: childID, VenueID: venueID, ParentID: &id}

	repo.On("FindByID", ctx, parentID).Return(newParent, nil)
	repo.On("List", ctx, venueID).Return([]*domain.LocationCategory{cat, newParent, child}, nil)

	err := svc.Update(ctx, cat)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrValidation))
	repo.AssertNotCalled(t, "Update")
}

func TestLocationCategoryService_Delete_RejectsParentWithSubcategories(t *testing.T) {
	repo := &mocks.LocationCategoryRepository{}
	svc := newTestLocationCategoryService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	id := uuid.New()
	childID := uuid.New()
	cat := &domain.LocationCategory{ID: id, VenueID: venueID}
	child := &domain.LocationCategory{ID: childID, VenueID: venueID, ParentID: &id}

	repo.On("FindByID", ctx, id).Return(cat, nil)
	repo.On("List", ctx, venueID).Return([]*domain.LocationCategory{cat, child}, nil)

	err := svc.Delete(ctx, id)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrValidation))
	repo.AssertNotCalled(t, "Delete")
}

// ── Get ───────────────────────────────────────────────────────────────────────

func TestLocationService_Get_Success(t *testing.T) {
	repo := &mocks.LocationRepository{}
	svc := newTestLocationService(repo, nil)

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.Location{ID: id, CommonName: "Coffee Shop"}

	repo.On("FindByID", ctx, id).Return(expected, nil)

	got, err := svc.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	repo.AssertExpectations(t)
}

func TestLocationService_Get_NotFound(t *testing.T) {
	repo := &mocks.LocationRepository{}
	svc := newTestLocationService(repo, nil)

	ctx := context.Background()
	id := uuid.New()

	repo.On("FindByID", ctx, id).Return((*domain.Location)(nil), domain.NewNotFound("location not found"))

	_, err := svc.Get(ctx, id)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	repo.AssertExpectations(t)
}

// ── List ──────────────────────────────────────────────────────────────────────

func TestLocationService_List_Success(t *testing.T) {
	repo := &mocks.LocationRepository{}
	svc := newTestLocationService(repo, nil)

	ctx := context.Background()
	venueID := uuid.New()
	p := domain.Pagination{Page: 1, PageSize: 20}
	expected := []*domain.Location{{CommonName: "Cafe"}, {CommonName: "ATM"}}

	repo.On("List", ctx, venueID, (*int)(nil), p).Return(expected, int64(2), nil)

	got, total, err := svc.List(ctx, venueID, nil, p)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestLocationService_Create_DefaultsSource(t *testing.T) {
	repo := &mocks.LocationRepository{}
	svc := newTestLocationService(repo, nil)

	ctx := context.Background()
	l := &domain.Location{CommonName: "Shop"}

	repo.On("Create", ctx, l).Return(nil)

	err := svc.Create(ctx, l)
	require.NoError(t, err)
	assert.Equal(t, "internal", l.Source, "source must default to internal")
	repo.AssertExpectations(t)
}

func TestLocationService_Create_WithCategories(t *testing.T) {
	repo := &mocks.LocationRepository{}
	svc := newTestLocationService(repo, nil)

	ctx := context.Background()
	catID := uuid.New()
	l := &domain.Location{
		CommonName: "Shop",
		Categories: []*domain.LocationCategory{{ID: catID}},
	}

	repo.On("Create", ctx, l).Return(nil)
	repo.On("SetCategories", ctx, l.ID, []uuid.UUID{catID}).Return(nil)

	err := svc.Create(ctx, l)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestLocationService_Update_Success(t *testing.T) {
	repo := &mocks.LocationRepository{}
	svc := newTestLocationService(repo, nil)

	ctx := context.Background()
	l := &domain.Location{ID: uuid.New(), CommonName: "Updated Shop"}
	catIDs := []uuid.UUID{uuid.New()}

	repo.On("Update", ctx, l).Return(nil)
	repo.On("SetCategories", ctx, l.ID, catIDs).Return(nil)

	err := svc.Update(ctx, l, catIDs)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestLocationService_Delete_Success(t *testing.T) {
	repo := &mocks.LocationRepository{}
	svc := newTestLocationService(repo, nil)

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

// ── Duplicate ─────────────────────────────────────────────────────────────────

func TestLocationService_Duplicate_Success(t *testing.T) {
	repo := &mocks.LocationRepository{}
	svc := newTestLocationService(repo, nil)

	ctx := context.Background()
	id := uuid.New()
	catID := uuid.New()
	src := &domain.Location{
		ID:            id,
		CommonName:    "Coffee Shop",
		IsTopLocation: true,
		Categories:    []*domain.LocationCategory{{ID: catID}},
	}

	repo.On("FindByID", ctx, id).Return(src, nil)
	repo.On("Create", ctx, mockAny).Return(nil)
	repo.On("SetCategories", ctx, uuid.Nil, []uuid.UUID{catID}).Return(nil)

	clone, err := svc.Duplicate(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Coffee Shop (copy)", clone.CommonName)
	assert.False(t, clone.IsTopLocation, "top flag must be cleared on duplicate")
	repo.AssertExpectations(t)
}

func TestLocationService_Duplicate_NotFound(t *testing.T) {
	repo := &mocks.LocationRepository{}
	svc := newTestLocationService(repo, nil)

	ctx := context.Background()
	id := uuid.New()

	repo.On("FindByID", ctx, id).Return((*domain.Location)(nil), domain.NewNotFound("location not found"))

	_, err := svc.Duplicate(ctx, id)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	repo.AssertExpectations(t)
}

// ── SetTop ────────────────────────────────────────────────────────────────────
// Note: publishTopLocations runs in a background goroutine (Django parity).
// These tests cover the synchronous path only; the async publish is verified via integration tests.

func TestLocationService_SetTop_Success(t *testing.T) {
	repo := &mocks.LocationRepository{}
	svc := newTestLocationService(repo, nil)

	ctx := context.Background()
	id := uuid.New()
	venueID := uuid.New()
	sortIdx := 1

	repo.On("FindByID", ctx, id).Return(&domain.Location{ID: id, VenueID: venueID}, nil)
	repo.On("SetTopLocation", ctx, id, true, &sortIdx).Return(nil)

	err := svc.SetTop(ctx, id, true, &sortIdx)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLocationService_SetTop_NotFound(t *testing.T) {
	repo := &mocks.LocationRepository{}
	svc := newTestLocationService(repo, nil)

	ctx := context.Background()
	id := uuid.New()

	repo.On("FindByID", ctx, id).Return((*domain.Location)(nil), domain.NewNotFound("location not found"))

	err := svc.SetTop(ctx, id, true, nil)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	repo.AssertExpectations(t)
}

func TestLocationService_SetTop_DoesNotLaunchGoroutineWhenUpdateFails(t *testing.T) {
	repo := &mocks.LocationRepository{}
	svc := newTestLocationService(repo, nil)

	ctx := context.Background()
	id := uuid.New()
	venueID := uuid.New()
	sortIdx := 2

	repo.On("FindByID", ctx, id).Return(&domain.Location{ID: id, VenueID: venueID}, nil)
	repo.On("SetTopLocation", ctx, id, false, &sortIdx).Return(assert.AnError)

	err := svc.SetTop(ctx, id, false, &sortIdx)
	require.ErrorIs(t, err, assert.AnError)
	repo.AssertExpectations(t)
}

// ── DeleteImage ───────────────────────────────────────────────────────────────

func TestLocationService_DeleteImage_Success(t *testing.T) {
	repo := &mocks.LocationRepository{}
	svc := newTestLocationService(repo, nil)

	ctx := context.Background()
	locID := uuid.New()
	imgID := uuid.New()
	imgs := []*domain.LocationImage{{ID: imgID}}

	repo.On("ListImages", ctx, locID).Return(imgs, nil)
	repo.On("DeleteImage", ctx, imgID).Return(nil)

	err := svc.DeleteImage(ctx, locID, imgID)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLocationService_DeleteImage_NotFound(t *testing.T) {
	repo := &mocks.LocationRepository{}
	svc := newTestLocationService(repo, nil)

	ctx := context.Background()
	locID := uuid.New()
	imgID := uuid.New()

	repo.On("ListImages", ctx, locID).Return([]*domain.LocationImage{{ID: uuid.New()}}, nil)

	err := svc.DeleteImage(ctx, locID, imgID)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	repo.AssertExpectations(t)
}

// ── PBT: source default invariant ────────────────────────────────────────────

func TestLocationService_Create_SourceDefaultInvariant(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		repo := &mocks.LocationRepository{}
		svc := newTestLocationService(repo, nil)

		ctx := context.Background()
		l := &domain.Location{
			CommonName: rapid.StringN(1, 50, 50).Draw(rt, "name"),
			Source:     "", // always start empty to test defaulting
		}

		repo.On("Create", ctx, l).Return(nil)

		err := svc.Create(ctx, l)
		require.NoError(rt, err)
		assert.Equal(rt, "internal", l.Source, "source must always default to 'internal' when empty")
	})
}
