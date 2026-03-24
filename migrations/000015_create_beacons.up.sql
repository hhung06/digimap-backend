CREATE TABLE IF NOT EXISTS beacons (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id    UUID        REFERENCES venues(id) ON DELETE CASCADE,
    level_id    UUID        REFERENCES levels(id) ON DELETE SET NULL,
    element_id  UUID,
    name        TEXT,
    hw_id       TEXT,
    vendor_key  TEXT,
    lot_key     TEXT,
    uuid_val    TEXT,
    mac         TEXT,
    radius      INT         NOT NULL DEFAULT 0,
    battery     INT         NOT NULL DEFAULT 0,
    position_x  FLOAT       NOT NULL DEFAULT 0,
    position_y  FLOAT       NOT NULL DEFAULT 0,
    is_enable   BOOLEAN     NOT NULL DEFAULT TRUE,
    major       INT,
    minor       INT,
    voltage     INT,
    tx_power    INT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_beacons_venue ON beacons (venue_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_beacons_level ON beacons (level_id) WHERE deleted_at IS NULL;

CREATE TRIGGER set_beacons_updated_at
    BEFORE UPDATE ON beacons
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
