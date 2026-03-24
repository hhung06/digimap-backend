CREATE TABLE videos (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id    UUID REFERENCES venues(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL DEFAULT '',
    description TEXT,
    url         VARCHAR(1000),
    thumbnail   VARCHAR(1000),
    duration    INT,
    status      VARCHAR(50) NOT NULL DEFAULT 'draft',
    published_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_videos_venue_id ON videos(venue_id) WHERE deleted_at IS NULL;

SELECT create_updated_at_trigger('videos');
