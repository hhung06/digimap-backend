# Database Migration Plan: Normalize pgloader Schema to Go Backend Schema

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Transform the pgloader-migrated PostgreSQL database (with `indoormap_api_*` tables) into the exact schema the Go backend expects, merge `appuser` + `visitor` into a new `app_users` table, and register all 24 existing Go migrations as applied so `make migrate-up` only runs future migrations.

**Architecture:** Pure SQL transformation scripts in `resources/migration_scripts/`, executed in strict order. A `_user_id_map` temp table bridges `auth_user` integer PKs to new UUIDs. All FK constraints are dropped upfront (Task 2), transformations run freely (Tasks 3–13), then constraints and triggers are re-installed (Task 14). Migration `000025` introduces `app_users` as the first net-new table.

**Tech Stack:** PostgreSQL 14+, psql, golang-migrate

---

## File Map

| File | Purpose |
|------|---------|
| `resources/migration_scripts/01_prep.sql` | Install extensions/functions from migration 000001, drop Django framework tables |
| `resources/migration_scripts/02_drop_fk_constraints.sql` | Drop every FK constraint so renames/type changes can proceed freely |
| `resources/migration_scripts/03_rename_tables.sql` | Rename all `indoormap_api_*` tables to final names in one pass |
| `resources/migration_scripts/04_core_tables.sql` | Transform customers, auth_user→users (with UUID mapping), venue_user_roles, venue_invitations |
| `resources/migration_scripts/05_venues.sql` | Transform venues (heavy camelCase renames + new columns) |
| `resources/migration_scripts/06_levels.sql` | Transform map_groups, perspectives, levels, geo_references |
| `resources/migration_scripts/07_locations.sql` | Transform location_categories, locations, location_category_links, location_images, promotions |
| `resources/migration_scripts/08_products.sql` | Transform product_categories, products, product_category_links, product_attachments |
| `resources/migration_scripts/09_events.sql` | Transform event_types, event_tags, events, event_tag_links, event_location_links, event_images |
| `resources/migration_scripts/10_surveys.sql` | Transform surveys, questions, options, survey_responses, survey_answers |
| `resources/migration_scripts/11_content.sql` | Transform advertisements, articles, article_images, coupons, beacons, connections, connection_levels, qrcodes |
| `resources/migration_scripts/12_misc.sql` | Transform tags, search_queries, snapshots |
| `resources/migration_scripts/13_restore_constraints.sql` | Re-add all FK constraints, install updated_at triggers, set uuidv7() defaults |
| `resources/migration_scripts/14_drop_leftover_tables.sql` | Drop all remaining `indoormap_api_*` tables with no new schema equivalent |
| `resources/migration_scripts/15_mark_migrations_applied.sql` | Populate schema_migrations so golang-migrate skips 000001–000024 |
| `migrations/000025_create_app_users.up.sql` | Create app_users table (merged appuser + visitor) |
| `migrations/000025_create_app_users.down.sql` | Drop app_users table |

---

### Task 1: Prerequisites — extensions, functions, drop framework tables

**Files:**
- Create: `resources/migration_scripts/01_prep.sql`

- [ ] **Step 1: Create `01_prep.sql`**

```sql
-- ── Extensions & functions (from migration 000001) ────────────────────────
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN NEW.updated_at = NOW(); RETURN NEW; END;
$$;

CREATE OR REPLACE FUNCTION create_updated_at_trigger(tbl TEXT)
RETURNS VOID LANGUAGE plpgsql AS $$
BEGIN
  EXECUTE format(
    'CREATE TRIGGER set_updated_at BEFORE UPDATE ON %I
     FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()', tbl);
END;
$$;

-- ── Drop Django framework tables ──────────────────────────────────────────
DROP TABLE IF EXISTS auth_group_permissions CASCADE;
DROP TABLE IF EXISTS auth_user_groups CASCADE;
DROP TABLE IF EXISTS auth_user_user_permissions CASCADE;
DROP TABLE IF EXISTS authtoken_token CASCADE;
DROP TABLE IF EXISTS auth_group CASCADE;
DROP TABLE IF EXISTS auth_permission CASCADE;
DROP TABLE IF EXISTS debug_toolbar_historyentry CASCADE;
DROP TABLE IF EXISTS django_admin_log CASCADE;
DROP TABLE IF EXISTS django_content_type CASCADE;
DROP TABLE IF EXISTS django_migrations CASCADE;
DROP TABLE IF EXISTS django_session CASCADE;
DROP TABLE IF EXISTS rest_framework_api_key_apikey CASCADE;
DROP TABLE IF EXISTS token_blacklist_blacklistedtoken CASCADE;
DROP TABLE IF EXISTS token_blacklist_outstandingtoken CASCADE;

-- ── Drop tables with no new schema equivalent ─────────────────────────────
DROP TABLE IF EXISTS indoormap_api_accesshistory CASCADE;
DROP TABLE IF EXISTS indoormap_api_adtrack CASCADE;
DROP TABLE IF EXISTS indoormap_api_articlerelatedproducts CASCADE;
DROP TABLE IF EXISTS indoormap_api_asset CASCADE;
DROP TABLE IF EXISTS indoormap_api_libraryasset CASCADE;
DROP TABLE IF EXISTS indoormap_api_categoryad CASCADE;
DROP TABLE IF EXISTS indoormap_api_coupon_gifts CASCADE;
DROP TABLE IF EXISTS indoormap_api_promogift CASCADE;
DROP TABLE IF EXISTS indoormap_api_forcesyncversion CASCADE;
DROP TABLE IF EXISTS indoormap_api_forceversion CASCADE;
DROP TABLE IF EXISTS indoormap_api_language CASCADE;
DROP TABLE IF EXISTS indoormap_api_leveltype CASCADE;
DROP TABLE IF EXISTS indoormap_api_location_place_tags CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationcategorytranslation CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtranslation CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtemplate CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtemplate_common_categories CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtemplate_place_tags CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtemplateimage CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtemplatetranslation CASCADE;
DROP TABLE IF EXISTS indoormap_api_new CASCADE;
DROP TABLE IF EXISTS indoormap_api_plugin CASCADE;
DROP TABLE IF EXISTS indoormap_api_productplaza CASCADE;
DROP TABLE IF EXISTS indoormap_api_productplazadescriptionimage CASCADE;
DROP TABLE IF EXISTS indoormap_api_productplazaimage CASCADE;
DROP TABLE IF EXISTS indoormap_api_productplazapendingvideo CASCADE;
DROP TABLE IF EXISTS indoormap_api_profile CASCADE;
DROP TABLE IF EXISTS indoormap_api_segment CASCADE;
DROP TABLE IF EXISTS indoormap_api_syncdata CASCADE;
DROP TABLE IF EXISTS indoormap_api_venue_sync_flag CASCADE;
DROP TABLE IF EXISTS indoormap_api_venuetheme CASCADE;
DROP TABLE IF EXISTS indoormap_api_venuetranslation CASCADE;
DROP TABLE IF EXISTS indoormap_api_tagtranslation CASCADE;
DROP TABLE IF EXISTS indoormap_api_usercoupon CASCADE;
DROP TABLE IF EXISTS indoormap_api_visitor CASCADE;   -- migrated in Task 14
DROP TABLE IF EXISTS indoormap_api_appuser CASCADE;   -- migrated in Task 14
DROP TABLE IF EXISTS indoormap_api_templatecategory CASCADE;
DROP TABLE IF EXISTS indoormap_api_templatecategorytranslation CASCADE;
```

- [ ] **Step 2: Run and verify**

```bash
psql $DATABASE_URL -f resources/migration_scripts/01_prep.sql
# Expected: no ERRORs, a series of DROP TABLE and CREATE statements
psql $DATABASE_URL -c "\dt indoormap_api_access*"
# Expected: no rows (table dropped)
```

- [ ] **Step 3: Commit**

```bash
git add resources/migration_scripts/01_prep.sql
git commit -m "chore: add db migration script 01 - prep and framework table cleanup"
```

---

### Task 2: Drop all FK constraints

**Files:**
- Create: `resources/migration_scripts/02_drop_fk_constraints.sql`

This lets Tasks 3–13 freely rename tables and change column types without constraint violations.

- [ ] **Step 1: Generate the drop script by querying pg_constraint, then create the file**

Run this query to see all FK constraints that need dropping:
```sql
SELECT 'ALTER TABLE '||quote_ident(tc.table_name)||
       ' DROP CONSTRAINT '||quote_ident(tc.constraint_name)||';'
FROM information_schema.table_constraints tc
WHERE tc.constraint_type = 'FOREIGN KEY'
  AND tc.table_schema = 'public'
ORDER BY tc.table_name;
```

- [ ] **Step 2: Create `02_drop_fk_constraints.sql`** with explicit drops for every known FK:

