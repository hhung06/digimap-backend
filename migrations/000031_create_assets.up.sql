CREATE TABLE IF NOT EXISTS assets (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id     UUID        NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    name         TEXT        NOT NULL,
    key          TEXT        NOT NULL,
    content_type TEXT        NOT NULL DEFAULT '',
    size_bytes   BIGINT      NOT NULL DEFAULT 0,
    url          TEXT        NOT NULL DEFAULT '',
    created_by   UUID        REFERENCES users(id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_assets_venue ON assets (venue_id) WHERE deleted_at IS NULL;

SELECT create_updated_at_trigger('assets');
