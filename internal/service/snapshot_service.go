package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/cdn"
	"github.com/hhung06/digimap-backend/internal/platform/crypto"
	"github.com/hhung06/digimap-backend/internal/platform/storage"
	"github.com/hhung06/digimap-backend/internal/repository"
	"github.com/hhung06/digimap-backend/internal/service/bundle"
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
	repo                repository.SnapshotRepository
	venueRepo           repository.VenueRepository
	languageRepo        repository.LanguageRepository
	locationRepo        repository.LocationRepository
	locationCategoryRepo repository.LocationCategoryRepository
	productRepo         repository.ProductRepository
	storer              storage.Storer
	invalidator         cdn.Invalidator
	appVersions         *AppVersionService
	env                 string
}

// NewSnapshotService creates a SnapshotService.
func NewSnapshotService(
	repo repository.SnapshotRepository,
	venueRepo repository.VenueRepository,
	languageRepo repository.LanguageRepository,
	locationRepo repository.LocationRepository,
	locationCategoryRepo repository.LocationCategoryRepository,
	productRepo repository.ProductRepository,
	storer storage.Storer,
	invalidator cdn.Invalidator,
	appVersions *AppVersionService,
	env string,
) SnapshotService {
	return &snapshotService{
		repo:                repo,
		venueRepo:           venueRepo,
		languageRepo:        languageRepo,
		locationRepo:        locationRepo,
		locationCategoryRepo: locationCategoryRepo,
		productRepo:         productRepo,
		storer:              storer,
		invalidator:         invalidator,
		appVersions:         appVersions,
		env:                 env,
	}
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

// Publish flips a draft snapshot to Public and launches the v2 bundle upload in the background.
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

	go s.publishV2(context.Background(), snap)

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

// AutoPublish republishes the latest existing draft snapshot for a venue.
// Mirrors Django's auto-publish behavior in publish_venue_v2.py.
func (s *snapshotService) AutoPublish(ctx context.Context, venueID, _ uuid.UUID) (*domain.Snapshot, error) {
	snap, err := s.repo.LatestDraft(ctx, venueID)
	if err != nil {
		return nil, fmt.Errorf("find latest draft: %w", err)
	}
	if snap == nil {
		return nil, domain.NewNotFound("no draft snapshot found for auto-publish")
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

	go s.publishV2(context.Background(), snap)

	return snap, nil
}

func (s *snapshotService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// publishV2 assembles per-language bundles from the DB, encrypts them, uploads to S3,
// and invalidates CloudFront. Mirrors publish_venue_v2.py.
// Runs in a goroutine; errors are logged, not returned.
func (s *snapshotService) publishV2(ctx context.Context, snap *domain.Snapshot) {
	venue, err := s.venueRepo.FindByID(ctx, snap.VenueID)
	if err != nil {
		fmt.Printf("[snapshot] publishV2: load venue %s: %v\n", snap.VenueID, err)
		return
	}

	langs, err := s.languageRepo.ListEnabled(ctx, snap.VenueID)
	if err != nil {
		fmt.Printf("[snapshot] publishV2: list languages venue=%s: %v\n", snap.VenueID, err)
		return
	}
	if len(langs) == 0 {
		fmt.Printf("[snapshot] publishV2: no enabled languages for venue=%s, skipping\n", snap.VenueID)
		return
	}

	// Load shared data once — per-language differences are in serialization only.
	allPagination := domain.Pagination{Page: 1, PageSize: 100000}

	locations, _, err := s.locationRepo.List(ctx, snap.VenueID, nil, allPagination)
	if err != nil {
		fmt.Printf("[snapshot] publishV2: list locations: %v\n", err)
		return
	}

	locationCats, err := s.locationCategoryRepo.List(ctx, snap.VenueID)
	if err != nil {
		fmt.Printf("[snapshot] publishV2: list location categories: %v\n", err)
		return
	}

	products, _, err := s.productRepo.List(ctx, snap.VenueID, allPagination)
	if err != nil {
		fmt.Printf("[snapshot] publishV2: list products: %v\n", err)
		return
	}

	productCats, err := s.productRepo.ListCategories(ctx, snap.VenueID)
	if err != nil {
		fmt.Printf("[snapshot] publishV2: list product categories: %v\n", err)
		return
	}

	// Fetch the snapshot draft JSON (the map metadata uploaded by the admin).
	draftKey := storage.SnapshotDraftKey(s.env, snap.VenueID, snap.ID)
	rawMeta, err := s.storer.GetObject(ctx, draftKey)
	if err != nil {
		fmt.Printf("[snapshot] publishV2: fetch draft metadata key=%s: %v\n", draftKey, err)
		rawMeta = []byte("{}")
	}
	var metadata map[string]any
	if err := json.Unmarshal(rawMeta, &metadata); err != nil {
		metadata = map[string]any{}
	}

	encMeta := map[string]string{"encrypted": "AES", "compressed": "gzip"}

	var invalidationPaths []string
	var failedLangs []string

	for _, lang := range langs {
		langBundle := bundle.AssembleLanguageBundle(
			venue, lang.Code, locations, locationCats, products, productCats, metadata,
		)

		cipherText, err := crypto.EncryptBundle(venue.PublicKey, langBundle)
		if err != nil {
			fmt.Printf("[snapshot] publishV2: encrypt lang=%s: %v\n", lang.Code, err)
			failedLangs = append(failedLangs, lang.Code)
			continue
		}

		key := storage.DigimapV2Key(s.env, snap.VenueID, lang.Code)
		if err := s.storer.PutEncrypted(ctx, key, []byte(cipherText), encMeta); err != nil {
			fmt.Printf("[snapshot] publishV2: upload lang=%s: %v\n", lang.Code, err)
			failedLangs = append(failedLangs, lang.Code)
			continue
		}

		invalidationPaths = append(invalidationPaths, "/"+key)
	}

	if len(failedLangs) > 0 {
		fmt.Printf("[snapshot] publishV2: failed languages=%v venue=%s snapshot=%s\n", failedLangs, snap.VenueID, snap.ID)
	}

	// Batched CloudFront invalidation for all language bundles.
	if len(invalidationPaths) > 0 {
		if _, err := s.invalidator.Invalidate(ctx, invalidationPaths); err != nil {
			fmt.Printf("[snapshot] publishV2: CF invalidation: %v\n", err)
		}
	}

	// Bump the force-sync version (also uploads latest-bundle JSON + CF invalidation).
	if _, err := s.appVersions.Bump(ctx, snap.VenueID); err != nil {
		fmt.Printf("[snapshot] publishV2: bump version venue=%s: %v\n", snap.VenueID, err)
	}
}
