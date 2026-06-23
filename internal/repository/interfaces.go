package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hhung06/digimap-backend/internal/domain"
)

// UserRepository handles persistence for User and VenueUserRole entities.
type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, u *domain.User) error
	Update(ctx context.Context, u *domain.User) error
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error

	// Venue roles
	GetVenueRole(ctx context.Context, venueID, userID uuid.UUID) (*domain.VenueUserRole, error)
	ListVenueUsers(ctx context.Context, venueID uuid.UUID) ([]*domain.User, error)
	UpsertVenueRole(ctx context.Context, r *domain.VenueUserRole) error
	DeleteVenueRole(ctx context.Context, venueID, userID uuid.UUID) error

	// Invitations
	CreateInvitation(ctx context.Context, inv *domain.VenueInvitation) error
	FindInvitationByID(ctx context.Context, id uuid.UUID) (*domain.VenueInvitation, error)
	FindInvitationByToken(ctx context.Context, token string) (*domain.VenueInvitation, error)
	UpdateInvitationStatus(ctx context.Context, id uuid.UUID, status domain.InvitationStatus, acceptedAt, cancelledAt *time.Time) error
	ListInvitations(ctx context.Context, venueID uuid.UUID) ([]*domain.VenueInvitation, error)
	DeleteInvitation(ctx context.Context, id uuid.UUID) error
}

// CustomerRepository handles persistence for Customer entities.
type CustomerRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	List(ctx context.Context, p domain.Pagination) ([]*domain.Customer, int64, error)
	Create(ctx context.Context, c *domain.Customer) error
	Update(ctx context.Context, c *domain.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// VenueRepository handles persistence for Venue entities.
type VenueRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Venue, error)
	FindByPublicKey(ctx context.Context, publicKey string) (*domain.Venue, error)
	FindByPrivateKey(ctx context.Context, privateKey string) (*domain.Venue, error)
	List(ctx context.Context, customerID uuid.UUID, p domain.Pagination) ([]*domain.Venue, int64, error)
	ListAll(ctx context.Context, p domain.Pagination) ([]*domain.Venue, int64, error)
	Create(ctx context.Context, v *domain.Venue) error
	Update(ctx context.Context, v *domain.Venue) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateKeys(ctx context.Context, id uuid.UUID, publicKey, privateKey string) error
	SetThemeID(ctx context.Context, venueID uuid.UUID, themeID *uuid.UUID) error

	GetCustomerID(ctx context.Context, venueID uuid.UUID) (uuid.UUID, error)
}

// LevelRepository handles persistence for Level, MapGroup, Perspective, and GeoReference.
type LevelRepository interface {
	// Map groups
	FindMapGroupByID(ctx context.Context, id uuid.UUID) (*domain.MapGroup, error)
	ListMapGroups(ctx context.Context, venueID uuid.UUID) ([]*domain.MapGroup, error)
	CreateMapGroup(ctx context.Context, mg *domain.MapGroup) error
	UpdateMapGroup(ctx context.Context, mg *domain.MapGroup) error
	DeleteMapGroup(ctx context.Context, id uuid.UUID) error

	// Levels
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Level, error)
	List(ctx context.Context, venueID uuid.UUID) ([]*domain.Level, error)
	Create(ctx context.Context, l *domain.Level) error
	Update(ctx context.Context, l *domain.Level) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Perspectives
	UpsertPerspective(ctx context.Context, p *domain.Perspective) error

	// Geo references
	ListGeoReferences(ctx context.Context, levelID uuid.UUID) ([]*domain.GeoReference, error)
	CreateGeoReference(ctx context.Context, g *domain.GeoReference) error
	DeleteGeoReference(ctx context.Context, id uuid.UUID) error
}

