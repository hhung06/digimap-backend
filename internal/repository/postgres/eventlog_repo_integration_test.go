//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/hhung06/digimap-backend/internal/domain"
	postgresrepo "github.com/hhung06/digimap-backend/internal/repository/postgres"
)

// setupTestDB spins up a PostgreSQL container and returns a connection pool.
// The container is automatically terminated when the test finishes.
func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections"),
		),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = pgContainer.Terminate(ctx) })

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	// Apply minimal schema for event_logs
	_, err = pool.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS "pgcrypto";

		CREATE TABLE IF NOT EXISTS event_logs (
			id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
			venue_id    UUID,
			name        TEXT        NOT NULL,
			params      JSONB       NOT NULL DEFAULT '{}',
			device_id   TEXT,
			user_id     TEXT,
			user_agent  TEXT,
			ip_address  TEXT,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
		) PARTITION BY RANGE (created_at);

		CREATE TABLE event_logs_default PARTITION OF event_logs DEFAULT;
	`)
	require.NoError(t, err)

	return pool
}

func TestEventLogRepository_Create(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgresrepo.NewEventLogRepository(pool)

	ctx := context.Background()
	venueID := uuid.New()
	e := &domain.EventLog{
		VenueID:   &venueID,
		Name:      "page_view",
		UserAgent: "TestAgent/1.0",
		IPAddress: "127.0.0.1",
	}

	err := repo.Create(ctx, e)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, e.ID)
	assert.False(t, e.CreatedAt.IsZero())
}

func TestEventLogRepository_ListByVenue(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgresrepo.NewEventLogRepository(pool)

	ctx := context.Background()
	venueID := uuid.New()
	otherVenueID := uuid.New()

	// Insert events for two venues
	for i := 0; i < 3; i++ {
		e := &domain.EventLog{VenueID: &venueID, Name: "click"}
		require.NoError(t, repo.Create(ctx, e))
		time.Sleep(time.Millisecond) // ensure distinct created_at ordering
	}
	other := &domain.EventLog{VenueID: &otherVenueID, Name: "view"}
	require.NoError(t, repo.Create(ctx, other))

	p := domain.Pagination{Page: 1, PageSize: 10}
	logs, total, err := repo.ListByVenue(ctx, venueID, p)

	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, logs, 3)
	for _, l := range logs {
		assert.Equal(t, venueID, *l.VenueID)
	}
}

func TestEventLogRepository_ListByVenue_Pagination(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgresrepo.NewEventLogRepository(pool)

	ctx := context.Background()
	venueID := uuid.New()

	for i := 0; i < 5; i++ {
		e := &domain.EventLog{VenueID: &venueID, Name: "event"}
		require.NoError(t, repo.Create(ctx, e))
	}

	p := domain.Pagination{Page: 1, PageSize: 2}
	logs, total, err := repo.ListByVenue(ctx, venueID, p)

	require.NoError(t, err)
	assert.Equal(t, 5, total) // total is always the full count
	assert.Len(t, logs, 2)    // page size respected
}
