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
	Reconcile(ctx context.Context, venueID uuid.UUID, items []*domain.Language) ([]*domain.Language, error)
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
	if l.IsDefault && !l.Enabled {
		return domain.NewValidation(map[string]string{"enabled": "the default language must be enabled"})
	}
	if err := s.repo.Create(ctx, l); err != nil {
		return err
	}
	if l.IsDefault {
		return s.repo.SetDefault(ctx, l.VenueID, l.ID)
	}
	return nil
}

func (s *languageService) Update(ctx context.Context, l *domain.Language) error {
	if l.IsDefault && !l.Enabled {
		return domain.NewValidation(map[string]string{"enabled": "cannot disable the default language"})
	}
	if err := s.repo.Update(ctx, l); err != nil {
		return err
	}
	if l.IsDefault {
		return s.repo.SetDefault(ctx, l.VenueID, l.ID)
	}
	return nil
}

func (s *languageService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *languageService) Reconcile(ctx context.Context, venueID uuid.UUID, items []*domain.Language) ([]*domain.Language, error) {
	if len(items) == 0 {
		return nil, domain.NewValidation(map[string]string{"languages": "at least one language is required"})
	}
	seenCodes := make(map[string]bool, len(items))
	defaultCount := 0
	for _, item := range items {
		if seenCodes[item.Code] {
			return nil, domain.NewValidation(map[string]string{"languages": "duplicate language code: " + item.Code})
		}
		seenCodes[item.Code] = true
		if item.IsDefault {
			defaultCount++
		}
	}
	if defaultCount > 1 {
		return nil, domain.NewValidation(map[string]string{"languages": "only one language may be marked as default"})
	}
	if defaultCount == 0 {
		items[0].IsDefault = true
	}
	for _, item := range items {
		item.VenueID = venueID
	}
	return s.repo.ReplaceAll(ctx, venueID, items)
}