// LocationCategoryRepository handles location categories.
type LocationCategoryRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.LocationCategory, error)
	FindByNameAndVenue(ctx context.Context, venueID uuid.UUID, name, source string) (*domain.LocationCategory, error)
	List(ctx context.Context, venueID uuid.UUID) ([]*domain.LocationCategory, error)
	Create(ctx context.Context, c *domain.LocationCategory) error
	Update(ctx context.Context, c *domain.LocationCategory) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// LocationRepository handles locations, images, and promotions.
type LocationRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Location, error)
	List(ctx context.Context, venueID uuid.UUID, typeFilter *int, p domain.Pagination) ([]*domain.Location, int64, error)
	SearchByName(ctx context.Context, venueID uuid.UUID, q string, limit int) ([]*domain.Location, error)
	Create(ctx context.Context, l *domain.Location) error
	Update(ctx context.Context, l *domain.Location) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetCategories(ctx context.Context, locationID uuid.UUID, categoryIDs []uuid.UUID) error
	SetTopLocation(ctx context.Context, id uuid.UUID, isTop bool, sortIndex *int) error
	ListTopLocations(ctx context.Context, venueID uuid.UUID) ([]*domain.Location, error)
	GeoSearch(ctx context.Context, lat, lng, radiusKm float64, venueID *uuid.UUID) ([]*domain.Location, error)
	FindByExternalID(ctx context.Context, venueID uuid.UUID, externalID string, locationType int) (*domain.Location, error)

	// Images
	ListImages(ctx context.Context, locationID uuid.UUID) ([]*domain.LocationImage, error)
	CreateImage(ctx context.Context, img *domain.LocationImage) error
	DeleteImage(ctx context.Context, id uuid.UUID) error

	// Memos — locations with common_location_type = LocationTypeMemo (6)
	ListMemos(ctx context.Context, venueID uuid.UUID) ([]*domain.Location, error)
	FindMemoByID(ctx context.Context, id uuid.UUID) (*domain.Location, error)
}

// ProductRepository handles products, categories, and attachments.
type ProductRepository interface {
	// Categories
	FindCategoryByID(ctx context.Context, id uuid.UUID) (*domain.ProductCategory, error)
	ListCategories(ctx context.Context, venueID uuid.UUID) ([]*domain.ProductCategory, error)
	CreateCategory(ctx context.Context, c *domain.ProductCategory) error
	UpdateCategory(ctx context.Context, c *domain.ProductCategory) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error

	// Products
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Product, int64, error)
	SearchByName(ctx context.Context, venueID uuid.UUID, q string, limit int) ([]*domain.Product, error)
	Create(ctx context.Context, p *domain.Product) error
	Update(ctx context.Context, p *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetCategories(ctx context.Context, productID uuid.UUID, categoryIDs []uuid.UUID) error

	FindByCode(ctx context.Context, venueID uuid.UUID, code, source string) (*domain.Product, error)

	// Attachments
	ListAttachments(ctx context.Context, productID uuid.UUID) ([]*domain.ProductAttachment, error)
	CreateAttachment(ctx context.Context, a *domain.ProductAttachment) error
	DeleteAttachment(ctx context.Context, id uuid.UUID) error
}

// NotificationRepository handles push notification records.
type NotificationRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error)
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Notification, int64, error)
	ListDueScheduled(ctx context.Context, now time.Time) ([]*domain.Notification, error)
	Create(ctx context.Context, n *domain.Notification) error
	Update(ctx context.Context, n *domain.Notification) error
	Delete(ctx context.Context, id uuid.UUID) error
	MarkSent(ctx context.Context, id uuid.UUID, publishedAt time.Time) error
	MarkFailed(ctx context.Context, id uuid.UUID, errInfos []byte) error
}

