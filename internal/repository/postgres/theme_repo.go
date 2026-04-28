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

const themeSelectCols = `id, venue_id, name, primary_color, secondary_color, created_at, updated_at`

func (r *themeRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Theme, error) {
	q := `SELECT ` + themeSelectCols + ` FROM themes WHERE id = $1 AND deleted_at IS NULL`
	t, err := scanTheme(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("theme not found")
	}
	return t, err
}

func (r *themeRepo) List(ctx context.Context, venueID uuid.UUID) ([]*domain.Theme, error) {
	q := `SELECT ` + themeSelectCols + ` FROM themes WHERE venue_id = $1 AND deleted_at IS NULL ORDER BY name`
	rows, err := r.pool.Query(ctx, q, venueID)
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

func (r *themeRepo) Create(ctx context.Context, t *domain.Theme) error {
	if t.ID == uuid.Nil {
		t.ID = newID()
	}
	const q = `
		INSERT INTO themes (id, venue_id, name, primary_color, secondary_color)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, q, t.ID, t.VenueID, t.Name, t.PrimaryColor, t.SecondaryColor).Scan(&t.CreatedAt, &t.UpdatedAt)
}

func (r *themeRepo) Update(ctx context.Context, t *domain.Theme) error {
	const q = `
		UPDATE themes SET name = $2, primary_color = $3, secondary_color = $4
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`
	err := r.pool.QueryRow(ctx, q, t.ID, t.Name, t.PrimaryColor, t.SecondaryColor).Scan(&t.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("theme not found")
	}
	return err
}

func (r *themeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "themes", id.String())
}

func scanTheme(row scanner) (*domain.Theme, error) {
	var t domain.Theme
	err := row.Scan(&t.ID, &t.VenueID, &t.Name, &t.PrimaryColor, &t.SecondaryColor, &t.CreatedAt, &t.UpdatedAt)
	return &t, err
}
