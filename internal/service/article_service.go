package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type ArticleService interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Article, int, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Article, error)
	Create(ctx context.Context, a *domain.Article) error
	Update(ctx context.Context, a *domain.Article) error
	UpdateWithMedia(ctx context.Context, a *domain.Article, replacement ArticleMediaReplacement) error
	Delete(ctx context.Context, id uuid.UUID) error
	CreateImage(ctx context.Context, img *domain.ArticleImage) error
	DeleteImage(ctx context.Context, id uuid.UUID) error
}

type ArticleMediaReplacement struct {
	Replace bool
	Remove  bool
	Uploads []MediaUpload
}

type articleService struct {
	repo     repository.ArticleRepository
	mediaSvc MediaService
}

func NewArticleService(repo repository.ArticleRepository, mediaSvc ...MediaService) ArticleService {
	var media MediaService
	if len(mediaSvc) > 0 {
		media = mediaSvc[0]
	}
	return &articleService{repo: repo, mediaSvc: media}
}

func (s *articleService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Article, int, error) {
	return s.repo.List(ctx, venueID, p)
}

func (s *articleService) Get(ctx context.Context, id uuid.UUID) (*domain.Article, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *articleService) Create(ctx context.Context, a *domain.Article) error {
	return s.repo.Create(ctx, a)
}

func (s *articleService) Update(ctx context.Context, a *domain.Article) error {
	return s.repo.Update(ctx, a)
}

func (s *articleService) UpdateWithMedia(ctx context.Context, a *domain.Article, replacement ArticleMediaReplacement) error {
	if !replacement.Replace {
		return s.repo.Update(ctx, a)
	}
	target := articleImagesMediaTarget(a.ID)
	oldKeys := articleImageKeys(a.Images)
	newKeys := make([]string, 0, len(replacement.Uploads))
	if len(replacement.Uploads) > 0 && s.mediaSvc == nil {
		return &domain.AppError{Err: domain.ErrInternal, Message: "media service is not configured"}
	}
	for _, upload := range replacement.Uploads {
		key, err := s.mediaSvc.Upload(ctx, target, upload)
		if err != nil {
			deleteArticleMediaKeys(ctx, s.mediaSvc, target, newKeys)
			return err
		}
		newKeys = append(newKeys, key)
	}

	if err := s.repo.UpdateWithImages(ctx, a, domain.ArticleMediaChange{Replace: true, Keys: newKeys}); err != nil {
		deleteArticleMediaKeys(ctx, s.mediaSvc, target, newKeys)
		return err
	}

	deleteArticleMediaKeys(ctx, s.mediaSvc, target, oldKeys)
	return nil
}

func (s *articleService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *articleService) CreateImage(ctx context.Context, img *domain.ArticleImage) error {
	return s.repo.CreateImage(ctx, img)
}

func (s *articleService) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteImage(ctx, id)
}

func articleImagesMediaTarget(articleID uuid.UUID) MediaTarget {
	return MediaTarget{Entity: "articles", RecordID: articleID, Field: "images"}
}

func articleImageKeys(images []*domain.ArticleImage) []string {
	keys := make([]string, 0, len(images))
	for _, img := range images {
		if img != nil {
			keys = append(keys, img.Image)
		}
	}
	return keys
}

func deleteArticleMediaKeys(ctx context.Context, mediaSvc MediaService, target MediaTarget, keys []string) {
	if mediaSvc == nil {
		return
	}
	for _, key := range keys {
		_ = mediaSvc.DeleteOwned(ctx, target, key)
	}
}
