package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	LevelBundleStatePending = 0 // reserved
	LevelBundleStateDraft   = 1
	LevelBundleStatePublic  = 2
)

// LevelBundle is the bundled map data artifact for a single level,
// associated with a parent Snapshot. Cleaned up automatically when
// the parent snapshot is hard-deleted (ON DELETE CASCADE).
type LevelBundle struct {
	ID         uuid.UUID
	SnapshotID uuid.UUID // FK → snapshots(id) ON DELETE CASCADE
	VenueID    uuid.UUID // FK → venues(id)
	LevelID    uuid.UUID // FK → levels(id)
	State      int       // 0=pending 1=draft 2=public
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
