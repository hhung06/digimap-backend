package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
)

type notificationServiceSpy struct {
	created []*domain.Notification
	err     error
}

func (s *notificationServiceSpy) List(context.Context, uuid.UUID, domain.Pagination) ([]*domain.Notification, int64, error) {
	return nil, 0, nil
}

func (s *notificationServiceSpy) ListPushTypes(context.Context, uuid.UUID) ([]map[string]any, error) {
	return nil, nil
}

func (s *notificationServiceSpy) Get(context.Context, uuid.UUID) (*domain.Notification, error) {
	return nil, nil
}

func (s *notificationServiceSpy) Create(_ context.Context, n *domain.Notification) error {
	s.created = append(s.created, n)
	return s.err
}

func (s *notificationServiceSpy) Update(context.Context, *domain.Notification) error {
	return nil
}

func (s *notificationServiceSpy) Delete(context.Context, uuid.UUID) error {
	return nil
}

func (s *notificationServiceSpy) Send(context.Context, uuid.UUID) error {
	return nil
}

func (s *notificationServiceSpy) SendDueScheduled(context.Context, time.Time) (int, error) {
	return 0, nil
}

func TestWebhookHandler_JMAPushNotification_UsesNotificationService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	venueID := uuid.New()
	venueRepo := &mocks.VenueRepository{}
	notificationSvc := &notificationServiceSpy{}
	h := newWebhookHandler(venueRepo, nil, nil, nil, notificationSvc)

	venueRepo.
		On("FindByPrivateKey", mock.Anything, "secret").
		Return(&domain.Venue{ID: venueID, PrivateKey: "secret"}, nil)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/jma/push", bytes.NewBufferString(`{
		"expo_id": "foodex",
		"device_tokens": ["token-1"],
		"title": "Hello",
		"content": "World",
		"data": {"screen": "home"}
	}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("TOKEN", "secret")

	h.JMAPushNotification(c)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Len(t, notificationSvc.created, 1)
	require.Equal(t, domain.NotifTypeImmediate, notificationSvc.created[0].SendType)
	require.Equal(t, []byte(`["token-1"]`), []byte(notificationSvc.created[0].DeviceTokens))
	venueRepo.AssertExpectations(t)
}
