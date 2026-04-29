ALTER TABLE coupons
    DROP COLUMN IF EXISTS redeemed_at,
    DROP COLUMN IF EXISTS redeemed_by;
