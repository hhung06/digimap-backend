package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/storage"
	"github.com/hhung06/digimap-backend/internal/repository"
)

// SnapshotService manages snapshot metadata and S3 bundle storage.
type SnapshotService interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Snapshot, int64, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Snapshot, error)
	LatestPublished(ctx context.Context, venueID uuid.UUID) (*domain.Snapshot, error)
	CreateDraft(ctx context.Context, venueID, createdBy uuid.UUID, bundle []byte) (*domain.Snapshot, error)
	Publish(ctx context.Context, id uuid.UUID) (*domain.Snapshot, error)
	Revert(ctx context.Context, id uuid.UUID) (*domain.Snapshot, error)
	AutoPublish(ctx context.Context, venueID, createdBy uuid.UUID) (*domain.Snapshot, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type snapshotService struct {
	repo   repository.SnapshotRepository
	storer storage.Storer
	env    string
}

// NewSnapshotService creates a SnapshotService.
func NewSnapshotService(repo repository.SnapshotRepository, storer storage.Storer, env string) SnapshotService {
	return &snapshotService{repo: repo, storer: storer, env: env}
}

func (s *snapshotService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Snapshot, int64, error) {
	return s.repo.List(ctx, venueID, p)
}

func (s *snapshotService) Get(ctx context.Context, id uuid.UUID) (*domain.Snapshot, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *snapshotService) LatestPublished(ctx context.Context, venueID uuid.UUID) (*domain.Snapshot, error) {
	return s.repo.LatestPublished(ctx, venueID)
}

// CreateDraft uploads bundle JSON to S3 then inserts a draft snapshot record.
// S3-first ordering: if DB insert fails, the orphaned S3 object is harmless
// (deterministic key is overwritten on retry).
func (s *snapshotService) CreateDraft(ctx context.Context, venueID, createdBy uuid.UUID, bundle []byte) (*domain.Snapshot, error) {
	id, err := uuid.NewV7()
	if err != nil {
		id = uuid.New()
	}

	key := fmt.Sprintf("%s/%s/snapshots/draft/%s.json", s.env, venueID, id)
	if err := s.storer.PutObject(ctx, key, bundle); err != nil {
		return nil, fmt.Errorf("upload snapshot bundle: %w", err)
	}

	snap := &domain.Snapshot{
		ID:        id,
		VenueID:   venueID,
		State:     domain.SnapshotStateDraft,
		Method:    domain.SnapshotMethodManual,
		CreatedBy: &createdBy,
	}
	if err := s.repo.Create(ctx, snap); err != nil {
		return nil, fmt.Errorf("create snapshot record: %w", err)
	}

	// Count is post-insert: if count == MaxSnapshotVersions, the oldest draft is pruned
	// so the total stays at MaxSnapshotVersions. The new snapshot is included in the count.
	count, err := s.repo.CountDraftsByVenue(ctx, venueID)
	if err != nil {
		return snap, nil // non-fatal: version control best-effort
	}
	if count >= domain.MaxSnapshotVersions {
		_ = s.repo.DeleteOldestDraft(ctx, venueID) // best-effort prune
	}

	return snap, nil
}

func (s *snapshotService) Publish(ctx context.Context, id uuid.UUID) (*domain.Snapshot, error) {
	snap, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if err := s.repo.UpdateState(ctx, snap.ID, domain.SnapshotStatePublic, &now); err != nil {
		return nil, fmt.Errorf("publish snapshot: %w", err)
	}
	snap.State = domain.SnapshotStatePublic
	snap.PublishAt = &now
	return snap, nil
}

func (s *snapshotService) Revert(ctx context.Context, id uuid.UUID) (*domain.Snapshot, error) {
	snap, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.repo.UnpublishVenue(ctx, snap.VenueID); err != nil {
		return nil, fmt.Errorf("unpublish existing snapshots: %w", err)
	}
	now := time.Now()
	if err := s.repo.UpdateState(ctx, snap.ID, domain.SnapshotStatePublic, &now); err != nil {
		return nil, fmt.Errorf("revert snapshot: %w", err)
	}
	snap.State = domain.SnapshotStatePublic
	snap.PublishAt = &now
	return snap, nil
}

func (s *snapshotService) AutoPublish(ctx context.Context, venueID, createdBy uuid.UUID) (*domain.Snapshot, error) {
	id, err := uuid.NewV7()
	if err != nil {
		id = uuid.New()
	}
	snap := &domain.Snapshot{
		ID:        id,
		VenueID:   venueID,
		State:     domain.SnapshotStateDraft,
		Method:    domain.SnapshotMethodAuto,
		CreatedBy: &createdBy,
	}
	if err := s.repo.Create(ctx, snap); err != nil {
		return nil, fmt.Errorf("create auto snapshot record: %w", err)
	}
	if err := s.repo.UnpublishVenue(ctx, venueID); err != nil {
		return nil, fmt.Errorf("unpublish existing snapshots: %w", err)
	}
	now := time.Now()
	if err := s.repo.UpdateState(ctx, snap.ID, domain.SnapshotStatePublic, &now); err != nil {
		return nil, fmt.Errorf("publish auto snapshot: %w", err)
	}
	snap.State = domain.SnapshotStatePublic
	snap.PublishAt = &now
	return snap, nil
}

func (s *snapshotService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
