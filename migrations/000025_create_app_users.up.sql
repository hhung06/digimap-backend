CREATE TABLE app_users (
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),
    venue_id            UUID REFERENCES venues(id) ON DELETE SET NULL,
    external_id         TEXT,
    source              TEXT NOT NULL,
    type                SMALLINT NOT NULL DEFAULT 0,

    first_name          TEXT,
    last_name           TEXT,
    first_name_en       TEXT,
    last_name_en        TEXT,
    email               TEXT,
    phone               TEXT,

    company_name        TEXT,
    company_name_en     TEXT,
    department          TEXT,
    position            TEXT,
    position_en         TEXT,
    section             SMALLINT,
    app                 TEXT,
    token               TEXT,
    token_expiration    TIMESTAMPTZ,
    survey_flag         SMALLINT NOT NULL DEFAULT 0,
    staff_lead_flag     SMALLINT NOT NULL DEFAULT 0,

    visitor_type        SMALLINT,
    interests           JSONB,
    other_interests     TEXT,
    is_consented        BOOLEAN NOT NULL DEFAULT FALSE,
    ip_address          TEXT,
    user_agent          TEXT,
    business_name       TEXT,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

CREATE INDEX app_users_venue_id_idx    ON app_users(venue_id)    WHERE deleted_at IS NULL;
CREATE INDEX app_users_external_id_idx ON app_users(external_id) WHERE external_id IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX app_users_email_idx       ON app_users(email)       WHERE email IS NOT NULL AND deleted_at IS NULL;
SELECT create_updated_at_trigger('app_users');
