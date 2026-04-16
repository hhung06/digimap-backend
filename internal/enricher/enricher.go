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
