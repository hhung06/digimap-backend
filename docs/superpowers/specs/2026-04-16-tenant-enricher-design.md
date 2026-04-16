# Tenant Enricher — Design

**Date:** 2026-04-16
**Scope:** Per-tenant computed/derived field extension for Location, Product, Notification, Survey, Advertisement responses
**Status:** Approved

---

## Background

Different customers (tenants) need extra computed or derived fields appended to certain API responses. The base fields are always returned unchanged — tenant logic only adds on top. The same API is consumed by multiple platforms (web app, mobile app), so field naming must stay consistent across all tenants; only additive extras are allowed.

The existing `Customer` entity is the tenant unit. Each `Venue` has a `CustomerID` FK, so the tenant is always implicit from the venue on any venue-scoped request.

---

## Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Enrichment approach | Registered `EnricherFunc` per `(customerID, resource)` | Compile-time, explicit, no DB config table needed |
| Tenant resolution | Registry-internal via `VenueRepository.GetCustomerID` | No middleware needed; DB call skipped entirely when registry is empty |
| Extra fields shape | Flat `map[string]any` merged into base DTO | Simpler than typed interfaces; base fields remain typed; extras are additive only |
| List vs detail cost | One enricher call per request, extras applied uniformly to all items | Extras are venue-scoped (not item-scoped), so per-item calls would be redundant |
| Enricher error handling | Log and continue — return base response unchanged | A failed enricher must never break the base response |
| Tenant registration | One file per tenant in `cmd/tenants/`, registered in `cmd/serve.go` | Adding a tenant = new file + one line in serve.go; no other files touched |
| Extensibility | `Resource` is a typed `string` alias | Adding a new resource = one new `const` line, no registry or interface changes |

---

## Package: `internal/enricher`

### `enricher.go`

```go
type Resource string

const (
    ResourceLocation     Resource = "location"
    ResourceProduct      Resource = "product"
    ResourceNotification Resource = "notification"
    ResourceSurvey       Resource = "survey"
    ResourceAd           Resource = "advertisement"
)

// EnricherFunc computes extra fields for a specific tenant+resource+venue combination.
// The same extras are applied to every item in a list response (one call per request).
// Return nil, nil if no extras are needed.
type EnricherFunc func(ctx context.Context, venueID uuid.UUID) (map[string]any, error)

type enricherKey struct {
    customerID uuid.UUID
    resource   Resource
}

type Registry struct {
    m         map[enricherKey]EnricherFunc
    venueRepo repository.VenueRepository
}

func NewRegistry(venueRepo repository.VenueRepository) *Registry

func (r *Registry) Register(customerID uuid.UUID, resource Resource, fn EnricherFunc)

// EnrichForVenue resolves the tenant from venueID and calls the registered enricher.
// Returns nil, nil immediately if no enrichers are registered (zero DB cost).
func (r *Registry) EnrichForVenue(ctx context.Context, venueID uuid.UUID, resource Resource) (map[string]any, error)
```

`VenueRepository` gains one new method:
```go
GetCustomerID(ctx context.Context, venueID uuid.UUID) (uuid.UUID, error)
// SELECT customer_id FROM venues WHERE id = $1 AND deleted_at IS NULL
```

### `merge.go`

```go
// MergeInto merges extras on top of base and returns the result as any.
// Returns base unchanged (no JSON round-trip) if extras is nil or empty.
func MergeInto(base any, extras map[string]any) any
```

---

## Handler Integration

The same pattern applies to all five handlers. Constructor gains `*enricher.Registry`:

```go
type locationHandler struct {
    categorySvc service.LocationCategoryService
    locationSvc service.LocationService
    enrichers   *enricher.Registry
}
```

**List endpoint** — one enricher call, applied to all items:
```go
extras, _ := h.enrichers.EnrichForVenue(ctx, venueID, enricher.ResourceLocation)
items := make([]any, len(locs))
for i, loc := range locs {
    items[i] = enricher.MergeInto(dto.LocationToResponse(loc), extras)
}
c.JSON(http.StatusOK, dto.Paginated(items, total, p.Page, p.PageSize))
```

**Detail endpoint:**
```go
extras, _ := h.enrichers.EnrichForVenue(ctx, venueID, enricher.ResourceLocation)
c.JSON(http.StatusOK, dto.OK(enricher.MergeInto(dto.LocationToResponse(loc), extras)))
```

Enricher errors are logged and ignored — the base response is always returned.

---

## Registration

```
cmd/tenants/
    acmecorp.go
    othertenant.go
```

Each file exports a single `Register` function:

```go
// cmd/tenants/acmecorp.go
func RegisterAcmeCorp(r *enricher.Registry) {
    r.Register(acmeCorpID, enricher.ResourceLocation, func(ctx context.Context, venueID uuid.UUID) (map[string]any, error) {
        return map[string]any{
            "crm_zone":     "west",
            "custom_label": "Gate",
        }, nil
    })
}
```

In `cmd/serve.go`:
```go
enricherRegistry := enricher.NewRegistry(venueRepo)
tenants.RegisterAcmeCorp(enricherRegistry)
```

---

## Files to Create / Modify

| Action | File |
|---|---|
| Create | `internal/enricher/enricher.go` |
| Create | `internal/enricher/merge.go` |
| Create | `internal/enricher/enricher_test.go` |
| Modify | `internal/repository/interfaces.go` — add `GetCustomerID` to `VenueRepository` |
| Modify | `internal/repository/postgres/venue_repo.go` — implement `GetCustomerID` |
| Modify | `internal/repository/mocks/mocks.go` — add `GetCustomerID` mock method |
| Modify | `internal/handler/location_handler.go` — add enricher, update list/detail endpoints |
| Modify | `internal/handler/ad_handler.go` — same |
| Modify | `internal/handler/notification_handler.go` — same |
| Modify | `internal/handler/survey_handler.go` — same |
| Modify | `internal/handler/product_handler.go` — same |
| Modify | `internal/handler/router.go` — pass registry to affected handler constructors |
| Modify | `cmd/serve.go` — create registry, wire venueRepo, call tenant registrations |
| Create | `cmd/tenants/` — one file per tenant (initially empty/placeholder) |

---

## Out of Scope

- Per-item enrichment (extras are venue-scoped, uniform across list items)
- Database-driven enricher configuration
- Request enrichment (mutating inbound fields per tenant)
- Enrichment for resources outside the five listed
