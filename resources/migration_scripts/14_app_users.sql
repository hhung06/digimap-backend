BEGIN;
SET search_path TO digimap_db, public;

-- Create app_users table (equivalent to Go migration 000025)
CREATE TABLE IF NOT EXISTS app_users (
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),
    venue_id            UUID,  -- FK to venues(id) added after char(32)->UUID conversion
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

CREATE INDEX IF NOT EXISTS app_users_venue_id_idx    ON app_users(venue_id)    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS app_users_external_id_idx ON app_users(external_id) WHERE external_id IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS app_users_email_idx       ON app_users(email)       WHERE email IS NOT NULL AND deleted_at IS NULL;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'app_users'
                   AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('app_users');
  END IF;
END $$;

-- Migrate staff users from indoormap_api_appuser
INSERT INTO app_users (
    id, venue_id, external_id, source, type,
    first_name, last_name, first_name_en, last_name_en,
    email, phone,
    company_name, company_name_en, department, position, position_en, section,
    app, token, token_expiration, survey_flag, staff_lead_flag,
    created_at, updated_at, deleted_at
)
SELECT
    uuidv7(),
    NULL,
    external_id,
    'staff',
    COALESCE(type, 0),
    first_name, last_name, first_name_en, last_name_en,
    mail, tel,
    company_name, company_name_en, department, position, position_en,
    COALESCE(section::smallint, 0::smallint),
    app, token, token_expiration,
    COALESCE(survey_flag, 0::smallint),
    COALESCE(staff_lead_flag, 0::smallint),
    COALESCE(created_at, NOW()),
    COALESCE(updated_at, NOW()),
    deleted_at
FROM indoormap_api_appuser;

-- Migrate visitors from indoormap_api_visitor
-- id and venue_id are char(32) hex strings from MySQL; convert to UUID with dashes
INSERT INTO app_users (
    id, venue_id, source, type,
    first_name, email, phone,
    visitor_type, interests, other_interests,
    is_consented, ip_address, user_agent, business_name,
    created_at, updated_at, deleted_at
)
SELECT
    (substring(id::text,1,8)||'-'||substring(id::text,9,4)||'-'||
     substring(id::text,13,4)||'-'||substring(id::text,17,4)||'-'||
     substring(id::text,21,12))::uuid,
    CASE WHEN venue_id IS NOT NULL AND venue_id != ''
         THEN (substring(venue_id::text,1,8)||'-'||substring(venue_id::text,9,4)||'-'||
               substring(venue_id::text,13,4)||'-'||substring(venue_id::text,17,4)||'-'||
               substring(venue_id::text,21,12))::uuid
         ELSE NULL END,
    'visitor',
    0::smallint,
    full_name, email,
    encode(phone_number, 'escape'),
    visitor_type::smallint,
    CASE WHEN interests IS NULL OR interests::text IN ('', 'null') THEN NULL
         ELSE interests::jsonb END,
    other_interests,
    COALESCE(is_consented, false),
    ip_address, user_agent, business_name,
    COALESCE(created_at, NOW()),
    COALESCE(updated_at, NOW()),
    deleted_at
FROM indoormap_api_visitor;

-- Drop source tables
DROP TABLE IF EXISTS indoormap_api_appuser;
DROP TABLE IF EXISTS indoormap_api_visitor;

COMMIT;
