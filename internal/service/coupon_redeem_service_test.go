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
	svc := service.NewCouponService(repo)

	ctx := context.Background()
	couponID := uuid.New()
	appUserID := uuid.New()
	coupon := &domain.Coupon{ID: couponID, CouponCode: "SAVE10", RedeemedAt: nil}

	repo.On("FindByID", ctx, couponID).Return(coupon, nil)
	repo.On("Redeem", ctx, couponID, appUserID).Return(nil)

	err := svc.RedeemCoupon(ctx, couponID, appUserID)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCouponService_RedeemCoupon_AlreadyRedeemed(t *testing.T) {
	repo := &mocks.CouponRepository{}
	svc := service.NewCouponService(repo)

	ctx := context.Background()
	couponID := uuid.New()
	appUserID := uuid.New()
	redeemedAt := time.Now()
	coupon := &domain.Coupon{ID: couponID, CouponCode: "SAVE10", RedeemedAt: &redeemedAt}

	repo.On("FindByID", ctx, couponID).Return(coupon, nil)

	err := svc.RedeemCoupon(ctx, couponID, appUserID)
	require.Error(t, err)

	appErr, ok := err.(*domain.AppError)
	require.True(t, ok)
	assert.Equal(t, domain.ErrConflict, appErr.Err)
	repo.AssertExpectations(t)
}
