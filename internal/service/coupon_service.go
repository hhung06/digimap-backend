package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type CouponService interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Coupon, int, error)
	// ListForUser returns coupons assigned to a specific app user, with IsUsed populated.
	ListForUser(ctx context.Context, venueID, userID uuid.UUID, p domain.Pagination) ([]*domain.Coupon, int, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Coupon, error)
	Create(ctx context.Context, c *domain.Coupon) error
	Update(ctx context.Context, c *domain.Coupon) error
	Delete(ctx context.Context, id uuid.UUID) error
	// AssignToUser creates a coupon_users row linking a coupon to an app user.
	AssignToUser(ctx context.Context, couponID, userID uuid.UUID, deviceID string) error
	// RedeemCoupon marks the user's coupon_users row as used.
	RedeemCoupon(ctx context.Context, couponID uuid.UUID, appUserID uuid.UUID) error
}

type couponService struct {
	repo     repository.CouponRepository
	userRepo repository.CouponUserRepository
}

func NewCouponService(repo repository.CouponRepository, userRepo repository.CouponUserRepository) CouponService {
	return &couponService{repo: repo, userRepo: userRepo}
}

func (s *couponService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Coupon, int, error) {
	return s.repo.List(ctx, venueID, p)
}

func (s *couponService) ListForUser(ctx context.Context, venueID, userID uuid.UUID, p domain.Pagination) ([]*domain.Coupon, int, error) {
	return s.repo.ListByUser(ctx, venueID, userID, p)
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

func (s *couponService) AssignToUser(ctx context.Context, couponID, userID uuid.UUID, deviceID string) error {
	cu := &domain.CouponUser{
		CouponID: couponID,
		UserID:   &userID,
		DeviceID: deviceID,
		IsUsed:   false,
	}
	return s.userRepo.Create(ctx, cu)
}

func (s *couponService) RedeemCoupon(ctx context.Context, couponID uuid.UUID, appUserID uuid.UUID) error {
	cu, err := s.userRepo.FindByCouponAndUser(ctx, couponID, appUserID)
	if err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			return err
		}
		c, findErr := s.repo.FindByID(ctx, couponID)
		if findErr != nil {
			return findErr
		}
		if c.RedeemedAt != nil {
			return domain.NewConflict("coupon already redeemed")
		}
		return s.repo.Redeem(ctx, couponID, appUserID)
	}
	if cu.IsUsed {
		return domain.NewConflict("coupon already redeemed")
	}
	return s.userRepo.MarkUsed(ctx, cu.ID)
}
