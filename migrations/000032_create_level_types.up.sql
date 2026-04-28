CREATE TABLE IF NOT EXISTS level_types (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id   UUID        NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    name       TEXT        NOT NULL,
    icon       TEXT        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_level_types_venue ON level_types (venue_id) WHERE deleted_at IS NULL;

SELECT create_updated_at_trigger('level_types');
