package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type snapshotRepo struct {
	pool *pgxpool.Pool
}

// NewSnapshotRepository returns a SnapshotRepository backed by PostgreSQL.
func NewSnapshotRepository(pool *pgxpool.Pool) repository.SnapshotRepository {
	return &snapshotRepo{pool: pool}
}

const snapshotSelectCols = `id, venue_id, state, method, created_by, publish_at, created_at, updated_at`

func (r *snapshotRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Snapshot, error) {
	q := `SELECT ` + snapshotSelectCols + ` FROM snapshots WHERE id = $1 AND deleted_at IS NULL`
	s, err := scanSnapshot(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("snapshot not found")
	}
	return s, err
}

func (r *snapshotRepo) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Snapshot, int64, error) {
	var total int64
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM snapshots WHERE venue_id = $1 AND deleted_at IS NULL`, venueID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	q := `SELECT ` + snapshotSelectCols + ` FROM snapshots WHERE venue_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, venueID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var snapshots []*domain.Snapshot
	for rows.Next() {
		s, err := scanSnapshot(rows)
		if err != nil {
			return nil, 0, err
		}
		snapshots = append(snapshots, s)
	}
	return snapshots, total, rows.Err()
}

func (r *snapshotRepo) LatestPublished(ctx context.Context, venueID uuid.UUID) (*domain.Snapshot, error) {
	q := `SELECT ` + snapshotSelectCols + ` FROM snapshots WHERE venue_id = $1 AND state = $2 AND deleted_at IS NULL ORDER BY publish_at DESC LIMIT 1`
	s, err := scanSnapshot(r.pool.QueryRow(ctx, q, venueID, domain.SnapshotStatePublic))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("no published snapshot found")
	}
	return s, err
}

func (r *snapshotRepo) Create(ctx context.Context, s *domain.Snapshot) error {
	if s.ID == uuid.Nil {
		s.ID = newID()
	}
	const q = `
		INSERT INTO snapshots (id, venue_id, state, method, created_by, publish_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, q,
		s.ID, s.VenueID, s.State, s.Method, s.CreatedBy, s.PublishAt,
	).Scan(&s.CreatedAt, &s.UpdatedAt)
}

func (r *snapshotRepo) UpdateState(ctx context.Context, id uuid.UUID, state int, publishAt *time.Time) error {
	const q = `UPDATE snapshots SET state = $2, publish_at = $3 WHERE id = $1 AND deleted_at IS NULL RETURNING updated_at`
	var updatedAt time.Time
	err := r.pool.QueryRow(ctx, q, id, state, publishAt).Scan(&updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("snapshot not found")
	}
	return err
}

func (r *snapshotRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "snapshots", id.String())
}

func (r *snapshotRepo) CountDraftsByVenue(ctx context.Context, venueID uuid.UUID) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM snapshots WHERE venue_id = $1 AND state = $2 AND deleted_at IS NULL`,
		venueID, domain.SnapshotStateDraft,
	).Scan(&count)
	return count, err
}

func (r *snapshotRepo) UnpublishVenue(ctx context.Context, venueID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE snapshots SET state = $1 WHERE venue_id = $2 AND state = $3 AND deleted_at IS NULL`,
		domain.SnapshotStateDraft, venueID, domain.SnapshotStatePublic,
	)
	return err
}

func (r *snapshotRepo) DeleteOldestDraft(ctx context.Context, venueID uuid.UUID) error {
	// Hard-delete (not soft-delete) so ON DELETE CASCADE fires for level_bundles.
	_, err := r.pool.Exec(ctx, `
		DELETE FROM snapshots
		WHERE id = (
			SELECT id FROM snapshots
			WHERE venue_id = $1 AND state = $2 AND deleted_at IS NULL
			ORDER BY created_at ASC
			LIMIT 1
		)`, venueID, domain.SnapshotStateDraft)
	return err
}

func scanSnapshot(row scanner) (*domain.Snapshot, error) {
	var s domain.Snapshot
	err := row.Scan(
		&s.ID, &s.VenueID, &s.State, &s.Method,
		&s.CreatedBy, &s.PublishAt,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
