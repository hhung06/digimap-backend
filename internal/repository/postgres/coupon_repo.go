package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type couponRepository struct{ pool *pgxpool.Pool }

func NewCouponRepository(pool *pgxpool.Pool) *couponRepository {
	return &couponRepository{pool: pool}
}

const couponSelectCols = `
    c.id, c.venue_id, c.external_id, c.coupon_name, c.coupon_code,
    c.status, c.issued_at, c.expired_at, c.redeemed_at, c.redeemed_by,
    c.localization, c.created_at, c.updated_at`

func (r *couponRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Coupon, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM coupons WHERE venue_id=$1 AND deleted_at IS NULL`, venueID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx,
		`SELECT `+couponSelectCols+` FROM coupons c WHERE c.venue_id=$1 AND c.deleted_at IS NULL
         ORDER BY c.created_at DESC LIMIT $2 OFFSET $3`,
		venueID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []*domain.Coupon
	for rows.Next() {
		c, err := scanCoupon(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, c)
	}
	return out, total, rows.Err()
}

func (r *couponRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Coupon, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+couponSelectCols+` FROM coupons c WHERE c.id=$1 AND c.deleted_at IS NULL`, id)
	return scanCoupon(row)
}

func (r *couponRepository) Create(ctx context.Context, c *domain.Coupon) error {
	c.ID = newID()
	return r.pool.QueryRow(ctx,
		`INSERT INTO coupons (id,venue_id,external_id,coupon_name,coupon_code,status,issued_at,expired_at,localization)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
         RETURNING created_at,updated_at`,
		c.ID, c.VenueID, c.ExternalID, c.CouponName, c.CouponCode,
		c.Status, c.IssuedAt, c.ExpiredAt, c.Localization,
	).Scan(&c.CreatedAt, &c.UpdatedAt)
}

func (r *couponRepository) Update(ctx context.Context, c *domain.Coupon) error {
	return r.pool.QueryRow(ctx,
		`UPDATE coupons SET external_id=$2,coupon_name=$3,coupon_code=$4,status=$5,
         issued_at=$6,expired_at=$7,localization=$8
         WHERE id=$1 AND deleted_at IS NULL RETURNING updated_at`,
		c.ID, c.ExternalID, c.CouponName, c.CouponCode, c.Status,
		c.IssuedAt, c.ExpiredAt, c.Localization,
	).Scan(&c.UpdatedAt)
}

func (r *couponRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE coupons SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func (r *couponRepository) Redeem(ctx context.Context, id uuid.UUID, appUserID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE coupons SET redeemed_at=NOW(), redeemed_by=$2 WHERE id=$1 AND deleted_at IS NULL`,
		id, appUserID)
	return err
}

func scanCoupon(row scanner) (*domain.Coupon, error) {
	var c domain.Coupon
	if err := row.Scan(
		&c.ID, &c.VenueID, &c.ExternalID, &c.CouponName, &c.CouponCode,
		&c.Status, &c.IssuedAt, &c.ExpiredAt, &c.RedeemedAt, &c.RedeemedBy,
		&c.Localization, &c.CreatedAt, &c.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &c, nil
}
