package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

// AnalyticsService handles event log ingestion and search query tracking.
type AnalyticsService interface {
	LogEvent(ctx context.Context, e *domain.EventLog) error
	ListEventLogs(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.EventLog, int, error)
	TrackSearch(ctx context.Context, venueID uuid.UUID, term, origin string) error
	ListSearchQueries(ctx context.Context, venueID uuid.UUID, filter repository.SearchQueryFilter, p domain.Pagination) ([]*domain.SearchQuery, int, error)
}

type analyticsService struct {
	eventRepo  repository.EventLogRepository
	searchRepo repository.SearchQueryRepository
	venueRepo  repository.VenueRepository
	cache      *redis.Client
}

// NewAnalyticsService creates an AnalyticsService.
func NewAnalyticsService(
	eventRepo repository.EventLogRepository,
	searchRepo repository.SearchQueryRepository,
	venueRepo repository.VenueRepository,
	cache *redis.Client,
) AnalyticsService {
	return &analyticsService{
		eventRepo:  eventRepo,
		searchRepo: searchRepo,
		venueRepo:  venueRepo,
		cache:      cache,
	}
}

// deriveAppID maps a venue's external_id prefix to the known app identifier.
func deriveAppID(externalID string) string {
	lower := strings.ToLower(externalID)
	switch {
	case strings.HasPrefix(lower, "foodex"):
		return "foodex"
	case strings.HasPrefix(lower, "hcj"):
		return "hcj"
	default:
		return ""
	}
}

func (s *analyticsService) LogEvent(ctx context.Context, e *domain.EventLog) error {
	return s.eventRepo.Create(ctx, e)
}

func (s *analyticsService) ListEventLogs(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.EventLog, int, error) {
	return s.eventRepo.ListByVenue(ctx, venueID, p)
}

func (s *analyticsService) TrackSearch(ctx context.Context, venueID uuid.UUID, term, origin string) error {
	venue, err := s.venueRepo.FindByID(ctx, venueID)
	if err != nil {
		return err
	}
	return s.searchRepo.Upsert(ctx, venueID, term, origin, deriveAppID(venue.ExternalID))
}

func (s *analyticsService) ListSearchQueries(ctx context.Context, venueID uuid.UUID, filter repository.SearchQueryFilter, p domain.Pagination) ([]*domain.SearchQuery, int, error) {
	return s.searchRepo.List(ctx, venueID, filter, p)
}
