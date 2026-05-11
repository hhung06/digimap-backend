CREATE TABLE coupons (
    id          UUID PRIMARY KEY DEFAULT uuidv7(),
    venue_id    UUID REFERENCES venues(id) ON DELETE CASCADE,
    external_id VARCHAR(255),
    coupon_name VARCHAR(255),
    coupon_code VARCHAR(100),
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    issued_at   TIMESTAMPTZ,
    expired_at  TIMESTAMPTZ,
    localization JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    redeemed_at TIMESTAMPTZ,
    redeemed_by UUID REFERENCES app_users(id) ON DELETE SET NULL
);

CREATE INDEX idx_coupons_venue_id ON coupons(venue_id) WHERE deleted_at IS NULL;

SELECT create_updated_at_trigger('coupons');
