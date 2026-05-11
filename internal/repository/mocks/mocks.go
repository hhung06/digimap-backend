// Package mocks provides testify/mock implementations of all repository interfaces.
// Import this package in service unit tests — never in production code.
package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/firebase"
	"github.com/hhung06/digimap-backend/internal/repository"
)

// ── UserRepository ────────────────────────────────────────────────────────────

// UserRepository is a mock implementation of repository.UserRepository.
type UserRepository struct{ mock.Mock }

func (m *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if u, ok := args.Get(0).(*domain.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if u, ok := args.Get(0).(*domain.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepository) Create(ctx context.Context, u *domain.User) error {
	return m.Called(ctx, u).Error(0)
}

func (m *UserRepository) Update(ctx context.Context, u *domain.User) error {
	return m.Called(ctx, u).Error(0)
}

func (m *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	return m.Called(ctx, id, passwordHash).Error(0)
}

func (m *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *UserRepository) GetVenueRole(ctx context.Context, venueID, userID uuid.UUID) (*domain.VenueUserRole, error) {
	args := m.Called(ctx, venueID, userID)
	if r, ok := args.Get(0).(*domain.VenueUserRole); ok {
		return r, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepository) ListVenueUsers(ctx context.Context, venueID uuid.UUID) ([]*domain.User, error) {
	args := m.Called(ctx, venueID)
	if users, ok := args.Get(0).([]*domain.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepository) UpsertVenueRole(ctx context.Context, r *domain.VenueUserRole) error {
	return m.Called(ctx, r).Error(0)
}

func (m *UserRepository) DeleteVenueRole(ctx context.Context, venueID, userID uuid.UUID) error {
	return m.Called(ctx, venueID, userID).Error(0)
}

func (m *UserRepository) CreateInvitation(ctx context.Context, inv *domain.VenueInvitation) error {
	return m.Called(ctx, inv).Error(0)
}

func (m *UserRepository) FindInvitationByID(ctx context.Context, id uuid.UUID) (*domain.VenueInvitation, error) {
	args := m.Called(ctx, id)
	if inv, ok := args.Get(0).(*domain.VenueInvitation); ok {
		return inv, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepository) FindInvitationByToken(ctx context.Context, token string) (*domain.VenueInvitation, error) {
	args := m.Called(ctx, token)
	if inv, ok := args.Get(0).(*domain.VenueInvitation); ok {
		return inv, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepository) UpdateInvitationStatus(ctx context.Context, id uuid.UUID, status domain.InvitationStatus, acceptedAt, cancelledAt *time.Time) error {
	return m.Called(ctx, id, status, acceptedAt, cancelledAt).Error(0)
}

func (m *UserRepository) ListInvitations(ctx context.Context, venueID uuid.UUID) ([]*domain.VenueInvitation, error) {
	args := m.Called(ctx, venueID)
	if invs, ok := args.Get(0).([]*domain.VenueInvitation); ok {
		return invs, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepository) DeleteInvitation(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// ── TokenRepository ───────────────────────────────────────────────────────────

// TokenRepository is a mock implementation of repository.TokenRepository.
type TokenRepository struct{ mock.Mock }

func (m *TokenRepository) CreateRefreshToken(ctx context.Context, t *domain.RefreshToken) error {
	return m.Called(ctx, t).Error(0)
}

func (m *TokenRepository) FindRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	args := m.Called(ctx, tokenHash)
	if t, ok := args.Get(0).(*domain.RefreshToken); ok {
		return t, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *TokenRepository) RevokeRefreshToken(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *TokenRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	return m.Called(ctx, userID).Error(0)
}

func (m *TokenRepository) CreateResetToken(ctx context.Context, t *domain.ResetPasswordToken) error {
	return m.Called(ctx, t).Error(0)
}

func (m *TokenRepository) FindResetToken(ctx context.Context, tokenHash string) (*domain.ResetPasswordToken, error) {
	args := m.Called(ctx, tokenHash)
	if t, ok := args.Get(0).(*domain.ResetPasswordToken); ok {
		return t, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *TokenRepository) MarkResetTokenUsed(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// ── EventLogRepository ────────────────────────────────────────────────────────

// EventLogRepository is a mock implementation of repository.EventLogRepository.
type EventLogRepository struct{ mock.Mock }

func (m *EventLogRepository) Create(ctx context.Context, e *domain.EventLog) error {
	return m.Called(ctx, e).Error(0)
}

func (m *EventLogRepository) ListByVenue(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.EventLog, int, error) {
	args := m.Called(ctx, venueID, p)
	if logs, ok := args.Get(0).([]*domain.EventLog); ok {
		return logs, args.Int(1), args.Error(2)
	}
	return nil, args.Int(1), args.Error(2)
}

// ── SearchQueryRepository ─────────────────────────────────────────────────────

// SearchQueryRepository is a mock implementation of repository.SearchQueryRepository.
type SearchQueryRepository struct{ mock.Mock }

func (m *SearchQueryRepository) Upsert(ctx context.Context, venueID uuid.UUID, term, origin, appID string) error {
	return m.Called(ctx, venueID, term, origin, appID).Error(0)
}

func (m *SearchQueryRepository) List(ctx context.Context, venueID uuid.UUID, filter repository.SearchQueryFilter, p domain.Pagination) ([]*domain.SearchQuery, int, error) {
	args := m.Called(ctx, venueID, filter, p)
	if qs, ok := args.Get(0).([]*domain.SearchQuery); ok {
		return qs, args.Int(1), args.Error(2)
	}
	return nil, args.Int(1), args.Error(2)
}

func (m *SearchQueryRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.SearchQuery, error) {
	args := m.Called(ctx, id)
	if q, ok := args.Get(0).(*domain.SearchQuery); ok {
		return q, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SearchQueryRepository) Create(ctx context.Context, q *domain.SearchQuery) error {
	return m.Called(ctx, q).Error(0)
}

func (m *SearchQueryRepository) Update(ctx context.Context, q *domain.SearchQuery) error {
	return m.Called(ctx, q).Error(0)
}

func (m *SearchQueryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// ── BeaconRepository ──────────────────────────────────────────────────────────

// BeaconRepository is a mock implementation of repository.BeaconRepository.
type BeaconRepository struct{ mock.Mock }

func (m *BeaconRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Beacon, error) {
	args := m.Called(ctx, id)
	if b, ok := args.Get(0).(*domain.Beacon); ok {
		return b, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *BeaconRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Beacon, int64, error) {
	args := m.Called(ctx, venueID, p)
	if bs, ok := args.Get(0).([]*domain.Beacon); ok {
		return bs, int64(args.Int(1)), args.Error(2)
	}
	return nil, int64(args.Int(1)), args.Error(2)
}

func (m *BeaconRepository) Create(ctx context.Context, b *domain.Beacon) error {
	return m.Called(ctx, b).Error(0)
}

func (m *BeaconRepository) Update(ctx context.Context, b *domain.Beacon) error {
	return m.Called(ctx, b).Error(0)
}

func (m *BeaconRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// ── CustomerRepository ────────────────────────────────────────────────────────

// CustomerRepository is a mock implementation of repository.CustomerRepository.
type CustomerRepository struct{ mock.Mock }

func (m *CustomerRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	args := m.Called(ctx, id)
	if c, ok := args.Get(0).(*domain.Customer); ok {
		return c, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CustomerRepository) List(ctx context.Context, p domain.Pagination) ([]*domain.Customer, int64, error) {
	args := m.Called(ctx, p)
	if cs, ok := args.Get(0).([]*domain.Customer); ok {
		return cs, int64(args.Int(1)), args.Error(2)
	}
	return nil, int64(args.Int(1)), args.Error(2)
}

func (m *CustomerRepository) Create(ctx context.Context, c *domain.Customer) error {
	return m.Called(ctx, c).Error(0)
}

func (m *CustomerRepository) Update(ctx context.Context, c *domain.Customer) error {
	return m.Called(ctx, c).Error(0)
}

func (m *CustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// ── AdvertisementRepository ───────────────────────────────────────────────────

// AdvertisementRepository is a mock implementation of repository.AdvertisementRepository.
type AdvertisementRepository struct{ mock.Mock }

func (m *AdvertisementRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Advertisement, int, error) {
	args := m.Called(ctx, venueID, p)
	if ads, ok := args.Get(0).([]*domain.Advertisement); ok {
		return ads, args.Int(1), args.Error(2)
	}
	return nil, args.Int(1), args.Error(2)
}

func (m *AdvertisementRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Advertisement, error) {
	args := m.Called(ctx, id)
	if a, ok := args.Get(0).(*domain.Advertisement); ok {
		return a, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AdvertisementRepository) Create(ctx context.Context, a *domain.Advertisement) error {
	return m.Called(ctx, a).Error(0)
}

func (m *AdvertisementRepository) Update(ctx context.Context, a *domain.Advertisement) error {
	return m.Called(ctx, a).Error(0)
}

func (m *AdvertisementRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// ── CouponRepository ──────────────────────────────────────────────────────────

// CouponRepository is a mock implementation of repository.CouponRepository.
type CouponRepository struct{ mock.Mock }

func (m *CouponRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Coupon, int, error) {
	args := m.Called(ctx, venueID, p)
	if cs, ok := args.Get(0).([]*domain.Coupon); ok {
		return cs, args.Int(1), args.Error(2)
	}
	return nil, args.Int(1), args.Error(2)
}

func (m *CouponRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Coupon, error) {
	args := m.Called(ctx, id)
	if c, ok := args.Get(0).(*domain.Coupon); ok {
		return c, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CouponRepository) Create(ctx context.Context, c *domain.Coupon) error {
	return m.Called(ctx, c).Error(0)
}

func (m *CouponRepository) Update(ctx context.Context, c *domain.Coupon) error {
	return m.Called(ctx, c).Error(0)
}

func (m *CouponRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *CouponRepository) Redeem(ctx context.Context, id uuid.UUID, appUserID uuid.UUID) error {
	return m.Called(ctx, id, appUserID).Error(0)
}

// ── VideoRepository ───────────────────────────────────────────────────────────

// VideoRepository is a mock implementation of repository.VideoRepository.
type VideoRepository struct{ mock.Mock }

func (m *VideoRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Video, int, error) {
	args := m.Called(ctx, venueID, p)
	if vs, ok := args.Get(0).([]*domain.Video); ok {
		return vs, args.Int(1), args.Error(2)
	}
	return nil, args.Int(1), args.Error(2)
}

func (m *VideoRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Video, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.Video); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *VideoRepository) Create(ctx context.Context, v *domain.Video) error {
	return m.Called(ctx, v).Error(0)
}

func (m *VideoRepository) Update(ctx context.Context, v *domain.Video) error {
	return m.Called(ctx, v).Error(0)
}

func (m *VideoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// ── EmailSender ───────────────────────────────────────────────────────────────

// EmailSender is a mock implementation of email.Sender.
type EmailSender struct{ mock.Mock }

func (m *EmailSender) SendPasswordReset(ctx context.Context, toEmail, resetToken string) error {
	return m.Called(ctx, toEmail, resetToken).Error(0)
}

func (m *EmailSender) SendInvitation(ctx context.Context, toEmail, venueName, inviteToken string) error {
	return m.Called(ctx, toEmail, venueName, inviteToken).Error(0)
}

// ── SnapshotRepository ────────────────────────────────────────────────────────

// SnapshotRepository is a mock implementation of repository.SnapshotRepository.
type SnapshotRepository struct{ mock.Mock }

func (m *SnapshotRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Snapshot, error) {
	args := m.Called(ctx, id)
	if s, ok := args.Get(0).(*domain.Snapshot); ok {
		return s, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SnapshotRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Snapshot, int64, error) {
	args := m.Called(ctx, venueID, p)
	if s, ok := args.Get(0).([]*domain.Snapshot); ok {
		return s, args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *SnapshotRepository) LatestPublished(ctx context.Context, venueID uuid.UUID) (*domain.Snapshot, error) {
	args := m.Called(ctx, venueID)
	if s, ok := args.Get(0).(*domain.Snapshot); ok {
		return s, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SnapshotRepository) Create(ctx context.Context, s *domain.Snapshot) error {
	return m.Called(ctx, s).Error(0)
}

func (m *SnapshotRepository) UpdateState(ctx context.Context, id uuid.UUID, state int, publishAt *time.Time) error {
	return m.Called(ctx, id, state, publishAt).Error(0)
}

func (m *SnapshotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *SnapshotRepository) CountDraftsByVenue(ctx context.Context, venueID uuid.UUID) (int64, error) {
	args := m.Called(ctx, venueID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *SnapshotRepository) DeleteOldestDraft(ctx context.Context, venueID uuid.UUID) error {
	return m.Called(ctx, venueID).Error(0)
}

func (m *SnapshotRepository) UnpublishVenue(ctx context.Context, venueID uuid.UUID) error {
	args := m.Called(ctx, venueID)
	return args.Error(0)
}

func (m *SnapshotRepository) LatestDraft(ctx context.Context, venueID uuid.UUID) (*domain.Snapshot, error) {
	args := m.Called(ctx, venueID)
	if v, ok := args.Get(0).(*domain.Snapshot); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

// ── LevelBundleRepository ─────────────────────────────────────────────────────

// LevelBundleRepository is a mock implementation of repository.LevelBundleRepository.
type LevelBundleRepository struct{ mock.Mock }

func (m *LevelBundleRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.LevelBundle, error) {
	args := m.Called(ctx, id)
	if b, ok := args.Get(0).(*domain.LevelBundle); ok {
		return b, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LevelBundleRepository) ListBySnapshot(ctx context.Context, snapshotID uuid.UUID) ([]*domain.LevelBundle, error) {
	args := m.Called(ctx, snapshotID)
	if b, ok := args.Get(0).([]*domain.LevelBundle); ok {
		return b, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LevelBundleRepository) Create(ctx context.Context, b *domain.LevelBundle) error {
	return m.Called(ctx, b).Error(0)
}

func (m *LevelBundleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// ── VenueRepository ───────────────────────────────────────────────────────────

// VenueRepository is a mock implementation of repository.VenueRepository.
type VenueRepository struct{ mock.Mock }

func (m *VenueRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Venue, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.Venue); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *VenueRepository) FindByPublicKey(ctx context.Context, publicKey string) (*domain.Venue, error) {
	args := m.Called(ctx, publicKey)
	if v, ok := args.Get(0).(*domain.Venue); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *VenueRepository) FindByPrivateKey(ctx context.Context, privateKey string) (*domain.Venue, error) {
	args := m.Called(ctx, privateKey)
	if v, ok := args.Get(0).(*domain.Venue); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *VenueRepository) List(ctx context.Context, customerID uuid.UUID, p domain.Pagination) ([]*domain.Venue, int64, error) {
	args := m.Called(ctx, customerID, p)
	if v, ok := args.Get(0).([]*domain.Venue); ok {
		return v, args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *VenueRepository) ListAll(ctx context.Context, p domain.Pagination) ([]*domain.Venue, int64, error) {
	args := m.Called(ctx, p)
	if v, ok := args.Get(0).([]*domain.Venue); ok {
		return v, args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *VenueRepository) Create(ctx context.Context, v *domain.Venue) error {
	return m.Called(ctx, v).Error(0)
}

func (m *VenueRepository) Update(ctx context.Context, v *domain.Venue) error {
	return m.Called(ctx, v).Error(0)
}

func (m *VenueRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *VenueRepository) UpdateKeys(ctx context.Context, id uuid.UUID, publicKey, privateKey string) error {
	return m.Called(ctx, id, publicKey, privateKey).Error(0)
}

func (m *VenueRepository) GetCustomerID(ctx context.Context, venueID uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, venueID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

// ── StorerMock ────────────────────────────────────────────────────────────────

// StorerMock is a testify/mock implementation of storage.Storer.
// Place it here so service tests can use it without importing the platform package.
type StorerMock struct{ mock.Mock }

func (m *StorerMock) PresignUpload(ctx context.Context, key, contentType string, ttl time.Duration) (string, error) {
	args := m.Called(ctx, key, contentType, ttl)
	return args.String(0), args.Error(1)
}

func (m *StorerMock) PresignDownload(ctx context.Context, key string, ttl time.Duration) (string, error) {
	args := m.Called(ctx, key, ttl)
	return args.String(0), args.Error(1)
}

func (m *StorerMock) PutObject(ctx context.Context, key string, data []byte) error {
	return m.Called(ctx, key, data).Error(0)
}

func (m *StorerMock) PutEncrypted(ctx context.Context, key string, body []byte, meta map[string]string) error {
	return m.Called(ctx, key, body, meta).Error(0)
}

func (m *StorerMock) GetObject(ctx context.Context, key string) ([]byte, error) {
	args := m.Called(ctx, key)
	if b, ok := args.Get(0).([]byte); ok {
		return b, args.Error(1)
	}
	return nil, args.Error(1)
}

// ── AssetRepository ───────────────────────────────────────────────────────────

type AssetRepository struct{ mock.Mock }

func (m *AssetRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	args := m.Called(ctx, id)
	if a, ok := args.Get(0).(*domain.Asset); ok {
		return a, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AssetRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Asset, int64, error) {
	args := m.Called(ctx, venueID, p)
	if a, ok := args.Get(0).([]*domain.Asset); ok {
		return a, args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *AssetRepository) Create(ctx context.Context, a *domain.Asset) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *AssetRepository) Update(ctx context.Context, a *domain.Asset) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *AssetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ── LevelTypeRepository ───────────────────────────────────────────────────────

type LevelTypeRepository struct{ mock.Mock }

func (m *LevelTypeRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.LevelType, error) {
	args := m.Called(ctx, id)
	if lt, ok := args.Get(0).(*domain.LevelType); ok {
		return lt, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LevelTypeRepository) List(ctx context.Context) ([]*domain.LevelType, error) {
	args := m.Called(ctx)
	if lt, ok := args.Get(0).([]*domain.LevelType); ok {
		return lt, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LevelTypeRepository) Create(ctx context.Context, lt *domain.LevelType) error {
	return m.Called(ctx, lt).Error(0)
}

func (m *LevelTypeRepository) Update(ctx context.Context, lt *domain.LevelType) error {
	return m.Called(ctx, lt).Error(0)
}

func (m *LevelTypeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// ── ThemeRepository ───────────────────────────────────────────────────────────

type ThemeRepository struct{ mock.Mock }

func (m *ThemeRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Theme, error) {
	args := m.Called(ctx, id)
	if t, ok := args.Get(0).(*domain.Theme); ok {
		return t, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ThemeRepository) List(ctx context.Context, venueID uuid.UUID) ([]*domain.Theme, error) {
	args := m.Called(ctx, venueID)
	if t, ok := args.Get(0).([]*domain.Theme); ok {
		return t, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ThemeRepository) Create(ctx context.Context, t *domain.Theme) error {
	return m.Called(ctx, t).Error(0)
}

func (m *ThemeRepository) Update(ctx context.Context, t *domain.Theme) error {
	return m.Called(ctx, t).Error(0)
}

func (m *ThemeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// ── ProductPlazaRepository ────────────────────────────────────────────────────

// ── LocationRepository ────────────────────────────────────────────────────────

type LocationRepository struct{ mock.Mock }

func (m *LocationRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Location, error) {
	args := m.Called(ctx, id)
	if l, ok := args.Get(0).(*domain.Location); ok {
		return l, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LocationRepository) List(ctx context.Context, venueID uuid.UUID, typeFilter *int, p domain.Pagination) ([]*domain.Location, int64, error) {
	args := m.Called(ctx, venueID, typeFilter, p)
	if locs, ok := args.Get(0).([]*domain.Location); ok {
		return locs, args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *LocationRepository) Create(ctx context.Context, l *domain.Location) error {
	return m.Called(ctx, l).Error(0)
}

func (m *LocationRepository) Update(ctx context.Context, l *domain.Location) error {
	return m.Called(ctx, l).Error(0)
}

func (m *LocationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *LocationRepository) SetCategories(ctx context.Context, locationID uuid.UUID, categoryIDs []uuid.UUID) error {
	return m.Called(ctx, locationID, categoryIDs).Error(0)
}

func (m *LocationRepository) SetTopLocation(ctx context.Context, id uuid.UUID, isTop bool, sortIndex *int) error {
	return m.Called(ctx, id, isTop, sortIndex).Error(0)
}

func (m *LocationRepository) ListTopLocations(ctx context.Context, venueID uuid.UUID) ([]*domain.Location, error) {
	args := m.Called(ctx, venueID)
	if items, ok := args.Get(0).([]*domain.Location); ok {
		return items, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LocationRepository) GeoSearch(ctx context.Context, lat, lng, radiusKm float64, venueID *uuid.UUID) ([]*domain.Location, error) {
	args := m.Called(ctx, lat, lng, radiusKm, venueID)
	if locs, ok := args.Get(0).([]*domain.Location); ok {
		return locs, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LocationRepository) ListImages(ctx context.Context, locationID uuid.UUID) ([]*domain.LocationImage, error) {
	args := m.Called(ctx, locationID)
	if imgs, ok := args.Get(0).([]*domain.LocationImage); ok {
		return imgs, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LocationRepository) CreateImage(ctx context.Context, img *domain.LocationImage) error {
	return m.Called(ctx, img).Error(0)
}

func (m *LocationRepository) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *LocationRepository) FindByExternalID(ctx context.Context, venueID uuid.UUID, externalID string, locationType int) (*domain.Location, error) {
	args := m.Called(ctx, venueID, externalID, locationType)
	if l, ok := args.Get(0).(*domain.Location); ok {
		return l, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LocationRepository) SearchByName(ctx context.Context, venueID uuid.UUID, q string, limit int) ([]*domain.Location, error) {
	args := m.Called(ctx, venueID, q, limit)
	if locs, ok := args.Get(0).([]*domain.Location); ok {
		return locs, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LocationRepository) ListMemos(ctx context.Context, venueID uuid.UUID) ([]*domain.Location, error) {
	args := m.Called(ctx, venueID)
	if locs, ok := args.Get(0).([]*domain.Location); ok {
		return locs, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LocationRepository) FindMemoByID(ctx context.Context, id uuid.UUID) (*domain.Location, error) {
	args := m.Called(ctx, id)
	if loc, ok := args.Get(0).(*domain.Location); ok {
		return loc, args.Error(1)
	}
	return nil, args.Error(1)
}

// ── ProductPlazaRepository ────────────────────────────────────────────────────

type ProductPlazaRepository struct{ mock.Mock }

func (m *ProductPlazaRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ProductPlaza, error) {
	args := m.Called(ctx, id)
	if p, ok := args.Get(0).(*domain.ProductPlaza); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductPlazaRepository) List(ctx context.Context, venueID uuid.UUID) ([]*domain.ProductPlaza, error) {
	args := m.Called(ctx, venueID)
	if p, ok := args.Get(0).([]*domain.ProductPlaza); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductPlazaRepository) Create(ctx context.Context, p *domain.ProductPlaza) error {
	return m.Called(ctx, p).Error(0)
}

func (m *ProductPlazaRepository) Update(ctx context.Context, p *domain.ProductPlaza) error {
	return m.Called(ctx, p).Error(0)
}

func (m *ProductPlazaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// ── ProductRepository ─────────────────────────────────────────────────────────

type ProductRepository struct{ mock.Mock }

func (m *ProductRepository) FindCategoryByID(ctx context.Context, id uuid.UUID) (*domain.ProductCategory, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.ProductCategory); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductRepository) ListCategories(ctx context.Context, venueID uuid.UUID) ([]*domain.ProductCategory, error) {
	args := m.Called(ctx, venueID)
	if v, ok := args.Get(0).([]*domain.ProductCategory); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductRepository) CreateCategory(ctx context.Context, c *domain.ProductCategory) error {
	return m.Called(ctx, c).Error(0)
}

func (m *ProductRepository) UpdateCategory(ctx context.Context, c *domain.ProductCategory) error {
	return m.Called(ctx, c).Error(0)
}

func (m *ProductRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *ProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.Product); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Product, int64, error) {
	args := m.Called(ctx, venueID, p)
	if v, ok := args.Get(0).([]*domain.Product); ok {
		return v, args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *ProductRepository) Create(ctx context.Context, p *domain.Product) error {
	return m.Called(ctx, p).Error(0)
}

func (m *ProductRepository) Update(ctx context.Context, p *domain.Product) error {
	return m.Called(ctx, p).Error(0)
}

func (m *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *ProductRepository) SetCategories(ctx context.Context, productID uuid.UUID, categoryIDs []uuid.UUID) error {
	return m.Called(ctx, productID, categoryIDs).Error(0)
}

func (m *ProductRepository) FindByCode(ctx context.Context, venueID uuid.UUID, code, source string) (*domain.Product, error) {
	args := m.Called(ctx, venueID, code, source)
	if v, ok := args.Get(0).(*domain.Product); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductRepository) ListAttachments(ctx context.Context, productID uuid.UUID) ([]*domain.ProductAttachment, error) {
	args := m.Called(ctx, productID)
	if v, ok := args.Get(0).([]*domain.ProductAttachment); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductRepository) CreateAttachment(ctx context.Context, a *domain.ProductAttachment) error {
	return m.Called(ctx, a).Error(0)
}

func (m *ProductRepository) DeleteAttachment(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *ProductRepository) SearchByName(ctx context.Context, venueID uuid.UUID, q string, limit int) ([]*domain.Product, error) {
	args := m.Called(ctx, venueID, q, limit)
	if v, ok := args.Get(0).([]*domain.Product); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

// ── AppUserRepository ─────────────────────────────────────────────────────────

type AppUserRepository struct{ mock.Mock }

func (m *AppUserRepository) FindByToken(ctx context.Context, venueID uuid.UUID, token string) (*domain.AppUser, error) {
	args := m.Called(ctx, venueID, token)
	if v, ok := args.Get(0).(*domain.AppUser); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AppUserRepository) Create(ctx context.Context, u *domain.AppUser) error {
	return m.Called(ctx, u).Error(0)
}

// ── LanguageRepository ────────────────────────────────────────────────────────

type LanguageRepository struct{ mock.Mock }

func (m *LanguageRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Language, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.Language); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LanguageRepository) List(ctx context.Context, venueID uuid.UUID) ([]*domain.Language, error) {
	args := m.Called(ctx, venueID)
	if v, ok := args.Get(0).([]*domain.Language); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LanguageRepository) Create(ctx context.Context, l *domain.Language) error {
	return m.Called(ctx, l).Error(0)
}

func (m *LanguageRepository) Update(ctx context.Context, l *domain.Language) error {
	return m.Called(ctx, l).Error(0)
}

func (m *LanguageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *LanguageRepository) ListEnabled(ctx context.Context, venueID uuid.UUID) ([]*domain.Language, error) {
	args := m.Called(ctx, venueID)
	if v, ok := args.Get(0).([]*domain.Language); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

// ── AppVersionRepository ──────────────────────────────────────────────────────

type AppVersionRepository struct{ mock.Mock }

func (m *AppVersionRepository) Upsert(ctx context.Context, venueID, version uuid.UUID) error {
	return m.Called(ctx, venueID, version).Error(0)
}

func (m *AppVersionRepository) Get(ctx context.Context, venueID uuid.UUID) (*domain.AppVersion, error) {
	args := m.Called(ctx, venueID)
	if v, ok := args.Get(0).(*domain.AppVersion); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

// ── LevelRepository ───────────────────────────────────────────────────────────

type LevelRepository struct{ mock.Mock }

func (m *LevelRepository) FindMapGroupByID(ctx context.Context, id uuid.UUID) (*domain.MapGroup, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.MapGroup); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LevelRepository) ListMapGroups(ctx context.Context, venueID uuid.UUID) ([]*domain.MapGroup, error) {
	args := m.Called(ctx, venueID)
	if v, ok := args.Get(0).([]*domain.MapGroup); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LevelRepository) CreateMapGroup(ctx context.Context, mg *domain.MapGroup) error {
	return m.Called(ctx, mg).Error(0)
}

func (m *LevelRepository) UpdateMapGroup(ctx context.Context, mg *domain.MapGroup) error {
	return m.Called(ctx, mg).Error(0)
}

func (m *LevelRepository) DeleteMapGroup(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *LevelRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Level, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.Level); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LevelRepository) List(ctx context.Context, venueID uuid.UUID) ([]*domain.Level, error) {
	args := m.Called(ctx, venueID)
	if v, ok := args.Get(0).([]*domain.Level); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LevelRepository) Create(ctx context.Context, l *domain.Level) error {
	return m.Called(ctx, l).Error(0)
}

func (m *LevelRepository) Update(ctx context.Context, l *domain.Level) error {
	return m.Called(ctx, l).Error(0)
}

func (m *LevelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *LevelRepository) UpsertPerspective(ctx context.Context, p *domain.Perspective) error {
	return m.Called(ctx, p).Error(0)
}

func (m *LevelRepository) ListGeoReferences(ctx context.Context, levelID uuid.UUID) ([]*domain.GeoReference, error) {
	args := m.Called(ctx, levelID)
	if v, ok := args.Get(0).([]*domain.GeoReference); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *LevelRepository) CreateGeoReference(ctx context.Context, g *domain.GeoReference) error {
	return m.Called(ctx, g).Error(0)
}

func (m *LevelRepository) DeleteGeoReference(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// ── NotificationRepository ────────────────────────────────────────────────────

type NotificationRepository struct{ mock.Mock }

func (m *NotificationRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.Notification); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *NotificationRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Notification, int64, error) {
	args := m.Called(ctx, venueID, p)
	if v, ok := args.Get(0).([]*domain.Notification); ok {
		return v, args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *NotificationRepository) ListDueScheduled(ctx context.Context, now time.Time) ([]*domain.Notification, error) {
	args := m.Called(ctx, now)
	if v, ok := args.Get(0).([]*domain.Notification); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *NotificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	return m.Called(ctx, n).Error(0)
}

func (m *NotificationRepository) Update(ctx context.Context, n *domain.Notification) error {
	return m.Called(ctx, n).Error(0)
}

func (m *NotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *NotificationRepository) MarkSent(ctx context.Context, id uuid.UUID, publishedAt time.Time) error {
	return m.Called(ctx, id, publishedAt).Error(0)
}

func (m *NotificationRepository) MarkFailed(ctx context.Context, id uuid.UUID, errInfos []byte) error {
	return m.Called(ctx, id, errInfos).Error(0)
}

// ── SurveyRepository ──────────────────────────────────────────────────────────

type SurveyRepository struct{ mock.Mock }

func (m *SurveyRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Survey, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.Survey); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SurveyRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Survey, int64, error) {
	args := m.Called(ctx, venueID, p)
	if v, ok := args.Get(0).([]*domain.Survey); ok {
		return v, args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *SurveyRepository) ListDueActivation(ctx context.Context, now time.Time) ([]*domain.Survey, error) {
	args := m.Called(ctx, now)
	if v, ok := args.Get(0).([]*domain.Survey); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SurveyRepository) ListDueClosure(ctx context.Context, now time.Time) ([]*domain.Survey, error) {
	args := m.Called(ctx, now)
	if v, ok := args.Get(0).([]*domain.Survey); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SurveyRepository) Create(ctx context.Context, s *domain.Survey) error {
	return m.Called(ctx, s).Error(0)
}

func (m *SurveyRepository) Update(ctx context.Context, s *domain.Survey) error {
	return m.Called(ctx, s).Error(0)
}

func (m *SurveyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *SurveyRepository) FindQuestionByID(ctx context.Context, id uuid.UUID) (*domain.Question, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.Question); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SurveyRepository) ListQuestions(ctx context.Context, surveyID uuid.UUID) ([]*domain.Question, error) {
	args := m.Called(ctx, surveyID)
	if v, ok := args.Get(0).([]*domain.Question); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SurveyRepository) CreateQuestion(ctx context.Context, q *domain.Question) error {
	return m.Called(ctx, q).Error(0)
}

func (m *SurveyRepository) UpdateQuestion(ctx context.Context, q *domain.Question) error {
	return m.Called(ctx, q).Error(0)
}

func (m *SurveyRepository) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *SurveyRepository) CreateOption(ctx context.Context, o *domain.Option) error {
	return m.Called(ctx, o).Error(0)
}

func (m *SurveyRepository) UpdateOption(ctx context.Context, o *domain.Option) error {
	return m.Called(ctx, o).Error(0)
}

func (m *SurveyRepository) DeleteOption(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *SurveyRepository) ListActive(ctx context.Context, venueID uuid.UUID, publishTypes []int) ([]*domain.Survey, error) {
	args := m.Called(ctx, venueID, publishTypes)
	if v, ok := args.Get(0).([]*domain.Survey); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SurveyRepository) ListResponses(ctx context.Context, surveyID uuid.UUID, p domain.Pagination) ([]*domain.SurveyResponse, int64, error) {
	args := m.Called(ctx, surveyID, p)
	if v, ok := args.Get(0).([]*domain.SurveyResponse); ok {
		return v, args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *SurveyRepository) CreateResponse(ctx context.Context, r *domain.SurveyResponse) error {
	return m.Called(ctx, r).Error(0)
}

// ── EventRepository ───────────────────────────────────────────────────────────

type EventRepository struct{ mock.Mock }

func (m *EventRepository) FindTagByID(ctx context.Context, id uuid.UUID) (*domain.EventTag, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.EventTag); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *EventRepository) ListTags(ctx context.Context) ([]*domain.EventTag, error) {
	args := m.Called(ctx)
	if v, ok := args.Get(0).([]*domain.EventTag); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *EventRepository) CreateTag(ctx context.Context, t *domain.EventTag) error {
	return m.Called(ctx, t).Error(0)
}

func (m *EventRepository) UpdateTag(ctx context.Context, t *domain.EventTag) error {
	return m.Called(ctx, t).Error(0)
}

func (m *EventRepository) DeleteTag(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *EventRepository) FindEventTypeByID(ctx context.Context, id uuid.UUID) (*domain.EventType, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.EventType); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *EventRepository) ListEventTypes(ctx context.Context, venueID uuid.UUID) ([]*domain.EventType, error) {
	args := m.Called(ctx, venueID)
	if v, ok := args.Get(0).([]*domain.EventType); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *EventRepository) CreateEventType(ctx context.Context, t *domain.EventType) error {
	return m.Called(ctx, t).Error(0)
}

func (m *EventRepository) UpdateEventType(ctx context.Context, t *domain.EventType) error {
	return m.Called(ctx, t).Error(0)
}

func (m *EventRepository) DeleteEventType(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *EventRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.Event); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *EventRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Event, int64, error) {
	args := m.Called(ctx, venueID, p)
	if v, ok := args.Get(0).([]*domain.Event); ok {
		return v, args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *EventRepository) Create(ctx context.Context, e *domain.Event) error {
	return m.Called(ctx, e).Error(0)
}

func (m *EventRepository) Update(ctx context.Context, e *domain.Event) error {
	return m.Called(ctx, e).Error(0)
}

func (m *EventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *EventRepository) SetTags(ctx context.Context, eventID uuid.UUID, tagIDs []uuid.UUID) error {
	return m.Called(ctx, eventID, tagIDs).Error(0)
}

func (m *EventRepository) SetLocations(ctx context.Context, eventID uuid.UUID, locationIDs []uuid.UUID) error {
	return m.Called(ctx, eventID, locationIDs).Error(0)
}

func (m *EventRepository) CreateImage(ctx context.Context, img *domain.EventImage) error {
	return m.Called(ctx, img).Error(0)
}

func (m *EventRepository) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// ── ConnectionRepository ──────────────────────────────────────────────────────

type ConnectionRepository struct{ mock.Mock }

func (m *ConnectionRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Connection, int, error) {
	args := m.Called(ctx, venueID, p)
	if v, ok := args.Get(0).([]*domain.Connection); ok {
		return v, args.Int(1), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *ConnectionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Connection, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.Connection); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ConnectionRepository) Create(ctx context.Context, c *domain.Connection) error {
	return m.Called(ctx, c).Error(0)
}

func (m *ConnectionRepository) Update(ctx context.Context, c *domain.Connection) error {
	return m.Called(ctx, c).Error(0)
}

func (m *ConnectionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *ConnectionRepository) ListLevels(ctx context.Context, connectionID uuid.UUID) ([]*domain.ConnectionLevel, error) {
	args := m.Called(ctx, connectionID)
	if v, ok := args.Get(0).([]*domain.ConnectionLevel); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ConnectionRepository) AddLevel(ctx context.Context, cl *domain.ConnectionLevel) error {
	return m.Called(ctx, cl).Error(0)
}

func (m *ConnectionRepository) RemoveLevel(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// ── ArticleRepository ─────────────────────────────────────────────────────────

type ArticleRepository struct{ mock.Mock }

func (m *ArticleRepository) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Article, int, error) {
	args := m.Called(ctx, venueID, p)
	if v, ok := args.Get(0).([]*domain.Article); ok {
		return v, args.Int(1), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *ArticleRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Article, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.Article); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ArticleRepository) Create(ctx context.Context, a *domain.Article) error {
	return m.Called(ctx, a).Error(0)
}

func (m *ArticleRepository) Update(ctx context.Context, a *domain.Article) error {
	return m.Called(ctx, a).Error(0)
}

func (m *ArticleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *ArticleRepository) CreateImage(ctx context.Context, img *domain.ArticleImage) error {
	return m.Called(ctx, img).Error(0)
}

func (m *ArticleRepository) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// ── TagRepository ─────────────────────────────────────────────────────────────

type TagRepository struct{ mock.Mock }

func (m *TagRepository) List(ctx context.Context, p domain.Pagination) ([]*domain.Tag, int, error) {
	args := m.Called(ctx, p)
	if v, ok := args.Get(0).([]*domain.Tag); ok {
		return v, args.Int(1), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *TagRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Tag, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*domain.Tag); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *TagRepository) Create(ctx context.Context, t *domain.Tag) error {
	return m.Called(ctx, t).Error(0)
}

func (m *TagRepository) Update(ctx context.Context, t *domain.Tag) error {
	return m.Called(ctx, t).Error(0)
}

func (m *TagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *TagRepository) AttachTag(ctx context.Context, et *domain.EntityTag) error {
	return m.Called(ctx, et).Error(0)
}

func (m *TagRepository) DetachTag(ctx context.Context, tagID uuid.UUID, entityType string, entityID uuid.UUID) error {
	return m.Called(ctx, tagID, entityType, entityID).Error(0)
}

func (m *TagRepository) ListEntityTags(ctx context.Context, entityType string, entityID uuid.UUID) ([]*domain.Tag, error) {
	args := m.Called(ctx, entityType, entityID)
	if v, ok := args.Get(0).([]*domain.Tag); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

// ── MockPusher ────────────────────────────────────────────────────────────────

// MockPusher is a testify/mock implementation of firebase.Pusher.
type MockPusher struct{ mock.Mock }

func (m *MockPusher) Send(ctx context.Context, msg firebase.Message) (string, error) {
	args := m.Called(ctx, msg)
	return args.String(0), args.Error(1)
}

func (m *MockPusher) SendMulticast(ctx context.Context, tokens []string, title, body string, data map[string]string) (int, int, error) {
	args := m.Called(ctx, tokens, title, body, data)
	return args.Int(0), args.Int(1), args.Error(2)
}
