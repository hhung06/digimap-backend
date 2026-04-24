-- Script 16: Convert all char(32) / bigserial PKs and FK columns to UUID
SET search_path TO digimap_db, public;

-- =============================================================================
-- Helper function and procedures
-- =============================================================================

CREATE OR REPLACE FUNCTION hex_to_uuid(hex text) RETURNS uuid LANGUAGE sql IMMUTABLE AS $$
  SELECT CASE WHEN hex IS NULL OR length(btrim(hex)) != 32 THEN NULL
         ELSE (substring(hex,1,8)||'-'||substring(hex,9,4)||'-'||
               substring(hex,13,4)||'-'||substring(hex,17,4)||'-'||
               substring(hex,21,12))::uuid END
$$;

-- Convert char/varchar PK column (hex string) to UUID
CREATE OR REPLACE PROCEDURE conv_pk(tbl TEXT) LANGUAGE plpgsql AS $$
DECLARE cname TEXT;
BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns
    WHERE table_schema = current_schema() AND table_name = tbl AND column_name = 'id'
    AND data_type IN ('character', 'character varying')) THEN RETURN; END IF;
  EXECUTE format('ALTER TABLE %I ADD COLUMN _new_id UUID', tbl);
  EXECUTE format('UPDATE %I SET _new_id = hex_to_uuid(id::text)', tbl);
  -- Any id that wasn't a 32-char hex (NULL result) gets a fresh UUID
  EXECUTE format('UPDATE %I SET _new_id = uuidv7() WHERE _new_id IS NULL', tbl);
  SELECT constraint_name INTO cname FROM information_schema.table_constraints
    WHERE table_schema = current_schema() AND table_name = tbl AND constraint_type = 'PRIMARY KEY';
  IF cname IS NOT NULL THEN EXECUTE format('ALTER TABLE %I DROP CONSTRAINT %I', tbl, cname); END IF;
  EXECUTE format('ALTER TABLE %I DROP COLUMN id', tbl);
  EXECUTE format('ALTER TABLE %I RENAME COLUMN _new_id TO id', tbl);
  EXECUTE format('ALTER TABLE %I ADD PRIMARY KEY (id)', tbl);
END $$;

-- Convert bigserial PK to UUID (new random id — serial has no meaningful hex to preserve)
CREATE OR REPLACE PROCEDURE conv_pk_serial(tbl TEXT) LANGUAGE plpgsql AS $$
DECLARE cname TEXT;
BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns
    WHERE table_schema = current_schema() AND table_name = tbl AND column_name = 'id'
    AND data_type IN ('bigint', 'integer')) THEN RETURN; END IF;
  EXECUTE format('ALTER TABLE %I ADD COLUMN _new_id UUID DEFAULT uuidv7()', tbl);
  SELECT constraint_name INTO cname FROM information_schema.table_constraints
    WHERE table_schema = current_schema() AND table_name = tbl AND constraint_type = 'PRIMARY KEY';
  IF cname IS NOT NULL THEN EXECUTE format('ALTER TABLE %I DROP CONSTRAINT %I', tbl, cname); END IF;
  EXECUTE format('ALTER TABLE %I DROP COLUMN id', tbl);
  EXECUTE format('ALTER TABLE %I RENAME COLUMN _new_id TO id', tbl);
  EXECUTE format('ALTER TABLE %I ADD PRIMARY KEY (id)', tbl);
END $$;

-- Convert char/varchar FK column (hex string) to UUID
CREATE OR REPLACE PROCEDURE conv_fk(tbl TEXT, col TEXT) LANGUAGE plpgsql AS $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns
    WHERE table_schema = current_schema() AND table_name = tbl AND column_name = col
    AND data_type IN ('character', 'character varying')) THEN RETURN; END IF;
  EXECUTE format('ALTER TABLE %I ADD COLUMN _new_col UUID', tbl);
  EXECUTE format('UPDATE %I SET _new_col = hex_to_uuid(%I::text)', tbl, col);
  EXECUTE format('ALTER TABLE %I DROP COLUMN %I', tbl, col);
  EXECUTE format('ALTER TABLE %I RENAME COLUMN _new_col TO %I', tbl, col);
END $$;

-- Add FK constraint if it doesn't exist
CREATE OR REPLACE PROCEDURE add_fk(
  con TEXT, tbl TEXT, col TEXT, ref_tbl TEXT, ref_col TEXT, on_delete TEXT
) LANGUAGE plpgsql AS $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
    WHERE constraint_name = con AND table_schema = current_schema()) THEN
    EXECUTE format(
      'ALTER TABLE %I ADD CONSTRAINT %I FOREIGN KEY (%I) REFERENCES %I(%I) ON DELETE %s',
      tbl, con, col, ref_tbl, ref_col, on_delete);
  END IF;
