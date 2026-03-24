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

type userRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepository returns a UserRepository backed by PostgreSQL.
func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &userRepo{pool: pool}
}

// ── User CRUD ─────────────────────────────────────────────────────────────────

func (r *userRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, first_name, last_name, phone, avatar_url,
		       is_active, is_system_admin, last_login_at, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL`

	u, err := scanUser(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("user not found")
	}
	return u, err
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, first_name, last_name, phone, avatar_url,
		       is_active, is_system_admin, last_login_at, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL`

	u, err := scanUser(r.pool.QueryRow(ctx, q, email))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("user not found")
	}
	return u, err
}

func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	const q = `
		INSERT INTO users (id, email, password_hash, first_name, last_name, phone, avatar_url,
		                   is_active, is_system_admin)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at`

	if u.ID == uuid.Nil {
		u.ID = newID()
	}
	err := r.pool.QueryRow(ctx, q,
		u.ID, u.Email, u.PasswordHash, u.FirstName, u.LastName,
		u.Phone, u.AvatarURL, u.IsActive, u.IsSystemAdmin,
	).Scan(&u.CreatedAt, &u.UpdatedAt)

	if IsUniqueViolation(err) {
		return domain.NewConflict("email already registered")
	}
	return err
}

func (r *userRepo) Update(ctx context.Context, u *domain.User) error {
	const q = `
		UPDATE users
		SET first_name = $2, last_name = $3, phone = $4, avatar_url = $5, is_active = $6
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, q,
		u.ID, u.FirstName, u.LastName, u.Phone, u.AvatarURL, u.IsActive,
	).Scan(&u.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("user not found")
	}
	return err
}

func (r *userRepo) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET last_login_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	return err
}

// ── Venue roles ───────────────────────────────────────────────────────────────

func (r *userRepo) GetVenueRole(ctx context.Context, venueID, userID uuid.UUID) (*domain.VenueUserRole, error) {
	const q = `
		SELECT id, venue_id, user_id, role, created_at, updated_at, deleted_at
		FROM venue_user_roles
		WHERE venue_id = $1 AND user_id = $2 AND deleted_at IS NULL`

	var vr domain.VenueUserRole
	var roleStr string
	err := r.pool.QueryRow(ctx, q, venueID, userID).Scan(
		&vr.ID, &vr.VenueID, &vr.UserID, &roleStr,
		&vr.CreatedAt, &vr.UpdatedAt, &vr.DeletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("venue role not found")
	}
	if err != nil {
		return nil, err
	}
	vr.Role = domain.RoleFromString(roleStr)
	return &vr, nil
}

func (r *userRepo) UpsertVenueRole(ctx context.Context, vr *domain.VenueUserRole) error {
	if vr.ID == uuid.Nil {
		vr.ID = newID()
	}
	const q = `
		INSERT INTO venue_user_roles (id, venue_id, user_id, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (venue_id, user_id)
		DO UPDATE SET role = EXCLUDED.role, deleted_at = NULL, updated_at = NOW()
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		vr.ID, vr.VenueID, vr.UserID, vr.Role.String(),
	).Scan(&vr.CreatedAt, &vr.UpdatedAt)
}

