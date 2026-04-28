package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type CreateAssetInput struct {
	VenueID     uuid.UUID
	Name        string
	Key         string
	ContentType string
	SizeBytes   int64
	URL         string
	CreatedBy   *uuid.UUID
}

type AssetService interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Asset, int64, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Asset, error)
	Create(ctx context.Context, in CreateAssetInput) (*domain.Asset, error)
	Update(ctx context.Context, id uuid.UUID, name string) (*domain.Asset, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type assetService struct{ repo repository.AssetRepository }

func NewAssetService(repo repository.AssetRepository) AssetService {
	return &assetService{repo: repo}
}

func (s *assetService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Asset, int64, error) {
	return s.repo.List(ctx, venueID, p)
}

func (s *assetService) Get(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *assetService) Create(ctx context.Context, in CreateAssetInput) (*domain.Asset, error) {
	a := &domain.Asset{
		VenueID:     in.VenueID,
		Name:        in.Name,
		Key:         in.Key,
		ContentType: in.ContentType,
		SizeBytes:   in.SizeBytes,
		URL:         in.URL,
		CreatedBy:   in.CreatedBy,
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *assetService) Update(ctx context.Context, id uuid.UUID, name string) (*domain.Asset, error) {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	a.Name = name
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *assetService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