```sql
-- customers
ALTER TABLE indoormap_api_venue DROP CONSTRAINT IF EXISTS indoormap_api_venue_customer_id_fk;

-- venue-scoped tables (venue_id FKs)
ALTER TABLE indoormap_api_venueuserrole DROP CONSTRAINT IF EXISTS indoormap_api_venueuserrole_venue_id_fk;
ALTER TABLE indoormap_api_venueinvitation DROP CONSTRAINT IF EXISTS indoormap_api_venueinvitation_venue_id_fk;
ALTER TABLE indoormap_api_venueinvitation DROP CONSTRAINT IF EXISTS indoormap_api_venueinvitation_invited_by_id_fk;
ALTER TABLE indoormap_api_mapgroup DROP CONSTRAINT IF EXISTS indoormap_api_mapgroup_venue_id_fk;
ALTER TABLE indoormap_api_level DROP CONSTRAINT IF EXISTS indoormap_api_level_venue_id_fk;
ALTER TABLE indoormap_api_level DROP CONSTRAINT IF EXISTS indoormap_api_level_mapgroup_id_fk;
ALTER TABLE indoormap_api_level DROP CONSTRAINT IF EXISTS indoormap_api_level_perspective_id_fk;
ALTER TABLE indoormap_api_georeference DROP CONSTRAINT IF EXISTS indoormap_api_georeference_level_id_fk;
ALTER TABLE indoormap_api_locationcategory DROP CONSTRAINT IF EXISTS indoormap_api_locationcategory_venue_id_fk;
ALTER TABLE indoormap_api_location DROP CONSTRAINT IF EXISTS indoormap_api_location_venue_id_fk;
ALTER TABLE indoormap_api_location DROP CONSTRAINT IF EXISTS indoormap_api_location_level_id_fk;
ALTER TABLE indoormap_api_location DROP CONSTRAINT IF EXISTS indoormap_api_location_main_category_id_fk;
ALTER TABLE indoormap_api_location_common_categories DROP CONSTRAINT IF EXISTS indoormap_api_location_common_categories_location_id_fk;
ALTER TABLE indoormap_api_location_common_categories DROP CONSTRAINT IF EXISTS indoormap_api_location_common_categories_category_id_fk;
ALTER TABLE indoormap_api_locationimage DROP CONSTRAINT IF EXISTS indoormap_api_locationimage_location_id_fk;
ALTER TABLE indoormap_api_productcategory DROP CONSTRAINT IF EXISTS indoormap_api_productcategory_venue_id_fk;
ALTER TABLE indoormap_api_product DROP CONSTRAINT IF EXISTS indoormap_api_product_venue_id_fk;
ALTER TABLE indoormap_api_product DROP CONSTRAINT IF EXISTS indoormap_api_product_location_id_fk;
ALTER TABLE indoormap_api_product DROP CONSTRAINT IF EXISTS indoormap_api_product_main_category_id_fk;
ALTER TABLE indoormap_api_product_categories DROP CONSTRAINT IF EXISTS indoormap_api_product_categories_product_id_fk;
ALTER TABLE indoormap_api_product_categories DROP CONSTRAINT IF EXISTS indoormap_api_product_categories_category_id_fk;
ALTER TABLE indoormap_api_productattachment DROP CONSTRAINT IF EXISTS indoormap_api_productattachment_product_id_fk;
ALTER TABLE indoormap_api_eventtype DROP CONSTRAINT IF EXISTS indoormap_api_eventtype_venue_id_fk;
ALTER TABLE indoormap_api_event DROP CONSTRAINT IF EXISTS indoormap_api_event_venue_id_fk;
ALTER TABLE indoormap_api_event DROP CONSTRAINT IF EXISTS indoormap_api_event_type_id_fk;
ALTER TABLE indoormap_api_event_tags DROP CONSTRAINT IF EXISTS indoormap_api_event_tags_event_id_fk;
ALTER TABLE indoormap_api_event_tags DROP CONSTRAINT IF EXISTS indoormap_api_event_tags_tag_id_fk;
ALTER TABLE indoormap_api_event_locations DROP CONSTRAINT IF EXISTS indoormap_api_event_locations_event_id_fk;
ALTER TABLE indoormap_api_event_locations DROP CONSTRAINT IF EXISTS indoormap_api_event_locations_location_id_fk;
ALTER TABLE indoormap_api_eventimage DROP CONSTRAINT IF EXISTS indoormap_api_eventimage_event_id_fk;
ALTER TABLE indoormap_api_survey DROP CONSTRAINT IF EXISTS indoormap_api_survey_venue_id_fk;
ALTER TABLE indoormap_api_question DROP CONSTRAINT IF EXISTS indoormap_api_question_survey_id_fk;
ALTER TABLE indoormap_api_option DROP CONSTRAINT IF EXISTS indoormap_api_option_question_id_fk;
ALTER TABLE indoormap_api_surveyresponse DROP CONSTRAINT IF EXISTS indoormap_api_surveyresponse_survey_id_fk;
ALTER TABLE indoormap_api_surveyanswer DROP CONSTRAINT IF EXISTS indoormap_api_surveyanswer_response_id_fk;
ALTER TABLE indoormap_api_surveyanswer DROP CONSTRAINT IF EXISTS indoormap_api_surveyanswer_question_id_fk;
ALTER TABLE indoormap_api_surveyanswer DROP CONSTRAINT IF EXISTS indoormap_api_surveyanswer_option_id_fk;
ALTER TABLE indoormap_api_advertisement DROP CONSTRAINT IF EXISTS indoormap_api_advertisement_venue_id_fk;
ALTER TABLE indoormap_api_advertisement DROP CONSTRAINT IF EXISTS indoormap_api_advertisement_location_id_fk;
ALTER TABLE indoormap_api_article DROP CONSTRAINT IF EXISTS indoormap_api_article_venue_id_fk;
ALTER TABLE indoormap_api_article DROP CONSTRAINT IF EXISTS indoormap_api_article_location_id_fk;
ALTER TABLE indoormap_api_articleimage DROP CONSTRAINT IF EXISTS indoormap_api_articleimage_article_id_fk;
ALTER TABLE indoormap_api_coupon DROP CONSTRAINT IF EXISTS indoormap_api_coupon_venue_id_fk;
ALTER TABLE indoormap_api_beacon DROP CONSTRAINT IF EXISTS indoormap_api_beacon_venue_id_fk;
ALTER TABLE indoormap_api_beacon DROP CONSTRAINT IF EXISTS indoormap_api_beacon_level_id_fk;
ALTER TABLE indoormap_api_connection DROP CONSTRAINT IF EXISTS indoormap_api_connection_venue_id_fk;
ALTER TABLE indoormap_api_connectionlevel DROP CONSTRAINT IF EXISTS indoormap_api_connectionlevel_connection_id_fk;
ALTER TABLE indoormap_api_connectionlevel DROP CONSTRAINT IF EXISTS indoormap_api_connectionlevel_level_id_fk;
ALTER TABLE indoormap_api_qrcode DROP CONSTRAINT IF EXISTS indoormap_api_qrcode_venue_id_fk;
ALTER TABLE indoormap_api_qrcode DROP CONSTRAINT IF EXISTS indoormap_api_qrcode_level_id_fk;
ALTER TABLE indoormap_api_qrcode DROP CONSTRAINT IF EXISTS indoormap_api_qrcode_location_id_fk;
ALTER TABLE indoormap_api_tag DROP CONSTRAINT IF EXISTS indoormap_api_tag_venue_id_fk;
ALTER TABLE indoormap_api_searchquery DROP CONSTRAINT IF EXISTS indoormap_api_searchquery_venue_id_fk;
ALTER TABLE indoormap_api_snapshot DROP CONSTRAINT IF EXISTS indoormap_api_snapshot_venue_id_fk;
ALTER TABLE indoormap_api_snapshot DROP CONSTRAINT IF EXISTS indoormap_api_snapshot_publish_by_id_fk;

-- Use IF EXISTS + CASCADE on all; actual names may differ — run the query above first
-- to get exact names and update this file before running.
```

- [ ] **Step 3: Run**

```bash
psql $DATABASE_URL -f resources/migration_scripts/02_drop_fk_constraints.sql
# Expected: ALTER TABLE for each line, no ERRORs
```

- [ ] **Step 4: Verify no FKs remain**

```sql
SELECT COUNT(*) FROM information_schema.table_constraints
WHERE constraint_type = 'FOREIGN KEY' AND table_schema = 'public';
-- Expected: 0
```

- [ ] **Step 5: Commit**

```bash
git add resources/migration_scripts/02_drop_fk_constraints.sql
git commit -m "chore: add db migration script 02 - drop all FK constraints"
```

---

### Task 3: Rename all tables in one pass

