package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type productRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) repository.ProductRepository {
	return &productRepo{pool: pool}
}

// ── Categories ────────────────────────────────────────────────────────────────

func (r *productRepo) FindCategoryByID(ctx context.Context, id uuid.UUID) (*domain.ProductCategory, error) {
	const q = `
		SELECT id, venue_id, external_id, name, source, localization, created_at, updated_at, deleted_at
		FROM product_categories WHERE id = $1 AND deleted_at IS NULL`

	c, err := scanProductCategory(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("product category not found")
	}
	return c, err
}

func (r *productRepo) ListCategories(ctx context.Context, venueID uuid.UUID) ([]*domain.ProductCategory, error) {
	const q = `
		SELECT id, venue_id, external_id, name, source, localization, created_at, updated_at, deleted_at
		FROM product_categories WHERE venue_id = $1 AND deleted_at IS NULL ORDER BY name`

	rows, err := r.pool.Query(ctx, q, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []*domain.ProductCategory
	for rows.Next() {
		c, err := scanProductCategory(rows)
		if err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

func (r *productRepo) CreateCategory(ctx context.Context, c *domain.ProductCategory) error {
	if c.ID == uuid.Nil {
		c.ID = newID()
	}
	const q = `
		INSERT INTO product_categories (id, venue_id, external_id, name, source, localization)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		c.ID, c.VenueID, c.ExternalID, c.Name, c.Source, jsonOrNil(c.Localization),
	).Scan(&c.CreatedAt, &c.UpdatedAt)
}

func (r *productRepo) UpdateCategory(ctx context.Context, c *domain.ProductCategory) error {
	const q = `
		UPDATE product_categories
		SET external_id=$2, name=$3, source=$4, localization=$5
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, q,
		c.ID, c.ExternalID, c.Name, c.Source, jsonOrNil(c.Localization),
	).Scan(&c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("product category not found")
	}
	return err
}

func (r *productRepo) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "product_categories", id.String())
}

// ── Products ──────────────────────────────────────────────────────────────────

func (r *productRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	const q = `
		SELECT id, venue_id, location_id, main_category_id, image, name, external_id, size, price,
		       origin_country, expiration, description, custom, localization, source,
		       created_at, updated_at, deleted_at
		FROM products WHERE id = $1 AND deleted_at IS NULL`

	p, err := scanProduct(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("product not found")
	}
	if err != nil {
		return nil, err
	}
	cats, err := r.loadProductCategories(ctx, id)
	if err != nil {
		return nil, err
	}
	p.Categories = cats
	attachments, err := r.ListAttachments(ctx, id)
	if err != nil {
		return nil, err
	}
	p.Attachments = attachments
	return p, nil
}

func (r *productRepo) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Product, int64, error) {
	const countQ = `SELECT COUNT(*) FROM products WHERE venue_id = $1 AND deleted_at IS NULL`
	const q = `
		SELECT id, venue_id, location_id, main_category_id, image, name, external_id, size, price,
		       origin_country, expiration, description, custom, localization, source,
		       created_at, updated_at, deleted_at
		FROM products WHERE venue_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	var total int64
	if err := r.pool.QueryRow(ctx, countQ, venueID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx, q, venueID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		prod, err := scanProduct(rows)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, prod)
	}
	return products, total, rows.Err()
}

func (r *productRepo) Create(ctx context.Context, p *domain.Product) error {
	if p.ID == uuid.Nil {
		p.ID = newID()
	}
	const q = `
		INSERT INTO products (
			id, venue_id, location_id, main_category_id, image, name, external_id, size, price,
			origin_country, expiration, description, custom, localization, source
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		p.ID, uuidOrNil(p.VenueID), p.LocationID, p.MainCategoryID,
		p.Image, p.Name, p.ExternalID, p.Size,
		p.Price, p.Country, p.Expiration,
		p.Description, jsonOrNil(p.Custom), jsonOrNil(p.Localization), p.Source,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
}

func (r *productRepo) Update(ctx context.Context, p *domain.Product) error {
	const q = `
		UPDATE products SET
			location_id=$2, main_category_id=$3, image=$4, name=$5, external_id=$6, size=$7,
			price=$8, origin_country=$9, expiration=$10, description=$11,
			custom=$12, localization=$13, source=$14
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, q,
		p.ID, p.LocationID, p.MainCategoryID,
		p.Image, p.Name, p.ExternalID, p.Size,
		p.Price, p.Country, p.Expiration,
		p.Description, jsonOrNil(p.Custom), jsonOrNil(p.Localization), p.Source,
	).Scan(&p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("product not found")
	}
	return err
}

func (r *productRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "products", id.String())
}

func (r *productRepo) FindByCode(ctx context.Context, venueID uuid.UUID, code, source string) (*domain.Product, error) {
	const q = `
		SELECT id, venue_id, location_id, main_category_id, image, name, external_id, size, price,
		       origin_country, expiration, description, custom, localization, source,
		       created_at, updated_at, deleted_at
		FROM products WHERE venue_id = $1 AND external_id = $2 AND source = $3 AND deleted_at IS NULL`

	p, err := scanProduct(r.pool.QueryRow(ctx, q, venueID, code, source))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("product not found")
	}
	return p, err
}

func (r *productRepo) SetCategories(ctx context.Context, productID uuid.UUID, categoryIDs []uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM product_category_links WHERE product_id = $1`, productID)
	if err != nil {
		return err
	}
	for _, catID := range categoryIDs {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO product_category_links (product_id, category_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			productID, catID)
		if err != nil {
			return err
		}
	}
	return nil
}

