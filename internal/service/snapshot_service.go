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
	"github.com/hhung06/digimap-backend/internal/platform/search"
	"github.com/hhung06/digimap-backend/internal/platform/storage"
	"github.com/hhung06/digimap-backend/internal/repository"
	"github.com/hhung06/digimap-backend/internal/service/bundle"
	"github.com/hhung06/digimap-backend/internal/service/searchindex"
	applog "github.com/hhung06/digimap-backend/log"
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
	repo                 repository.SnapshotRepository
	venueRepo            repository.VenueRepository
	languageRepo         repository.LanguageRepository
	locationRepo         repository.LocationRepository
	locationCategoryRepo repository.LocationCategoryRepository
	productRepo          repository.ProductRepository
	themeRepo            repository.ThemeRepository
	storer               storage.Storer // snapshot bucket: draft/publish JSON blobs
	assetStorer          storage.Storer // assets bucket: encrypted bundle files
	invalidator          cdn.Invalidator
	appVersions          *AppVersionService
	env                  string
	searcher             search.Searcher
	logger               applog.Logger
}

// NewSnapshotService creates a SnapshotService.
// storer targets the snapshot bucket (draft/publish JSON blobs).
// assetStorer targets the assets bucket (encrypted bundle files written during publish).
func NewSnapshotService(
	repo repository.SnapshotRepository,
	venueRepo repository.VenueRepository,
	languageRepo repository.LanguageRepository,
	locationRepo repository.LocationRepository,
	locationCategoryRepo repository.LocationCategoryRepository,
	productRepo repository.ProductRepository,
	themeRepo repository.ThemeRepository,
	storer storage.Storer,
	assetStorer storage.Storer,
	invalidator cdn.Invalidator,
	appVersions *AppVersionService,
	env string,
	searcher search.Searcher,
	logger applog.Logger,
) SnapshotService {
	return &snapshotService{
		repo:                 repo,
		venueRepo:            venueRepo,
		languageRepo:         languageRepo,
		locationRepo:         locationRepo,
		locationCategoryRepo: locationCategoryRepo,
		productRepo:          productRepo,
		themeRepo:            themeRepo,
		storer:               storer,
		assetStorer:          assetStorer,
		invalidator:          invalidator,
		appVersions:          appVersions,
		env:                  env,
		searcher:             searcher,
		logger:               logger,
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

	key := storage.SnapshotDraftKey(s.env, venueID, id)
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
		pruned, err := s.repo.DeleteOldestDraft(ctx, venueID) // best-effort prune
		if err == nil && pruned != nil {
			pruneKey := storage.SnapshotDraftKey(s.env, pruned.VenueID, pruned.ID)
			if err := s.storer.DeleteObject(ctx, pruneKey); err != nil {
				fmt.Printf("[snapshot] CreateDraft prune: s3 delete key=%s: %v\n", pruneKey, err)
			}
		}
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

	go s.publishArtifacts(context.Background(), snap)

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
	// Re-push the four public split artifacts from the reverted snapshot's draft JSON.
	// Matches Django's revert: upload_bundle_data_v2(state=public), no per-language regen.
	go s.uploadBundleData(context.Background(), snap)
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

	go s.publishArtifacts(context.Background(), snap)

	return snap, nil
}

func (s *snapshotService) Delete(ctx context.Context, id uuid.UUID) error {
	snap, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	// Best-effort S3 cleanup; do not block DB delete on S3 error.
	key := storage.SnapshotDraftKey(s.env, snap.VenueID, snap.ID)
	if err := s.storer.DeleteObject(ctx, key); err != nil {
		fmt.Printf("[snapshot] Delete: s3 delete key=%s: %v\n", key, err)
	}
	return s.repo.Delete(ctx, id)
}

// publishArtifacts builds and uploads the full publish payload: four viewer split bundles
// (base/overview/location_simple/metadata) plus per-language digiapp.{lang} and digimap.{lang}
// overlays. Mirrors Django's upload_bundle_data_v2 + publish_venue_v2 + export_localized_data.
// Called by Publish and AutoPublish. Runs in a goroutine; errors are logged, not returned.
func (s *snapshotService) publishArtifacts(ctx context.Context, snap *domain.Snapshot) {
	venue, err := s.venueRepo.FindByID(ctx, snap.VenueID)
	if err != nil {
		fmt.Printf("[snapshot] publishArtifacts: load venue %s: %v\n", snap.VenueID, err)
		return
	}

	langs, err := s.languageRepo.ListEnabled(ctx, snap.VenueID)
	if err != nil {
		fmt.Printf("[snapshot] publishArtifacts: list languages venue=%s: %v\n", snap.VenueID, err)
		return
	}
	if len(langs) == 0 {
		fmt.Printf("[snapshot] publishArtifacts: no enabled languages for venue=%s, skipping\n", snap.VenueID)
		return
	}

	// Load shared data once — per-language differences are in serialization only.
	allPagination := domain.Pagination{Page: 1, PageSize: 100000}

	locations, _, err := s.locationRepo.List(ctx, snap.VenueID, nil, nil, allPagination)
	if err != nil {
		fmt.Printf("[snapshot] publishArtifacts: list locations: %v\n", err)
		return
	}

	locationCats, err := s.locationCategoryRepo.List(ctx, snap.VenueID)
	if err != nil {
		fmt.Printf("[snapshot] publishArtifacts: list location categories: %v\n", err)
		return
	}

	products, _, err := s.productRepo.List(ctx, snap.VenueID, allPagination)
	if err != nil {
		fmt.Printf("[snapshot] publishArtifacts: list products: %v\n", err)
		return
	}

	productCats, err := s.productRepo.ListCategories(ctx, snap.VenueID)
	if err != nil {
		fmt.Printf("[snapshot] publishArtifacts: list product categories: %v\n", err)
		return
	}

	// Fetch the snapshot draft JSON (the map metadata uploaded by the admin).
	draftKey := storage.SnapshotDraftKey(s.env, snap.VenueID, snap.ID)
	rawMeta, err := s.storer.GetObject(ctx, draftKey)
	if err != nil {
		fmt.Printf("[snapshot] publishArtifacts: fetch draft metadata key=%s: %v\n", draftKey, err)
		rawMeta = []byte("{}")
	}
	var metadata map[string]any
	if err := json.Unmarshal(rawMeta, &metadata); err != nil {
		metadata = map[string]any{}
	}

	// Extract the four viewer-bundle fields the client pre-split before sending.
	// The FE sends { base, overview, location_simple, metadata } as the draft bundle;
	// the server reads them directly rather than splitting the merged SDK object.
	base := asMap(metadata["base"])
	over := asMap(metadata["overview"])
	locSim := asMap(metadata["location_simple"])
	meta := asMap(metadata["metadata"])

	// Resolve the active venue theme: prefer DB record; fall back to the client-sent
	// overview.theme (which lives inside the pre-split bundle, not at the top level).
	if venueTheme, err := s.themeRepo.FindVenueTheme(ctx, snap.VenueID); err == nil && venueTheme != nil {
		var themeData any
		if json.Unmarshal(venueTheme.Data, &themeData) == nil && themeData != nil {
			over["theme"] = themeData
		}
	}
	// else: leave over["theme"] as the client sent it (matches Django views.py:1213-1215).

	encMeta := map[string]string{"encrypted": "AES", "compressed": "gzip"}

	var invalidationPaths []string

	// Upload the four viewer-fetched split bundles.
	// Encoding: gzip(BestCompression) → AES-CBC → base64, matching Django's upload_bundle_data_v2.
	type splitEntry struct {
		key     string
		payload map[string]any
	}
	for _, e := range []splitEntry{
		{storage.BaseKey(s.env, snap.VenueID), base},
		{storage.OverviewKey(s.env, snap.VenueID), over},
		{storage.LocationSimpleKey(s.env, snap.VenueID), locSim},
		{storage.MetadataKey(s.env, snap.VenueID), meta},
	} {
		cipher, err := crypto.EncryptBytes(venue.PublicKey, e.payload)
		if err != nil {
			fmt.Printf("[snapshot] publishArtifacts: encrypt split key=%s: %v\n", e.key, err)
			continue
		}
		if err := s.assetStorer.PutEncrypted(ctx, e.key, []byte(cipher), encMeta); err != nil {
			fmt.Printf("[snapshot] publishArtifacts: upload split key=%s: %v\n", e.key, err)
			continue
		}
		invalidationPaths = append(invalidationPaths, "/"+e.key)
	}

	var failedLangs []string

	for _, lang := range langs {
		// Existing .digiapp.{lang} six-key bundle — kept for legacy consumers.
		// Pass meta (the metadata sub-bundle) as the base; it holds locations, categories, etc.
		langBundle := bundle.AssembleLanguageBundle(
			venue, lang.Code, locations, locationCats, products, productCats, meta,
		)

		cipherText, err := crypto.EncryptBundle(venue.PublicKey, langBundle)
		if err != nil {
			fmt.Printf("[snapshot] publishArtifacts: encrypt lang=%s: %v\n", lang.Code, err)
			failedLangs = append(failedLangs, lang.Code)
			continue
		}

		key := storage.DigimapV2Key(s.env, snap.VenueID, lang.Code)
		if err := s.assetStorer.PutEncrypted(ctx, key, []byte(cipherText), encMeta); err != nil {
			fmt.Printf("[snapshot] publishArtifacts: upload lang=%s: %v\n", lang.Code, err)
			failedLangs = append(failedLangs, lang.Code)
			continue
		}
		invalidationPaths = append(invalidationPaths, "/"+key)

		// Per-language .digimap.{lang} overlay — what the viewer fetches for non-English locales.
		// Encoding: AES-CBC only (no gzip), matching Django's export_localized_data.py.
		overlay := bundle.AssembleLocalizedOverlay(locations, locationCats)
		overlayCipher, err := crypto.EncryptJSON(venue.PublicKey, overlay)
		if err != nil {
			fmt.Printf("[snapshot] publishArtifacts: encrypt overlay lang=%s: %v\n", lang.Code, err)
		} else {
			lkey := storage.LocalizedKey(s.env, snap.VenueID, lang.Code)
			if err := s.assetStorer.PutObject(ctx, lkey, []byte(overlayCipher)); err != nil {
				fmt.Printf("[snapshot] publishArtifacts: upload overlay lang=%s: %v\n", lang.Code, err)
			} else {
				invalidationPaths = append(invalidationPaths, "/"+lkey)
			}
		}
	}

	if len(failedLangs) > 0 {
		fmt.Printf("[snapshot] publishArtifacts: failed languages=%v venue=%s snapshot=%s\n", failedLangs, snap.VenueID, snap.ID)
	}

	// Index exhibitors and products into OpenSearch (no-op if searcher is a log stub).
	if s.searcher != nil {
		searchindex.IndexVenueAsync(ctx, s.logger, s.searcher, venue, locations, products, s.env)
	}

	// Upload custom themes for this venue and add their paths to the invalidation batch.
	// Mirrors Django's upload_venue_themes() called during publish.
	themePaths := s.uploadCustomThemes(ctx, snap.VenueID)
	invalidationPaths = append(invalidationPaths, themePaths...)

	// Batched CloudFront invalidation for all language bundles.
	if len(invalidationPaths) > 0 {
		if _, err := s.invalidator.Invalidate(ctx, invalidationPaths); err != nil {
			fmt.Printf("[snapshot] publishArtifacts: CF invalidation: %v\n", err)
		}
	}

	// Bump the force-sync version (also uploads latest-bundle JSON + CF invalidation).
	if _, err := s.appVersions.Bump(ctx, snap.VenueID); err != nil {
		fmt.Printf("[snapshot] publishArtifacts: bump version venue=%s: %v\n", snap.VenueID, err)
	}
}

// uploadBundleData builds and uploads the four viewer split bundles (base/overview/
// location_simple/metadata) for the given snapshot and invalidates CloudFront.
// Mirrors Django's upload_bundle_data_v2. Called by Revert (splits-only, no per-language regen).
// Runs in a goroutine; errors are logged, not returned.
func (s *snapshotService) uploadBundleData(ctx context.Context, snap *domain.Snapshot) {
	venue, err := s.venueRepo.FindByID(ctx, snap.VenueID)
	if err != nil {
		fmt.Printf("[snapshot] uploadBundleData: load venue %s: %v\n", snap.VenueID, err)
		return
	}

	draftKey := storage.SnapshotDraftKey(s.env, snap.VenueID, snap.ID)
	rawMeta, err := s.storer.GetObject(ctx, draftKey)
	if err != nil {
		fmt.Printf("[snapshot] uploadBundleData: fetch draft key=%s: %v\n", draftKey, err)
		rawMeta = []byte("{}")
	}
	var metadata map[string]any
	if err := json.Unmarshal(rawMeta, &metadata); err != nil {
		metadata = map[string]any{}
	}

	base := asMap(metadata["base"])
	over := asMap(metadata["overview"])
	locSim := asMap(metadata["location_simple"])
	meta := asMap(metadata["metadata"])

	if venueTheme, err := s.themeRepo.FindVenueTheme(ctx, snap.VenueID); err == nil && venueTheme != nil {
		var themeData any
		if json.Unmarshal(venueTheme.Data, &themeData) == nil && themeData != nil {
			over["theme"] = themeData
		}
	}

	encMeta := map[string]string{"encrypted": "AES", "compressed": "gzip"}

	type splitEntry struct {
		key     string
		payload map[string]any
	}
	var invalidationPaths []string
	for _, e := range []splitEntry{
		{storage.BaseKey(s.env, snap.VenueID), base},
		{storage.OverviewKey(s.env, snap.VenueID), over},
		{storage.LocationSimpleKey(s.env, snap.VenueID), locSim},
		{storage.MetadataKey(s.env, snap.VenueID), meta},
	} {
		cipher, err := crypto.EncryptBytes(venue.PublicKey, e.payload)
		if err != nil {
			fmt.Printf("[snapshot] uploadBundleData: encrypt split key=%s: %v\n", e.key, err)
			continue
		}
		if err := s.assetStorer.PutEncrypted(ctx, e.key, []byte(cipher), encMeta); err != nil {
			fmt.Printf("[snapshot] uploadBundleData: upload split key=%s: %v\n", e.key, err)
			continue
		}
		invalidationPaths = append(invalidationPaths, "/"+e.key)
	}

	if len(invalidationPaths) > 0 {
		if _, err := s.invalidator.Invalidate(ctx, invalidationPaths); err != nil {
			fmt.Printf("[snapshot] uploadBundleData: CF invalidation: %v\n", err)
		}
	}
}

// uploadCustomThemes uploads all custom themes for a venue to S3, updates their
// storage_path in the DB, and returns S3 paths for CloudFront invalidation.
// Mirrors Django's upload_venue_themes() called during venue publish.
func (s *snapshotService) uploadCustomThemes(ctx context.Context, venueID uuid.UUID) []string {
	themes, err := s.themeRepo.List(ctx, venueID)
	if err != nil {
		fmt.Printf("[snapshot] uploadCustomThemes: list themes venue=%s: %v\n", venueID, err)
		return nil
	}

	var paths []string
	for _, t := range themes {
		if t.Scope != domain.ThemeScopeCustom {
			continue
		}
		key := storage.CustomThemeKey(s.env, venueID, t.Name)
		if err := s.assetStorer.PutObject(ctx, key, []byte(t.Data)); err != nil {
			fmt.Printf("[snapshot] uploadCustomThemes: upload theme=%s: %v\n", t.ID, err)
			continue
		}
		t.StoragePath = key
		if err := s.themeRepo.Update(ctx, t); err != nil {
			fmt.Printf("[snapshot] uploadCustomThemes: update storage_path theme=%s: %v\n", t.ID, err)
		}
		paths = append(paths, "/"+key)
	}
	return paths
}

// asMap safely casts v to map[string]any. Returns an empty map if v is nil or not the right type.
// Used to extract pre-split viewer bundles from the client-sent draft JSON.
func asMap(v any) map[string]any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return map[string]any{}
}