**Files:**
- Create: `resources/migration_scripts/03_rename_tables.sql`

- [ ] **Step 1: Create `03_rename_tables.sql`**

```sql
ALTER TABLE indoormap_api_customer               RENAME TO customers;
ALTER TABLE auth_user                            RENAME TO users;
ALTER TABLE indoormap_api_venue                  RENAME TO venues;
ALTER TABLE indoormap_api_venueuserrole          RENAME TO venue_user_roles;
ALTER TABLE indoormap_api_venueinvitation        RENAME TO venue_invitations;
ALTER TABLE indoormap_api_mapgroup               RENAME TO map_groups;
ALTER TABLE indoormap_api_perspective            RENAME TO perspectives;
ALTER TABLE indoormap_api_level                  RENAME TO levels;
ALTER TABLE indoormap_api_georeference           RENAME TO geo_references;
ALTER TABLE indoormap_api_locationcategory       RENAME TO location_categories;
ALTER TABLE indoormap_api_location               RENAME TO locations;
ALTER TABLE indoormap_api_location_common_categories RENAME TO location_category_links;
ALTER TABLE indoormap_api_locationimage          RENAME TO location_images;
ALTER TABLE indoormap_api_locationpromotion      RENAME TO promotions;
ALTER TABLE indoormap_api_productcategory        RENAME TO product_categories;
ALTER TABLE indoormap_api_product                RENAME TO products;
ALTER TABLE indoormap_api_product_categories     RENAME TO product_category_links;
ALTER TABLE indoormap_api_productattachment      RENAME TO product_attachments;
ALTER TABLE indoormap_api_eventtype              RENAME TO event_types;
ALTER TABLE indoormap_api_eventtag               RENAME TO event_tags;
ALTER TABLE indoormap_api_event                  RENAME TO events;
ALTER TABLE indoormap_api_event_tags             RENAME TO event_tag_links;
ALTER TABLE indoormap_api_event_locations        RENAME TO event_location_links;
ALTER TABLE indoormap_api_eventimage             RENAME TO event_images;
ALTER TABLE indoormap_api_survey                 RENAME TO surveys;
ALTER TABLE indoormap_api_question               RENAME TO questions;
ALTER TABLE indoormap_api_option                 RENAME TO options;
ALTER TABLE indoormap_api_surveyresponse         RENAME TO survey_responses;
ALTER TABLE indoormap_api_surveyanswer           RENAME TO survey_answers;
ALTER TABLE indoormap_api_beacon                 RENAME TO beacons;
ALTER TABLE indoormap_api_connection             RENAME TO connections;
ALTER TABLE indoormap_api_connectionlevel        RENAME TO connection_levels;
ALTER TABLE indoormap_api_advertisement          RENAME TO advertisements;
ALTER TABLE indoormap_api_article                RENAME TO articles;
ALTER TABLE indoormap_api_articleimage           RENAME TO article_images;
ALTER TABLE indoormap_api_coupon                 RENAME TO coupons;
ALTER TABLE indoormap_api_tag                    RENAME TO tags;
ALTER TABLE indoormap_api_searchquery            RENAME TO search_queries;
ALTER TABLE indoormap_api_snapshot               RENAME TO snapshots;
ALTER TABLE indoormap_api_qrcode                 RENAME TO qrcodes;
ALTER TABLE indoormap_api_resetpasswordtoken     RENAME TO reset_password_tokens;
```

- [ ] **Step 2: Run and verify**

```bash
psql $DATABASE_URL -f resources/migration_scripts/03_rename_tables.sql
psql $DATABASE_URL -c "\dt" | grep -c "indoormap_api_"
# Expected: 0 (no indoormap_api_ tables remaining except appuser/visitor handled in Task 14)
```

- [ ] **Step 3: Commit**

```bash
git add resources/migration_scripts/03_rename_tables.sql
git commit -m "chore: add db migration script 03 - rename all tables"
```

---

### Task 4: Transform core tables — customers, users, venue_user_roles, venue_invitations

**Files:**
- Create: `resources/migration_scripts/04_core_tables.sql`

- [ ] **Step 1: Create `04_core_tables.sql`**

```sql
-- ══════════════════════════════════════════════════════
-- customers
-- ══════════════════════════════════════════════════════
ALTER TABLE customers RENAME COLUMN image TO logo_url;
ALTER TABLE customers DROP COLUMN IF EXISTS url;
ALTER TABLE customers DROP COLUMN IF EXISTS description;
ALTER TABLE customers DROP COLUMN IF EXISTS restored_at;
ALTER TABLE customers DROP COLUMN IF EXISTS transaction_id;
ALTER TABLE customers DROP COLUMN IF EXISTS address;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS metadata JSONB;

-- ══════════════════════════════════════════════════════
-- users  (was auth_user)
-- ══════════════════════════════════════════════════════

-- Build UUID mapping: old integer id → new UUID
-- (needed so snapshot.created_by and survey.created_by can reference new UUIDs)
CREATE TABLE IF NOT EXISTS _user_id_map (
    old_id BIGINT PRIMARY KEY,
    new_id UUID NOT NULL DEFAULT gen_random_uuid()
);
INSERT INTO _user_id_map (old_id)
SELECT id FROM users
ON CONFLICT DO NOTHING;

-- Add new UUID column alongside old integer id
ALTER TABLE users ADD COLUMN IF NOT EXISTS new_id UUID;
UPDATE users u SET new_id = m.new_id FROM _user_id_map m WHERE m.old_id = u.id;

-- Rename / add / drop columns
ALTER TABLE users RENAME COLUMN password    TO password_hash;
ALTER TABLE users RENAME COLUMN last_login  TO last_login_at;
ALTER TABLE users RENAME COLUMN is_superuser TO is_system_admin;
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url  TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone       TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at  TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE users DROP COLUMN IF EXISTS is_staff;
ALTER TABLE users DROP COLUMN IF EXISTS username;
ALTER TABLE users DROP COLUMN IF EXISTS date_joined;

-- Swap integer PK for UUID
ALTER TABLE users DROP CONSTRAINT IF EXISTS auth_user_pkey;
ALTER TABLE users DROP COLUMN IF EXISTS id;
ALTER TABLE users RENAME COLUMN new_id TO id;
ALTER TABLE users ADD PRIMARY KEY (id);
ALTER TABLE users ALTER COLUMN email TYPE TEXT;
ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);

-- ══════════════════════════════════════════════════════
-- venue_user_roles  (was indoormap_api_venueuserrole)
-- ══════════════════════════════════════════════════════
-- Create the enum type expected by migration 000005
DO $$ BEGIN
  CREATE TYPE venue_role AS ENUM ('owner', 'editor', 'viewer');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

ALTER TABLE venue_user_roles DROP COLUMN IF EXISTS restored_at;
ALTER TABLE venue_user_roles DROP COLUMN IF EXISTS transaction_id;
-- role column: smallint in old schema → venue_role enum
-- Map: assume 1=viewer,2=editor,3=owner (adjust if different)
ALTER TABLE venue_user_roles ADD COLUMN IF NOT EXISTS role_new venue_role;
UPDATE venue_user_roles SET role_new = CASE role::int
  WHEN 1 THEN 'viewer'::venue_role
  WHEN 2 THEN 'editor'::venue_role
  WHEN 3 THEN 'owner'::venue_role
  ELSE 'viewer'::venue_role END;
ALTER TABLE venue_user_roles DROP COLUMN role;
ALTER TABLE venue_user_roles RENAME COLUMN role_new TO role;
ALTER TABLE venue_user_roles ALTER COLUMN role SET NOT NULL;

-- ══════════════════════════════════════════════════════
-- venue_invitations  (was indoormap_api_venueinvitation)
-- ══════════════════════════════════════════════════════
DO $$ BEGIN
  CREATE TYPE invitation_status AS ENUM ('pending', 'accepted', 'cancelled');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

ALTER TABLE venue_invitations DROP COLUMN IF EXISTS restored_at;
ALTER TABLE venue_invitations DROP COLUMN IF EXISTS transaction_id;
-- status column: smallint → invitation_status enum
ALTER TABLE venue_invitations ADD COLUMN IF NOT EXISTS status_new invitation_status;
UPDATE venue_invitations SET status_new = CASE status::int
  WHEN 0 THEN 'pending'::invitation_status
  WHEN 1 THEN 'accepted'::invitation_status
  WHEN 2 THEN 'cancelled'::invitation_status
  ELSE 'pending'::invitation_status END;
ALTER TABLE venue_invitations DROP COLUMN IF EXISTS status;
ALTER TABLE venue_invitations RENAME COLUMN status_new TO status;
ALTER TABLE venue_invitations ALTER COLUMN status SET NOT NULL DEFAULT 'pending';
-- Map invited_by_id (old integer) to new user UUID
ALTER TABLE venue_invitations ADD COLUMN IF NOT EXISTS invited_by_new UUID;
UPDATE venue_invitations vi
SET invited_by_new = m.new_id
FROM _user_id_map m
WHERE vi.invited_by_id = m.old_id;
ALTER TABLE venue_invitations DROP COLUMN IF EXISTS invited_by_id;
ALTER TABLE venue_invitations RENAME COLUMN invited_by_new TO invited_by;
```

