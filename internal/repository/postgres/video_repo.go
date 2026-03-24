package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type videoRepository struct{ pool *pgxpool.Pool }

func NewVideoRepository(pool *pgxpool.Pool) *videoRepository {
	return &videoRepository{pool: pool}
}

const videoSelectCols = `
    v.id, v.venue_id, v.title, v.description, v.url,
    v.thumbnail, v.duration, v.status, v.published_at,
    v.created_at, v.updated_at`

func (r *videoRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Video, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM videos WHERE venue_id=$1 AND deleted_at IS NULL`, venueID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx,
		`SELECT `+videoSelectCols+` FROM videos v WHERE v.venue_id=$1 AND v.deleted_at IS NULL
         ORDER BY v.created_at DESC LIMIT $2 OFFSET $3`,
		venueID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []*domain.Video
	for rows.Next() {
		v, err := scanVideo(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, v)
	}
	return out, total, rows.Err()
}

func (r *videoRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Video, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+videoSelectCols+` FROM videos v WHERE v.id=$1 AND v.deleted_at IS NULL`, id)
	return scanVideo(row)
}

func (r *videoRepository) Create(ctx context.Context, v *domain.Video) error {
	v.ID = newID()
	return r.pool.QueryRow(ctx,
		`INSERT INTO videos (id,venue_id,title,description,url,thumbnail,duration,status,published_at)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
         RETURNING created_at,updated_at`,
		v.ID, v.VenueID, v.Title, v.Description, v.URL,
		v.Thumbnail, v.Duration, v.Status, v.PublishedAt,
	).Scan(&v.CreatedAt, &v.UpdatedAt)
}

func (r *videoRepository) Update(ctx context.Context, v *domain.Video) error {
	return r.pool.QueryRow(ctx,
		`UPDATE videos SET title=$2,description=$3,url=$4,thumbnail=$5,duration=$6,status=$7,published_at=$8
         WHERE id=$1 AND deleted_at IS NULL RETURNING updated_at`,
		v.ID, v.Title, v.Description, v.URL, v.Thumbnail, v.Duration, v.Status, v.PublishedAt,
	).Scan(&v.UpdatedAt)
}

func (r *videoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE videos SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func scanVideo(row scanner) (*domain.Video, error) {
	var v domain.Video
	if err := row.Scan(
		&v.ID, &v.VenueID, &v.Title, &v.Description, &v.URL,
		&v.Thumbnail, &v.Duration, &v.Status, &v.PublishedAt,
		&v.CreatedAt, &v.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &v, nil
}
