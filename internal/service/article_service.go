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
	CreateWithMedia(ctx context.Context, a *domain.Article, uploads []MediaUpload) error
	Update(ctx context.Context, a *domain.Article) error
	UpdateWithMedia(ctx context.Context, a *domain.Article, replacement ArticleMediaReplacement) error
	Delete(ctx context.Context, id uuid.UUID) error
	CreateImage(ctx context.Context, img *domain.ArticleImage) error
	DeleteImage(ctx context.Context, id uuid.UUID) error
}

type ArticleMediaReplacement struct {
	Replace      bool
	Remove       bool
	KeepImageIDs []uuid.UUID
	Uploads      []MediaUpload
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

func (s *articleService) CreateWithMedia(ctx context.Context, a *domain.Article, uploads []MediaUpload) error {
	if err := s.repo.Create(ctx, a); err != nil {
		return err
	}
	if len(uploads) == 0 {
		return nil
	}
	target := articleImagesMediaTarget(a.ID)
	newKeys, err := s.uploadArticleMedia(ctx, target, uploads)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateWithImages(ctx, a, domain.ArticleMediaChange{Replace: true, Keys: newKeys}); err != nil {
		deleteArticleMediaKeys(ctx, s.mediaSvc, target, newKeys)
		return err
	}
	return nil
}

func (s *articleService) Update(ctx context.Context, a *domain.Article) error {
	return s.repo.Update(ctx, a)
}

func (s *articleService) UpdateWithMedia(ctx context.Context, a *domain.Article, replacement ArticleMediaReplacement) error {
	if !replacement.Replace {
		return s.repo.Update(ctx, a)
	}
	target := articleImagesMediaTarget(a.ID)
	keptKeys, dropKeys := partitionArticleImageKeys(a.Images, replacement.KeepImageIDs, replacement.Remove)
	newKeys, err := s.uploadArticleMedia(ctx, target, replacement.Uploads)
	if err != nil {
		return err
	}
	finalKeys := append(keptKeys, newKeys...)

	if err := s.repo.UpdateWithImages(ctx, a, domain.ArticleMediaChange{Replace: true, Keys: finalKeys}); err != nil {
		deleteArticleMediaKeys(ctx, s.mediaSvc, target, newKeys)
		return err
	}

	deleteArticleMediaKeys(ctx, s.mediaSvc, target, dropKeys)
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

func partitionArticleImageKeys(images []*domain.ArticleImage, keepIDs []uuid.UUID, remove bool) ([]string, []string) {
	if remove || len(keepIDs) == 0 {
		return nil, articleImageKeys(images)
	}
	keep := make(map[uuid.UUID]struct{}, len(keepIDs))
	for _, id := range keepIDs {
		keep[id] = struct{}{}
	}
	keptKeys := make([]string, 0, len(images))
	dropKeys := make([]string, 0, len(images))
	for _, img := range images {
		if img == nil {
			continue
		}
		if _, ok := keep[img.ID]; ok {
			keptKeys = append(keptKeys, img.Image)
			continue
		}
		dropKeys = append(dropKeys, img.Image)
	}
	return keptKeys, dropKeys
}

func (s *articleService) uploadArticleMedia(ctx context.Context, target MediaTarget, uploads []MediaUpload) ([]string, error) {
	newKeys := make([]string, 0, len(uploads))
	if len(uploads) > 0 && s.mediaSvc == nil {
		return nil, &domain.AppError{Err: domain.ErrInternal, Message: "media service is not configured"}
	}
	for _, upload := range uploads {
		key, err := s.mediaSvc.Upload(ctx, target, upload)
		if err != nil {
			deleteArticleMediaKeys(ctx, s.mediaSvc, target, newKeys)
			return nil, err
		}
		newKeys = append(newKeys, key)
	}
	return newKeys, nil
}

func deleteArticleMediaKeys(ctx context.Context, mediaSvc MediaService, target MediaTarget, keys []string) {
	if mediaSvc == nil {
		return
	}
	for _, key := range keys {
		_ = mediaSvc.DeleteOwned(ctx, target, key)
	}
}
