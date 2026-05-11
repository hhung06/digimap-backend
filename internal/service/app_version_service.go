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

// AppVersionService manages the force-sync version pointer per venue.
// Mirrors Django's ForceSyncVersion bump logic in publish_venue_v2.py:962 and
// api/locations/views.py:591 (update_version).
type AppVersionService struct {
	repo       repository.AppVersionRepository
	storer     storage.Storer
	invalidator cdn.Invalidator
	env        string
}

func NewAppVersionService(
	repo repository.AppVersionRepository,
	storer storage.Storer,
	invalidator cdn.Invalidator,
	env string,
) *AppVersionService {
	return &AppVersionService{repo: repo, storer: storer, invalidator: invalidator, env: env}
}

// Bump generates a new UUID version, persists it, uploads latest-bundle JSON to S3,
// and invalidates CloudFront. Non-fatal errors are returned so callers can log and continue.
func (s *AppVersionService) Bump(ctx context.Context, venueID uuid.UUID) (uuid.UUID, error) {
	version := uuid.New()

	if err := s.repo.Upsert(ctx, venueID, version); err != nil {
		return uuid.Nil, fmt.Errorf("app_version upsert: %w", err)
	}

	key := storage.LatestBundleKey(s.env, venueID)
	payload, _ := json.Marshal(map[string]string{"version": version.String()})
	if err := s.storer.PutObject(ctx, key, payload); err != nil {
		return version, fmt.Errorf("app_version put latest-bundle: %w", err)
	}

	if _, err := s.invalidator.Invalidate(ctx, []string{"/" + key}); err != nil {
		// Non-fatal — Django also treats CF invalidation failure as non-fatal (publish_venue_v2.py:926-931).
		fmt.Printf("app_version cf invalidation failed venue=%s: %v\n", venueID, err)
	}

	return version, nil
}

// Get returns the current version for a venue (auto-creates if absent).
func (s *AppVersionService) Get(ctx context.Context, venueID uuid.UUID) (*domain.AppVersion, error) {
	return s.repo.Get(ctx, venueID)
}
