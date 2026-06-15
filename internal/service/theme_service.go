package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/cdn"
	"github.com/hhung06/digimap-backend/internal/platform/storage"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type ThemeService interface {
	ListGlobal(ctx context.Context) ([]*domain.Theme, error)
	List(ctx context.Context, venueID uuid.UUID) ([]*domain.Theme, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Theme, error)
	CreateGlobal(ctx context.Context, name string, data json.RawMessage) (*domain.Theme, error)
	Create(ctx context.Context, venueID uuid.UUID, name string, data json.RawMessage) (*domain.Theme, error)
	Update(ctx context.Context, id uuid.UUID, name string, data json.RawMessage) (*domain.Theme, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SetVenueTheme(ctx context.Context, venueID uuid.UUID, themeID *uuid.UUID) error
}

type themeService struct {
	repo        repository.ThemeRepository
	venueRepo   repository.VenueRepository
	storer      storage.Storer
	invalidator cdn.Invalidator
	env         string
}

func NewThemeService(
	repo repository.ThemeRepository,
	venueRepo repository.VenueRepository,
	storer storage.Storer,
	invalidator cdn.Invalidator,
	env string,
) ThemeService {
	return &themeService{repo: repo, venueRepo: venueRepo, storer: storer, invalidator: invalidator, env: env}
}

func (s *themeService) ListGlobal(ctx context.Context) ([]*domain.Theme, error) {
	return s.repo.ListGlobal(ctx)
}

func (s *themeService) List(ctx context.Context, venueID uuid.UUID) ([]*domain.Theme, error) {
	return s.repo.List(ctx, venueID)
}

func (s *themeService) Get(ctx context.Context, id uuid.UUID) (*domain.Theme, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *themeService) CreateGlobal(ctx context.Context, name string, data json.RawMessage) (*domain.Theme, error) {
	key := storage.GlobalThemeKey(s.env, name)
	t := &domain.Theme{Scope: domain.ThemeScopeGlobal, Name: name, Data: data, StoragePath: key}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	s.uploadAndInvalidate(ctx, key, data)
	return t, nil
}

func (s *themeService) Create(ctx context.Context, venueID uuid.UUID, name string, data json.RawMessage) (*domain.Theme, error) {
	t := &domain.Theme{VenueID: &venueID, Scope: domain.ThemeScopeCustom, Name: name, Data: data}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *themeService) Update(ctx context.Context, id uuid.UUID, name string, data json.RawMessage) (*domain.Theme, error) {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	t.Name = name
	t.Data = data
	if t.Scope == domain.ThemeScopeGlobal {
		t.StoragePath = storage.GlobalThemeKey(s.env, name)
	}
	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	if t.Scope == domain.ThemeScopeGlobal {
		s.uploadAndInvalidate(ctx, t.StoragePath, data)
	}
	return t, nil
}

func (s *themeService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}
	used, err := s.repo.IsUsedByVenues(ctx, id)
	if err != nil {
		return err
	}
	if used {
		return domain.NewConflict("theme is currently used by one or more venues")
	}
	return s.repo.Delete(ctx, id)
}

func (s *themeService) SetVenueTheme(ctx context.Context, venueID uuid.UUID, themeID *uuid.UUID) error {
	if themeID != nil {
		t, err := s.repo.FindByID(ctx, *themeID)
		if err != nil {
			return err
		}
		if t.Scope == domain.ThemeScopeCustom && (t.VenueID == nil || *t.VenueID != venueID) {
			return domain.NewValidation(map[string]string{"theme_id": "custom theme does not belong to this venue"})
		}
	}
	return s.venueRepo.SetThemeID(ctx, venueID, themeID)
}

// uploadAndInvalidate uploads theme JSON to S3 and invalidates CloudFront.
// Errors are logged but not fatal — matching Django's fire-and-forget pattern.
func (s *themeService) uploadAndInvalidate(ctx context.Context, key string, data json.RawMessage) {
	if err := s.storer.PutObject(ctx, key, []byte(data)); err != nil {
		fmt.Printf("theme upload failed key=%s: %v\n", key, err)
		return
	}
	if _, err := s.invalidator.Invalidate(ctx, []string{"/" + key}); err != nil {
		fmt.Printf("theme cf invalidation failed key=%s: %v\n", key, err)
	}
}
