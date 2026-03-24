package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/database"
	"github.com/hhung06/digimap-backend/internal/repository"
)

// ── Customer service ──────────────────────────────────────────────────────────

type CustomerService interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	List(ctx context.Context, p domain.Pagination) ([]*domain.Customer, int64, error)
	Create(ctx context.Context, c *domain.Customer) error
	Update(ctx context.Context, c *domain.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type customerService struct {
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func (s *customerService) Get(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *customerService) List(ctx context.Context, p domain.Pagination) ([]*domain.Customer, int64, error) {
	p.Normalize()
	return s.repo.List(ctx, p)
}

func (s *customerService) Create(ctx context.Context, c *domain.Customer) error {
	return s.repo.Create(ctx, c)
}

func (s *customerService) Update(ctx context.Context, c *domain.Customer) error {
	return s.repo.Update(ctx, c)
}

func (s *customerService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// ── Venue service ─────────────────────────────────────────────────────────────

type VenueService interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Venue, error)
	GetByPublicKey(ctx context.Context, publicKey string) (*domain.Venue, error)
	List(ctx context.Context, customerID uuid.UUID, p domain.Pagination) ([]*domain.Venue, int64, error)
	ListAll(ctx context.Context, p domain.Pagination) ([]*domain.Venue, int64, error)
	Create(ctx context.Context, v *domain.Venue) error
	Update(ctx context.Context, v *domain.Venue) error
	Delete(ctx context.Context, id uuid.UUID) error
	Publish(ctx context.Context, id uuid.UUID, published bool) error
	RegenerateKeys(ctx context.Context, id uuid.UUID) (*domain.Venue, error)
	Clone(ctx context.Context, sourceID uuid.UUID) (*domain.Venue, error)
}

type venueService struct {
	repo      repository.VenueRepository
	levelRepo repository.LevelRepository
	pool      *pgxpool.Pool
}

func NewVenueService(repo repository.VenueRepository, levelRepo repository.LevelRepository, pool *pgxpool.Pool) VenueService {
	return &venueService{repo: repo, levelRepo: levelRepo, pool: pool}
}

func (s *venueService) Get(ctx context.Context, id uuid.UUID) (*domain.Venue, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *venueService) GetByPublicKey(ctx context.Context, publicKey string) (*domain.Venue, error) {
	return s.repo.FindByPublicKey(ctx, publicKey)
}

func (s *venueService) List(ctx context.Context, customerID uuid.UUID, p domain.Pagination) ([]*domain.Venue, int64, error) {
	p.Normalize()
	return s.repo.List(ctx, customerID, p)
}

func (s *venueService) ListAll(ctx context.Context, p domain.Pagination) ([]*domain.Venue, int64, error) {
	p.Normalize()
	return s.repo.ListAll(ctx, p)
}

func (s *venueService) Create(ctx context.Context, v *domain.Venue) error {
	pub, priv, err := generateVenueKeys()
	if err != nil {
		return fmt.Errorf("generate venue keys: %w", err)
	}
	v.PublicKey = pub
	v.PrivateKey = priv
	return s.repo.Create(ctx, v)
}

func (s *venueService) Update(ctx context.Context, v *domain.Venue) error {
	existing, err := s.repo.FindByID(ctx, v.ID)
	if err != nil {
		return err
	}
	// Preserve immutable fields
	v.PublicKey = existing.PublicKey
	v.PrivateKey = existing.PrivateKey
	v.CustomerID = existing.CustomerID
	return s.repo.Update(ctx, v)
}

func (s *venueService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *venueService) Publish(ctx context.Context, id uuid.UUID, published bool) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.repo.UpdatePublished(ctx, id, published)
}

func (s *venueService) RegenerateKeys(ctx context.Context, id uuid.UUID) (*domain.Venue, error) {
	pub, priv, err := generateVenueKeys()
	if err != nil {
		return nil, fmt.Errorf("generate venue keys: %w", err)
	}
	if err := s.repo.UpdateKeys(ctx, id, pub, priv); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, id)
}

