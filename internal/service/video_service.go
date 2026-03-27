package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type VideoService interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Video, int, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Video, error)
	Create(ctx context.Context, v *domain.Video) error
	Update(ctx context.Context, v *domain.Video) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type videoService struct {
	repo repository.VideoRepository
}

func NewVideoService(repo repository.VideoRepository) VideoService {
	return &videoService{repo: repo}
}

func (s *videoService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Video, int, error) {
	return s.repo.List(ctx, venueID, p)
}

func (s *videoService) Get(ctx context.Context, id uuid.UUID) (*domain.Video, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *videoService) Create(ctx context.Context, v *domain.Video) error {
	return s.repo.Create(ctx, v)
}

func (s *videoService) Update(ctx context.Context, v *domain.Video) error {
	return s.repo.Update(ctx, v)
}

func (s *videoService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
