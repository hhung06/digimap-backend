package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type CreateProductPlazaInput struct {
	VenueID      uuid.UUID
	Name         string
	Description  string
	Localization json.RawMessage
	LocationID   *uuid.UUID
}

type ProductPlazaService interface {
	List(ctx context.Context, venueID uuid.UUID) ([]*domain.ProductPlaza, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.ProductPlaza, error)
	Create(ctx context.Context, in CreateProductPlazaInput) (*domain.ProductPlaza, error)
	Update(ctx context.Context, id uuid.UUID, in CreateProductPlazaInput) (*domain.ProductPlaza, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type productPlazaService struct{ repo repository.ProductPlazaRepository }

func NewProductPlazaService(repo repository.ProductPlazaRepository) ProductPlazaService {
	return &productPlazaService{repo: repo}
}

func (s *productPlazaService) List(ctx context.Context, venueID uuid.UUID) ([]*domain.ProductPlaza, error) {
	return s.repo.List(ctx, venueID)
}

func (s *productPlazaService) Get(ctx context.Context, id uuid.UUID) (*domain.ProductPlaza, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *productPlazaService) Create(ctx context.Context, in CreateProductPlazaInput) (*domain.ProductPlaza, error) {
	p := &domain.ProductPlaza{
		VenueID:      in.VenueID,
		Name:         in.Name,
		Description:  in.Description,
		Localization: in.Localization,
		LocationID:   in.LocationID,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *productPlazaService) Update(ctx context.Context, id uuid.UUID, in CreateProductPlazaInput) (*domain.ProductPlaza, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	p.Name = in.Name
	p.Description = in.Description
	p.LocationID = in.LocationID
	if in.Localization != nil {
		p.Localization = in.Localization
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *productPlazaService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
