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
	RedeemedAt   *time.Time
	RedeemedBy   *uuid.UUID
	Localization json.RawMessage
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time

	// IsUsed is a transient field populated in user-context queries (ListByUser).
	// It reflects the coupon_users.is_used value for the requesting user.
	IsUsed bool
}

type CouponUser struct {
	ID        uuid.UUID
	CouponID  uuid.UUID
	UserID    *uuid.UUID
	DeviceID  string
	IsUsed    bool
	UsedAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
