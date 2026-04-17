BEGIN;

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
    CASE WHEN interests IS NULL OR interests = '' THEN NULL
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
