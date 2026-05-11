package domain

import (
	"time"

	"github.com/google/uuid"
)

// AppVersion tracks the current published bundle version for a venue.
// Mirrors Django's ForceSyncVersion model (indoormap-backend/indoormap_api/app/models.py:57).
// Clients poll /app/v1/latest-bundle?venue=<uuid> to detect new bundles.
type AppVersion struct {
	ID        uuid.UUID
	VenueID   uuid.UUID
	Version   uuid.UUID
	UpdatedAt time.Time
}