END $$;

-- =============================================================================
-- SECTION 1: Drop all FK constraints (re-added at end with UUID column types)
-- =============================================================================

ALTER TABLE venues                  DROP CONSTRAINT IF EXISTS venues_customer_id_fk;
ALTER TABLE venue_user_roles        DROP CONSTRAINT IF EXISTS venue_user_roles_venue_id_fk;
ALTER TABLE venue_invitations       DROP CONSTRAINT IF EXISTS venue_invitations_venue_id_fk;
ALTER TABLE map_groups              DROP CONSTRAINT IF EXISTS map_groups_venue_id_fk;
ALTER TABLE levels                  DROP CONSTRAINT IF EXISTS levels_venue_id_fk;
ALTER TABLE levels                  DROP CONSTRAINT IF EXISTS levels_map_group_id_fk;
ALTER TABLE levels                  DROP CONSTRAINT IF EXISTS levels_perspective_id_fk;
ALTER TABLE geo_references          DROP CONSTRAINT IF EXISTS geo_references_level_id_fk;
ALTER TABLE location_categories     DROP CONSTRAINT IF EXISTS location_categories_venue_id_fk;
ALTER TABLE locations               DROP CONSTRAINT IF EXISTS locations_venue_id_fk;
ALTER TABLE locations               DROP CONSTRAINT IF EXISTS locations_level_id_fk;
ALTER TABLE locations               DROP CONSTRAINT IF EXISTS locations_main_category_id_fk;
ALTER TABLE location_category_links DROP CONSTRAINT IF EXISTS location_category_links_location_id_fk;
ALTER TABLE location_category_links DROP CONSTRAINT IF EXISTS location_category_links_category_id_fk;
ALTER TABLE location_images         DROP CONSTRAINT IF EXISTS location_images_location_id_fk;
ALTER TABLE product_categories      DROP CONSTRAINT IF EXISTS product_categories_venue_id_fk;
ALTER TABLE products                DROP CONSTRAINT IF EXISTS products_venue_id_fk;
ALTER TABLE products                DROP CONSTRAINT IF EXISTS products_location_id_fk;
ALTER TABLE products                DROP CONSTRAINT IF EXISTS products_main_category_id_fk;
ALTER TABLE product_category_links  DROP CONSTRAINT IF EXISTS product_category_links_product_id_fk;
ALTER TABLE product_category_links  DROP CONSTRAINT IF EXISTS product_category_links_category_id_fk;
ALTER TABLE product_attachments     DROP CONSTRAINT IF EXISTS product_attachments_product_id_fk;
ALTER TABLE event_types             DROP CONSTRAINT IF EXISTS event_types_venue_id_fk;
ALTER TABLE events                  DROP CONSTRAINT IF EXISTS events_venue_id_fk;
ALTER TABLE events                  DROP CONSTRAINT IF EXISTS events_type_id_fk;
ALTER TABLE event_tag_links         DROP CONSTRAINT IF EXISTS event_tag_links_event_id_fk;
ALTER TABLE event_tag_links         DROP CONSTRAINT IF EXISTS event_tag_links_tag_id_fk;
ALTER TABLE event_location_links    DROP CONSTRAINT IF EXISTS event_location_links_event_id_fk;
ALTER TABLE event_location_links    DROP CONSTRAINT IF EXISTS event_location_links_location_id_fk;
ALTER TABLE event_images            DROP CONSTRAINT IF EXISTS event_images_event_id_fk;
ALTER TABLE surveys                 DROP CONSTRAINT IF EXISTS surveys_venue_id_fk;
ALTER TABLE surveys                 DROP CONSTRAINT IF EXISTS surveys_created_by_fk;
ALTER TABLE questions               DROP CONSTRAINT IF EXISTS questions_survey_id_fk;
ALTER TABLE options                 DROP CONSTRAINT IF EXISTS options_question_id_fk;
ALTER TABLE survey_responses        DROP CONSTRAINT IF EXISTS survey_responses_survey_id_fk;
ALTER TABLE survey_answers          DROP CONSTRAINT IF EXISTS survey_answers_response_id_fk;
ALTER TABLE survey_answers          DROP CONSTRAINT IF EXISTS survey_answers_question_id_fk;
ALTER TABLE survey_answers          DROP CONSTRAINT IF EXISTS survey_answers_option_id_fk;
ALTER TABLE beacons                 DROP CONSTRAINT IF EXISTS beacons_venue_id_fk;
ALTER TABLE beacons                 DROP CONSTRAINT IF EXISTS beacons_level_id_fk;
ALTER TABLE connections             DROP CONSTRAINT IF EXISTS connections_venue_id_fk;
ALTER TABLE connection_levels       DROP CONSTRAINT IF EXISTS connection_levels_connection_id_fk;
ALTER TABLE connection_levels       DROP CONSTRAINT IF EXISTS connection_levels_level_id_fk;
ALTER TABLE advertisements          DROP CONSTRAINT IF EXISTS advertisements_venue_id_fk;
ALTER TABLE advertisements          DROP CONSTRAINT IF EXISTS advertisements_location_id_fk;
ALTER TABLE articles                DROP CONSTRAINT IF EXISTS articles_venue_id_fk;
ALTER TABLE articles                DROP CONSTRAINT IF EXISTS articles_location_id_fk;
ALTER TABLE article_images          DROP CONSTRAINT IF EXISTS article_images_article_id_fk;
ALTER TABLE coupons                 DROP CONSTRAINT IF EXISTS coupons_venue_id_fk;
ALTER TABLE snapshots               DROP CONSTRAINT IF EXISTS snapshots_venue_id_fk;
ALTER TABLE snapshots               DROP CONSTRAINT IF EXISTS snapshots_created_by_fk;

