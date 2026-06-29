package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type AdvertisementService interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Advertisement, int, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Advertisement, error)
	Create(ctx context.Context, a *domain.Advertisement) error
	CreateWithMedia(ctx context.Context, a *domain.Advertisement, upload *MediaUpload) error
	Update(ctx context.Context, a *domain.Advertisement) error
	UpdateWithMedia(ctx context.Context, a *domain.Advertisement, replacement AdvertisementMediaReplacement) error
	Delete(ctx context.Context, id uuid.UUID) error
	Publish(ctx context.Context, id uuid.UUID) error
}

type AdvertisementMediaReplacement struct {
	OldKey string
	Upload *MediaUpload
}

type adService struct {
	repo     repository.AdvertisementRepository
	mediaSvc MediaService
}

func NewAdvertisementService(repo repository.AdvertisementRepository, mediaSvc ...MediaService) AdvertisementService {
	var media MediaService
	if len(mediaSvc) > 0 {
		media = mediaSvc[0]
	}
	return &adService{repo: repo, mediaSvc: media}
}

func (s *adService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Advertisement, int, error) {
	return s.repo.List(ctx, venueID, p)
}

func (s *adService) Get(ctx context.Context, id uuid.UUID) (*domain.Advertisement, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *adService) Create(ctx context.Context, a *domain.Advertisement) error {
	return s.repo.Create(ctx, a)
}

func (s *adService) CreateWithMedia(ctx context.Context, a *domain.Advertisement, upload *MediaUpload) error {
	if upload == nil {
		return s.repo.Create(ctx, a)
	}
	if s.mediaSvc == nil {
		return &domain.AppError{Err: domain.ErrInternal, Message: "media service is not configured"}
	}
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	target := advertisementImageMediaTarget(a.ID)
	key, err := s.mediaSvc.Upload(ctx, target, *upload)
	if err != nil {
		return err
	}
	a.ContentImage = &key
	if err := s.repo.Create(ctx, a); err != nil {
		_ = s.mediaSvc.DeleteOwned(ctx, target, key)
		return err
	}
	return nil
}

func (s *adService) Update(ctx context.Context, a *domain.Advertisement) error {
	return s.repo.Update(ctx, a)
}

func (s *adService) UpdateWithMedia(ctx context.Context, a *domain.Advertisement, replacement AdvertisementMediaReplacement) error {
	if replacement.Upload == nil {
		return s.repo.Update(ctx, a)
	}
	if s.mediaSvc == nil {
		return &domain.AppError{Err: domain.ErrInternal, Message: "media service is not configured"}
	}
	target := advertisementImageMediaTarget(a.ID)
	key, err := s.mediaSvc.Upload(ctx, target, *replacement.Upload)
	if err != nil {
		return err
	}
	a.ContentImage = &key
	if err := s.repo.Update(ctx, a); err != nil {
		_ = s.mediaSvc.DeleteOwned(ctx, target, key)
		return err
	}
	if replacement.OldKey != "" && replacement.OldKey != key {
		_ = s.mediaSvc.DeleteOwned(ctx, target, replacement.OldKey)
	}
	return nil
}

func (s *adService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *adService) Publish(ctx context.Context, id uuid.UUID) error {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now()
	a.Status = "published"
	a.PublishedAt = &now
	return s.repo.Update(ctx, a)
}

func advertisementImageMediaTarget(adID uuid.UUID) MediaTarget {
	return MediaTarget{Entity: "ads", RecordID: adID, Field: "content_image"}
}
