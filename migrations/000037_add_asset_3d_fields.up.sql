ALTER TABLE assets
    ADD COLUMN IF NOT EXISTS asset_type  TEXT    NOT NULL DEFAULT '2d',
    ADD COLUMN IF NOT EXISTS description TEXT    NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS file_type   TEXT    NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS thumbnail   TEXT,
    ADD COLUMN IF NOT EXISTS material    TEXT,
    ADD COLUMN IF NOT EXISTS width       DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS height      DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS status      TEXT    NOT NULL DEFAULT 'unpublished';

CREATE INDEX IF NOT EXISTS idx_assets_venue_type ON assets (venue_id, asset_type) WHERE deleted_at IS NULL;
