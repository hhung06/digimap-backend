package tenants

import "github.com/hhung06/digimap-backend/internal/enricher"

// RegisterAll registers all tenant-specific enrichers with the registry.
// Add one call per tenant as new tenants are onboarded.
func RegisterAll(_ *enricher.Registry) {}