- [ ] **Step 2: Run and verify**

```bash
psql $DATABASE_URL -f resources/migration_scripts/04_core_tables.sql

psql $DATABASE_URL -c "SELECT COUNT(*) FROM users WHERE id IS NULL;"
# Expected: 0

psql $DATABASE_URL -c "SELECT COUNT(*) FROM _user_id_map;"
# Expected: equals SELECT COUNT(*) FROM users
```

- [ ] **Step 3: Commit**

```bash
git add resources/migration_scripts/04_core_tables.sql
git commit -m "chore: add db migration script 04 - transform core tables"
```

---

### Task 5: Transform venues

**Files:**
- Create: `resources/migration_scripts/05_venues.sql`

- [ ] **Step 1: Create `05_venues.sql`**

```sql
-- Column renames (camelCase → snake_case)
ALTER TABLE venues RENAME COLUMN largelogo      TO large_logo;
ALTER TABLE venues RENAME COLUMN mediumlogo     TO medium_logo;
ALTER TABLE venues RENAME COLUMN originallogo   TO original_logo;
ALTER TABLE venues RENAME COLUMN smalllogo      TO small_logo;
ALTER TABLE venues RENAME COLUMN publish        TO is_published;
ALTER TABLE venues RENAME COLUMN externalid     TO external_id;
ALTER TABLE venues RENAME COLUMN workhours      TO work_hours;
ALTER TABLE venues RENAME COLUMN seodescription TO seo_description;
ALTER TABLE venues RENAME COLUMN seokeywords    TO seo_keywords;
ALTER TABLE venues RENAME COLUMN seotitle       TO seo_title;
ALTER TABLE venues RENAME COLUMN subdomains     TO sub_domains;
ALTER TABLE venues RENAME COLUMN appconfigs     TO app_configs;
ALTER TABLE venues RENAME COLUMN appdomains     TO app_domains;
ALTER TABLE venues RENAME COLUMN bodytag        TO body_tag;
ALTER TABLE venues RENAME COLUMN headtag        TO head_tag;
ALTER TABLE venues RENAME COLUMN endat          TO end_at;
ALTER TABLE venues RENAME COLUMN startat        TO start_at;
ALTER TABLE venues RENAME COLUMN latitude       TO lat;
ALTER TABLE venues RENAME COLUMN longitude      TO lng;

-- Drop columns not in new schema
ALTER TABLE venues DROP COLUMN IF EXISTS countrycode;
ALTER TABLE venues DROP COLUMN IF EXISTS defaultmap;
ALTER TABLE venues DROP COLUMN IF EXISTS restored_at;
ALTER TABLE venues DROP COLUMN IF EXISTS transaction_id;

-- Add new columns
ALTER TABLE venues ADD COLUMN IF NOT EXISTS slug        TEXT;
ALTER TABLE venues ADD COLUMN IF NOT EXISTS type        SMALLINT NOT NULL DEFAULT 0;
ALTER TABLE venues ADD COLUMN IF NOT EXISTS theme       JSONB;
ALTER TABLE venues ADD COLUMN IF NOT EXISTS plugins     JSONB;
ALTER TABLE venues ADD COLUMN IF NOT EXISTS translations JSONB;

-- Generate slug from name (lowercase, replace spaces with hyphens, append short id suffix)
UPDATE venues SET slug = lower(regexp_replace(name, '[^a-zA-Z0-9]+', '-', 'g'))
    || '-' || substring(id::text, 1, 8)
WHERE slug IS NULL;
ALTER TABLE venues ALTER COLUMN slug SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS venues_slug_unique ON venues(slug) WHERE deleted_at IS NULL;

-- Ensure timezone has a default
ALTER TABLE venues ALTER COLUMN timezone SET DEFAULT 'UTC';
UPDATE venues SET timezone = 'UTC' WHERE timezone IS NULL OR timezone = '';

-- Type precision: double precision → numeric(10,7) for lat/lng
ALTER TABLE venues ALTER COLUMN lat  TYPE NUMERIC(10,7);
ALTER TABLE venues ALTER COLUMN lng  TYPE NUMERIC(10,7);
```

- [ ] **Step 2: Run and verify**

```bash
psql $DATABASE_URL -f resources/migration_scripts/05_venues.sql

psql $DATABASE_URL -c "SELECT COUNT(*) FROM venues WHERE slug IS NULL;"
# Expected: 0

psql $DATABASE_URL -c "\d venues" | grep -E "slug|large_logo|is_published"
# Expected: all three columns present
```

- [ ] **Step 3: Commit**

```bash
git add resources/migration_scripts/05_venues.sql
git commit -m "chore: add db migration script 05 - transform venues"
```

---

### Task 6: Transform level domain — map_groups, perspectives, levels, geo_references

**Files:**
- Create: `resources/migration_scripts/06_levels.sql`

- [ ] **Step 1: Create `06_levels.sql`**

```sql
-- ── map_groups ────────────────────────────────────────
ALTER TABLE map_groups RENAME COLUMN sortindex  TO sort_index;
ALTER TABLE map_groups RENAME COLUMN shortname  TO short_name;
ALTER TABLE map_groups DROP COLUMN IF EXISTS restored_at;
ALTER TABLE map_groups DROP COLUMN IF EXISTS transaction_id;
ALTER TABLE map_groups ALTER COLUMN sort_index TYPE INTEGER USING sort_index::integer;

-- ── perspectives ──────────────────────────────────────
ALTER TABLE perspectives RENAME COLUMN camerazoom            TO camera_zoom;
ALTER TABLE perspectives RENAME COLUMN cameratype            TO camera_type;
ALTER TABLE perspectives RENAME COLUMN cameramaxzoom         TO camera_max_zoom;
ALTER TABLE perspectives RENAME COLUMN cameraminzoom         TO camera_min_zoom;
ALTER TABLE perspectives RENAME COLUMN cameratargetbearing   TO camera_target_bearing;
ALTER TABLE perspectives RENAME COLUMN cameratargetcenterlat TO camera_target_center_lat;
ALTER TABLE perspectives RENAME COLUMN cameratargetcenterlng TO camera_target_center_lng;
ALTER TABLE perspectives RENAME COLUMN cameratargetpitch     TO camera_target_pitch;
ALTER TABLE perspectives RENAME COLUMN cameratargetzoom      TO camera_target_zoom;
ALTER TABLE perspectives DROP COLUMN IF EXISTS restored_at;
ALTER TABLE perspectives DROP COLUMN IF EXISTS transaction_id;

-- ── levels ────────────────────────────────────────────
ALTER TABLE levels RENAME COLUMN externalid  TO external_id;
ALTER TABLE levels RENAME COLUMN shortname   TO short_name;
ALTER TABLE levels RENAME COLUMN mapgroup_id TO map_group_id;
ALTER TABLE levels RENAME COLUMN publish     TO is_published;
ALTER TABLE levels DROP COLUMN IF EXISTS restored_at;
ALTER TABLE levels DROP COLUMN IF EXISTS transaction_id;
ALTER TABLE levels DROP COLUMN IF EXISTS type_id;  -- was FK to leveltype lookup; type is now a direct smallint

-- ── geo_references ─────────────────────────────────────
ALTER TABLE geo_references RENAME COLUMN controlx TO control_x;
ALTER TABLE geo_references RENAME COLUMN controly TO control_y;
ALTER TABLE geo_references RENAME COLUMN targetx  TO target_x;
ALTER TABLE geo_references RENAME COLUMN targety  TO target_y;
ALTER TABLE geo_references DROP COLUMN IF EXISTS restored_at;
ALTER TABLE geo_references DROP COLUMN IF EXISTS transaction_id;
```

- [ ] **Step 2: Run and verify**

```bash
psql $DATABASE_URL -f resources/migration_scripts/06_levels.sql
psql $DATABASE_URL -c "\d levels" | grep -E "map_group_id|is_published|external_id"
# Expected: all three columns present
```

- [ ] **Step 3: Commit**

```bash
git add resources/migration_scripts/06_levels.sql
git commit -m "chore: add db migration script 06 - transform level domain"
```

---

### Task 7: Transform location domain

**Files:**
- Create: `resources/migration_scripts/07_locations.sql`

- [ ] **Step 1: Create `07_locations.sql`**

