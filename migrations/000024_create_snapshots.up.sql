CREATE TABLE snapshots (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    venue_id    UUID        NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    state       SMALLINT    NOT NULL DEFAULT 0,
    method      SMALLINT    NOT NULL DEFAULT 1,
    created_by  UUID        REFERENCES users(id) ON DELETE SET NULL,
    publish_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
SELECT create_updated_at_trigger('snapshots');
CREATE INDEX idx_snapshots_venue_state ON snapshots (venue_id, state) WHERE deleted_at IS NULL;

CREATE TABLE level_bundles (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    snapshot_id UUID        NOT NULL REFERENCES snapshots(id) ON DELETE CASCADE,
    venue_id    UUID        NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    level_id    UUID        NOT NULL REFERENCES levels(id) ON DELETE RESTRICT,
    state       SMALLINT    NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
SELECT create_updated_at_trigger('level_bundles');
CREATE INDEX idx_level_bundles_snapshot ON level_bundles (snapshot_id) WHERE deleted_at IS NULL;
