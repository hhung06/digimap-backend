-- Ensure trigger helper functions exist (may be missing from older DB setups)
CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION create_updated_at_trigger(tbl TEXT)
RETURNS VOID AS $$
BEGIN
    EXECUTE format(
        'CREATE TRIGGER set_%s_updated_at BEFORE UPDATE ON %I FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at()',
        tbl, tbl
    );
END;
$$ LANGUAGE plpgsql;

-- Event tags (global, not venue-specific)
CREATE TABLE IF NOT EXISTS event_tags (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    name        TEXT,
    localization JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_event_tags_deleted ON event_tags (deleted_at) WHERE deleted_at IS NULL;
CREATE TRIGGER set_event_tags_updated_at
    BEFORE UPDATE ON event_tags
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- Event types (per venue)
CREATE TABLE IF NOT EXISTS event_types (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    venue_id    UUID        REFERENCES venues(id) ON DELETE CASCADE,
    name        TEXT,
    localization JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_event_types_venue ON event_types (venue_id) WHERE deleted_at IS NULL;
CREATE TRIGGER set_event_types_updated_at
    BEFORE UPDATE ON event_types
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- Events
CREATE TABLE IF NOT EXISTS events (
    id               UUID        PRIMARY KEY DEFAULT uuidv7(),
    venue_id         UUID        REFERENCES venues(id) ON DELETE CASCADE,
    type_id          UUID        REFERENCES event_types(id) ON DELETE SET NULL,
    title            TEXT,
    description      TEXT,
    banner_image     TEXT,
    icon_image       TEXT,
    start_time       TIMESTAMPTZ,
    end_time         TIMESTAMPTZ,
    show_start_time  TIMESTAMPTZ,
    show_end_time    TIMESTAMPTZ,
    content_detail   TEXT,
    content_url      TEXT,
    localization     JSONB,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_events_venue ON events (venue_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_events_type  ON events (type_id)  WHERE deleted_at IS NULL;
CREATE TRIGGER set_events_updated_at
    BEFORE UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- Event <-> Tag M2M
CREATE TABLE IF NOT EXISTS event_tag_links (
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    tag_id   UUID NOT NULL REFERENCES event_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (event_id, tag_id)
);

-- Event <-> Location M2M
CREATE TABLE IF NOT EXISTS event_location_links (
    event_id    UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    PRIMARY KEY (event_id, location_id)
);

-- Event images
CREATE TABLE IF NOT EXISTS event_images (
    id         UUID        PRIMARY KEY DEFAULT uuidv7(),
    event_id   UUID        NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    image      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_event_images_event ON event_images (event_id) WHERE deleted_at IS NULL;
CREATE TRIGGER set_event_images_updated_at
    BEFORE UPDATE ON event_images
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