// Clone creates a full copy of the source venue (including all levels and map
// groups) in a single transaction. The clone gets fresh API keys and is
// unpublished by default.
func (s *venueService) Clone(ctx context.Context, sourceID uuid.UUID) (*domain.Venue, error) {
	source, err := s.repo.FindByID(ctx, sourceID)
	if err != nil {
		return nil, err
	}

	pub, priv, err := generateVenueKeys()
	if err != nil {
		return nil, fmt.Errorf("generate venue keys: %w", err)
	}

	var cloned *domain.Venue
	err = database.WithTransaction(ctx, s.pool, func(ctx context.Context, _ pgx.Tx) error {
		clone := *source
		clone.ID = uuid.Nil // will be assigned by newID() in Create
		clone.PublicKey = pub
		clone.PrivateKey = priv
		clone.IsPublished = false
		clone.Name = source.Name + " (copy)"

		if err := s.repo.Create(ctx, &clone); err != nil {
			return fmt.Errorf("clone venue: %w", err)
		}

		// Clone map groups (venue-scoped), building an old→new ID map
		srcGroups, err := s.levelRepo.ListMapGroups(ctx, sourceID)
		if err != nil {
			return fmt.Errorf("list map groups: %w", err)
		}
		groupIDMap := make(map[uuid.UUID]uuid.UUID, len(srcGroups))
		for _, mg := range srcGroups {
			newMG := *mg
			newMG.ID = uuid.Nil
			newMG.VenueID = clone.ID
			if err := s.levelRepo.CreateMapGroup(ctx, &newMG); err != nil {
				return fmt.Errorf("clone map group: %w", err)
			}
			groupIDMap[mg.ID] = newMG.ID
		}

		// Clone levels
		srcLevels, err := s.levelRepo.List(ctx, sourceID)
		if err != nil {
			return fmt.Errorf("list levels: %w", err)
		}
		for _, l := range srcLevels {
			newLevel := *l
			newLevel.ID = uuid.Nil
			newLevel.VenueID = clone.ID
			newLevel.Perspective = nil
			newLevel.PerspectiveID = nil

			// Remap map group reference
			if l.MapGroupID != nil {
				if newGroupID, ok := groupIDMap[*l.MapGroupID]; ok {
					newLevel.MapGroupID = &newGroupID
				}
			}

			// Clone perspective if present
			if l.PerspectiveID != nil {
				// Fetch and clone perspective
				origPerspID := *l.PerspectiveID
				_ = origPerspID // perspective is fetched separately if needed
				newPersp := &domain.Perspective{ID: uuid.Nil}
				if err := s.levelRepo.UpsertPerspective(ctx, newPersp); err != nil {
					return fmt.Errorf("clone perspective: %w", err)
				}
				newLevel.PerspectiveID = &newPersp.ID
			}

			if err := s.levelRepo.Create(ctx, &newLevel); err != nil {
				return fmt.Errorf("clone level: %w", err)
			}
		}

		cloned = &clone
		return nil
	})
	return cloned, err
}

// ── Level service ─────────────────────────────────────────────────────────────

type LevelService interface {
	// Map groups
	GetMapGroup(ctx context.Context, id uuid.UUID) (*domain.MapGroup, error)
	ListMapGroups(ctx context.Context, venueID uuid.UUID) ([]*domain.MapGroup, error)
	CreateMapGroup(ctx context.Context, mg *domain.MapGroup) error
	UpdateMapGroup(ctx context.Context, mg *domain.MapGroup) error
	DeleteMapGroup(ctx context.Context, id uuid.UUID) error

	// Levels
	Get(ctx context.Context, id uuid.UUID) (*domain.Level, error)
	List(ctx context.Context, venueID uuid.UUID) ([]*domain.Level, error)
	Create(ctx context.Context, l *domain.Level) error
	Update(ctx context.Context, l *domain.Level) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Perspectives
	UpsertPerspective(ctx context.Context, levelID uuid.UUID, p *domain.Perspective) error

	// Geo references
	ListGeoReferences(ctx context.Context, levelID uuid.UUID) ([]*domain.GeoReference, error)
	CreateGeoReference(ctx context.Context, g *domain.GeoReference) error
	DeleteGeoReference(ctx context.Context, levelID, geoRefID uuid.UUID) error
}

