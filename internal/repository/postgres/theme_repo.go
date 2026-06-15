package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type themeRepo struct{ pool *pgxpool.Pool }

func NewThemeRepository(pool *pgxpool.Pool) repository.ThemeRepository {
	return &themeRepo{pool: pool}
}

const themeSelectCols = `id, venue_id, scope, name, data, storage_path, created_at, updated_at`

func prefixedThemeCols(alias string) string {
	return alias + `.id, ` + alias + `.venue_id, ` + alias + `.scope, ` + alias + `.name, ` +
		alias + `.data, ` + alias + `.storage_path, ` + alias + `.created_at, ` + alias + `.updated_at`
}

func (r *themeRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Theme, error) {
	q := `SELECT ` + themeSelectCols + ` FROM themes WHERE id = $1 AND deleted_at IS NULL`
	t, err := scanTheme(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("theme not found")
	}
	return t, err
}

func (r *themeRepo) FindVenueTheme(ctx context.Context, venueID uuid.UUID) (*domain.Theme, error) {
	q := `SELECT ` + prefixedThemeCols("t") + `
		FROM themes t
		INNER JOIN venues v ON v.theme_id = t.id
		WHERE v.id = $1 AND v.deleted_at IS NULL AND t.deleted_at IS NULL`
	t, err := scanTheme(r.pool.QueryRow(ctx, q, venueID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil // no theme selected — not an error
	}
	return t, err
}

func (r *themeRepo) ListGlobal(ctx context.Context) ([]*domain.Theme, error) {
	q := `SELECT ` + themeSelectCols + ` FROM themes WHERE scope = 'global' AND deleted_at IS NULL ORDER BY name`
	return queryThemes(r.pool, ctx, q)
}

func (r *themeRepo) List(ctx context.Context, venueID uuid.UUID) ([]*domain.Theme, error) {
	// Return custom themes for this venue first, then all global themes.
	q := `SELECT ` + themeSelectCols + `
		FROM themes
		WHERE deleted_at IS NULL
		  AND ((scope = 'custom' AND venue_id = $1) OR scope = 'global')
		ORDER BY CASE WHEN scope = 'custom' THEN 0 ELSE 1 END, name`
	return queryThemes(r.pool, ctx, q, venueID)
}

func (r *themeRepo) Create(ctx context.Context, t *domain.Theme) error {
	if t.ID == uuid.Nil {
		t.ID = newID()
	}
	const q = `
		INSERT INTO themes (id, venue_id, scope, name, data, storage_path)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, q, t.ID, t.VenueID, t.Scope, t.Name, t.Data, t.StoragePath).
		Scan(&t.CreatedAt, &t.UpdatedAt)
}

func (r *themeRepo) Update(ctx context.Context, t *domain.Theme) error {
	const q = `
		UPDATE themes SET name = $2, data = $3, storage_path = $4
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`
	err := r.pool.QueryRow(ctx, q, t.ID, t.Name, t.Data, t.StoragePath).Scan(&t.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("theme not found")
	}
	return err
}

func (r *themeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "themes", id.String())
}

func (r *themeRepo) IsUsedByVenues(ctx context.Context, id uuid.UUID) (bool, error) {
	const q = `SELECT EXISTS(SELECT 1 FROM venues WHERE theme_id = $1 AND deleted_at IS NULL)`
	var used bool
	err := r.pool.QueryRow(ctx, q, id).Scan(&used)
	return used, err
}

func queryThemes(pool *pgxpool.Pool, ctx context.Context, q string, args ...any) ([]*domain.Theme, error) {
	rows, err := pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.Theme
	for rows.Next() {
		t, err := scanTheme(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}

func scanTheme(row scanner) (*domain.Theme, error) {
	var t domain.Theme
	err := row.Scan(&t.ID, &t.VenueID, &t.Scope, &t.Name, &t.Data, &t.StoragePath, &t.CreatedAt, &t.UpdatedAt)
	return &t, err
}
