package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type tagRepository struct{ pool *pgxpool.Pool }

func NewTagRepository(pool *pgxpool.Pool) *tagRepository {
	return &tagRepository{pool: pool}
}

const tagSelectCols = `t.id, t.name, t.localization, t.created_at, t.updated_at`

func (r *tagRepository) List(ctx context.Context, p domain.Pagination) ([]*domain.Tag, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM tags WHERE deleted_at IS NULL`,
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx,
		`SELECT `+tagSelectCols+` FROM tags t WHERE t.deleted_at IS NULL
         ORDER BY t.name LIMIT $1 OFFSET $2`,
		p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []*domain.Tag
	for rows.Next() {
		t, err := scanTag(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, t)
	}
	return out, total, rows.Err()
}

func (r *tagRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Tag, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+tagSelectCols+` FROM tags t WHERE t.id=$1 AND t.deleted_at IS NULL`, id)
	return scanTag(row)
}

func (r *tagRepository) Create(ctx context.Context, t *domain.Tag) error {
	t.ID = newID()
	return r.pool.QueryRow(ctx,
		`INSERT INTO tags (id,name,localization) VALUES ($1,$2,$3)
         RETURNING created_at,updated_at`,
		t.ID, t.Name, t.Localization,
	).Scan(&t.CreatedAt, &t.UpdatedAt)
}

func (r *tagRepository) Update(ctx context.Context, t *domain.Tag) error {
	return r.pool.QueryRow(ctx,
		`UPDATE tags SET name=$2,localization=$3 WHERE id=$1 AND deleted_at IS NULL RETURNING updated_at`,
		t.ID, t.Name, t.Localization,
	).Scan(&t.UpdatedAt)
}

func (r *tagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE tags SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func (r *tagRepository) AttachTag(ctx context.Context, et *domain.EntityTag) error {
	et.ID = newID()
	return r.pool.QueryRow(ctx,
		`INSERT INTO entity_tags (id,tag_id,entity_type,entity_id) VALUES ($1,$2,$3,$4)
         ON CONFLICT DO NOTHING RETURNING created_at`,
		et.ID, et.TagID, et.EntityType, et.EntityID,
	).Scan(&et.CreatedAt)
}

func (r *tagRepository) DetachTag(ctx context.Context, tagID uuid.UUID, entityType string, entityID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM entity_tags WHERE tag_id=$1 AND entity_type=$2 AND entity_id=$3`,
		tagID, entityType, entityID)
	return err
}

func (r *tagRepository) ListEntityTags(ctx context.Context, entityType string, entityID uuid.UUID) ([]*domain.Tag, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+tagSelectCols+` FROM tags t
         JOIN entity_tags et ON et.tag_id=t.id
         WHERE et.entity_type=$1 AND et.entity_id=$2 AND t.deleted_at IS NULL`,
		entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Tag
	for rows.Next() {
		t, err := scanTag(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func scanTag(row scanner) (*domain.Tag, error) {
	var t domain.Tag
	if err := row.Scan(&t.ID, &t.Name, &t.Localization, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}
