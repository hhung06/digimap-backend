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

type tokenRepo struct {
	pool *pgxpool.Pool
}

// NewTokenRepository returns a TokenRepository backed by PostgreSQL.
func NewTokenRepository(pool *pgxpool.Pool) repository.TokenRepository {
	return &tokenRepo{pool: pool}
}

// ── Refresh tokens ────────────────────────────────────────────────────────────

func (r *tokenRepo) CreateRefreshToken(ctx context.Context, t *domain.RefreshToken) error {
	if t.ID == uuid.Nil {
		t.ID = newID()
	}
	const q = `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`

	return r.pool.QueryRow(ctx, q, t.ID, t.UserID, t.TokenHash, t.ExpiresAt).Scan(&t.CreatedAt)
}

func (r *tokenRepo) FindRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	const q = `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW()`

	var t domain.RefreshToken
	err := r.pool.QueryRow(ctx, q, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.RevokedAt, &t.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewUnauthorized("refresh token invalid or expired")
	}
	return &t, err
}

func (r *tokenRepo) RevokeRefreshToken(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1`,
		id,
	)
	return err
}

func (r *tokenRepo) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`,
		userID,
	)
	return err
}

// ── Password reset tokens ─────────────────────────────────────────────────────

func (r *tokenRepo) CreateResetToken(ctx context.Context, t *domain.ResetPasswordToken) error {
	if t.ID == uuid.Nil {
		t.ID = newID()
	}
	const q = `
		INSERT INTO reset_password_tokens (id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`

	return r.pool.QueryRow(ctx, q, t.ID, t.UserID, t.TokenHash, t.ExpiresAt).Scan(&t.CreatedAt)
}

func (r *tokenRepo) FindResetToken(ctx context.Context, tokenHash string) (*domain.ResetPasswordToken, error) {
	const q = `
		SELECT id, user_id, token_hash, expires_at, used_at, created_at
		FROM reset_password_tokens
		WHERE token_hash = $1 AND used_at IS NULL AND expires_at > NOW()`

	var t domain.ResetPasswordToken
	err := r.pool.QueryRow(ctx, q, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewBadRequest("reset token invalid or expired")
	}
	return &t, err
}

func (r *tokenRepo) MarkResetTokenUsed(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE reset_password_tokens SET used_at = NOW() WHERE id = $1`,
		id,
	)
	return err
}
