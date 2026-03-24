-- Location categories (venue-scoped)
CREATE TABLE location_categories (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id     UUID        NOT NULL REFERENCES venues (id),
    external_id  TEXT,
    name         TEXT,
    short_name   TEXT,
    color        TEXT,
    icon         TEXT,
    icon_default TEXT,
    sort_index   INTEGER     NOT NULL DEFAULT 0,
    visible      BOOLEAN     NOT NULL DEFAULT TRUE,
    description  TEXT,
    type         TEXT,
    image        TEXT,
    localization JSONB,
    source       TEXT        NOT NULL DEFAULT 'internal',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_location_categories_venue_id ON location_categories (venue_id) WHERE deleted_at IS NULL;

CREATE TRIGGER set_location_categories_updated_at
    BEFORE UPDATE ON location_categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Amenities (location templates / master amenity definitions)
CREATE TABLE amenities (
    id                              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    common_name                     TEXT        NOT NULL DEFAULT '',
    common_short_name               TEXT,
    common_description              TEXT,
    common_color                    TEXT,
    common_location_type            SMALLINT    NOT NULL DEFAULT 0,
    common_latitude                 NUMERIC(20, 17),
    common_longitude                NUMERIC(20, 17),
    common_address                  TEXT,
    common_location_state           SMALLINT,
    common_location_state_start_date DATE,
    common_location_state_end_date  DATE,
    common_logo                     TEXT,
    common_social_website           TEXT,
    common_social_twitter           TEXT,
    common_social_tiktok            TEXT,
    common_social_facebook          TEXT,
    common_social_instagram         TEXT,
    common_contact_email            TEXT,
    common_contact_phone            TEXT,
    place_work_hours                JSONB,
    booth_number                    TEXT,
    booth_event_date                DATE,
    booth_size                      TEXT,
    booth_services_offered          TEXT,
    booth_products_showcased        TEXT,
    person_full_name                TEXT,
    person_job_title                TEXT,
    room_number                     TEXT,
    room_department                 TEXT,
    room_bed_count                  INTEGER,
    room_equipment_details          TEXT,
    localization                    JSONB,
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                      TIMESTAMPTZ
);

CREATE TRIGGER set_amenities_updated_at
    BEFORE UPDATE ON amenities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Venue ↔ amenity links (which amenities are available at a venue)
CREATE TABLE venue_amenities (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id    UUID        NOT NULL REFERENCES venues (id),
    amenity_id  UUID        NOT NULL REFERENCES amenities (id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    UNIQUE (venue_id, amenity_id)
);

CREATE INDEX idx_venue_amenities_venue_id   ON venue_amenities (venue_id)   WHERE deleted_at IS NULL;
CREATE INDEX idx_venue_amenities_amenity_id ON venue_amenities (amenity_id) WHERE deleted_at IS NULL;

CREATE TRIGGER set_venue_amenities_updated_at
    BEFORE UPDATE ON venue_amenities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
