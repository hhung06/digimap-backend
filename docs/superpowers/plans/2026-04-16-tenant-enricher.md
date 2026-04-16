# Tenant Enricher Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a compile-time per-tenant enricher registry that appends extra computed fields to Location, Product, Notification, Survey, and Advertisement responses without changing base field names.

**Architecture:** A `Registry` in `internal/enricher` maps `(customerID, resource)` → `EnricherFunc`. The registry resolves the tenant by querying `venues.customer_id` (fast path: skips DB entirely when no enrichers are registered or the registry is nil). Handlers call `EnrichForVenue` after building the base DTO and merge extras using `MergeInto`. Tenant registrations live in `cmd/tenants/`. Task order ensures a clean build after every commit: enricher package → GetCustomerID → router/serve wiring → individual handler updates.

**Tech Stack:** Go 1.21+, Gin, pgx/v5, testify/assert

---

## File Map

| Action | File | Purpose |
|---|---|---|
| Create | `internal/enricher/enricher.go` | `Resource`, `EnricherFunc`, `VenueCustomerResolver`, `Registry` |
| Create | `internal/enricher/merge.go` | `MergeInto` helper |
| Create | `internal/enricher/enricher_test.go` | Unit tests for registry and merge |
| Modify | `internal/repository/interfaces.go` | Add `GetCustomerID` to `VenueRepository` |
| Modify | `internal/repository/postgres/venue_repo.go` | Implement `GetCustomerID` |
| Modify | `internal/handler/router.go` | Add `EnricherRegistry` to `Dependencies`, pass to constructors |
| Modify | `cmd/serve.go` | Create registry, wire into deps |
| Create | `cmd/tenants/tenants.go` | Placeholder for tenant registrations |
| Modify | `internal/handler/location_handler.go` | Add `enrichers`, enrich `ListLocations` + `GetLocation` |
| Modify | `internal/handler/product_handler.go` | Add `enrichers`, enrich `List` + `Get` |
| Modify | `internal/handler/notification_handler.go` | Add `enrichers`, enrich `List` + `Get` |
| Modify | `internal/handler/survey_handler.go` | Add `enrichers`, enrich `List` + `Get` |
| Modify | `internal/handler/ad_handler.go` | Add `enrichers`, enrich `List` + `Get` |

---

### Task 1: Enricher package — core types, registry, merge

**Files:**
- Create: `internal/enricher/enricher.go`
- Create: `internal/enricher/merge.go`
- Create: `internal/enricher/enricher_test.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/enricher/enricher_test.go`:

