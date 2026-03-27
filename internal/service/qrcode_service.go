package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type QRCodeService interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.QRCode, int, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.QRCode, error)
	Create(ctx context.Context, q *domain.QRCode) error
	Update(ctx context.Context, q *domain.QRCode) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type qrcodeService struct {
	repo repository.QRCodeRepository
}

func NewQRCodeService(repo repository.QRCodeRepository) QRCodeService {
	return &qrcodeService{repo: repo}
}

func (s *qrcodeService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.QRCode, int, error) {
	return s.repo.List(ctx, venueID, p)
}

func (s *qrcodeService) Get(ctx context.Context, id uuid.UUID) (*domain.QRCode, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *qrcodeService) Create(ctx context.Context, q *domain.QRCode) error {
	return s.repo.Create(ctx, q)
}

func (s *qrcodeService) Update(ctx context.Context, q *domain.QRCode) error {
	return s.repo.Update(ctx, q)
}

func (s *qrcodeService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
