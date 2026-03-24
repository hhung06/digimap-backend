package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
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
	return s.repo.FindByID(ctx, id)
}

func (s *locationCategoryService) List(ctx context.Context, venueID uuid.UUID) ([]*domain.LocationCategory, error) {
	return s.repo.List(ctx, venueID)
}

func (s *locationCategoryService) Create(ctx context.Context, c *domain.LocationCategory) error {
	if c.Source == "" {
		c.Source = "internal"
	}
	return s.repo.Create(ctx, c)
}

func (s *locationCategoryService) Update(ctx context.Context, c *domain.LocationCategory) error {
	return s.repo.Update(ctx, c)
}

func (s *locationCategoryService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// ── Amenity service ───────────────────────────────────────────────────────────

type AmenityService interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Amenity, error)
	List(ctx context.Context, p domain.Pagination) ([]*domain.Amenity, int64, error)
	ListByVenue(ctx context.Context, venueID uuid.UUID) ([]*domain.Amenity, error)
	Create(ctx context.Context, a *domain.Amenity) error
	Update(ctx context.Context, a *domain.Amenity) error
	Delete(ctx context.Context, id uuid.UUID) error
	LinkToVenue(ctx context.Context, venueID, amenityID uuid.UUID) error
	UnlinkFromVenue(ctx context.Context, venueID, amenityID uuid.UUID) error
}

type amenityService struct {
	repo repository.AmenityRepository
}

func NewAmenityService(repo repository.AmenityRepository) AmenityService {
	return &amenityService{repo: repo}
}

func (s *amenityService) Get(ctx context.Context, id uuid.UUID) (*domain.Amenity, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *amenityService) List(ctx context.Context, p domain.Pagination) ([]*domain.Amenity, int64, error) {
	p.Normalize()
	return s.repo.List(ctx, p)
}

func (s *amenityService) ListByVenue(ctx context.Context, venueID uuid.UUID) ([]*domain.Amenity, error) {
	return s.repo.ListByVenue(ctx, venueID)
}

func (s *amenityService) Create(ctx context.Context, a *domain.Amenity) error {
	return s.repo.Create(ctx, a)
}

func (s *amenityService) Update(ctx context.Context, a *domain.Amenity) error {
	return s.repo.Update(ctx, a)
}

func (s *amenityService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *amenityService) LinkToVenue(ctx context.Context, venueID, amenityID uuid.UUID) error {
	va := &domain.VenueAmenity{VenueID: venueID, AmenityID: amenityID}
	return s.repo.LinkToVenue(ctx, va)
}

func (s *amenityService) UnlinkFromVenue(ctx context.Context, venueID, amenityID uuid.UUID) error {
	return s.repo.UnlinkFromVenue(ctx, venueID, amenityID)
}

// ── Location service ──────────────────────────────────────────────────────────

type LocationService interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Location, error)
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Location, int64, error)
	Create(ctx context.Context, l *domain.Location) error
	Update(ctx context.Context, l *domain.Location, categoryIDs []uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	Duplicate(ctx context.Context, id uuid.UUID) (*domain.Location, error)
	SetTop(ctx context.Context, id uuid.UUID, isTop bool, sortIndex *int) error

	// Images
	ListImages(ctx context.Context, locationID uuid.UUID) ([]*domain.LocationImage, error)
	CreateImage(ctx context.Context, img *domain.LocationImage) error
	DeleteImage(ctx context.Context, locationID, imageID uuid.UUID) error

	// Promotions
	GetPromotion(ctx context.Context, id uuid.UUID) (*domain.Promotion, error)
	ListPromotions(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Promotion, int64, error)
	CreatePromotion(ctx context.Context, p *domain.Promotion) error
	UpdatePromotion(ctx context.Context, p *domain.Promotion) error
	DeletePromotion(ctx context.Context, id uuid.UUID) error
}

type locationService struct {
	repo repository.LocationRepository
}

func NewLocationService(repo repository.LocationRepository) LocationService {
	return &locationService{repo: repo}
}

func (s *locationService) Get(ctx context.Context, id uuid.UUID) (*domain.Location, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *locationService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Location, int64, error) {
	p.Normalize()
	return s.repo.List(ctx, venueID, p)
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
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.repo.SetTopLocation(ctx, id, isTop, sortIndex)
}

func (s *locationService) ListImages(ctx context.Context, locationID uuid.UUID) ([]*domain.LocationImage, error) {
	return s.repo.ListImages(ctx, locationID)
}

func (s *locationService) CreateImage(ctx context.Context, img *domain.LocationImage) error {
	return s.repo.CreateImage(ctx, img)
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

func (s *locationService) GetPromotion(ctx context.Context, id uuid.UUID) (*domain.Promotion, error) {
	return s.repo.FindPromotionByID(ctx, id)
}

func (s *locationService) ListPromotions(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Promotion, int64, error) {
	p.Normalize()
	return s.repo.ListPromotions(ctx, venueID, p)
}

func (s *locationService) CreatePromotion(ctx context.Context, p *domain.Promotion) error {
	if p.DisplayType == "" {
		p.DisplayType = "random"
	}
	return s.repo.CreatePromotion(ctx, p)
}

func (s *locationService) UpdatePromotion(ctx context.Context, p *domain.Promotion) error {
	return s.repo.UpdatePromotion(ctx, p)
}

func (s *locationService) DeletePromotion(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeletePromotion(ctx, id)
}