```go
package enricher_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/internal/enricher"
)

// stubVenueRepo satisfies enricher.VenueCustomerResolver for tests.
type stubVenueRepo struct {
	customerID uuid.UUID
	err        error
}

func (s *stubVenueRepo) GetCustomerID(_ context.Context, _ uuid.UUID) (uuid.UUID, error) {
	return s.customerID, s.err
}

func TestRegistry_EnrichForVenue_NilRegistry(t *testing.T) {
	var reg *enricher.Registry
	extras, err := reg.EnrichForVenue(context.Background(), uuid.New(), enricher.ResourceLocation)
	require.NoError(t, err)
	assert.Nil(t, extras)
}

func TestRegistry_EnrichForVenue_EmptyRegistry(t *testing.T) {
	called := false
	repo := &trackingRepo{called: &called}
	reg := enricher.NewRegistry(repo)

	extras, err := reg.EnrichForVenue(context.Background(), uuid.New(), enricher.ResourceLocation)
	require.NoError(t, err)
	assert.Nil(t, extras)
	assert.False(t, called, "GetCustomerID must not be called on empty registry")
}

func TestRegistry_EnrichForVenue_NoMatchingEnricher(t *testing.T) {
	customerID := uuid.New()
	venueID := uuid.New()
	otherID := uuid.New()

	reg := enricher.NewRegistry(&stubVenueRepo{customerID: customerID})
	reg.Register(otherID, enricher.ResourceLocation, func(_ context.Context, _ uuid.UUID) (map[string]any, error) {
		return map[string]any{"key": "val"}, nil
	})

	extras, err := reg.EnrichForVenue(context.Background(), venueID, enricher.ResourceLocation)
	require.NoError(t, err)
	assert.Nil(t, extras)
}

func TestRegistry_EnrichForVenue_ReturnsExtras(t *testing.T) {
	customerID := uuid.New()
	venueID := uuid.New()

	reg := enricher.NewRegistry(&stubVenueRepo{customerID: customerID})
	reg.Register(customerID, enricher.ResourceLocation, func(_ context.Context, _ uuid.UUID) (map[string]any, error) {
		return map[string]any{"crm_zone": "west"}, nil
	})

	extras, err := reg.EnrichForVenue(context.Background(), venueID, enricher.ResourceLocation)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"crm_zone": "west"}, extras)
}

func TestRegistry_EnrichForVenue_RepoError(t *testing.T) {
	venueID := uuid.New()

	reg := enricher.NewRegistry(&stubVenueRepo{err: errors.New("db error")})
	reg.Register(uuid.New(), enricher.ResourceLocation, func(_ context.Context, _ uuid.UUID) (map[string]any, error) {
		return map[string]any{"key": "val"}, nil
	})

	extras, err := reg.EnrichForVenue(context.Background(), venueID, enricher.ResourceLocation)
	require.NoError(t, err) // repo errors are swallowed — base response unaffected
	assert.Nil(t, extras)
}

func TestMergeInto_NoExtras(t *testing.T) {
	type base struct {
		Name string `json:"name"`
	}
	result := enricher.MergeInto(base{Name: "test"}, nil)
	assert.Equal(t, base{Name: "test"}, result)
}

func TestMergeInto_EmptyExtras(t *testing.T) {
	type base struct {
		Name string `json:"name"`
	}
	result := enricher.MergeInto(base{Name: "test"}, map[string]any{})
	assert.Equal(t, base{Name: "test"}, result)
}

func TestMergeInto_WithExtras(t *testing.T) {
	type base struct {
		Name string `json:"name"`
	}
	result := enricher.MergeInto(base{Name: "test"}, map[string]any{"crm_zone": "west", "custom_label": "Gate"})

	m, ok := result.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "test", m["name"])
	assert.Equal(t, "west", m["crm_zone"])
	assert.Equal(t, "Gate", m["custom_label"])
}

// trackingRepo records whether GetCustomerID was called.
type trackingRepo struct {
	called *bool
}

func (t *trackingRepo) GetCustomerID(_ context.Context, _ uuid.UUID) (uuid.UUID, error) {
	*t.called = true
	return uuid.Nil, nil
}
```

- [ ] **Step 2: Run tests — verify they fail**

```bash
go test ./internal/enricher/... -v
```
Expected: `cannot find package` or `no Go files`

- [ ] **Step 3: Create `internal/enricher/enricher.go`**