// SurveyRepository handles surveys, questions, options, and responses.
type SurveyRepository interface {
	// Surveys
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Survey, error)
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Survey, int64, error)
	ListDueActivation(ctx context.Context, now time.Time) ([]*domain.Survey, error)
	ListDueClosure(ctx context.Context, now time.Time) ([]*domain.Survey, error)
	Create(ctx context.Context, s *domain.Survey) error
	Update(ctx context.Context, s *domain.Survey) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Questions
	FindQuestionByID(ctx context.Context, id uuid.UUID) (*domain.Question, error)
	ListQuestions(ctx context.Context, surveyID uuid.UUID) ([]*domain.Question, error)
	CreateQuestion(ctx context.Context, q *domain.Question) error
	UpdateQuestion(ctx context.Context, q *domain.Question) error
	DeleteQuestion(ctx context.Context, id uuid.UUID) error

	// Options
	CreateOption(ctx context.Context, o *domain.Option) error
	UpdateOption(ctx context.Context, o *domain.Option) error
	DeleteOption(ctx context.Context, id uuid.UUID) error

	// Active surveys for visitor/app consumption
	ListActive(ctx context.Context, venueID uuid.UUID, publishTypes []int) ([]*domain.Survey, error)

	// Responses
	ListResponses(ctx context.Context, surveyID uuid.UUID, p domain.Pagination) ([]*domain.SurveyResponse, int64, error)
	CreateResponse(ctx context.Context, r *domain.SurveyResponse) error
}

// BeaconRepository handles beacon CRUD.
type BeaconRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Beacon, error)
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Beacon, int64, error)
	Create(ctx context.Context, b *domain.Beacon) error
	Update(ctx context.Context, b *domain.Beacon) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// EventRepository handles events, event types, tags, and images.
type EventRepository interface {
	// Tags (global)
	FindTagByID(ctx context.Context, id uuid.UUID) (*domain.EventTag, error)
	ListTags(ctx context.Context) ([]*domain.EventTag, error)
	CreateTag(ctx context.Context, t *domain.EventTag) error
	UpdateTag(ctx context.Context, t *domain.EventTag) error
	DeleteTag(ctx context.Context, id uuid.UUID) error

	// Event types (per venue)
	FindEventTypeByID(ctx context.Context, id uuid.UUID) (*domain.EventType, error)
	ListEventTypes(ctx context.Context, venueID uuid.UUID) ([]*domain.EventType, error)
	CreateEventType(ctx context.Context, t *domain.EventType) error
	UpdateEventType(ctx context.Context, t *domain.EventType) error
	DeleteEventType(ctx context.Context, id uuid.UUID) error

	// Events
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Event, int64, error)
	Create(ctx context.Context, e *domain.Event) error
	Update(ctx context.Context, e *domain.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetTags(ctx context.Context, eventID uuid.UUID, tagIDs []uuid.UUID) error
	SetLocations(ctx context.Context, eventID uuid.UUID, locationIDs []uuid.UUID) error

	// Images
	CreateImage(ctx context.Context, img *domain.EventImage) error
	DeleteImage(ctx context.Context, id uuid.UUID) error
}

// TokenRepository handles persistence for refresh and reset-password tokens.
type TokenRepository interface {
	// Refresh tokens
	CreateRefreshToken(ctx context.Context, t *domain.RefreshToken) error
	FindRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id uuid.UUID) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error

	// Password reset tokens
	CreateResetToken(ctx context.Context, t *domain.ResetPasswordToken) error
	FindResetToken(ctx context.Context, tokenHash string) (*domain.ResetPasswordToken, error)
	MarkResetTokenUsed(ctx context.Context, id uuid.UUID) error
}

// ConnectionRepository handles connections and their level links.
type ConnectionRepository interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Connection, int, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Connection, error)
	Create(ctx context.Context, c *domain.Connection) error
	Update(ctx context.Context, c *domain.Connection) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListLevels(ctx context.Context, connectionID uuid.UUID) ([]*domain.ConnectionLevel, error)
	AddLevel(ctx context.Context, cl *domain.ConnectionLevel) error
	RemoveLevel(ctx context.Context, id uuid.UUID) error
}