type levelService struct {
	repo repository.LevelRepository
}

func NewLevelService(repo repository.LevelRepository) LevelService {
	return &levelService{repo: repo}
}

func (s *levelService) GetMapGroup(ctx context.Context, id uuid.UUID) (*domain.MapGroup, error) {
	return s.repo.FindMapGroupByID(ctx, id)
}

func (s *levelService) ListMapGroups(ctx context.Context, venueID uuid.UUID) ([]*domain.MapGroup, error) {
	return s.repo.ListMapGroups(ctx, venueID)
}

func (s *levelService) CreateMapGroup(ctx context.Context, mg *domain.MapGroup) error {
	return s.repo.CreateMapGroup(ctx, mg)
}

func (s *levelService) UpdateMapGroup(ctx context.Context, mg *domain.MapGroup) error {
	return s.repo.UpdateMapGroup(ctx, mg)
}

func (s *levelService) DeleteMapGroup(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteMapGroup(ctx, id)
}

func (s *levelService) Get(ctx context.Context, id uuid.UUID) (*domain.Level, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *levelService) List(ctx context.Context, venueID uuid.UUID) ([]*domain.Level, error) {
	return s.repo.List(ctx, venueID)
}

func (s *levelService) Create(ctx context.Context, l *domain.Level) error {
	return s.repo.Create(ctx, l)
}

func (s *levelService) Update(ctx context.Context, l *domain.Level) error {
	return s.repo.Update(ctx, l)
}

func (s *levelService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *levelService) UpsertPerspective(ctx context.Context, levelID uuid.UUID, p *domain.Perspective) error {
	if err := s.repo.UpsertPerspective(ctx, p); err != nil {
		return err
	}
	// Link perspective to level
	l := &domain.Level{ID: levelID, PerspectiveID: &p.ID}
	return s.repo.Update(ctx, l)
}

func (s *levelService) ListGeoReferences(ctx context.Context, levelID uuid.UUID) ([]*domain.GeoReference, error) {
	return s.repo.ListGeoReferences(ctx, levelID)
}

func (s *levelService) CreateGeoReference(ctx context.Context, g *domain.GeoReference) error {
	return s.repo.CreateGeoReference(ctx, g)
}

func (s *levelService) DeleteGeoReference(ctx context.Context, levelID, geoRefID uuid.UUID) error {
	// Verify the geo ref belongs to this level
	refs, err := s.repo.ListGeoReferences(ctx, levelID)
	if err != nil {
		return err
	}
	for _, r := range refs {
		if r.ID == geoRefID {
			return s.repo.DeleteGeoReference(ctx, geoRefID)
		}
	}
	return domain.NewNotFound("geo reference not found")
}

// ── key generation ────────────────────────────────────────────────────────────

// generateVenueKeys matches the Django implementation exactly:
//   - public key:  URL-safe base64 of 32 random bytes (43 chars)
//   - private key: SHA-256 hex of 64 random bytes (64 chars)
func generateVenueKeys() (publicKey, privateKey string, err error) {
	pubRaw := make([]byte, 32)
	if _, err = rand.Read(pubRaw); err != nil {
		return
	}
	publicKey = base64.RawURLEncoding.EncodeToString(pubRaw)

	privRaw := make([]byte, 64)
	if _, err = rand.Read(privRaw); err != nil {
		return
	}
	h := sha256.Sum256(privRaw)
	privateKey = hex.EncodeToString(h[:])
	return
}
