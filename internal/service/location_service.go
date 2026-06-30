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
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
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
	List(ctx context.Context, venueID uuid.UUID, typeFilter *int, categoryID *uuid.UUID, p domain.Pagination) ([]*domain.Location, int64, error)
	SearchByName(ctx context.Context, venueID uuid.UUID, q string, limit int) ([]*domain.Location, error)
	Create(ctx context.Context, l *domain.Location) error
	CreateWithMedia(ctx context.Context, l *domain.Location, replacement LocationMediaReplacement) error
	Update(ctx context.Context, l *domain.Location, categoryIDs []uuid.UUID) error
	UpdateWithMedia(ctx context.Context, l *domain.Location, categoryIDs []uuid.UUID, replacement LocationMediaReplacement) error
	Delete(ctx context.Context, id uuid.UUID) error
	Duplicate(ctx context.Context, id uuid.UUID) (*domain.Location, error)
	SetTop(ctx context.Context, id uuid.UUID, isTop bool, sortIndex *int) error
	SetTopWithMedia(ctx context.Context, id uuid.UUID, isTop bool, sortIndex *int, upload *MediaUpload) error
	GeoSearch(ctx context.Context, lat, lng, radiusKm float64, venueID *uuid.UUID) ([]*domain.Location, error)

	// Images
	ListImages(ctx context.Context, locationID uuid.UUID) ([]*domain.LocationImage, error)
	CreateImage(ctx context.Context, img *domain.LocationImage) error
	DeleteImage(ctx context.Context, locationID, imageID uuid.UUID) error
}

type LocationLogoUploads struct {
	Original MediaUpload
	Large    MediaUpload
	Medium   MediaUpload
	Small    MediaUpload
}

type LocationMediaReplacement struct {
	CommonLogo      *LocationLogoUploads
	ClearCommonLogo bool
	ReplaceImages   bool
	KeepImageIDs    []uuid.UUID
	Uploads         []MediaUpload
}

type locationService struct {
	repo        repository.LocationRepository
	venueRepo   repository.VenueRepository
	storer      storage.Storer
	invalidator cdn.Invalidator
	appVersions *AppVersionService
	env         string
	mediaSvc    MediaService
}

func NewLocationService(
	repo repository.LocationRepository,
	venueRepo repository.VenueRepository,
	storer storage.Storer,
	invalidator cdn.Invalidator,
	appVersions *AppVersionService,
	env string,
	mediaSvc ...MediaService,
) LocationService {
	var media MediaService
	if len(mediaSvc) > 0 {
		media = mediaSvc[0]
	}
	return &locationService{
		repo:        repo,
		venueRepo:   venueRepo,
		storer:      storer,
		invalidator: invalidator,
		appVersions: appVersions,
		env:         env,
		mediaSvc:    media,
	}
}

