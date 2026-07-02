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

func setupArticleTestDB(t *testing.T) *pgxpool.Pool {
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
		CREATE TABLE articles (
			id UUID PRIMARY KEY,
			venue_id UUID,
			external_id VARCHAR(255),
			location_id UUID,
			placement VARCHAR(50) NOT NULL DEFAULT 'article',
			navigate VARCHAR(50),
			title VARCHAR(255) NOT NULL DEFAULT '',
			label VARCHAR(255),
			content TEXT,
			status VARCHAR(50) NOT NULL DEFAULT 'draft',
			published_at TIMESTAMPTZ,
			published_period_start DATE,
			published_period_end DATE,
			localization JSONB,
			related_products JSONB NOT NULL DEFAULT '[]',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ
		);

		CREATE TABLE locations (
			id UUID PRIMARY KEY,
			venue_id UUID,
			common_name TEXT NOT NULL DEFAULT '',
			deleted_at TIMESTAMPTZ
		);

		CREATE TABLE article_images (
			id UUID PRIMARY KEY,
			article_id UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
			image VARCHAR(1000) NOT NULL DEFAULT '' CHECK (image <> 'fail-key'),
			sort_order INT NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ
		);
	`)
	require.NoError(t, err)

	return pool
}

func TestArticleRepositoryUpdateWithImagesPreservesImagesWhenReplaceFalse(t *testing.T) {
	pool := setupArticleTestDB(t)
	repo := postgresrepo.NewArticleRepository(pool)
	ctx := context.Background()
	a := createArticleWithImages(t, ctx, repo, []string{"old-a", "old-b"})
	a.Title = "Scalar update"

	err := repo.UpdateWithImages(ctx, a, domain.ArticleMediaChange{Replace: false})

	require.NoError(t, err)
	assert.Equal(t, []string{"old-a", "old-b"}, activeArticleImageKeys(t, ctx, pool, a.ID))
}

func TestArticleRepositoryListLoadsImagesForResponseURLs(t *testing.T) {
	pool := setupArticleTestDB(t)
	repo := postgresrepo.NewArticleRepository(pool)
	ctx := context.Background()
	venueID := uuid.New()
	a := &domain.Article{VenueID: &venueID, Title: "Listed", Placement: "article", Status: "draft"}
	require.NoError(t, repo.Create(ctx, a))
	require.NoError(t, repo.CreateImage(ctx, &domain.ArticleImage{ArticleID: a.ID, Image: "article-cover", SortOrder: 0}))

	got, total, err := repo.List(ctx, venueID, domain.Pagination{Page: 1, PageSize: 20})

	require.NoError(t, err)
	assert.Equal(t, 1, total)
	require.Len(t, got, 1)
	require.Len(t, got[0].Images, 1)
	assert.Equal(t, "article-cover", got[0].Images[0].Image)
}

func TestArticleRepositoryListLoadsLocationForResponseName(t *testing.T) {
	pool := setupArticleTestDB(t)
	repo := postgresrepo.NewArticleRepository(pool)
	ctx := context.Background()
	venueID := uuid.New()
	locationID := uuid.New()
	_, err := pool.Exec(ctx,
		`INSERT INTO locations (id, venue_id, common_name) VALUES ($1, $2, $3)`,
		locationID, venueID, "Premium Lounge",
	)
	require.NoError(t, err)
	a := &domain.Article{
		VenueID:    &venueID,
		LocationID: &locationID,
		Title:      "Listed",
		Placement:  "article",
		Status:     "draft",
	}
	require.NoError(t, repo.Create(ctx, a))

	got, total, err := repo.List(ctx, venueID, domain.Pagination{Page: 1, PageSize: 20})

	require.NoError(t, err)
	assert.Equal(t, 1, total)
	require.Len(t, got, 1)
	require.NotNil(t, got[0].Location)
	assert.Equal(t, locationID, got[0].Location.ID)
	assert.Equal(t, "Premium Lounge", got[0].Location.Name)
}

func TestArticleRepositoryPersistsRelatedProducts(t *testing.T) {
	pool := setupArticleTestDB(t)
	repo := postgresrepo.NewArticleRepository(pool)
	ctx := context.Background()
	venueID := uuid.New()
	firstProductID := uuid.New()
	secondProductID := uuid.New()
	a := &domain.Article{
		VenueID:         &venueID,
		Title:           "Listed",
		Placement:       "article",
		Status:          "draft",
		RelatedProducts: []uuid.UUID{firstProductID, secondProductID},
	}
	require.NoError(t, repo.Create(ctx, a))

	got, err := repo.FindByID(ctx, a.ID)

	require.NoError(t, err)
	assert.Equal(t, []uuid.UUID{firstProductID, secondProductID}, got.RelatedProducts)
}

func TestArticleRepositoryUpdateWithImagesReplacesActiveImagesInKeyOrder(t *testing.T) {
	pool := setupArticleTestDB(t)
	repo := postgresrepo.NewArticleRepository(pool)
	ctx := context.Background()
	a := createArticleWithImages(t, ctx, repo, []string{"old-a", "old-b"})
	a.Title = "Replaced"

	err := repo.UpdateWithImages(ctx, a, domain.ArticleMediaChange{Replace: true, Keys: []string{"new-a", "new-b"}})

	require.NoError(t, err)
	assert.Equal(t, []string{"new-a", "new-b"}, activeArticleImageKeys(t, ctx, pool, a.ID))
	require.Len(t, a.Images, 2)
	assert.Equal(t, "new-a", a.Images[0].Image)
	assert.Equal(t, 0, a.Images[0].SortOrder)
	assert.Equal(t, "new-b", a.Images[1].Image)
	assert.Equal(t, 1, a.Images[1].SortOrder)
}

func TestArticleRepositoryUpdateWithImagesRollsBackScalarAndImagesOnInsertFailure(t *testing.T) {
	pool := setupArticleTestDB(t)
	repo := postgresrepo.NewArticleRepository(pool)
	ctx := context.Background()
	a := createArticleWithImages(t, ctx, repo, []string{"old-a"})
	a.Title = "Should roll back"

	err := repo.UpdateWithImages(ctx, a, domain.ArticleMediaChange{Replace: true, Keys: []string{"fail-key"}})

	require.Error(t, err)
	assert.Equal(t, "Original", articleTitle(t, ctx, pool, a.ID))
	assert.Equal(t, []string{"old-a"}, activeArticleImageKeys(t, ctx, pool, a.ID))
}

type articleRepoForIntegration interface {
	Create(context.Context, *domain.Article) error
	CreateImage(context.Context, *domain.ArticleImage) error
	FindByID(context.Context, uuid.UUID) (*domain.Article, error)
}

func createArticleWithImages(t *testing.T, ctx context.Context, repo articleRepoForIntegration, keys []string) *domain.Article {
	t.Helper()

	a := &domain.Article{ID: uuid.New(), Title: "Original", Placement: "article", Status: "draft"}
	require.NoError(t, repo.Create(ctx, a))
	for i, key := range keys {
		require.NoError(t, repo.CreateImage(ctx, &domain.ArticleImage{ArticleID: a.ID, Image: key, SortOrder: i}))
	}
	found, err := repo.FindByID(ctx, a.ID)
	require.NoError(t, err)
	return found
}

func activeArticleImageKeys(t *testing.T, ctx context.Context, pool *pgxpool.Pool, articleID uuid.UUID) []string {
	t.Helper()

	rows, err := pool.Query(ctx, `SELECT image FROM article_images WHERE article_id=$1 AND deleted_at IS NULL ORDER BY sort_order`, articleID)
	require.NoError(t, err)
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		require.NoError(t, rows.Scan(&key))
		keys = append(keys, key)
	}
	require.NoError(t, rows.Err())
	return keys
}

func articleTitle(t *testing.T, ctx context.Context, pool *pgxpool.Pool, articleID uuid.UUID) string {
	t.Helper()

	var title string
	require.NoError(t, pool.QueryRow(ctx, `SELECT title FROM articles WHERE id=$1`, articleID).Scan(&title))
	return title
}