```sql
-- ── location_categories ───────────────────────────────
ALTER TABLE location_categories RENAME COLUMN icondefault TO icon_default;
ALTER TABLE location_categories RENAME COLUMN externalid  TO external_id;
ALTER TABLE location_categories RENAME COLUMN sortindex   TO sort_index;
ALTER TABLE location_categories DROP COLUMN IF EXISTS restored_at;
ALTER TABLE location_categories DROP COLUMN IF EXISTS transaction_id;

-- ── locations ─────────────────────────────────────────
ALTER TABLE locations RENAME COLUMN externalid TO external_id;
-- Drop columns not in new schema
ALTER TABLE locations DROP COLUMN IF EXISTS restored_at;
ALTER TABLE locations DROP COLUMN IF EXISTS transaction_id;
ALTER TABLE locations DROP COLUMN IF EXISTS top_logo;
ALTER TABLE locations DROP COLUMN IF EXISTS top_logo_type;
ALTER TABLE locations DROP COLUMN IF EXISTS start_time;   -- event-style fields removed
ALTER TABLE locations DROP COLUMN IF EXISTS end_time;
ALTER TABLE locations DROP COLUMN IF EXISTS common_show_short_name;
ALTER TABLE locations DROP COLUMN IF EXISTS common_sub_type;

-- ── location_category_links ───────────────────────────
-- Old table had a surrogate bigserial id + (location_id, locationcategory_id)
-- New schema uses composite PK (location_id, category_id) — no surrogate id
ALTER TABLE location_category_links RENAME COLUMN locationcategory_id TO category_id;
ALTER TABLE location_category_links DROP COLUMN IF EXISTS id;
ALTER TABLE location_category_links ADD PRIMARY KEY (location_id, category_id);

-- ── location_images ───────────────────────────────────
ALTER TABLE location_images DROP COLUMN IF EXISTS restored_at;
ALTER TABLE location_images DROP COLUMN IF EXISTS transaction_id;

-- ── promotions ────────────────────────────────────────
ALTER TABLE promotions DROP COLUMN IF EXISTS restored_at;
ALTER TABLE promotions DROP COLUMN IF EXISTS transaction_id;
```

- [ ] **Step 2: Run and verify**

```bash
psql $DATABASE_URL -f resources/migration_scripts/07_locations.sql
psql $DATABASE_URL -c "SELECT COUNT(*) FROM location_category_links WHERE location_id IS NULL;"
# Expected: 0
```

- [ ] **Step 3: Commit**

```bash
git add resources/migration_scripts/07_locations.sql
git commit -m "chore: add db migration script 07 - transform location domain"
```

---

### Task 8: Transform product domain

**Files:**
- Create: `resources/migration_scripts/08_products.sql`

- [ ] **Step 1: Create `08_products.sql`**

```sql
-- ── product_categories ────────────────────────────────
ALTER TABLE product_categories DROP COLUMN IF EXISTS restored_at;
ALTER TABLE product_categories DROP COLUMN IF EXISTS transaction_id;

-- ── products ──────────────────────────────────────────
ALTER TABLE products DROP COLUMN IF EXISTS restored_at;
ALTER TABLE products DROP COLUMN IF EXISTS transaction_id;

-- ── product_category_links ────────────────────────────
-- Same pattern as location_category_links: drop surrogate id, make composite PK
ALTER TABLE product_category_links RENAME COLUMN productcategory_id TO category_id;
ALTER TABLE product_category_links DROP COLUMN IF EXISTS id;
ALTER TABLE product_category_links ADD PRIMARY KEY (product_id, category_id);

-- ── product_attachments ───────────────────────────────
ALTER TABLE product_attachments DROP COLUMN IF EXISTS restored_at;
ALTER TABLE product_attachments DROP COLUMN IF EXISTS transaction_id;
```

- [ ] **Step 2: Run and verify**

```bash
psql $DATABASE_URL -f resources/migration_scripts/08_products.sql
psql $DATABASE_URL -c "\d product_category_links"
# Expected: primary key on (product_id, category_id), no id column
```

- [ ] **Step 3: Commit**

```bash
git add resources/migration_scripts/08_products.sql
git commit -m "chore: add db migration script 08 - transform product domain"
```

---

### Task 9: Transform event domain

**Files:**
- Create: `resources/migration_scripts/09_events.sql`

- [ ] **Step 1: Create `09_events.sql`**

```sql
-- ── event_types ───────────────────────────────────────
ALTER TABLE event_types DROP COLUMN IF EXISTS restored_at;
ALTER TABLE event_types DROP COLUMN IF EXISTS transaction_id;

-- ── event_tags ────────────────────────────────────────
ALTER TABLE event_tags DROP COLUMN IF EXISTS restored_at;
ALTER TABLE event_tags DROP COLUMN IF EXISTS transaction_id;

-- ── events ────────────────────────────────────────────
ALTER TABLE events RENAME COLUMN bannerimage    TO banner_image;
ALTER TABLE events RENAME COLUMN starttime      TO start_time;
ALTER TABLE events RENAME COLUMN endtime        TO end_time;
ALTER TABLE events RENAME COLUMN contentdetail  TO content_detail;
ALTER TABLE events RENAME COLUMN contenturl     TO content_url;
ALTER TABLE events RENAME COLUMN iconimage      TO icon_image;
ALTER TABLE events RENAME COLUMN showendtime    TO show_end_time;
ALTER TABLE events RENAME COLUMN showstarttime  TO show_start_time;
ALTER TABLE events DROP COLUMN IF EXISTS restored_at;
ALTER TABLE events DROP COLUMN IF EXISTS transaction_id;

-- ── event_tag_links ───────────────────────────────────
ALTER TABLE event_tag_links RENAME COLUMN eventtag_id TO tag_id;
ALTER TABLE event_tag_links DROP COLUMN IF EXISTS id;
ALTER TABLE event_tag_links ADD PRIMARY KEY (event_id, tag_id);

-- ── event_location_links ──────────────────────────────
ALTER TABLE event_location_links DROP COLUMN IF EXISTS id;
ALTER TABLE event_location_links ADD PRIMARY KEY (event_id, location_id);

-- ── event_images ──────────────────────────────────────
ALTER TABLE event_images DROP COLUMN IF EXISTS restored_at;
ALTER TABLE event_images DROP COLUMN IF EXISTS transaction_id;
```

- [ ] **Step 2: Run and verify**

```bash
psql $DATABASE_URL -f resources/migration_scripts/09_events.sql
psql $DATABASE_URL -c "\d event_tag_links"
# Expected: primary key on (event_id, tag_id)
```

- [ ] **Step 3: Commit**

```bash
git add resources/migration_scripts/09_events.sql
git commit -m "chore: add db migration script 09 - transform event domain"
```

---

### Task 10: Transform survey domain

**Files:**
- Create: `resources/migration_scripts/10_surveys.sql`

- [ ] **Step 1: Create `10_surveys.sql`**

```sql
-- ── surveys ───────────────────────────────────────────
-- created_by_id was bigint FK to indoormap_api_appuser (not auth_user).
-- appuser will not become a users row. Set created_by to NULL for all rows
-- (no valid mapping exists).
ALTER TABLE surveys DROP COLUMN IF EXISTS created_by_id;
ALTER TABLE surveys ADD COLUMN IF NOT EXISTS created_by UUID;  -- nullable, no FK yet
ALTER TABLE surveys DROP COLUMN IF EXISTS restored_at;
ALTER TABLE surveys DROP COLUMN IF EXISTS transaction_id;

-- ── questions ─────────────────────────────────────────
ALTER TABLE questions DROP COLUMN IF EXISTS restored_at;
ALTER TABLE questions DROP COLUMN IF EXISTS transaction_id;
ALTER TABLE questions ALTER COLUMN question_number TYPE INTEGER USING question_number::integer;

-- ── options ───────────────────────────────────────────
ALTER TABLE options DROP COLUMN IF EXISTS restored_at;
ALTER TABLE options DROP COLUMN IF EXISTS transaction_id;
ALTER TABLE options ALTER COLUMN option_number TYPE INTEGER USING option_number::integer;

-- ── survey_responses ──────────────────────────────────
-- user_id was FK to indoormap_api_appuser — drop it (no mapping)
ALTER TABLE survey_responses DROP COLUMN IF EXISTS user_id;
ALTER TABLE survey_responses DROP COLUMN IF EXISTS restored_at;
ALTER TABLE survey_responses DROP COLUMN IF EXISTS transaction_id;

-- ── survey_answers ────────────────────────────────────
ALTER TABLE survey_answers DROP COLUMN IF EXISTS restored_at;
ALTER TABLE survey_answers DROP COLUMN IF EXISTS transaction_id;
```

- [ ] **Step 2: Run and verify**

