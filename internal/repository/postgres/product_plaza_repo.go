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

type productPlazaRepo struct{ pool *pgxpool.Pool }

func NewProductPlazaRepository(pool *pgxpool.Pool) repository.ProductPlazaRepository {
	return &productPlazaRepo{pool: pool}
}

const productPlazaSelectCols = `id, venue_id, name, description, location_id, created_at, updated_at`

func (r *productPlazaRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.ProductPlaza, error) {
	q := `SELECT ` + productPlazaSelectCols + ` FROM product_plazas WHERE id = $1 AND deleted_at IS NULL`
	p, err := scanProductPlaza(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("product plaza not found")
	}
	return p, err
}

func (r *productPlazaRepo) List(ctx context.Context, venueID uuid.UUID) ([]*domain.ProductPlaza, error) {
	q := `SELECT ` + productPlazaSelectCols + ` FROM product_plazas WHERE venue_id = $1 AND deleted_at IS NULL ORDER BY name`
	rows, err := r.pool.Query(ctx, q, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.ProductPlaza
	for rows.Next() {
		p, err := scanProductPlaza(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, rows.Err()
}

func (r *productPlazaRepo) Create(ctx context.Context, p *domain.ProductPlaza) error {
	if p.ID == uuid.Nil {
		p.ID = newID()
	}
	const q = `
		INSERT INTO product_plazas (id, venue_id, name, description, location_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, q, p.ID, p.VenueID, p.Name, p.Description, p.LocationID).Scan(&p.CreatedAt, &p.UpdatedAt)
}

func (r *productPlazaRepo) Update(ctx context.Context, p *domain.ProductPlaza) error {
	const q = `
		UPDATE product_plazas SET name = $2, description = $3, location_id = $4
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`
	err := r.pool.QueryRow(ctx, q, p.ID, p.Name, p.Description, p.LocationID).Scan(&p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("product plaza not found")
	}
	return err
}

func (r *productPlazaRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "product_plazas", id.String())
}

func scanProductPlaza(row scanner) (*domain.ProductPlaza, error) {
	var p domain.ProductPlaza
	err := row.Scan(&p.ID, &p.VenueID, &p.Name, &p.Description, &p.LocationID, &p.CreatedAt, &p.UpdatedAt)
	return &p, err
}
