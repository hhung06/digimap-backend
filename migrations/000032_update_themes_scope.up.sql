-- Add scope/global support to themes and theme selection on venues

ALTER TABLE themes
    ADD COLUMN IF NOT EXISTS scope        TEXT NOT NULL DEFAULT 'custom',
    ADD COLUMN IF NOT EXISTS storage_path TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS data         JSONB NOT NULL DEFAULT '{}';

-- Migrate existing color data into the data JSONB blob
UPDATE themes
SET data = jsonb_build_object(
    'primary_color',   primary_color,
    'secondary_color', secondary_color
)
WHERE primary_color <> '' OR secondary_color <> '';

ALTER TABLE themes
    DROP COLUMN IF EXISTS primary_color,
    DROP COLUMN IF EXISTS secondary_color;

ALTER TABLE themes
    ALTER COLUMN venue_id DROP NOT NULL;

ALTER TABLE themes
    ADD CONSTRAINT themes_scope_check CHECK (scope IN ('global', 'custom')),
    ADD CONSTRAINT valid_theme_scope_venue CHECK (
        (scope = 'global' AND venue_id IS NULL) OR
        (scope = 'custom' AND venue_id IS NOT NULL)
    );

CREATE INDEX IF NOT EXISTS idx_themes_scope ON themes (scope) WHERE deleted_at IS NULL;

-- Allow venues to reference which theme they use
ALTER TABLE venues
    ADD COLUMN IF NOT EXISTS theme_id UUID REFERENCES themes(id) ON DELETE SET NULL;