func (r *userRepo) DeleteVenueRole(ctx context.Context, venueID, userID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE venue_user_roles SET deleted_at = NOW(), updated_at = NOW()
		 WHERE venue_id = $1 AND user_id = $2 AND deleted_at IS NULL`,
		venueID, userID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.NewNotFound("venue role not found")
	}
	return nil
}

func (r *userRepo) ListVenueUsers(ctx context.Context, venueID uuid.UUID) ([]*domain.User, error) {
	const q = `
		SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.phone, u.avatar_url,
		       u.is_active, u.is_system_admin, u.last_login_at, u.created_at, u.updated_at, u.deleted_at
		FROM users u
		JOIN venue_user_roles vr ON vr.user_id = u.id
		WHERE vr.venue_id = $1 AND vr.deleted_at IS NULL AND u.deleted_at IS NULL
		ORDER BY u.created_at`

	rows, err := r.pool.Query(ctx, q, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// ── Invitations ───────────────────────────────────────────────────────────────

func (r *userRepo) CreateInvitation(ctx context.Context, inv *domain.VenueInvitation) error {
	if inv.ID == uuid.Nil {
		inv.ID = newID()
	}
	const q = `
		INSERT INTO venue_invitations (id, venue_id, email, role, token, invited_by, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at`

	return r.pool.QueryRow(ctx, q,
		inv.ID, inv.VenueID, inv.Email, inv.Role.String(),
		inv.Token, inv.InvitedBy, inv.ExpiresAt,
	).Scan(&inv.CreatedAt)
}

func (r *userRepo) FindInvitationByToken(ctx context.Context, token string) (*domain.VenueInvitation, error) {
	const q = `
		SELECT id, venue_id, email, role, token, invited_by, status, accepted_at, cancelled_at, expires_at, created_at, deleted_at
		FROM venue_invitations
		WHERE token = $1 AND deleted_at IS NULL`
	return scanInvitationRow(r.pool.QueryRow(ctx, q, token))
}

func (r *userRepo) FindInvitationByID(ctx context.Context, id uuid.UUID) (*domain.VenueInvitation, error) {
	const q = `
		SELECT id, venue_id, email, role, token, invited_by, status, accepted_at, cancelled_at, expires_at, created_at, deleted_at
		FROM venue_invitations
		WHERE id = $1 AND deleted_at IS NULL`
	return scanInvitationRow(r.pool.QueryRow(ctx, q, id))
}

func (r *userRepo) UpdateInvitationStatus(ctx context.Context, id uuid.UUID, status domain.InvitationStatus, acceptedAt, cancelledAt *time.Time) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE venue_invitations SET status = $2, accepted_at = $3, cancelled_at = $4 WHERE id = $1 AND deleted_at IS NULL`,
		id, string(status), acceptedAt, cancelledAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.NewNotFound("invitation not found")
	}
	return nil
}

func (r *userRepo) ListInvitations(ctx context.Context, venueID uuid.UUID) ([]*domain.VenueInvitation, error) {
	const q = `
		SELECT id, venue_id, email, role, token, invited_by, status, accepted_at, cancelled_at, expires_at, created_at, deleted_at
		FROM venue_invitations
		WHERE venue_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, q, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invs []*domain.VenueInvitation
	for rows.Next() {
		inv, err := scanInvitationRow(rows)
		if err != nil {
			return nil, err
		}
		invs = append(invs, inv)
	}
	return invs, rows.Err()
}

func (r *userRepo) DeleteInvitation(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "venue_invitations", id.String())
}

// ── helpers ───────────────────────────────────────────────────────────────────

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(row scanner) (*domain.User, error) {
	var u domain.User
	var phone, avatarURL *string
	var lastLogin *time.Time
	err := row.Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
		&phone, &avatarURL, &u.IsActive, &u.IsSystemAdmin,
		&lastLogin, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	if phone != nil {
		u.Phone = *phone
	}
	if avatarURL != nil {
		u.AvatarURL = *avatarURL
	}
	u.LastLoginAt = lastLogin
	return &u, nil
}

func scanInvitationRow(row scanner) (*domain.VenueInvitation, error) {
	var inv domain.VenueInvitation
	var roleStr, statusStr string
	err := row.Scan(
		&inv.ID, &inv.VenueID, &inv.Email, &roleStr, &inv.Token,
		&inv.InvitedBy, &statusStr, &inv.AcceptedAt, &inv.CancelledAt,
		&inv.ExpiresAt, &inv.CreatedAt, &inv.DeletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("invitation not found")
	}
	if err != nil {
		return nil, err
	}
	inv.Role = domain.RoleFromString(roleStr)
	inv.Status = domain.InvitationStatus(statusStr)
	return &inv, nil
}
