CREATE TABLE venues (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id     UUID        NOT NULL REFERENCES customers (id),
    name            TEXT        NOT NULL,
    slug            TEXT        UNIQUE,
    external_id     TEXT,
    type            SMALLINT    NOT NULL DEFAULT 0,
    public_key      CHAR(43)    UNIQUE NOT NULL,
    private_key     CHAR(64)    UNIQUE NOT NULL,
    address         TEXT,
    city            TEXT,
    state           TEXT,
    country         TEXT,
    postal          TEXT,
    lat             NUMERIC(10, 7),
    lng             NUMERIC(10, 7),
    timezone        TEXT        NOT NULL DEFAULT 'UTC',
    telephone       TEXT,
    work_hours      TEXT,
    description     TEXT,
    is_published    BOOLEAN     NOT NULL DEFAULT FALSE,
    theme           JSONB,
    plugins         JSONB,
    translations    JSONB,
    localization    JSONB,
    custom_data     JSONB,
    app_configs     JSONB,
    app_domains     JSONB,
    sub_domains     TEXT,
    seo_title       TEXT,
    seo_description TEXT,
    seo_keywords    TEXT,
    head_tag        TEXT,
    body_tag        TEXT,
    original_logo   TEXT,
    small_logo      TEXT,
    medium_logo     TEXT,
    large_logo      TEXT,
    start_at        TIMESTAMPTZ,
    end_at          TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_venues_customer_id  ON venues (customer_id)  WHERE deleted_at IS NULL;
CREATE INDEX idx_venues_public_key   ON venues (public_key)   WHERE deleted_at IS NULL;
CREATE INDEX idx_venues_is_published ON venues (is_published) WHERE deleted_at IS NULL;
CREATE INDEX idx_venues_slug         ON venues (slug)         WHERE deleted_at IS NULL;
CREATE INDEX idx_venues_name_trgm    ON venues USING GIN (name gin_trgm_ops) WHERE deleted_at IS NULL;

CREATE TRIGGER set_venues_updated_at
    BEFORE UPDATE ON venues
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
