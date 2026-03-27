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
	Delete(ctx context.Context, id uuid.UUID) error
	CreateImage(ctx context.Context, img *domain.ArticleImage) error
	DeleteImage(ctx context.Context, id uuid.UUID) error
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
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

func (s *articleService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *articleService) CreateImage(ctx context.Context, img *domain.ArticleImage) error {
	return s.repo.CreateImage(ctx, img)
}

func (s *articleService) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteImage(ctx, id)
}
