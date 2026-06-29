package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type adRepository struct{ pool *pgxpool.Pool }

func NewAdRepository(pool *pgxpool.Pool) *adRepository {
	return &adRepository{pool: pool}
}

const adSelectCols = `
    a.id, a.venue_id, a.location_id, a.type, a.status,
    a.navigate, a.content_image, a.content_cta_url, a.placement,
    a.size_width, a.size_height, a.reward_type, a.reward_amount, a.display_duration,
    a.published_at, a.start_at, a.end_at,
    a.created_at, a.updated_at`

func (r *adRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Advertisement, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM advertisements WHERE venue_id=$1 AND deleted_at IS NULL`, venueID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx,
		`SELECT `+adSelectCols+` FROM advertisements a WHERE a.venue_id=$1 AND a.deleted_at IS NULL
         ORDER BY a.created_at DESC LIMIT $2 OFFSET $3`,
		venueID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []*domain.Advertisement
	for rows.Next() {
		a, err := scanAd(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, a)
	}
	return out, total, rows.Err()
}

func (r *adRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Advertisement, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+adSelectCols+` FROM advertisements a WHERE a.id=$1 AND a.deleted_at IS NULL`, id)
	return scanAd(row)
}

func (r *adRepository) Create(ctx context.Context, a *domain.Advertisement) error {
	if a.ID == uuid.Nil {
		a.ID = newID()
	}
	return r.pool.QueryRow(ctx,
		`INSERT INTO advertisements
         (id,venue_id,location_id,type,status,navigate,content_image,content_cta_url,placement,
          size_width,size_height,reward_type,reward_amount,display_duration,published_at,start_at,end_at)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
         RETURNING created_at,updated_at`,
		a.ID, a.VenueID, a.LocationID, a.Type, a.Status,
		a.Navigate, a.ContentImage, a.ContentCTAURL, a.Placement,
		a.SizeWidth, a.SizeHeight, a.RewardType, a.RewardAmount, a.DisplayDuration,
		a.PublishedAt, a.StartAt, a.EndAt,
	).Scan(&a.CreatedAt, &a.UpdatedAt)
}

func (r *adRepository) Update(ctx context.Context, a *domain.Advertisement) error {
	return r.pool.QueryRow(ctx,
		`UPDATE advertisements SET
         location_id=$2,type=$3,status=$4,navigate=$5,content_image=$6,content_cta_url=$7,
         placement=$8,size_width=$9,size_height=$10,reward_type=$11,reward_amount=$12,
         display_duration=$13,published_at=$14,start_at=$15,end_at=$16
         WHERE id=$1 AND deleted_at IS NULL RETURNING updated_at`,
		a.ID, a.LocationID, a.Type, a.Status, a.Navigate, a.ContentImage, a.ContentCTAURL,
		a.Placement, a.SizeWidth, a.SizeHeight, a.RewardType, a.RewardAmount,
		a.DisplayDuration, a.PublishedAt, a.StartAt, a.EndAt,
	).Scan(&a.UpdatedAt)
}

func (r *adRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE advertisements SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func scanAd(row scanner) (*domain.Advertisement, error) {
	var a domain.Advertisement
	if err := row.Scan(
		&a.ID, &a.VenueID, &a.LocationID, &a.Type, &a.Status,
		&a.Navigate, &a.ContentImage, &a.ContentCTAURL, &a.Placement,
		&a.SizeWidth, &a.SizeHeight, &a.RewardType, &a.RewardAmount, &a.DisplayDuration,
		&a.PublishedAt, &a.StartAt, &a.EndAt,
		&a.CreatedAt, &a.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &a, nil
}