-- Also drop app_users venue_id FK (added in script 14 without reference, safe no-op)
ALTER TABLE app_users DROP CONSTRAINT IF EXISTS app_users_venue_id_fkey;

-- =============================================================================
-- SECTION 2: Convert PKs and FK columns (dependency order: parents first)
-- =============================================================================

-- ── customers ────────────────────────────────────────────────────────────────
CALL conv_pk('customers');

-- ── perspectives ─────────────────────────────────────────────────────────────
CALL conv_pk('perspectives');

-- ── event_tags ───────────────────────────────────────────────────────────────
CALL conv_pk('event_tags');

-- ── venues ───────────────────────────────────────────────────────────────────
CALL conv_pk('venues');
CALL conv_fk('venues', 'customer_id');

-- ── map_groups ───────────────────────────────────────────────────────────────
CALL conv_pk('map_groups');
CALL conv_fk('map_groups', 'venue_id');

-- ── event_types ──────────────────────────────────────────────────────────────
CALL conv_pk('event_types');
CALL conv_fk('event_types', 'venue_id');

-- ── location_categories ──────────────────────────────────────────────────────
CALL conv_pk('location_categories');
CALL conv_fk('location_categories', 'venue_id');

-- ── product_categories ───────────────────────────────────────────────────────
CALL conv_pk('product_categories');
CALL conv_fk('product_categories', 'venue_id');

-- ── levels ───────────────────────────────────────────────────────────────────
CALL conv_pk('levels');
CALL conv_fk('levels', 'venue_id');
CALL conv_fk('levels', 'map_group_id');
CALL conv_fk('levels', 'perspective_id');

-- ── geo_references ───────────────────────────────────────────────────────────
CALL conv_pk('geo_references');
CALL conv_fk('geo_references', 'level_id');

-- ── locations ────────────────────────────────────────────────────────────────
CALL conv_pk('locations');
CALL conv_fk('locations', 'venue_id');
CALL conv_fk('locations', 'level_id');
CALL conv_fk('locations', 'main_category_id');

-- ── location_category_links (id already dropped, composite PK on location_id+category_id)
CALL conv_fk('location_category_links', 'location_id');
CALL conv_fk('location_category_links', 'category_id');

-- ── location_images ──────────────────────────────────────────────────────────
CALL conv_pk('location_images');
CALL conv_fk('location_images', 'location_id');

-- ── products ─────────────────────────────────────────────────────────────────
CALL conv_pk('products');
CALL conv_fk('products', 'venue_id');
CALL conv_fk('products', 'location_id');
CALL conv_fk('products', 'main_category_id');

-- ── product_category_links (id already dropped, composite PK on product_id+category_id)
CALL conv_fk('product_category_links', 'product_id');
CALL conv_fk('product_category_links', 'category_id');

-- ── product_attachments ──────────────────────────────────────────────────────
CALL conv_pk('product_attachments');
CALL conv_fk('product_attachments', 'product_id');

