package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/firebase"
	"github.com/hhung06/digimap-backend/internal/repository"
)

// NotificationService manages notification CRUD and FCM dispatch.
type NotificationService interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Notification, int64, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Notification, error)
	Create(ctx context.Context, n *domain.Notification) error
	Update(ctx context.Context, n *domain.Notification) error
	Delete(ctx context.Context, id uuid.UUID) error
	Send(ctx context.Context, id uuid.UUID) error
}

type notificationService struct {
	repo   repository.NotificationRepository
	pusher firebase.Pusher
}

// NewNotificationService creates a NotificationService.
func NewNotificationService(repo repository.NotificationRepository, pusher firebase.Pusher) NotificationService {
	return &notificationService{repo: repo, pusher: pusher}
}

func (s *notificationService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Notification, int64, error) {
	return s.repo.List(ctx, venueID, p)
}

func (s *notificationService) Get(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *notificationService) Create(ctx context.Context, n *domain.Notification) error {
	n.Status = domain.NotifStatusUnsent
	n.SendStatus = domain.NotifSendPending
	return s.repo.Create(ctx, n)
}

func (s *notificationService) Update(ctx context.Context, n *domain.Notification) error {
	return s.repo.Update(ctx, n)
}

func (s *notificationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// Send dispatches the notification immediately via FCM.
func (s *notificationService) Send(ctx context.Context, id uuid.UUID) error {
	n, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if n.SendStatus == domain.NotifSendSuccess {
		return domain.NewConflict("notification already sent")
	}

	var sendErr error
	if n.Topic != "" {
		_, sendErr = s.pusher.Send(ctx, firebase.Message{
			Topic: n.Topic,
			Title: n.Title,
			Body:  n.Content,
		})
	} else if len(n.DeviceTokens) > 0 {
		var tokens []string
		if jsonErr := json.Unmarshal(n.DeviceTokens, &tokens); jsonErr == nil && len(tokens) > 0 {
			_, _, sendErr = s.pusher.SendMulticast(ctx, tokens, n.Title, n.Content, nil)
		}
	}

	if sendErr != nil {
		errInfos, _ := json.Marshal([]string{fmt.Sprintf("FCM error: %v", sendErr)})
		_ = s.repo.MarkFailed(ctx, id, errInfos)
		return fmt.Errorf("fcm send: %w", sendErr)
	}

	return s.repo.MarkSent(ctx, id, time.Now())
}
