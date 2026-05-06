package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func newTestTagService(repo *mocks.TagRepository) service.TagService {
	return service.NewTagService(repo)
}

func TestTagService_List_Success(t *testing.T) {
	repo := &mocks.TagRepository{}
	svc := newTestTagService(repo)

	ctx := context.Background()
	p := domain.Pagination{Page: 1, PageSize: 20}
	expected := []*domain.Tag{{Name: "sale"}, {Name: "new"}}

	repo.On("List", ctx, p).Return(expected, 2, nil)

	got, total, err := svc.List(ctx, p)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, 2, total)
	repo.AssertExpectations(t)
}

func TestTagService_Get_Success(t *testing.T) {
	repo := &mocks.TagRepository{}
	svc := newTestTagService(repo)

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.Tag{ID: id, Name: "sale"}

	repo.On("FindByID", ctx, id).Return(expected, nil)

	got, err := svc.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	repo.AssertExpectations(t)
}

func TestTagService_Get_NotFound(t *testing.T) {
	repo := &mocks.TagRepository{}
	svc := newTestTagService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("FindByID", ctx, id).Return((*domain.Tag)(nil), domain.NewNotFound("tag not found"))

	_, err := svc.Get(ctx, id)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	repo.AssertExpectations(t)
}

func TestTagService_Create_Success(t *testing.T) {
	repo := &mocks.TagRepository{}
	svc := newTestTagService(repo)

	ctx := context.Background()
	tag := &domain.Tag{Name: "promo"}

	repo.On("Create", ctx, tag).Return(nil)

	err := svc.Create(ctx, tag)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestTagService_Update_Success(t *testing.T) {
	repo := &mocks.TagRepository{}
	svc := newTestTagService(repo)

	ctx := context.Background()
	tag := &domain.Tag{ID: uuid.New(), Name: "updated-tag"}

	repo.On("Update", ctx, tag).Return(nil)

	err := svc.Update(ctx, tag)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestTagService_Delete_Success(t *testing.T) {
	repo := &mocks.TagRepository{}
	svc := newTestTagService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestTagService_AttachTag_Success(t *testing.T) {
	repo := &mocks.TagRepository{}
	svc := newTestTagService(repo)

	ctx := context.Background()
	et := &domain.EntityTag{TagID: uuid.New(), EntityType: "location", EntityID: uuid.New()}

	repo.On("AttachTag", ctx, et).Return(nil)

	err := svc.AttachTag(ctx, et)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestTagService_DetachTag_Success(t *testing.T) {
	repo := &mocks.TagRepository{}
	svc := newTestTagService(repo)

	ctx := context.Background()
	tagID := uuid.New()
	entityID := uuid.New()

	repo.On("DetachTag", ctx, tagID, "location", entityID).Return(nil)

	err := svc.DetachTag(ctx, tagID, "location", entityID)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}
