package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type FeaturedZoneService interface {
	List(ctx context.Context, venueID uuid.UUID) ([]*domain.FeaturedZone, error)
	ListActive(ctx context.Context, venueID uuid.UUID) ([]*domain.FeaturedZone, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.FeaturedZone, error)
	Create(ctx context.Context, z *domain.FeaturedZone) error
	Update(ctx context.Context, z *domain.FeaturedZone) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type featuredZoneService struct {
	repo repository.FeaturedZoneRepository
}

func NewFeaturedZoneService(repo repository.FeaturedZoneRepository) FeaturedZoneService {
	return &featuredZoneService{repo: repo}
}

func (s *featuredZoneService) List(ctx context.Context, venueID uuid.UUID) ([]*domain.FeaturedZone, error) {
	return s.repo.List(ctx, venueID)
}

func (s *featuredZoneService) ListActive(ctx context.Context, venueID uuid.UUID) ([]*domain.FeaturedZone, error) {
	return s.repo.ListActive(ctx, venueID)
}

func (s *featuredZoneService) Get(ctx context.Context, id uuid.UUID) (*domain.FeaturedZone, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *featuredZoneService) Create(ctx context.Context, z *domain.FeaturedZone) error {
	return s.repo.Create(ctx, z)
}

func (s *featuredZoneService) Update(ctx context.Context, z *domain.FeaturedZone) error {
	return s.repo.Update(ctx, z)
}

func (s *featuredZoneService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
