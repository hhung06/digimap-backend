package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Coupon struct {
	ID           uuid.UUID
	VenueID      *uuid.UUID
	ExternalID   string
	CouponName   string
	CouponCode   string
	Status       string
	IssuedAt     *time.Time
	ExpiredAt    *time.Time
	Localization json.RawMessage
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
