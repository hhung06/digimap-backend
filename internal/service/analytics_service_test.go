package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
	"github.com/hhung06/digimap-backend/internal/repository/mocks"
	"github.com/hhung06/digimap-backend/internal/service"
)

func newTestAnalyticsService(eventRepo *mocks.EventLogRepository, searchRepo *mocks.SearchQueryRepository, venueRepo *mocks.VenueRepository) service.AnalyticsService {
	// nil redis client — analytics service only uses it for future caching, not in current logic
	return service.NewAnalyticsService(eventRepo, searchRepo, venueRepo, nil)
}

// ── LogEvent ──────────────────────────────────────────────────────────────────

func TestAnalyticsService_LogEvent_Success(t *testing.T) {
	eventRepo := &mocks.EventLogRepository{}
	searchRepo := &mocks.SearchQueryRepository{}
	svc := newTestAnalyticsService(eventRepo, searchRepo, &mocks.VenueRepository{})

	ctx := context.Background()
	venueID := uuid.New()
	e := &domain.EventLog{
		VenueID:   &venueID,
		Name:      "page_view",
		UserAgent: "Mozilla/5.0",
		IPAddress: "1.2.3.4",
	}

	eventRepo.On("Create", ctx, e).Return(nil)

	err := svc.LogEvent(ctx, e)
	require.NoError(t, err)
	eventRepo.AssertExpectations(t)
}

func TestAnalyticsService_LogEvent_RepoError(t *testing.T) {
	eventRepo := &mocks.EventLogRepository{}
	searchRepo := &mocks.SearchQueryRepository{}
	svc := newTestAnalyticsService(eventRepo, searchRepo, &mocks.VenueRepository{})

	ctx := context.Background()
	venueID := uuid.New()
	e := &domain.EventLog{VenueID: &venueID, Name: "click"}

	eventRepo.On("Create", ctx, e).Return(domain.NewNotFound("venue not found"))

	err := svc.LogEvent(ctx, e)
	require.Error(t, err)
	eventRepo.AssertExpectations(t)
}

// ── ListEventLogs ─────────────────────────────────────────────────────────────

func TestAnalyticsService_ListEventLogs(t *testing.T) {
	eventRepo := &mocks.EventLogRepository{}
	searchRepo := &mocks.SearchQueryRepository{}
	svc := newTestAnalyticsService(eventRepo, searchRepo, &mocks.VenueRepository{})

	ctx := context.Background()
	venueID := uuid.New()
	p := domain.Pagination{Page: 1, PageSize: 20}

	expected := []*domain.EventLog{
		{Name: "page_view"},
		{Name: "search"},
	}

	eventRepo.On("ListByVenue", ctx, venueID, p).Return(expected, 2, nil)

	logs, total, err := svc.ListEventLogs(ctx, venueID, p)
	require.NoError(t, err)
	assert.Len(t, logs, 2)
	assert.Equal(t, 2, total)
	eventRepo.AssertExpectations(t)
}

// ── TrackSearch ───────────────────────────────────────────────────────────────

func TestAnalyticsService_TrackSearch_Success(t *testing.T) {
	eventRepo := &mocks.EventLogRepository{}
	searchRepo := &mocks.SearchQueryRepository{}
	venueRepo := &mocks.VenueRepository{}
	svc := newTestAnalyticsService(eventRepo, searchRepo, venueRepo)

	ctx := context.Background()
	venueID := uuid.New()
	venue := &domain.Venue{ID: venueID, ExternalID: "foodex_mar_2026"}

	venueRepo.On("FindByID", ctx, venueID).Return(venue, nil)
	searchRepo.On("Upsert", ctx, venueID, "coffee shop", "product", "foodex").Return(nil)

	err := svc.TrackSearch(ctx, venueID, "coffee shop", "product")
	require.NoError(t, err)
	venueRepo.AssertExpectations(t)
	searchRepo.AssertExpectations(t)
}

func TestAnalyticsService_TrackSearch_EmptyTerm(t *testing.T) {
	eventRepo := &mocks.EventLogRepository{}
	searchRepo := &mocks.SearchQueryRepository{}
	venueRepo := &mocks.VenueRepository{}
	svc := newTestAnalyticsService(eventRepo, searchRepo, venueRepo)

	ctx := context.Background()
	venueID := uuid.New()
	venue := &domain.Venue{ID: venueID, ExternalID: ""}

	venueRepo.On("FindByID", ctx, venueID).Return(venue, nil)
	// Service delegates directly to repo — empty term is repo's concern
	searchRepo.On("Upsert", ctx, venueID, "", "", "").Return(nil)

	err := svc.TrackSearch(ctx, venueID, "", "")
	require.NoError(t, err)
	venueRepo.AssertExpectations(t)
	searchRepo.AssertExpectations(t)
}

// ── ListSearchQueries ─────────────────────────────────────────────────────────

func TestAnalyticsService_ListSearchQueries(t *testing.T) {
	eventRepo := &mocks.EventLogRepository{}
	searchRepo := &mocks.SearchQueryRepository{}
	svc := newTestAnalyticsService(eventRepo, searchRepo, &mocks.VenueRepository{})

	ctx := context.Background()
	venueID := uuid.New()
	p := domain.Pagination{Page: 1, PageSize: 10}
	filter := repository.SearchQueryFilter{}

	venueIDCopy := venueID
	expected := []*domain.SearchQuery{
		{VenueID: &venueIDCopy, SearchTerm: "coffee", SearchCount: 5},
		{VenueID: &venueIDCopy, SearchTerm: "restroom", SearchCount: 3},
	}

	searchRepo.On("List", ctx, venueID, filter, p).Return(expected, 2, nil)

	queries, total, err := svc.ListSearchQueries(ctx, venueID, filter, p)
	require.NoError(t, err)
	assert.Len(t, queries, 2)
	assert.Equal(t, 2, total)
	assert.Equal(t, "coffee", queries[0].SearchTerm)
	searchRepo.AssertExpectations(t)
}
