package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type qrcodeRepository struct{ pool *pgxpool.Pool }

func NewQRCodeRepository(pool *pgxpool.Pool) *qrcodeRepository {
	return &qrcodeRepository{pool: pool}
}

const qrcodeSelectCols = `
    q.id, q.venue_id, q.level_id, q.location_id,
    q.lat, q.lng, q.angle, q.link, q.base64_image,
    q.created_at, q.updated_at`

func (r *qrcodeRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.QRCode, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM qrcodes WHERE venue_id=$1 AND deleted_at IS NULL`, venueID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx,
		`SELECT `+qrcodeSelectCols+` FROM qrcodes q WHERE q.venue_id=$1 AND q.deleted_at IS NULL
         ORDER BY q.created_at DESC LIMIT $2 OFFSET $3`,
		venueID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []*domain.QRCode
	for rows.Next() {
		q, err := scanQRCode(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, q)
	}
	return out, total, rows.Err()
}

func (r *qrcodeRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.QRCode, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+qrcodeSelectCols+` FROM qrcodes q WHERE q.id=$1 AND q.deleted_at IS NULL`, id)
	return scanQRCode(row)
}

func (r *qrcodeRepository) Create(ctx context.Context, q *domain.QRCode) error {
	q.ID = newID()
	return r.pool.QueryRow(ctx,
		`INSERT INTO qrcodes (id,venue_id,level_id,location_id,lat,lng,angle,link,base64_image)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
         RETURNING created_at,updated_at`,
		q.ID, q.VenueID, q.LevelID, q.LocationID, q.Lat, q.Lng, q.Angle, q.Link, q.Base64Image,
	).Scan(&q.CreatedAt, &q.UpdatedAt)
}

func (r *qrcodeRepository) Update(ctx context.Context, q *domain.QRCode) error {
	return r.pool.QueryRow(ctx,
		`UPDATE qrcodes SET level_id=$2,location_id=$3,lat=$4,lng=$5,angle=$6,link=$7,base64_image=$8
         WHERE id=$1 AND deleted_at IS NULL RETURNING updated_at`,
		q.ID, q.LevelID, q.LocationID, q.Lat, q.Lng, q.Angle, q.Link, q.Base64Image,
	).Scan(&q.UpdatedAt)
}

func (r *qrcodeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE qrcodes SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func scanQRCode(row scanner) (*domain.QRCode, error) {
	var q domain.QRCode
	if err := row.Scan(
		&q.ID, &q.VenueID, &q.LevelID, &q.LocationID,
		&q.Lat, &q.Lng, &q.Angle, &q.Link, &q.Base64Image,
		&q.CreatedAt, &q.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &q, nil
}
