package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type connectionRepository struct{ pool *pgxpool.Pool }

func NewConnectionRepository(pool *pgxpool.Pool) *connectionRepository {
	return &connectionRepository{pool: pool}
}

const connectionSelectCols = `
    c.id, c.venue_id, c.external_id, c.name, c.type,
    c.x, c.y, c.state, c.status, c.accessible, c.active,
    c.created_at, c.updated_at`

func (r *connectionRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Connection, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM connections WHERE venue_id=$1 AND deleted_at IS NULL`, venueID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx,
		`SELECT `+connectionSelectCols+` FROM connections c WHERE c.venue_id=$1 AND c.deleted_at IS NULL
         ORDER BY c.created_at DESC LIMIT $2 OFFSET $3`,
		venueID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []*domain.Connection
	for rows.Next() {
		c, err := scanConnection(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, c)
	}
	return out, total, rows.Err()
}

func (r *connectionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Connection, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+connectionSelectCols+` FROM connections c WHERE c.id=$1 AND c.deleted_at IS NULL`, id)
	return scanConnection(row)
}

func (r *connectionRepository) Create(ctx context.Context, c *domain.Connection) error {
	c.ID = newID()
	return r.pool.QueryRow(ctx,
		`INSERT INTO connections (id,venue_id,external_id,name,type,x,y,state,status,accessible,active)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
         RETURNING created_at,updated_at`,
		c.ID, c.VenueID, c.ExternalID, c.Name, c.Type,
		c.X, c.Y, c.State, c.Status, c.Accessible, c.Active,
	).Scan(&c.CreatedAt, &c.UpdatedAt)
}

func (r *connectionRepository) Update(ctx context.Context, c *domain.Connection) error {
	return r.pool.QueryRow(ctx,
		`UPDATE connections SET external_id=$2,name=$3,type=$4,x=$5,y=$6,state=$7,status=$8,accessible=$9,active=$10
         WHERE id=$1 AND deleted_at IS NULL RETURNING updated_at`,
		c.ID, c.ExternalID, c.Name, c.Type, c.X, c.Y, c.State, c.Status, c.Accessible, c.Active,
	).Scan(&c.UpdatedAt)
}

func (r *connectionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE connections SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func (r *connectionRepository) ListLevels(ctx context.Context, connectionID uuid.UUID) ([]*domain.ConnectionLevel, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,connection_id,level_id,element_id,active,created_at,updated_at
         FROM connection_levels WHERE connection_id=$1 AND deleted_at IS NULL`, connectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.ConnectionLevel
	for rows.Next() {
		var cl domain.ConnectionLevel
		if err := rows.Scan(&cl.ID, &cl.ConnectionID, &cl.LevelID, &cl.ElementID, &cl.Active, &cl.CreatedAt, &cl.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, &cl)
	}
	return out, rows.Err()
}

func (r *connectionRepository) AddLevel(ctx context.Context, cl *domain.ConnectionLevel) error {
	cl.ID = newID()
	return r.pool.QueryRow(ctx,
		`INSERT INTO connection_levels (id,connection_id,level_id,element_id,active)
         VALUES ($1,$2,$3,$4,$5) RETURNING created_at,updated_at`,
		cl.ID, cl.ConnectionID, cl.LevelID, cl.ElementID, cl.Active,
	).Scan(&cl.CreatedAt, &cl.UpdatedAt)
}

func (r *connectionRepository) RemoveLevel(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE connection_levels SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func scanConnection(row scanner) (*domain.Connection, error) {
	var c domain.Connection
	if err := row.Scan(
		&c.ID, &c.VenueID, &c.ExternalID, &c.Name, &c.Type,
		&c.X, &c.Y, &c.State, &c.Status, &c.Accessible, &c.Active,
		&c.CreatedAt, &c.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &c, nil
}
