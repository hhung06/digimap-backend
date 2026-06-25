package domain

import (
	"time"

	"github.com/google/uuid"
)

type LevelType struct {
	ID        uuid.UUID
	Name      string
	Icon      string
	Value     int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
