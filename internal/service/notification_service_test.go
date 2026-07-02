package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/firebase"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func newTestNotificationService(repo *mocks.NotificationRepository, pusher *mocks.MockPusher) service.NotificationService {
	return service.NewNotificationService(repo, pusher)
}

func strPtr(s string) *string { return &s }

// ── Create ────────────────────────────────────────────────────────────────────

func TestNotificationService_Create_DraftWithTopic(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	svc := newTestNotificationService(repo, nil)

	ctx := context.Background()
	venueID := uuid.New()
	n := &domain.Notification{
		Title:    strPtr("Flash Sale"),
		Content:  strPtr("50% off"),
		VenueID:  &venueID,
		Topic:    strPtr("venue-promo"),
		Kind:     domain.NotifKindNormal,
		SendType: domain.NotifTypeDraft,
	}

	repo.On("Create", ctx, n).Return(nil)

	err := svc.Create(ctx, n)
	require.NoError(t, err)
	assert.Equal(t, domain.NotifStatusUnsent, n.Status)
	assert.Equal(t, domain.NotifSendPending, n.SendStatus)
	repo.AssertExpectations(t)
}

func TestNotificationService_Create_ImmediateAutoDispatches(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	pusher := &mocks.MockPusher{}
	svc := newTestNotificationService(repo, pusher)

	ctx := context.Background()
	venueID := uuid.New()
	id := uuid.New()
	n := &domain.Notification{
		ID:       id,
		Title:    strPtr("Flash Sale"),
		Content:  strPtr("50% off"),
		VenueID:  &venueID,
		Topic:    strPtr("venue-promo"),
		Kind:     domain.NotifKindNormal,
		SendType: domain.NotifTypeImmediate,
	}
	stored := &domain.Notification{
		ID:         id,
		Title:      strPtr("Flash Sale"),
		Content:    strPtr("50% off"),
		Topic:      strPtr("venue-promo"),
		SendStatus: domain.NotifSendPending,
	}

	repo.On("Create", ctx, n).Return(nil)
	repo.On("FindByID", ctx, id).Return(stored, nil)
	pusher.On("Send", ctx, firebase.Message{Topic: "venue-promo", Title: "Flash Sale", Body: "50% off"}).Return("msg-id", nil)
	repo.On("MarkSent", ctx, id, mockAny).Return(nil)

	err := svc.Create(ctx, n)
	require.NoError(t, err)
	repo.AssertExpectations(t)
	pusher.AssertExpectations(t)
}

func TestNotificationService_Create_ScheduledRequiresScheduledAt(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	svc := newTestNotificationService(repo, nil)

	ctx := context.Background()
	venueID := uuid.New()
	n := &domain.Notification{
		Title:    strPtr("Promo"),
		Content:  strPtr("Deal"),
		VenueID:  &venueID,
		Topic:    strPtr("venue-promo"),
		Kind:     domain.NotifKindNormal,
		SendType: domain.NotifTypeScheduled,
		// ScheduledAt intentionally missing
	}

	err := svc.Create(ctx, n)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrValidation))
	repo.AssertNotCalled(t, "Create")
}

func TestNotificationService_Create_ScheduledRequiresFutureScheduledAt(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	svc := newTestNotificationService(repo, nil)

	ctx := context.Background()
	past := time.Now().Add(-time.Hour)
	venueID := uuid.New()
	n := &domain.Notification{
		Title:       strPtr("Promo"),
		Content:     strPtr("Deal"),
		VenueID:     &venueID,
		Topic:       strPtr("venue-promo"),
		Kind:        domain.NotifKindNormal,
		SendType:    domain.NotifTypeScheduled,
		ScheduledAt: &past,
	}

	err := svc.Create(ctx, n)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrValidation))
	repo.AssertNotCalled(t, "Create")
}

func TestNotificationService_Create_RequiresDeliveryTarget(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	svc := newTestNotificationService(repo, nil)

	ctx := context.Background()
	venueID := uuid.New()
	n := &domain.Notification{
		Title:    strPtr("Promo"),
		Content:  strPtr("Deal"),
		VenueID:  &venueID,
		Kind:     domain.NotifKindNormal,
		SendType: domain.NotifTypeDraft,
		// No Topic, DeviceTokens, or SegmentFilters
	}

	err := svc.Create(ctx, n)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrValidation))
	repo.AssertNotCalled(t, "Create")
}

func TestNotificationService_Create_RejectsUnsupportedSegmentFilters(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	svc := newTestNotificationService(repo, nil)

	ctx := context.Background()
	venueID := uuid.New()
	n := &domain.Notification{
		Title:          strPtr("Promo"),
		Content:        strPtr("Deal"),
		VenueID:        &venueID,
		Kind:           domain.NotifKindNormal,
		SendType:       domain.NotifTypeDraft,
		SegmentFilters: []byte(`[{"foo":"bar"}]`),
	}

	err := svc.Create(ctx, n)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrValidation))
	repo.AssertNotCalled(t, "Create")
}

