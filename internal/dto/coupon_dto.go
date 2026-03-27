package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type CouponResponse struct {
	ID           uuid.UUID       `json:"id"`
	VenueID      *uuid.UUID      `json:"venue_id,omitempty"`
	ExternalID   string          `json:"external_id,omitempty"`
	CouponName   string          `json:"coupon_name,omitempty"`
	CouponCode   string          `json:"coupon_code,omitempty"`
	Status       string          `json:"status"`
	IssuedAt     *time.Time      `json:"issued_at,omitempty"`
	ExpiredAt    *time.Time      `json:"expired_at,omitempty"`
	Localization json.RawMessage `json:"localization,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type CouponRequest struct {
	ExternalID   string          `json:"external_id"`
	CouponName   string          `json:"coupon_name"`
	CouponCode   string          `json:"coupon_code"`
	Status       string          `json:"status"`
	IssuedAt     *time.Time      `json:"issued_at"`
	ExpiredAt    *time.Time      `json:"expired_at"`
	Localization json.RawMessage `json:"localization"`
}

func CouponToResponse(c *domain.Coupon) CouponResponse {
	return CouponResponse{
		ID: c.ID, VenueID: c.VenueID, ExternalID: c.ExternalID,
		CouponName: c.CouponName, CouponCode: c.CouponCode,
		Status: c.Status, IssuedAt: c.IssuedAt, ExpiredAt: c.ExpiredAt,
		Localization: c.Localization,
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
}
