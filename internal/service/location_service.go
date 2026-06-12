package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/cdn"
	"github.com/hhung06/digimap-backend/internal/platform/crypto"
	"github.com/hhung06/digimap-backend/internal/platform/storage"
	"github.com/hhung06/digimap-backend/internal/repository"
)

// ── Location category service ─────────────────────────────────────────────────

type LocationCategoryService interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.LocationCategory, error)
	List(ctx context.Context, venueID uuid.UUID) ([]*domain.LocationCategory, error)
	Create(ctx context.Context, c *domain.LocationCategory) error
	Update(ctx context.Context, c *domain.LocationCategory) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type locationCategoryService struct {
	repo repository.LocationCategoryRepository
}

func NewLocationCategoryService(repo repository.LocationCategoryRepository) LocationCategoryService {
	return &locationCategoryService{repo: repo}
}

func (s *locationCategoryService) Get(ctx context.Context, id uuid.UUID) (*domain.LocationCategory, error) {
	cat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	cats, err := s.repo.List(ctx, cat.VenueID)
	if err != nil {
		return nil, err
	}
	for _, candidate := range cats {
		if candidate.ID == id {
			return candidate, nil
		}
	}
	return cat, nil
}

func (s *locationCategoryService) List(ctx context.Context, venueID uuid.UUID) ([]*domain.LocationCategory, error) {
	return s.repo.List(ctx, venueID)
}

func (s *locationCategoryService) Create(ctx context.Context, c *domain.LocationCategory) error {
	if c.Source == "" {
		c.Source = "internal"
	}
	if err := s.validateHierarchy(ctx, c); err != nil {
		return err
	}
	return s.repo.Create(ctx, c)
}

func (s *locationCategoryService) Update(ctx context.Context, c *domain.LocationCategory) error {
	if err := s.validateHierarchy(ctx, c); err != nil {
		return err
	}
	return s.repo.Update(ctx, c)
}

func (s *locationCategoryService) Delete(ctx context.Context, id uuid.UUID) error {
	cat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	cats, err := s.repo.List(ctx, cat.VenueID)
	if err != nil {
		return err
	}
	for _, child := range cats {
		if child.ParentID != nil && *child.ParentID == id {
			return domain.NewValidation(map[string]string{"category": "parent category with subcategories cannot be deleted"})
		}
	}
	return s.repo.Delete(ctx, id)
}

func (s *locationCategoryService) validateHierarchy(ctx context.Context, c *domain.LocationCategory) error {
	if c.ParentID == nil {
		return nil
	}
	if *c.ParentID == uuid.Nil {
		return domain.NewValidation(map[string]string{"parent_id": "must be a valid category id"})
	}
	if c.ID != uuid.Nil && *c.ParentID == c.ID {
		return domain.NewValidation(map[string]string{"parent_id": "category cannot be its own parent"})
	}

	parent, err := s.repo.FindByID(ctx, *c.ParentID)
	if err != nil {
		return err
	}
	if parent.VenueID != c.VenueID {
		return domain.NewValidation(map[string]string{"parent_id": "parent category must belong to the same venue"})
	}
	if parent.ParentID != nil {
		return domain.NewValidation(map[string]string{"parent_id": "subcategory cannot be used as a parent"})
	}

	if c.ID == uuid.Nil {
		return nil
	}
	cats, err := s.repo.List(ctx, c.VenueID)
	if err != nil {
		return err
	}
	for _, child := range cats {
		if child.ParentID != nil && *child.ParentID == c.ID {
			return domain.NewValidation(map[string]string{"parent_id": "category with subcategories cannot become a subcategory"})
		}
	}
	return nil
}

// ── Location service ──────────────────────────────────────────────────────────

type LocationService interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Location, error)
	List(ctx context.Context, venueID uuid.UUID, typeFilter *int, p domain.Pagination) ([]*domain.Location, int64, error)
	SearchByName(ctx context.Context, venueID uuid.UUID, q string, limit int) ([]*domain.Location, error)
	Create(ctx context.Context, l *domain.Location) error
	Update(ctx context.Context, l *domain.Location, categoryIDs []uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	Duplicate(ctx context.Context, id uuid.UUID) (*domain.Location, error)
	SetTop(ctx context.Context, id uuid.UUID, isTop bool, sortIndex *int) error
	GeoSearch(ctx context.Context, lat, lng, radiusKm float64, venueID *uuid.UUID) ([]*domain.Location, error)

	// Images
	ListImages(ctx context.Context, locationID uuid.UUID) ([]*domain.LocationImage, error)
	CreateImage(ctx context.Context, img *domain.LocationImage) error
	DeleteImage(ctx context.Context, locationID, imageID uuid.UUID) error
}

type locationService struct {
	repo        repository.LocationRepository
	venueRepo   repository.VenueRepository
	storer      storage.Storer
	invalidator cdn.Invalidator
	appVersions *AppVersionService
	env         string
}

func NewLocationService(
	repo repository.LocationRepository,
	venueRepo repository.VenueRepository,
	storer storage.Storer,
	invalidator cdn.Invalidator,
	appVersions *AppVersionService,
	env string,
) LocationService {
	return &locationService{
		repo:        repo,
		venueRepo:   venueRepo,
		storer:      storer,
		invalidator: invalidator,
		appVersions: appVersions,
		env:         env,
	}
}

func (s *locationService) Get(ctx context.Context, id uuid.UUID) (*domain.Location, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *locationService) List(ctx context.Context, venueID uuid.UUID, typeFilter *int, p domain.Pagination) ([]*domain.Location, int64, error) {
	p.Normalize()
	return s.repo.List(ctx, venueID, typeFilter, p)
}

