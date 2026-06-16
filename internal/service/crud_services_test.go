package service_test

// Tests for thin CRUD delegation services: Ad, Coupon, Video.
// All three follow the same List/Get/Create/Update/Delete pattern delegating directly to their repo.

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

// ── AdvertisementService ──────────────────────────────────────────────────────

func TestAdService_List(t *testing.T) {
	repo := &mocks.AdvertisementRepository{}
	svc := service.NewAdvertisementService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	p := domain.Pagination{Page: 1, PageSize: 20}
	expected := []*domain.Advertisement{{Type: "banner"}, {Type: "popup"}}

	repo.On("List", ctx, venueID, p).Return(expected, 2, nil)

	ads, total, err := svc.List(ctx, venueID, p)
	require.NoError(t, err)
	assert.Len(t, ads, 2)
	assert.Equal(t, 2, total)
	repo.AssertExpectations(t)
}

func TestAdService_Get(t *testing.T) {
	repo := &mocks.AdvertisementRepository{}
	svc := service.NewAdvertisementService(repo)

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.Advertisement{ID: id, Type: "banner", Status: "draft"}

	repo.On("FindByID", ctx, id).Return(expected, nil)

	a, err := svc.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, a.ID)
	repo.AssertExpectations(t)
}

func TestAdService_Create(t *testing.T) {
	repo := &mocks.AdvertisementRepository{}
	svc := service.NewAdvertisementService(repo)

	ctx := context.Background()
	a := &domain.Advertisement{Type: "banner", Status: "draft"}

	repo.On("Create", ctx, a).Return(nil)

	err := svc.Create(ctx, a)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAdService_Update(t *testing.T) {
	repo := &mocks.AdvertisementRepository{}
	svc := service.NewAdvertisementService(repo)

	ctx := context.Background()
	a := &domain.Advertisement{ID: uuid.New(), Status: "published"}

	repo.On("Update", ctx, a).Return(nil)

	err := svc.Update(ctx, a)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAdService_Delete(t *testing.T) {
	repo := &mocks.AdvertisementRepository{}
	svc := service.NewAdvertisementService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

// ── CouponService ─────────────────────────────────────────────────────────────

func TestCouponService_List(t *testing.T) {
	repo := &mocks.CouponRepository{}
	svc := service.NewCouponService(repo, &mocks.CouponUserRepository{})

	ctx := context.Background()
	venueID := uuid.New()
	p := domain.Pagination{Page: 1, PageSize: 10}
	expected := []*domain.Coupon{{CouponCode: "SAVE10"}}

	repo.On("List", ctx, venueID, p).Return(expected, 1, nil)

	coupons, total, err := svc.List(ctx, venueID, p)
	require.NoError(t, err)
	assert.Len(t, coupons, 1)
	assert.Equal(t, 1, total)
	repo.AssertExpectations(t)
}

func TestCouponService_Get(t *testing.T) {
	repo := &mocks.CouponRepository{}
	svc := service.NewCouponService(repo, &mocks.CouponUserRepository{})

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.Coupon{ID: id, CouponCode: "SAVE10"}

	repo.On("FindByID", ctx, id).Return(expected, nil)

	c, err := svc.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "SAVE10", c.CouponCode)
	repo.AssertExpectations(t)
}

func TestCouponService_Create(t *testing.T) {
	repo := &mocks.CouponRepository{}
	svc := service.NewCouponService(repo, &mocks.CouponUserRepository{})

	ctx := context.Background()
	c := &domain.Coupon{CouponCode: "NEW20"}

	repo.On("Create", ctx, c).Return(nil)

	err := svc.Create(ctx, c)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCouponService_Delete(t *testing.T) {
	repo := &mocks.CouponRepository{}
	svc := service.NewCouponService(repo, &mocks.CouponUserRepository{})

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

// ── VideoService ──────────────────────────────────────────────────────────────

func TestVideoService_List(t *testing.T) {
	repo := &mocks.VideoRepository{}
	svc := service.NewVideoService(repo)

	ctx := context.Background()
	venueID := uuid.New()
	p := domain.Pagination{Page: 1, PageSize: 10}
	expected := []*domain.Video{{Title: "Tour Video"}}

	repo.On("List", ctx, venueID, p).Return(expected, 1, nil)

	videos, total, err := svc.List(ctx, venueID, p)
	require.NoError(t, err)
	assert.Len(t, videos, 1)
	assert.Equal(t, 1, total)
	repo.AssertExpectations(t)
}

func TestVideoService_Get(t *testing.T) {
	repo := &mocks.VideoRepository{}
	svc := service.NewVideoService(repo)

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.Video{ID: id, Title: "Welcome"}

	repo.On("FindByID", ctx, id).Return(expected, nil)

	v, err := svc.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Welcome", v.Title)
	repo.AssertExpectations(t)
}

func TestVideoService_Create(t *testing.T) {
	repo := &mocks.VideoRepository{}
	svc := service.NewVideoService(repo)

	ctx := context.Background()
	v := &domain.Video{Title: "New Video"}

	repo.On("Create", ctx, v).Return(nil)

	err := svc.Create(ctx, v)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestVideoService_Delete(t *testing.T) {
	repo := &mocks.VideoRepository{}
	svc := service.NewVideoService(repo)

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}
