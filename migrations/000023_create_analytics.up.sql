-- event_logs: partitioned by month via range on created_at
CREATE TABLE event_logs (
    id          UUID NOT NULL DEFAULT uuidv7(),
    venue_id    UUID,
    name        VARCHAR(50) NOT NULL DEFAULT 'view',
    params      JSONB NOT NULL DEFAULT '{}',
    device_id   VARCHAR(255) NOT NULL DEFAULT '',
    user_id     VARCHAR(255),
    user_agent  VARCHAR(512) NOT NULL DEFAULT '',
    ip_address  VARCHAR(45) NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Default partition catches overflow rows
CREATE TABLE event_logs_default PARTITION OF event_logs DEFAULT;

-- Create the first two monthly partitions automatically
DO $$
DECLARE
    start_month DATE := DATE_TRUNC('month', NOW());
    end_month   DATE := start_month + INTERVAL '1 month';
    next_start  DATE := end_month;
    next_end    DATE := end_month + INTERVAL '1 month';
    tbl1        TEXT := 'event_logs_' || TO_CHAR(start_month, 'YYYY_MM');
    tbl2        TEXT := 'event_logs_' || TO_CHAR(next_start, 'YYYY_MM');
BEGIN
    EXECUTE FORMAT(
        'CREATE TABLE IF NOT EXISTS %I PARTITION OF event_logs FOR VALUES FROM (%L) TO (%L)',
        tbl1, start_month, end_month
    );
    EXECUTE FORMAT(
        'CREATE TABLE IF NOT EXISTS %I PARTITION OF event_logs FOR VALUES FROM (%L) TO (%L)',
        tbl2, next_start, next_end
    );
END $$;

CREATE INDEX idx_event_logs_venue_id    ON event_logs (venue_id, created_at DESC);
CREATE INDEX idx_event_logs_name        ON event_logs (name, created_at DESC);
CREATE INDEX idx_event_logs_device_id   ON event_logs (device_id);

-- search_queries
CREATE TABLE search_queries (
    id           UUID PRIMARY KEY DEFAULT uuidv7(),
    venue_id     UUID,
    app_id       VARCHAR(255),
    origin       VARCHAR(255),
    search_term  VARCHAR(255) NOT NULL,
    search_count INT NOT NULL DEFAULT 0,
    last_searched TIMESTAMPTZ,
    is_promoted  BOOLEAN NOT NULL DEFAULT FALSE,
    reference    JSONB,
    status       VARCHAR(20) NOT NULL DEFAULT 'published',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_search_queries_venue_id    ON search_queries (venue_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_search_queries_search_term ON search_queries (search_term) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_search_queries_venue_term_origin
    ON search_queries (venue_id, search_term, COALESCE(origin, ''))
    WHERE deleted_at IS NULL;

SELECT create_updated_at_trigger('search_queries');
