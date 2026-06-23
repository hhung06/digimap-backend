package service_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/storage"
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

type articleMediaServiceSpy struct {
	uploadKeys []string
	uploadErr  error
	deleted    []articleMediaDeleteCall
}

type articleMediaDeleteCall struct {
	target service.MediaTarget
	key    string
}

func (s *articleMediaServiceSpy) Upload(ctx context.Context, target service.MediaTarget, upload service.MediaUpload) (string, error) {
	if s.uploadErr != nil {
		return "", s.uploadErr
	}
	if upload.Reader != nil {
		_, _ = io.ReadAll(upload.Reader)
	}
	if len(s.uploadKeys) == 0 {
		return storage.MediaKey("develop", target.Entity, target.RecordID, target.Field, uuid.New(), "png"), nil
	}
	key := s.uploadKeys[0]
	s.uploadKeys = s.uploadKeys[1:]
	return key, nil
}

func (s *articleMediaServiceSpy) URL(context.Context, service.MediaTarget, string) *string {
	return nil
}

func (s *articleMediaServiceSpy) DeleteOwned(_ context.Context, target service.MediaTarget, key string) error {
	s.deleted = append(s.deleted, articleMediaDeleteCall{target: target, key: key})
	return nil
}

func TestArticleService_UpdateWithMediaUploadFailureDoesNotMutateDB(t *testing.T) {
	repo := &mocks.ArticleRepository{}
	media := &articleMediaServiceSpy{uploadErr: assert.AnError}
	svc := service.NewArticleService(repo, media)
	ctx := context.Background()
	a := &domain.Article{ID: uuid.New(), Title: "Updated"}

	err := svc.UpdateWithMedia(ctx, a, service.ArticleMediaReplacement{
		Replace: true,
		Uploads: []service.MediaUpload{{
			Filename: "hero.png", ContentType: "image/png", Size: 3, Reader: bytes.NewReader([]byte("png")),
		}},
	})

	require.Error(t, err)
	repo.AssertNotCalled(t, "UpdateWithImages")
	repo.AssertExpectations(t)
	assert.Empty(t, media.deleted)
}

func TestArticleService_UpdateWithMediaDBFailureDeletesUploadedKeys(t *testing.T) {
	repo := &mocks.ArticleRepository{}
	media := &articleMediaServiceSpy{}
	svc := service.NewArticleService(repo, media)
	ctx := context.Background()
	articleID := uuid.New()
	target := service.MediaTarget{Entity: "articles", RecordID: articleID, Field: "images"}
	newKey := storage.MediaKey("develop", target.Entity, target.RecordID, target.Field, uuid.New(), "png")
	a := &domain.Article{ID: articleID, Title: "Updated"}

	repo.On("UpdateWithImages", ctx, a, domain.ArticleMediaChange{Replace: true, Keys: []string{newKey}}).Return(assert.AnError)
	media.uploadKeys = []string{newKey}

	err := svc.UpdateWithMedia(ctx, a, service.ArticleMediaReplacement{
		Replace: true,
		Uploads: []service.MediaUpload{{
			Filename: "hero.png", ContentType: "image/png", Size: 3, Reader: bytes.NewReader([]byte("png")),
		}},
	})

	require.ErrorIs(t, err, assert.AnError)
	require.Len(t, media.deleted, 1)
	assert.Equal(t, target, media.deleted[0].target)
	assert.Equal(t, newKey, media.deleted[0].key)
	repo.AssertExpectations(t)
}

