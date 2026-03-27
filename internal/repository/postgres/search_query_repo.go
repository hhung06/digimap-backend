package postgres

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type searchQueryRepository struct{ pool *pgxpool.Pool }

func NewSearchQueryRepository(pool *pgxpool.Pool) *searchQueryRepository {
	return &searchQueryRepository{pool: pool}
}

const searchQuerySelectCols = `
    q.id, q.venue_id, q.app_id, q.origin, q.search_term,
    q.search_count, q.last_searched, q.is_promoted, q.reference,
    q.status, q.created_at, q.updated_at`

func (r *searchQueryRepository) Upsert(ctx context.Context, venueID uuid.UUID, term string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO search_queries (id, venue_id, search_term, search_count, last_searched)
         VALUES ($1, $2, $3, 1, NOW())
         ON CONFLICT (venue_id, search_term) DO UPDATE
         SET search_count = search_queries.search_count + 1,
             last_searched = NOW(),
             updated_at = NOW()`,
		newID(), venueID, term)
	return err
}

func (r *searchQueryRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.SearchQuery, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM search_queries WHERE venue_id=$1 AND deleted_at IS NULL`, venueID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx,
		`SELECT `+searchQuerySelectCols+` FROM search_queries q
         WHERE q.venue_id=$1 AND q.deleted_at IS NULL
         ORDER BY q.search_count DESC, q.last_searched DESC LIMIT $2 OFFSET $3`,
		venueID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []*domain.SearchQuery
	for rows.Next() {
		q, err := scanSearchQuery(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, q)
	}
	return out, total, rows.Err()
}

func (r *searchQueryRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.SearchQuery, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+searchQuerySelectCols+` FROM search_queries q WHERE q.id=$1 AND q.deleted_at IS NULL`, id)
	return scanSearchQuery(row)
}

func (r *searchQueryRepository) Create(ctx context.Context, q *domain.SearchQuery) error {
	q.ID = newID()
	ref := q.Reference
	if ref == nil {
		ref = json.RawMessage("{}")
	}
	return r.pool.QueryRow(ctx,
		`INSERT INTO search_queries
         (id,venue_id,app_id,origin,search_term,search_count,last_searched,is_promoted,reference,status)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
         RETURNING created_at,updated_at`,
		q.ID, q.VenueID, q.AppID, q.Origin, q.SearchTerm,
		q.SearchCount, q.LastSearched, q.IsPromoted, ref, q.Status,
	).Scan(&q.CreatedAt, &q.UpdatedAt)
}

func (r *searchQueryRepository) Update(ctx context.Context, q *domain.SearchQuery) error {
	ref := q.Reference
	if ref == nil {
		ref = json.RawMessage("{}")
	}
	return r.pool.QueryRow(ctx,
		`UPDATE search_queries SET
         app_id=$2,origin=$3,search_term=$4,search_count=$5,
         last_searched=$6,is_promoted=$7,reference=$8,status=$9
         WHERE id=$1 AND deleted_at IS NULL RETURNING updated_at`,
		q.ID, q.AppID, q.Origin, q.SearchTerm, q.SearchCount,
		q.LastSearched, q.IsPromoted, ref, q.Status,
	).Scan(&q.UpdatedAt)
}

func (r *searchQueryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE search_queries SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func scanSearchQuery(row scanner) (*domain.SearchQuery, error) {
	var q domain.SearchQuery
	if err := row.Scan(
		&q.ID, &q.VenueID, &q.AppID, &q.Origin, &q.SearchTerm,
		&q.SearchCount, &q.LastSearched, &q.IsPromoted, &q.Reference,
		&q.Status, &q.CreatedAt, &q.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &q, nil
}
