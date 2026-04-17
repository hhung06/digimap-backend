-- Script 05: Transform venues table

-- Rename camelCase/concatenated columns to snake_case
ALTER TABLE venues RENAME COLUMN largelogo      TO large_logo;
ALTER TABLE venues RENAME COLUMN mediumlogo     TO medium_logo;
ALTER TABLE venues RENAME COLUMN originallogo   TO original_logo;
ALTER TABLE venues RENAME COLUMN smalllogo      TO small_logo;
ALTER TABLE venues RENAME COLUMN publish        TO is_published;
ALTER TABLE venues RENAME COLUMN externalid     TO external_id;
ALTER TABLE venues RENAME COLUMN workhours      TO work_hours;
ALTER TABLE venues RENAME COLUMN seodescription TO seo_description;
ALTER TABLE venues RENAME COLUMN seokeywords    TO seo_keywords;
ALTER TABLE venues RENAME COLUMN seotitle       TO seo_title;
ALTER TABLE venues RENAME COLUMN subdomains     TO sub_domains;
ALTER TABLE venues RENAME COLUMN appconfigs     TO app_configs;
ALTER TABLE venues RENAME COLUMN appdomains     TO app_domains;
ALTER TABLE venues RENAME COLUMN bodytag        TO body_tag;
ALTER TABLE venues RENAME COLUMN headtag        TO head_tag;
ALTER TABLE venues RENAME COLUMN endat          TO end_at;
ALTER TABLE venues RENAME COLUMN startat        TO start_at;
ALTER TABLE venues RENAME COLUMN latitude       TO lat;
ALTER TABLE venues RENAME COLUMN longitude      TO lng;

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
