CREATE TABLE coupon_users (
    id         UUID PRIMARY KEY DEFAULT uuidv7(),
    coupon_id  UUID NOT NULL REFERENCES coupons(id) ON DELETE CASCADE,
    user_id    UUID REFERENCES app_users(id) ON DELETE SET NULL,
    device_id  TEXT NOT NULL DEFAULT '',
    is_used    BOOLEAN NOT NULL DEFAULT FALSE,
    used_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_coupon_users_coupon_id ON coupon_users(coupon_id);
CREATE INDEX idx_coupon_users_user_id   ON coupon_users(user_id) WHERE user_id IS NOT NULL;
CREATE UNIQUE INDEX uq_coupon_users_coupon_user
    ON coupon_users(coupon_id, user_id)
    WHERE user_id IS NOT NULL;
CREATE UNIQUE INDEX uq_coupon_users_coupon_device
    ON coupon_users(coupon_id, device_id)
    WHERE device_id <> '';

SELECT create_updated_at_trigger('coupon_users');