// AdvertisementRepository handles advertisements.
type AdvertisementRepository interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Advertisement, int, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Advertisement, error)
	Create(ctx context.Context, a *domain.Advertisement) error
	Update(ctx context.Context, a *domain.Advertisement) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ArticleRepository handles articles and their images.
type ArticleRepository interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Article, int, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Article, error)
	Create(ctx context.Context, a *domain.Article) error
	Update(ctx context.Context, a *domain.Article) error
	UpdateWithImages(ctx context.Context, a *domain.Article, change domain.ArticleMediaChange) error
	Delete(ctx context.Context, id uuid.UUID) error
	CreateImage(ctx context.Context, img *domain.ArticleImage) error
	DeleteImage(ctx context.Context, id uuid.UUID) error
}

// CouponRepository handles coupons.
type CouponRepository interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Coupon, int, error)
	// ListByUser returns coupons assigned to a specific app user for a venue,
	// with IsUsed populated from coupon_users.
	ListByUser(ctx context.Context, venueID, userID uuid.UUID, p domain.Pagination) ([]*domain.Coupon, int, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Coupon, error)
	Create(ctx context.Context, c *domain.Coupon) error
	Update(ctx context.Context, c *domain.Coupon) error
	Delete(ctx context.Context, id uuid.UUID) error
	Redeem(ctx context.Context, id uuid.UUID, appUserID uuid.UUID) error
}

// CouponUserRepository handles per-user coupon assignment and redemption.
type CouponUserRepository interface {
	Create(ctx context.Context, cu *domain.CouponUser) error
	FindByCouponAndUser(ctx context.Context, couponID, userID uuid.UUID) (*domain.CouponUser, error)
	MarkUsed(ctx context.Context, id uuid.UUID) error
	// ExistsForVenueUser returns true if the user already has any coupon for the venue.
	ExistsForVenueUser(ctx context.Context, venueID, userID uuid.UUID) (bool, error)
}

// VideoRepository handles videos.
type VideoRepository interface {
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Video, int, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Video, error)
	Create(ctx context.Context, v *domain.Video) error
	Update(ctx context.Context, v *domain.Video) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// TagRepository handles global tags and entity-tag links.
type TagRepository interface {
	List(ctx context.Context, p domain.Pagination) ([]*domain.Tag, int, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Tag, error)
	Create(ctx context.Context, t *domain.Tag) error
	Update(ctx context.Context, t *domain.Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
	AttachTag(ctx context.Context, et *domain.EntityTag) error
	DetachTag(ctx context.Context, tagID uuid.UUID, entityType string, entityID uuid.UUID) error
	ListEntityTags(ctx context.Context, entityType string, entityID uuid.UUID) ([]*domain.Tag, error)
}

// EventLogRepository handles high-volume analytics event writes.
type EventLogRepository interface {
	Create(ctx context.Context, e *domain.EventLog) error
	ListByVenue(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.EventLog, int, error)
}

// SearchQueryFilter narrows search_queries results on the read path.
type SearchQueryFilter struct {
	Origin     *string
	IsPromoted *bool
}

// SearchQueryRepository handles search term tracking and promoted keywords.
type SearchQueryRepository interface {
	Upsert(ctx context.Context, venueID uuid.UUID, term, origin, appID string) error
	List(ctx context.Context, venueID uuid.UUID, filter SearchQueryFilter, p domain.Pagination) ([]*domain.SearchQuery, int, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.SearchQuery, error)
	Create(ctx context.Context, q *domain.SearchQuery) error
	Update(ctx context.Context, q *domain.SearchQuery) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// SnapshotRepository handles snapshot metadata persistence.
type SnapshotRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Snapshot, error)
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Snapshot, int64, error)
	LatestPublished(ctx context.Context, venueID uuid.UUID) (*domain.Snapshot, error)
	Create(ctx context.Context, s *domain.Snapshot) error
	UpdateState(ctx context.Context, id uuid.UUID, state int, publishAt *time.Time) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountDraftsByVenue(ctx context.Context, venueID uuid.UUID) (int64, error)
	// DeleteOldestDraft hard-deletes (not soft-delete) the draft snapshot with
	// the oldest created_at for the given venue, triggering ON DELETE CASCADE
	// to remove its associated LevelBundle rows. Returns the deleted snapshot
	// so callers can remove the corresponding S3 object. Returns (nil, nil) if
	// no draft exists.
	DeleteOldestDraft(ctx context.Context, venueID uuid.UUID) (*domain.Snapshot, error)
	UnpublishVenue(ctx context.Context, venueID uuid.UUID) error
	// LatestDraft returns the most recently created draft snapshot for a venue (nil if none).
	LatestDraft(ctx context.Context, venueID uuid.UUID) (*domain.Snapshot, error)
}