func TestNotificationService_Create_RejectsInvalidSendType(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	svc := newTestNotificationService(repo, nil)

	ctx := context.Background()
	venueID := uuid.New()
	n := &domain.Notification{
		Title:    strPtr("Promo"),
		Content:  strPtr("Deal"),
		VenueID:  &venueID,
		Topic:    strPtr("venue-promo"),
		Kind:     domain.NotifKindNormal,
		SendType: 99,
	}

	err := svc.Create(ctx, n)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrValidation))
	repo.AssertNotCalled(t, "Create")
}

func TestNotificationService_Create_RejectsInvalidKind(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	svc := newTestNotificationService(repo, nil)

	ctx := context.Background()
	venueID := uuid.New()
	n := &domain.Notification{
		Title:    strPtr("Promo"),
		Content:  strPtr("Deal"),
		VenueID:  &venueID,
		Topic:    strPtr("venue-promo"),
		Kind:     99,
		SendType: domain.NotifTypeDraft,
	}

	err := svc.Create(ctx, n)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrValidation))
	repo.AssertNotCalled(t, "Create")
}

func TestNotificationService_Create_ScheduledFuture(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	svc := newTestNotificationService(repo, nil)

	ctx := context.Background()
	future := time.Now().Add(time.Hour)
	venueID := uuid.New()
	n := &domain.Notification{
		Title:       strPtr("Promo"),
		Content:     strPtr("Deal"),
		VenueID:     &venueID,
		Topic:       strPtr("venue-promo"),
		Kind:        domain.NotifKindNormal,
		SendType:    domain.NotifTypeScheduled,
		ScheduledAt: &future,
	}

	repo.On("Create", ctx, n).Return(nil)

	err := svc.Create(ctx, n)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

// ── Get ───────────────────────────────────────────────────────────────────────

func TestNotificationService_Get_Success(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	svc := newTestNotificationService(repo, nil)

	ctx := context.Background()
	id := uuid.New()
	expected := &domain.Notification{ID: id, Title: strPtr("Flash Sale")}

	repo.On("FindByID", ctx, id).Return(expected, nil)

	got, err := svc.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	repo.AssertExpectations(t)
}

func TestNotificationService_Get_NotFound(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	svc := newTestNotificationService(repo, nil)

	ctx := context.Background()
	id := uuid.New()

	repo.On("FindByID", ctx, id).Return((*domain.Notification)(nil), domain.NewNotFound("notification not found"))

	_, err := svc.Get(ctx, id)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
	repo.AssertExpectations(t)
}

// ── List ──────────────────────────────────────────────────────────────────────

func TestNotificationService_List_Success(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	svc := newTestNotificationService(repo, nil)

	ctx := context.Background()
	venueID := uuid.New()
	p := domain.Pagination{Page: 1, PageSize: 20}
	expected := []*domain.Notification{{Title: strPtr("Alert")}}

	repo.On("List", ctx, venueID, p).Return(expected, int64(1), nil)

	got, total, err := svc.List(ctx, venueID, p)
	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}

func TestNotificationService_ListPushTypes_ComposesFoodexOptions(t *testing.T) {
	t.Setenv("BFACE_BASIC_AUTHENTICATION", "")
	t.Setenv("FOODEX_X_API_KEY", "")

	repo := &mocks.NotificationRepository{}
	venueRepo := &mocks.VenueRepository{}
	svc := service.NewNotificationService(repo, nil, venueRepo)

	ctx := context.Background()
	venueID := uuid.New()
	venueRepo.On("FindByID", ctx, venueID).Return(&domain.Venue{
		ID:         venueID,
		ExternalID: "foodex_2026",
	}, nil)

	pushTypes, err := svc.ListPushTypes(ctx, venueID)
	require.NoError(t, err)

	var industryOptions []map[string]any
	var seminarOptions []map[string]any
	for _, pushType := range pushTypes {
		switch pushType["key"] {
		case "topic_visitor_industry":
			industryOptions, _ = pushType["options"].([]map[string]any)
		case "seminar_registered":
			seminarOptions, _ = pushType["options"].([]map[string]any)
		}
	}

	require.NotEmpty(t, industryOptions)
	assert.Equal(t, "商社・卸 ➞ 卸・問屋", industryOptions[0]["label"])
	assert.Empty(t, seminarOptions)
	venueRepo.AssertExpectations(t)
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestNotificationService_Delete_Success(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	svc := newTestNotificationService(repo, nil)

	ctx := context.Background()
	id := uuid.New()

	repo.On("Delete", ctx, id).Return(nil)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

// ── Send ──────────────────────────────────────────────────────────────────────

func TestNotificationService_Send_ByTopic_Success(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	pusher := &mocks.MockPusher{}
	svc := newTestNotificationService(repo, pusher)

	ctx := context.Background()
	id := uuid.New()
	n := &domain.Notification{
		ID:         id,
		Title:      strPtr("Flash Sale"),
		Content:    strPtr("50% off today!"),
		Topic:      strPtr("venue-123"),
		SendStatus: domain.NotifSendPending,
	}

	repo.On("FindByID", ctx, id).Return(n, nil)
	pusher.On("Send", ctx, firebase.Message{Topic: "venue-123", Title: "Flash Sale", Body: "50% off today!"}).Return("msg-id", nil)
	repo.On("MarkSent", ctx, id, mockAny).Return(nil)

	err := svc.Send(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
	pusher.AssertExpectations(t)
}

func TestNotificationService_Send_AlreadySent(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	pusher := &mocks.MockPusher{}
	svc := newTestNotificationService(repo, pusher)

	ctx := context.Background()
	id := uuid.New()
	n := &domain.Notification{
		ID:         id,
		SendStatus: domain.NotifSendSuccess,
	}

	repo.On("FindByID", ctx, id).Return(n, nil)

	err := svc.Send(ctx, id)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrConflict))
	// Pusher must NOT be called
	pusher.AssertNotCalled(t, "Send")
	pusher.AssertNotCalled(t, "SendMulticast")
	repo.AssertExpectations(t)
}

func TestNotificationService_Send_FCMError(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	pusher := &mocks.MockPusher{}
	svc := newTestNotificationService(repo, pusher)

	ctx := context.Background()
	id := uuid.New()
	n := &domain.Notification{
		ID:         id,
		Title:      strPtr("Flash Sale"),
		Content:    strPtr("50% off!"),
		Topic:      strPtr("venue-123"),
		SendStatus: domain.NotifSendPending,
	}

	repo.On("FindByID", ctx, id).Return(n, nil)
	pusher.On("Send", ctx, mockAny).Return("", errors.New("fcm unavailable"))
	repo.On("MarkFailed", ctx, id, mockAny).Return(nil)

	err := svc.Send(ctx, id)
	require.Error(t, err)
	repo.AssertExpectations(t)
	pusher.AssertExpectations(t)
}

func TestNotificationService_Send_BySegmentFiltersTopics_Success(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	pusher := &mocks.MockPusher{}
	svc := newTestNotificationService(repo, pusher)

	ctx := context.Background()
	id := uuid.New()
	n := &domain.Notification{
		ID:             id,
		Title:          strPtr("Survey"),
		Content:        strPtr("Please answer"),
		SendStatus:     domain.NotifSendPending,
		SegmentFilters: []byte(`[{"key":"visitors","type":"text","value":"vip"},{"key":"all_users","type":null}]`),
	}

	repo.On("FindByID", ctx, id).Return(n, nil)
	pusher.On("Send", ctx, firebase.Message{Topic: "visitors_vip", Title: "Survey", Body: "Please answer"}).Return("msg-1", nil)
	pusher.On("Send", ctx, firebase.Message{Topic: "all_users", Title: "Survey", Body: "Please answer"}).Return("msg-2", nil)
	repo.On("MarkSent", ctx, id, mockAny).Return(nil)

	err := svc.Send(ctx, id)
	require.NoError(t, err)
	repo.AssertExpectations(t)
	pusher.AssertExpectations(t)
}

func TestNotificationService_SendDueScheduled_SendsReadyNotifications(t *testing.T) {
	repo := &mocks.NotificationRepository{}
	pusher := &mocks.MockPusher{}
	svc := newTestNotificationService(repo, pusher)

	ctx := context.Background()
	now := time.Now()
	id := uuid.New()
	n := &domain.Notification{
		ID:         id,
		Title:      strPtr("Scheduled"),
		Content:    strPtr("Due now"),
		Topic:      strPtr("venue-123"),
		SendType:   domain.NotifTypeScheduled,
		SendStatus: domain.NotifSendPending,
	}

	repo.On("ListDueScheduled", ctx, now).Return([]*domain.Notification{{ID: id}}, nil)
	repo.On("FindByID", ctx, id).Return(n, nil)
	pusher.On("Send", ctx, firebase.Message{Topic: "venue-123", Title: "Scheduled", Body: "Due now"}).Return("msg-id", nil)
	repo.On("MarkSent", ctx, id, mockAny).Return(nil)

	count, err := svc.SendDueScheduled(ctx, now)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
	repo.AssertExpectations(t)
	pusher.AssertExpectations(t)
}
