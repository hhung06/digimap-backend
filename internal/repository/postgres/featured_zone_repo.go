package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type featuredZoneRepository struct{ pool *pgxpool.Pool }

func NewFeaturedZoneRepository(pool *pgxpool.Pool) *featuredZoneRepository {
	return &featuredZoneRepository{pool: pool}
}

const featuredZoneSelectCols = `
    z.id, z.venue_id, z.name, z.description, z.image_url,
    z.sort_index, z.is_active, z.created_at, z.updated_at`

func (r *featuredZoneRepository) List(ctx context.Context, venueID uuid.UUID) ([]*domain.FeaturedZone, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+featuredZoneSelectCols+` FROM featured_zones z
         WHERE z.venue_id=$1 AND z.deleted_at IS NULL ORDER BY z.sort_index ASC, z.created_at ASC`,
		venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.FeaturedZone
	for rows.Next() {
		z, err := scanFeaturedZone(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, z)
	}
	return out, rows.Err()
}

func (r *featuredZoneRepository) ListActive(ctx context.Context, venueID uuid.UUID) ([]*domain.FeaturedZone, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+featuredZoneSelectCols+` FROM featured_zones z
         WHERE z.venue_id=$1 AND z.is_active=true AND z.deleted_at IS NULL
         ORDER BY z.sort_index ASC, z.created_at ASC`,
		venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.FeaturedZone
	for rows.Next() {
		z, err := scanFeaturedZone(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, z)
	}
	return out, rows.Err()
}

func (r *featuredZoneRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.FeaturedZone, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+featuredZoneSelectCols+` FROM featured_zones z WHERE z.id=$1 AND z.deleted_at IS NULL`, id)
	return scanFeaturedZone(row)
}

func (r *featuredZoneRepository) Create(ctx context.Context, z *domain.FeaturedZone) error {
	z.ID = newID()
	return r.pool.QueryRow(ctx,
		`INSERT INTO featured_zones (id, venue_id, name, description, image_url, sort_index, is_active)
         VALUES ($1,$2,$3,$4,$5,$6,$7)
         RETURNING created_at, updated_at`,
		z.ID, z.VenueID, z.Name, z.Description, z.ImageURL, z.SortIndex, z.IsActive,
	).Scan(&z.CreatedAt, &z.UpdatedAt)
}

func (r *featuredZoneRepository) Update(ctx context.Context, z *domain.FeaturedZone) error {
	return r.pool.QueryRow(ctx,
		`UPDATE featured_zones SET name=$2, description=$3, image_url=$4, sort_index=$5, is_active=$6
         WHERE id=$1 AND deleted_at IS NULL RETURNING updated_at`,
		z.ID, z.Name, z.Description, z.ImageURL, z.SortIndex, z.IsActive,
	).Scan(&z.UpdatedAt)
}

func (r *featuredZoneRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "featured_zones", id.String())
}

func scanFeaturedZone(row scanner) (*domain.FeaturedZone, error) {
	var z domain.FeaturedZone
	if err := row.Scan(
		&z.ID, &z.VenueID, &z.Name, &z.Description, &z.ImageURL,
		&z.SortIndex, &z.IsActive, &z.CreatedAt, &z.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &z, nil
}