-- ── events ───────────────────────────────────────────────────────────────────
CALL conv_pk('events');
CALL conv_fk('events', 'venue_id');
CALL conv_fk('events', 'type_id');

-- ── event_tag_links (id already dropped, composite PK on event_id+tag_id)
CALL conv_fk('event_tag_links', 'event_id');
CALL conv_fk('event_tag_links', 'tag_id');

-- ── event_location_links (id already dropped, composite PK on event_id+location_id)
CALL conv_fk('event_location_links', 'event_id');
CALL conv_fk('event_location_links', 'location_id');

-- ── event_images ─────────────────────────────────────────────────────────────
CALL conv_pk('event_images');
CALL conv_fk('event_images', 'event_id');

-- ── surveys (id is varchar(255) — same hex format, conv_pk handles it)
CALL conv_pk('surveys');
CALL conv_fk('surveys', 'venue_id');
-- surveys.created_by is already UUID (set to NULL in script 10, column type already UUID)

-- ── questions (survey_id is varchar(255) referencing surveys.id) ──────────────
CALL conv_pk('questions');
CALL conv_fk('questions', 'survey_id');

-- ── options ──────────────────────────────────────────────────────────────────
CALL conv_pk('options');
CALL conv_fk('options', 'question_id');

-- ── survey_responses (survey_id is varchar(255)) ──────────────────────────────
CALL conv_pk('survey_responses');
CALL conv_fk('survey_responses', 'survey_id');
-- user_id was dropped in script 10; no integer FK to recover

-- ── survey_answers ───────────────────────────────────────────────────────────
CALL conv_pk('survey_answers');
CALL conv_fk('survey_answers', 'response_id');
CALL conv_fk('survey_answers', 'question_id');
CALL conv_fk('survey_answers', 'option_id');

-- ── beacons ──────────────────────────────────────────────────────────────────
CALL conv_pk('beacons');
CALL conv_fk('beacons', 'venue_id');
CALL conv_fk('beacons', 'level_id');
CALL conv_fk('beacons', 'element_id');  -- char(32), not an FK reference

-- ── connections ──────────────────────────────────────────────────────────────
CALL conv_pk('connections');
CALL conv_fk('connections', 'venue_id');

-- ── connection_levels (bigserial id) ─────────────────────────────────────────
CALL conv_pk_serial('connection_levels');
CALL conv_fk('connection_levels', 'connection_id');
CALL conv_fk('connection_levels', 'level_id');
CALL conv_fk('connection_levels', 'element_id');  -- char(32), not an FK reference

-- ── articles ─────────────────────────────────────────────────────────────────
CALL conv_pk('articles');
CALL conv_fk('articles', 'venue_id');
CALL conv_fk('articles', 'location_id');
-- articles.product_id and product_category_id were dropped in script 11

-- ── article_images ───────────────────────────────────────────────────────────
CALL conv_pk('article_images');
CALL conv_fk('article_images', 'article_id');

-- ── advertisements ───────────────────────────────────────────────────────────
CALL conv_pk('advertisements');
CALL conv_fk('advertisements', 'venue_id');
CALL conv_fk('advertisements', 'location_id');
-- advertisements.article_id was dropped in script 11

-- ── coupons ──────────────────────────────────────────────────────────────────
CALL conv_pk('coupons');
CALL conv_fk('coupons', 'venue_id');
-- coupons.survey_id was dropped in script 11

-- ── snapshots ────────────────────────────────────────────────────────────────
CALL conv_pk('snapshots');
CALL conv_fk('snapshots', 'venue_id');
-- snapshots.created_by was already remapped via _user_id_map in script 12

-- =============================================================================
-- SECTION 3: venue_user_roles (bigserial id + profile_id INTEGER → user_id UUID)
-- =============================================================================

-- Convert bigserial id → UUID
CALL conv_pk_serial('venue_user_roles');

-- Convert venue_id char(32) → UUID
CALL conv_fk('venue_user_roles', 'venue_id');

-- profile_id (INTEGER = auth_user.id) cannot be mapped to UUID: _user_id_map was
-- already dropped in script 13. Rename to user_id UUID (NULL), then delete those rows
-- since user_id NOT NULL is required by the Go schema. Venue roles must be re-assigned.
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
    WHERE table_schema = current_schema() AND table_name = 'venue_user_roles'
    AND column_name = 'profile_id') THEN
    ALTER TABLE venue_user_roles ADD COLUMN user_id UUID;
    ALTER TABLE venue_user_roles DROP COLUMN profile_id;
  END IF;
