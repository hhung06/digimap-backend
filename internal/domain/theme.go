package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	ThemeScopeGlobal = "global"
	ThemeScopeCustom = "custom"
)

type Theme struct {
	ID          uuid.UUID
	VenueID     *uuid.UUID
	Scope       string
	Name        string
	Data        json.RawMessage
	StoragePath string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