```go
package enricher

import (
	"context"

	"github.com/google/uuid"
)

// Resource identifies which domain entity is being enriched.
type Resource string

const (
	ResourceLocation     Resource = "location"
	ResourceProduct      Resource = "product"
	ResourceNotification Resource = "notification"
	ResourceSurvey       Resource = "survey"
	ResourceAd           Resource = "advertisement"
)

// VenueCustomerResolver resolves a venue's owning customer.
// Registry depends only on this narrow interface, not the full VenueRepository.
type VenueCustomerResolver interface {
	GetCustomerID(ctx context.Context, venueID uuid.UUID) (uuid.UUID, error)
}

// EnricherFunc computes extra fields for a tenant+venue combination.
// The same extras map is applied to every item in a list response (one call per request).
// Return nil, nil when no extras are needed.
type EnricherFunc func(ctx context.Context, venueID uuid.UUID) (map[string]any, error)

type enricherKey struct {
	customerID uuid.UUID
	resource   Resource
}

// Registry maps (customerID, resource) → EnricherFunc.
// Register all enrichers at startup before the server begins serving requests.
type Registry struct {
	m         map[enricherKey]EnricherFunc
	venueRepo VenueCustomerResolver
}

// NewRegistry creates an empty Registry backed by the given VenueCustomerResolver.
func NewRegistry(venueRepo VenueCustomerResolver) *Registry {
	return &Registry{
		m:         make(map[enricherKey]EnricherFunc),
		venueRepo: venueRepo,
	}
}

// Register associates fn with the given customer and resource.
func (r *Registry) Register(customerID uuid.UUID, resource Resource, fn EnricherFunc) {
	r.m[enricherKey{customerID, resource}] = fn
}

// EnrichForVenue resolves the tenant from venueID and calls the registered enricher.
// Safe to call on a nil *Registry — returns nil, nil immediately.
// Also returns nil, nil when no enrichers are registered (zero DB cost) or no enricher
// matches the resolved customer+resource pair.
func (r *Registry) EnrichForVenue(ctx context.Context, venueID uuid.UUID, resource Resource) (map[string]any, error) {
	if r == nil || len(r.m) == 0 {
		return nil, nil
	}
	customerID, err := r.venueRepo.GetCustomerID(ctx, venueID)
	if err != nil {
		return nil, nil
	}
	fn, ok := r.m[enricherKey{customerID, resource}]
	if !ok {
		return nil, nil
	}
	return fn(ctx, venueID)
}
```

- [ ] **Step 4: Create `internal/enricher/merge.go`**

```go
package enricher

import "encoding/json"

// MergeInto merges extras on top of the base DTO and returns the result as any.
// Returns base unchanged (no JSON round-trip cost) if extras is nil or empty.
func MergeInto(base any, extras map[string]any) any {
	if len(extras) == 0 {
		return base
	}
	b, err := json.Marshal(base)
	if err != nil {
		return base
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return base
	}
	for k, v := range extras {
		m[k] = v
	}
	return m
}
```

- [ ] **Step 5: Run tests — verify all pass**

```bash
go test ./internal/enricher/... -v
```
Expected: all 8 tests PASS

- [ ] **Step 6: Commit**

```bash
git add internal/enricher/
git commit -m "feat: add tenant enricher registry and MergeInto helper"
```

---

### Task 2: Add `GetCustomerID` to VenueRepository

**Files:**
- Modify: `internal/repository/interfaces.go`
- Modify: `internal/repository/postgres/venue_repo.go`

- [ ] **Step 1: Add `GetCustomerID` to the `VenueRepository` interface**

In `internal/repository/interfaces.go`, add one line to `VenueRepository`:

```go
type VenueRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Venue, error)
	FindByPublicKey(ctx context.Context, publicKey string) (*domain.Venue, error)
	List(ctx context.Context, customerID uuid.UUID, p domain.Pagination) ([]*domain.Venue, int64, error)
	ListAll(ctx context.Context, p domain.Pagination) ([]*domain.Venue, int64, error)
	Create(ctx context.Context, v *domain.Venue) error
	Update(ctx context.Context, v *domain.Venue) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateKeys(ctx context.Context, id uuid.UUID, publicKey, privateKey string) error
	UpdatePublished(ctx context.Context, id uuid.UUID, published bool) error
	GetCustomerID(ctx context.Context, venueID uuid.UUID) (uuid.UUID, error) // new
}
```

- [ ] **Step 2: Verify build fails with missing method**

```bash
make build 2>&1 | grep "does not implement"
```
Expected: `venueRepo does not implement repository.VenueRepository (missing GetCustomerID method)`

- [ ] **Step 3: Implement `GetCustomerID` in the postgres repo**

Add to the bottom of `internal/repository/postgres/venue_repo.go`:

