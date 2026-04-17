-- Script 08: Transform product_categories, products, product_category_links, product_attachments

-- ── product_categories ───────────────────────────────────────────────────────
-- Actual columns: name, venue_id, localization, source, externalid, restored_at, transaction_id
ALTER TABLE product_categories DROP COLUMN IF EXISTS restored_at;
ALTER TABLE product_categories DROP COLUMN IF EXISTS transaction_id;

-- ── products ─────────────────────────────────────────────────────────────────
-- Actual columns: name, code, size, price, origin_country, expiration, description,
--                 custom, location_id, main_category_id, localization, source, image,
--                 restored_at, transaction_id, venue_id
ALTER TABLE products DROP COLUMN IF EXISTS restored_at;
ALTER TABLE products DROP COLUMN IF EXISTS transaction_id;

-- ── product_category_links ───────────────────────────────────────────────────
-- Actual columns: id (bigserial PK), product_id, productcategory_id
ALTER TABLE product_category_links RENAME COLUMN productcategory_id TO category_id;
ALTER TABLE product_category_links DROP COLUMN IF EXISTS id;
ALTER TABLE product_category_links ADD PRIMARY KEY (product_id, category_id);

-- ── product_attachments ──────────────────────────────────────────────────────
-- Actual columns: title, file_type, file, product_id, source_url, restored_at, transaction_id
ALTER TABLE product_attachments DROP COLUMN IF EXISTS restored_at;
ALTER TABLE product_attachments DROP COLUMN IF EXISTS transaction_id;
