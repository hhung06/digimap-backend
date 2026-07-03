package postgres

import (
	"context"
	"errors"
	"fmt"

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

const assetSelectCols = `id, venue_id, name, key, content_type, size_bytes, url,
	asset_type, description, file_type, thumbnail, material, width, height, status,
	created_by, created_at, updated_at`

func (r *assetRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	q := `SELECT ` + assetSelectCols + ` FROM assets WHERE id = $1 AND deleted_at IS NULL`
	a, err := scanAsset(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("asset not found")
	}
	return a, err
}

func (r *assetRepo) FindByIDs(ctx context.Context, venueID uuid.UUID, ids []uuid.UUID) ([]*domain.Asset, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	rows, err := r.pool.Query(ctx,
		`SELECT `+assetSelectCols+` FROM assets WHERE venue_id = $1 AND id = ANY($2) AND deleted_at IS NULL`,
		venueID, ids,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var assets []*domain.Asset
	for rows.Next() {
		a, err := scanAsset(rows)
		if err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, rows.Err()
}

func (r *assetRepo) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination, assetType string) ([]*domain.Asset, int64, error) {
	var total int64
	var countErr error
	if assetType != "" {
		countErr = r.pool.QueryRow(ctx,
			`SELECT COUNT(*) FROM assets WHERE venue_id = $1 AND asset_type = $2 AND deleted_at IS NULL`,
			venueID, assetType,
		).Scan(&total)
	} else {
		countErr = r.pool.QueryRow(ctx,
			`SELECT COUNT(*) FROM assets WHERE venue_id = $1 AND deleted_at IS NULL`,
			venueID,
		).Scan(&total)
	}
	if countErr != nil {
		return nil, 0, countErr
	}

	var rows pgx.Rows
	var err error
	if assetType != "" {
		rows, err = r.pool.Query(ctx,
			`SELECT `+assetSelectCols+` FROM assets WHERE venue_id = $1 AND asset_type = $2 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $3 OFFSET $4`,
			venueID, assetType, p.PageSize, p.Offset(),
		)
	} else {
		rows, err = r.pool.Query(ctx,
			`SELECT `+assetSelectCols+` FROM assets WHERE venue_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
			venueID, p.PageSize, p.Offset(),
		)
	}
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

func (r *assetRepo) ListLibrary(ctx context.Context, env, status, assetType string) ([]*domain.Asset, error) {
	// Library assets are identified by key prefix "{env}/library/".
	args := []any{"%" + env + "/library/%"}
	where := "key LIKE $1 AND deleted_at IS NULL"
	if status != "" {
		args = append(args, status)
		where += fmt.Sprintf(" AND status = $%d", len(args))
	}
	if assetType != "" {
		args = append(args, assetType)
		where += fmt.Sprintf(" AND asset_type = $%d", len(args))
	}

	rows, err := r.pool.Query(ctx,
		`SELECT `+assetSelectCols+` FROM assets WHERE `+where+` ORDER BY created_at DESC`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*domain.Asset
	for rows.Next() {
		a, err := scanAsset(rows)
		if err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, rows.Err()
}

func (r *assetRepo) ListAll(ctx context.Context, assetType string) ([]*domain.Asset, error) {
	var rows pgx.Rows
	var err error
	if assetType != "" {
		rows, err = r.pool.Query(ctx,
			`SELECT `+assetSelectCols+` FROM assets WHERE asset_type = $1 AND deleted_at IS NULL ORDER BY created_at DESC`,
			assetType,
		)
	} else {
		rows, err = r.pool.Query(ctx,
			`SELECT `+assetSelectCols+` FROM assets WHERE deleted_at IS NULL ORDER BY created_at DESC`,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*domain.Asset
	for rows.Next() {
		a, err := scanAsset(rows)
		if err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, rows.Err()
}

func (r *assetRepo) Create(ctx context.Context, a *domain.Asset) error {
	if a.ID == uuid.Nil {
		a.ID = newID()
	}
	if a.AssetType == "" {
		a.AssetType = domain.AssetType2D
	}
	if a.Status == "" {
		a.Status = domain.AssetStatusUnpublished
	}
	const q = `
		INSERT INTO assets (id, venue_id, name, key, content_type, size_bytes, url,
		                    asset_type, description, file_type, thumbnail, material, width, height, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, q,
		a.ID, a.VenueID, a.Name, a.Key, a.ContentType, a.SizeBytes, a.URL,
		a.AssetType, a.Description, a.FileType, nullStr(a.Thumbnail), nullStr(a.Material), a.Width, a.Height, a.Status,
		a.CreatedBy,
	).Scan(&a.CreatedAt, &a.UpdatedAt)
}

func (r *assetRepo) Update(ctx context.Context, a *domain.Asset) error {
	const q = `
		UPDATE assets
		SET name = $2, key = $3, content_type = $4, size_bytes = $5, url = $6, asset_type = $7,
		    description = $8, file_type = $9, thumbnail = $10, material = $11, width = $12,
		    height = $13, status = $14
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`
	err := r.pool.QueryRow(ctx, q,
		a.ID, a.Name, a.Key, a.ContentType, a.SizeBytes, a.URL, a.AssetType, a.Description, a.FileType,
		nullStr(a.Thumbnail), nullStr(a.Material), a.Width, a.Height, a.Status,
	).Scan(&a.UpdatedAt)
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
	var thumbnail, material *string
	err := row.Scan(
		&a.ID, &a.VenueID, &a.Name, &a.Key, &a.ContentType,
		&a.SizeBytes, &a.URL,
		&a.AssetType, &a.Description, &a.FileType, &thumbnail, &material, &a.Width, &a.Height, &a.Status,
		&a.CreatedBy, &a.CreatedAt, &a.UpdatedAt,
	)
	if thumbnail != nil {
		a.Thumbnail = *thumbnail
	}
	if material != nil {
		a.Material = *material
	}
	return &a, err
}