func (s *locationService) Get(ctx context.Context, id uuid.UUID) (*domain.Location, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *locationService) List(ctx context.Context, venueID uuid.UUID, typeFilter *int, categoryID *uuid.UUID, p domain.Pagination) ([]*domain.Location, int64, error) {
	p.Normalize()
	return s.repo.List(ctx, venueID, typeFilter, categoryID, p)
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

func (s *locationService) CreateWithMedia(ctx context.Context, l *domain.Location, replacement LocationMediaReplacement) error {
	if err := s.Create(ctx, l); err != nil {
		return err
	}
	categoryIDs := make([]uuid.UUID, 0, len(l.Categories))
	for _, cat := range l.Categories {
		if cat != nil {
			categoryIDs = append(categoryIDs, cat.ID)
		}
	}
	return s.UpdateWithMedia(ctx, l, categoryIDs, replacement)
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

func (s *locationService) UpdateWithMedia(ctx context.Context, l *domain.Location, categoryIDs []uuid.UUID, replacement LocationMediaReplacement) error {
	oldLogoKeys := locationLogoKeys(l)
	newLogoKeys := map[MediaTarget]string{}
	if replacement.CommonLogo != nil {
		if s.mediaSvc == nil {
			return &domain.AppError{Err: domain.ErrInternal, Message: "media service is not configured"}
		}
		logoKeys, err := s.uploadLocationLogo(ctx, l.ID, *replacement.CommonLogo)
		if err != nil {
			return err
		}
		newLogoKeys = logoKeys
		l.CommonLogo = logoKeys[locationMediaTarget(l.ID, "common_logo")]
		l.CommonLargeLogo = logoKeys[locationMediaTarget(l.ID, "common_large_logo")]
		l.CommonMediumLogo = logoKeys[locationMediaTarget(l.ID, "common_medium_logo")]
		l.CommonSmallLogo = logoKeys[locationMediaTarget(l.ID, "common_small_logo")]
	} else if replacement.ClearCommonLogo {
		l.CommonLogo = ""
		l.CommonLargeLogo = ""
		l.CommonMediumLogo = ""
		l.CommonSmallLogo = ""
	}

	if err := s.repo.Update(ctx, l); err != nil {
		deleteLocationMediaKeys(ctx, s.mediaSvc, newLogoKeys)
		return err
	}
	if err := s.repo.SetCategories(ctx, l.ID, categoryIDs); err != nil {
		deleteLocationMediaKeys(ctx, s.mediaSvc, newLogoKeys)
		return err
	}
	if replacement.ReplaceImages {
		if err := s.replaceLocationImages(ctx, l, replacement.KeepImageIDs, replacement.Uploads); err != nil {
			return err
		}
	}

	if replacement.CommonLogo != nil || replacement.ClearCommonLogo {
		deleteLocationMediaKeys(ctx, s.mediaSvc, oldLogoKeys)
	}
	return nil
}

func locationMediaTarget(locationID uuid.UUID, field string) MediaTarget {
	return MediaTarget{Entity: "locations", RecordID: locationID, Field: field}
}

func locationLogoKeys(l *domain.Location) map[MediaTarget]string {
	if l == nil {
		return nil
	}
	return map[MediaTarget]string{
		locationMediaTarget(l.ID, "common_logo"):        l.CommonLogo,
		locationMediaTarget(l.ID, "common_large_logo"):  l.CommonLargeLogo,
		locationMediaTarget(l.ID, "common_medium_logo"): l.CommonMediumLogo,
		locationMediaTarget(l.ID, "common_small_logo"):  l.CommonSmallLogo,
	}
}

func (s *locationService) uploadLocationLogo(ctx context.Context, locationID uuid.UUID, uploads LocationLogoUploads) (map[MediaTarget]string, error) {
	ordered := []struct {
		field  string
		upload MediaUpload
	}{
		{field: "common_logo", upload: uploads.Original},
		{field: "common_large_logo", upload: uploads.Large},
		{field: "common_medium_logo", upload: uploads.Medium},
		{field: "common_small_logo", upload: uploads.Small},
	}
	keys := map[MediaTarget]string{}
	for _, item := range ordered {
		target := locationMediaTarget(locationID, item.field)
		key, err := s.mediaSvc.Upload(ctx, target, item.upload)
		if err != nil {
			deleteLocationMediaKeys(ctx, s.mediaSvc, keys)
			return nil, err
		}
		keys[target] = key
	}
	return keys, nil
}

func (s *locationService) replaceLocationImages(ctx context.Context, l *domain.Location, keepIDs []uuid.UUID, uploads []MediaUpload) error {
	kept, dropped := partitionLocationImages(l.Images, keepIDs)
	_ = kept
	for _, img := range dropped {
		if err := s.repo.DeleteImage(ctx, img.ID); err != nil {
			return err
		}
	}
	if len(uploads) == 0 {
		deleteLocationImageMedia(ctx, s.mediaSvc, dropped)
		return nil
	}
	if s.mediaSvc == nil {
		return &domain.AppError{Err: domain.ErrInternal, Message: "media service is not configured"}
	}
	target := locationMediaTarget(l.ID, "pictures")
	newKeys := []string{}
	for _, upload := range uploads {
		key, err := s.mediaSvc.Upload(ctx, target, upload)
		if err != nil {
			deleteLocationMediaKeyList(ctx, s.mediaSvc, target, newKeys)
			return err
		}
		newKeys = append(newKeys, key)
		if err := s.repo.CreateImage(ctx, &domain.LocationImage{LocationID: l.ID, Original: key}); err != nil {
			deleteLocationMediaKeyList(ctx, s.mediaSvc, target, newKeys)
			return err
		}
	}
	deleteLocationImageMedia(ctx, s.mediaSvc, dropped)
	return nil
}

func partitionLocationImages(images []*domain.LocationImage, keepIDs []uuid.UUID) ([]*domain.LocationImage, []*domain.LocationImage) {
	keep := make(map[uuid.UUID]struct{}, len(keepIDs))
	for _, id := range keepIDs {
		keep[id] = struct{}{}
	}
	kept := []*domain.LocationImage{}
	dropped := []*domain.LocationImage{}
	for _, img := range images {
		if img == nil {
			continue
		}
		if _, ok := keep[img.ID]; ok {
			kept = append(kept, img)
			continue
		}
		dropped = append(dropped, img)
	}
	return kept, dropped
}

func deleteLocationMediaKeys(ctx context.Context, mediaSvc MediaService, keys map[MediaTarget]string) {
	if mediaSvc == nil {
		return
	}
	for target, key := range keys {
		if key != "" {
			_ = mediaSvc.DeleteOwned(ctx, target, key)
		}
	}
}

func deleteLocationMediaKeyList(ctx context.Context, mediaSvc MediaService, target MediaTarget, keys []string) {
	if mediaSvc == nil {
		return
	}
	for _, key := range keys {
		if key != "" {
			_ = mediaSvc.DeleteOwned(ctx, target, key)
		}
	}
}

func deleteLocationImageMedia(ctx context.Context, mediaSvc MediaService, images []*domain.LocationImage) {
	keys := make([]string, 0, len(images))
	for _, img := range images {
		if img != nil && img.Original != "" {
			keys = append(keys, img.Original)
		}
	}
	deleteLocationMediaKeyList(ctx, mediaSvc, locationMediaTarget(uuidFromLocationImages(images), "pictures"), keys)
}

func uuidFromLocationImages(images []*domain.LocationImage) uuid.UUID {
	for _, img := range images {
		if img != nil {
			return img.LocationID
		}
	}
	return uuid.Nil
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

func (s *locationService) SetTopWithMedia(ctx context.Context, id uuid.UUID, isTop bool, sortIndex *int, upload *MediaUpload) error {
	loc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if upload != nil {
		if s.mediaSvc == nil {
			return &domain.AppError{Err: domain.ErrInternal, Message: "media service is not configured"}
		}
		target := MediaTarget{Entity: "locations", RecordID: id, Field: "top_logo"}
		key, err := s.mediaSvc.Upload(ctx, target, *upload)
		if err != nil {
			return err
		}
		oldKey := loc.TopLogo
		loc.TopLogo = key
		loc.TopLogoType = upload.ContentType
		if err := s.repo.Update(ctx, loc); err != nil {
			_ = s.mediaSvc.DeleteOwned(ctx, target, key)
			return err
		}
		if oldKey != "" && oldKey != key {
			_ = s.mediaSvc.DeleteOwned(ctx, target, oldKey)
		}
	}
	if err := s.repo.SetTopLocation(ctx, id, isTop, sortIndex); err != nil {
		return err
	}
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
