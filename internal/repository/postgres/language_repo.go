package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/database"
)

type languageRepository struct{ pool *pgxpool.Pool }

func NewLanguageRepository(pool *pgxpool.Pool) *languageRepository {
	return &languageRepository{pool: pool}
}

const languageSelectCols = `
    l.id, l.venue_id, l.code, l.name, l.is_default, l.enabled,
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
		`INSERT INTO languages (id, venue_id, code, name, is_default, enabled)
         VALUES ($1,$2,$3,$4,$5,$6)
         RETURNING created_at, updated_at`,
		l.ID, l.VenueID, l.Code, l.Name, l.IsDefault, l.Enabled,
	).Scan(&l.CreatedAt, &l.UpdatedAt)
}

func (r *languageRepository) Update(ctx context.Context, l *domain.Language) error {
	return r.pool.QueryRow(ctx,
		`UPDATE languages SET code=$2, name=$3, is_default=$4, enabled=$5
         WHERE id=$1 AND deleted_at IS NULL RETURNING updated_at`,
		l.ID, l.Code, l.Name, l.IsDefault, l.Enabled,
	).Scan(&l.UpdatedAt)
}

func (r *languageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "languages", id.String())
}

func (r *languageRepository) ListEnabled(ctx context.Context, venueID uuid.UUID) ([]*domain.Language, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+languageSelectCols+` FROM languages l
         WHERE l.venue_id=$1 AND l.deleted_at IS NULL AND l.enabled = true ORDER BY l.created_at ASC`,
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

// SetDefault marks langID as the sole default language for venueID, clearing
// is_default on every other language in the same statement.
func (r *languageRepository) SetDefault(ctx context.Context, venueID, langID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE languages SET is_default = (id = $2), updated_at = NOW()
         WHERE venue_id = $1 AND deleted_at IS NULL`,
		venueID, langID)
	return err
}

// ReplaceAll reconciles a venue's full language set against items in one
// transaction: codes missing from items are soft-deleted, codes present with
// changed name/is_default are updated, and new codes are inserted.
func (r *languageRepository) ReplaceAll(ctx context.Context, venueID uuid.UUID, items []*domain.Language) ([]*domain.Language, error) {
	err := database.WithTransaction(ctx, r.pool, func(ctx context.Context, tx pgx.Tx) error {
		existing, err := listLanguagesForUpdate(ctx, tx, venueID)
		if err != nil {
			return err
		}
		existingByCode := make(map[string]*domain.Language, len(existing))
		for _, l := range existing {
			existingByCode[l.Code] = l
		}
		wantedCodes := make(map[string]bool, len(items))
		for _, item := range items {
			wantedCodes[item.Code] = true
		}

		for _, l := range existing {
			if wantedCodes[l.Code] {
				continue
			}
			if _, err := tx.Exec(ctx,
				`UPDATE languages SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1`,
				l.ID,
			); err != nil {
				return err
			}
		}

		for _, item := range items {
			if existingRow, ok := existingByCode[item.Code]; ok {
				if existingRow.Name == item.Name && existingRow.IsDefault == item.IsDefault && existingRow.Enabled == item.Enabled {
					continue
				}
				if _, err := tx.Exec(ctx,
					`UPDATE languages SET name=$2, is_default=$3, enabled=$4, updated_at=NOW() WHERE id=$1`,
					existingRow.ID, item.Name, item.IsDefault, item.Enabled,
				); err != nil {
					return err
				}
				continue
			}
			if _, err := tx.Exec(ctx,
				`INSERT INTO languages (id, venue_id, code, name, is_default, enabled) VALUES ($1,$2,$3,$4,$5,$6)`,
				newID(), venueID, item.Code, item.Name, item.IsDefault, item.Enabled,
			); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.List(ctx, venueID)
}

func listLanguagesForUpdate(ctx context.Context, tx pgx.Tx, venueID uuid.UUID) ([]*domain.Language, error) {
	rows, err := tx.Query(ctx,
		`SELECT `+languageSelectCols+` FROM languages l
         WHERE l.venue_id=$1 AND l.deleted_at IS NULL ORDER BY l.created_at ASC FOR UPDATE`,
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

func scanLanguage(row scanner) (*domain.Language, error) {
	var l domain.Language
	if err := row.Scan(
		&l.ID, &l.VenueID, &l.Code, &l.Name, &l.IsDefault, &l.Enabled,
		&l.CreatedAt, &l.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &l, nil
}
