CREATE TABLE IF NOT EXISTS notifications (
    id                         UUID        PRIMARY KEY DEFAULT uuidv7(),
    venue_id                   UUID        REFERENCES venues(id) ON DELETE CASCADE,
    survey_id                  UUID,       -- forward ref resolved after surveys table created
    title                      TEXT,
    content                    TEXT,
    topic                      TEXT,
    type                       SMALLINT    NOT NULL DEFAULT 1,  -- 1=normal 2=survey
    status                     SMALLINT    NOT NULL DEFAULT 2,  -- 1=sent 2=unsent
    send_status                SMALLINT    NOT NULL DEFAULT 0,  -- 0=pending 1=success 2=failed
    send_type                  SMALLINT    NOT NULL DEFAULT 1,  -- 1=draft 2=scheduled 3=immediate
    data                       JSONB,
    link_url                   TEXT,
    scheduled_at               TIMESTAMPTZ,
    target_app                 TEXT        NOT NULL DEFAULT 'all',
    segment_filters            JSONB,
    device_tokens              JSONB,
    error_infos                JSONB,
    retry_count                INT         NOT NULL DEFAULT 0,
    retry_at                   TIMESTAMPTZ,
    published_at               TIMESTAMPTZ,
    created_by                 UUID        REFERENCES users(id) ON DELETE SET NULL,
    created_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                 TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_notifications_venue     ON notifications (venue_id)              WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_notifications_status    ON notifications (status, send_status)   WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_notifications_scheduled ON notifications (send_type, scheduled_at) WHERE deleted_at IS NULL;

CREATE TRIGGER set_notifications_updated_at
    BEFORE UPDATE ON notifications
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