END $$;
DELETE FROM venue_user_roles WHERE user_id IS NULL;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
    WHERE table_schema = current_schema() AND table_name = 'venue_user_roles'
    AND column_name = 'user_id') THEN
    ALTER TABLE venue_user_roles ALTER COLUMN user_id SET NOT NULL;
  END IF;
END $$;

-- =============================================================================
-- SECTION 4: venue_invitations (char(32) id + integer user FK columns)
-- =============================================================================

CALL conv_pk('venue_invitations');
CALL conv_fk('venue_invitations', 'venue_id');

-- invited_by_id and invited_profile_id are integer profile IDs — mapping is lost.
-- Drop integer columns, add invited_by UUID. Delete all rows (stale invitations).
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
    WHERE table_schema = current_schema() AND table_name = 'venue_invitations'
    AND column_name = 'invited_by_id') THEN
    ALTER TABLE venue_invitations DROP COLUMN IF EXISTS invited_by_id;
    ALTER TABLE venue_invitations DROP COLUMN IF EXISTS invited_profile_id;
    ALTER TABLE venue_invitations ADD COLUMN IF NOT EXISTS invited_by UUID;
  END IF;
END $$;
-- All invitation rows have NULL invited_by; delete them (invited_by NOT NULL in schema)
DELETE FROM venue_invitations WHERE invited_by IS NULL;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
    WHERE table_schema = current_schema() AND table_name = 'venue_invitations'
    AND column_name = 'invited_by') THEN
    ALTER TABLE venue_invitations ALTER COLUMN invited_by SET NOT NULL;
  END IF;
END $$;

-- =============================================================================
-- SECTION 5: Set uuidv7() default on all newly converted id columns
-- =============================================================================

DO $$
DECLARE tbl TEXT;
BEGIN
  FOR tbl IN
    SELECT DISTINCT table_name FROM information_schema.columns
    WHERE table_schema = current_schema()
      AND column_name = 'id'
      AND data_type = 'uuid'
      AND column_default IS DISTINCT FROM 'uuidv7()'
      AND table_name NOT IN ('schema_migrations')
  LOOP
    BEGIN
      EXECUTE format('ALTER TABLE %I ALTER COLUMN id SET DEFAULT uuidv7()', tbl);
    EXCEPTION WHEN OTHERS THEN NULL; END;
  END LOOP;
END $$;

-- =============================================================================
-- SECTION 6: Re-add FK constraints (now referencing UUID columns)
-- =============================================================================

