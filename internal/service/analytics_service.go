package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

// AnalyticsService handles event log ingestion and search query tracking.
type AnalyticsService interface {
	LogEvent(ctx context.Context, e *domain.EventLog) error
	ListEventLogs(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.EventLog, int, error)
	TrackSearch(ctx context.Context, venueID uuid.UUID, term string) error
	ListSearchQueries(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.SearchQuery, int, error)
}

type analyticsService struct {
	eventRepo  repository.EventLogRepository
	searchRepo repository.SearchQueryRepository
	cache      *redis.Client
}

// NewAnalyticsService creates an AnalyticsService.
func NewAnalyticsService(
	eventRepo repository.EventLogRepository,
	searchRepo repository.SearchQueryRepository,
	cache *redis.Client,
) AnalyticsService {
	return &analyticsService{
		eventRepo:  eventRepo,
		searchRepo: searchRepo,
		cache:      cache,
	}
}

func (s *analyticsService) LogEvent(ctx context.Context, e *domain.EventLog) error {
	return s.eventRepo.Create(ctx, e)
}

func (s *analyticsService) ListEventLogs(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.EventLog, int, error) {
	return s.eventRepo.ListByVenue(ctx, venueID, p)
}

func (s *analyticsService) TrackSearch(ctx context.Context, venueID uuid.UUID, term string) error {
	return s.searchRepo.Upsert(ctx, venueID, term)
}

func (s *analyticsService) ListSearchQueries(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.SearchQuery, int, error) {
	return s.searchRepo.List(ctx, venueID, p)
}
