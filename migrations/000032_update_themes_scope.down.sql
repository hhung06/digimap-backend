ALTER TABLE venues DROP COLUMN IF EXISTS theme_id;

DROP INDEX IF EXISTS idx_themes_scope;

ALTER TABLE themes
    DROP CONSTRAINT IF EXISTS valid_theme_scope_venue,
    DROP CONSTRAINT IF EXISTS themes_scope_check;

ALTER TABLE themes ALTER COLUMN venue_id SET NOT NULL;

ALTER TABLE themes
    ADD COLUMN IF NOT EXISTS primary_color   TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS secondary_color TEXT NOT NULL DEFAULT '';

UPDATE themes
SET primary_color   = COALESCE(data->>'primary_color', ''),
    secondary_color = COALESCE(data->>'secondary_color', '');

ALTER TABLE themes
    DROP COLUMN IF EXISTS data,
    DROP COLUMN IF EXISTS storage_path,
    DROP COLUMN IF EXISTS scope;