// LevelBundleRepository handles per-level bundle data persistence.
type LevelBundleRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.LevelBundle, error)
	ListBySnapshot(ctx context.Context, snapshotID uuid.UUID) ([]*domain.LevelBundle, error)
	Create(ctx context.Context, b *domain.LevelBundle) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// AssetRepository handles asset metadata persistence.
type AssetRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Asset, error)
	List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Asset, int64, error)
	Create(ctx context.Context, a *domain.Asset) error
	Update(ctx context.Context, a *domain.Asset) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// LevelTypeRepository handles level type CRUD.
type LevelTypeRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.LevelType, error)
	List(ctx context.Context) ([]*domain.LevelType, error)
	Create(ctx context.Context, lt *domain.LevelType) error
	Update(ctx context.Context, lt *domain.LevelType) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ThemeRepository handles theme CRUD.
type ThemeRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Theme, error)
	// FindVenueTheme returns the theme currently selected by a venue (via venues.theme_id).
	// Returns nil, nil when the venue has no theme set.
	FindVenueTheme(ctx context.Context, venueID uuid.UUID) (*domain.Theme, error)
	ListGlobal(ctx context.Context) ([]*domain.Theme, error)
	List(ctx context.Context, venueID uuid.UUID) ([]*domain.Theme, error)
	Create(ctx context.Context, t *domain.Theme) error
	Update(ctx context.Context, t *domain.Theme) error
	Delete(ctx context.Context, id uuid.UUID) error
	IsUsedByVenues(ctx context.Context, id uuid.UUID) (bool, error)
}

// ProductPlazaRepository handles product plaza CRUD.
type ProductPlazaRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.ProductPlaza, error)
	List(ctx context.Context, venueID uuid.UUID) ([]*domain.ProductPlaza, error)
	Create(ctx context.Context, p *domain.ProductPlaza) error
	Update(ctx context.Context, p *domain.ProductPlaza) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// AppUserRepository handles visitor/app user persistence.
type AppUserRepository interface {
	FindByToken(ctx context.Context, venueID uuid.UUID, token string) (*domain.AppUser, error)
	Create(ctx context.Context, u *domain.AppUser) error
}

// LanguageRepository handles supported-language CRUD per venue.
type LanguageRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Language, error)
	List(ctx context.Context, venueID uuid.UUID) ([]*domain.Language, error)
	// ListEnabled returns only active languages for a venue (used by the v2 bundle publisher).
	ListEnabled(ctx context.Context, venueID uuid.UUID) ([]*domain.Language, error)
	Create(ctx context.Context, l *domain.Language) error
	Update(ctx context.Context, l *domain.Language) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// AppVersionRepository persists the force-sync version for each venue.
// Mirrors Django's ForceSyncVersion model (indoormap-backend/indoormap_api/app/models.py:57).
type AppVersionRepository interface {
	// Upsert creates or updates the version row for the given venue.
	Upsert(ctx context.Context, venueID, version uuid.UUID) error
	// Get returns the current version for a venue, creating a default row if none exists.
	Get(ctx context.Context, venueID uuid.UUID) (*domain.AppVersion, error)
}
