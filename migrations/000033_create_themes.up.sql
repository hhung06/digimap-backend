CREATE TABLE IF NOT EXISTS themes (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id        UUID        NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    name            TEXT        NOT NULL,
    primary_color   TEXT        NOT NULL DEFAULT '',
    secondary_color TEXT        NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_themes_venue ON themes (venue_id) WHERE deleted_at IS NULL;

SELECT create_updated_at_trigger('themes');