```bash
psql $DATABASE_URL -f resources/migration_scripts/10_surveys.sql
psql $DATABASE_URL -c "SELECT COUNT(*) FROM surveys;"
# Expected: > 0 (rows preserved)
```

- [ ] **Step 3: Commit**

```bash
git add resources/migration_scripts/10_surveys.sql
git commit -m "chore: add db migration script 10 - transform survey domain"
```

---

### Task 11: Transform content domain

**Files:**
- Create: `resources/migration_scripts/11_content.sql`

- [ ] **Step 1: Create `11_content.sql`**

```sql
-- ── advertisements ────────────────────────────────────
ALTER TABLE advertisements RENAME COLUMN displayduration TO display_duration;
ALTER TABLE advertisements DROP COLUMN IF EXISTS article_id;
ALTER TABLE advertisements DROP COLUMN IF EXISTS restored_at;
ALTER TABLE advertisements DROP COLUMN IF EXISTS transaction_id;

-- ── articles ──────────────────────────────────────────
ALTER TABLE articles DROP COLUMN IF EXISTS product_id;
ALTER TABLE articles DROP COLUMN IF EXISTS product_category_id;
ALTER TABLE articles DROP COLUMN IF EXISTS application_language;
ALTER TABLE articles DROP COLUMN IF EXISTS company_name;
ALTER TABLE articles DROP COLUMN IF EXISTS email;
ALTER TABLE articles DROP COLUMN IF EXISTS full_name;
ALTER TABLE articles DROP COLUMN IF EXISTS landline;
ALTER TABLE articles DROP COLUMN IF EXISTS mobile;
ALTER TABLE articles DROP COLUMN IF EXISTS created_by;
ALTER TABLE articles DROP COLUMN IF EXISTS restored_at;
ALTER TABLE articles DROP COLUMN IF EXISTS transaction_id;

-- ── article_images ────────────────────────────────────
ALTER TABLE article_images RENAME COLUMN "order" TO sort_order;
ALTER TABLE article_images DROP COLUMN IF EXISTS restored_at;
ALTER TABLE article_images DROP COLUMN IF EXISTS transaction_id;

-- ── coupons ───────────────────────────────────────────
ALTER TABLE coupons DROP COLUMN IF EXISTS survey_id;
ALTER TABLE coupons DROP COLUMN IF EXISTS restored_at;
ALTER TABLE coupons DROP COLUMN IF EXISTS transaction_id;

-- ── beacons ───────────────────────────────────────────
ALTER TABLE beacons RENAME COLUMN hwid       TO hw_id;
ALTER TABLE beacons RENAME COLUMN vendorkey  TO vendor_key;
ALTER TABLE beacons RENAME COLUMN lotkey     TO lot_key;
ALTER TABLE beacons RENAME COLUMN uuid       TO uuid_val;
ALTER TABLE beacons RENAME COLUMN positionx  TO position_x;
ALTER TABLE beacons RENAME COLUMN positiony  TO position_y;
ALTER TABLE beacons RENAME COLUMN isenable   TO is_enable;
ALTER TABLE beacons RENAME COLUMN txpower    TO tx_power;
ALTER TABLE beacons DROP COLUMN IF EXISTS level;   -- redundant integer field
ALTER TABLE beacons DROP COLUMN IF EXISTS restored_at;
ALTER TABLE beacons DROP COLUMN IF EXISTS transaction_id;

-- ── connections ───────────────────────────────────────
ALTER TABLE connections RENAME COLUMN externalid TO external_id;
ALTER TABLE connections DROP COLUMN IF EXISTS restored_at;
ALTER TABLE connections DROP COLUMN IF EXISTS transaction_id;

-- ── connection_levels ─────────────────────────────────
ALTER TABLE connection_levels DROP COLUMN IF EXISTS restored_at;
ALTER TABLE connection_levels DROP COLUMN IF EXISTS transaction_id;

-- ── qrcodes ───────────────────────────────────────────
ALTER TABLE qrcodes DROP COLUMN IF EXISTS restored_at;
ALTER TABLE qrcodes DROP COLUMN IF EXISTS transaction_id;

-- ── reset_password_tokens ─────────────────────────────
ALTER TABLE reset_password_tokens DROP COLUMN IF EXISTS restored_at;
ALTER TABLE reset_password_tokens DROP COLUMN IF EXISTS transaction_id;
-- Map user_id (old integer) to new UUID
ALTER TABLE reset_password_tokens ADD COLUMN IF NOT EXISTS user_id_new UUID;
UPDATE reset_password_tokens rpt
SET user_id_new = m.new_id
FROM _user_id_map m
WHERE rpt.user_id = m.old_id;
ALTER TABLE reset_password_tokens DROP COLUMN IF EXISTS user_id;
ALTER TABLE reset_password_tokens RENAME COLUMN user_id_new TO user_id;
```

- [ ] **Step 2: Run and verify**

```bash
psql $DATABASE_URL -f resources/migration_scripts/11_content.sql
psql $DATABASE_URL -c "\d beacons" | grep -E "hw_id|uuid_val|is_enable"
# Expected: all three present
```

- [ ] **Step 3: Commit**

```bash
git add resources/migration_scripts/11_content.sql
git commit -m "chore: add db migration script 11 - transform content domain"
```

---

### Task 12: Transform misc tables — tags, search_queries, snapshots

**Files:**
- Create: `resources/migration_scripts/12_misc.sql`

- [ ] **Step 1: Create `12_misc.sql`**

```sql
-- ── tags ──────────────────────────────────────────────
ALTER TABLE tags DROP COLUMN IF EXISTS restored_at;
ALTER TABLE tags DROP COLUMN IF EXISTS transaction_id;

-- ── search_queries ────────────────────────────────────
ALTER TABLE search_queries DROP COLUMN IF EXISTS restored_at;
ALTER TABLE search_queries DROP COLUMN IF EXISTS transaction_id;
ALTER TABLE search_queries ALTER COLUMN search_count TYPE INTEGER USING search_count::integer;

-- ── snapshots ─────────────────────────────────────────
ALTER TABLE snapshots DROP COLUMN IF EXISTS restored_at;
ALTER TABLE snapshots DROP COLUMN IF EXISTS transaction_id;
-- state: varchar(6) → smallint (map known values)
ALTER TABLE snapshots ADD COLUMN IF NOT EXISTS state_new SMALLINT;
UPDATE snapshots SET state_new = CASE state
  WHEN 'draft'     THEN 0
  WHEN 'publish'   THEN 1
  WHEN 'published' THEN 1
  ELSE 0 END;
ALTER TABLE snapshots DROP COLUMN state;
ALTER TABLE snapshots RENAME COLUMN state_new TO state;
ALTER TABLE snapshots ALTER COLUMN state SET NOT NULL;
ALTER TABLE snapshots ALTER COLUMN state SET DEFAULT 0;
-- Rename publish_by_id → created_by and remap to UUID
ALTER TABLE snapshots ADD COLUMN IF NOT EXISTS created_by UUID;
UPDATE snapshots s
SET created_by = m.new_id
FROM _user_id_map m
WHERE s.publish_by_id = m.old_id;
ALTER TABLE snapshots DROP COLUMN IF EXISTS publish_by_id;
```

- [ ] **Step 2: Run and verify**

```bash
psql $DATABASE_URL -f resources/migration_scripts/12_misc.sql
psql $DATABASE_URL -c "\d snapshots" | grep -E "state|created_by"
# Expected: state smallint, created_by uuid
```

- [ ] **Step 3: Commit**

```bash
git add resources/migration_scripts/12_misc.sql
git commit -m "chore: add db migration script 12 - transform misc tables"
```

---

### Task 13: Restore FK constraints, install triggers, set uuidv7 defaults

**Files:**
- Create: `resources/migration_scripts/13_restore_constraints.sql`

- [ ] **Step 1: Create `13_restore_constraints.sql`**

