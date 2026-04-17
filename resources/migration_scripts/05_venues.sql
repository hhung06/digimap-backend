-- Script 05: Transform venues table

-- Rename camelCase/concatenated columns to snake_case
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='largelogo') THEN
    ALTER TABLE venues RENAME COLUMN largelogo TO large_logo;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='mediumlogo') THEN
    ALTER TABLE venues RENAME COLUMN mediumlogo TO medium_logo;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='originallogo') THEN
    ALTER TABLE venues RENAME COLUMN originallogo TO original_logo;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='smalllogo') THEN
    ALTER TABLE venues RENAME COLUMN smalllogo TO small_logo;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='publish') THEN
    ALTER TABLE venues RENAME COLUMN publish TO is_published;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='externalid') THEN
    ALTER TABLE venues RENAME COLUMN externalid TO external_id;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='workhours') THEN
    ALTER TABLE venues RENAME COLUMN workhours TO work_hours;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='seodescription') THEN
    ALTER TABLE venues RENAME COLUMN seodescription TO seo_description;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='seokeywords') THEN
    ALTER TABLE venues RENAME COLUMN seokeywords TO seo_keywords;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='seotitle') THEN
    ALTER TABLE venues RENAME COLUMN seotitle TO seo_title;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='subdomains') THEN
    ALTER TABLE venues RENAME COLUMN subdomains TO sub_domains;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='appconfigs') THEN
    ALTER TABLE venues RENAME COLUMN appconfigs TO app_configs;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='appdomains') THEN
    ALTER TABLE venues RENAME COLUMN appdomains TO app_domains;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='bodytag') THEN
    ALTER TABLE venues RENAME COLUMN bodytag TO body_tag;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='headtag') THEN
    ALTER TABLE venues RENAME COLUMN headtag TO head_tag;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='endat') THEN
    ALTER TABLE venues RENAME COLUMN endat TO end_at;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='startat') THEN
    ALTER TABLE venues RENAME COLUMN startat TO start_at;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='latitude') THEN
    ALTER TABLE venues RENAME COLUMN latitude TO lat;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venues' AND column_name='longitude') THEN
    ALTER TABLE venues RENAME COLUMN longitude TO lng;
  END IF;
END $$;

-- Drop unused columns
ALTER TABLE venues DROP COLUMN IF EXISTS countrycode;
ALTER TABLE venues DROP COLUMN IF EXISTS defaultmap;
ALTER TABLE venues DROP COLUMN IF EXISTS restored_at;
ALTER TABLE venues DROP COLUMN IF EXISTS transaction_id;

-- NOTE: type column already exists as smallint NOT NULL in the source DB — no need to add it.

-- Add new columns
ALTER TABLE venues ADD COLUMN IF NOT EXISTS slug         TEXT;
ALTER TABLE venues ADD COLUMN IF NOT EXISTS theme        JSONB;
ALTER TABLE venues ADD COLUMN IF NOT EXISTS plugins      JSONB;
ALTER TABLE venues ADD COLUMN IF NOT EXISTS translations JSONB;

-- Change lat/lng precision from double precision to NUMERIC(10,7)
ALTER TABLE venues ALTER COLUMN lat  TYPE NUMERIC(10,7) USING lat::NUMERIC(10,7);
ALTER TABLE venues ALTER COLUMN lng  TYPE NUMERIC(10,7) USING lng::NUMERIC(10,7);

-- Generate slug from name + id prefix
UPDATE venues
SET slug = lower(regexp_replace(name, '[^a-zA-Z0-9]+', '-', 'g')) || '-' || substring(id::text, 1, 8)
WHERE slug IS NULL;

ALTER TABLE venues ALTER COLUMN slug SET NOT NULL;

-- Unique index on slug (excluding soft-deleted rows)
CREATE UNIQUE INDEX IF NOT EXISTS venues_slug_unique
    ON venues (slug)
    WHERE deleted_at IS NULL;