```go
func (r *venueRepo) GetCustomerID(ctx context.Context, venueID uuid.UUID) (uuid.UUID, error) {
	var customerID uuid.UUID
	err := r.pool.QueryRow(ctx,
		`SELECT customer_id FROM venues WHERE id = $1 AND deleted_at IS NULL`,
		venueID,
	).Scan(&customerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, domain.NewNotFound("venue not found")
	}
	return customerID, err
}
```

- [ ] **Step 4: Build and test**

```bash
make build && make test
```
Expected: clean build, all tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/repository/interfaces.go internal/repository/postgres/venue_repo.go
git commit -m "feat: add GetCustomerID to VenueRepository"
```

---

### Task 3: Wire registry into router and serve.go

Do this before the handler tasks so every subsequent commit produces a clean build. The `EnricherRegistry` field starts populated with an empty registry (no tenant enrichers registered yet), so all enricher calls return nil extras — existing behaviour is unchanged.

**Files:**
- Modify: `internal/handler/router.go`
- Modify: `cmd/serve.go`
- Create: `cmd/tenants/tenants.go`

- [ ] **Step 1: Add `EnricherRegistry` to `Dependencies` in `router.go`**

In `internal/handler/router.go`, add the field to the `Dependencies` struct:

```go
type Dependencies struct {
	AuthService             service.AuthService
	CustomerService         service.CustomerService
	VenueService            service.VenueService
	LevelService            service.LevelService
	LocationCategoryService service.LocationCategoryService
	LocationService         service.LocationService
	ProductService          service.ProductService
	StorageService          service.StorageService
	EventService            service.EventService
	UserService             service.UserService
	NotificationService     service.NotificationService
	SurveyService           service.SurveyService
	BeaconService           service.BeaconService
	ConnectionService       service.ConnectionService
	AdvertisementService    service.AdvertisementService
	ArticleService          service.ArticleService
	CouponService           service.CouponService
	VideoService            service.VideoService
	TagService              service.TagService
	AnalyticsService        service.AnalyticsService
	SnapshotService         service.SnapshotService
	LevelBundleService      service.LevelBundleService
	EnricherRegistry        *enricher.Registry
	UserRepo                repository.UserRepository
	RedisClient             *redis.Client
	DB                      *pgxpool.Pool
}
```

Add the import:
```go
"github.com/hhung06/digimap-backend/internal/enricher"
```

- [ ] **Step 2: Pass `deps.EnricherRegistry` to the five affected handler constructors in `NewRouter`**

Find each constructor call and add the registry argument:

```go
locH := newLocationHandler(deps.LocationCategoryService, deps.LocationService, deps.EnricherRegistry)
```

```go
prodH := newProductHandler(deps.ProductService, deps.StorageService, deps.EnricherRegistry)
```

```go
notifH := newNotificationHandler(deps.NotificationService, deps.EnricherRegistry)
```

```go
surveyH := newSurveyHandler(deps.SurveyService, deps.EnricherRegistry)
```

```go
adH := newAdHandler(deps.AdvertisementService, deps.EnricherRegistry)
```

- [ ] **Step 3: Create `cmd/tenants/tenants.go`**

```go
// Package tenants registers per-tenant response enrichers.
// Add one exported Register function per tenant and call it from RegisterAll.
package tenants

import "github.com/hhung06/digimap-backend/internal/enricher"