```sql
-- ── Re-add FK constraints ─────────────────────────────
ALTER TABLE venues ADD CONSTRAINT venues_customer_id_fk
  FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL;

ALTER TABLE venue_user_roles ADD CONSTRAINT venue_user_roles_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
ALTER TABLE venue_user_roles ADD CONSTRAINT venue_user_roles_user_id_fk
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE venue_invitations ADD CONSTRAINT venue_invitations_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
ALTER TABLE venue_invitations ADD CONSTRAINT venue_invitations_invited_by_fk
  FOREIGN KEY (invited_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE map_groups ADD CONSTRAINT map_groups_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;

ALTER TABLE levels ADD CONSTRAINT levels_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
ALTER TABLE levels ADD CONSTRAINT levels_map_group_id_fk
  FOREIGN KEY (map_group_id) REFERENCES map_groups(id) ON DELETE SET NULL;
ALTER TABLE levels ADD CONSTRAINT levels_perspective_id_fk
  FOREIGN KEY (perspective_id) REFERENCES perspectives(id) ON DELETE SET NULL;

ALTER TABLE geo_references ADD CONSTRAINT geo_references_level_id_fk
  FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE CASCADE;

ALTER TABLE location_categories ADD CONSTRAINT location_categories_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;

ALTER TABLE locations ADD CONSTRAINT locations_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
ALTER TABLE locations ADD CONSTRAINT locations_level_id_fk
  FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE SET NULL;
ALTER TABLE locations ADD CONSTRAINT locations_main_category_id_fk
  FOREIGN KEY (main_category_id) REFERENCES location_categories(id) ON DELETE SET NULL;

ALTER TABLE location_category_links ADD CONSTRAINT location_category_links_location_id_fk
  FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE;
ALTER TABLE location_category_links ADD CONSTRAINT location_category_links_category_id_fk
  FOREIGN KEY (category_id) REFERENCES location_categories(id) ON DELETE CASCADE;

ALTER TABLE location_images ADD CONSTRAINT location_images_location_id_fk
  FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE;

ALTER TABLE product_categories ADD CONSTRAINT product_categories_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;

ALTER TABLE products ADD CONSTRAINT products_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
ALTER TABLE products ADD CONSTRAINT products_location_id_fk
  FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL;
ALTER TABLE products ADD CONSTRAINT products_main_category_id_fk
  FOREIGN KEY (main_category_id) REFERENCES product_categories(id) ON DELETE SET NULL;

ALTER TABLE product_category_links ADD CONSTRAINT product_category_links_product_id_fk
  FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
ALTER TABLE product_category_links ADD CONSTRAINT product_category_links_category_id_fk
  FOREIGN KEY (category_id) REFERENCES product_categories(id) ON DELETE CASCADE;

ALTER TABLE product_attachments ADD CONSTRAINT product_attachments_product_id_fk
  FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;

ALTER TABLE event_types ADD CONSTRAINT event_types_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;

ALTER TABLE events ADD CONSTRAINT events_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
ALTER TABLE events ADD CONSTRAINT events_type_id_fk
  FOREIGN KEY (type_id) REFERENCES event_types(id) ON DELETE SET NULL;

ALTER TABLE event_tag_links ADD CONSTRAINT event_tag_links_event_id_fk
  FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE;
ALTER TABLE event_tag_links ADD CONSTRAINT event_tag_links_tag_id_fk
  FOREIGN KEY (tag_id) REFERENCES event_tags(id) ON DELETE CASCADE;

ALTER TABLE event_location_links ADD CONSTRAINT event_location_links_event_id_fk
  FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE;
ALTER TABLE event_location_links ADD CONSTRAINT event_location_links_location_id_fk
  FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE;

ALTER TABLE event_images ADD CONSTRAINT event_images_event_id_fk
  FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE;

ALTER TABLE surveys ADD CONSTRAINT surveys_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
ALTER TABLE surveys ADD CONSTRAINT surveys_created_by_fk
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE questions ADD CONSTRAINT questions_survey_id_fk
  FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE CASCADE;

ALTER TABLE options ADD CONSTRAINT options_question_id_fk
  FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE;

ALTER TABLE survey_responses ADD CONSTRAINT survey_responses_survey_id_fk
  FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE CASCADE;

ALTER TABLE survey_answers ADD CONSTRAINT survey_answers_response_id_fk
  FOREIGN KEY (response_id) REFERENCES survey_responses(id) ON DELETE CASCADE;
ALTER TABLE survey_answers ADD CONSTRAINT survey_answers_question_id_fk
  FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE;
ALTER TABLE survey_answers ADD CONSTRAINT survey_answers_option_id_fk
  FOREIGN KEY (option_id) REFERENCES options(id) ON DELETE SET NULL;

ALTER TABLE advertisements ADD CONSTRAINT advertisements_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
ALTER TABLE advertisements ADD CONSTRAINT advertisements_location_id_fk
  FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE;

ALTER TABLE articles ADD CONSTRAINT articles_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
ALTER TABLE articles ADD CONSTRAINT articles_location_id_fk
  FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL;

ALTER TABLE article_images ADD CONSTRAINT article_images_article_id_fk
  FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE;

ALTER TABLE coupons ADD CONSTRAINT coupons_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;

ALTER TABLE beacons ADD CONSTRAINT beacons_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
ALTER TABLE beacons ADD CONSTRAINT beacons_level_id_fk
  FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE SET NULL;

ALTER TABLE connections ADD CONSTRAINT connections_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;

ALTER TABLE connection_levels ADD CONSTRAINT connection_levels_connection_id_fk
  FOREIGN KEY (connection_id) REFERENCES connections(id) ON DELETE CASCADE;
ALTER TABLE connection_levels ADD CONSTRAINT connection_levels_level_id_fk
  FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE CASCADE;

ALTER TABLE qrcodes ADD CONSTRAINT qrcodes_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
ALTER TABLE qrcodes ADD CONSTRAINT qrcodes_level_id_fk
  FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE SET NULL;
ALTER TABLE qrcodes ADD CONSTRAINT qrcodes_location_id_fk
  FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL;

ALTER TABLE snapshots ADD CONSTRAINT snapshots_venue_id_fk
  FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
ALTER TABLE snapshots ADD CONSTRAINT snapshots_created_by_fk
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE reset_password_tokens ADD CONSTRAINT reset_password_tokens_user_id_fk
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- ── Install updated_at triggers ───────────────────────
SELECT create_updated_at_trigger('customers');
SELECT create_updated_at_trigger('users');
SELECT create_updated_at_trigger('venues');
SELECT create_updated_at_trigger('venue_user_roles');
SELECT create_updated_at_trigger('venue_invitations');
SELECT create_updated_at_trigger('map_groups');
SELECT create_updated_at_trigger('perspectives');
SELECT create_updated_at_trigger('levels');
SELECT create_updated_at_trigger('geo_references');
SELECT create_updated_at_trigger('location_categories');
SELECT create_updated_at_trigger('locations');
SELECT create_updated_at_trigger('location_images');
SELECT create_updated_at_trigger('promotions');
SELECT create_updated_at_trigger('product_categories');
SELECT create_updated_at_trigger('products');
SELECT create_updated_at_trigger('product_attachments');
SELECT create_updated_at_trigger('event_types');
SELECT create_updated_at_trigger('event_tags');
SELECT create_updated_at_trigger('events');
SELECT create_updated_at_trigger('event_images');
SELECT create_updated_at_trigger('surveys');
SELECT create_updated_at_trigger('questions');
SELECT create_updated_at_trigger('options');
SELECT create_updated_at_trigger('survey_responses');
SELECT create_updated_at_trigger('survey_answers');
SELECT create_updated_at_trigger('beacons');
SELECT create_updated_at_trigger('connections');
SELECT create_updated_at_trigger('connection_levels');
SELECT create_updated_at_trigger('advertisements');
SELECT create_updated_at_trigger('articles');
SELECT create_updated_at_trigger('article_images');
SELECT create_updated_at_trigger('coupons');
SELECT create_updated_at_trigger('tags');
SELECT create_updated_at_trigger('search_queries');
SELECT create_updated_at_trigger('snapshots');
SELECT create_updated_at_trigger('qrcodes');

-- ── Set uuidv7() as default for all id columns ────────
-- (from migration 000023 — only run after uuidv7 function exists)
-- The uuidv7() function is defined in migration 000023.
-- We run that migration's ALTER TABLE statements here instead of via migrate.

-- Create uuidv7 function first (from migration 000023)
-- (paste the full uuidv7() function body from migrations/000023_uuidv7_defaults.up.sql here)
\i migrations/000023_uuidv7_defaults.up.sql

-- ── Drop temporary mapping table ──────────────────────
DROP TABLE IF EXISTS _user_id_map;
```

- [ ] **Step 2: Run and verify**

```bash
psql $DATABASE_URL -f resources/migration_scripts/13_restore_constraints.sql

# Verify FK count matches expectation
psql $DATABASE_URL -c "SELECT COUNT(*) FROM information_schema.table_constraints
  WHERE constraint_type = 'FOREIGN KEY' AND table_schema = 'public';"
# Expected: ~40+

# Verify triggers installed
psql $DATABASE_URL -c "SELECT COUNT(*) FROM information_schema.triggers
  WHERE trigger_name = 'set_updated_at';"
# Expected: matches number of SELECT create_updated_at_trigger calls above
```

- [ ] **Step 3: Commit**

```bash
git add resources/migration_scripts/13_restore_constraints.sql
git commit -m "chore: add db migration script 13 - restore FK constraints and install triggers"
```

---

### Task 14: Create app_users migration and migrate appuser + visitor data