func TestArticleService_UpdateWithMediaSuccessDeletesOldKeysAfterCommit(t *testing.T) {
	repo := &mocks.ArticleRepository{}
	media := &articleMediaServiceSpy{}
	svc := service.NewArticleService(repo, media)
	ctx := context.Background()
	articleID := uuid.New()
	target := service.MediaTarget{Entity: "articles", RecordID: articleID, Field: "images"}
	oldKey := storage.MediaKey("develop", target.Entity, target.RecordID, target.Field, uuid.New(), "png")
	newKey := storage.MediaKey("develop", target.Entity, target.RecordID, target.Field, uuid.New(), "png")
	a := &domain.Article{
		ID:     articleID,
		Title:  "Updated",
		Images: []*domain.ArticleImage{{ArticleID: articleID, Image: oldKey}},
	}

	repo.On("UpdateWithImages", ctx, a, domain.ArticleMediaChange{Replace: true, Keys: []string{newKey}}).Return(nil)
	media.uploadKeys = []string{newKey}

	err := svc.UpdateWithMedia(ctx, a, service.ArticleMediaReplacement{
		Replace: true,
		Uploads: []service.MediaUpload{{
			Filename: "hero.png", ContentType: "image/png", Size: 3, Reader: bytes.NewReader([]byte("png")),
		}},
	})

	require.NoError(t, err)
	require.Len(t, media.deleted, 1)
	assert.Equal(t, target, media.deleted[0].target)
	assert.Equal(t, oldKey, media.deleted[0].key)
	repo.AssertExpectations(t)
}

func TestArticleService_UpdateWithMediaSameNameUploadsUseMediaServiceUniqueKeys(t *testing.T) {
	repo := &mocks.ArticleRepository{}
	media := &articleMediaServiceSpy{}
	svc := service.NewArticleService(repo, media)
	ctx := context.Background()
	articleID := uuid.New()
	target := service.MediaTarget{Entity: "articles", RecordID: articleID, Field: "images"}
	firstKey := storage.MediaKey("develop", target.Entity, target.RecordID, target.Field, uuid.New(), "png")
	secondKey := storage.MediaKey("develop", target.Entity, target.RecordID, target.Field, uuid.New(), "png")
	a := &domain.Article{ID: articleID, Title: "Updated"}

	repo.On("UpdateWithImages", ctx, a, domain.ArticleMediaChange{Replace: true, Keys: []string{firstKey, secondKey}}).Return(nil)
	media.uploadKeys = []string{firstKey, secondKey}

	err := svc.UpdateWithMedia(ctx, a, service.ArticleMediaReplacement{
		Replace: true,
		Uploads: []service.MediaUpload{
			{Filename: "hero.png", ContentType: "image/png", Size: 3, Reader: bytes.NewReader([]byte("one"))},
			{Filename: "hero.png", ContentType: "image/png", Size: 3, Reader: bytes.NewReader([]byte("two"))},
		},
	})

	require.NoError(t, err)
	assert.NotEqual(t, firstKey, secondKey)
	repo.AssertExpectations(t)
}

func TestArticleService_UpdateWithMediaOldKeyDeletionUsesArticleImagesTarget(t *testing.T) {
	repo := &mocks.ArticleRepository{}
	media := &articleMediaServiceSpy{}
	svc := service.NewArticleService(repo, media)
	ctx := context.Background()
	articleID := uuid.New()
	foreignTarget := service.MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "images"}
	foreignKey := storage.MediaKey("develop", foreignTarget.Entity, foreignTarget.RecordID, foreignTarget.Field, uuid.New(), "png")
	newTarget := service.MediaTarget{Entity: "articles", RecordID: articleID, Field: "images"}
	newKey := storage.MediaKey("develop", newTarget.Entity, newTarget.RecordID, newTarget.Field, uuid.New(), "png")
	a := &domain.Article{
		ID:     articleID,
		Title:  "Updated",
		Images: []*domain.ArticleImage{{ArticleID: articleID, Image: foreignKey}},
	}

	repo.On("UpdateWithImages", ctx, a, domain.ArticleMediaChange{Replace: true, Keys: []string{newKey}}).Return(nil)
	media.uploadKeys = []string{newKey}

	err := svc.UpdateWithMedia(ctx, a, service.ArticleMediaReplacement{
		Replace: true,
		Uploads: []service.MediaUpload{{Filename: "hero.png", ContentType: "image/png", Size: 3, Reader: bytes.NewReader([]byte("png"))}},
	})

	require.NoError(t, err)
	require.Len(t, media.deleted, 1)
	assert.Equal(t, service.MediaTarget{Entity: "articles", RecordID: articleID, Field: "images"}, media.deleted[0].target)
	assert.Equal(t, foreignKey, media.deleted[0].key)
	repo.AssertExpectations(t)
}
