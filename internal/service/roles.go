package service

import "github.com/hhung06/digimap-backend/internal/domain"

// Re-export domain role constants so callers wiring up the router only need
// to import the service package (avoids a second domain import in router.go).
const (
	RoleViewer      = domain.RoleViewer
	RoleEditor      = domain.RoleEditor
	RoleOwner       = domain.RoleOwner
	RoleSystemAdmin = domain.RoleSystemAdmin
)
