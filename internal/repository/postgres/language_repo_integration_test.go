//go:build integration

package postgres_test

import (
	"context"
	"testing"

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

func setupLanguageTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("testcontainers unavailable: %v", r)
		}
	}()

	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections"),
		),
	)
	if err != nil {
		t.Skipf("testcontainers unavailable: %v", err)
	}
	t.Cleanup(func() { _ = pgContainer.Terminate(ctx) })

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	_, err = pool.Exec(ctx, `
		CREATE TABLE languages (
			id UUID PRIMARY KEY,
			venue_id UUID NOT NULL,
			code VARCHAR(10) NOT NULL,
			name VARCHAR(100) NOT NULL,
			is_default BOOLEAN NOT NULL DEFAULT false,
			enabled BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ,
			UNIQUE (venue_id, code)
		);
	`)
	require.NoError(t, err)

	return pool
}

func TestLanguageRepositoryReplaceAllInsertsUpdatesAndDeletes(t *testing.T) {
	pool := setupLanguageTestDB(t)
	repo := postgresrepo.NewLanguageRepository(pool)
	ctx := context.Background()
	venueID := uuid.New()

	require.NoError(t, repo.Create(ctx, &domain.Language{VenueID: venueID, Code: "en", Name: "English", IsDefault: true}))
	require.NoError(t, repo.Create(ctx, &domain.Language{VenueID: venueID, Code: "ja", Name: "Japanese"}))

	desired := []*domain.Language{
		{VenueID: venueID, Code: "en", Name: "English (US)", IsDefault: false},
		{VenueID: venueID, Code: "fr", Name: "French", IsDefault: true},
	}

	got, err := repo.ReplaceAll(ctx, venueID, desired)

	require.NoError(t, err)
	require.Len(t, got, 2)
	byCode := make(map[string]*domain.Language, len(got))
	for _, l := range got {
		byCode[l.Code] = l
	}
	require.Contains(t, byCode, "en")
	require.Contains(t, byCode, "fr")
	assert.NotContains(t, byCode, "ja")
	assert.Equal(t, "English (US)", byCode["en"].Name)
	assert.False(t, byCode["en"].IsDefault)
	assert.True(t, byCode["fr"].IsDefault)
}

func TestLanguageRepositoryReplaceAllRollsBackOnConstraintViolation(t *testing.T) {
	pool := setupLanguageTestDB(t)
	repo := postgresrepo.NewLanguageRepository(pool)
	ctx := context.Background()
	venueID := uuid.New()
	require.NoError(t, repo.Create(ctx, &domain.Language{VenueID: venueID, Code: "en", Name: "English", IsDefault: true}))

	desired := []*domain.Language{
		{VenueID: venueID, Code: "fr", Name: "French", IsDefault: true},
		{VenueID: venueID, Code: "fr", Name: "French (dup)", IsDefault: false},
	}

	_, err := repo.ReplaceAll(ctx, venueID, desired)

	require.Error(t, err)
	remaining, err := repo.List(ctx, venueID)
	require.NoError(t, err)
	require.Len(t, remaining, 1)
	assert.Equal(t, "en", remaining[0].Code)
}