// ── Attachments ───────────────────────────────────────────────────────────────

func (r *productRepo) ListAttachments(ctx context.Context, productID uuid.UUID) ([]*domain.ProductAttachment, error) {
	const q = `
		SELECT id, product_id, title, file_type, file, source_url, created_at, updated_at, deleted_at
		FROM product_attachments WHERE product_id = $1 AND deleted_at IS NULL ORDER BY file_type DESC`

	rows, err := r.pool.Query(ctx, q, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []*domain.ProductAttachment
	for rows.Next() {
		a, err := scanProductAttachment(rows)
		if err != nil {
			return nil, err
		}
		attachments = append(attachments, a)
	}
	return attachments, rows.Err()
}

func (r *productRepo) CreateAttachment(ctx context.Context, a *domain.ProductAttachment) error {
	if a.ID == uuid.Nil {
		a.ID = newID()
	}
	const q = `
		INSERT INTO product_attachments (id, product_id, title, file_type, file, source_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		a.ID, a.ProductID, a.Title, a.FileType, a.File, a.SourceURL,
	).Scan(&a.CreatedAt, &a.UpdatedAt)
}

func (r *productRepo) DeleteAttachment(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "product_attachments", id.String())
}

func (r *productRepo) SearchByName(ctx context.Context, venueID uuid.UUID, q string, limit int) ([]*domain.Product, error) {
	const query = `
		SELECT id, venue_id, location_id, main_category_id, image, name, external_id, size, price,
		       origin_country, expiration, description, custom, localization, source,
		       created_at, updated_at, deleted_at
		FROM products WHERE venue_id = $1 AND name ILIKE $2 AND deleted_at IS NULL
		ORDER BY name LIMIT $3`
	rows, err := r.pool.Query(ctx, query, venueID, "%"+q+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Product
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ── scan helpers ──────────────────────────────────────────────────────────────

func (r *productRepo) loadProductCategories(ctx context.Context, productID uuid.UUID) ([]*domain.ProductCategory, error) {
	const q = `
		SELECT pc.id, pc.venue_id, pc.external_id, pc.name, pc.source, pc.localization,
		       pc.created_at, pc.updated_at, pc.deleted_at
		FROM product_categories pc
		JOIN product_category_links pcl ON pcl.category_id = pc.id
		WHERE pcl.product_id = $1 AND pc.deleted_at IS NULL`

	rows, err := r.pool.Query(ctx, q, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []*domain.ProductCategory
	for rows.Next() {
		c, err := scanProductCategory(rows)
		if err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

func scanProductCategory(row pgx.Row) (*domain.ProductCategory, error) {
	var c domain.ProductCategory
	var localization []byte
	var deletedAt *time.Time
	err := row.Scan(
		&c.ID, &c.VenueID, &c.ExternalID, &c.Name, &c.Source, &localization,
		&c.CreatedAt, &c.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}
	c.Localization = json.RawMessage(localization)
	c.DeletedAt = deletedAt
	return &c, nil
}

func scanProduct(row pgx.Row) (*domain.Product, error) {
	var p domain.Product
	var custom, localization []byte
	var deletedAt *time.Time
	err := row.Scan(
		&p.ID, &p.VenueID, &p.LocationID, &p.MainCategoryID,
		&p.Image, &p.Name, &p.ExternalID, &p.Size, &p.Price,
		&p.Country, &p.Expiration, &p.Description, &custom, &localization, &p.Source,
		&p.CreatedAt, &p.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}
	p.Custom = json.RawMessage(custom)
	p.Localization = json.RawMessage(localization)
	p.DeletedAt = deletedAt
	return &p, nil
}

func scanProductAttachment(row pgx.Row) (*domain.ProductAttachment, error) {
	var a domain.ProductAttachment
	var deletedAt *time.Time
	err := row.Scan(
		&a.ID, &a.ProductID, &a.Title, &a.FileType, &a.File, &a.SourceURL,
		&a.CreatedAt, &a.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}
	a.DeletedAt = deletedAt
	return &a, nil
}
