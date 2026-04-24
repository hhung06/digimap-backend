-- Script 07: Transform location_categories, locations, location_category_links, location_images
SET search_path TO digimap_db, public;

-- ── location_categories ──────────────────────────────────────────────────────
-- Actual columns: name, icon, color, sortindex, icondefault, visible, localization,
--                 source, type, shortname, restored_at, transaction_id, externalid, description, image
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='location_categories' AND column_name='icondefault') THEN
    ALTER TABLE location_categories RENAME COLUMN icondefault TO icon_default;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='location_categories' AND column_name='sortindex') THEN
    ALTER TABLE location_categories RENAME COLUMN sortindex TO sort_index;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='location_categories' AND column_name='externalid') THEN
    ALTER TABLE location_categories RENAME COLUMN externalid TO external_id;
  END IF;
END $$;
ALTER TABLE location_categories DROP COLUMN IF EXISTS restored_at;
ALTER TABLE location_categories DROP COLUMN IF EXISTS transaction_id;
ALTER TABLE location_categories ALTER COLUMN sort_index TYPE INTEGER USING sort_index::INTEGER;

-- ── locations ────────────────────────────────────────────────────────────────
-- Actual columns: common_name, externalid, common_description, etc.
--   Columns to drop: restored_at, transaction_id, top_logo, top_logo_type,
--                    start_time, end_time, common_show_short_name, common_sub_type
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='locations' AND column_name='externalid') THEN
    ALTER TABLE locations RENAME COLUMN externalid TO external_id;
  END IF;
END $$;
ALTER TABLE locations DROP COLUMN IF EXISTS restored_at;
ALTER TABLE locations DROP COLUMN IF EXISTS transaction_id;
ALTER TABLE locations DROP COLUMN IF EXISTS top_logo;
ALTER TABLE locations DROP COLUMN IF EXISTS top_logo_type;
ALTER TABLE locations DROP COLUMN IF EXISTS start_time;
ALTER TABLE locations DROP COLUMN IF EXISTS end_time;
ALTER TABLE locations DROP COLUMN IF EXISTS common_show_short_name;
ALTER TABLE locations DROP COLUMN IF EXISTS common_sub_type;

-- ── location_category_links ──────────────────────────────────────────────────
-- Actual columns: id (bigserial PK), location_id, locationcategory_id
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='location_category_links' AND column_name='locationcategory_id') THEN
    ALTER TABLE location_category_links RENAME COLUMN locationcategory_id TO category_id;
  END IF;
END $$;
ALTER TABLE location_category_links DROP COLUMN IF EXISTS id;
ALTER TABLE location_category_links ADD PRIMARY KEY (location_id, category_id);

-- ── location_images ──────────────────────────────────────────────────────────
-- Actual columns: location_id, large, medium, original, small, restored_at, transaction_id
ALTER TABLE location_images DROP COLUMN IF EXISTS restored_at;
ALTER TABLE location_images DROP COLUMN IF EXISTS transaction_id;
