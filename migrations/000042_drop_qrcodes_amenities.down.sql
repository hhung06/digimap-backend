-- Recreate the tables as defined in 000008 (amenities) and 000016 (qrcodes).
CREATE TABLE amenities (
    id                              UUID        PRIMARY KEY DEFAULT uuidv7(),
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

CREATE TABLE venue_amenities (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    venue_id    UUID        NOT NULL REFERENCES venues (id),
    amenity_id  UUID        NOT NULL REFERENCES amenities (id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    UNIQUE (venue_id, amenity_id)
);

CREATE INDEX idx_venue_amenities_venue_id   ON venue_amenities (venue_id)   WHERE deleted_at IS NULL;
CREATE INDEX idx_venue_amenities_amenity_id ON venue_amenities (amenity_id) WHERE deleted_at IS NULL;

CREATE TABLE qrcodes (
    id           UUID PRIMARY KEY DEFAULT uuidv7(),
    venue_id     UUID REFERENCES venues(id) ON DELETE CASCADE,
    level_id     UUID REFERENCES levels(id) ON DELETE SET NULL,
    location_id  UUID REFERENCES locations(id) ON DELETE SET NULL,
    lat          FLOAT,
    lng          FLOAT,
    angle        FLOAT,
    link         VARCHAR(255),
    base64_image TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_qrcodes_venue_id ON qrcodes(venue_id) WHERE deleted_at IS NULL;

SELECT create_updated_at_trigger('amenities');
SELECT create_updated_at_trigger('venue_amenities');
SELECT create_updated_at_trigger('qrcodes');
