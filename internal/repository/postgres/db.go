package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
)

// newID returns a new UUIDv7 (time-ordered).
// Falls back to UUIDv4 only if the entropy source fails — which is essentially impossible.
func newID() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New()
	}
	return id
}

// pgErrCode extracts the PostgreSQL error code from an error, returning "" if not a PgError.
func pgErrCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

// IsUniqueViolation reports whether err is a PostgreSQL unique-constraint violation.
func IsUniqueViolation(err error) bool {
	return pgErrCode(err) == "23505"
}

// IsForeignKeyViolation reports whether err is a PostgreSQL foreign-key violation.
func IsForeignKeyViolation(err error) bool {
	return pgErrCode(err) == "23503"
}

// nullStr converts an empty string to nil for nullable TEXT columns.
func nullStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// derefStr sets dst to the value of src if src is non-nil.
func derefStr(dst *string, src *string) {
	if src != nil {
		*dst = *src
	}
}

// softDelete executes a soft-delete UPDATE on any table that has deleted_at.
func softDelete(ctx context.Context, pool *pgxpool.Pool, table, id string) error {
	sql := fmt.Sprintf(
		`UPDATE %s SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
		table,
	)
	tag, err := pool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