func (s *locationService) SearchByName(ctx context.Context, venueID uuid.UUID, q string, limit int) ([]*domain.Location, error) {
	return s.repo.SearchByName(ctx, venueID, q, limit)
}

func (s *locationService) Create(ctx context.Context, l *domain.Location) error {
	if l.Source == "" {
		l.Source = "internal"
	}
	if l.IsSearchable == false {
		l.IsSearchable = true // default
	}
	if err := s.repo.Create(ctx, l); err != nil {
		return err
	}
	// Set initial categories if provided
	if len(l.Categories) > 0 {
		ids := make([]uuid.UUID, len(l.Categories))
		for i, c := range l.Categories {
			ids[i] = c.ID
		}
		return s.repo.SetCategories(ctx, l.ID, ids)
	}
	return nil
}

func (s *locationService) Update(ctx context.Context, l *domain.Location, categoryIDs []uuid.UUID) error {
	if err := s.repo.Update(ctx, l); err != nil {
		return err
	}
	return s.repo.SetCategories(ctx, l.ID, categoryIDs)
}

func (s *locationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// Duplicate creates a copy of the location with the same venue/level assignment.
func (s *locationService) Duplicate(ctx context.Context, id uuid.UUID) (*domain.Location, error) {
	src, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	clone := *src
	clone.ID = uuid.Nil
	clone.CommonName = src.CommonName + " (copy)"
	clone.IsTopLocation = false
	clone.TopLocationSortIndex = nil

	if err := s.repo.Create(ctx, &clone); err != nil {
		return nil, err
	}
	// Copy categories
	if len(src.Categories) > 0 {
		ids := make([]uuid.UUID, len(src.Categories))
		for i, c := range src.Categories {
			ids[i] = c.ID
		}
		if err := s.repo.SetCategories(ctx, clone.ID, ids); err != nil {
			return nil, err
		}
	}
	return &clone, nil
}

func (s *locationService) SetTop(ctx context.Context, id uuid.UUID, isTop bool, sortIndex *int) error {
	loc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.SetTopLocation(ctx, id, isTop, sortIndex); err != nil {
		return err
	}
	// Run publish in the background so the HTTP response is not blocked.
	// Mirrors Django's threading.Thread approach (api/locations/views.py:368).
	venueID := loc.VenueID
	go func() {
		if err := s.publishTopLocations(context.Background(), venueID); err != nil {
			fmt.Printf("top-location publish error venue=%s: %v\n", venueID, err)
		}
	}()
	return nil
}

type topLocationBundleItem struct {
	ID            uuid.UUID `json:"id"`
	TopLogo       string    `json:"top_logo"`
	TopLogoType   string    `json:"top_logo_type"`
	IsTopLocation bool      `json:"is_top_location"`
}

// publishTopLocations compresses, AES-encrypts, uploads to S3, invalidates CloudFront,
// and bumps the force-sync version — matching Django's upload_top_location pipeline
// (indoormap-backend/api/locations/views.py:383).
func (s *locationService) publishTopLocations(ctx context.Context, venueID uuid.UUID) error {
	venue, err := s.venueRepo.FindByID(ctx, venueID)
	if err != nil {
		return fmt.Errorf("load venue: %w", err)
	}

	locations, err := s.repo.ListTopLocations(ctx, venueID)
	if err != nil {
		return err
	}
	if len(locations) == 0 {
		return nil
	}

	items := make([]topLocationBundleItem, len(locations))
	for i, loc := range locations {
		items[i] = topLocationBundleItem{
			ID:            loc.ID,
			TopLogo:       loc.TopLogo,
			TopLogoType:   loc.TopLogoType,
			IsTopLocation: loc.IsTopLocation,
		}
	}

	ciphertext, err := crypto.EncryptBytes(venue.PublicKey, items)
	if err != nil {
		return fmt.Errorf("encrypt top-location bundle: %w", err)
	}

	key := storage.TopLocationKey(s.env, venueID)
	meta := map[string]string{"encrypted": "AES", "compressed": "gzip"}
	if err := s.storer.PutEncrypted(ctx, key, []byte(ciphertext), meta); err != nil {
		return fmt.Errorf("upload top-location bundle: %w", err)
	}

	if _, err := s.invalidator.Invalidate(ctx, []string{"/" + key}); err != nil {
		fmt.Printf("top-location cf invalidation failed venue=%s: %v\n", venueID, err)
	}

	if _, err := s.appVersions.Bump(ctx, venueID); err != nil {
		fmt.Printf("top-location version bump failed venue=%s: %v\n", venueID, err)
	}

	return nil
}

func (s *locationService) ListImages(ctx context.Context, locationID uuid.UUID) ([]*domain.LocationImage, error) {
	return s.repo.ListImages(ctx, locationID)
}

func (s *locationService) CreateImage(ctx context.Context, img *domain.LocationImage) error {
	return s.repo.CreateImage(ctx, img)
}

func (s *locationService) GeoSearch(ctx context.Context, lat, lng, radiusKm float64, venueID *uuid.UUID) ([]*domain.Location, error) {
	return s.repo.GeoSearch(ctx, lat, lng, radiusKm, venueID)
}

func (s *locationService) DeleteImage(ctx context.Context, locationID, imageID uuid.UUID) error {
	imgs, err := s.repo.ListImages(ctx, locationID)
	if err != nil {
		return err
	}
	for _, img := range imgs {
		if img.ID == imageID {
			return s.repo.DeleteImage(ctx, imageID)
		}
	}
	return domain.NewNotFound("image not found")
}
