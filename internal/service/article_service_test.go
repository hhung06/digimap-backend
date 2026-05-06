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

func newTestArticleService(repo *mocks.ArticleRepository) service.ArticleService {
	return service.NewArticleService(repo)
}

func TestArticleService_List_Success(t *testing.T) {
	repo := &mocks.ArticleRepository{}
	svc := newTestArticleService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	p := domain.Pagination{Page: 1, PageSize: 20}
	expected := []*domain.Article{{Title: "News A"}, {Title: "News B"}}

	repo.On("List", ctx, venueID, p).Return(expected, 2, nil)

	got, total, err := svc.List(ctx, venueID, p)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, 2, total)
	repo.AssertExpectations(t)
}

func TestArticleService_Get_Success(t *testing.T) {
	repo := &mocks.ArticleRepository{}
	svc := newTestArticleService(repo)

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.Article{ID: id, Title: "News A"}

	repo.On("FindByID", ctx, id).Return(expected, nil)

	got, err := svc.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	repo.AssertExpectations(t)
}

func TestArticleService_Get_NotFound(t *testing.T) {
	repo := &mocks.ArticleRepository{}
	svc := newTestArticleService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("FindByID", ctx, id).Return((*domain.Article)(nil), domain.NewNotFound("article not found"))

	_, err := svc.Get(ctx, id)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	repo.AssertExpectations(t)
}

func TestArticleService_Create_Success(t *testing.T) {
	repo := &mocks.ArticleRepository{}
	svc := newTestArticleService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	a := &domain.Article{Title: "New Article", VenueID: &venueID}

	repo.On("Create", ctx, a).Return(nil)

	err := svc.Create(ctx, a)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestArticleService_Update_Success(t *testing.T) {
	repo := &mocks.ArticleRepository{}
	svc := newTestArticleService(repo)

	ctx := context.Background()
	a := &domain.Article{ID: uuid.New(), Title: "Updated Article"}

	repo.On("Update", ctx, a).Return(nil)

	err := svc.Update(ctx, a)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestArticleService_Delete_Success(t *testing.T) {
	repo := &mocks.ArticleRepository{}
	svc := newTestArticleService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestArticleService_CreateImage_Success(t *testing.T) {
	repo := &mocks.ArticleRepository{}
	svc := newTestArticleService(repo)

	ctx := context.Background()
	img := &domain.ArticleImage{ArticleID: uuid.New()}

	repo.On("CreateImage", ctx, img).Return(nil)

	err := svc.CreateImage(ctx, img)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestArticleService_DeleteImage_Success(t *testing.T) {
	repo := &mocks.ArticleRepository{}
	svc := newTestArticleService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("DeleteImage", ctx, id).Return(nil)

	err := svc.DeleteImage(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}
