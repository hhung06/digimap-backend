package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type customerRepo struct {
	pool *pgxpool.Pool
}

// NewCustomerRepository returns a CustomerRepository backed by PostgreSQL.
func NewCustomerRepository(pool *pgxpool.Pool) repository.CustomerRepository {
	return &customerRepo{pool: pool}
}

func (r *customerRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	const q = `
		SELECT id, name, image, phone, email, address, url, description,
		       created_at, updated_at, deleted_at
		FROM customers
		WHERE id = $1 AND deleted_at IS NULL`

	c, err := scanCustomer(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("customer not found")
	}
	return c, err
}

func (r *customerRepo) List(ctx context.Context, p domain.Pagination) ([]*domain.Customer, int64, error) {
	const countQ = `SELECT COUNT(*) FROM customers WHERE deleted_at IS NULL`
	const q = `
		SELECT id, name, image, phone, email, address, url, description,
		       created_at, updated_at, deleted_at
		FROM customers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	var total int64
	if err := r.pool.QueryRow(ctx, countQ).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx, q, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var customers []*domain.Customer
	for rows.Next() {
		c, err := scanCustomer(rows)
		if err != nil {
			return nil, 0, err
		}
		customers = append(customers, c)
	}
	return customers, total, rows.Err()
}

func (r *customerRepo) Create(ctx context.Context, c *domain.Customer) error {
	if c.ID == uuid.Nil {
		c.ID = newID()
	}
	const q = `
		INSERT INTO customers (id, name, image, phone, email, address, url, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		c.ID, c.Name, nullStr(c.Image), nullStr(c.Phone), nullStr(c.Email),
		nullStr(c.Address), nullStr(c.URL), nullStr(c.Description),
	).Scan(&c.CreatedAt, &c.UpdatedAt)
}

func (r *customerRepo) Update(ctx context.Context, c *domain.Customer) error {
	const q = `
		UPDATE customers
		SET name = $2, image = $3, phone = $4, email = $5, address = $6, url = $7, description = $8
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, q,
		c.ID, c.Name, nullStr(c.Image), nullStr(c.Phone), nullStr(c.Email),
		nullStr(c.Address), nullStr(c.URL), nullStr(c.Description),
	).Scan(&c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("customer not found")
	}
	return err
}

func (r *customerRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "customers", id.String())
}

// ── helpers ───────────────────────────────────────────────────────────────────

func scanCustomer(row pgx.Row) (*domain.Customer, error) {
	var c domain.Customer
	var image, phone, email, address, url, description *string
	var deletedAt *time.Time
	err := row.Scan(
		&c.ID, &c.Name, &image, &phone, &email, &address, &url, &description,
		&c.CreatedAt, &c.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}
	derefStr(&c.Image, image)
	derefStr(&c.Phone, phone)
	derefStr(&c.Email, email)
	derefStr(&c.Address, address)
	derefStr(&c.URL, url)
	derefStr(&c.Description, description)
	c.DeletedAt = deletedAt
	return &c, nil
}
