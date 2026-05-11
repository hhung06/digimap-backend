CREATE TABLE advertisements (
    id               UUID PRIMARY KEY DEFAULT uuidv7(),
    venue_id         UUID REFERENCES venues(id) ON DELETE CASCADE,
    location_id      UUID REFERENCES locations(id) ON DELETE CASCADE,
    type             VARCHAR(50) NOT NULL DEFAULT 'dialog',
    status           VARCHAR(50) NOT NULL DEFAULT 'draft',
    navigate         VARCHAR(50),
    content_image_url VARCHAR(1000),
    content_cta_url  VARCHAR(255),
    placement        VARCHAR(255) NOT NULL DEFAULT 'home_screen',
    size_width       INT,
    size_height      INT,
    reward_type      VARCHAR(255),
    reward_amount    INT,
    display_duration INT,
    published_at     TIMESTAMPTZ,
    start_at         TIMESTAMPTZ,
    end_at           TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_ads_venue_id ON advertisements(venue_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ads_type ON advertisements(type) WHERE deleted_at IS NULL;

SELECT create_updated_at_trigger('advertisements');
