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

type assetRepo struct{ pool *pgxpool.Pool }

func NewAssetRepository(pool *pgxpool.Pool) repository.AssetRepository {
	return &assetRepo{pool: pool}
}

const assetSelectCols = `id, venue_id, name, key, content_type, size_bytes, url, created_by, created_at, updated_at`

func (r *assetRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	q := `SELECT ` + assetSelectCols + ` FROM assets WHERE id = $1 AND deleted_at IS NULL`
	a, err := scanAsset(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("asset not found")
	}
	return a, err
}

func (r *assetRepo) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Asset, int64, error) {
	var total int64
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM assets WHERE venue_id = $1 AND deleted_at IS NULL`, venueID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	q := `SELECT ` + assetSelectCols + ` FROM assets WHERE venue_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, venueID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var assets []*domain.Asset
	for rows.Next() {
		a, err := scanAsset(rows)
		if err != nil {
			return nil, 0, err
		}
		assets = append(assets, a)
	}
	return assets, total, rows.Err()
}

func (r *assetRepo) Create(ctx context.Context, a *domain.Asset) error {
	if a.ID == uuid.Nil {
		a.ID = newID()
	}
	const q = `
		INSERT INTO assets (id, venue_id, name, key, content_type, size_bytes, url, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, q,
		a.ID, a.VenueID, a.Name, a.Key, a.ContentType, a.SizeBytes, a.URL, a.CreatedBy,
	).Scan(&a.CreatedAt, &a.UpdatedAt)
}

func (r *assetRepo) Update(ctx context.Context, a *domain.Asset) error {
	const q = `
		UPDATE assets SET name = $2, url = $3
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`
	err := r.pool.QueryRow(ctx, q, a.ID, a.Name, a.URL).Scan(&a.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("asset not found")
	}
	return err
}

func (r *assetRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "assets", id.String())
}

func scanAsset(row scanner) (*domain.Asset, error) {
	var a domain.Asset
	err := row.Scan(
		&a.ID, &a.VenueID, &a.Name, &a.Key, &a.ContentType,
		&a.SizeBytes, &a.URL, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt,
	)
	return &a, err
}
