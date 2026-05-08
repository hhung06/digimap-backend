package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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
	SendDueScheduled(ctx context.Context, now time.Time) (int, error)
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

func (s *notificationService) SendDueScheduled(ctx context.Context, now time.Time) (int, error) {
	items, err := s.repo.ListDueScheduled(ctx, now)
	if err != nil {
		return 0, err
	}
	sent := 0
	for _, item := range items {
		if err := s.Send(ctx, item.ID); err != nil {
			return sent, err
		}
		sent++
	}
	return sent, nil
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
	} else {
		sendErr = s.sendBySegmentFilters(ctx, n)
	}

	if sendErr != nil {
		errInfos, _ := json.Marshal([]string{fmt.Sprintf("FCM error: %v", sendErr)})
		_ = s.repo.MarkFailed(ctx, id, errInfos)
		return fmt.Errorf("fcm send: %w", sendErr)
	}

	return s.repo.MarkSent(ctx, id, time.Now())
}

func (s *notificationService) sendBySegmentFilters(ctx context.Context, n *domain.Notification) error {
	topics, err := topicsFromSegmentFilters(n.SegmentFilters)
	if err != nil {
		return err
	}
	if len(topics) == 0 {
		return domain.NewValidation(map[string]string{"segment_filters": "notification has no supported delivery target"})
	}
	for _, topic := range topics {
		if _, err := s.pusher.Send(ctx, firebase.Message{
			Topic: topic,
			Title: n.Title,
			Body:  n.Content,
		}); err != nil {
			return err
		}
	}
	return nil
}

func topicsFromSegmentFilters(raw json.RawMessage) ([]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var list []map[string]any
	if err := json.Unmarshal(raw, &list); err == nil {
		return segmentListToTopics(list), nil
	}

	var single map[string]any
	if err := json.Unmarshal(raw, &single); err == nil {
		return segmentListToTopics([]map[string]any{single}), nil
	}

	return nil, domain.NewValidation(map[string]string{"segment_filters": "invalid segment_filters payload"})
}

func segmentListToTopics(segments []map[string]any) []string {
	topics := make([]string, 0, len(segments))
	seen := make(map[string]struct{}, len(segments))
	for _, segment := range segments {
		topic := topicFromSegment(segment)
		if topic == "" {
			continue
		}
		if _, ok := seen[topic]; ok {
			continue
		}
		seen[topic] = struct{}{}
		topics = append(topics, topic)
	}
	return topics
}

func topicFromSegment(segment map[string]any) string {
	key, _ := segment["key"].(string)
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}

	segmentType, _ := segment["type"].(string)
	value := strings.TrimSpace(fmt.Sprint(segment["value"]))
	if value == "<nil>" {
		value = ""
	}
	if value == "" || (segmentType != "text" && segmentType != "number" && segmentType != "date") {
		return key
	}
	return key + "_" + value
}
