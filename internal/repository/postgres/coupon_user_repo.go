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

type couponUserRepository struct{ pool *pgxpool.Pool }

func NewCouponUserRepository(pool *pgxpool.Pool) repository.CouponUserRepository {
	return &couponUserRepository{pool: pool}
}

func (r *couponUserRepository) Create(ctx context.Context, cu *domain.CouponUser) error {
	cu.ID = newID()
	return r.pool.QueryRow(ctx, `
		INSERT INTO coupon_users (id, coupon_id, user_id, device_id, is_used)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`,
		cu.ID, cu.CouponID, cu.UserID, cu.DeviceID, cu.IsUsed,
	).Scan(&cu.CreatedAt, &cu.UpdatedAt)
}

func (r *couponUserRepository) FindByCouponAndUser(ctx context.Context, couponID, userID uuid.UUID) (*domain.CouponUser, error) {
	var cu domain.CouponUser
	err := r.pool.QueryRow(ctx, `
		SELECT id, coupon_id, user_id, device_id, is_used, used_at, created_at, updated_at
		FROM coupon_users WHERE coupon_id = $1 AND user_id = $2`,
		couponID, userID,
	).Scan(&cu.ID, &cu.CouponID, &cu.UserID, &cu.DeviceID, &cu.IsUsed, &cu.UsedAt, &cu.CreatedAt, &cu.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewNotFound("coupon not assigned to this user")
		}
		return nil, err
	}
	return &cu, nil
}

func (r *couponUserRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE coupon_users SET is_used = TRUE, used_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *couponUserRepository) ExistsForVenueUser(ctx context.Context, venueID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM coupon_users cu
			INNER JOIN coupons c ON c.id = cu.coupon_id
			WHERE c.venue_id = $1 AND cu.user_id = $2 AND c.deleted_at IS NULL
		)`, venueID, userID,
	).Scan(&exists)
	return exists, err
}
