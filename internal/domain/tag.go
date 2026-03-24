package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID           uuid.UUID
	Name         string
	Localization json.RawMessage
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type EntityTag struct {
	ID         uuid.UUID
	TagID      uuid.UUID
	EntityType string
	EntityID   uuid.UUID
	CreatedAt  time.Time
}
