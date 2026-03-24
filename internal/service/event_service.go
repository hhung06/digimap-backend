package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

// EventService manages event tags, event types, events, and images.
type EventService interface {
	// Tags (global)
	ListTags(ctx context.Context) ([]*domain.EventTag, error)
	CreateTag(ctx context.Context, t *domain.EventTag) error
	UpdateTag(ctx context.Context, t *domain.EventTag) error
	DeleteTag(ctx context.Context, id uuid.UUID) error

	// Event types (per venue)
	ListEventTypes(ctx context.Context, venueID uuid.UUID) ([]*domain.EventType, error)
	CreateEventType(ctx context.Context, t *domain.EventType) error
	UpdateEventType(ctx context.Context, t *domain.EventType) error
	DeleteEventType(ctx context.Context, id uuid.UUID) error

	// Events
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Event, int64, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	Create(ctx context.Context, e *domain.Event, tagIDs []uuid.UUID, locationIDs []uuid.UUID) error
	Update(ctx context.Context, e *domain.Event, tagIDs []uuid.UUID, locationIDs []uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Images
	CreateImage(ctx context.Context, img *domain.EventImage) error
	DeleteImage(ctx context.Context, id uuid.UUID) error
}

type eventService struct {
	repo repository.EventRepository
}

// NewEventService creates an EventService backed by the given repository.
func NewEventService(repo repository.EventRepository) EventService {
	return &eventService{repo: repo}
}

func (s *eventService) ListTags(ctx context.Context) ([]*domain.EventTag, error) {
	return s.repo.ListTags(ctx)
}

func (s *eventService) CreateTag(ctx context.Context, t *domain.EventTag) error {
	return s.repo.CreateTag(ctx, t)
}

func (s *eventService) UpdateTag(ctx context.Context, t *domain.EventTag) error {
	return s.repo.UpdateTag(ctx, t)
}

func (s *eventService) DeleteTag(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTag(ctx, id)
}

func (s *eventService) ListEventTypes(ctx context.Context, venueID uuid.UUID) ([]*domain.EventType, error) {
	return s.repo.ListEventTypes(ctx, venueID)
}

func (s *eventService) CreateEventType(ctx context.Context, t *domain.EventType) error {
	return s.repo.CreateEventType(ctx, t)
}

func (s *eventService) UpdateEventType(ctx context.Context, t *domain.EventType) error {
	return s.repo.UpdateEventType(ctx, t)
}

func (s *eventService) DeleteEventType(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteEventType(ctx, id)
}

func (s *eventService) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Event, int64, error) {
	return s.repo.List(ctx, venueID, p)
}

func (s *eventService) Get(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *eventService) Create(ctx context.Context, e *domain.Event, tagIDs []uuid.UUID, locationIDs []uuid.UUID) error {
	if err := s.repo.Create(ctx, e); err != nil {
		return err
	}
	if len(tagIDs) > 0 {
		if err := s.repo.SetTags(ctx, e.ID, tagIDs); err != nil {
			return err
		}
	}
	if len(locationIDs) > 0 {
		if err := s.repo.SetLocations(ctx, e.ID, locationIDs); err != nil {
			return err
		}
	}
	return nil
}

func (s *eventService) Update(ctx context.Context, e *domain.Event, tagIDs []uuid.UUID, locationIDs []uuid.UUID) error {
	if err := s.repo.Update(ctx, e); err != nil {
		return err
	}
	if err := s.repo.SetTags(ctx, e.ID, tagIDs); err != nil {
		return err
	}
	return s.repo.SetLocations(ctx, e.ID, locationIDs)
}

func (s *eventService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *eventService) CreateImage(ctx context.Context, img *domain.EventImage) error {
	return s.repo.CreateImage(ctx, img)
}

func (s *eventService) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteImage(ctx, id)
}
