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

type levelBundleRepo struct {
	pool *pgxpool.Pool
}

// NewLevelBundleRepository returns a LevelBundleRepository backed by PostgreSQL.
func NewLevelBundleRepository(pool *pgxpool.Pool) repository.LevelBundleRepository {
	return &levelBundleRepo{pool: pool}
}

const levelBundleSelectCols = `id, snapshot_id, venue_id, level_id, state, created_at, updated_at`

func (r *levelBundleRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.LevelBundle, error) {
	q := `SELECT ` + levelBundleSelectCols + ` FROM level_bundles WHERE id = $1 AND deleted_at IS NULL`
	b, err := scanLevelBundle(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("level bundle not found")
	}
	return b, err
}

func (r *levelBundleRepo) ListBySnapshot(ctx context.Context, snapshotID uuid.UUID) ([]*domain.LevelBundle, error) {
	q := `SELECT ` + levelBundleSelectCols + ` FROM level_bundles WHERE snapshot_id = $1 AND deleted_at IS NULL ORDER BY created_at`
	rows, err := r.pool.Query(ctx, q, snapshotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bundles []*domain.LevelBundle
	for rows.Next() {
		b, err := scanLevelBundle(rows)
		if err != nil {
			return nil, err
		}
		bundles = append(bundles, b)
	}
	return bundles, rows.Err()
}

func (r *levelBundleRepo) Create(ctx context.Context, b *domain.LevelBundle) error {
	if b.ID == uuid.Nil {
		b.ID = newID()
	}
	const q = `
		INSERT INTO level_bundles (id, snapshot_id, venue_id, level_id, state)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, q,
		b.ID, b.SnapshotID, b.VenueID, b.LevelID, b.State,
	).Scan(&b.CreatedAt, &b.UpdatedAt)
}

func (r *levelBundleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "level_bundles", id.String())
}

func scanLevelBundle(row scanner) (*domain.LevelBundle, error) {
	var b domain.LevelBundle
	err := row.Scan(
		&b.ID, &b.SnapshotID, &b.VenueID, &b.LevelID,
		&b.State, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
