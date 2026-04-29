package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type languageRepository struct{ pool *pgxpool.Pool }

func NewLanguageRepository(pool *pgxpool.Pool) *languageRepository {
	return &languageRepository{pool: pool}
}

const languageSelectCols = `
    l.id, l.venue_id, l.code, l.name, l.is_default,
    l.created_at, l.updated_at`

func (r *languageRepository) List(ctx context.Context, venueID uuid.UUID) ([]*domain.Language, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+languageSelectCols+` FROM languages l
         WHERE l.venue_id=$1 AND l.deleted_at IS NULL ORDER BY l.created_at ASC`,
		venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Language
	for rows.Next() {
		l, err := scanLanguage(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *languageRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Language, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+languageSelectCols+` FROM languages l WHERE l.id=$1 AND l.deleted_at IS NULL`, id)
	return scanLanguage(row)
}

func (r *languageRepository) Create(ctx context.Context, l *domain.Language) error {
	l.ID = newID()
	return r.pool.QueryRow(ctx,
		`INSERT INTO languages (id, venue_id, code, name, is_default)
         VALUES ($1,$2,$3,$4,$5)
         RETURNING created_at, updated_at`,
		l.ID, l.VenueID, l.Code, l.Name, l.IsDefault,
	).Scan(&l.CreatedAt, &l.UpdatedAt)
}

func (r *languageRepository) Update(ctx context.Context, l *domain.Language) error {
	return r.pool.QueryRow(ctx,
		`UPDATE languages SET code=$2, name=$3, is_default=$4
         WHERE id=$1 AND deleted_at IS NULL RETURNING updated_at`,
		l.ID, l.Code, l.Name, l.IsDefault,
	).Scan(&l.UpdatedAt)
}

func (r *languageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "languages", id.String())
}

func scanLanguage(row scanner) (*domain.Language, error) {
	var l domain.Language
	if err := row.Scan(
		&l.ID, &l.VenueID, &l.Code, &l.Name, &l.IsDefault,
		&l.CreatedAt, &l.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &l, nil
}
