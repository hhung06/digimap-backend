CREATE TABLE connections (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id    UUID REFERENCES venues(id) ON DELETE CASCADE,
    external_id VARCHAR(255),
    name        VARCHAR(255),
    type        SMALLINT NOT NULL DEFAULT 1,
    x           FLOAT NOT NULL DEFAULT 0,
    y           FLOAT NOT NULL DEFAULT 0,
    state       SMALLINT NOT NULL DEFAULT 1,
    status      SMALLINT NOT NULL DEFAULT 1,
    accessible  BOOLEAN NOT NULL DEFAULT TRUE,
    active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE TABLE connection_levels (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    connection_id UUID NOT NULL REFERENCES connections(id) ON DELETE CASCADE,
    level_id      UUID REFERENCES levels(id) ON DELETE CASCADE,
    element_id    UUID,
    active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE TABLE qrcodes (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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

CREATE INDEX idx_connections_venue_id ON connections(venue_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_connection_levels_connection_id ON connection_levels(connection_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_qrcodes_venue_id ON qrcodes(venue_id) WHERE deleted_at IS NULL;

SELECT create_updated_at_trigger('connections');
SELECT create_updated_at_trigger('connection_levels');
SELECT create_updated_at_trigger('qrcodes');
