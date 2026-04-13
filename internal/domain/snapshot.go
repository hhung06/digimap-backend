package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	SnapshotStatePending = 0 // reserved — never written in Phase 1
	SnapshotStateDraft   = 1
	SnapshotStatePublic  = 2

	SnapshotMethodManual = 1
	SnapshotMethodAuto   = 2

	MaxSnapshotVersions = 5
)

// Snapshot is a point-in-time capture of a venue's map bundle data.
// The bundle JSON is stored in S3; this struct holds only metadata.
type Snapshot struct {
	ID        uuid.UUID
	VenueID   uuid.UUID  // FK → venues(id)
	State     int        // 0=pending 1=draft 2=public
	Method    int        // 1=manual 2=auto
	CreatedBy *uuid.UUID // FK → users(id)
	PublishAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
