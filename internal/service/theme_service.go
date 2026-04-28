package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type ThemeService interface {
	List(ctx context.Context, venueID uuid.UUID) ([]*domain.Theme, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Theme, error)
	Create(ctx context.Context, venueID uuid.UUID, name, primaryColor, secondaryColor string) (*domain.Theme, error)
	Update(ctx context.Context, id uuid.UUID, name, primaryColor, secondaryColor string) (*domain.Theme, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type themeService struct{ repo repository.ThemeRepository }

func NewThemeService(repo repository.ThemeRepository) ThemeService {
	return &themeService{repo: repo}
}

func (s *themeService) List(ctx context.Context, venueID uuid.UUID) ([]*domain.Theme, error) {
	return s.repo.List(ctx, venueID)
}

func (s *themeService) Get(ctx context.Context, id uuid.UUID) (*domain.Theme, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *themeService) Create(ctx context.Context, venueID uuid.UUID, name, primaryColor, secondaryColor string) (*domain.Theme, error) {
	t := &domain.Theme{VenueID: venueID, Name: name, PrimaryColor: primaryColor, SecondaryColor: secondaryColor}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *themeService) Update(ctx context.Context, id uuid.UUID, name, primaryColor, secondaryColor string) (*domain.Theme, error) {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	t.Name = name
	t.PrimaryColor = primaryColor
	t.SecondaryColor = secondaryColor
	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *themeService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