CALL add_fk('venues_customer_id_fk',                    'venues',                  'customer_id',    'customers',         'id', 'SET NULL');
CALL add_fk('venue_user_roles_venue_id_fk',             'venue_user_roles',        'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('venue_user_roles_user_id_fk',              'venue_user_roles',        'user_id',        'users',             'id', 'CASCADE');
CALL add_fk('venue_invitations_venue_id_fk',            'venue_invitations',       'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('venue_invitations_invited_by_fk',          'venue_invitations',       'invited_by',     'users',             'id', 'CASCADE');
CALL add_fk('map_groups_venue_id_fk',                   'map_groups',              'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('levels_venue_id_fk',                       'levels',                  'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('levels_map_group_id_fk',                   'levels',                  'map_group_id',   'map_groups',        'id', 'SET NULL');
CALL add_fk('levels_perspective_id_fk',                 'levels',                  'perspective_id', 'perspectives',      'id', 'SET NULL');
CALL add_fk('geo_references_level_id_fk',               'geo_references',          'level_id',       'levels',            'id', 'CASCADE');
CALL add_fk('location_categories_venue_id_fk',          'location_categories',     'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('locations_venue_id_fk',                    'locations',               'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('locations_level_id_fk',                    'locations',               'level_id',       'levels',            'id', 'SET NULL');
CALL add_fk('locations_main_category_id_fk',            'locations',               'main_category_id','location_categories','id','SET NULL');
CALL add_fk('location_category_links_location_id_fk',   'location_category_links', 'location_id',    'locations',         'id', 'CASCADE');
CALL add_fk('location_category_links_category_id_fk',   'location_category_links', 'category_id',    'location_categories','id','CASCADE');
CALL add_fk('location_images_location_id_fk',           'location_images',         'location_id',    'locations',         'id', 'CASCADE');
CALL add_fk('product_categories_venue_id_fk',           'product_categories',      'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('products_venue_id_fk',                     'products',                'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('products_location_id_fk',                  'products',                'location_id',    'locations',         'id', 'SET NULL');
CALL add_fk('products_main_category_id_fk',             'products',                'main_category_id','product_categories','id','SET NULL');
CALL add_fk('product_category_links_product_id_fk',     'product_category_links',  'product_id',     'products',          'id', 'CASCADE');
CALL add_fk('product_category_links_category_id_fk',    'product_category_links',  'category_id',    'product_categories','id','CASCADE');
CALL add_fk('product_attachments_product_id_fk',        'product_attachments',     'product_id',     'products',          'id', 'CASCADE');
CALL add_fk('event_types_venue_id_fk',                  'event_types',             'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('events_venue_id_fk',                       'events',                  'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('events_type_id_fk',                        'events',                  'type_id',        'event_types',       'id', 'SET NULL');
CALL add_fk('event_tag_links_event_id_fk',              'event_tag_links',         'event_id',       'events',            'id', 'CASCADE');
CALL add_fk('event_tag_links_tag_id_fk',                'event_tag_links',         'tag_id',         'event_tags',        'id', 'CASCADE');
CALL add_fk('event_location_links_event_id_fk',         'event_location_links',    'event_id',       'events',            'id', 'CASCADE');
CALL add_fk('event_location_links_location_id_fk',      'event_location_links',    'location_id',    'locations',         'id', 'CASCADE');
CALL add_fk('event_images_event_id_fk',                 'event_images',            'event_id',       'events',            'id', 'CASCADE');
CALL add_fk('surveys_venue_id_fk',                      'surveys',                 'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('surveys_created_by_fk',                    'surveys',                 'created_by',     'users',             'id', 'SET NULL');
CALL add_fk('questions_survey_id_fk',                   'questions',               'survey_id',      'surveys',           'id', 'CASCADE');
CALL add_fk('options_question_id_fk',                   'options',                 'question_id',    'questions',         'id', 'CASCADE');
CALL add_fk('survey_responses_survey_id_fk',            'survey_responses',        'survey_id',      'surveys',           'id', 'CASCADE');
CALL add_fk('survey_answers_response_id_fk',            'survey_answers',          'response_id',    'survey_responses',  'id', 'CASCADE');
CALL add_fk('survey_answers_question_id_fk',            'survey_answers',          'question_id',    'questions',         'id', 'CASCADE');
CALL add_fk('survey_answers_option_id_fk',              'survey_answers',          'option_id',      'options',           'id', 'SET NULL');
CALL add_fk('beacons_venue_id_fk',                      'beacons',                 'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('beacons_level_id_fk',                      'beacons',                 'level_id',       'levels',            'id', 'SET NULL');
CALL add_fk('connections_venue_id_fk',                  'connections',             'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('connection_levels_connection_id_fk',       'connection_levels',       'connection_id',  'connections',       'id', 'CASCADE');
CALL add_fk('connection_levels_level_id_fk',            'connection_levels',       'level_id',       'levels',            'id', 'CASCADE');
CALL add_fk('articles_venue_id_fk',                     'articles',                'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('articles_location_id_fk',                  'articles',                'location_id',    'locations',         'id', 'SET NULL');
CALL add_fk('article_images_article_id_fk',             'article_images',          'article_id',     'articles',          'id', 'CASCADE');
CALL add_fk('advertisements_venue_id_fk',               'advertisements',          'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('advertisements_location_id_fk',            'advertisements',          'location_id',    'locations',         'id', 'SET NULL');
CALL add_fk('coupons_venue_id_fk',                      'coupons',                 'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('snapshots_venue_id_fk',                    'snapshots',               'venue_id',       'venues',            'id', 'CASCADE');
CALL add_fk('snapshots_created_by_fk',                  'snapshots',               'created_by',     'users',             'id', 'SET NULL');

-- Wire app_users.venue_id → venues.id now that venues.id is UUID
CALL add_fk('app_users_venue_id_fk', 'app_users', 'venue_id', 'venues', 'id', 'SET NULL');

-- =============================================================================
-- SECTION 7: Cleanup helpers
-- =============================================================================

DROP PROCEDURE IF EXISTS conv_pk(TEXT);
DROP PROCEDURE IF EXISTS conv_pk_serial(TEXT);
DROP PROCEDURE IF EXISTS conv_fk(TEXT, TEXT);
DROP PROCEDURE IF EXISTS add_fk(TEXT, TEXT, TEXT, TEXT, TEXT, TEXT);
DROP FUNCTION  IF EXISTS hex_to_uuid(TEXT);
