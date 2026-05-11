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

type MemoService interface {
	List(ctx context.Context, venueID uuid.UUID) ([]*domain.Location, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Location, error)
	Create(ctx context.Context, l *domain.Location) error
	Update(ctx context.Context, l *domain.Location) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type memoService struct {
	repo        repository.LocationRepository
	venueRepo   repository.VenueRepository
	storer      storage.Storer
	invalidator cdn.Invalidator
	appVersions *AppVersionService
	env         string
}

func NewMemoService(
	repo repository.LocationRepository,
	venueRepo repository.VenueRepository,
	storer storage.Storer,
	invalidator cdn.Invalidator,
	appVersions *AppVersionService,
	env string,
) MemoService {
	return &memoService{
		repo:        repo,
		venueRepo:   venueRepo,
		storer:      storer,
		invalidator: invalidator,
		appVersions: appVersions,
		env:         env,
	}
}

func (s *memoService) List(ctx context.Context, venueID uuid.UUID) ([]*domain.Location, error) {
	return s.repo.ListMemos(ctx, venueID)
}

func (s *memoService) Get(ctx context.Context, id uuid.UUID) (*domain.Location, error) {
	return s.repo.FindMemoByID(ctx, id)
}

func (s *memoService) Create(ctx context.Context, l *domain.Location) error {
	l.CommonLocationType = domain.LocationTypeMemo
	if err := s.repo.Create(ctx, l); err != nil {
		return err
	}
	go s.publishMemos(context.Background(), l.VenueID)
	return nil
}

func (s *memoService) Update(ctx context.Context, l *domain.Location) error {
	l.CommonLocationType = domain.LocationTypeMemo
	if err := s.repo.Update(ctx, l); err != nil {
		return err
	}
	go s.publishMemos(context.Background(), l.VenueID)
	return nil
}

func (s *memoService) Delete(ctx context.Context, id uuid.UUID) error {
	memo, err := s.repo.FindMemoByID(ctx, id)
	if err != nil {
		return err
	}
	venueID := memo.VenueID
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	go s.publishMemos(context.Background(), venueID)
	return nil
}

// publishMemos compresses, AES-encrypts, uploads to S3, invalidates CloudFront, and bumps
// the force-sync version — matching Django's LocationMemoViewSet.upload_memo_to_s3
// (indoormap-backend/api/locations/views.py:538).
func (s *memoService) publishMemos(ctx context.Context, venueID uuid.UUID) {
	venue, err := s.venueRepo.FindByID(ctx, venueID)
	if err != nil {
		fmt.Printf("memo publish: load venue %s: %v\n", venueID, err)
		return
	}

	memos, err := s.repo.ListMemos(ctx, venueID)
	if err != nil {
		fmt.Printf("memo publish: list memos %s: %v\n", venueID, err)
		return
	}

	ciphertext, err := crypto.EncryptBytes(venue.PublicKey, memos)
	if err != nil {
		fmt.Printf("memo publish: encrypt %s: %v\n", venueID, err)
		return
	}

	key := storage.MemoKey(s.env, venueID)
	meta := map[string]string{"encrypted": "AES", "compressed": "gzip"}
	if err := s.storer.PutEncrypted(ctx, key, []byte(ciphertext), meta); err != nil {
		fmt.Printf("memo publish: put encrypted %s: %v\n", venueID, err)
		return
	}

	if _, err := s.invalidator.Invalidate(ctx, []string{"/" + key}); err != nil {
		fmt.Printf("memo publish: cf invalidation %s: %v\n", venueID, err)
	}

	if _, err := s.appVersions.Bump(ctx, venueID); err != nil {
		fmt.Printf("memo publish: version bump %s: %v\n", venueID, err)
	}
}
