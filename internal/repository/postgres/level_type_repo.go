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

type levelTypeRepo struct{ pool *pgxpool.Pool }

func NewLevelTypeRepository(pool *pgxpool.Pool) repository.LevelTypeRepository {
	return &levelTypeRepo{pool: pool}
}

const levelTypeSelectCols = `id, name, icon, value, created_at, updated_at`

func (r *levelTypeRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.LevelType, error) {
	q := `SELECT ` + levelTypeSelectCols + ` FROM level_types WHERE id = $1 AND deleted_at IS NULL`
	lt, err := scanLevelType(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("level type not found")
	}
	return lt, err
}

func (r *levelTypeRepo) List(ctx context.Context) ([]*domain.LevelType, error) {
	q := `SELECT ` + levelTypeSelectCols + ` FROM level_types WHERE deleted_at IS NULL ORDER BY name`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.LevelType
	for rows.Next() {
		lt, err := scanLevelType(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, lt)
	}
	return result, rows.Err()
}

func (r *levelTypeRepo) Create(ctx context.Context, lt *domain.LevelType) error {
	if lt.ID == uuid.Nil {
		lt.ID = newID()
	}
	const q = `
		INSERT INTO level_types (id, venue_id, name, icon, value)
		VALUES ($1, NULL, $2, $3, $4)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, q, lt.ID, lt.Name, lt.Icon, lt.Value).Scan(&lt.CreatedAt, &lt.UpdatedAt)
}

func (r *levelTypeRepo) Update(ctx context.Context, lt *domain.LevelType) error {
	const q = `
		UPDATE level_types SET name = $2, icon = $3, value = $4
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`
	err := r.pool.QueryRow(ctx, q, lt.ID, lt.Name, lt.Icon, lt.Value).Scan(&lt.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("level type not found")
	}
	return err
}

func (r *levelTypeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "level_types", id.String())
}

func scanLevelType(row scanner) (*domain.LevelType, error) {
	var lt domain.LevelType
	err := row.Scan(&lt.ID, &lt.Name, &lt.Icon, &lt.Value, &lt.CreatedAt, &lt.UpdatedAt)
	return &lt, err
}
