-- Map groups (venue-scoped grouping of levels)
CREATE TABLE map_groups (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    venue_id    UUID        NOT NULL REFERENCES venues (id),
    type        TEXT,
    name        TEXT,
    short_name  TEXT,
    sort_index  INTEGER     NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_map_groups_venue_id ON map_groups (venue_id) WHERE deleted_at IS NULL;

CREATE TRIGGER set_map_groups_updated_at
    BEFORE UPDATE ON map_groups
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Perspectives (camera settings; one-to-one with levels)
CREATE TABLE perspectives (
    id                       UUID        PRIMARY KEY DEFAULT uuidv7(),
    name                     TEXT,
    camera_zoom              DOUBLE PRECISION NOT NULL DEFAULT 0,
    camera_type              SMALLINT    NOT NULL DEFAULT 0,
    camera_max_zoom          DOUBLE PRECISION NOT NULL DEFAULT 0,
    camera_min_zoom          DOUBLE PRECISION NOT NULL DEFAULT 0,
    camera_target_center_lng DOUBLE PRECISION NOT NULL DEFAULT 0,
    camera_target_center_lat DOUBLE PRECISION NOT NULL DEFAULT 0,
    camera_target_zoom       DOUBLE PRECISION NOT NULL DEFAULT 0,
    camera_target_bearing    DOUBLE PRECISION NOT NULL DEFAULT 0,
    camera_target_pitch      DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at               TIMESTAMPTZ
);

CREATE TRIGGER set_perspectives_updated_at
    BEFORE UPDATE ON perspectives
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Levels (floors/maps within a venue)
CREATE TABLE levels (
    id             UUID        PRIMARY KEY DEFAULT uuidv7(),
    venue_id       UUID        REFERENCES venues (id) ON DELETE SET NULL,
    map_group_id   UUID        REFERENCES map_groups (id) ON DELETE SET NULL,
    perspective_id UUID        UNIQUE REFERENCES perspectives (id) ON DELETE SET NULL,
    name           TEXT,
    short_name     TEXT,
    external_id    TEXT,
    type           SMALLINT    NOT NULL DEFAULT 0,
    latitude       DOUBLE PRECISION NOT NULL DEFAULT 0,
    longitude      DOUBLE PRECISION NOT NULL DEFAULT 0,
    bearing        DOUBLE PRECISION NOT NULL DEFAULT 0,
    width          INTEGER     NOT NULL DEFAULT 0,
    height         INTEGER     NOT NULL DEFAULT 0,
    scale          DOUBLE PRECISION NOT NULL DEFAULT 0,
    level_width    DOUBLE PRECISION NOT NULL DEFAULT 0,
    level_height   DOUBLE PRECISION NOT NULL DEFAULT 0,
    file_ids       TEXT,
    elevation      INTEGER,
    is_published   BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ
);

CREATE INDEX idx_levels_venue_id      ON levels (venue_id)    WHERE deleted_at IS NULL;
CREATE INDEX idx_levels_map_group_id  ON levels (map_group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_levels_is_published  ON levels (is_published) WHERE deleted_at IS NULL;

CREATE TRIGGER set_levels_updated_at
    BEFORE UPDATE ON levels
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Geo references (control-point pairs for coordinate transforms)
CREATE TABLE geo_references (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    level_id    UUID        REFERENCES levels (id) ON DELETE CASCADE,
    control_x   INTEGER     NOT NULL DEFAULT 0,
    control_y   INTEGER     NOT NULL DEFAULT 0,
    target_x    DOUBLE PRECISION NOT NULL DEFAULT 0,
    target_y    DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_geo_references_level_id ON geo_references (level_id) WHERE deleted_at IS NULL;

CREATE TRIGGER set_geo_references_updated_at
    BEFORE UPDATE ON geo_references
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
