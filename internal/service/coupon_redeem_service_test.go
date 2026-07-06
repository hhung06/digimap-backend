package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func TestCouponService_RedeemCoupon_Success(t *testing.T) {
	repo := &mocks.CouponRepository{}
	userRepo := &mocks.CouponUserRepository{}
	svc := service.NewCouponService(repo, userRepo)

	ctx := context.Background()
	venueID := uuid.New()
	couponID := uuid.New()
	appUserID := uuid.New()
	cuID := uuid.New()
	coupon := &domain.Coupon{ID: couponID, VenueID: &venueID}
	cu := &domain.CouponUser{ID: cuID, CouponID: couponID, IsUsed: false}

	repo.On("FindByID", ctx, couponID).Return(coupon, nil)
	userRepo.On("FindByCouponAndUser", ctx, couponID, appUserID).Return(cu, nil)
	userRepo.On("MarkUsed", ctx, cuID).Return(nil)

	err := svc.RedeemCoupon(ctx, venueID, couponID, appUserID)
	require.NoError(t, err)
	repo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestCouponService_RedeemCoupon_AlreadyRedeemed(t *testing.T) {
	repo := &mocks.CouponRepository{}
	userRepo := &mocks.CouponUserRepository{}
	svc := service.NewCouponService(repo, userRepo)

	ctx := context.Background()
	venueID := uuid.New()
	couponID := uuid.New()
	appUserID := uuid.New()
	coupon := &domain.Coupon{ID: couponID, VenueID: &venueID}
	cu := &domain.CouponUser{ID: uuid.New(), CouponID: couponID, IsUsed: true}

	repo.On("FindByID", ctx, couponID).Return(coupon, nil)
	userRepo.On("FindByCouponAndUser", ctx, couponID, appUserID).Return(cu, nil)

	err := svc.RedeemCoupon(ctx, venueID, couponID, appUserID)
	require.Error(t, err)

	appErr, ok := err.(*domain.AppError)
	require.True(t, ok)
	assert.Equal(t, domain.ErrConflict, appErr.Err)
	repo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestCouponService_RedeemCoupon_FallsBackToCouponLevelRedemption(t *testing.T) {
	repo := &mocks.CouponRepository{}
	userRepo := &mocks.CouponUserRepository{}
	svc := service.NewCouponService(repo, userRepo)

	ctx := context.Background()
	venueID := uuid.New()
	couponID := uuid.New()
	appUserID := uuid.New()
	coupon := &domain.Coupon{ID: couponID, VenueID: &venueID, CouponCode: strPtr("SAVE10")}

	repo.On("FindByID", ctx, couponID).Return(coupon, nil)
	userRepo.On("FindByCouponAndUser", ctx, couponID, appUserID).Return((*domain.CouponUser)(nil), domain.NewNotFound("coupon not assigned to this user"))
	repo.On("Redeem", ctx, couponID, appUserID).Return(nil)

	err := svc.RedeemCoupon(ctx, venueID, couponID, appUserID)
	require.NoError(t, err)
	repo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestCouponService_RedeemCoupon_FallbackDetectsAlreadyRedeemedCoupon(t *testing.T) {
	repo := &mocks.CouponRepository{}
	userRepo := &mocks.CouponUserRepository{}
	svc := service.NewCouponService(repo, userRepo)

	ctx := context.Background()
	venueID := uuid.New()
	couponID := uuid.New()
	appUserID := uuid.New()
	redeemedAt := time.Now()
	coupon := &domain.Coupon{ID: couponID, VenueID: &venueID, CouponCode: strPtr("SAVE10"), RedeemedAt: &redeemedAt}

	repo.On("FindByID", ctx, couponID).Return(coupon, nil)
	userRepo.On("FindByCouponAndUser", ctx, couponID, appUserID).Return((*domain.CouponUser)(nil), domain.NewNotFound("coupon not assigned to this user"))

	err := svc.RedeemCoupon(ctx, venueID, couponID, appUserID)
	require.Error(t, err)

	appErr, ok := err.(*domain.AppError)
	require.True(t, ok)
	assert.Equal(t, domain.ErrConflict, appErr.Err)
	repo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestCouponService_RedeemCoupon_RejectsCrossVenueCoupon(t *testing.T) {
	repo := &mocks.CouponRepository{}
	userRepo := &mocks.CouponUserRepository{}
	svc := service.NewCouponService(repo, userRepo)

	ctx := context.Background()
	callerVenueID := uuid.New()
	ownerVenueID := uuid.New()
	couponID := uuid.New()
	appUserID := uuid.New()
	coupon := &domain.Coupon{ID: couponID, VenueID: &ownerVenueID}

	repo.On("FindByID", ctx, couponID).Return(coupon, nil)

	err := svc.RedeemCoupon(ctx, callerVenueID, couponID, appUserID)
	require.Error(t, err)

	appErr, ok := err.(*domain.AppError)
	require.True(t, ok)
	assert.Equal(t, domain.ErrNotFound, appErr.Err)
	userRepo.AssertNotCalled(t, "FindByCouponAndUser")
	repo.AssertExpectations(t)
}
