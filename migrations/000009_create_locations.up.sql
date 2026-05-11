CREATE TABLE locations (
    id                               UUID        PRIMARY KEY DEFAULT uuidv7(),
    venue_id                         UUID        REFERENCES venues (id) ON DELETE CASCADE,
    level_id                         UUID        REFERENCES levels (id) ON DELETE SET NULL,
    main_category_id                 UUID        REFERENCES location_categories (id) ON DELETE SET NULL,
    external_id                      TEXT,
    common_hidden                    BOOLEAN     NOT NULL DEFAULT FALSE,
    common_name                      TEXT        NOT NULL DEFAULT '',
    common_short_name                TEXT,
    common_description               TEXT,
    common_color                     TEXT,
    common_location_type             SMALLINT    NOT NULL DEFAULT 0,
    common_latitude                  NUMERIC(20, 17),
    common_longitude                 NUMERIC(20, 17),
    common_address                   TEXT,
    common_location_state            SMALLINT,
    common_location_state_start_date DATE,
    common_location_state_end_date   DATE,
    common_logo                      TEXT,
    common_large_logo                TEXT,
    common_medium_logo               TEXT,
    common_small_logo                TEXT,
    common_social_website            TEXT,
    common_social_twitter            TEXT,
    common_social_tiktok             TEXT,
    common_social_facebook           TEXT,
    common_social_instagram          TEXT,
    common_contact_email             TEXT,
    common_contact_phone             TEXT,
    common_show_short_name           BOOLEAN     NOT NULL DEFAULT FALSE,
    top_logo                         TEXT,
    top_logo_type                    TEXT,
    place_work_hours                 JSONB,
    booth_number                     TEXT,
    booth_event_date                 DATE,
    booth_size                       TEXT,
    booth_services_offered           TEXT,
    booth_products_showcased         TEXT,
    person_full_name                 TEXT,
    person_job_title                 TEXT,
    room_number                      TEXT,
    room_department                  TEXT,
    room_bed_count                   INTEGER,
    room_equipment_details           TEXT,
    is_top_location                  BOOLEAN     NOT NULL DEFAULT FALSE,
    top_location_sort_index          INTEGER,
    icon_default                     TEXT        DEFAULT 'special',
    custom                           JSONB,
    localization                     JSONB,
    source                           TEXT        NOT NULL DEFAULT 'internal',
    start_time                       TIMESTAMPTZ,
    end_time                         TIMESTAMPTZ,
    is_searchable                    BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at                       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                       TIMESTAMPTZ,
    common_location_sub_type         SMALLINT    NOT NULL DEFAULT 0
);

CREATE INDEX idx_locations_venue_id         ON locations (venue_id)          WHERE deleted_at IS NULL;
CREATE INDEX idx_locations_level_id         ON locations (level_id)          WHERE deleted_at IS NULL;
CREATE INDEX idx_locations_main_category_id ON locations (main_category_id)  WHERE deleted_at IS NULL;
CREATE INDEX idx_locations_is_top           ON locations (is_top_location)   WHERE deleted_at IS NULL AND is_top_location = TRUE;
CREATE INDEX idx_locations_name_trgm        ON locations USING GIN (common_name gin_trgm_ops) WHERE deleted_at IS NULL;

CREATE TRIGGER set_locations_updated_at
    BEFORE UPDATE ON locations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Location ↔ category M2M (additional categories beyond main_category)
CREATE TABLE location_category_links (
    location_id UUID NOT NULL REFERENCES locations (id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES location_categories (id) ON DELETE CASCADE,
    PRIMARY KEY (location_id, category_id)
);

CREATE INDEX idx_location_category_links_category_id ON location_category_links (category_id);

-- Location images
CREATE TABLE location_images (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    location_id UUID        NOT NULL REFERENCES locations (id) ON DELETE CASCADE,
    original    TEXT,
    small       TEXT,
    medium      TEXT,
    large       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_location_images_location_id ON location_images (location_id) WHERE deleted_at IS NULL;

CREATE TRIGGER set_location_images_updated_at
    BEFORE UPDATE ON location_images
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
