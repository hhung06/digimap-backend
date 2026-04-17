-- Script 11: Transform advertisements, articles, article_images, coupons,
--             beacons, connections, connection_levels

-- ── advertisements ───────────────────────────────────────────────────────────
-- Actual columns: content_cta_url, type, end_at, venue_id, content_image_url,
--                 displayduration, location_id, placement, reward_amount, reward_type,
--                 size_height, size_width, restored_at, transaction_id, published_at,
--                 start_at, status, navigate, article_id
ALTER TABLE advertisements RENAME COLUMN displayduration TO display_duration;
ALTER TABLE advertisements DROP COLUMN IF EXISTS article_id;
ALTER TABLE advertisements DROP COLUMN IF EXISTS restored_at;
ALTER TABLE advertisements DROP COLUMN IF EXISTS transaction_id;

-- ── articles ─────────────────────────────────────────────────────────────────
-- Actual columns: external_id, title, content, status, published_at, localization,
--                 product_id, product_category_id, venue_id, placement, navigate,
--                 published_period_end, published_period_start, created_by (varchar),
--                 location_id, application_language, company_name, email, full_name,
--                 label, landline, mobile, deleted_at, restored_at, transaction_id
ALTER TABLE articles DROP COLUMN IF EXISTS product_id;
ALTER TABLE articles DROP COLUMN IF EXISTS product_category_id;
ALTER TABLE articles DROP COLUMN IF EXISTS application_language;
ALTER TABLE articles DROP COLUMN IF EXISTS company_name;
ALTER TABLE articles DROP COLUMN IF EXISTS email;
ALTER TABLE articles DROP COLUMN IF EXISTS full_name;
ALTER TABLE articles DROP COLUMN IF EXISTS landline;
ALTER TABLE articles DROP COLUMN IF EXISTS mobile;
ALTER TABLE articles DROP COLUMN IF EXISTS created_by;
ALTER TABLE articles DROP COLUMN IF EXISTS restored_at;
ALTER TABLE articles DROP COLUMN IF EXISTS transaction_id;

-- ── article_images ───────────────────────────────────────────────────────────
-- Actual columns: image, "order" (bigint), article_id, deleted_at, restored_at, transaction_id
ALTER TABLE article_images RENAME COLUMN "order" TO sort_order;
ALTER TABLE article_images DROP COLUMN IF EXISTS restored_at;
ALTER TABLE article_images DROP COLUMN IF EXISTS transaction_id;

-- ── coupons ──────────────────────────────────────────────────────────────────
-- Actual columns: coupon_code, status, issued_at, expired_at, survey_id,
--                 venue_id, coupon_name, external_id, localization,
--                 deleted_at, restored_at, transaction_id
ALTER TABLE coupons DROP COLUMN IF EXISTS survey_id;
ALTER TABLE coupons DROP COLUMN IF EXISTS restored_at;
ALTER TABLE coupons DROP COLUMN IF EXISTS transaction_id;

-- ── beacons ──────────────────────────────────────────────────────────────────
-- Actual columns: hwid, vendorkey, lotkey, uuid, mac, radius, battery, name,
--                 positionx, positiony, level (integer), isenable, major, minor,
--                 voltage, txpower, venue_id, level_id, element_id, restored_at, transaction_id
ALTER TABLE beacons RENAME COLUMN hwid      TO hw_id;
ALTER TABLE beacons RENAME COLUMN vendorkey TO vendor_key;
ALTER TABLE beacons RENAME COLUMN lotkey    TO lot_key;
ALTER TABLE beacons RENAME COLUMN uuid      TO uuid_val;
ALTER TABLE beacons RENAME COLUMN positionx TO position_x;
ALTER TABLE beacons RENAME COLUMN positiony TO position_y;
ALTER TABLE beacons RENAME COLUMN isenable  TO is_enable;
ALTER TABLE beacons RENAME COLUMN txpower   TO tx_power;
-- Drop redundant integer 'level' column (level_id FK is the proper reference)
ALTER TABLE beacons DROP COLUMN IF EXISTS level;
ALTER TABLE beacons DROP COLUMN IF EXISTS restored_at;
ALTER TABLE beacons DROP COLUMN IF EXISTS transaction_id;

-- ── connections ──────────────────────────────────────────────────────────────
-- Actual columns: name, type, accessible, venue_id, state, status, x, y,
--                 active, externalid, restored_at, transaction_id
ALTER TABLE connections RENAME COLUMN externalid TO external_id;
ALTER TABLE connections DROP COLUMN IF EXISTS restored_at;
ALTER TABLE connections DROP COLUMN IF EXISTS transaction_id;

-- ── connection_levels ────────────────────────────────────────────────────────
-- Actual columns: active, connection_id, level_id, element_id, restored_at, transaction_id
ALTER TABLE connection_levels DROP COLUMN IF EXISTS restored_at;
ALTER TABLE connection_levels DROP COLUMN IF EXISTS transaction_id;