**Files:**
- Create: `migrations/000025_create_app_users.up.sql`
- Create: `migrations/000025_create_app_users.down.sql`
- Create: `resources/migration_scripts/14_app_users.sql`

The merged table unifies staff app-users (from `indoormap_api_appuser`) and anonymous event visitors (from `indoormap_api_visitor`). A `source` discriminator column (`'staff'` | `'visitor'`) distinguishes them.

- [ ] **Step 1: Create `migrations/000025_create_app_users.up.sql`**

```sql
CREATE TABLE app_users (
    id                  UUID PRIMARY KEY DEFAULT uuidv7(),
    venue_id            UUID REFERENCES venues(id) ON DELETE SET NULL,
    external_id         TEXT,
    source              TEXT NOT NULL,           -- 'staff' | 'visitor'
    type                SMALLINT NOT NULL DEFAULT 0,

    -- Identity
    first_name          TEXT,
    last_name           TEXT,
    first_name_en       TEXT,
    last_name_en        TEXT,
    email               TEXT,
    phone               TEXT,

    -- Staff-specific
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

    -- Visitor-specific
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

CREATE INDEX app_users_venue_id_idx    ON app_users(venue_id);
CREATE INDEX app_users_external_id_idx ON app_users(external_id) WHERE external_id IS NOT NULL;
CREATE INDEX app_users_email_idx       ON app_users(email)       WHERE email IS NOT NULL;
SELECT create_updated_at_trigger('app_users');
```

- [ ] **Step 2: Create `migrations/000025_create_app_users.down.sql`**

```sql
DROP TABLE IF EXISTS app_users;
```

- [ ] **Step 3: Create `resources/migration_scripts/14_app_users.sql`** to migrate the data

```sql
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
    gen_random_uuid(),
    NULL,                   -- appuser has no direct venue_id
    external_id,
    'staff',
    COALESCE(type, 0),
    first_name, last_name, first_name_en, last_name_en,
    mail, tel,
    company_name, company_name_en, department, position, position_en,
    COALESCE(section, 0),
    app, token, token_expiration,
    COALESCE(survey_flag, 0),
    COALESCE(staff_lead_flag, 0),
    COALESCE(created_at, NOW()),
    COALESCE(updated_at, NOW()),
    deleted_at
FROM indoormap_api_appuser;

-- Migrate visitors from indoormap_api_visitor
-- Note: phone_number was stored as bytea (encrypted) — cast to text; will be garbled
-- if encrypted at rest. Accept this — the original encryption key is unavailable.
INSERT INTO app_users (
    id, venue_id, source, type,
    first_name, email, phone,
    visitor_type, interests, other_interests,
    is_consented, ip_address, user_agent, business_name,
    created_at, updated_at, deleted_at
)
SELECT
    -- reformat char(32) UUID to proper UUID
    (substring(id,1,8)||'-'||substring(id,9,4)||'-'||substring(id,13,4)||
     '-'||substring(id,17,4)||'-'||substring(id,21,12))::uuid,
    CASE WHEN venue_id IS NOT NULL
         THEN (substring(venue_id,1,8)||'-'||substring(venue_id,9,4)||'-'||
               substring(venue_id,13,4)||'-'||substring(venue_id,17,4)||'-'||
               substring(venue_id,21,12))::uuid
         ELSE NULL END,
    'visitor',
    COALESCE(visitor_type, 0),
    full_name, email,
    encode(phone_number, 'escape'),   -- bytea → text best-effort
    visitor_type,
    interests::jsonb,
    other_interests,
    is_consented,
    ip_address, user_agent, business_name,
    COALESCE(created_at, NOW()),
    COALESCE(updated_at, NOW()),
    deleted_at
FROM indoormap_api_visitor;

-- Drop source tables now that data is migrated
DROP TABLE IF EXISTS indoormap_api_appuser;
DROP TABLE IF EXISTS indoormap_api_visitor;
```

- [ ] **Step 4: Run migration 000025 first, then data migration**

```bash
make migrate-up
# (runs 000025_create_app_users — creates the table)

psql $DATABASE_URL -f resources/migration_scripts/14_app_users.sql
```

- [ ] **Step 5: Verify**

```bash
psql $DATABASE_URL -c "
  SELECT source, COUNT(*) FROM app_users GROUP BY source;
"
# Expected: rows for 'staff' and 'visitor' matching original table counts

psql $DATABASE_URL -c "
  SELECT COUNT(*) FROM information_schema.tables
  WHERE table_name IN ('indoormap_api_appuser','indoormap_api_visitor');
"
# Expected: 0
```

- [ ] **Step 6: Commit**

```bash
git add migrations/000025_create_app_users.up.sql migrations/000025_create_app_users.down.sql \
        resources/migration_scripts/14_app_users.sql
git commit -m "feat: add app_users table (merged staff + visitor) and migration script"
```

---

### Task 15: Mark Go migrations as applied and verify

**Files:**
- Create: `resources/migration_scripts/15_mark_migrations_applied.sql`

- [ ] **Step 1: Create `15_mark_migrations_applied.sql`**

golang-migrate tracks applied migrations in a `schema_migrations` table. Populate it so `make migrate-up` only runs future migrations (000026+).

```sql
CREATE TABLE IF NOT EXISTS schema_migrations (
    version  BIGINT  NOT NULL PRIMARY KEY,
    dirty    BOOLEAN NOT NULL
);

INSERT INTO schema_migrations (version, dirty) VALUES
  (1, false),
  (2, false),
  (3, false),
  (4, false),
  (5, false),
  (6, false),
  (7, false),
  (8, false),
  (9, false),
  (10, false),
  (11, false),
  (12, false),
  (13, false),
  (14, false),
  (15, false),
  (16, false),
  (17, false),
  (18, false),
  (19, false),
  (20, false),
  (21, false),
  (22, false),
  (23, false),
  (24, false),
  (25, false)
ON CONFLICT (version) DO UPDATE SET dirty = false;
```

- [ ] **Step 2: Run**

```bash
psql $DATABASE_URL -f resources/migration_scripts/15_mark_migrations_applied.sql
```

- [ ] **Step 3: Verify golang-migrate sees no pending migrations**

```bash
make migrate-up
# Expected: "no change" — no migrations to apply
```

- [ ] **Step 4: Run a full sanity check**

```bash
# Go backend should compile and start without errors
make serve

# Check key table shapes match what the Go repo layer expects
psql $DATABASE_URL -c "\d venues" | grep -E "slug|large_logo|is_published|lat|lng"
psql $DATABASE_URL -c "\d users"  | grep -E "password_hash|is_system_admin|deleted_at"
psql $DATABASE_URL -c "\d app_users" | grep -E "source|venue_id|is_consented"

# Row counts — confirm data not lost
psql $DATABASE_URL -c "
  SELECT 'venues'     , COUNT(*) FROM venues      UNION ALL
  SELECT 'locations'  , COUNT(*) FROM locations    UNION ALL
  SELECT 'products'   , COUNT(*) FROM products     UNION ALL
  SELECT 'events'     , COUNT(*) FROM events       UNION ALL
  SELECT 'surveys'    , COUNT(*) FROM surveys      UNION ALL
  SELECT 'app_users'  , COUNT(*) FROM app_users;
"
```

- [ ] **Step 5: Commit**

```bash
git add resources/migration_scripts/15_mark_migrations_applied.sql
git commit -m "chore: add db migration script 15 - mark Go migrations 1-25 as applied"
```

---

## Execution Order Summary

```
psql $DATABASE_URL -f resources/migration_scripts/01_prep.sql
psql $DATABASE_URL -f resources/migration_scripts/02_drop_fk_constraints.sql
psql $DATABASE_URL -f resources/migration_scripts/03_rename_tables.sql
psql $DATABASE_URL -f resources/migration_scripts/04_core_tables.sql
psql $DATABASE_URL -f resources/migration_scripts/05_venues.sql
psql $DATABASE_URL -f resources/migration_scripts/06_levels.sql
psql $DATABASE_URL -f resources/migration_scripts/07_locations.sql
psql $DATABASE_URL -f resources/migration_scripts/08_products.sql
psql $DATABASE_URL -f resources/migration_scripts/09_events.sql
psql $DATABASE_URL -f resources/migration_scripts/10_surveys.sql
psql $DATABASE_URL -f resources/migration_scripts/11_content.sql
psql $DATABASE_URL -f resources/migration_scripts/12_misc.sql
psql $DATABASE_URL -f resources/migration_scripts/13_restore_constraints.sql
make migrate-up    # runs 000025 only (creates app_users table)
psql $DATABASE_URL -f resources/migration_scripts/14_app_users.sql
psql $DATABASE_URL -f resources/migration_scripts/15_mark_migrations_applied.sql
```

> **Important:** Run each script in a transaction (`psql -1`) so a failure rolls back cleanly. Test against a copy of the database first.
