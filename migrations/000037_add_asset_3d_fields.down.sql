DROP INDEX IF EXISTS idx_assets_venue_type;

ALTER TABLE assets
    DROP COLUMN IF EXISTS asset_type,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS file_type,
    DROP COLUMN IF EXISTS thumbnail,
    DROP COLUMN IF EXISTS material,
    DROP COLUMN IF EXISTS width,
    DROP COLUMN IF EXISTS height,
    DROP COLUMN IF EXISTS status;