// RegisterAll calls all active tenant enricher registration functions.
// Add one line per tenant:
//
//	RegisterAcmeCorp(r)
func RegisterAll(_ *enricher.Registry) {}
```

- [ ] **Step 4: Update `cmd/serve.go` to create the registry**

After the last repository line (`levelBundleRepo`) and before the platform services section, add:

```go
// ── Enricher registry ─────────────────────────────────────────────────────
enricherRegistry := enricher.NewRegistry(venueRepo)
tenants.RegisterAll(enricherRegistry)
```

Add `EnricherRegistry` to the `deps` struct:

```go
deps := handler.Dependencies{
	// ... all existing fields ...
	EnricherRegistry:        enricherRegistry,
	// ...
}
```

Add imports to `cmd/serve.go`:

```go
"github.com/hhung06/digimap-backend/internal/enricher"
"github.com/hhung06/digimap-backend/cmd/tenants"
```

- [ ] **Step 5: Build — must be clean (handler constructors not yet updated, but `EnrichForVenue` is nil-safe)**

```bash
make build
```
Expected: **fails** — router.go now passes 3/4 args to constructors that still only accept 2/3. This is expected — the handler updates come next. If you want a fully clean build at this step, skip Step 2 here and do it together with each handler task instead.

> **Note:** The cleanest path is to do Steps 1, 3, 4 now and do Step 2 (the constructor call updates) together with each handler task (Tasks 4–8). Either approach is fine — just be aware.

- [ ] **Step 6: Commit**

```bash
git add internal/handler/router.go cmd/serve.go cmd/tenants/
git commit -m "feat: wire enricher registry into router dependencies and serve"
```

---

### Task 4: Enrich `location_handler` — List and Get

**Files:**
- Modify: `internal/handler/location_handler.go`
- Modify: `internal/handler/router.go` (if Step 2 of Task 3 was deferred)

- [ ] **Step 1: Add `enrichers` field and update the constructor**

In `internal/handler/location_handler.go`, replace the struct and constructor:

```go
type locationHandler struct {
	categorySvc service.LocationCategoryService
	locationSvc service.LocationService
	enrichers   *enricher.Registry
}

func newLocationHandler(
	categorySvc service.LocationCategoryService,
	locationSvc service.LocationService,
	enrichers *enricher.Registry,
) *locationHandler {
	return &locationHandler{categorySvc: categorySvc, locationSvc: locationSvc, enrichers: enrichers}
}
```

Add import:
```go
"github.com/hhung06/digimap-backend/internal/enricher"
```

- [ ] **Step 2: Update `ListLocations`**

Replace the current `ListLocations` method:

```go
func (h *locationHandler) ListLocations(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, err.Error()))
		return
	}
	p := paginationFromQuery(c)
	locations, total, err := h.locationSvc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceLocation)
	items := make([]any, len(locations))
	for i, l := range locations {
		items[i] = enricher.MergeInto(dto.LocationToResponse(l), extras)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}
```

- [ ] **Step 3: Update `GetLocation`**

Replace the current `GetLocation` method:

```go
func (h *locationHandler) GetLocation(c *gin.Context) {
	venueID, _ := parseVenueID(c)
	id, err := uuid.Parse(c.Param("locationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid location id"))
		return
	}
	l, err := h.locationSvc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceLocation)
	c.JSON(http.StatusOK, dto.OK(enricher.MergeInto(dto.LocationToResponse(l), extras)))
}
```

- [ ] **Step 4: If Task 3 Step 2 was deferred, update the router call now**

In `internal/handler/router.go`, find and update:
```go
locH := newLocationHandler(deps.LocationCategoryService, deps.LocationService, deps.EnricherRegistry)
```

- [ ] **Step 5: Build and test**

```bash
make build && make test
```
Expected: clean build, all tests PASS

- [ ] **Step 6: Commit**

```bash
git add internal/handler/location_handler.go internal/handler/router.go
git commit -m "feat: add enricher support to location handler"
```

---

### Task 5: Enrich `product_handler` — List and Get

**Files:**
- Modify: `internal/handler/product_handler.go`
- Modify: `internal/handler/router.go` (if Task 3 Step 2 was deferred)

- [ ] **Step 1: Add `enrichers` field and update the constructor**

```go
type productHandler struct {
	svc        service.ProductService
	storageSvc service.StorageService
	enrichers  *enricher.Registry
}

