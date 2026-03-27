package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type AdvertisementService interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Advertisement, int, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Advertisement, error)
	Create(ctx context.Context, a *domain.Advertisement) error
	Update(ctx context.Context, a *domain.Advertisement) error
	Delete(ctx context.Context, id uuid.UUID) error
	Publish(ctx context.Context, id uuid.UUID) error
}

type adService struct {
	repo repository.AdvertisementRepository
}

func NewAdvertisementService(repo repository.AdvertisementRepository) AdvertisementService {
	return &adService{repo: repo}
}

func (s *adService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Advertisement, int, error) {
	return s.repo.List(ctx, venueID, p)
}

func (s *adService) Get(ctx context.Context, id uuid.UUID) (*domain.Advertisement, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *adService) Create(ctx context.Context, a *domain.Advertisement) error {
	return s.repo.Create(ctx, a)
}

func (s *adService) Update(ctx context.Context, a *domain.Advertisement) error {
	return s.repo.Update(ctx, a)
}

func (s *adService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *adService) Publish(ctx context.Context, id uuid.UUID) error {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now()
	a.Status = "published"
	a.PublishedAt = &now
	return s.repo.Update(ctx, a)
}
