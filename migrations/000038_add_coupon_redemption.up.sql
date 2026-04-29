ALTER TABLE coupons
    ADD COLUMN redeemed_at TIMESTAMPTZ,
    ADD COLUMN redeemed_by UUID REFERENCES app_users(id) ON DELETE SET NULL;
