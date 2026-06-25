package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/storage"
	"github.com/hhung06/digimap-backend/internal/repository"
)

// LevelBundleService manages per-level map data bundles.
type LevelBundleService interface {
	ListBySnapshot(ctx context.Context, snapshotID uuid.UUID) ([]*domain.LevelBundle, error)
	Create(ctx context.Context, snapshotID, venueID, levelID uuid.UUID, bundle []byte) (*domain.LevelBundle, error)
	GetContent(ctx context.Context, id uuid.UUID) ([]byte, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type levelBundleService struct {
	repo         repository.LevelBundleRepository
	snapshotRepo repository.SnapshotRepository
	storer       storage.Storer
	env          string
}

// NewLevelBundleService creates a LevelBundleService.
func NewLevelBundleService(
	repo repository.LevelBundleRepository,
	snapshotRepo repository.SnapshotRepository,
	storer storage.Storer,
	env string,
) LevelBundleService {
	return &levelBundleService{
		repo:         repo,
		snapshotRepo: snapshotRepo,
		storer:       storer,
		env:          env,
	}
}

func (s *levelBundleService) ListBySnapshot(ctx context.Context, snapshotID uuid.UUID) ([]*domain.LevelBundle, error) {
	return s.repo.ListBySnapshot(ctx, snapshotID)
}

// Create uploads the bundle to S3 then inserts a LevelBundle record.
// The new bundle's state is inherited from its parent snapshot.
func (s *levelBundleService) Create(ctx context.Context, snapshotID, venueID, levelID uuid.UUID, bundle []byte) (*domain.LevelBundle, error) {
	parent, err := s.snapshotRepo.FindByID(ctx, snapshotID)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		id = uuid.New()
	}

	key := storage.LevelBundleKey(s.env, venueID, snapshotID, levelID)
	if err := s.storer.PutObject(ctx, key, bundle); err != nil {
		return nil, fmt.Errorf("upload level bundle: %w", err)
	}

	b := &domain.LevelBundle{
		ID:         id,
		SnapshotID: snapshotID,
		VenueID:    venueID,
		LevelID:    levelID,
		State:      parent.State,
	}
	if err := s.repo.Create(ctx, b); err != nil {
		return nil, fmt.Errorf("create level bundle record: %w", err)
	}
	return b, nil
}

// GetContent fetches the raw stored bundle blob from S3 for the given bundle ID.
func (s *levelBundleService) GetContent(ctx context.Context, id uuid.UUID) ([]byte, error) {
	b, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	key := storage.LevelBundleKey(s.env, b.VenueID, b.SnapshotID, b.LevelID)
	data, err := s.storer.GetObject(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("fetch level bundle content: %w", err)
	}
	return data, nil
}

func (s *levelBundleService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
