package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type LanguageService interface {
	List(ctx context.Context, venueID uuid.UUID) ([]*domain.Language, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Language, error)
	Create(ctx context.Context, l *domain.Language) error
	Update(ctx context.Context, l *domain.Language) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type languageService struct {
	repo repository.LanguageRepository
}

func NewLanguageService(repo repository.LanguageRepository) LanguageService {
	return &languageService{repo: repo}
}

func (s *languageService) List(ctx context.Context, venueID uuid.UUID) ([]*domain.Language, error) {
	return s.repo.List(ctx, venueID)
}

func (s *languageService) Get(ctx context.Context, id uuid.UUID) (*domain.Language, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *languageService) Create(ctx context.Context, l *domain.Language) error {
	return s.repo.Create(ctx, l)
}

func (s *languageService) Update(ctx context.Context, l *domain.Language) error {
	return s.repo.Update(ctx, l)
}

func (s *languageService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