func newProductHandler(svc service.ProductService, storageSvc service.StorageService, enrichers *enricher.Registry) *productHandler {
	return &productHandler{svc: svc, storageSvc: storageSvc, enrichers: enrichers}
}
```

Add import:
```go
"github.com/hhung06/digimap-backend/internal/enricher"
```

- [ ] **Step 2: Update `List`**

Replace the current `List` method:

```go
func (h *productHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	products, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceProduct)
	items := make([]any, len(products))
	for i, prod := range products {
		items[i] = enricher.MergeInto(dto.ProductToResponse(prod), extras)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}
```

- [ ] **Step 3: Update `Get`**

Replace the current `Get` method:

```go
func (h *productHandler) Get(c *gin.Context) {
	venueID, _ := parseVenueID(c)
	id, err := uuid.Parse(c.Param("productID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid product id"))
		return
	}
	prod, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceProduct)
	c.JSON(http.StatusOK, dto.OK(enricher.MergeInto(dto.ProductToResponse(prod), extras)))
}
```

- [ ] **Step 4: If Task 3 Step 2 was deferred, update the router call now**

```go
prodH := newProductHandler(deps.ProductService, deps.StorageService, deps.EnricherRegistry)
```

- [ ] **Step 5: Build and test**

```bash
make build && make test
```
Expected: clean build, all tests PASS

- [ ] **Step 6: Commit**

```bash
git add internal/handler/product_handler.go internal/handler/router.go
git commit -m "feat: add enricher support to product handler"
```

---

### Task 6: Enrich `notification_handler` — List and Get

**Files:**
- Modify: `internal/handler/notification_handler.go`
- Modify: `internal/handler/router.go` (if Task 3 Step 2 was deferred)

- [ ] **Step 1: Add `enrichers` field and update the constructor**

```go
type notificationHandler struct {
	svc       service.NotificationService
	enrichers *enricher.Registry
}

func newNotificationHandler(svc service.NotificationService, enrichers *enricher.Registry) *notificationHandler {
	return &notificationHandler{svc: svc, enrichers: enrichers}
}
```

Add import (alongside the existing `middleware` import):
```go
"github.com/hhung06/digimap-backend/internal/enricher"
```

- [ ] **Step 2: Update `List`**

Replace the current `List` method:

```go
func (h *notificationHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	ns, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceNotification)
	items := make([]any, len(ns))
	for i, n := range ns {
		items[i] = enricher.MergeInto(dto.NotificationToResponse(n), extras)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}
```

- [ ] **Step 3: Update `Get`**

Replace the current `Get` method:

```go
func (h *notificationHandler) Get(c *gin.Context) {
	venueID, _ := parseVenueID(c)
	id, err := uuid.Parse(c.Param("notifID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid notification id"))
		return
	}
	n, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceNotification)
	c.JSON(http.StatusOK, dto.OK(enricher.MergeInto(dto.NotificationToResponse(n), extras)))
}
```

- [ ] **Step 4: If Task 3 Step 2 was deferred, update the router call now**

```go
notifH := newNotificationHandler(deps.NotificationService, deps.EnricherRegistry)
```

- [ ] **Step 5: Build and test**

```bash
make build && make test
```
Expected: clean build, all tests PASS

- [ ] **Step 6: Commit**

```bash
git add internal/handler/notification_handler.go internal/handler/router.go
git commit -m "feat: add enricher support to notification handler"
```

---

### Task 7: Enrich `survey_handler` — List and Get

**Files:**
- Modify: `internal/handler/survey_handler.go`
- Modify: `internal/handler/router.go` (if Task 3 Step 2 was deferred)

- [ ] **Step 1: Add `enrichers` field and update the constructor**

```go
type surveyHandler struct {
	svc       service.SurveyService
	enrichers *enricher.Registry
}

