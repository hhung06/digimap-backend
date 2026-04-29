package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)


type CouponService interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Coupon, int, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Coupon, error)
	Create(ctx context.Context, c *domain.Coupon) error
	Update(ctx context.Context, c *domain.Coupon) error
	Delete(ctx context.Context, id uuid.UUID) error
	RedeemCoupon(ctx context.Context, couponID uuid.UUID, appUserID uuid.UUID) error
}

type couponService struct {
	repo repository.CouponRepository
}

func NewCouponService(repo repository.CouponRepository) CouponService {
	return &couponService{repo: repo}
}

func (s *couponService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Coupon, int, error) {
	return s.repo.List(ctx, venueID, p)
}

func (s *couponService) Get(ctx context.Context, id uuid.UUID) (*domain.Coupon, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *couponService) Create(ctx context.Context, c *domain.Coupon) error {
	return s.repo.Create(ctx, c)
}

func (s *couponService) Update(ctx context.Context, c *domain.Coupon) error {
	return s.repo.Update(ctx, c)
}

func (s *couponService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *couponService) RedeemCoupon(ctx context.Context, couponID uuid.UUID, appUserID uuid.UUID) error {
	c, err := s.repo.FindByID(ctx, couponID)
	if err != nil {
		return err
	}
	if c.RedeemedAt != nil {
		return domain.NewConflict("coupon already redeemed")
	}
	return s.repo.Redeem(ctx, couponID, appUserID)
}
