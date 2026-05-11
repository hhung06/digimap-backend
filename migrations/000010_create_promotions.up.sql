CREATE TABLE promotions (
    id                    UUID        PRIMARY KEY DEFAULT uuidv7(),
    venue_id              UUID        REFERENCES venues (id) ON DELETE CASCADE,
    location_id           UUID        REFERENCES locations (id) ON DELETE CASCADE,
    external_id           TEXT,
    promo_image           TEXT,
    introduction          TEXT,
    gift_content          TEXT,
    detail_url            TEXT,
    booth_number          TEXT,
    expected_gift_count   INTEGER,
    distribution_start    TIMESTAMPTZ,
    distribution_end      TIMESTAMPTZ,
    display_type          TEXT        NOT NULL DEFAULT 'random',
    localization          JSONB,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ
);

CREATE INDEX idx_promotions_venue_id    ON promotions (venue_id)    WHERE deleted_at IS NULL;
CREATE INDEX idx_promotions_location_id ON promotions (location_id) WHERE deleted_at IS NULL;

CREATE TRIGGER set_promotions_updated_at
    BEFORE UPDATE ON promotions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