func newSurveyHandler(svc service.SurveyService, enrichers *enricher.Registry) *surveyHandler {
	return &surveyHandler{svc: svc, enrichers: enrichers}
}
```

Add import (alongside the existing `middleware` import):
```go
"github.com/hhung06/digimap-backend/internal/enricher"
```

- [ ] **Step 2: Update `List`**

Replace the current `List` method:

```go
func (h *surveyHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	surveys, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceSurvey)
	items := make([]any, len(surveys))
	for i, s := range surveys {
		items[i] = enricher.MergeInto(dto.SurveyToResponse(s), extras)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
}
```

- [ ] **Step 3: Update `Get`**

Replace the current `Get` method:

```go
func (h *surveyHandler) Get(c *gin.Context) {
	venueID, _ := parseVenueID(c)
	id, err := uuid.Parse(c.Param("surveyID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid survey id"))
		return
	}
	s, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceSurvey)
	c.JSON(http.StatusOK, dto.OK(enricher.MergeInto(dto.SurveyToResponse(s), extras)))
}
```

- [ ] **Step 4: If Task 3 Step 2 was deferred, update the router call now**

```go
surveyH := newSurveyHandler(deps.SurveyService, deps.EnricherRegistry)
```

- [ ] **Step 5: Build and test**

```bash
make build && make test
```
Expected: clean build, all tests PASS

- [ ] **Step 6: Commit**

```bash
git add internal/handler/survey_handler.go internal/handler/router.go
git commit -m "feat: add enricher support to survey handler"
```

---

### Task 8: Enrich `ad_handler` — List and Get

**Files:**
- Modify: `internal/handler/ad_handler.go`
- Modify: `internal/handler/router.go` (if Task 3 Step 2 was deferred)

- [ ] **Step 1: Add `enrichers` field and update the constructor**

```go
type adHandler struct {
	svc       service.AdvertisementService
	enrichers *enricher.Registry
}

func newAdHandler(svc service.AdvertisementService, enrichers *enricher.Registry) *adHandler {
	return &adHandler{svc: svc, enrichers: enrichers}
}
```

Add import:
```go
"github.com/hhung06/digimap-backend/internal/enricher"
```

- [ ] **Step 2: Update `List`**

Replace the current `List` method. Note: `total` from `AdvertisementService.List` is `int`, not `int64` — keep the `int64(total)` cast.

```go
func (h *adHandler) List(c *gin.Context) {
	venueID, err := parseVenueID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid venue id"))
		return
	}
	p := paginationFromQuery(c)
	ads, total, err := h.svc.List(c.Request.Context(), venueID, p)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceAd)
	items := make([]any, len(ads))
	for i, a := range ads {
		items[i] = enricher.MergeInto(dto.AdvertisementToResponse(a), extras)
	}
	c.JSON(http.StatusOK, dto.Paginated(items, int64(total), p.Page, p.PageSize))
}
```

- [ ] **Step 3: Update `Get`**

Replace the current `Get` method:

```go
func (h *adHandler) Get(c *gin.Context) {
	venueID, _ := parseVenueID(c)
	id, err := uuid.Parse(c.Param("adID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(dto.CodeValidationError, "invalid ad id"))
		return
	}
	a, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	extras, _ := h.enrichers.EnrichForVenue(c.Request.Context(), venueID, enricher.ResourceAd)
	c.JSON(http.StatusOK, dto.OK(enricher.MergeInto(dto.AdvertisementToResponse(a), extras)))
}
```

- [ ] **Step 4: If Task 3 Step 2 was deferred, update the router call now**

```go
adH := newAdHandler(deps.AdvertisementService, deps.EnricherRegistry)
```

- [ ] **Step 5: Final build and full test suite**

```bash
make build && make test
```
Expected: clean build, all tests PASS including `internal/enricher`

- [ ] **Step 6: Commit**

```bash
git add internal/handler/ad_handler.go internal/handler/router.go
git commit -m "feat: add enricher support to ad handler"
```
