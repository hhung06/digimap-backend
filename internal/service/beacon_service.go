package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

// BeaconService manages BLE beacon CRUD.
type BeaconService interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Beacon, int64, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Beacon, error)
	Create(ctx context.Context, b *domain.Beacon) error
	Update(ctx context.Context, b *domain.Beacon) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type beaconService struct {
	repo repository.BeaconRepository
}

// NewBeaconService creates a BeaconService.
func NewBeaconService(repo repository.BeaconRepository) BeaconService {
	return &beaconService{repo: repo}
}

func (s *beaconService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Beacon, int64, error) {
	return s.repo.List(ctx, venueID, p)
}

func (s *beaconService) Get(ctx context.Context, id uuid.UUID) (*domain.Beacon, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *beaconService) Create(ctx context.Context, b *domain.Beacon) error {
	return s.repo.Create(ctx, b)
}

func (s *beaconService) Update(ctx context.Context, b *domain.Beacon) error {
	return s.repo.Update(ctx, b)
}

func (s *beaconService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
