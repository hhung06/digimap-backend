-- Script 09: Transform event_types, event_tags, events, event_tag_links, event_location_links, event_images

-- ── event_types ──────────────────────────────────────────────────────────────
-- Actual columns: name, venue_id, localization, restored_at, transaction_id
ALTER TABLE event_types DROP COLUMN IF EXISTS restored_at;
ALTER TABLE event_types DROP COLUMN IF EXISTS transaction_id;

-- ── event_tags ───────────────────────────────────────────────────────────────
-- Actual columns: name, localization, restored_at, transaction_id
ALTER TABLE event_tags DROP COLUMN IF EXISTS restored_at;
ALTER TABLE event_tags DROP COLUMN IF EXISTS transaction_id;

-- ── events ───────────────────────────────────────────────────────────────────
-- Actual columns: bannerimage, title, starttime, endtime, type_id, contentdetail,
--                 contenturl, description, iconimage, showendtime, showstarttime,
--                 venue_id, localization, restored_at, transaction_id
ALTER TABLE events RENAME COLUMN bannerimage   TO banner_image;
ALTER TABLE events RENAME COLUMN starttime     TO start_time;
ALTER TABLE events RENAME COLUMN endtime       TO end_time;
ALTER TABLE events RENAME COLUMN contentdetail TO content_detail;
ALTER TABLE events RENAME COLUMN contenturl    TO content_url;
ALTER TABLE events RENAME COLUMN iconimage     TO icon_image;
ALTER TABLE events RENAME COLUMN showendtime   TO show_end_time;
ALTER TABLE events RENAME COLUMN showstarttime TO show_start_time;
ALTER TABLE events DROP COLUMN IF EXISTS restored_at;
ALTER TABLE events DROP COLUMN IF EXISTS transaction_id;

-- ── event_tag_links ──────────────────────────────────────────────────────────
-- Actual columns: id (bigserial PK), event_id, eventtag_id
ALTER TABLE event_tag_links RENAME COLUMN eventtag_id TO tag_id;
ALTER TABLE event_tag_links DROP COLUMN IF EXISTS id;
ALTER TABLE event_tag_links ADD PRIMARY KEY (event_id, tag_id);

-- ── event_location_links ─────────────────────────────────────────────────────
-- Actual columns: id (bigserial PK), event_id, location_id
ALTER TABLE event_location_links DROP COLUMN IF EXISTS id;
ALTER TABLE event_location_links ADD PRIMARY KEY (event_id, location_id);

-- ── event_images ─────────────────────────────────────────────────────────────
-- Actual columns: image, event_id, restored_at, transaction_id
ALTER TABLE event_images DROP COLUMN IF EXISTS restored_at;
ALTER TABLE event_images DROP COLUMN IF EXISTS transaction_id;
